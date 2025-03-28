package s3

import (
	"encoding/xml"
	"strconv"
	"strings"
)

var _ error = (*ResponseError)(nil)

// ResponseError defines a general S3 response error.
//
// https://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html
type ResponseError struct {
	XMLName   xml.Name `json:"-" xml:"Error"`
	Code      string   `json:"code" xml:"Code"`
	Message   string   `json:"message" xml:"Message"`
	RequestId string   `json:"requestId" xml:"RequestId"`
	Resource  string   `json:"resource" xml:"Resource"`
	Raw       []byte   `json:"-" xml:"-"`
	Status    int      `json:"status" xml:"Status"`
}

// Error implements the std error interface.
func (err *ResponseError) Error() string {
	var strBuilder strings.Builder

	strBuilder.WriteString(strconv.Itoa(err.Status))
	strBuilder.WriteString(" ")

	if err.Code != "" {
		strBuilder.WriteString(err.Code)
	} else {
		strBuilder.WriteString("S3ResponseError")
	}

	if err.Message != "" {
		strBuilder.WriteString(": ")
		strBuilder.WriteString(err.Message)
	}

	if len(err.Raw) > 0 {
		strBuilder.WriteString("\n(RAW: ")
		strBuilder.Write(err.Raw)
		strBuilder.WriteString(")")
	}

	return strBuilder.String()
}
