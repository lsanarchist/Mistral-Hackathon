package llm

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsightsCache(t *testing.T) {
	t.Run("Test cache hit and miss", func(t *testing.T) {
		// Create temporary cache directory
		cacheDir, err := os.MkdirTemp("", "triageprof-cache-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(cacheDir)

		// Create cache config
		cacheConfig := CacheConfig{
			Enabled:        true,
			CacheDir:       cacheDir,
			MaxCacheSizeMB: 10,
			MaxCacheAgeHours: 24,
		}

		// Create cache
		cache := NewInsightsCache(cacheConfig)
		require.NotNil(t, cache)

		// Create test data
		bundle := &model.ProfileBundle{
			Metadata: model.Metadata{
				Timestamp:   time.Now(),
				DurationSec: 30,
				Service:     "test-service",
				Scenario:    "test-scenario",
				GitSha:      "abc123",
			},
			Target: model.Target{
				Type:    "go",
				BaseURL: "http://localhost:6060",
			},
		}

		findings := &model.FindingsBundle{
			Findings: []model.Finding{
				{
					Category: "Performance",
					Title:    "Test Finding",
					Severity: "High",
					Score:    85,
					Top:      []model.StackFrame{},
                Evidence: []model.EvidenceItem{
                    {
                        Type:        "profile",
                        Description: "Profile evidence",
                        Value:       "profile.pb.gz",
                        Weight:      1.0,
                    },
                },
				},
			},
		}

		insights := &model.InsightsBundle{
			SchemaVersion: "2.0",
			GeneratedAt:   time.Now(),
			ExecutiveSummary: model.ExecutiveSummary{
				Overview:        "Test overview",
				OverallSeverity: "Medium",
				Confidence:      85,
			},
		}

		ctx := context.Background()

		// Test cache miss
		cachedInsights, found := cache.GetCachedInsights(ctx, bundle, findings)
		assert.False(t, found)
		assert.Nil(t, cachedInsights)

		// Cache the insights
		err = cache.CacheInsights(ctx, bundle, findings, insights)
		require.NoError(t, err)

		// Test cache hit
		cachedInsights, found = cache.GetCachedInsights(ctx, bundle, findings)
		assert.True(t, found)
		require.NotNil(t, cachedInsights)
		assert.Equal(t, insights.SchemaVersion, cachedInsights.SchemaVersion)
		assert.Equal(t, insights.ExecutiveSummary.Overview, cachedInsights.ExecutiveSummary.Overview)

		// Test cache stats
		stats := cache.GetCacheStats()
		assert.Equal(t, 1, stats.Hits)
		assert.Equal(t, 1, stats.Misses)
		assert.Equal(t, 1, stats.Writes)
		assert.Equal(t, 0, stats.Evictions)
		assert.Greater(t, stats.BytesSaved, int64(0))
	})

	t.Run("Test cache expiration", func(t *testing.T) {
		// Create temporary cache directory
		cacheDir, err := os.MkdirTemp("", "triageprof-cache-expiry-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(cacheDir)

		// Create cache config with very short expiration (0 hours = immediate)
		cacheConfig := CacheConfig{
			Enabled:          true,
			CacheDir:         cacheDir,
			MaxCacheSizeMB:   10,
			MaxCacheAgeHours: 0, // Immediate expiration
		}

		// Create cache
		cache := NewInsightsCache(cacheConfig)

		// Create test data
		bundle := &model.ProfileBundle{
			Metadata: model.Metadata{
				Service: "test-service",
			},
		}

		findings := &model.FindingsBundle{}

		insights := &model.InsightsBundle{
			SchemaVersion: "2.0",
			ExecutiveSummary: model.ExecutiveSummary{
				Overview: "Test insights",
			},
		}

		ctx := context.Background()

		// Cache the insights
		err = cache.CacheInsights(ctx, bundle, findings, insights)
		require.NoError(t, err)

		// Small delay to ensure cache is older than expiration time
		time.Sleep(1 * time.Millisecond)

		// Test cache miss due to expiration
		cachedInsights, found := cache.GetCachedInsights(ctx, bundle, findings)
		assert.False(t, found, "Should miss due to expiration")
		assert.Nil(t, cachedInsights)

		// Check eviction stat
		stats := cache.GetCacheStats()
		assert.Equal(t, 1, stats.Evictions, "Should have one eviction")
	})

	t.Run("Test cache size limit", func(t *testing.T) {
		// Create temporary cache directory
		cacheDir, err := os.MkdirTemp("", "triageprof-cache-size-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(cacheDir)

		// Create cache config with very small size limit
		cacheConfig := CacheConfig{
			Enabled:        true,
			CacheDir:       cacheDir,
			MaxCacheSizeMB: 1, // Very small limit (1MB)
			MaxCacheAgeHours: 24,
		}

		// Create cache
		cache := NewInsightsCache(cacheConfig)

		// Create base test data
		baseBundle := &model.ProfileBundle{
			Metadata: model.Metadata{
				Service: "test-service",
			},
		}

		// Cache many large insights to trigger size limit
		for i := 0; i < 20; i++ {
			bundle := *baseBundle
			bundle.Metadata.Service = fmt.Sprintf("service-%d", i)

			findings := &model.FindingsBundle{}

			// Create very large insights
			largeInsights := &model.InsightsBundle{
				SchemaVersion: "2.0",
				GeneratedAt:   time.Now(),
				ExecutiveSummary: model.ExecutiveSummary{
					Overview: "Large insights " + strings.Repeat("x", 200000), // Very large content (200KB each)
				},
			}

			err = cache.CacheInsights(context.Background(), &bundle, findings, largeInsights)
			require.NoError(t, err)
		}

		// Check that some evictions occurred due to size limit
		stats := cache.GetCacheStats()
		assert.Greater(t, stats.Evictions, 0, "Should have evictions due to size limit")
		assert.Equal(t, 20, stats.Writes)
	})

	t.Run("Test cache disabled", func(t *testing.T) {
		cacheConfig := CacheConfig{
			Enabled: false,
		}

		cache := NewInsightsCache(cacheConfig)

		bundle := &model.ProfileBundle{}
		findings := &model.FindingsBundle{}
		insights := &model.InsightsBundle{
			SchemaVersion: "2.0",
		}

		// Test that cache operations are no-ops when disabled
		err := cache.CacheInsights(context.Background(), bundle, findings, insights)
		require.NoError(t, err)

		cachedInsights, found := cache.GetCachedInsights(context.Background(), bundle, findings)
		assert.False(t, found)
		assert.Nil(t, cachedInsights)

		stats := cache.GetCacheStats()
		assert.Equal(t, 0, stats.Hits)
		assert.Equal(t, 0, stats.Misses)
		assert.Equal(t, 0, stats.Writes)
	})

	t.Run("Test clear cache", func(t *testing.T) {
		// Create temporary cache directory
		cacheDir, err := os.MkdirTemp("", "triageprof-cache-clear-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(cacheDir)

		// Create cache config
		cacheConfig := CacheConfig{
			Enabled:        true,
			CacheDir:       cacheDir,
			MaxCacheSizeMB: 10,
			MaxCacheAgeHours: 24,
		}

		// Create cache
		cache := NewInsightsCache(cacheConfig)

		// Add some cache entries
		for i := 0; i < 5; i++ {
			bundle := &model.ProfileBundle{
				Metadata: model.Metadata{
					Service: fmt.Sprintf("service-%d", i),
				},
			}
			findings := &model.FindingsBundle{}
			insights := &model.InsightsBundle{
				SchemaVersion: "2.0",
			}

			err = cache.CacheInsights(context.Background(), bundle, findings, insights)
			require.NoError(t, err)
		}

		// Verify cache has entries
		files, err := os.ReadDir(cacheDir)
		require.NoError(t, err)
		assert.Greater(t, len(files), 0)

		// Clear cache
		err = cache.ClearCache()
		require.NoError(t, err)

		// Verify cache is empty
		files, err = os.ReadDir(cacheDir)
		require.NoError(t, err)
		assert.Equal(t, 0, len(files))

		// Verify stats are reset
		stats := cache.GetCacheStats()
		assert.Equal(t, 0, stats.Hits)
		assert.Equal(t, 0, stats.Misses)
		assert.Equal(t, 0, stats.Writes)
	})
}