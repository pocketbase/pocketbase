package osutils

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase/tools/list"
)

// MoveDirContent moves the src dir content, that is not listed in the exclude list,
// to dest dir (it will be created if missing).
//
// The rootExclude argument is used to specify a list of src root entries to exclude.
//
// Note that this method doesn't delete the old src dir.
//
// It is an alternative to os.Rename() for the cases where we can't
// rename/delete the src dir (see https://github.com/pocketbase/pocketbase/issues/2519).
func MoveDirContent(src string, dest string, rootExclude ...string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// make sure that the dest dir exist
	manuallyCreatedDestDir := false
	if _, err := os.Stat(dest); err != nil {
		if err := os.Mkdir(dest, os.ModePerm); err != nil {
			return err
		}
		manuallyCreatedDestDir = true
	}

	moved := map[string]string{}

	tryRollback := func() []error {
		errs := []error{}

		for old, new := range moved {
			if err := os.Rename(new, old); err != nil {
				errs = append(errs, err)
			}
		}

		// try to delete manually the created dest dir if all moved files were restored
		if manuallyCreatedDestDir && len(errs) == 0 {
			if err := os.Remove(dest); err != nil {
				errs = append(errs, err)
			}
		}

		return errs
	}

	for _, entry := range entries {
		basename := entry.Name()

		if list.ExistInSlice(basename, rootExclude) {
			continue
		}

		old := filepath.Join(src, basename)
		new := filepath.Join(dest, basename)

		if err := os.Rename(old, new); err != nil {
			if errs := tryRollback(); len(errs) > 0 {
				errs = append(errs, err)
				err = errors.Join(errs...)
			}

			return err
		}

		moved[old] = new
	}

	return nil
}
