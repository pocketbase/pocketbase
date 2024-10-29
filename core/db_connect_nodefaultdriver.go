//go:build nodefaultdriver

package core

import "github.com/pocketbase/dbx"

func DefaultDBConnect(dbPath string) (*dbx.DB, error) {
	panic("DBConnect config option must be set when the nodefaultdriver tag is used!")
}
