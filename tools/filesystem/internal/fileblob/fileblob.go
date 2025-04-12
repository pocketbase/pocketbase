// Package fileblob provides a blob.Bucket driver implementation.
//
// NB! To minimize breaking changes with older PocketBase releases,
// the driver is a stripped down and adapted version of the previously
// used gocloud.dev/blob/fileblob, hence many of the below doc comments,
// struct options and interface implementations are the same.
//
// To avoid partial writes, fileblob writes to a temporary file and then renames
// the temporary file to the final path on Close. By default, it creates these
// temporary files in `os.TempDir`. If `os.TempDir` is on a different mount than
// your base bucket path, the `os.Rename` will fail with `invalid cross-device link`.
// To avoid this, either configure the temp dir to use by setting the environment
// variable `TMPDIR`, or set `Options.NoTempDir` to `true` (fileblob will create
// the temporary files next to the actual files instead of in a temporary directory).
//
// By default fileblob stores blob metadata in "sidecar" files under the original
// filename with an additional ".attrs" suffix.
// This behaviour can be changed via `Options.Metadata`;
// writing of those metadata files can be suppressed by setting it to
// `MetadataDontWrite` or its equivalent "metadata=skip" in the URL for the opener.
// In either case, absent any stored metadata many `blob.Attributes` fields
// will be set to default values.
//
// The blob abstraction supports all UTF-8 strings; to make this work with services lacking
// full UTF-8 support, strings must be escaped (during writes) and unescaped
// (during reads). The following escapes are performed for fileblob:
//   - Blob keys: ASCII characters 0-31 are escaped to "__0x<hex>__".
//     If os.PathSeparator != "/", it is also escaped.
//     Additionally, the "/" in "../", the trailing "/" in "//", and a trailing
//     "/" is key names are escaped in the same way.
//     On Windows, the characters "<>:"|?*" are also escaped.
//
// Example:
//
//	drv, _ := fileblob.New("/path/to/dir", nil)
//	bucket := blob.NewBucket(drv)
package fileblob

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/tools/filesystem/blob"
)

const defaultPageSize = 1000

type metadataOption string // Not exported as subject to change.

// Settings for Options.Metadata.
const (
	// Metadata gets written to a separate file.
	MetadataInSidecar metadataOption = ""

	// Writes won't carry metadata, as per the package docstring.
	MetadataDontWrite metadataOption = "skip"
)

// Options sets options for constructing a *blob.Bucket backed by fileblob.
type Options struct {
	// Refers to the strategy for how to deal with metadata (such as blob.Attributes).
	// For supported values please see the Metadata* constants.
	// If left unchanged, 'MetadataInSidecar' will be used.
	Metadata metadataOption

	// The FileMode to use when creating directories for the top-level directory
	// backing the bucket (when CreateDir is true), and for subdirectories for keys.
	// Defaults to 0777.
	DirFileMode os.FileMode

	// If true, create the directory backing the Bucket if it does not exist
	// (using os.MkdirAll).
	CreateDir bool

	// If true, don't use os.TempDir for temporary files, but instead place them
	// next to the actual files. This may result in "stranded" temporary files
	// (e.g., if the application is killed before the file cleanup runs).
	//
	// If your bucket directory is on a different mount than os.TempDir, you will
	// need to set this to true, as os.Rename will fail across mount points.
	NoTempDir bool
}

// New creates a new instance of the fileblob driver backed by the
// filesystem and rooted at dir, which must exist.
func New(dir string, opts *Options) (blob.Driver, error) {
	if opts == nil {
		opts = &Options{}
	}
	if opts.DirFileMode == 0 {
		opts.DirFileMode = os.FileMode(0o777)
	}

	absdir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s into an absolute path: %v", dir, err)
	}

	// Optionally, create the directory if it does not already exist.
	info, err := os.Stat(absdir)
	if err != nil && opts.CreateDir && os.IsNotExist(err) {
		err = os.MkdirAll(absdir, opts.DirFileMode)
		if err != nil {
			return nil, fmt.Errorf("tried to create directory but failed: %v", err)
		}
		info, err = os.Stat(absdir)
	}
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", absdir)
	}

	return &driver{dir: absdir, opts: opts}, nil
}

type driver struct {
	opts *Options
	dir  string
}

// Close implements [blob/Driver.Close].
func (drv *driver) Close() error {
	return nil
}

// NormalizeError implements [blob/Driver.NormalizeError].
func (drv *driver) NormalizeError(err error) error {
	if os.IsNotExist(err) {
		return errors.Join(err, blob.ErrNotFound)
	}

	return err
}

// path returns the full path for a key.
func (drv *driver) path(key string) (string, error) {
	path := filepath.Join(drv.dir, escapeKey(key))

	if strings.HasSuffix(path, attrsExt) {
		return "", errAttrsExt
	}

	return path, nil
}

// forKey returns the full path, os.FileInfo, and attributes for key.
func (drv *driver) forKey(key string) (string, os.FileInfo, *xattrs, error) {
	path, err := drv.path(key)
	if err != nil {
		return "", nil, nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", nil, nil, err
	}

	if info.IsDir() {
		return "", nil, nil, os.ErrNotExist
	}

	xa, err := getAttrs(path)
	if err != nil {
		return "", nil, nil, err
	}

	return path, info, &xa, nil
}

// ListPaged implements [blob/Driver.ListPaged].
func (drv *driver) ListPaged(ctx context.Context, opts *blob.ListOptions) (*blob.ListPage, error) {
	var pageToken string
	if len(opts.PageToken) > 0 {
		pageToken = string(opts.PageToken)
	}

	pageSize := opts.PageSize
	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	// If opts.Delimiter != "", lastPrefix contains the last "directory" key we
	// added. It is used to avoid adding it again; all files in this "directory"
	// are collapsed to the single directory entry.
	var lastPrefix string
	var lastKeyAdded string

	// If the Prefix contains a "/", we can set the root of the Walk
	// to the path specified by the Prefix as any files below the path will not
	// match the Prefix.
	// Note that we use "/" explicitly and not os.PathSeparator, as the opts.Prefix
	// is in the unescaped form.
	root := drv.dir
	if i := strings.LastIndex(opts.Prefix, "/"); i > -1 {
		root = filepath.Join(root, opts.Prefix[:i])
	}

	var result blob.ListPage

	// Do a full recursive scan of the root directory.
	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			// Couldn't read this file/directory for some reason; just skip it.
			return nil
		}

		// Skip the self-generated attribute files.
		if strings.HasSuffix(path, attrsExt) {
			return nil
		}

		// os.Walk returns the root directory; skip it.
		if path == drv.dir {
			return nil
		}

		// Strip the <drv.dir> prefix from path.
		prefixLen := len(drv.dir)
		// Include the separator for non-root.
		if drv.dir != "/" {
			prefixLen++
		}
		path = path[prefixLen:]

		// Unescape the path to get the key.
		key := unescapeKey(path)

		// Skip all directories. If opts.Delimiter is set, we'll create
		// pseudo-directories later.
		// Note that returning nil means that we'll still recurse into it;
		// we're just not adding a result for the directory itself.
		if info.IsDir() {
			key += "/"
			// Avoid recursing into subdirectories if the directory name already
			// doesn't match the prefix; any files in it are guaranteed not to match.
			if len(key) > len(opts.Prefix) && !strings.HasPrefix(key, opts.Prefix) {
				return filepath.SkipDir
			}
			// Similarly, avoid recursing into subdirectories if we're making
			// "directories" and all of the files in this subdirectory are guaranteed
			// to collapse to a "directory" that we've already added.
			if lastPrefix != "" && strings.HasPrefix(key, lastPrefix) {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip files/directories that don't match the Prefix.
		if !strings.HasPrefix(key, opts.Prefix) {
			return nil
		}

		var md5 []byte
		if xa, err := getAttrs(path); err == nil {
			// Note: we only have the MD5 hash for blobs that we wrote.
			// For other blobs, md5 will remain nil.
			md5 = xa.MD5
		}

		fi, err := info.Info()
		if err != nil {
			return err
		}

		obj := &blob.ListObject{
			Key:     key,
			ModTime: fi.ModTime(),
			Size:    fi.Size(),
			MD5:     md5,
		}

		// If using Delimiter, collapse "directories".
		if opts.Delimiter != "" {
			// Strip the prefix, which may contain Delimiter.
			keyWithoutPrefix := key[len(opts.Prefix):]
			// See if the key still contains Delimiter.
			// If no, it's a file and we just include it.
			// If yes, it's a file in a "sub-directory" and we want to collapse
			// all files in that "sub-directory" into a single "directory" result.
			if idx := strings.Index(keyWithoutPrefix, opts.Delimiter); idx != -1 {
				prefix := opts.Prefix + keyWithoutPrefix[0:idx+len(opts.Delimiter)]
				// We've already included this "directory"; don't add it.
				if prefix == lastPrefix {
					return nil
				}
				// Update the object to be a "directory".
				obj = &blob.ListObject{
					Key:   prefix,
					IsDir: true,
				}
				lastPrefix = prefix
			}
		}

		// If there's a pageToken, skip anything before it.
		if pageToken != "" && obj.Key <= pageToken {
			return nil
		}

		// If we've already got a full page of results, set NextPageToken and stop.
		// Unless the current object is a directory, in which case there may
		// still be objects coming that are alphabetically before it (since
		// we appended the delimiter). In that case, keep going; we'll trim the
		// extra entries (if any) before returning.
		if len(result.Objects) == pageSize && !obj.IsDir {
			result.NextPageToken = []byte(result.Objects[pageSize-1].Key)
			return io.EOF
		}

		result.Objects = append(result.Objects, obj)

		// Normally, objects are added in the correct order (by Key).
		// However, sometimes adding the file delimiter messes that up
		// (e.g., if the file delimiter is later in the alphabet than the last character of a key).
		// Detect if this happens and swap if needed.
		if len(result.Objects) > 1 && obj.Key < lastKeyAdded {
			i := len(result.Objects) - 1
			result.Objects[i-1], result.Objects[i] = result.Objects[i], result.Objects[i-1]
			lastKeyAdded = result.Objects[i].Key
		} else {
			lastKeyAdded = obj.Key
		}

		return nil
	})
	if err != nil && err != io.EOF {
		return nil, err
	}

	if len(result.Objects) > pageSize {
		result.Objects = result.Objects[0:pageSize]
		result.NextPageToken = []byte(result.Objects[pageSize-1].Key)
	}

	return &result, nil
}

// Attributes implements [blob/Driver.Attributes].
func (drv *driver) Attributes(ctx context.Context, key string) (*blob.Attributes, error) {
	_, info, xa, err := drv.forKey(key)
	if err != nil {
		return nil, err
	}

	return &blob.Attributes{
		CacheControl:       xa.CacheControl,
		ContentDisposition: xa.ContentDisposition,
		ContentEncoding:    xa.ContentEncoding,
		ContentLanguage:    xa.ContentLanguage,
		ContentType:        xa.ContentType,
		Metadata:           xa.Metadata,
		// CreateTime left as the zero time.
		ModTime: info.ModTime(),
		Size:    info.Size(),
		MD5:     xa.MD5,
		ETag:    fmt.Sprintf("\"%x-%x\"", info.ModTime().UnixNano(), info.Size()),
	}, nil
}

// NewRangeReader implements [blob/Driver.NewRangeReader].
func (drv *driver) NewRangeReader(ctx context.Context, key string, offset, length int64) (blob.DriverReader, error) {
	path, info, xa, err := drv.forKey(key)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if offset > 0 {
		if _, err := f.Seek(offset, io.SeekStart); err != nil {
			return nil, err
		}
	}

	r := io.Reader(f)
	if length >= 0 {
		r = io.LimitReader(r, length)
	}

	return &reader{
		r: r,
		c: f,
		attrs: &blob.ReaderAttributes{
			ContentType: xa.ContentType,
			ModTime:     info.ModTime(),
			Size:        info.Size(),
		},
	}, nil
}

func createTemp(path string, noTempDir bool) (*os.File, error) {
	// Use a custom createTemp function rather than os.CreateTemp() as
	// os.CreateTemp() sets the permissions of the tempfile to 0600, rather than
	// 0666, making it inconsistent with the directories and attribute files.
	try := 0
	for {
		// Append the current time with nanosecond precision and .tmp to the
		// base path. If the file already exists try again. Nanosecond changes enough
		// between each iteration to make a conflict unlikely. Using the full
		// time lowers the chance of a collision with a file using a similar
		// pattern, but has undefined behavior after the year 2262.
		var name string
		if noTempDir {
			name = path
		} else {
			name = filepath.Join(os.TempDir(), filepath.Base(path))
		}
		name += "." + strconv.FormatInt(time.Now().UnixNano(), 16) + ".tmp"

		f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o666)
		if os.IsExist(err) {
			if try++; try < 10000 {
				continue
			}
			return nil, &os.PathError{Op: "createtemp", Path: path + ".*.tmp", Err: os.ErrExist}
		}

		return f, err
	}
}

// NewTypedWriter implements [blob/Driver.NewTypedWriter].
func (drv *driver) NewTypedWriter(ctx context.Context, key, contentType string, opts *blob.WriterOptions) (blob.DriverWriter, error) {
	path, err := drv.path(key)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(filepath.Dir(path), drv.opts.DirFileMode)
	if err != nil {
		return nil, err
	}

	f, err := createTemp(path, drv.opts.NoTempDir)
	if err != nil {
		return nil, err
	}

	if drv.opts.Metadata == MetadataDontWrite {
		w := &writer{
			ctx:  ctx,
			File: f,
			path: path,
		}
		return w, nil
	}

	var metadata map[string]string
	if len(opts.Metadata) > 0 {
		metadata = opts.Metadata
	}

	return &writerWithSidecar{
		ctx:        ctx,
		f:          f,
		path:       path,
		contentMD5: opts.ContentMD5,
		md5hash:    md5.New(),
		attrs: xattrs{
			CacheControl:       opts.CacheControl,
			ContentDisposition: opts.ContentDisposition,
			ContentEncoding:    opts.ContentEncoding,
			ContentLanguage:    opts.ContentLanguage,
			ContentType:        contentType,
			Metadata:           metadata,
		},
	}, nil
}

// Copy implements [blob/Driver.Copy].
func (drv *driver) Copy(ctx context.Context, dstKey, srcKey string) error {
	// Note: we could use NewRangeReader here, but since we need to copy all of
	// the metadata (from xa), it's more efficient to do it directly.
	srcPath, _, xa, err := drv.forKey(srcKey)
	if err != nil {
		return err
	}

	f, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// We'll write the copy using Writer, to avoid re-implementing making of a
	// temp file, cleaning up after partial failures, etc.
	wopts := blob.WriterOptions{
		CacheControl:       xa.CacheControl,
		ContentDisposition: xa.ContentDisposition,
		ContentEncoding:    xa.ContentEncoding,
		ContentLanguage:    xa.ContentLanguage,
		Metadata:           xa.Metadata,
	}

	// Create a cancelable context so we can cancel the write if there are problems.
	writeCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	w, err := drv.NewTypedWriter(writeCtx, dstKey, xa.ContentType, &wopts)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)
	if err != nil {
		cancel() // cancel before Close cancels the write
		w.Close()
		return err
	}

	return w.Close()
}

// Delete implements [blob/Driver.Delete].
func (b *driver) Delete(ctx context.Context, key string) error {
	path, err := b.path(key)
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil {
		return err
	}

	err = os.Remove(path + attrsExt)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// -------------------------------------------------------------------

type reader struct {
	r     io.Reader
	c     io.Closer
	attrs *blob.ReaderAttributes
}

func (r *reader) Read(p []byte) (int, error) {
	if r.r == nil {
		return 0, io.EOF
	}
	return r.r.Read(p)
}

func (r *reader) Close() error {
	if r.c == nil {
		return nil
	}
	return r.c.Close()
}

// Attributes implements [blob/DriverReader.Attributes].
func (r *reader) Attributes() *blob.ReaderAttributes {
	return r.attrs
}

// -------------------------------------------------------------------

// writerWithSidecar implements the strategy of storing metadata in a distinct file.
type writerWithSidecar struct {
	ctx        context.Context
	md5hash    hash.Hash
	f          *os.File
	path       string
	attrs      xattrs
	contentMD5 []byte
}

func (w *writerWithSidecar) Write(p []byte) (n int, err error) {
	n, err = w.f.Write(p)
	if err != nil {
		// Don't hash the unwritten tail twice when writing is resumed.
		w.md5hash.Write(p[:n])
		return n, err
	}

	if _, err := w.md5hash.Write(p); err != nil {
		return n, err
	}

	return n, nil
}

func (w *writerWithSidecar) Close() error {
	err := w.f.Close()
	if err != nil {
		return err
	}

	// Always delete the temp file. On success, it will have been
	// renamed so the Remove will fail.
	defer func() {
		_ = os.Remove(w.f.Name())
	}()

	// Check if the write was cancelled.
	if err := w.ctx.Err(); err != nil {
		return err
	}

	md5sum := w.md5hash.Sum(nil)
	w.attrs.MD5 = md5sum

	// Write the attributes file.
	if err := setAttrs(w.path, w.attrs); err != nil {
		return err
	}

	// Rename the temp file to path.
	if err := os.Rename(w.f.Name(), w.path); err != nil {
		_ = os.Remove(w.path + attrsExt)
		return err
	}

	return nil
}

// writer is a file with a temporary name until closed.
//
// Embedding os.File allows the likes of io.Copy to use optimizations,
// which is why it is not folded into writerWithSidecar.
type writer struct {
	*os.File
	ctx  context.Context
	path string
}

func (w *writer) Close() error {
	err := w.File.Close()
	if err != nil {
		return err
	}

	// Always delete the temp file. On success, it will have been renamed so
	// the Remove will fail.
	tempname := w.Name()
	defer os.Remove(tempname)

	// Check if the write was cancelled.
	if err := w.ctx.Err(); err != nil {
		return err
	}

	// Rename the temp file to path.
	return os.Rename(tempname, w.path)
}

// -------------------------------------------------------------------

// escapeKey does all required escaping for UTF-8 strings to work the filesystem.
func escapeKey(s string) string {
	s = blob.HexEscape(s, func(r []rune, i int) bool {
		c := r[i]
		switch {
		case c < 32:
			return true
		// We're going to replace '/' with os.PathSeparator below. In order for this
		// to be reversible, we need to escape raw os.PathSeparators.
		case os.PathSeparator != '/' && c == os.PathSeparator:
			return true
		// For "../", escape the trailing slash.
		case i > 1 && c == '/' && r[i-1] == '.' && r[i-2] == '.':
			return true
		// For "//", escape the trailing slash.
		case i > 0 && c == '/' && r[i-1] == '/':
			return true
		// Escape the trailing slash in a key.
		case c == '/' && i == len(r)-1:
			return true
		// https://docs.microsoft.com/en-us/windows/desktop/fileio/naming-a-file
		case os.PathSeparator == '\\' && (c == '>' || c == '<' || c == ':' || c == '"' || c == '|' || c == '?' || c == '*'):
			return true
		}
		return false
	})

	// Replace "/" with os.PathSeparator if needed, so that the local filesystem
	// can use subdirectories.
	if os.PathSeparator != '/' {
		s = strings.ReplaceAll(s, "/", string(os.PathSeparator))
	}

	return s
}

// unescapeKey reverses escapeKey.
func unescapeKey(s string) string {
	if os.PathSeparator != '/' {
		s = strings.ReplaceAll(s, string(os.PathSeparator), "/")
	}

	return blob.HexUnescape(s)
}
