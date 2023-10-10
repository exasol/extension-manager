# Extension Manager 0.5.2, released 2023-10-10

Code name: Small improvements

## Summary

This release contains the following notable changes:

* It updates the upload process for the extension registry to verify that the extension URLs are valid.
* It verifies that no instances exist before uninstalling an extension.
* It adds design, requirements and user guide for the integration testing framework.
* It adds a custom JavaScript console logger that forwards the log levels of `console.log()`, `console.warn()` and `console.error()` to the Go `logrus` functions `Print()`, `Warn()` and `Error()`. Before, all log messages were output using `Print()`.
* It speeds up tests that use a mock extension by avoiding compiling TypeScript to JavaScript.
* Controller method `GetInstalledExtensions()` now additionally returns the extension ID.

The release deprecates function `extensionController.CreateWithConfig()` in favor of `extensionController.CreateWithValidatedConfig()`. This new function validates the given configuration and returns an error in case it finds an issue.

## Features

* #129: Added verification for extension URLs before uploading to registry
* #142: Added validation of configuration when creating a new controller
* #130: Added verification that no instance exists before uninstalling an extension
* #115: Added custom JavaScript console logger that forwards log levels to the Go logger
* #148: Added extension ID to result of controller method `GetInstalledExtensions()`

## Documentation

* #9: Add design, requirements and user guide for integration testing framework

## Refactoring

* #139: Converted `extensionForTesting` to JavaScript to speedup tests
* #114: Fixed many linter warnings

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `golang:1.19` to `1.20`
* Updated `github.com/dop251/goja:v0.0.0-20230707174833-636fdf960de1` to `v0.0.0-20230919151941-fc55792775de`
* Updated `github.com/dop251/goja_nodejs:v0.0.0-20230602164024-804a84515562` to `v0.0.0-20230914102007-198ba9a8b098`
* Updated `github.com/exasol/exasol-driver-go:v1.0.0` to `v1.0.2`

#### Test Dependency Updates

* Updated `golang.org/x/mod:v0.12.0` to `v0.13.0`

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `io.swagger.core.v3:swagger-annotations:2.2.15` to `2.2.16`

#### Test Dependency Updates

* Updated `nl.jqno.equalsverifier:equalsverifier:3.15.1` to `3.15.2`

#### Plugin Dependency Updates

* Updated `com.exasol:project-keeper-maven-plugin:2.9.11` to `2.9.12`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:exasol-test-setup-abstraction-java:2.0.3` to `2.0.4`
* Updated `com.exasol:extension-manager-client-java:0.5.1` to `0.5.2`
* Updated `com.exasol:hamcrest-resultset-matcher:1.6.0` to `1.6.1`
* Updated `com.exasol:test-db-builder-java:3.5.0` to `3.5.1`

#### Test Dependency Updates

* Updated `com.exasol:udf-debugging-java:0.6.10` to `0.6.11`
* Updated `org.mockito:mockito-junit-jupiter:5.5.0` to `5.6.0`

#### Plugin Dependency Updates

* Updated `com.exasol:project-keeper-maven-plugin:2.9.11` to `2.9.12`

### Registry

#### Compile Dependency Updates

* Updated `aws-cdk-lib:2.87.0` to `2.95.1`
* Updated `constructs:^10.2.69` to `^10.2.70`

#### Development Dependency Updates

* Updated `@types/node:20.4.2` to `^20.6.0`
* Updated `@types/jest:^29.5.3` to `^29.5.4`
* Updated `typescript:~5.1.6` to `~5.2.2`
* Updated `aws-cdk:2.87.0` to `2.95.1`
* Updated `jest:^29.6.1` to `^29.6.4`

### Registry-upload

#### Compile Dependency Updates

* Added `@aws-sdk/client-cloudfront:^3.409.0`
* Added `follow-redirects:^1.15.2`
* Added `@aws-sdk/client-s3:^3.409.0`
* Added `@aws-sdk/client-cloudformation:^3.409.0`

#### Development Dependency Updates

* Added `eslint:^8.49.0`
* Added `@types/follow-redirects:^1.14.1`
* Added `@typescript-eslint/parser:^6.6.0`
* Added `@types/node:^20.6.0`
* Added `typescript:~5.2.2`
* Added `@types/prettier:2.7.3`
* Added `@typescript-eslint/eslint-plugin:^6.6.0`
* Added `ts-node:^10.9.1`
