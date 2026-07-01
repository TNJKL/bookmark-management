package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenPass_GeneratePassword(t *testing.T) {
	//phai co Parallel() de chay test song song
	t.Parallel()

	testCases := []struct {
		name           string // Tên test case
		expectedLength int    // Input
		expectedError  error  // Expected output
	}{
		{
			name:           "normal case - length 12",
			expectedLength: 12,
			expectedError:  nil,
		},
		{
			name:           "normal case - length 16",
			expectedLength: 16,
			expectedError:  nil,
		},
		{
			name:           "normal case - length 10000",
			expectedLength: 10000,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//van phai goi Parallel o day de chay song song testcases
			t.Parallel()            // Chạy song song
			testSvc := NewGenPass() // Tạo service mới cho mỗi test
			password, err := testSvc.GeneratePassword(tc.expectedLength)
			assert.Equal(t, tc.expectedError, err)            // Kiểm tra error
			assert.Equal(t, tc.expectedLength, len(password)) // Kiểm tra kết quả
		})
	}

}
