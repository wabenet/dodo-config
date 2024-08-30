package includes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wabenet/dodo-config/pkg/includes"
)

func TestInclude(t *testing.T) {
	resolved, err := includes.ResolveIncludes("test/dodo.yaml")
	require.NoError(t, err)

	assert.Equal(t, len(resolved), 2)
	assert.Equal(t, resolved[0], "test/dodo.yaml")
	assert.Equal(t, resolved[1], "test/included.yaml")
}
