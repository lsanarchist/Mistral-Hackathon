package llm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderCreationOnly(t *testing.T) {
	t.Run("Mistral provider creation", func(t *testing.T) {
		config := ProviderConfig{
			ProviderName: "mistral",
			APIKey:       "test-key",
			Model:        "test-model",
			Timeout:      10 * time.Second,
			MaxResponse:  100,
			DryRun:       true,
		}

		provider, err := NewProvider(config)
		require.NoError(t, err)
		assert.NotNil(t, provider)
		assert.Equal(t, "mistral", provider.Name())
		assert.Equal(t, "test-model", provider.Model())
	})

	t.Run("OpenAI provider creation", func(t *testing.T) {
		config := ProviderConfig{
			ProviderName: "openai",
			APIKey:       "test-key",
			Model:        "gpt-4",
			Timeout:      15 * time.Second,
			MaxResponse:  200,
			DryRun:       true,
		}

		provider, err := NewProvider(config)
		require.NoError(t, err)
		assert.NotNil(t, provider)
		assert.Equal(t, "openai", provider.Name())
		assert.Equal(t, "gpt-4", provider.Model())
	})

	t.Run("Default provider creation", func(t *testing.T) {
		config := ProviderConfig{
			ProviderName: "",
			APIKey:       "test-key",
			Model:        "test-model",
			Timeout:      10 * time.Second,
			DryRun:       true,
		}

		provider, err := NewProvider(config)
		require.NoError(t, err)
		assert.NotNil(t, provider)
		assert.Equal(t, "mistral", provider.Name())
	})

	t.Run("Unknown provider", func(t *testing.T) {
		config := ProviderConfig{
			ProviderName: "unknown",
			APIKey:       "test-key",
			Model:        "test-model",
			Timeout:      10 * time.Second,
			DryRun:       true,
		}

		provider, err := NewProvider(config)
		require.Error(t, err)
		assert.Nil(t, provider)
		assert.ErrorIs(t, err, ErrUnknownProvider)
	})
}

func TestDefaultConfigOnly(t *testing.T) {
	config := DefaultConfig()
	assert.Equal(t, "mistral", config.ProviderName)
	assert.Equal(t, "devstral-small-latest", config.Model)
	assert.Equal(t, 20*time.Second, config.Timeout)
	assert.Equal(t, 4096, config.MaxResponse)
	assert.Equal(t, 12000, config.MaxPrompt)
	assert.False(t, config.DryRun)
}