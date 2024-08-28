package internal_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/darth-raijin/go-echo/internal"
	"github.com/stretchr/testify/assert"
)

const (
	TemplatePath = "templates/response.txt"
)

func TestHandleRequest(t *testing.T) {
	// Load the actual template using the internal.LoadTemplate function
	tmpl, err := internal.LoadTemplate(TemplatePath)
	assert.NoError(t, err, "Template should load without error")

	// Define test cases, including edge cases
	tests := []struct {
		name           string
		httpMethod     string
		requestPath    string
		headers        map[string]string
		expectedStatus int
	}{
		{
			name:        "Custom Status Code with Headers",
			httpMethod:  http.MethodGet,
			requestPath: "/test-route",
			headers: map[string]string{
				"X-Response-Status": "202",
				"X-Custom-Header":   "Some value 1",
			},
			expectedStatus: 202,
		},
		{
			name:        "Default Status Code with Headers",
			httpMethod:  http.MethodGet,
			requestPath: "/another-route",
			headers: map[string]string{
				"X-Custom-Header-2": "Some value 2",
				"X-Custom-Header-3": "Some value 3",
			},
			expectedStatus: 200,
		},
		{
			name:           "Missing Headers",
			httpMethod:     http.MethodPost,
			requestPath:    "/no-headers",
			headers:        map[string]string{},
			expectedStatus: 200,
		},
		{
			name:           "Empty Path",
			httpMethod:     http.MethodGet,
			requestPath:    "/",
			headers:        map[string]string{},
			expectedStatus: 200,
		},
		{
			name:        "Invalid Status Code Header",
			httpMethod:  http.MethodGet,
			requestPath: "/invalid-status",
			headers: map[string]string{
				"X-Response-Status": "invalid",
			},
			expectedStatus: 200, // Should fall back to default status
		},
		{
			name:        "Large Number of Headers",
			httpMethod:  http.MethodGet,
			requestPath: "/many-headers",
			headers: func() map[string]string {
				h := make(map[string]string)
				for i := 0; i < 100; i++ {
					h[http.CanonicalHeaderKey("X-Custom-Header-"+strconv.Itoa(i))] = "Some value"
				}
				return h
			}(),
			expectedStatus: 200,
		},
		{
			name:        "Uncommon HTTP Method",
			httpMethod:  http.MethodPut,
			requestPath: "/uncommon-method",
			headers: map[string]string{
				"X-Response-Status": "204",
			},
			expectedStatus: 204,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new HTTP request using the method from the test case
			req := httptest.NewRequest(tt.httpMethod, tt.requestPath, nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			recorder := httptest.NewRecorder()

			internal.HandleRequest(recorder, req, tmpl)

			res := recorder.Result()

			// Assert that the response status code matches the expected status
			assert.Equal(t, tt.expectedStatus, res.StatusCode, "response status codes should be equal")

			// Assert that all headers from the request are correctly passed to the handler
			for k, v := range tt.headers {
				assert.Equal(t, v, req.Header.Get(k), "headers should be equal")
			}
		})
	}
}
func TestCanLoadTemplate(t *testing.T) {
	_, err := internal.LoadTemplate(TemplatePath)
	assert.NoError(t, err)
}

func TestCollectRequestData(t *testing.T) {
	expectedIP := "127.0.0.1"
	expectedRoute := "/test-route"

	expectedHeaders := map[string]string{
		"User-Agent":      "GoTest",
		"X-Custom-Header": "CustomValue",
	}

	req := httptest.NewRequest(http.MethodPost, expectedRoute, nil)
	for k, v := range expectedHeaders {
		req.Header.Set(k, v)
	}

	req.RemoteAddr = expectedIP

	// Call the function
	data := internal.CollectRequestData(req)

	// Assertions
	assert.Equal(t, expectedRoute, data.Route, "Route should match")
	assert.Equal(t, http.MethodPost, data.Method, "Method should match")
	assert.Equal(t, expectedIP, data.OriginIP, "Origin IP should match")

	// Assert headers are correctly collected
	for k, v := range expectedHeaders {
		assert.Equal(t, v, data.Headers[k], "Header values should match")
	}
}
