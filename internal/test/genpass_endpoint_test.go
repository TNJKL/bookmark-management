package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestGenPassEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "success case",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("GET", "/genpass", nil)
				respRecorder := httptest.NewRecorder()

				api.ServerHTTP(respRecorder, req)
				return respRecorder
			},

			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"password"`,
		},

		{
			name: "wrong endpoint method",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("POST", "/genpass", nil)
				respRecorder := httptest.NewRecorder()

				api.ServerHTTP(respRecorder, req)
				return respRecorder
			},

			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testAPI := api.NewEngine(&api.Config{}, nil, nil)
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)

		})
	}

}
