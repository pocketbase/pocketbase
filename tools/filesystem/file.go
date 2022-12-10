package filesystem

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/security"
)

// FileReader defines an interface for a file resource reader.
type FileReader interface {
	Open() (io.ReadSeekCloser, error)
}

// File defines a single file [io.ReadSeekCloser] resource.
//
// The file could be from a local path, multipipart/formdata header, etc.
type File struct {
	Name         string
	OriginalName string
	Size         int64
	Reader       FileReader
}

// NewFileFromPath creates a new File instance from the provided local file path.
func NewFileFromPath(path string) (*File, error) {
	f := &File{}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f.Reader = &PathReader{Path: path}
	f.Size = info.Size()
	f.OriginalName = info.Name()
	f.Name = normalizeName(f.Reader, f.OriginalName)

	return f, nil
}

// NewFileFromMultipart creates a new File instace from the provided multipart header.
func NewFileFromMultipart(mh *multipart.FileHeader) (*File, error) {
	f := &File{}

	f.Reader = &MultipartReader{Header: mh}
	f.Size = mh.Size
	f.OriginalName = mh.Filename
	f.Name = normalizeName(f.Reader, f.OriginalName)

	return f, nil
}

// -------------------------------------------------------------------

var _ FileReader = (*MultipartReader)(nil)

// MultipartReader defines a [multipart.FileHeader] reader.
type MultipartReader struct {
	Header *multipart.FileHeader
}

// Open implements the [filesystem.FileReader] interface.
func (r *MultipartReader) Open() (io.ReadSeekCloser, error) {
	return r.Header.Open()
}

// -------------------------------------------------------------------

var _ FileReader = (*PathReader)(nil)

type PathReader struct {
	Path string
}

// Open implements the [filesystem.FileReader] interface.
func (r *PathReader) Open() (io.ReadSeekCloser, error) {
	return os.Open(r.Path)
}

// -------------------------------------------------------------------

var extInvalidCharsRegex = regexp.MustCompile(`[^\w\.\*\-\+\=\#]+`)

func normalizeName(fr FileReader, name string) string {
	// extension
	// ---
	originalExt := filepath.Ext(name)
	cleanExt := extInvalidCharsRegex.ReplaceAllString(originalExt, "")
	if cleanExt == "" {
		// try to detect the extension from the file content
		cleanExt, _ = detectExtension(fr)
	}

	// name
	// ---
	cleanName := inflector.Snakecase(strings.TrimSuffix(name, originalExt))
	if length := len(cleanName); length < 3 {
		// the name is too short so we concatenate an additional random part
		cleanName += security.RandomString(10)
	} else if length > 100 {
		// keep only the first 100 characters (it is multibyte safe after Snakecase)
		cleanName = cleanName[:100]
	}

	return fmt.Sprintf(
		"%s_%s%s",
		cleanName,
		security.RandomString(10), // ensure that there is always a random part
		cleanExt,
	)
}

func detectExtension(fr FileReader) (string, error) {
	// try to detect the extension from the mime type
	r, err := fr.Open()
	if err != nil {
		return "", err
	}
	defer r.Close()

	mt, _ := mimetype.DetectReader(r)
	if err != nil {
		return "", err
	}

	return mt.Extension(), nil
}
