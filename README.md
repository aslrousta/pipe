# Go Pipe

![GitHub](https://img.shields.io/github/license/aslrousta/pipe)
[![GoDoc](https://godoc.org/github.com/aslrousta/pipe?status.svg)](https://godoc.org/github.com/aslrousta/pipe)

Go Pipe is a simple implementation of _pipe_ operator in Go.

## What is it all about?

One of the most interesting features of functional languages is the ability to
easily compose functions. The function composition is the operation that takes
two functions and produces another function in the way that the result of the
second function is passed as the input argument of first one.

Function composition is so common in functional programming which most of the
functional languages have defined a special operator, typically called a _pipe_
operator, to make composition even simpler. For example, in Haskell we have the
`.` (dot) operator, and Elixir defines the `|>` (pipe) operator:

```elixir
append = fn list, item ->
  list
  |> Enum.reverse
  |> prepend.(item)
  |> Enum.reverse
end
```

Actually, in any programming language which supports functional style we can
achieve the same result by applying functions consequently. For example in Go,
we can write:

```go
func h(x int) int {
    f := func(x int) int {
        if x < 0 {
            return -x
        } else {
            return x
        }
    }

    g := func(x int) int {
        return x * x
    }

    return f(g(x))
}
```

But the problem arises when one of the func's returns multiple values, as is
common with returning errors. In that case we need to handle error specifically
and write:

```go
func h(x int) (int, error) {
    f := func(x int) (int, error) {
        if x < 0 {
            return 0, errors.New("x should not be nagtive")
        }
        return x, nil
    }

    g := func(x int) int {
        return x * x
    }

    y, err := f(x)
    if err != nil {
        return 0, err
    }

    return g(y), nil
}
```

This is not much complex for two func's, but it can easily become a burden if
the number of func's grows.

To overcome this complexity I've written this library to make life simpler for
those who love functional style like me. Go Pipe provides a `Pipe` func which
composes one or more funcs into a pipeline func, which can be invoked anytime
withing the code.

Using `Pipe` func is very simple:

```go
import . "github.com/aslrousta/pipe"

func h(x int) int {
    var result int

    pipe := Pipe(
        func(x int) int {
            if x < 0 {
                return -x
            } else {
                return x
            }
        },
        func(x int) {
            result = x * x
        },
    )

    pipe(x)
    return result
}
```

And, it can handle errors automatically:

```go
import . "github.com/aslrousta/pipe"

func h(x int) (int, error) {
    var result int

    pipe := Pipe(
        func(x int) (int, error) {
            if x < 0 {
                return 0, errors.New("x should not be nagtive")
            }
            return x, nil
        },
        func(x int) {
            result = x * x
        },
    )

    err := pipe(x)
    return result, err
}
```

## License

Go Pipe source is released to the public under MIT license.
