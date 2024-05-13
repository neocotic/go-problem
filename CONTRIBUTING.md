# Contributing

If you have any questions about [go-problem](https://github.com/neocotic/go-problem) please feel free to
[raise an issue](https://github.com/neocotic/go-problem/issues/new).

Please [search existing issues](https://github.com/neocotic/go-problem/issues) for the same feature and/or issue before
raising a new issue. Commenting on an existing issue is usually preferred over raising duplicate issues.

Please ensure that all files conform to the coding standards, using the same coding style as the rest of the code base.
All unit tests should be updated and passing as well. All of this can easily be checked via command-line:

``` bash
# install/update package dependencies
$ make download tidy
# format code
$ make format
# run test suite
$ make test
```

You must have at least [Golang](https://go.dev) version 1.21 or newer installed.

All pull requests should be made to the `main` branch.

Don't forget to add your details to the list of
[AUTHORS.md](https://github.com/neocotic/go-problem/blob/main/AUTHORS.md) if you want your contribution to be recognized
by others.
