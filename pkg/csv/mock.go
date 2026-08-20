package csv

import (
	"bytes"
	"io"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// CreateTestMultipartRequest helps create a multipart from request with a CSV file
func CreateTestMultipartRequest(t *testing.T, content string) (*multipart.Writer, *bytes.Buffer) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.csv")
	assert.NoError(t, err)

	_, err = io.Copy(part, strings.NewReader(content))
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	return writer, body
}
