package rgbmatrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	config := LoadConfig()

	assert.Equal(t, config.Options.Rows, 64)
}
