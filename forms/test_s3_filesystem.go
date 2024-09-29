package forms

import (
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/security"
)

const (
	s3FilesystemStorage = "storage"
	s3FilesystemBackups = "backups"
)

// TestS3Filesystem defines a S3 filesystem connection test.
type TestS3Filesystem struct {
	app core.App

	// The name of the filesystem - storage or backups
	Filesystem string `form:"filesystem" json:"filesystem"`
}

// NewTestS3Filesystem creates and initializes new TestS3Filesystem form.
func NewTestS3Filesystem(app core.App) *TestS3Filesystem {
	return &TestS3Filesystem{app: app}
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *TestS3Filesystem) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Filesystem,
			validation.Required,
			validation.In(s3FilesystemStorage, s3FilesystemBackups),
		),
	)
}

// Submit validates and performs a S3 filesystem connection test.
func (form *TestS3Filesystem) Submit() error {
	if err := form.Validate(); err != nil {
		return err
	}

	var s3Config core.S3Config

	if form.Filesystem == s3FilesystemBackups {
		s3Config = form.app.Settings().Backups.S3
	} else {
		s3Config = form.app.Settings().S3
	}

	if !s3Config.Enabled {
		return errors.New("S3 storage filesystem is not enabled")
	}

	fsys, err := filesystem.NewS3(
		s3Config.Bucket,
		s3Config.Region,
		s3Config.Endpoint,
		s3Config.AccessKey,
		s3Config.Secret,
		s3Config.ForcePathStyle,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize the S3 filesystem: %w", err)
	}
	defer fsys.Close()

	testPrefix := "pb_settings_test_" + security.PseudorandomString(5)
	testFileKey := testPrefix + "/test.txt"

	// try to upload a test file
	if err := fsys.Upload([]byte("test"), testFileKey); err != nil {
		return fmt.Errorf("failed to upload a test file: %w", err)
	}

	// test prefix deletion (ensures that both bucket list and delete works)
	if errs := fsys.DeletePrefix(testPrefix); len(errs) > 0 {
		return fmt.Errorf("failed to delete a test file: %w", errs[0])
	}

	return nil
}
