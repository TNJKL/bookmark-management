package service

import (
	"errors"
	"testing"

	repoMocks "github.com/TNJKL/bookmark-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHealthCheck_HealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		//test case name
		name string

		//input
		serviceName string
		instanceID  string
		setupMock   func(m *repoMocks.HealthRepository)

		//expected output
		expectedMessage     string
		expectedServiceName string
		expectedInstanceID  string
		expectedError       error
	}{
		{
			name:        "happy path",
			serviceName: "bookmark_service",
			instanceID:  "67026e45-34fd-449c-aa29-d18c7686ab00",
			setupMock: func(m *repoMocks.HealthRepository) {
				m.On("Ping", mock.Anything).Return(nil)
			},
			expectedMessage:     "OK",
			expectedServiceName: "bookmark_service",
			expectedInstanceID:  "67026e45-34fd-449c-aa29-d18c7686ab00",
			expectedError:       nil,
		},
		{
			name:        "empty serviceName case",
			serviceName: "",
			instanceID:  "67026e45-34fd-449c-aa29-d18c7686ab00",
			setupMock: func(m *repoMocks.HealthRepository) {
				m.On("Ping", mock.Anything).Return(nil)
			},
			expectedMessage:     "OK",
			expectedServiceName: "",
			expectedInstanceID:  "67026e45-34fd-449c-aa29-d18c7686ab00",
			expectedError:       nil,
		},
		{
			name:        "redis connection error",
			serviceName: "bookmark_service",
			instanceID:  "67026e45-34fd-449c-aa29-d18c7686ab00",
			setupMock: func(m *repoMocks.HealthRepository) {
				m.On("Ping", mock.Anything).Return(errors.New("redis connection refused"))
			},
			expectedMessage:     "",
			expectedServiceName: "",
			expectedInstanceID:  "",
			expectedError:       errors.New("redis connection refused"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//goi Parallel de chay song song cac testcase
			t.Parallel()
			mockRepo := repoMocks.NewHealthRepository(t)
			tc.setupMock(mockRepo)
			testSvc := NewHealthCheck(tc.serviceName, tc.instanceID, mockRepo)
			result, err := testSvc.HealthCheck(t.Context())
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedMessage, result.Message)
				assert.Equal(t, tc.expectedServiceName, result.ServiceName)
				assert.Equal(t, tc.expectedInstanceID, result.InstanceID)
			}
		})
	}

}
