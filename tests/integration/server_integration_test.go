package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/0sansnom/load/internal/server"
)

func TestServerIntegration(t *testing.T) {
	// Mocking a backend server
	serverMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mocked Server Response"))
	}))
	defer serverMock.Close()

	// Create a new HTTP server with our mock server
	httpServer := server.New(nil, nil, nil) // Assuming you mock or bypass dependencies
	httpServer.SetHandler(http.NewServeMux())

	// Send a request to our server
	resp, err := http.Get(httpServer.URL + "/some-endpoint")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check the status code and response body
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// Assume the response from the mock backend will be "Mocked Server Response"
	assert.Contains(t, resp.Body, "Mocked Server Response")
}
