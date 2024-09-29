package filesystem

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
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
// The file could be from a local path, multipart/form-data header, etc.
type File struct {
	Reader       FileReader `form:"-" json:"-" xml:"-"`
	Name         string     `form:"name" json:"name" xml:"name"`
	OriginalName string     `form:"originalName" json:"originalName" xml:"originalName"`
	Size         int64      `form:"size" json:"size" xml:"size"`
}

// AsMap implements [core.mapExtractor] and returns a value suitable
// to be used in an API rule expression.
func (f *File) AsMap() map[string]any {
	return map[string]any{
		"name":         f.Name,
		"originalName": f.OriginalName,
		"size":         f.Size,
	}
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

// NewFileFromBytes creates a new File instance from the provided byte slice.
func NewFileFromBytes(b []byte, name string) (*File, error) {
	size := len(b)
	if size == 0 {
		return nil, errors.New("cannot create an empty file")
	}

	f := &File{}

	f.Reader = &BytesReader{b}
	f.Size = int64(size)
	f.OriginalName = name
	f.Name = normalizeName(f.Reader, f.OriginalName)

	return f, nil
}

// NewFileFromMultipart creates a new File from the provided multipart header.
func NewFileFromMultipart(mh *multipart.FileHeader) (*File, error) {
	f := &File{}

	f.Reader = &MultipartReader{Header: mh}
	f.Size = mh.Size
	f.OriginalName = mh.Filename
	f.Name = normalizeName(f.Reader, f.OriginalName)

	return f, nil
}

// NewFileFromURL creates a new File from the provided url by
// downloading the resource and load it as BytesReader.
//
// Example
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	file, err := filesystem.NewFileFromURL(ctx, "https://example.com/image.png")
func NewFileFromURL(ctx context.Context, url string) (*File, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 399 {
		return nil, fmt.Errorf("failed to download url %s (%d)", url, res.StatusCode)
	}

	var buf bytes.Buffer

	if _, err = io.Copy(&buf, res.Body); err != nil {
		return nil, err
	}

	return NewFileFromBytes(buf.Bytes(), path.Base(url))
}

// -------------------------------------------------------------------

var _ FileReader = (*MultipartReader)(nil)

// MultipartReader defines a FileReader from [multipart.FileHeader].
type MultipartReader struct {
	Header *multipart.FileHeader
}

// Open implements the [filesystem.FileReader] interface.
func (r *MultipartReader) Open() (io.ReadSeekCloser, error) {
	return r.Header.Open()
}

// -------------------------------------------------------------------

var _ FileReader = (*PathReader)(nil)

// PathReader defines a FileReader from a local file path.
type PathReader struct {
	Path string
}

// Open implements the [filesystem.FileReader] interface.
func (r *PathReader) Open() (io.ReadSeekCloser, error) {
	return os.Open(r.Path)
}

// -------------------------------------------------------------------

var _ FileReader = (*BytesReader)(nil)

// BytesReader defines a FileReader from bytes content.
type BytesReader struct {
	Bytes []byte
}

// Open implements the [filesystem.FileReader] interface.
func (r *BytesReader) Open() (io.ReadSeekCloser, error) {
	return &bytesReadSeekCloser{bytes.NewReader(r.Bytes)}, nil
}

type bytesReadSeekCloser struct {
	*bytes.Reader
}

// Close implements the [io.ReadSeekCloser] interface.
func (r *bytesReadSeekCloser) Close() error {
	return nil
}

// -------------------------------------------------------------------

var extInvalidCharsRegex = regexp.MustCompile(`[^\w\.\*\-\+\=\#]+`)

const randomAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

func normalizeName(fr FileReader, name string) string {
	// extension
	// ---
	originalExt := extractExtension(name)
	cleanExt := extInvalidCharsRegex.ReplaceAllString(originalExt, "")
	if cleanExt == "" {
		// try to detect the extension from the file content
		cleanExt, _ = detectExtension(fr)
	}
	if extLength := len(cleanExt); extLength > 20 {
		// keep only the last 20 characters (it is multibyte safe after the regex replace)
		cleanExt = "." + cleanExt[extLength-20:]
	}

	// name
	// ---
	cleanName := inflector.Snakecase(strings.TrimSuffix(name, originalExt))
	if length := len(cleanName); length < 3 {
		// the name is too short so we concatenate an additional random part
		cleanName += security.RandomStringWithAlphabet(10, randomAlphabet)
	} else if length > 100 {
		// keep only the first 100 characters (it is multibyte safe after Snakecase)
		cleanName = cleanName[:100]
	}

	return fmt.Sprintf(
		"%s_%s%s",
		cleanName,
		security.RandomStringWithAlphabet(10, randomAlphabet), // ensure that there is always a random part
		cleanExt,
	)
}

// extractExtension extracts the extension (with leading dot) from name.
//
// This differ from filepath.Ext() by supporting double extensions (eg. ".tar.gz").
//
// Returns an empty string if no match is found.
//
// Example:
// extractExtension("test.txt")      // .txt
// extractExtension("test.tar.gz")   // .tar.gz
// extractExtension("test.a.tar.gz") // .tar.gz
func extractExtension(name string) string {
	primaryDot := strings.LastIndex(name, ".")

	if primaryDot == -1 {
		return ""
	}

	// look for secondary extension
	secondaryDot := strings.LastIndex(name[:primaryDot], ".")
	if secondaryDot >= 0 {
		return name[secondaryDot:]
	}

	return name[primaryDot:]
}

// detectExtension tries to detect the extension from file mime type.
func detectExtension(fr FileReader) (string, error) {
	r, err := fr.Open()
	if err != nil {
		return "", err
	}
	defer r.Close()

	mt, err := mimetype.DetectReader(r)
	if err != nil {
		return "", err
	}

	return mt.Extension(), nil
}
