# Developers Guide

## Non-Parallel Tests

The tests of this project use the exasol-test-setup-abstraction-server. There the tests connect to an Exasol database
running in a docker container.
For performance reasons the test-setup-abstraction reuses that container.
This feature is not compatible with running tests in parallel.

Problems would be:

* Name conflicts, e.g. schema names
* Missing isolation, e.g. (`EXA_ALL_SCRIPTS`) contains objects from other tests
* Issues with the exasol-test-setup-abstraction-server (the download of the server jar is triggerd by the first test.
  The second one tries to use the unfinished jar)

For that reason we currently disabled parallel tests in the CI with `-p 1`.

To run test locally use:

```shell
go test -p 1 ./...
```

To run only unit tests use:

```shell
go test -short ./...
```

### Linter

To install golangci-lint on your machine, follow [these instruction](https://golangci-lint.run/usage/install/#local-installation). Then run

```shell
golangci-lint run
```
### Generate API Documentation

We generate Swagger API documentation and check it in at [generatedApiDocs/](../../generatedApiDocs/). To update the documentation first install the `swag` command:

```shell
go install github.com/swaggo/swag/cmd/swag@v1.8.2
```

Make sure to use the same version as specified in [.github/workflows/ci-build.yml](../../.github/workflows/ci-build.yml).

After making changes to the API follow these steps:

1. Run `go generate`
2. Commit changes in the [generatedApiDocs/](../../generatedApiDocs/) directory
