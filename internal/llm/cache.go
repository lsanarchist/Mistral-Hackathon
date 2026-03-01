package llm

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// CacheConfig defines configuration for the insights cache
type CacheConfig struct {
	Enabled          bool
	CacheDir         string
	MaxCacheSizeMB   int
	MaxCacheAgeHours int
}

// InsightsCache handles caching of LLM-generated insights
type InsightsCache struct {
	config      CacheConfig
	cacheDir    string
	stats       CacheStats
	mu          sync.Mutex
}

// CacheStats tracks cache performance
type CacheStats struct {
	Hits       int
	Misses     int
	Writes     int
	Evictions  int
	BytesSaved int64
}

// NewInsightsCache creates a new insights cache
func NewInsightsCache(config CacheConfig) *InsightsCache {
	cacheDir := config.CacheDir
	if cacheDir == "" {
		cacheDir = filepath.Join(os.TempDir(), "triageprof-insights-cache")
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		log.Printf("Warning: failed to create cache directory: %v", err)
	}

	return &InsightsCache{
		config:    config,
		cacheDir:  cacheDir,
		stats:     CacheStats{},
	}
}

// generateCacheKey creates a unique key based on profile bundle and findings
func generateCacheKey(bundle *model.ProfileBundle, findings *model.FindingsBundle) (string, error) {
	// Create a combined struct for hashing
	type cacheKeyData struct {
		Metadata    model.Metadata
		Target      model.Target
		Plugin      model.PluginRef
		Findings    []model.Finding
		GeneratedAt time.Time
	}

	// Extract relevant data for cache key
	var findingItems []model.Finding
	if findings != nil {
		findingItems = findings.Findings
	}

	keyData := cacheKeyData{
		Metadata:    bundle.Metadata,
		Target:      bundle.Target,
		Plugin:      bundle.Plugin,
		Findings:    findingItems,
		GeneratedAt: time.Now().Truncate(time.Hour), // Hour precision for cache
	}

	// Marshal to JSON for consistent hashing
	data, err := json.Marshal(keyData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cache key data: %v", err)
	}

	// Create SHA256 hash
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// GetCachedInsights retrieves cached insights if available
func (c *InsightsCache) GetCachedInsights(ctx context.Context, bundle *model.ProfileBundle, findings *model.FindingsBundle) (*model.InsightsBundle, bool) {
	if !c.config.Enabled {
		return nil, false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Generate cache key
	cacheKey, err := generateCacheKey(bundle, findings)
	if err != nil {
		log.Printf("Warning: failed to generate cache key: %v", err)
		c.stats.Misses++
		return nil, false
	}

	cacheFile := filepath.Join(c.cacheDir, cacheKey+".json")

	// Check if cache file exists
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		c.stats.Misses++
		return nil, false
	}

	// Check cache age
	if c.config.MaxCacheAgeHours >= 0 { // 0 means immediate expiration
		fileInfo, err := os.Stat(cacheFile)
		if err != nil {
			c.stats.Misses++
			return nil, false
		}

		cacheAge := time.Since(fileInfo.ModTime())
		if c.config.MaxCacheAgeHours == 0 || cacheAge.Hours() > float64(c.config.MaxCacheAgeHours) {
			// Cache entry is too old, remove it
			if err := os.Remove(cacheFile); err != nil {
				log.Printf("Warning: failed to remove expired cache entry: %v", err)
			}
			c.stats.Evictions++
			c.stats.Misses++
			log.Printf("Evicted expired cache entry: %s (age: %.1f hours)", cacheFile, cacheAge.Hours())
			return nil, false
		}
	}

	// Read cached insights
	cachedData, err := os.ReadFile(cacheFile)
	if err != nil {
		log.Printf("Warning: failed to read cache file: %v", err)
		c.stats.Misses++
		return nil, false
	}

	var insights model.InsightsBundle
	if err := json.Unmarshal(cachedData, &insights); err != nil {
		log.Printf("Warning: failed to parse cached insights: %v", err)
		c.stats.Misses++
		return nil, false
	}

	c.stats.Hits++
	c.stats.BytesSaved += int64(len(cachedData))
	log.Printf("Cache hit for insights (saved ~%d bytes)", len(cachedData))

	return &insights, true
}

// CacheInsights stores generated insights in the cache
func (c *InsightsCache) CacheInsights(ctx context.Context, bundle *model.ProfileBundle, findings *model.FindingsBundle, insights *model.InsightsBundle) error {
	if !c.config.Enabled || insights == nil {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Generate cache key
	cacheKey, err := generateCacheKey(bundle, findings)
	if err != nil {
		return fmt.Errorf("failed to generate cache key: %v", err)
	}

	cacheFile := filepath.Join(c.cacheDir, cacheKey+".json")

	// Marshal insights to JSON
	insightsData, err := json.Marshal(insights)
	if err != nil {
		return fmt.Errorf("failed to marshal insights for caching: %v", err)
	}

	// Write to cache file
	if err := os.WriteFile(cacheFile, insightsData, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %v", err)
	}

	c.stats.Writes++
	log.Printf("Cached insights for key: %s", cacheKey)

	// Enforce cache size limit
	if c.config.MaxCacheSizeMB > 0 {
		if err := c.enforceCacheSizeLimit(); err != nil {
			log.Printf("Warning: failed to enforce cache size limit: %v", err)
		}
	}

	return nil
}

// enforceCacheSizeLimit removes oldest cache entries if size limit is exceeded
func (c *InsightsCache) enforceCacheSizeLimit() error {
	files, err := os.ReadDir(c.cacheDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %v", err)
	}

	// Get file info with modification times
	var fileInfos []struct {
		name    string
		modTime time.Time
		size    int64
	}

	var totalSize int64
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			continue
		}

		fileInfos = append(fileInfos, struct {
			name    string
			modTime time.Time
			size    int64
		}{
			name:    file.Name(),
			modTime: fileInfo.ModTime(),
			size:    fileInfo.Size(),
		})
		totalSize += fileInfo.Size()
	}

	// Convert MB limit to bytes
	maxSizeBytes := int64(c.config.MaxCacheSizeMB) * 1024 * 1024

	// Remove oldest files until we're under the limit
	for totalSize > maxSizeBytes && len(fileInfos) > 0 {
		// Find oldest file
		oldestIndex := 0
		for i, info := range fileInfos {
			if info.modTime.Before(fileInfos[oldestIndex].modTime) {
				oldestIndex = i
			}
		}

		oldestFile := fileInfos[oldestIndex]
		cacheFile := filepath.Join(c.cacheDir, oldestFile.name)

		if err := os.Remove(cacheFile); err != nil {
			log.Printf("Warning: failed to remove cache file %s: %v", oldestFile.name, err)
		} else {
			totalSize -= oldestFile.size
			c.stats.Evictions++
			log.Printf("Evicted cache entry: %s (freed %d bytes)", oldestFile.name, oldestFile.size)
		}

		// Remove from slice
		fileInfos = append(fileInfos[:oldestIndex], fileInfos[oldestIndex+1:]...)
	}

	return nil
}

// GetCacheStats returns current cache statistics
func (c *InsightsCache) GetCacheStats() CacheStats {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.stats
}

// ClearCache removes all cached insights
func (c *InsightsCache) ClearCache() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	files, err := os.ReadDir(c.cacheDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %v", err)
	}

	var clearedCount int
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			cacheFile := filepath.Join(c.cacheDir, file.Name())
			if err := os.Remove(cacheFile); err != nil {
				log.Printf("Warning: failed to remove cache file %s: %v", file.Name(), err)
			} else {
				clearedCount++
			}
		}
	}

	log.Printf("Cleared %d cache entries", clearedCount)
	c.stats = CacheStats{}
	return nil
}