# Extension Manager 0.5.3, released 2023-??-??

Code name: Speedup listing extensions

## Summary

This release speeds up listing extensions, especially when there are many files in BucketFS.

**Notes:** Starting with this release EM is tested against Exasol version 8 instead of 7.1. This means that integration tests using `extension-manager-integration-test-java` will also need to run with Exasol 8.

## Bugfix

* #147: Improved speed of listing available extensions

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `github.com/exasol/exasol-driver-go:v1.0.2` to `v1.0.3`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.5.2` to `0.5.3`
