package pipe_test

import (
	"errors"
	"math/rand"
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

var result int

func BenchmarkPipe(b *testing.B) {
	sqr := func(x int) int { return x * x }
	inc := func(x int) int { return x + 1 }

	x := rand.Intn(1000)

	b.Run("SqrPlusOneDirect", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			result = inc(sqr(x))
		}
	})

	b.Run("SqrPlusOnePipe", func(b *testing.B) {
		pipe := Pipe(sqr, inc, func(x int) { result = x })

		for n := 0; n < b.N; n++ {
			pipe(x)
		}
	})
}
