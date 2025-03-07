package blob

import (
	"fmt"
	"strconv"
)

// Copied from gocloud.dev/blob to avoid nuances around the specific
// HEX escaping/unescaping rules.
//
// -------------------------------------------------------------------
// Copyright 2019 The Go Cloud Development Kit Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// -------------------------------------------------------------------

// HexEscape returns s, with all runes for which shouldEscape returns true
// escaped to "__0xXXX__", where XXX is the hex representation of the rune
// value. For example, " " would escape to "__0x20__".
//
// Non-UTF-8 strings will have their non-UTF-8 characters escaped to
// unicode.ReplacementChar; the original value is lost. Please file an
// issue if you need non-UTF8 support.
//
// Note: shouldEscape takes the whole string as a slice of runes and an
// index. Passing it a single byte or a single rune doesn't provide
// enough context for some escape decisions; for example, the caller might
// want to escape the second "/" in "//" but not the first one.
// We pass a slice of runes instead of the string or a slice of bytes
// because some decisions will be made on a rune basis (e.g., encode
// all non-ASCII runes).
func HexEscape(s string, shouldEscape func(s []rune, i int) bool) string {
	// Do a first pass to see which runes (if any) need escaping.
	runes := []rune(s)
	var toEscape []int
	for i := range runes {
		if shouldEscape(runes, i) {
			toEscape = append(toEscape, i)
		}
	}
	if len(toEscape) == 0 {
		return s
	}

	// Each escaped rune turns into at most 14 runes ("__0x7fffffff__"),
	// so allocate an extra 13 for each. We'll reslice at the end
	// if we didn't end up using them.
	escaped := make([]rune, len(runes)+13*len(toEscape))
	n := 0 // current index into toEscape
	j := 0 // current index into escaped
	for i, r := range runes {
		if n < len(toEscape) && i == toEscape[n] {
			// We were asked to escape this rune.
			for _, x := range fmt.Sprintf("__%#x__", r) {
				escaped[j] = x
				j++
			}
			n++
		} else {
			escaped[j] = r
			j++
		}
	}

	return string(escaped[0:j])
}

// unescape tries to unescape starting at r[i].
// It returns a boolean indicating whether the unescaping was successful,
// and (if true) the unescaped rune and the last index of r that was used
// during unescaping.
func unescape(r []rune, i int) (bool, rune, int) {
	// Look for "__0x".
	if r[i] != '_' {
		return false, 0, 0
	}
	i++
	if i >= len(r) || r[i] != '_' {
		return false, 0, 0
	}
	i++
	if i >= len(r) || r[i] != '0' {
		return false, 0, 0
	}
	i++
	if i >= len(r) || r[i] != 'x' {
		return false, 0, 0
	}
	i++

	// Capture the digits until the next "_" (if any).
	var hexdigits []rune
	for ; i < len(r) && r[i] != '_'; i++ {
		hexdigits = append(hexdigits, r[i])
	}

	// Look for the trailing "__".
	if i >= len(r) || r[i] != '_' {
		return false, 0, 0
	}
	i++
	if i >= len(r) || r[i] != '_' {
		return false, 0, 0
	}

	// Parse the hex digits into an int32.
	retval, err := strconv.ParseInt(string(hexdigits), 16, 32)
	if err != nil {
		return false, 0, 0
	}

	return true, rune(retval), i
}

// HexUnescape reverses HexEscape.
func HexUnescape(s string) string {
	var unescaped []rune

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if ok, newR, newI := unescape(runes, i); ok {
			// We unescaped some runes starting at i, resulting in the
			// unescaped rune newR. The last rune used was newI.
			if unescaped == nil {
				// This is the first rune we've encountered that
				// needed unescaping. Allocate a buffer and copy any
				// previous runes.
				unescaped = make([]rune, i)
				copy(unescaped, runes)
			}
			unescaped = append(unescaped, newR)
			i = newI
		} else if unescaped != nil {
			unescaped = append(unescaped, runes[i])
		}
	}

	if unescaped == nil {
		return s
	}

	return string(unescaped)
}
