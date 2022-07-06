package inflector

import (
	"regexp"
	"strings"
	"unicode"
)

var columnifyRemoveRegex = regexp.MustCompile(`[^\w\.\*\-\_\@\#]+`)
var snakecaseSplitRegex = regexp.MustCompile(`[\W_]+`)
var usernamifySplitRegex = regexp.MustCompile(`\W+`)

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

	s := []rune(str)
	sentence := string(unicode.ToUpper(s[0])) + string(s[1:])

	lastChar := string(s[len(s)-1:])
	if lastChar != "." && lastChar != "?" && lastChar != "!" {
		return sentence + "."
	}

	return sentence
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

// Usernamify generates a properly formatted username from the provided string.
// Returns "unknown" if `str` is empty or contains only non word characters.
//
// ```go
// Usernamify("John Doe,   hello") // "john.doe.hello"
// ```
func Usernamify(str string) string {
	// split at any non word character
	words := usernamifySplitRegex.Split(strings.ToLower(str), -1)

	// concatenate any non empty word with a dot
	var result strings.Builder
	for _, word := range words {
		if word == "" {
			continue
		}

		if result.Len() > 0 {
			result.WriteString(".")
		}

		result.WriteString(word)
	}

	if result.Len() == 0 {
		return "unknown"
	}

	return result.String()
}
