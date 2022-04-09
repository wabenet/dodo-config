package config_test

import (
	"testing"

	"github.com/dodo-cli/dodo-config/internal/config"
	"github.com/dodo-cli/dodo-config/pkg/cuetils"
	"github.com/dodo-cli/dodo-config/pkg/spec"
	"github.com/stretchr/testify/assert"
)

func TestInclude(t *testing.T) {
	cfg, err := ParseTestConfig()
	assert.Nil(t, err)

	assert.Equal(t, 1, len(cfg.Includes))

	value, err := cuetils.ReadYAMLFileWithSpec(spec.CueSpec, cfg.Includes[0])
	assert.Nil(t, err)

	included, err := config.ConfigFromValue(value)
	assert.Nil(t, err)

	_, ok := included.Backdrops["included_backdrop"]

	assert.True(t, ok)
}
