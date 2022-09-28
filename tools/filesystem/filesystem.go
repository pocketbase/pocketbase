package filesystem

import (
	"context"
	"errors"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"github.com/pocketbase/pocketbase/tools/list"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/blob/s3blob"
)

type System struct {
	ctx    context.Context
	bucket *blob.Bucket
}

// NewS3 initializes an S3 filesystem instance.
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

	cred := credentials.NewStaticCredentials(accessKey, secretKey, "")

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		Credentials:      cred,
		S3ForcePathStyle: aws.Bool(s3ForcePathStyle),
	})
	if err != nil {
		return nil, err
	}

	bucket, err := s3blob.OpenBucket(ctx, sess, bucketName, nil)
	if err != nil {
		return nil, err
	}

	return &System{ctx: ctx, bucket: bucket}, nil
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

	bucket, err := fileblob.OpenBucket(dirPath, nil)
	if err != nil {
		return nil, err
	}

	return &System{ctx: ctx, bucket: bucket}, nil
}

// Close releases any resources used for the related filesystem.
func (s *System) Close() error {
	return s.bucket.Close()
}

// Exists checks if file with fileKey path exists or not.
func (s *System) Exists(fileKey string) (bool, error) {
	return s.bucket.Exists(s.ctx, fileKey)
}

// Attributes returns the attributes for the file with fileKey path.
func (s *System) Attributes(fileKey string) (*blob.Attributes, error) {
	return s.bucket.Attributes(s.ctx, fileKey)
}

// Upload writes content into the fileKey location.
func (s *System) Upload(content []byte, fileKey string) error {
	opts := &blob.WriterOptions{
		ContentType: mimetype.Detect(content).String(),
	}

	w, writerErr := s.bucket.NewWriter(s.ctx, fileKey, opts)
	if writerErr != nil {
		return writerErr
	}

	if _, err := w.Write(content); err != nil {
		w.Close()
		return err
	}

	return w.Close()
}

// Delete deletes stored file at fileKey location.
func (s *System) Delete(fileKey string) error {
	return s.bucket.Delete(s.ctx, fileKey)
}

// DeletePrefix deletes everything starting with the specified prefix.
func (s *System) DeletePrefix(prefix string) []error {
	failed := []error{}

	if prefix == "" {
		failed = append(failed, errors.New("Prefix mustn't be empty."))
		return failed
	}

	dirsMap := map[string]struct{}{}
	dirsMap[prefix] = struct{}{}

	opts := blob.ListOptions{
		Prefix: prefix,
	}

	// delete all files with the prefix
	// ---
	iter := s.bucket.List(&opts)
	for {
		obj, err := iter.Next(s.ctx)
		if err == io.EOF {
			break
		}

		if err != nil {
			failed = append(failed, err)
			continue
		}

		if err := s.Delete(obj.Key); err != nil {
			failed = append(failed, err)
		} else {
			dirsMap[filepath.Dir(obj.Key)] = struct{}{}
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
	".svg": "image/svg+xml", // (see https://github.com/whatwg/mimesniff/issues/7)
	".css": "text/css",      // (see https://github.com/gabriel-vasile/mimetype/pull/113)
}

// Serve serves the file at fileKey location to an HTTP response.
func (s *System) Serve(response http.ResponseWriter, fileKey string, name string) error {
	r, readErr := s.bucket.NewReader(s.ctx, fileKey, nil)
	if readErr != nil {
		return readErr
	}
	defer r.Close()

	disposition := "attachment"
	realContentType := r.ContentType()
	if list.ExistInSlice(realContentType, inlineServeContentTypes) {
		disposition = "inline"
	}

	// make an exception for specific content types and force a
	// custom content type to send in the response so that it can be loaded directly.
	extContentType := realContentType
	if ct, found := manualExtensionContentTypes[filepath.Ext(name)]; found && extContentType != ct {
		extContentType = ct
	}

	response.Header().Set("Content-Disposition", disposition+"; filename="+name)
	response.Header().Set("Content-Type", extContentType)
	response.Header().Set("Content-Length", strconv.FormatInt(r.Size(), 10))
	response.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; sandbox")

	// All HTTP date/time stamps MUST be represented in Greenwich Mean Time (GMT)
	// (see https://www.w3.org/Protocols/rfc2616/rfc2616-sec3.html#sec3.3.1)
	//
	// NB! time.LoadLocation may fail on non-Unix systems (see https://github.com/pocketbase/pocketbase/issues/45)
	location, locationErr := time.LoadLocation("GMT")
	if locationErr == nil {
		response.Header().Set("Last-Modified", r.ModTime().In(location).Format("Mon, 02 Jan 06 15:04:05 MST"))
	}

	// copy from the read range to response.
	_, err := io.Copy(response, r)

	return err
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
func (s *System) CreateThumb(originalKey string, thumbKey, thumbSize string) error {
	sizeParts := ThumbSizeRegex.FindStringSubmatch(thumbSize)
	if len(sizeParts) != 4 {
		return errors.New("Thumb size must be in WxH, WxHt, WxHb or WxHf format.")
	}

	width, _ := strconv.Atoi(sizeParts[1])
	height, _ := strconv.Atoi(sizeParts[2])
	resizeType := sizeParts[3]

	if width == 0 && height == 0 {
		return errors.New("Thumb width and height cannot be zero at the same time.")
	}

	// fetch the original
	r, readErr := s.bucket.NewReader(s.ctx, originalKey, nil)
	if readErr != nil {
		return readErr
	}
	defer r.Close()

	// create imaging object from the original reader
	img, decodeErr := imaging.Decode(r, imaging.AutoOrientation(true))
	if decodeErr != nil {
		return decodeErr
	}

	var thumbImg *image.NRGBA

	if width == 0 || height == 0 {
		// force resize preserving aspect ratio
		thumbImg = imaging.Resize(img, width, height, imaging.CatmullRom)
	} else {
		switch resizeType {
		case "f":
			// fit
			thumbImg = imaging.Fit(img, width, height, imaging.CatmullRom)
		case "t":
			// fill and crop from top
			thumbImg = imaging.Fill(img, width, height, imaging.Top, imaging.CatmullRom)
		case "b":
			// fill and crop from bottom
			thumbImg = imaging.Fill(img, width, height, imaging.Bottom, imaging.CatmullRom)
		default:
			// fill and crop from center
			thumbImg = imaging.Fill(img, width, height, imaging.Center, imaging.CatmullRom)
		}
	}

	// open a thumb storage writer (aka. prepare for upload)
	w, writerErr := s.bucket.NewWriter(s.ctx, thumbKey, nil)
	if writerErr != nil {
		return writerErr
	}

	// try to detect the thumb format based on the original file name
	// (fallbacks to png on error)
	format, err := imaging.FormatFromFilename(thumbKey)
	if err != nil {
		format = imaging.PNG
	}

	// thumb encode (aka. upload)
	if err := imaging.Encode(w, thumbImg, format); err != nil {
		w.Close()
		return err
	}

	// check for close errors to ensure that the thumb was really saved
	return w.Close()
}
