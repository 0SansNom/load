package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/0sansnom/load/internal/proxy"
)

func TestProxyIntegration(t *testing.T) {
	// Create a mock backend server
	serverMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Proxy Test Passed"))
	}))
	defer serverMock.Close()

	// Create a proxy instance that uses the mock backend
	proxyServer := proxy.New(nil, nil) // Mock the load balancer if needed
	proxyServer.SetTarget(serverMock.URL)

	// Send a request through the proxy
	resp, err := http.Get(proxyServer.URL + "/test")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Validate the response from the proxy
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Body, "Proxy Test Passed")
}
