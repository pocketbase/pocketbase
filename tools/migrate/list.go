package migrate

import (
	"path/filepath"
	"runtime"
	"sort"

	"github.com/pocketbase/dbx"
)

type Migration struct {
	File string
	Up   func(db dbx.Builder) error
	Down func(db dbx.Builder) error
}

// MigrationsList defines a list with migration definitions
type MigrationsList struct {
	list []*Migration
}

// Item returns a single migration from the list by its index.
func (l *MigrationsList) Item(index int) *Migration {
	return l.list[index]
}

// Items returns the internal migrations list slice.
func (l *MigrationsList) Items() []*Migration {
	return l.list
}

// Register adds new migration definition to the list.
//
// If `optFilename` is not provided, it will try to get the name from its .go file.
//
// The list will be sorted automatically based on the migrations file name.
func (l *MigrationsList) Register(
	up func(db dbx.Builder) error,
	down func(db dbx.Builder) error,
	optFilename ...string,
) {
	var file string
	if len(optFilename) > 0 {
		file = optFilename[0]
	} else {
		_, path, _, _ := runtime.Caller(1)
		file = filepath.Base(path)
	}

	l.list = append(l.list, &Migration{
		File: file,
		Up:   up,
		Down: down,
	})

	sort.Slice(l.list, func(i int, j int) bool {
		return l.list[i].File < l.list[j].File
	})
}
