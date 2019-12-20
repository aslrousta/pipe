package pipe_test

import (
	"testing"

	. "github.com/aslrousta/pipe"
	. "github.com/stretchr/testify/assert"
)

func TestPipe(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		pipe := Pipe()

		NotNil(t, pipe)
		NoError(t, pipe())
	})

	t.Run("One", func(t *testing.T) {
		executed := false
		pipe := Pipe(func() { executed = true })

		NotNil(t, pipe)
		NoError(t, pipe())
		Equal(t, executed, true)
	})

	t.Run("More", func(t *testing.T) {
		var result int

		pipe := Pipe(
			func() int { return 10 },
			func(n int) int { return n + 12 },
			func(n int) { result = n / 2 },
		)

		NotNil(t, pipe)
		NoError(t, pipe())
		Equal(t, result, 11)
	})
}
