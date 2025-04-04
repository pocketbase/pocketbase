package jsvm

import (
	"bytes"
	"io"
	"mime/multipart"

	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/spf13/cast"
)

// FormData represents an interface similar to the browser's [FormData].
//
// The value of each FormData entry must be a string or [*filesystem.File] instance.
//
// It is intended to be used together with the JSVM `$http.send` when
// sending multipart/form-data requests.
//
// [FormData]: https://developer.mozilla.org/en-US/docs/Web/API/FormData.
type FormData map[string][]any

// Append appends a new value onto an existing key inside the current FormData,
// or adds the key if it does not already exist.
func (data FormData) Append(key string, value any) {
	data[key] = append(data[key], value)
}

// Set sets a new value for an existing key inside the current FormData,
// or adds the key/value if it does not already exist.
func (data FormData) Set(key string, value any) {
	data[key] = []any{value}
}

// Delete deletes a key and its value(s) from the current FormData.
func (data FormData) Delete(key string) {
	delete(data, key)
}

// Get returns the first value associated with a given key from
// within the current FormData.
//
// If you expect multiple values and want all of them,
// use the [FormData.GetAll] method instead.
func (data FormData) Get(key string) any {
	values, ok := data[key]
	if !ok || len(values) == 0 {
		return nil
	}

	return values[0]
}

// GetAll returns all the values associated with a given key
// from within the current FormData.
func (data FormData) GetAll(key string) []any {
	values, ok := data[key]
	if !ok {
		return nil
	}

	return values
}

// Has returns whether a FormData object contains a certain key.
func (data FormData) Has(key string) bool {
	values, ok := data[key]

	return ok && len(values) > 0
}

// Keys returns all keys contained in the current FormData.
func (data FormData) Keys() []string {
	result := make([]string, 0, len(data))

	for k := range data {
		result = append(result, k)
	}

	return result
}

// Keys returns all values contained in the current FormData.
func (data FormData) Values() []any {
	result := make([]any, 0, len(data))

	for _, values := range data {
		result = append(result, values...)
	}

	return result
}

// Entries returns a [key, value] slice pair for each FormData entry.
func (data FormData) Entries() [][]any {
	result := make([][]any, 0, len(data))

	for k, values := range data {
		for _, v := range values {
			result = append(result, []any{k, v})
		}
	}

	return result
}

// toMultipart converts the current FormData entries into multipart encoded data.
func (data FormData) toMultipart() (*bytes.Buffer, *multipart.Writer, error) {
	body := new(bytes.Buffer)

	mp := multipart.NewWriter(body)
	defer mp.Close()

	for k, values := range data {
		for _, rawValue := range values {
			switch v := rawValue.(type) {
			case *filesystem.File:
				err := func() error {
					mpw, err := mp.CreateFormFile(k, v.OriginalName)
					if err != nil {
						return err
					}

					file, err := v.Reader.Open()
					if err != nil {
						return err
					}
					defer file.Close()

					_, err = io.Copy(mpw, file)
					if err != nil {
						return err
					}

					return nil
				}()
				if err != nil {
					return nil, nil, err
				}
			default:
				err := mp.WriteField(k, cast.ToString(v))
				if err != nil {
					return nil, nil, err
				}
			}
		}
	}

	return body, mp, nil
}
