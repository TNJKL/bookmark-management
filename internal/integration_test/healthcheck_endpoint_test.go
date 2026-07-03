package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckEndpoint(t *testing.T) {
	//gọi Parallel để chạy song song các test case
	t.Parallel()

	testCases := []struct {
		name          string
		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "happpy path",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("GET", "/health-check", nil)
				responseRecorder := httptest.NewRecorder()

				api.ServerHTTP(responseRecorder, req)
				return responseRecorder
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"OK","service_name":"service_name_test","instance_id":"instance_test_id"}`,
		},
		{
			name: "wrong endpoint method",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("POST", "/health-check", nil)
				responseRecorder := httptest.NewRecorder()

				api.ServerHTTP(responseRecorder, req)
				return responseRecorder
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//gọi Parallel để chạy song song các test case
			t.Parallel()

			testAPI := api.NewEngine(&api.Config{ServiceName: "service_name_test", InstanceID: "instance_test_id"})
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)

		})
	}

}
