package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadManifest(t *testing.T) {
	t.Run("valid manifest", func(t *testing.T) {
		manifest, err := LoadManifest("../../plugins/manifests/go-pprof-http.json")
		require.NoError(t, err)
		assert.Equal(t, "go-pprof-http", manifest.Name)
		assert.Equal(t, "0.1.0", manifest.Version)
		assert.Equal(t, "1.0", manifest.SDKVersion)
		assert.Contains(t, manifest.Capabilities.Targets, "url")
		assert.Contains(t, manifest.Capabilities.Profiles, "cpu")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		// Create temporary invalid JSON file
		tempDir := t.TempDir()
		invalidPath := filepath.Join(tempDir, "invalid.json")
		err := os.WriteFile(invalidPath, []byte("{invalid json}"), 0644)
		require.NoError(t, err)

		_, err = LoadManifest(invalidPath)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse manifest file")
	})

	t.Run("unknown field", func(t *testing.T) {
		// Create temporary manifest with unknown field
		tempDir := t.TempDir()
		unknownPath := filepath.Join(tempDir, "unknown.json")
		content := `{"name": "test", "version": "1.0", "sdkVersion": "1.0", "unknownField": "value"}`
		err := os.WriteFile(unknownPath, []byte(content), 0644)
		require.NoError(t, err)

		_, err = LoadManifest(unknownPath)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown field")
	})

	t.Run("missing required field", func(t *testing.T) {
		// Create temporary manifest missing name field
		tempDir := t.TempDir()
		missingPath := filepath.Join(tempDir, "missing.json")
		content := `{"version": "1.0", "sdkVersion": "1.0", "capabilities": {"targets": ["url"], "profiles": ["cpu"]}}`
		err := os.WriteFile(missingPath, []byte(content), 0644)
		require.NoError(t, err)

		_, err = LoadManifest(missingPath)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing required field 'name'")
	})
}

func TestDiscoverManifests(t *testing.T) {
	t.Run("discover valid manifests", func(t *testing.T) {
		manifests, err := DiscoverManifests("../../plugins/manifests")
		require.NoError(t, err)
		assert.NotEmpty(t, manifests)
		assert.Equal(t, "go-pprof-http", manifests[0].Name)
	})

	t.Run("empty directory", func(t *testing.T) {
		// Create empty temp directory
		tempDir := t.TempDir()
		manifests, err := DiscoverManifests(tempDir)
		require.NoError(t, err)
		assert.Empty(t, manifests)
	})

	t.Run("mixed valid and invalid files", func(t *testing.T) {
		// Create temp directory with both valid and invalid manifests
		tempDir := t.TempDir()

		// Valid manifest
		validPath := filepath.Join(tempDir, "valid.json")
		validContent := `{"name": "valid-plugin", "version": "1.0", "sdkVersion": "1.0", "capabilities": {"targets": ["url"], "profiles": ["cpu"]}}`
		err := os.WriteFile(validPath, []byte(validContent), 0644)
		require.NoError(t, err)

		// Invalid manifest
		invalidPath := filepath.Join(tempDir, "invalid.json")
		err = os.WriteFile(invalidPath, []byte("{invalid}"), 0644)
		require.NoError(t, err)

		// Should only return valid manifests
		manifests, err := DiscoverManifests(tempDir)
		require.NoError(t, err)
		assert.Len(t, manifests, 1)
		assert.Equal(t, "valid-plugin", manifests[0].Name)
	})
}

func TestResolvePlugin(t *testing.T) {
	t.Run("resolve existing plugin", func(t *testing.T) {
		manifest, binaryPath, err := ResolvePlugin("../../plugins/manifests", "../../plugins/bin", "go-pprof-http")
		if err != nil && strings.Contains(err.Error(), "binary missing") {
			t.Skip("go-pprof-http binary not built; skipping test")
		}
		require.NoError(t, err)
		assert.Equal(t, "go-pprof-http", manifest.Name)
		assert.Contains(t, binaryPath, "go-pprof-http")
	})

	t.Run("plugin not found", func(t *testing.T) {
		_, _, err := ResolvePlugin("../../plugins/manifests", "../../plugins/bin", "nonexistent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "plugin nonexistent not found")
		assert.Contains(t, err.Error(), "Available plugins:")
	})

	t.Run("missing binary", func(t *testing.T) {
		// Create temp manifest directory
		tempDir := t.TempDir()
		manifestDir := filepath.Join(tempDir, "manifests")
		binDir := filepath.Join(tempDir, "bin")
		err := os.MkdirAll(manifestDir, 0755)
		require.NoError(t, err)
		err = os.MkdirAll(binDir, 0755)
		require.NoError(t, err)

		// Create manifest
		manifestPath := filepath.Join(manifestDir, "test.json")
		content := `{"name": "test-plugin", "version": "1.0", "sdkVersion": "1.0", "capabilities": {"targets": ["url"], "profiles": ["cpu"]}}`
		err = os.WriteFile(manifestPath, []byte(content), 0644)
		require.NoError(t, err)

		// Try to resolve (should fail due to missing binary)
		_, _, err = ResolvePlugin(manifestDir, binDir, "test-plugin")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "binary missing")
	})

	t.Run("SDK version mismatch", func(t *testing.T) {
		// Create temp manifest directory
		tempDir := t.TempDir()
		manifestDir := filepath.Join(tempDir, "manifests")
		binDir := filepath.Join(tempDir, "bin")
		err := os.MkdirAll(manifestDir, 0755)
		require.NoError(t, err)
		err = os.MkdirAll(binDir, 0755)
		require.NoError(t, err)

		// Create manifest with incompatible SDK version
		manifestPath := filepath.Join(manifestDir, "test.json")
		content := `{"name": "test-plugin", "version": "1.0", "sdkVersion": "2.0", "capabilities": {"targets": ["url"], "profiles": ["cpu"]}}`
		err = os.WriteFile(manifestPath, []byte(content), 0644)
		require.NoError(t, err)

		// Create dummy binary
		binaryPath := filepath.Join(binDir, "test-plugin")
		err = os.WriteFile(binaryPath, []byte("#!/bin/sh"), 0755)
		require.NoError(t, err)

		// Try to resolve (should fail due to SDK version mismatch)
		_, _, err = ResolvePlugin(manifestDir, binDir, "test-plugin")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "requires sdkVersion 2.0, but core supports 1.0")
	})
}

func TestValidateTarget(t *testing.T) {
	manifest := &Manifest{
		Name:       "test-plugin",
		Version:    "1.0",
		SDKVersion: "1.0",
		Capabilities: Capabilities{
			Targets:  []string{"url", "python"},
			Profiles: []string{"cpu", "heap"},
		},
	}

	t.Run("supported target", func(t *testing.T) {
		err := manifest.ValidateTarget("url")
		require.NoError(t, err)
	})

	t.Run("unsupported target", func(t *testing.T) {
		err := manifest.ValidateTarget("database")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "target type 'database' not supported")
		assert.Contains(t, err.Error(), "Supported targets: url, python")
	})
}

func TestValidateProfiles(t *testing.T) {
	manifest := &Manifest{
		Name:       "test-plugin",
		Version:    "1.0",
		SDKVersion: "1.0",
		Capabilities: Capabilities{
			Targets:  []string{"url"},
			Profiles: []string{"cpu", "heap"},
		},
	}

	t.Run("supported profiles", func(t *testing.T) {
		err := manifest.ValidateProfiles([]string{"cpu", "heap"})
		require.NoError(t, err)
	})

	t.Run("mixed supported and unsupported", func(t *testing.T) {
		err := manifest.ValidateProfiles([]string{"cpu", "mutex"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "profiles 'mutex' not supported")
		assert.Contains(t, err.Error(), "Supported profiles: cpu, heap")
	})

	t.Run("all unsupported", func(t *testing.T) {
		err := manifest.ValidateProfiles([]string{"mutex", "block"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "profiles 'mutex, block' not supported")
	})

	t.Run("empty list", func(t *testing.T) {
		err := manifest.ValidateProfiles([]string{})
		require.NoError(t, err)
	})
}
