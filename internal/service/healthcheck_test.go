package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_HealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		//test case name
		name string

		//input
		serviceName string
		instanceID  string

		//expected output
		expectedMessage     string
		expectedServiceName string
		expectedInstanceID  string
	}{
		{
			name:        "happy path",
			serviceName: "bookmark_service",
			instanceID:  "67026e45-34fd-449c-aa29-d18c7686ab00",

			expectedMessage:     "OK",
			expectedServiceName: "bookmark_service",
			expectedInstanceID:  "67026e45-34fd-449c-aa29-d18c7686ab00",
		},
		{
			name:        "empty serviceName case",
			serviceName: "",
			instanceID:  "67026e45-34fd-449c-aa29-d18c7686ab00",

			expectedMessage:     "OK",
			expectedServiceName: "",
			expectedInstanceID:  "67026e45-34fd-449c-aa29-d18c7686ab00",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//goi Parallel de chay song song cac testcase
			t.Parallel()
			testSvc := NewHealthCheck(tc.serviceName, tc.instanceID)
			result, err := testSvc.HealthCheck()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedMessage, result.Message)
			assert.Equal(t, tc.expectedServiceName, result.ServiceName)
			assert.Equal(t, tc.expectedInstanceID, result.InstanceID)
		})
	}

}
