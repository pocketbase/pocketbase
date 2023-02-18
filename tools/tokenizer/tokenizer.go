// Package tokenizer implements a rudimentary tokens parser of buffered
// io.Reader while respecting quotes and parenthesis boundaries.
//
// Example
//
//	tk := tokenizer.NewFromString("a, b, (c, d)")
//	result, _ := tk.ScanAll() // ["a", "b", "(c, d)"]
package tokenizer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// eof represents a marker rune for the end of the reader.
const eof = rune(0)

// DefaultSeparators is a list with the default token separator characters.
var DefaultSeparators = []rune{','}

// NewFromString creates new Tokenizer from the provided string.
func NewFromString(str string) *Tokenizer {
	return New(strings.NewReader(str))
}

// NewFromBytes creates new Tokenizer from the provided bytes slice.
func NewFromBytes(b []byte) *Tokenizer {
	return New(bytes.NewReader(b))
}

// New creates new Tokenizer from the provided reader with DefaultSeparators.
func New(r io.Reader) *Tokenizer {
	return &Tokenizer{
		r:             bufio.NewReader(r),
		separators:    DefaultSeparators,
		keepSeparator: false,
	}
}

// Tokenizer defines a struct that parses a reader into tokens while
// respecting quotes and parenthesis boundaries.
type Tokenizer struct {
	r *bufio.Reader

	separators    []rune
	keepSeparator bool
}

// Separators defines the provided separatos of the current Tokenizer.
func (s *Tokenizer) Separators(separators ...rune) {
	s.separators = separators
}

// KeepSeparator defines whether to keep the separator rune as part
// of the token (default to false).
func (s *Tokenizer) KeepSeparator(state bool) {
	s.keepSeparator = state
}

// Scan reads and returns the next available token from the Tokenizer's buffer (trimmed).
//
// Returns [io.EOF] error when there are no more tokens to scan.
func (s *Tokenizer) Scan() (string, error) {
	ch := s.read()

	if ch == eof {
		return "", io.EOF
	}

	if isWhitespaceRune(ch) {
		s.readWhiteSpaces()
	} else {
		s.unread()
	}

	token, err := s.readToken()
	if err != nil {
		return "", err
	}

	// read all remaining whitespaces
	s.readWhiteSpaces()

	return token, err
}

// ScanAll reads the entire Tokenizer's buffer and return all found tokens.
func (s *Tokenizer) ScanAll() ([]string, error) {
	tokens := []string{}

	for {
		token, err := s.Scan()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

// readToken reads a single token from the buffer and returns it.
func (s *Tokenizer) readToken() (string, error) {
	var buf bytes.Buffer
	var parenthesis int
	var quoteCh rune
	var prevCh rune

	for {
		ch := s.read()

		if ch == eof {
			break
		}

		if !isEscapeRune(prevCh) {
			if ch == '(' && quoteCh == eof {
				parenthesis++
			} else if ch == ')' && parenthesis > 0 && quoteCh == eof {
				parenthesis--
			} else if isQuoteRune(ch) {
				if quoteCh == ch {
					quoteCh = eof // reached closing quote
				} else if quoteCh == eof {
					quoteCh = ch // opening quote
				}
			}
		}

		if s.isSeperatorRune(ch) && parenthesis == 0 && quoteCh == eof {
			if s.keepSeparator {
				buf.WriteRune(ch)
			}
			break
		}

		prevCh = ch
		buf.WriteRune(ch)
	}

	if parenthesis > 0 || quoteCh != eof {
		return "", fmt.Errorf("unbalanced parenthesis or quoted expression: %q", buf.String())
	}

	return buf.String(), nil
}

// readWhiteSpaces consumes all contiguous whitespace runes.
func (s *Tokenizer) readWhiteSpaces() {
	for {
		ch := s.read()

		if ch == eof {
			break
		}

		if !s.isSeperatorRune(ch) {
			s.unread()
			break
		}
	}
}

// read reads the next rune from the buffered reader.
// Returns the `rune(0)` if an error or `io.EOF` occurs.
func (s *Tokenizer) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}

	return ch
}

// unread places the previously read rune back on the reader.
func (s *Tokenizer) unread() error {
	return s.r.UnreadRune()
}

// isSeperatorRune checks if a rune is a token part separator.
func (s *Tokenizer) isSeperatorRune(ch rune) bool {
	for _, r := range s.separators {
		if ch == r {
			return true
		}
	}

	return false
}

// isWhitespaceRune checks if a rune is a space, tab, or newline.
func isWhitespaceRune(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

// isQuoteRune checks if a rune is a quote.
func isQuoteRune(ch rune) bool {
	return ch == '\'' || ch == '"' || ch == '`'
}

// isEscapeRune checks if a rune is an escape character.
func isEscapeRune(ch rune) bool {
	return ch == '\\'
}
