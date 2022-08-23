package debug

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestDebug_SetLevel(t *testing.T) {
	debug := NewDebug()

	assert.Equal(t, debug.Level(), NoDebug)

	err := debug.SetLevel(Important)
	assert.Equal(t, Important, debug.Level())
	assert.Nil(t, err)

	err2 := debug.SetLevel(Detailed)
	assert.Equal(t, Detailed, debug.Level())
	assert.Nil(t, err2)

	err3 := debug.SetLevel(math.MaxInt)
	assert.Equal(t, ErrDebugLevelNotSupported, err3)

	err4 := debug.SetLevel(-10)
	assert.Equal(t, ErrDebugLevelNotSupported, err4)
}
