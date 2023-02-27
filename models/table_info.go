package models

import "github.com/pocketbase/pocketbase/tools/types"

type TableInfoRow struct {
	// the `db:"pk"` tag has special semantic so we cannot rename
	// the original field without specifying a custom mapper
	PK int

	Index        int           `db:"cid"`
	Name         string        `db:"name"`
	Type         string        `db:"type"`
	NotNull      bool          `db:"notnull"`
	DefaultValue types.JsonRaw `db:"dflt_value"`
}
