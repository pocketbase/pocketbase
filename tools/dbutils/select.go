package dbutils

import "regexp"

// Regexp for columns and tables (the same as the one in dbx).
var selectRegex = regexp.MustCompile(`(?i:\s+as\s+|\s+)([\w\-_\.]+)$`)

// AliasOrIdentifier returns the alias from a column or table identifier.
// Returns the identifier unmodified if no alias was found.
func AliasOrIdentifier(columnOrTableIdentifier string) string {
	matches := selectRegex.FindStringSubmatch(columnOrTableIdentifier)

	if len(matches) > 0 && matches[1] != "" {
		return matches[1]
	}

	return columnOrTableIdentifier
}
