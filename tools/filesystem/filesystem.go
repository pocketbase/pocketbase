package filesystem

import (
	"context"
	"errors"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/fatih/color"
	"github.com/gabriel-vasile/mimetype"
	"github.com/pocketbase/pocketbase/tools/filesystem/blob"
	"github.com/pocketbase/pocketbase/tools/filesystem/internal/fileblob"
	"github.com/pocketbase/pocketbase/tools/filesystem/internal/s3blob"
	"github.com/pocketbase/pocketbase/tools/filesystem/internal/s3blob/s3"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/routine"

	// manually register the webp decoder because disintegration/imaging does not support webp
	_ "golang.org/x/image/webp"
)

// note: the same as blob.ErrNotFound for backward compatibility with earlier versions
var ErrNotFound = blob.ErrNotFound

const MetadataOriginalName = "original-filename"

type DeleteEvent struct {
	hook.Event
	Filesystem *System
	FileKey    string
}

type NewWriterEvent struct {
	hook.Event
	Filesystem *System
	FileKey    string
	Options    *blob.WriterOptions
	Writer     *blob.Writer // filled only after e.Next()
}

// @todo consider renaming

type System struct {
	ctx    context.Context
	bucket *blob.Bucket

	// @todo consider with the refactoring to bind on driver level
	onNewWriter *hook.Hook[*NewWriterEvent]
	onDelete    *hook.Hook[*DeleteEvent]
}

// NewS3 initializes a new S3 filesystem instance.
//
// NB! Make sure to call `Close()` after you are done working with it.
func NewS3(
	bucketName string,
	region string,
	endpoint string,
	accessKey string,
	secretKey string,
	s3ForcePathStyle bool,
) (*System, error) {
	ctx := context.Background() // default context

	client := &s3.S3{
		Bucket:       bucketName,
		Region:       region,
		Endpoint:     endpoint,
		AccessKey:    accessKey,
		SecretKey:    secretKey,
		UsePathStyle: s3ForcePathStyle,
	}

	drv, err := s3blob.New(client)
	if err != nil {
		return nil, err
	}

	return &System{ctx: ctx, bucket: blob.NewBucket(drv)}, nil
}

// NewLocal initializes a new local filesystem instance.
//
// NB! Make sure to call `Close()` after you are done working with it.
func NewLocal(dirPath string) (*System, error) {
	ctx := context.Background() // default context

	// makes sure that the directory exist
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return nil, err
	}

	drv, err := fileblob.New(dirPath, &fileblob.Options{
		NoTempDir: true,
	})
	if err != nil {
		return nil, err
	}

	return &System{ctx: ctx, bucket: blob.NewBucket(drv)}, nil
}

// SetContext assigns the specified context to the current filesystem.
func (s *System) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// Close releases any resources used for the related filesystem.
func (s *System) Close() error {
	if s.onNewWriter != nil {
		s.onNewWriter.UnbindAll()
		s.onNewWriter = nil
	}

	if s.onDelete != nil {
		s.onDelete.UnbindAll()
		s.onDelete = nil
	}

	return s.bucket.Close()
}

// OnNewWriter is a low level hook that is triggered on every new writer initialization
// (aka. when attempting to create a new file with [system.NewWriter] or [system.Upload]).
//
// Note that currently it doesn't trigger on [System.Copy] but this may change in future releases.
func (s *System) OnNewWriter() *hook.Hook[*NewWriterEvent] {
	if s.onNewWriter == nil {
		s.onNewWriter = &hook.Hook[*NewWriterEvent]{}
	}

	return s.onNewWriter
}

// OnDelete is a low level hook that is triggered on every [System.Delete] call.
//
// Note that the hook doesn't fire when a file is being overwritten
// by a new one, because in that case [System.Delete] is not invoked.
func (s *System) OnDelete() *hook.Hook[*DeleteEvent] {
	if s.onDelete == nil {
		s.onDelete = &hook.Hook[*DeleteEvent]{}
	}

	return s.onDelete
}

// Exists checks if file with fileKey path exists or not.
func (s *System) Exists(fileKey string) (bool, error) {
	return s.bucket.Exists(s.ctx, fileKey)
}

// Attributes returns the attributes for the file with fileKey path.
//
// If the file doesn't exist it returns ErrNotFound.
func (s *System) Attributes(fileKey string) (*blob.Attributes, error) {
	return s.bucket.Attributes(s.ctx, fileKey)
}

// GetReader returns a file content reader for the given fileKey.
//
// NB! Make sure to call Close() on the file after you are done working with it.
//
// If the file doesn't exist returns ErrNotFound.
func (s *System) GetReader(fileKey string) (*blob.Reader, error) {
	return s.bucket.NewReader(s.ctx, fileKey)
}

// Deprecated: Please use GetReader(fileKey) instead.
func (s *System) GetFile(fileKey string) (*blob.Reader, error) {
	color.Yellow("Deprecated: Please replace GetFile with GetReader.")
	return s.GetReader(fileKey)
}

// GetReuploadableFile constructs a new reuploadable File value
// from the associated fileKey blob.Reader.
//
// If preserveName is false then the returned File.Name will have
// a new randomly generated suffix, otherwise it will reuse the original one.
//
// This method could be useful in case you want to clone an existing
// Record file and assign it to a new Record (e.g. in a Record duplicate action).
//
// If you simply want to copy an existing file to a new location you
// could check the Copy(srcKey, dstKey) method.
func (s *System) GetReuploadableFile(fileKey string, preserveName bool) (*File, error) {
	attrs, err := s.Attributes(fileKey)
	if err != nil {
		return nil, err
	}

	name := path.Base(fileKey)
	originalName := attrs.Metadata[MetadataOriginalName]
	if originalName == "" {
		originalName = name
	}

	file := &File{}
	file.Size = attrs.Size
	file.OriginalName = originalName
	file.Reader = openFuncAsReader(func() (io.ReadSeekCloser, error) {
		return s.GetReader(fileKey)
	})

	if preserveName {
		file.Name = name
	} else {
		file.Name = normalizeName(file.Reader, originalName)
	}

	return file, nil
}

// Copy copies the file stored at srcKey to dstKey.
//
// If srcKey file doesn't exist, it returns ErrNotFound.
//
// If dstKey file already exists, it is overwritten.
func (s *System) Copy(srcKey, dstKey string) error {
	return s.bucket.Copy(s.ctx, dstKey, srcKey)
}

// List returns a flat list with info for all files under the specified prefix.
func (s *System) List(prefix string) ([]*blob.ListObject, error) {
	files := []*blob.ListObject{}

	iter := s.bucket.List(&blob.ListOptions{
		Prefix: prefix,
	})

	for {
		obj, err := iter.Next(s.ctx)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, err
			}
			break
		}
		files = append(files, obj)
	}

	return files, nil
}

// Upload writes content into the fileKey location.
func (s *System) Upload(content []byte, fileKey string) error {
	opts := &blob.WriterOptions{
		ContentType: mimetype.Detect(content).String(),
	}

	w, writerErr := s.NewWriter(fileKey, opts)
	if writerErr != nil {
		return writerErr
	}

	if _, err := w.Write(content); err != nil {
		return errors.Join(err, w.Close())
	}

	return w.Close()
}

// UploadFile uploads the provided File to the fileKey location.
func (s *System) UploadFile(file *File, fileKey string) error {
	f, err := file.Reader.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	mt, err := mimetype.DetectReader(f)
	if err != nil {
		return err
	}

	// rewind
	f.Seek(0, io.SeekStart)

	originalName := file.OriginalName
	if len(originalName) > 255 {
		// keep only the first 255 chars as a very rudimentary measure
		// to prevent the metadata to grow too big in size
		originalName = originalName[:255]
	}
	opts := &blob.WriterOptions{
		ContentType: mt.String(),
		Metadata: map[string]string{
			MetadataOriginalName: originalName,
		},
	}

	w, err := s.NewWriter(fileKey, opts)
	if err != nil {
		return err
	}

	if _, err := w.ReadFrom(f); err != nil {
		w.Close()
		return err
	}

	return w.Close()
}

// UploadMultipart uploads the provided multipart file to the fileKey location.
func (s *System) UploadMultipart(fh *multipart.FileHeader, fileKey string) error {
	f, err := fh.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	mt, err := mimetype.DetectReader(f)
	if err != nil {
		return err
	}

	// rewind
	f.Seek(0, io.SeekStart)

	originalName := fh.Filename
	if len(originalName) > 255 {
		// keep only the first 255 chars as a very rudimentary measure
		// to prevent the metadata to grow too big in size
		originalName = originalName[:255]
	}
	opts := &blob.WriterOptions{
		ContentType: mt.String(),
		Metadata: map[string]string{
			MetadataOriginalName: originalName,
		},
	}

	w, err := s.NewWriter(fileKey, opts)
	if err != nil {
		return err
	}

	_, err = w.ReadFrom(f)
	if err != nil {
		w.Close()
		return err
	}

	return w.Close()
}

// NewWriter returns a new blob.Writer instance allowing direct file
// create from an io.Reader value.
//
// If a file with the specified fileKey already exists, it will be replaced.
//
// NB! Make sure to call `Close()` on the resulting writer after you are done working with it.
//
// Note: If you have a bytes slice, [filesystem.File], or a multipart header value,
// you can check also the Upload* related methods as they are more user-friendly.
//
// Example:
//
//	content := strings.NewReader("Lorem ipsum dolor sit amet...")
//
//	fsys, _ := filesystem.NewLocal("dir")
//	defer fsys.Close()
//
//	w, _ := fsys.NewWriter("example/file/key", nil)
//	w.ReadFrom(content)
//	w.Close()
func (s *System) NewWriter(fileKey string, opts *blob.WriterOptions) (*blob.Writer, error) {
	if s.onNewWriter == nil {
		return s.bucket.NewWriter(s.ctx, fileKey, opts)
	}

	event := new(NewWriterEvent)
	event.Filesystem = s
	event.FileKey = fileKey
	event.Options = opts

	err := s.onNewWriter.Trigger(event, func(e *NewWriterEvent) error {
		writer, err := e.Filesystem.bucket.NewWriter(e.Filesystem.ctx, e.FileKey, e.Options)
		if err != nil {
			return err
		}

		e.Writer = writer

		return nil
	})
	if err != nil {
		return nil, err
	}

	return event.Writer, nil
}

// Delete deletes stored file at fileKey location.
//
// If the file doesn't exist returns ErrNotFound.
func (s *System) Delete(fileKey string) error {
	if s.onDelete == nil {
		return s.bucket.Delete(s.ctx, fileKey)
	}

	event := new(DeleteEvent)
	event.Filesystem = s
	event.FileKey = fileKey

	return s.onDelete.Trigger(event, func(e *DeleteEvent) error {
		return e.Filesystem.bucket.Delete(e.Filesystem.ctx, e.FileKey)
	})
}

// DeletePrefix deletes everything starting with the specified prefix.
//
// The prefix could be subpath (ex. "/a/b/") or filename prefix (ex. "/a/b/file_").
func (s *System) DeletePrefix(prefix string) []error {
	failed := []error{}

	if prefix == "" {
		failed = append(failed, errors.New("prefix mustn't be empty"))
		return failed
	}

	dirsMap := map[string]struct{}{}

	var isPrefixDir bool

	// treat the prefix as directory only if it ends with trailing slash
	if strings.HasSuffix(prefix, "/") {
		isPrefixDir = true
		dirsMap[strings.TrimRight(prefix, "/")] = struct{}{}
	}

	// delete all files with the prefix
	// ---
	iter := s.bucket.List(&blob.ListOptions{
		Prefix: prefix,
	})
	for {
		obj, err := iter.Next(s.ctx)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				failed = append(failed, err)
			}
			break
		}

		if err := s.Delete(obj.Key); err != nil {
			failed = append(failed, err)
		} else if isPrefixDir {
			slashIdx := strings.LastIndex(obj.Key, "/")
			if slashIdx > -1 {
				dirsMap[obj.Key[:slashIdx]] = struct{}{}
			}
		}
	}
	// ---

	// try to delete the empty remaining dir objects
	// (this operation usually is optional and there is no need to strictly check the result)
	// ---
	// fill dirs slice
	dirs := make([]string, 0, len(dirsMap))
	for d := range dirsMap {
		dirs = append(dirs, d)
	}

	// sort the child dirs first, aka. ["a/b/c", "a/b", "a"]
	sort.SliceStable(dirs, func(i, j int) bool {
		return len(strings.Split(dirs[i], "/")) > len(strings.Split(dirs[j], "/"))
	})

	// delete dirs
	for _, d := range dirs {
		if d != "" {
			s.Delete(d)
		}
	}
	// ---

	return failed
}

// Checks if the provided dir prefix doesn't have any files.
//
// A trailing slash will be appended to a non-empty dir string argument
// to ensure that the checked prefix is a "directory".
//
// Returns "false" in case it has at least one file, otherwise - "true".
func (s *System) IsEmptyDir(dir string) bool {
	if dir != "" && !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	iter := s.bucket.List(&blob.ListOptions{
		Prefix: dir,
	})

	_, err := iter.Next(s.ctx)

	return err != nil && errors.Is(err, io.EOF)
}

var inlineServeContentTypes = []string{
	// image
	"image/png", "image/jpg", "image/jpeg", "image/gif", "image/webp", "image/x-icon", "image/bmp",
	// video
	"video/webm", "video/mp4", "video/3gpp", "video/quicktime", "video/x-ms-wmv",
	// audio
	"audio/basic", "audio/aiff", "audio/mpeg", "audio/midi", "audio/mp3", "audio/wave",
	"audio/wav", "audio/x-wav", "audio/x-mpeg", "audio/x-m4a", "audio/aac",
	// document
	"application/pdf", "application/x-pdf",
}

// manualExtensionContentTypes is a map of file extensions to content types.
var manualExtensionContentTypes = map[string]string{
	// https://github.com/whatwg/mimesniff/issues/7
	".svg": "image/svg+xml",

	// https://github.com/gabriel-vasile/mimetype/pull/113
	".css": "text/css",

	// https://github.com/pocketbase/pocketbase/issues/6597
	".js":  "text/javascript",
	".mjs": "text/javascript",

	// https://github.com/pocketbase/pocketbase/discussions/7467
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
}

// forceAttachmentParam is the name of the request query parameter to
// force "Content-Disposition: attachment" header.
const forceAttachmentParam = "download"

// Serve serves the file at fileKey location to an HTTP response.
//
// If the `download` query parameter is used the file will be always served for
// download no matter of its type (aka. with "Content-Disposition: attachment").
//
// Internally this method uses [http.ServeContent] so Range requests,
// If-Match, If-Unmodified-Since, etc. headers are handled transparently.
func (s *System) Serve(res http.ResponseWriter, req *http.Request, fileKey string, name string) error {
	br, readErr := s.GetReader(fileKey)
	if readErr != nil {
		return readErr
	}
	defer br.Close()

	var forceAttachment bool
	if raw := req.URL.Query().Get(forceAttachmentParam); raw != "" {
		forceAttachment, _ = strconv.ParseBool(raw)
	}

	disposition := "attachment"
	realContentType := br.ContentType()
	if !forceAttachment && list.ExistInSlice(realContentType, inlineServeContentTypes) {
		disposition = "inline"
	}

	// make an exception for specific content types and force a custom
	// content type to send in the response so that it can be loaded properly
	extContentType := realContentType
	if ct, found := manualExtensionContentTypes[filepath.Ext(fileKey)]; found {
		extContentType = ct
	}

	setHeaderIfMissing(res, "Content-Disposition", disposition+"; filename="+strconv.Quote(name))
	setHeaderIfMissing(res, "Content-Type", extContentType)
	setHeaderIfMissing(res, "Content-Security-Policy", "default-src 'none'; media-src 'self'; style-src 'unsafe-inline'; sandbox")

	// set a default cache-control header
	// (valid for 30 days but the cache is allowed to reuse the file for any requests
	// that are made in the last day while revalidating the res in the background)
	setHeaderIfMissing(res, "Cache-Control", "max-age=2592000, stale-while-revalidate=86400")

	http.ServeContent(res, req, name, br.ModTime(), br)

	return nil
}

// note: expects key to be in a canonical form (eg. "accept-encoding" should be "Accept-Encoding").
func setHeaderIfMissing(res http.ResponseWriter, key string, value string) {
	if _, ok := res.Header()[key]; !ok {
		res.Header().Set(key, value)
	}
}

var ThumbSizeRegex = regexp.MustCompile(`^(\d+)x(\d+)(t|b|f)?$`)

// CreateThumb creates a new thumb image for the file at originalKey location.
// The new thumb file is stored at thumbKey location.
//
// thumbSize is in the format:
// - 0xH  (eg. 0x100)    - resize to H height preserving the aspect ratio
// - Wx0  (eg. 300x0)    - resize to W width preserving the aspect ratio
// - WxH  (eg. 300x100)  - resize and crop to WxH viewbox (from center)
// - WxHt (eg. 300x100t) - resize and crop to WxH viewbox (from top)
// - WxHb (eg. 300x100b) - resize and crop to WxH viewbox (from bottom)
// - WxHf (eg. 300x100f) - fit inside a WxH viewbox (without cropping)
func (s *System) CreateThumb(originalKey, thumbKey, thumbSize string) error {
	// note: the wrapping is an extra precaution since there were several
	// golang.org/x/image panic related issues over the years
	return routine.SafeWrap(func() error {
		return s.createThumb(originalKey, thumbKey, thumbSize)
	})()
}

func (s *System) createThumb(originalKey, thumbKey, thumbSize string) error {
	sizeParts := ThumbSizeRegex.FindStringSubmatch(thumbSize)
	if len(sizeParts) != 4 {
		return errors.New("thumb size must be in WxH, WxHt, WxHb or WxHf format")
	}

	width, _ := strconv.Atoi(sizeParts[1])
	height, _ := strconv.Atoi(sizeParts[2])
	resizeType := sizeParts[3]

	if width == 0 && height == 0 {
		return errors.New("thumb width and height cannot be zero at the same time")
	}

	// fetch the original
	r, readErr := s.GetReader(originalKey)
	if readErr != nil {
		return readErr
	}
	defer r.Close()

	// create imaging object from the original reader
	// (note: only the first frame for animated image formats)
	img, decodeErr := imaging.Decode(r, imaging.AutoOrientation(true))
	if decodeErr != nil {
		return decodeErr
	}

	var thumbImg *image.NRGBA

	if width == 0 || height == 0 {
		// force resize preserving aspect ratio
		thumbImg = imaging.Resize(img, width, height, imaging.Linear)
	} else {
		switch resizeType {
		case "f":
			// fit
			thumbImg = imaging.Fit(img, width, height, imaging.Linear)
		case "t":
			// fill and crop from top
			thumbImg = imaging.Fill(img, width, height, imaging.Top, imaging.Linear)
		case "b":
			// fill and crop from bottom
			thumbImg = imaging.Fill(img, width, height, imaging.Bottom, imaging.Linear)
		default:
			// fill and crop from center
			thumbImg = imaging.Fill(img, width, height, imaging.Center, imaging.Linear)
		}
	}

	originalContentType := r.ContentType()

	opts := &blob.WriterOptions{
		ContentType: originalContentType,
	}

	var format imaging.Format

	switch originalContentType {
	case "image/jpeg":
		format = imaging.JPEG
	case "image/gif":
		format = imaging.GIF
	case "image/tiff":
		format = imaging.TIFF
	case "image/bmp":
		format = imaging.BMP
	default:
		// fallback to PNG (this includes webp!)
		opts.ContentType = "image/png"
		format = imaging.PNG
	}

	// open a thumb storage writer (aka. prepare for upload)
	w, err := s.NewWriter(thumbKey, opts)
	if err != nil {
		return err
	}

	// thumb encode (aka. upload)
	err = imaging.Encode(w, thumbImg, format)
	if err != nil {
		w.Close()
		return err
	}

	// check for close errors to ensure that the thumb was really saved
	return w.Close()
}
