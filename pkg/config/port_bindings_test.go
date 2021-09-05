package config_test

import (
	"testing"

	"github.com/dodo-cli/dodo-config/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestFullPortBindings(t *testing.T) {
	config, err := config.ParseConfig("test/dodo.yaml")
	assert.Nil(t, err)

	backdrop, ok := config["test_full_configs"]

	assert.True(t, ok)
	assert.Equal(t, 1, len(backdrop.Ports))

	for _, port := range backdrop.Ports {
		switch port.Target {
		case "80":
			assert.Equal(t, "8080", port.Published)
			assert.Equal(t, "192.168.0.1", port.HostIp)
		default:
			assert.Fail(t, "unexpected port config")
		}
	}
}

func TestPortBindingsWithList(t *testing.T) {
	config, err := config.ParseConfig("test/dodo.yaml")
	assert.Nil(t, err)

	backdrop, ok := config["test_with_lists"]

	assert.True(t, ok)
	assert.Equal(t, 1, len(backdrop.Ports))

	for _, port := range backdrop.Ports {
		switch port.Target {
		case "80":
			assert.Equal(t, "8080", port.Published)
		default:
			assert.Fail(t, "unexpected port config")
		}
	}
}
