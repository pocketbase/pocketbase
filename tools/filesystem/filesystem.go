package filesystem

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/disintegration/imaging"
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
) (*System, error) {
	ctx := context.Background() // default context

	cred := credentials.NewStaticCredentials(accessKey, secretKey, "")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Endpoint:    aws.String(endpoint),
		Credentials: cred,
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
	w, writerErr := s.bucket.NewWriter(s.ctx, fileKey, nil)
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
	dirs := []string{}
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

// Serve serves the file at fileKey location to an HTTP response.
func (s *System) Serve(response http.ResponseWriter, fileKey string, name string) error {
	r, readErr := s.bucket.NewReader(s.ctx, fileKey, nil)
	if readErr != nil {
		return readErr
	}
	defer r.Close()

	response.Header().Set("Content-Disposition", "attachment; filename="+name)
	response.Header().Set("Content-Type", r.ContentType())
	response.Header().Set("Content-Length", strconv.FormatInt(r.Size(), 10))

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

// CreateThumb creates a new thumb image for the file at originalKey location.
// The new thumb file is stored at thumbKey location.
//
// thumbSize is in the format "WxH", eg. "100x50".
func (s *System) CreateThumb(originalKey string, thumbKey, thumbSize string, cropCenter bool) error {
	thumbSizeParts := strings.SplitN(thumbSize, "x", 2)
	if len(thumbSizeParts) != 2 {
		return errors.New("Thumb size must be in WxH format.")
	}

	width, _ := strconv.Atoi(thumbSizeParts[0])
	height, _ := strconv.Atoi(thumbSizeParts[1])

	// fetch the original
	r, readErr := s.bucket.NewReader(s.ctx, originalKey, nil)
	if readErr != nil {
		return readErr
	}
	defer r.Close()

	// create imaging object from the origial reader
	img, decodeErr := imaging.Decode(r)
	if decodeErr != nil {
		return decodeErr
	}

	// determine crop anchor
	cropAnchor := imaging.Center
	if !cropCenter {
		cropAnchor = imaging.Top
	}

	// create thumb imaging object
	thumbImg := imaging.Fill(img, width, height, cropAnchor, imaging.CatmullRom)

	// open a thumb storage writer (aka. prepare for upload)
	w, writerErr := s.bucket.NewWriter(s.ctx, thumbKey, nil)
	if writerErr != nil {
		return writerErr
	}

	// thumb encode (aka. upload)
	if err := imaging.Encode(w, thumbImg, imaging.PNG); err != nil {
		w.Close()

		return err
	}

	// check for close errors to ensure that the thumb was really saved
	return w.Close()
}
