package cli

import (
	"bufio"
	"fmt"
	"strings"
)

func (c *Cli) Confirm(message string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(c.reader)
	var s string

	for {
		fmt.Fprintf(c.writer, "%s (%s) ", message, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}
