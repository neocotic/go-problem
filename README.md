# go-problem

[![Go Reference](https://img.shields.io/badge/go.dev-reference-007d9c?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/neocotic/go-problem)
[![Build Status](https://img.shields.io/github/actions/workflow/status/neocotic/go-problem/ci.yml?style=for-the-badge)](https://github.com/neocotic/go-problem/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/neocotic/go-problem?style=for-the-badge)](https://github.com/neocotic/go-problem)
[![License](https://img.shields.io/github/license/neocotic/go-problem?style=for-the-badge)](https://github.com/neocotic/go-problem/blob/main/LICENSE.md)

Flexible and customizable Go (golang) implementation of [RFC 9457](https://datatracker.ietf.org/doc/html/rfc9457);
Problem Details for HTTP APIs.

## Installation

Install using [go install](https://go.dev/ref/mod#go-install):

``` sh
go install github.com/neocotic/go-problem
```

Then import the package into your own code:

``` go
import "github.com/neocotic/go-problem"
```

## Documentation

Documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/neocotic/go-problem#section-documentation). It
contains an overview and reference.

### Example

Define reusable problem types and/or definitions:

``` go
var (
	NotFound = problem.Type{
		LogLevel: problem.LogLevelDebug,
		Status:   http.StatusNotFound,
		Title:    http.StatusText(http.StatusNotFound),
	}
	NotFoundDefinition = problem.Definition{
		Type: NotFound,
	}
)
```

Create a problem using the builder pattern:

``` go
problem.Build().
	Definition(NotFoundDefinition).
	Code(problem.MustBuildCode(404, "USER")).
	Detail("User not found").
	Instance("https://api.example.void/users/123").
	Wrap(err).
	Problem()
```

Or the option pattern:

``` go
problem.New(
	problem.FromDefinition(NotFoundDefinition),
	problem.WithCode(problem.MustBuildCode(404, "USER")),
	problem.WithDetail("User not found"),
	problem.WithInstance("https://api.example.void/users/123"),
	problem.Wrap(err),
)
```

There's a load of other functions to explore and `problem.Generator` can be used to customize the experience much
further. It even comes with functions for writing problems and errors to an HTTP response and middleware for panic
recovery with problem responses.

## Issues

If you have any problems or would like to see changes currently in development you can do so
[here](https://github.com/neocotic/go-problem/issues).

## Contributors

If you want to contribute, you're a legend! Information on how you can do so can be found in
[CONTRIBUTING.md](https://github.com/neocotic/go-problem/blob/main/CONTRIBUTING.md). We want your suggestions and pull
requests!

A list of contributors can be found in [AUTHORS.md](https://github.com/neocotic/go-problem/blob/main/AUTHORS.md).

## License

Copyright © 2024 neocotic

See [LICENSE.md](https://github.com/neocotic/go-problem/raw/main/LICENSE.md) for more information on our MIT license.
