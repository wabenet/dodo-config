package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wabenet/dodo-config/internal/config"
)

func TestFullDeviceMappings(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	assert.Nil(t, err)

	backdrop, ok := cfg["test_full_configs"]

	assert.True(t, ok)
	assert.Equal(t, 2, len(backdrop.Devices))

	for _, dev := range backdrop.Devices {
		switch dev.Target {
		case "/foo/bar":
			assert.Equal(t, "/dev/snd", dev.Source)
			assert.Equal(t, "rw", dev.Permissions)
		case "rule":
			assert.Equal(t, "c *:* rmw", dev.CgroupRule)
		default:
			assert.Fail(t, "unexpected port config")
		}
	}
}

func TestDeviceMappingsWithList(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	assert.Nil(t, err)

	backdrop, ok := cfg["test_with_lists"]

	assert.True(t, ok)
	assert.Equal(t, 1, len(backdrop.Devices))

	for _, dev := range backdrop.Devices {
		switch dev.Target {
		case "/dev/snd":
			assert.Equal(t, "/dev/snd", dev.Source)
		default:
			assert.Fail(t, "unexpected port config")
		}
	}
}
