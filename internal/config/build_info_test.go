package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wabenet/dodo-config/internal/config"
)

func TestBuildInfo(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	backdrop, ok := cfg["test_build"]
	assert.True(t, ok)

	build := backdrop.GetBuildInfo()
	assert.NotNil(t, backdrop.BuildInfo)

	assert.Equal(t, "testimage", build.GetImageName())
	assert.Equal(t, "/some/path", build.GetContext())
	assert.Equal(t, "/some/other/path", build.GetDockerfile())
	assert.Equal(t, []string{"FROM foo\n"}, build.GetInlineDockerfile())
}
