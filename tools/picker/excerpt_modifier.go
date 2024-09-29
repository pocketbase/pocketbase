package picker

import (
	"errors"
	"regexp"
	"strings"

	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/spf13/cast"
	"golang.org/x/net/html"
)

func init() {
	Modifiers["excerpt"] = func(args ...string) (Modifier, error) {
		return newExcerptModifier(args...)
	}
}

var whitespaceRegex = regexp.MustCompile(`\s+`)

var excludeTags = []string{
	"head", "style", "script", "iframe", "embed", "applet", "object",
	"svg", "img", "picture", "dialog", "template", "button", "form",
	"textarea", "input", "select", "option",
}

var inlineTags = []string{
	"a", "abbr", "acronym", "b", "bdo", "big", "br", "button",
	"cite", "code", "em", "i", "label", "q", "small", "span",
	"strong", "strike", "sub", "sup", "time",
}

var _ Modifier = (*excerptModifier)(nil)

type excerptModifier struct {
	max          int  // approximate max excerpt length
	withEllipsis bool // if enabled will add ellipsis when the plain text length > max
}

// newExcerptModifier validates the specified raw string arguments and
// initializes a new excerptModifier.
//
// This method is usually invoked in initModifer().
func newExcerptModifier(args ...string) (*excerptModifier, error) {
	totalArgs := len(args)

	if totalArgs == 0 {
		return nil, errors.New("max argument is required - expected (max, withEllipsis?)")
	}

	if totalArgs > 2 {
		return nil, errors.New("too many arguments - expected (max, withEllipsis?)")
	}

	max := cast.ToInt(args[0])
	if max == 0 {
		return nil, errors.New("max argument must be > 0")
	}

	var withEllipsis bool
	if totalArgs > 1 {
		withEllipsis = cast.ToBool(args[1])
	}

	return &excerptModifier{max, withEllipsis}, nil
}

// Modify implements the [Modifier.Modify] interface method.
//
// It returns a plain text excerpt/short-description from a formatted
// html string (non-string values are kept untouched).
func (m *excerptModifier) Modify(value any) (any, error) {
	strValue, ok := value.(string)
	if !ok {
		// not a string -> return as it is without applying the modifier
		// (we don't throw an error because the modifier could be applied for a missing expand field)
		return value, nil
	}

	var builder strings.Builder

	doc, err := html.Parse(strings.NewReader(strValue))
	if err != nil {
		return "", err
	}

	var hasPrevSpace bool

	// for all node types and more details check
	// https://pkg.go.dev/golang.org/x/net/html#Parse
	var stripTags func(*html.Node)
	stripTags = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			// collapse multiple spaces into one
			txt := whitespaceRegex.ReplaceAllString(n.Data, " ")

			if hasPrevSpace {
				txt = strings.TrimLeft(txt, " ")
			}

			if txt != "" {
				hasPrevSpace = strings.HasSuffix(txt, " ")

				builder.WriteString(txt)
			}
		}

		// excerpt max has been reached => no need to further iterate
		// (+2 for the extra whitespace suffix/prefix that will be trimmed later)
		if builder.Len() > m.max+2 {
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode || !list.ExistInSlice(c.Data, excludeTags) {
				isBlock := c.Type == html.ElementNode && !list.ExistInSlice(c.Data, inlineTags)

				if isBlock && !hasPrevSpace {
					builder.WriteString(" ")
					hasPrevSpace = true
				}

				stripTags(c)

				if isBlock && !hasPrevSpace {
					builder.WriteString(" ")
					hasPrevSpace = true
				}
			}
		}
	}
	stripTags(doc)

	result := strings.TrimSpace(builder.String())

	if len(result) > m.max {
		result = strings.TrimSpace(result[:m.max])

		if m.withEllipsis {
			result += "..."
		}
	}

	return result, nil
}
