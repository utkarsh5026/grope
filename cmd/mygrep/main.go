package main

import (
	"strings"

	// Uncomment this to pass the first stage
	// "bytes"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	var ok bool
	ok = ContainsPattern(line, pattern)
	return ok, nil
}

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		_, _ = fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}
}

func contains(start int, target string, pattern string) bool {
	ti := start // target index
	i := 0      // pattern index
	for i < len(pattern) {
		if ti >= len(target) {
			return false
		}

		if pattern[i] == target[ti] {
			ti++
			i++
			continue
		}

		curr := pattern[i]
		if curr == '\\' && i < len(pattern)-1 {
			next := pattern[i+1]
			if next == 'd' && !isDigit(target[ti]) {
				return false
			}

			if next == 'w' && !isAlphaNumeric(target[ti]) {
				return false
			}

			ti++
			i += 2
		} else if curr == '[' {
			end := strings.Index(pattern[i:], "]")
			if end == -1 {
				return false
			}

			chars := pattern[i+1 : i+end]
			if strings.HasPrefix(chars, "^") {
				if strings.ContainsRune(chars[1:], rune(target[ti])) {
					return false
				}
			} else if !strings.ContainsRune(chars, rune(target[ti])) {
				return false
			}
			ti++
			i += end + 1
		} else {
			return false
		}
	}

	return true
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isAlphaNumeric(c byte) bool {
	return isDigit(c) || isLetter(c) || c == '_'
}

func ContainsPattern(line []byte, pattern string) bool {
	for i := 0; i < len(line); i++ {
		if contains(i, string(line), pattern) {
			return true
		}
	}
	return false
}
