# Developers Guide

## Non-Parallel Tests

The tests of this project use the exasol-test-setup-abstraction-server. There the tests connect to an Exasol database
running in a docker container.
For performance reasons the test-setup-abstraction reuses that container.
This feature not compatible with parallel running tests.

Problems would be:

* Name conflicts
* Missing isolation (`EXA_ALL_SCRIPTS`) contains objects from other tests
* Issues with the exasol-test-setup-abstraction-server (the download of the server jar is triggerd by the first test.
  The second one tries to use the unfinished jar)

For that reason we currently disabled parallel tests in the CI.

To run test locally use:

```shell
go run -p 1 ./...
```
