package matcher

import (
	"strings"
	"unicode"
)

type Literal byte

const (
	StartsWith   = '^'
	EndsWith     = '$'
	Backslash    = '\\'
	Digit        = 'd'
	AlphaNumeric = 'w'
	LeftBracket  = '['
	RightBracket = ']'
	NotInClass   = '^'
	OneOrMore    = '+'
	ZeroOrOne    = '?'
	ZeroOrMore   = '*'
	AnyCharacter = '.'
	OrCharacter  = '|'
	LeftParen    = '('
	RightParen   = ')'
)

// Match checks if a given line matches a pattern.
//
// Parameters:
// - line: The byte slice representing the line to be checked.
// - pattern: The pattern string to be matched.
//
// Returns:
// - bool: True if the line matches the pattern, false otherwise.
func Match(line []byte, pattern string) bool {
	idx := MatchWithIdx(line, pattern)
	return idx != -1
}

// MatchWithIdx returns the index of the first character in the line that matches the pattern.
//
// Parameters:
// - line: The byte slice representing the line to be checked.
// - pattern: The pattern string to be matched.
//
// Returns:
// - int: The index of the first character in the line that matches the pattern, or -1 if no match is found.
func MatchWithIdx(line []byte, pattern string) int {
	if len(pattern) == 0 {
		return 0
	}

	if pattern[0] == StartsWith {
		if matchStart(line, pattern[1:]) {
			return 0
		}
		return -1
	}

	for i := range line {
		if matchStart(line[i:], pattern) {
			return i
		}
	}

	return -1
}

// matchStart checks if a given line matches a pattern from the start.
//
// Parameters:
// - line: The byte slice representing the line to be checked.
// - pattern: The pattern string to be matched.
//
// Returns:
// - bool: True if the line matches the pattern from the start, false otherwise.
func matchStart(line []byte, pattern string) bool {
	lineIdx := 0 // line index
	patIdx := 0  // pattern index

	for patIdx < len(pattern) {

		if lineIdx >= len(line) {
			return patIdx == len(pattern) || (patIdx == len(pattern)-1 && pattern[patIdx] == EndsWith)
		}

		switch {
		case pattern[patIdx] == Backslash:
			if patIdx+1 >= len(pattern) {
				return false
			}
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

		case pattern[patIdx] == LeftParen:
			endIdx := strings.IndexByte(pattern[patIdx:], RightParen)
			subPattern := pattern[patIdx+1 : patIdx+endIdx]

			if strings.Contains(subPattern, string(OrCharacter)) {
				alternatives := strings.Split(subPattern, string(OrCharacter))
				for _, alt := range alternatives {
					if matchStart(line[lineIdx:], alt) {
						return true
					}
				}
				return false
			}

		case patIdx+1 < len(pattern) && isQuantifier(pattern[patIdx+1]):
			quantifier := pattern[patIdx+1]
			count, ok := handleQuantifier(line[lineIdx:], pattern[patIdx], quantifier)
			if !ok {
				return false
			}
			lineIdx += count
			patIdx += 2

		case pattern[patIdx] == line[lineIdx] || pattern[patIdx] == AnyCharacter:
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

// matchRepetition counts the number of consecutive occurrences of a character in a line.
//
// Parameters:
// - line: The byte slice representing the line to be checked.
// - char: The character to be matched.
//
// Returns:
// - int: The count of consecutive occurrences of the character at the start of the line.
func matchRepetition(line []byte, char byte) int {
	count := 0
	for _, c := range line {
		if c != char {
			break
		}
		count++
	}
	return count
}

// isQuantifier checks if a given character is a quantifier.
//
// Parameters:
// - char: The character to be checked.
//
// Returns:
// - bool: True if the character is a quantifier (one of '+', '?', or '*'), false otherwise.
func isQuantifier(char byte) bool {
	return char == OneOrMore || char == ZeroOrOne || char == ZeroOrMore
}

// handleQuantifier processes a quantifier for a given character in a line.
//
// Parameters:
// - line: The byte slice representing the line to be checked.
// - char: The character to be matched.
// - quantifier: The quantifier character ('*', '+', or '?').
//
// Returns:
// - int: The count of consecutive occurrences of the character in the line.
// - bool: True if the quantifier condition is satisfied, false otherwise.
func handleQuantifier(line []byte, char byte, quantifier byte) (int, bool) {
	count := matchRepetition(line, char)

	switch quantifier {
	case ZeroOrMore:
		return count, true
	case OneOrMore:
		return count, count > 0
	case ZeroOrOne:
		if count > 1 {
			count = 1
		}
		return count, true
	default:
		return 0, false
	}
}
