package pipe_test

import (
	"testing"

	. "github.com/aslrousta/pipe"
	. "github.com/stretchr/testify/assert"
)

func TestPipe(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		pipe := Pipe()

		NotNil(t, pipe, "pipe is nil")
		NoError(t, pipe(), "unexpected error")
	})
}
