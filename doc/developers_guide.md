# Developers Guide

## Building

To build the binary, run

```shell
go build -o extension-manager main.go
```

## Testing

The different components of the project are responsible for testing different
things.

### extension-manager

The extension-manager project contains unit and integration tests that verify
* Loading and executing of JavaScript extensions
* Database interactions
* REST API interface
* Server-side parameter validation using `extension-parameter-validator`
* ...

Tests use dummy extensions, no real extensions.

### Extensions

Extensions are located in the repositories of the virtual schema implementations, e.g. `s3-document-files-virtual-schema`.

Tests for extensions are:
* Verify correct implementation of a specific version of the
  `extension-manager-interface` using the TypeScript compiler
* Unit tests written in TypeScript verify all execution paths of the extension
* Integration tests written in Java use a specific version of the
  `extension-manager` to verify that the extension
  * can be loaded
  * can install a virtual schema and check that it works
  * can update parameters of an existing virtual schema
  * can upgrade a virtual schema created with an older version
  * ...

### Restrictions as Document Virtual Schemas Only Support a Single Version

Document virtual schemas like `s3-document-files-virtual-schema` require a
`SET SCRIPT` that must have a specific name. As this script references a
specific virtual schema JAR archive, it is not possible to install multiple
version of the same virtual schema in the same database `SCHEMA`.

This means that in order to test a new version of a virtual schema, you need
to create a new `SCHEMA` with the required database objects.

### Non-Parallel Tests

The tests of this project use the exasol-test-setup-abstraction-server. There
the tests connect to an Exasol database running in a docker container.  For
performance reasons the test-setup-abstraction reuses that container.  This
feature is not compatible with running tests in parallel.

Problems would be:

* Name conflicts, e.g. schema names
* Missing isolation, e.g. (`EXA_ALL_SCRIPTS`) contains objects from other tests
* Issues with the exasol-test-setup-abstraction-server (the download of the
  server jar is triggered by the first test. The second one tries to use the
  unfinished jar)

For that reason parallel tests are currently disabled in the CI with `-p 1`.

To run test locally use:

```shell
go test -p 1 ./...
```

To run only unit tests use:

```shell
go test -short ./...
```

## Linter

To install golangci-lint on your machine, follow [these
instruction](https://golangci-lint.run/usage/install/#local-installation). Then
run

```shell
golangci-lint run
```

## Generate API Documentation

The developers generate Swagger API documentation and check it in at
[generatedApiDocs/](../generatedApiDocs/). To update the documentation
first install the `swag` command:

```shell
go install github.com/swaggo/swag/cmd/swag@v1.8.4
```

Make sure to use the same version as specified in
[.github/workflows/ci-build.yml](../.github/workflows/ci-build.yml).

After making changes to the API follow these steps:

1. Run `go generate`
2. Commit changes in the [generatedApiDocs/](../generatedApiDocs/) directory

## Using a Local Extension Interface

To use a local, non-published version of the extension interface in
integration tests, edit
[integrationTesting/extensionForTesting/package.json](../integrationTesting/extensionForTesting/package.json)
and replace the version of `"@exasol/extension-manager-interface"` with the
path to your local clone of
[extension-manager-interface](https://github.com/exasol/extension-manager-interface).

Make sure to not commit the modified `package.json`.
