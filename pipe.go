package pipe

import (
	"fmt"
	"reflect"
)

// Pipeline is the func type for the pipeline result.
type Pipeline func() error

func empty() error { return nil }

// Pipe accepts zero or more funcs fs and creates a pipeline.
//
// A pipeline syncs outputs and inputs of consequent funcs together, such that
// the output of i'th func is the input of (i+1)'th func. Each func can accept
// zero or one input argument and return zero or one value with an optional
// error.
//
// The first func is called a generator and can accept no input arguments but
// only return a value (optionally, with an error), and the last func is called
// a sink which only accepts an input argument and returns no value except an
// optional error.
//
// If a func in the pipeline fails with an error during the invocation, the pipe
// is broken immediately and the invocation returns an error.
func Pipe(fs ...interface{}) Pipeline {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("pipeline panicked: %v", r)
			}
		}()

		var inputs []reflect.Value
		for i, f := range fs {
			if f == nil {
				return fmt.Errorf("%s arg is nil", ord(i))
			}

			t := reflect.TypeOf(f)
			if t.Kind() != reflect.Func {
				return fmt.Errorf("%s arg is not a func", ord(i))
			}

			if t.NumIn() != len(inputs) {
				return fmt.Errorf(
					"%s func accepts %d args, %d passed",
					ord(i), t.NumIn(), len(inputs),
				)
			}

			outputs := reflect.ValueOf(f).Call(inputs)

			inputs = inputs[:0]
			for _, output := range outputs {
				if e, ok := output.Interface().(error); ok {
					if err != nil {
						err = fmt.Errorf("%s func failed: %w", ord(i), e)
						return
					}
				} else {
					inputs = append(inputs, output)
				}
			}
		}

		return
	}
}

func ord(index int) string {
	order := index + 1
	switch {
	case order > 10 && order < 20:
		return fmt.Sprintf("%dth", order)
	case order%10 == 1:
		return fmt.Sprintf("%dst", order)
	case order%10 == 2:
		return fmt.Sprintf("%dnd", order)
	case order%10 == 3:
		return fmt.Sprintf("%drd", order)
	default:
		return fmt.Sprintf("%dth", order)
	}
}
