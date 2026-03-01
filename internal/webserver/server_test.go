
package webserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginMarketplaceEndpoint(t *testing.T) {
	// Create WebSocket server
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false)
	defer server.Stop()

	// Test marketplace endpoint
	req := httptest.NewRequest("GET", "/plugins/marketplace", nil)
	w := httptest.NewRecorder()

	server.handlePluginMarketplace(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, response, "plugins")
	assert.Contains(t, response, "count")
	
	plugins := response["plugins"].([]interface{})
	assert.True(t, len(plugins) > 0, "Should have at least one plugin in marketplace")
}

func TestPluginInstallEndpoint(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false)
	defer server.Stop()

	// Test install endpoint with empty body
	req := httptest.NewRequest("POST", "/plugins/install", nil)
	w := httptest.NewRecorder()

	server.handleInstallPlugin(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPluginUpdateEndpoint(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false)
	defer server.Stop()

	// Test update endpoint with empty body
	req := httptest.NewRequest("POST", "/plugins/update", nil)
	w := httptest.NewRecorder()

	server.handleUpdatePlugin(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPluginUninstallEndpoint(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false)
	defer server.Stop()

	// Test uninstall endpoint with empty body
	req := httptest.NewRequest("POST", "/plugins/uninstall", nil)
	w := httptest.NewRecorder()

	server.handleUninstallPlugin(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
