package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wabenet/dodo-config/internal/config"
)

func TestFullVolume(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	assert.Nil(t, err)

	backdrop, ok := cfg["test_full_configs"]

	assert.True(t, ok)
	assert.NotEmpty(t, backdrop.Volumes)

	for _, vol := range backdrop.Volumes {
		if vol.Source == "/from/path" {
			assert.Equal(t, "/to/path", vol.Target)
			assert.True(t, vol.Readonly)

			return
		}
	}

	assert.Fail(t, "did not find expected volume config 'full'")
}

func TestPartialVolume(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	assert.Nil(t, err)

	backdrop, ok := cfg["test_full_configs"]

	assert.True(t, ok)
	assert.NotEmpty(t, backdrop.Volumes)

	for _, vol := range backdrop.Volumes {
		if vol.Source == "/some/mount" {
			return
		}
	}

	assert.Fail(t, "did not find expected volume config 'full'")
}

func TestVolumesWithLists(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	assert.Nil(t, err)

	backdrop, ok := cfg["test_with_lists"]

	assert.True(t, ok)
	assert.Equal(t, 1, len(backdrop.Volumes))

	for _, vol := range backdrop.Volumes {
		switch vol.Source {
		case "foo":
			assert.Equal(t, vol.Target, "bar")
			assert.True(t, vol.Readonly)
		default:
			assert.Fail(t, "unexpected volume config")
		}
	}
}
