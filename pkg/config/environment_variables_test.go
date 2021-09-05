package config_test

import (
	"testing"

	"github.com/dodo-cli/dodo-config/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestFullEnvironment(t *testing.T) {
	config, err := config.ParseConfig("test/dodo.yaml")
	assert.Nil(t, err)

	backdrop, ok := config["test_full_configs"]

	assert.True(t, ok)
	assert.NotEmpty(t, backdrop.Environment)

	for _, env := range backdrop.Environment {
		assert.NotEqual(t, "FULL", env.Key)

		if env.Key == "FOO" {
			assert.Equal(t, "BAR", env.Value)
			return
		}
	}

	assert.Fail(t, "did not find expected environment config 'FULL'")
}

func TestPartialEnvironment(t *testing.T) {
	config, err := config.ParseConfig("test/dodo.yaml")
	assert.Nil(t, err)

	backdrop, ok := config["test_full_configs"]

	assert.True(t, ok)
	assert.NotEmpty(t, backdrop.Environment)

	for _, env := range backdrop.Environment {
		if env.Key == "PARTIAL" {
			assert.Equal(t, "", env.Value)
			return
		}
	}

	assert.Fail(t, "did not find expected environment config 'PARTIAL'")
}

func TestEnvironmentWithList(t *testing.T) {
	config, err := config.ParseConfig("test/dodo.yaml")
	assert.Nil(t, err)

	backdrop, ok := config["test_with_lists"]

	assert.True(t, ok)
	assert.Equal(t, 2, len(backdrop.Environment))

	for _, env := range backdrop.Environment {
		switch env.Key {
		case "FOO":
			assert.Equal(t, "BAR", env.Value)
		case "SOMETHING":
			assert.Equal(t, "", env.Value)
		default:
			assert.Fail(t, "unexpected environment config")
		}
	}
}
