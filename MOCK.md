# gomock

[![Build Status][ci-badge]][ci-runs] [![Go Reference][reference-badge]][reference]

gomock is a mocking framework for the [Go programming language][golang]. It
integrates well with Go's built-in `testing` package, but can be used in other
contexts too.


## Installation

Install the `mockgen` tool.

```
go install go.uber.org/mock/mockgen@latest
```

To ensure it was installed correctly, use:

```
mockgen -version
```

If that fails, make sure your GOPATH/bin is in your PATH. You can add it with:

```
export PATH=$PATH:$(go env GOPATH)/bin
```

## Running mockgen

`mockgen` has two modes of operation: source and reflect.

### Source mode

Source mode generates mock interfaces from a source file.
It is enabled by using the -source flag. Other flags that
may be useful in this mode are -imports and -aux_files.

Example:

```bash
mockgen -source=foo.go [other options]
```

### Reflect mode

Reflect mode generates mock interfaces by building a program
that uses reflection to understand interfaces. It is enabled
by passing two non-flag arguments: an import path, and a
comma-separated list of symbols.

You can use "." to refer to the current path's package.

Example:

Example:

```bash
	mockgen -package mockcloud -destination mock/cloud/cloud.go  user/adapter/cloud CloudStorage,Provider

# Convenient for `go:generate`.
mockgen . CloudStorage,Provider
```

In this example, the destination is the relative path from the Go module. 
The next argument specifies the absolute path of the interface you want to mock.

### Flags

The `mockgen` command is used to generate source code for a mock
class given a Go source file containing interfaces to be mocked.
It supports the following flags:

- `-source`: A file containing interfaces to be mocked.

- `-destination`: A file to which to write the resulting source code. If you
  don't set this, the code is printed to standard output.

- `-package`: The package to use for the resulting mock class
  source code. If you don't set this, the package name is `mock_` concatenated
  with the package of the input file.

- `-imports`: A list of explicit imports that should be used in the resulting
  source code, specified as a comma-separated list of elements of the form
  `foo=bar/baz`, where `bar/baz` is the package being imported and `foo` is
  the identifier to use for the package in the generated source code.

- `-aux_files`: A list of additional files that should be consulted to
  resolve e.g. embedded interfaces defined in a different file. This is
  specified as a comma-separated list of elements of the form
  `foo=bar/baz.go`, where `bar/baz.go` is the source file and `foo` is the
  package name of that file used by the -source file.

- `-build_flags`: (reflect mode only) Flags passed verbatim to `go build`.

- `-mock_names`: A list of custom names for generated mocks. This is specified
  as a comma-separated list of elements of the form
  `Repository=MockSensorRepository,Endpoint=MockSensorEndpoint`, where
  `Repository` is the interface name and `MockSensorRepository` is the desired
  mock name (mock factory method and mock recorder will be named after the mock).
  If one of the interfaces has no custom name specified, then default naming
  convention will be used.

- `-self_package`: The full package import path for the generated code. The
  purpose of this flag is to prevent import cycles in the generated code by
  trying to include its own package. This can happen if the mock's package is
  set to one of its inputs (usually the main one) and the output is stdio so
  mockgen cannot detect the final output package. Setting this flag will then
  tell mockgen which import to exclude.

- `-copyright_file`: Copyright file used to add copyright header to the resulting source code.

- `-debug_parser`: Print out parser results only.

- `-exec_only`: (reflect mode) If set, execute this reflection program.

- `-prog_only`: (reflect mode) Only generate the reflection program; write it to stdout and exit.

- `-write_package_comment`: Writes package documentation comment (godoc) if true. (default true)

- `-write_generate_directive`: Add //go:generate directive to regenerate the mock. (default false)

- `-write_source_comment`: Writes original file (source mode) or interface names (reflect mode) comment if true. (default true)

- `-typed`: Generate Type-safe 'Return', 'Do', 'DoAndReturn' function. (default false)

- `-exclude_interfaces`: Comma-separated names of interfaces to be excluded

For an example of the use of `mockgen`, see the `sample/` directory. In simple
cases, you will need only the `-source` flag.

## Building Mocks

```go
type Foo interface {
  Bar(x int) int
}

func SUT(f Foo) {
 // ...
}

```

```go
func TestFoo(t *testing.T) {
  ctrl := gomock.NewController(t)

  m := NewMockFoo(ctrl)

  // Asserts that the first and only call to Bar() is passed 99.
  // Anything else will fail.
  m.
    EXPECT().
    Bar(gomock.Eq(99)).
    Return(101)

  SUT(m)
}
```

## Building Stubs

```go
type Foo interface {
  Bar(x int) int
}

func SUT(f Foo) {
 // ...
}

```

```go
func TestFoo(t *testing.T) {
  ctrl := gomock.NewController(t)

  m := NewMockFoo(ctrl)

  // Does not make any assertions. Executes the anonymous functions and returns
  // its result when Bar is invoked with 99.
  m.
    EXPECT().
    Bar(gomock.Eq(99)).
    DoAndReturn(func(_ int) int {
      time.Sleep(1*time.Second)
      return 101
    }).
    AnyTimes()

  // Does not make any assertions. Returns 103 when Bar is invoked with 101.
  m.
    EXPECT().
    Bar(gomock.Eq(101)).
    Return(103).
    AnyTimes()

  SUT(m)
}
```

## Modifying Failure Messages

When a matcher reports a failure, it prints the received (`Got`) vs the
expected (`Want`) value.

```shell
Got: [3]
Want: is equal to 2
Expected call at user_test.go:33 doesn't match the argument at index 1.
Got: [0 1 1 2 3]
Want: is equal to 1
```

### Modifying `Want`

The `Want` value comes from the matcher's `String()` method. If the matcher's
default output doesn't meet your needs, then it can be modified as follows:

```go
gomock.WantFormatter(
  gomock.StringerFunc(func() string { return "is equal to fifteen" }),
  gomock.Eq(15),
)
```

This modifies the `gomock.Eq(15)` matcher's output for `Want:` from `is equal
to 15` to `is equal to fifteen`.

### Modifying `Got`

The `Got` value comes from the object's `String()` method if it is available.
In some cases the output of an object is difficult to read (e.g., `[]byte`) and
it would be helpful for the test to print it differently. The following
modifies how the `Got` value is formatted:

```go
gomock.GotFormatterAdapter(
  gomock.GotFormatterFunc(func(i any) string {
    // Leading 0s
    return fmt.Sprintf("%02d", i)
  }),
  gomock.Eq(15),
)
```

If the received value is `3`, then it will be printed as `03`.

[golang]:              http://go.dev/
[ci-badge]:            https://github.com/uber-go/mock/actions/workflows/test.yaml/badge.svg
[ci-runs]:             https://github.com/uber-go/mock/actions
[reference-badge]:     https://pkg.go.dev/badge/go.uber.org/mock.svg
[reference]:           https://pkg.go.dev/go.uber.org/mock