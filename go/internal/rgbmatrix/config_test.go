package rgbmatrix

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	err := os.WriteFile(path, []byte("[options]\nrows = 64\ncols = 32\n"), 0o600)
	assert.NoError(t, err)

	config, err := LoadConfigFile(path)

	assert.NoError(t, err)
	assert.Equal(t, config.Options.Rows, 64)
	assert.Equal(t, config.Options.ChainLength, 1)
}
