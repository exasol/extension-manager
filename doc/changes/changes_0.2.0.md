# Extension Manager 0.2.0, released 2022-11-02

Code name: Add extension registry

## Summary

In this release we added a CDK stack for deploying the infrastructure of the Extension Registry to AWS.

We also moved all Go sources to the `pkg` directory. Projects that use this library will need to adapt the imports by replacing `"github.com/exasol/extension-manager/*"` with `"github.com/exasol/extension-manager/pkg/*"`.

## Features

* #80: Added prefix to log messages from JS extensions
* #82: Added infrastructure for extension registry

## Refactoring

* #86: Moved Go sources to `pkg` directory

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `github.com/dop251/goja:v0.0.0-20220906144433-c4d370b87b45` to `v0.0.0-20221019153710-09250e0eba20`
* Updated `github.com/dop251/goja_nodejs:v0.0.0-20220905124449-678b33ca5009` to `v0.0.0-20221009164102-3aa5028e57f6`
* Updated `github.com/exasol/exasol-driver-go:v0.4.5` to `v0.4.6`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.1.0` to `0.2.0`
