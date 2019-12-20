package pipe_test

import (
	"errors"
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
		Equal(t, true, executed)
	})

	t.Run("More", func(t *testing.T) {
		var result int

		pipe := Pipe(
			func(n int) int { return n + 12 },
			func(n int) { result = n / 2 },
		)

		NotNil(t, pipe)
		NoError(t, pipe(10))
		Equal(t, 11, result)
	})

	t.Run("Error", func(t *testing.T) {
		var result int

		pipe := Pipe(
			func(n int) (int, error) {
				if n > 0 {
					return n * n, nil
				} else {
					return 0, errors.New("n must be positive")
				}
			},
			func(n int) { result = n / 2 },
		)

		NotNil(t, pipe)
		NoError(t, pipe(6))
		Equal(t, 18, result)

		err := pipe(-1)
		EqualError(t, err, "1st func failed: n must be positive")
		EqualError(t, errors.Unwrap(err), "n must be positive")
	})
}
