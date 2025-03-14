package s3

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrUsedUploader = errors.New("the Uploader has been already used")

const (
	defaultMaxConcurrency int = 5
	defaultMinPartSize    int = 6 << 20
)

// Uploader handles the upload of a single S3 object.
//
// If the Payload size is less than the configured MinPartSize it sends
// a single (PutObject) request, otherwise performs chunked/multipart upload.
type Uploader struct {
	// S3 is the S3 client instance performing the upload object request (required).
	S3 *S3

	// Payload is the object content to upload (required).
	Payload io.Reader

	// Key is the destination key of the uploaded object (required).
	Key string

	// Metadata specifies the optional metadata to write with the object upload.
	Metadata map[string]string

	// MaxConcurrency specifies the max number of workers to use when
	// performing chunked/multipart upload.
	//
	// If zero or negative, defaults to 5.
	//
	// This option is used only when the Payload size is > MinPartSize.
	MaxConcurrency int

	// MinPartSize specifies the min Payload size required to perform
	// chunked/multipart upload.
	//
	// If zero or negative, defaults to ~6MB.
	MinPartSize int

	uploadId       string
	uploadedParts  []*mpPart
	lastPartNumber int
	mu             sync.Mutex // guards lastPartNumber and the uploadedParts slice
	used           bool
}

// Upload processes the current Uploader instance.
//
// Users can specify an optional optReqFuncs that will be passed down to all Upload internal requests
// (single upload, multipart init, multipart parts upload, multipart complete, multipart abort).
//
// Note that after this call the Uploader should be discarded (aka. no longer can be used).
func (u *Uploader) Upload(ctx context.Context, optReqFuncs ...func(*http.Request)) error {
	if u.used {
		return ErrUsedUploader
	}

	err := u.validateAndNormalize()
	if err != nil {
		return err
	}

	initPart, _, err := u.nextPart()
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	if len(initPart) < u.MinPartSize {
		return u.singleUpload(ctx, initPart, optReqFuncs...)
	}

	err = u.multipartInit(ctx, optReqFuncs...)
	if err != nil {
		return fmt.Errorf("multipart init error: %w", err)
	}

	err = u.multipartUpload(ctx, initPart, optReqFuncs...)
	if err != nil {
		return errors.Join(
			u.multipartAbort(ctx, optReqFuncs...),
			fmt.Errorf("multipart upload error: %w", err),
		)
	}

	err = u.multipartComplete(ctx, optReqFuncs...)
	if err != nil {
		return errors.Join(
			u.multipartAbort(ctx, optReqFuncs...),
			fmt.Errorf("multipart complete error: %w", err),
		)
	}

	return nil
}

// -------------------------------------------------------------------

func (u *Uploader) validateAndNormalize() error {
	if u.S3 == nil {
		return errors.New("Uploader.S3 must be a non-empty and properly initialized S3 client instance")
	}

	if u.Key == "" {
		return errors.New("Uploader.Key is required")
	}

	if u.Payload == nil {
		return errors.New("Uploader.Payload must be non-nill")
	}

	if u.MaxConcurrency <= 0 {
		u.MaxConcurrency = defaultMaxConcurrency
	}

	if u.MinPartSize <= 0 {
		u.MinPartSize = defaultMinPartSize
	}

	return nil
}

func (u *Uploader) singleUpload(ctx context.Context, part []byte, optReqFuncs ...func(*http.Request)) error {
	if u.used {
		return ErrUsedUploader
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.S3.URL(u.Key), bytes.NewReader(part))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Length", strconv.Itoa(len(part)))

	for k, v := range u.Metadata {
		req.Header.Set(metadataPrefix+k, v)
	}

	// apply optional request funcs
	for _, fn := range optReqFuncs {
		if fn != nil {
			fn(req)
		}
	}

	resp, err := u.S3.SignAndSend(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// -------------------------------------------------------------------

type mpPart struct {
	XMLName    xml.Name `xml:"Part"`
	ETag       string   `xml:"ETag"`
	PartNumber int      `xml:"PartNumber"`
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateMultipartUpload.html
func (u *Uploader) multipartInit(ctx context.Context, optReqFuncs ...func(*http.Request)) error {
	if u.used {
		return ErrUsedUploader
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.S3.URL(u.Key+"?uploads"), nil)
	if err != nil {
		return err
	}

	for k, v := range u.Metadata {
		req.Header.Set(metadataPrefix+k, v)
	}

	// apply optional request funcs
	for _, fn := range optReqFuncs {
		if fn != nil {
			fn(req)
		}
	}

	resp, err := u.S3.SignAndSend(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body := &struct {
		XMLName  xml.Name `xml:"InitiateMultipartUploadResult"`
		UploadId string   `xml:"UploadId"`
	}{}

	err = xml.NewDecoder(resp.Body).Decode(body)
	if err != nil {
		return err
	}

	u.uploadId = body.UploadId

	return nil
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_AbortMultipartUpload.html
func (u *Uploader) multipartAbort(ctx context.Context, optReqFuncs ...func(*http.Request)) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.used = true

	// ensure that the specified abort context is always valid to allow cleanup
	var abortCtx = ctx
	if abortCtx.Err() != nil {
		abortCtx = context.Background()
	}

	query := url.Values{"uploadId": []string{u.uploadId}}

	req, err := http.NewRequestWithContext(abortCtx, http.MethodDelete, u.S3.URL(u.Key+"?"+query.Encode()), nil)
	if err != nil {
		return err
	}

	// apply optional request funcs
	for _, fn := range optReqFuncs {
		if fn != nil {
			fn(req)
		}
	}

	resp, err := u.S3.SignAndSend(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_CompleteMultipartUpload.html
func (u *Uploader) multipartComplete(ctx context.Context, optReqFuncs ...func(*http.Request)) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.used = true

	// the list of parts must be sorted in ascending order
	slices.SortFunc(u.uploadedParts, func(a, b *mpPart) int {
		if a.PartNumber < b.PartNumber {
			return -1
		}
		if a.PartNumber > b.PartNumber {
			return 1
		}
		return 0
	})

	// build a request payload with the uploaded parts
	xmlParts := &struct {
		XMLName xml.Name `xml:"CompleteMultipartUpload"`
		Parts   []*mpPart
	}{
		Parts: u.uploadedParts,
	}
	rawXMLParts, err := xml.Marshal(xmlParts)
	if err != nil {
		return err
	}
	reqPayload := strings.NewReader(xml.Header + string(rawXMLParts))

	query := url.Values{"uploadId": []string{u.uploadId}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.S3.URL(u.Key+"?"+query.Encode()), reqPayload)
	if err != nil {
		return err
	}

	// apply optional request funcs
	for _, fn := range optReqFuncs {
		if fn != nil {
			fn(req)
		}
	}

	resp, err := u.S3.SignAndSend(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (u *Uploader) nextPart() ([]byte, int, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	part := make([]byte, u.MinPartSize)
	n, err := io.ReadFull(u.Payload, part)

	// normalize io.EOF errors and ensure that io.EOF is returned only when there were no read bytes
	if err != nil && (errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF)) {
		if n == 0 {
			err = io.EOF
		} else {
			err = nil
		}
	}

	u.lastPartNumber++

	return part[0:n], u.lastPartNumber, err
}

func (u *Uploader) multipartUpload(ctx context.Context, initPart []byte, optReqFuncs ...func(*http.Request)) error {
	var g errgroup.Group
	g.SetLimit(u.MaxConcurrency)

	totalParallel := u.MaxConcurrency

	if len(initPart) != 0 {
		totalParallel--
		initPartNumber := u.lastPartNumber
		g.Go(func() error {
			mp, err := u.uploadPart(ctx, initPartNumber, initPart, optReqFuncs...)
			if err != nil {
				return err
			}

			u.mu.Lock()
			u.uploadedParts = append(u.uploadedParts, mp)
			u.mu.Unlock()

			return nil
		})
	}

	for i := 0; i < totalParallel; i++ {
		g.Go(func() error {
			for {
				part, num, err := u.nextPart()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					return err
				}

				mp, err := u.uploadPart(ctx, num, part, optReqFuncs...)
				if err != nil {
					return err
				}

				u.mu.Lock()
				u.uploadedParts = append(u.uploadedParts, mp)
				u.mu.Unlock()
			}

			return nil
		})
	}

	return g.Wait()
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_UploadPart.html
func (u *Uploader) uploadPart(ctx context.Context, partNumber int, partData []byte, optReqFuncs ...func(*http.Request)) (*mpPart, error) {
	query := url.Values{}
	query.Set("uploadId", u.uploadId)
	query.Set("partNumber", strconv.Itoa(partNumber))

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.S3.URL(u.Key+"?"+query.Encode()), bytes.NewReader(partData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Length", strconv.Itoa(len(partData)))

	// apply optional request funcs
	for _, fn := range optReqFuncs {
		if fn != nil {
			fn(req)
		}
	}

	resp, err := u.S3.SignAndSend(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &mpPart{
		PartNumber: partNumber,
		ETag:       resp.Header.Get("ETag"),
	}, nil
}
