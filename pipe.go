package pipe

import (
	"fmt"
	"reflect"
)

// errType is the type of error interface.
var errType = reflect.TypeOf((*error)(nil)).Elem()

// Pipeline is the func type for the pipeline result.
type Pipeline func(...interface{}) error

func empty(...interface{}) error { return nil }

// Pipe accepts zero or more funcs fs and creates a pipeline.
//
// A pipeline syncs outputs and inputs of consequent funcs together, such that
// the output of i'th func is the input of (i+1)'th func. Each func can accept
// zero or one input argument and return zero or one value with an optional
// error.
//
// The last func is called a sink which only accepts an input argument and
// returns no value except an optional error; unless its return value will be
// ignored.
//
// If a func in the pipeline fails with an error during the invocation, the pipe
// is broken immediately and the invocation returns an error.
func Pipe(fs ...interface{}) Pipeline {
	if len(fs) == 0 {
		return empty
	}

	return func(args ...interface{}) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("pipeline panicked: %v", r)
			}
		}()

		var inputs []reflect.Value
		for _, arg := range args {
			inputs = append(inputs, reflect.ValueOf(arg))
		}

		for fIndex, f := range fs {
			outputs := reflect.ValueOf(f).Call(inputs)
			inputs = inputs[:0]

			funcType := reflect.TypeOf(f)

			for oIndex, output := range outputs {
				if funcType.Out(oIndex).Implements(errType) {
					if !output.IsNil() {
						err = fmt.Errorf("%s func failed: %w", ord(fIndex), output.Interface().(error))
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
