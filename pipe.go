package pipe

// Func is the func type for the pipeline result.
type Func func() error

func emptyFunc() error { return nil }

// Pipe accepts several funcs and creates a pipeline.
func Pipe(fs ...interface{}) Func {
	if len(fs) == 0 {
		return emptyFunc
	}

	return nil
}
