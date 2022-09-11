package mailer

import (
	"regexp"
	"strings"

	"github.com/pocketbase/pocketbase/tools/list"
	"golang.org/x/net/html"
)

var whitespaceRegex = regexp.MustCompile(`\s+`)

// Very rudimentary auto HTML to Text mail body converter.
//
// Caveats:
// - This method doesn't check for correctness of the HTML document.
// - Links will be converted to "[text](url)" format.
// - List items (<li>) are prefixed with "- ".
// - Indentation is stripped (both tabs and spaces).
// - Trailing spaces are preserved.
// - Multiple consequence newlines are collapsed as one unless multiple <br> tags are used.
func html2Text(htmlDocument string) (string, error) {
	var builder strings.Builder

	doc, err := html.Parse(strings.NewReader(htmlDocument))
	if err != nil {
		return "", err
	}

	tagsToSkip := []string{
		"style", "script", "iframe", "applet", "object", "svg", "img",
		"button", "form", "textarea", "input", "select", "option", "template",
	}

	inlineTags := []string{
		"a", "span", "small", "strike", "strong",
		"sub", "sup", "em", "b", "u", "i",
	}

	var canAddNewLine bool

	// see https://pkg.go.dev/golang.org/x/net/html#Parse
	var f func(*html.Node)
	f = func(n *html.Node) {
		// start link wrapping for producing "[text](link)" formatted string
		isLink := n.Type == html.ElementNode && n.Data == "a"
		if isLink {
			builder.WriteString("[")
		}

		switch n.Type {
		case html.TextNode:
			txt := whitespaceRegex.ReplaceAllString(n.Data, " ")

			// the prev node has new line so it is safe to trim the indentation
			if !canAddNewLine {
				txt = strings.TrimLeft(txt, " ")
			}

			if txt != "" {
				builder.WriteString(txt)
				canAddNewLine = true
			}
		case html.ElementNode:
			if n.Data == "br" {
				// always write new lines when <br> tag is used
				builder.WriteString("\r\n")
				canAddNewLine = false
			} else if canAddNewLine && !list.ExistInSlice(n.Data, inlineTags) {
				builder.WriteString("\r\n")
				canAddNewLine = false
			}

			// prefix list items with dash
			if n.Data == "li" {
				builder.WriteString("- ")
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode || !list.ExistInSlice(c.Data, tagsToSkip) {
				f(c)
			}
		}

		// end link wrapping
		if isLink {
			builder.WriteString("]")
			for _, a := range n.Attr {
				if a.Key == "href" {
					if a.Val != "" {
						builder.WriteString("(")
						builder.WriteString(a.Val)
						builder.WriteString(")")
					}
					break
				}
			}
		}
	}

	f(doc)

	return strings.TrimSpace(builder.String()), nil
}
