package config_test

import (
	"testing"

	"github.com/dodo-cli/dodo-config/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestInclude(t *testing.T) {
	cfg, err := config.ParseConfig("test/dodo.yaml")
	assert.Nil(t, err)

	assert.Equal(t, 1, len(cfg.Includes))

	included, err := config.ParseConfig(cfg.Includes[0])
	assert.Nil(t, err)

	_, ok := included.Backdrops["included_backdrop"]

	assert.True(t, ok)
}
