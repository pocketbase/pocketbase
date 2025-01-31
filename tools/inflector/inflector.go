package inflector

import (
	"regexp"
	"strings"
	"unicode"
)

var columnifyRemoveRegex = regexp.MustCompile(`[^\w\.\*\-\_\@\#]+`)
var snakecaseSplitRegex = regexp.MustCompile(`[\W_]+`)

// UcFirst converts the first character of a string into uppercase.
func UcFirst(str string) string {
	if str == "" {
		return ""
	}

	s := []rune(str)

	return string(unicode.ToUpper(s[0])) + string(s[1:])
}

// Columnify strips invalid db identifier characters.
func Columnify(str string) string {
	return columnifyRemoveRegex.ReplaceAllString(str, "")
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

// Snakecase removes all non word characters and converts any english text into a snakecase.
// "ABBREVIATIONS" are preserved, eg. "myTestDB" will become "my_test_db".
func Snakecase(str string) string {
	var result strings.Builder

	// split at any non word character and underscore
	words := snakecaseSplitRegex.Split(str, -1)

	for _, word := range words {
		if word == "" {
			continue
		}

		if result.Len() > 0 {
			result.WriteString("_")
		}

		for i, c := range word {
			if unicode.IsUpper(c) && i > 0 &&
				// is not a following uppercase character
				!unicode.IsUpper(rune(word[i-1])) {
				result.WriteString("_")
			}

			result.WriteRune(c)
		}
	}

	return strings.ToLower(result.String())
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
