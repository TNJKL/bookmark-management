package csv_test

import (
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/pkg/csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCSVItem struct {
	Name  string `csv:"name"`
	Value string `csv:"value"`
}

func getFileHeaderFromContent(t *testing.T, content string) *multipart.FileHeader {
	t.Helper()
	writer, body := csv.CreateTestMultipartRequest(t, content)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	_, header, err := req.FormFile("file")
	require.NoError(t, err)
	return header
}

func TestParseFromMultipartFile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		csvContent  string
		expected    []*testCSVItem
		expectedErr bool
	}{
		{
			name:       "happy path - parse valid csv into struct slice",
			csvContent: "name,value\nitem1,val1\nitem2,val2",
			expected: []*testCSVItem{
				{Name: "item1", Value: "val1"},
				{Name: "item2", Value: "val2"},
			},
			expectedErr: false,
		},
		{
			name:        "happy path - empty csv with headers only",
			csvContent:  "name,value\n",
			expected:    nil,
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fileHeader := getFileHeaderFromContent(t, tc.csvContent)

			var result []*testCSVItem
			err := csv.ParseFromMultipartFile(fileHeader, &result)

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
