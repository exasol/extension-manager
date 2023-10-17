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

* Updated `github.com/dop251/goja:v0.0.0-20230919151941-fc55792775de` to `v0.0.0-20231014103939-873a1496dc8e`
* Updated `github.com/exasol/exasol-driver-go:v1.0.2` to `v1.0.3`

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `com.fasterxml.jackson.core:jackson-annotations:2.15.2` to `2.15.3`
* Updated `com.fasterxml.jackson.core:jackson-core:2.15.2` to `2.15.3`
* Updated `com.fasterxml.jackson.core:jackson-databind:2.15.2` to `2.15.3`
* Updated `io.swagger.core.v3:swagger-annotations:2.2.16` to `2.2.17`

#### Plugin Dependency Updates

* Updated `io.swagger.codegen.v3:swagger-codegen-maven-plugin:3.0.46` to `3.0.47`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.5.2` to `0.5.3`

### Registry

#### Compile Dependency Updates

* Updated `aws-cdk-lib:2.95.1` to `2.101.1`
* Updated `constructs:^10.2.70` to `^10.3.0`

#### Development Dependency Updates

* Updated `@types/node:^20.6.0` to `^20.8.6`
* Updated `@types/jest:^29.5.4` to `^29.5.5`
* Updated `aws-cdk:2.95.1` to `2.101.1`
* Updated `jest:^29.6.4` to `^29.7.0`

### Registry-upload

#### Compile Dependency Updates

* Updated `@aws-sdk/client-cloudfront:^3.409.0` to `^3.429.0`
* Updated `follow-redirects:^1.15.2` to `^1.15.3`
* Updated `@aws-sdk/client-s3:^3.409.0` to `^3.429.0`
* Updated `@aws-sdk/client-cloudformation:^3.409.0` to `^3.429.0`

#### Development Dependency Updates

* Updated `eslint:^8.49.0` to `^8.51.0`
* Updated `@types/follow-redirects:^1.14.1` to `^1.14.2`
* Updated `@typescript-eslint/parser:^6.6.0` to `^6.8.0`
* Updated `@types/node:^20.6.0` to `^20.8.6`
* Updated `@typescript-eslint/eslint-plugin:^6.6.0` to `^6.8.0`
