# Extension Manager 0.2.0, released 2022-??-??

Code name:

## Summary

In this release we moved all Go sources to the `pkg` directory. Projects that use this library will need to adapt the imports by replacing `"github.com/exasol/extension-manager/*"` with `"github.com/exasol/extension-manager/pkg/*"`.

## Features

* #80: Added prefix to log messages from JS extensions

## Refactoring

* #86: Moved Go sources to `pkg` directory

## Dependency Updates

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.1.0` to `0.2.0`
