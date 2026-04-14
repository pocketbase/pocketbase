package inflector

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var columnifyRemoveRegex = regexp.MustCompile(`[^\w\.\*\-\_\@\#]+`)
var snakecaseSplitRegex = regexp.MustCompile(`[\W_]+`)

// UcFirst converts the first character of a string into uppercase.
func UcFirst(str string) string {
	if str == "" {
		return ""
	}

	r, size := utf8.DecodeRuneInString(str)
	if r == utf8.RuneError {
		return str
	}

	upper := unicode.ToUpper(r)
	if upper == r {
		return str // already uppercase, no allocation needed
	}

	return string(upper) + str[size:]
}

// isColumnifyAllowed reports whether c is a valid db identifier character:
// word characters (\w), '.', '*', '-', '_', '@', '#'.
func isColumnifyAllowed(c rune) bool {
	if c == '.' || c == '*' || c == '-' || c == '_' || c == '@' || c == '#' {
		return true
	}
	// \w equivalent: letters, digits, underscore (underscore already handled above)
	return unicode.IsLetter(c) || unicode.IsDigit(c)
}

// Columnify strips invalid db identifier characters.
func Columnify(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, c := range str {
		if isColumnifyAllowed(c) {
			b.WriteRune(c)
		}
	}
	if b.Len() == len(str) {
		return str // no characters were stripped, avoid allocation
	}
	return b.String()
}

// Sentenize converts and normalizes string into a sentence.
func Sentenize(str string) string {
	str = strings.TrimSpace(str)
	if str == "" {
		return ""
	}

	str = UcFirst(str)

	lastChar := str[len(str)-1:]
	if lastChar != "." && lastChar != "?" && lastChar != "!" {
		return str + "."
	}

	return str
}

// Sanitize sanitizes `str` by removing all characters satisfying `removePattern`.
// Returns an error if the pattern is not valid regex string.
func Sanitize(str string, removePattern string) (string, error) {
	exp, err := regexp.Compile(removePattern)
	if err != nil {
		return "", err
	}

	return exp.ReplaceAllString(str, ""), nil
}

// isWordChar reports whether c is a word character (letter or digit).
// Underscore is intentionally excluded as it is treated as a separator.
func isWordChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c)
}

// Snakecase removes all non word characters and converts any english text into a snakecase.
// "ABBREVIATIONS" are preserved, eg. "myTestDB" will become "my_test_db".
func Snakecase(str string) string {
	var result strings.Builder
	result.Grow(len(str) + 4) // slight over-estimate for underscores

	var prevWasWord bool
	var prevWasUpper bool
	var prevRune rune

	for _, c := range str {
		if !isWordChar(c) {
			// non-word character (includes underscore) acts as separator
			if prevWasWord {
				prevWasWord = false
			}
			prevWasUpper = false
			prevRune = c
			continue
		}

		isUpper := unicode.IsUpper(c)

		if !prevWasWord {
			// first word char after a separator or start
			if result.Len() > 0 {
				result.WriteByte('_')
			}
		} else if isUpper && !prevWasUpper {
			// camelCase boundary: lowercase->uppercase
			result.WriteByte('_')
		}

		result.WriteRune(unicode.ToLower(c))
		prevWasWord = true
		prevWasUpper = isUpper
		prevRune = c
	}

	_ = prevRune // suppress unused variable warning

	return result.String()
}

// Camelize converts the provided string to its "CamelCased" version
// (non alphanumeric characters are removed).
//
// For example:
//
//	inflector.Camelize("send_email") // "SendEmail"
func Camelize(str string) string {
	var result strings.Builder

	var isPrevSpecial bool

	for _, c := range str {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			isPrevSpecial = true
			continue
		}

		if isPrevSpecial || result.Len() == 0 {
			isPrevSpecial = false
			result.WriteRune(unicode.ToUpper(c))
		} else {
			result.WriteRune(c)
		}
	}

	return result.String()
}
