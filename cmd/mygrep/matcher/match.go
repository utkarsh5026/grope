package matcher

import (
	"strings"
	"unicode"
)

type Literal byte

const (
	StartsWith   = '^'
	Backslash    = '\\'
	Digit        = 'd'
	AlphaNumeric = 'w'
	LeftBracket  = '['
	RightBracket = ']'
	NotInClass   = '^'
)

func Match(line []byte, pattern string) bool {
	if len(pattern) == 0 {
		return true
	}

	if pattern[0] == StartsWith {
		return matchStart(line, pattern[1:])
	}

	for i := range line {
		if matchStart(line[i:], pattern) {
			return true
		}
	}

	return false
}

func matchStart(line []byte, pattern string) bool {
	lineIdx := 0 // line index
	patIdx := 0  // pattern index

	for patIdx < len(pattern) {
		if lineIdx >= len(line) {
			return false
		}

		switch {
		case pattern[patIdx] == Backslash:
			match := matchEscapeSequence(line[lineIdx], pattern[patIdx+1])
			if !match {
				return false
			}

			lineIdx++
			patIdx += 2

		case pattern[patIdx] == LeftBracket:
			endIdx := strings.IndexByte(pattern[patIdx:], RightBracket)
			if endIdx == -1 {
				return false
			}

			chars := pattern[patIdx+1 : patIdx+endIdx]
			match := matchCharacterClass(line[lineIdx], chars)

			if !match {
				return false
			}
			lineIdx++
			patIdx += endIdx + 1

		case pattern[patIdx] == line[lineIdx]:
			lineIdx++
			patIdx++

		default:
			return false
		}
	}

	return true
}

// matchEscapeSequence checks if a given character matches an escape sequence.
//
// Parameters:
// - char: The character to be matched.
// - escapeChar: The escape character that defines the type of match.
//
// Returns:
// - bool: True if the character matches the escape sequence, false otherwise.
func matchEscapeSequence(char byte, escapeChar byte) bool {
	switch escapeChar {
	case Digit:
		return unicode.IsDigit(rune(char))
	case AlphaNumeric:
		return unicode.IsLetter(rune(char)) || unicode.IsDigit(rune(char)) || char == '_'
	default:
		return char == escapeChar
	}
}

// matchCharacterClass checks if a given character matches a character class.
//
// Parameters:
// - char: The character to be matched.
// - class: The character class, which can optionally start with '^' to indicate negation.
//
// Returns:
// - bool: True if the character matches the character class, false otherwise.
func matchCharacterClass(char byte, class string) bool {
	if class[0] == NotInClass {
		return !strings.ContainsRune(class[1:], rune(char))
	}
	return strings.ContainsRune(class, rune(char))
}
