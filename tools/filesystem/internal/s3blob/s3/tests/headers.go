package tests

import (
	"net/http"
	"regexp"
	"strings"
)

// ExpectHeaders checks whether specified headers match the expectations.
// The expectations map entry key is the header name.
// The expectations map entry value is the first header value. If wrapped with `^...$`
// it is compared as regular expression.
func ExpectHeaders(headers http.Header, expectations map[string]string) bool {
	for h, expected := range expectations {
		v := headers.Get(h)

		pattern := expected
		if !strings.HasPrefix(pattern, "^") && !strings.HasSuffix(pattern, "$") {
			pattern = "^" + regexp.QuoteMeta(pattern) + "$"
		}

		expectedRegex, err := regexp.Compile(pattern)
		if err != nil {
			return false
		}

		if !expectedRegex.MatchString(v) {
			return false
		}
	}

	return true
}
