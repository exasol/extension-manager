# Extension Manager 0.2.0, released 2022-11-07

Code name: Add extension registry

## Summary

This release adds a CDK stack for deploying the infrastructure of the Extension Registry to AWS.

We also moved all Go sources to the `pkg` directory. Projects that use this library will need to adapt the imports by replacing `"github.com/exasol/extension-manager/*"` with `"github.com/exasol/extension-manager/pkg/*"`.

## Features

* #80: Added prefix to log messages from JS extensions
* #82: Added infrastructure for extension registry

## Refactoring

* #86: Moved Go sources to `pkg` directory

## Documentation

* #73: Updated requirements and added requirements tracing

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `github.com/stretchr/testify:v1.8.0` to `v1.8.1`
* Updated `github.com/dop251/goja:v0.0.0-20220906144433-c4d370b87b45` to `v0.0.0-20221106173738-3b8a68ca89b4`
* Updated `github.com/dop251/goja_nodejs:v0.0.0-20220905124449-678b33ca5009` to `v0.0.0-20221009164102-3aa5028e57f6`
* Updated `github.com/exasol/exasol-driver-go:v0.4.5` to `v0.4.6`
* Updated `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v0.2.4` to `v0.3.0`

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `com.fasterxml.jackson.core:jackson-annotations:2.14.0-rc1` to `2.14.0`
* Updated `com.fasterxml.jackson.core:jackson-core:2.14.0-rc1` to `2.14.0`
* Updated `com.fasterxml.jackson.core:jackson-databind:2.14.0-rc1` to `2.14.0`
* Updated `io.swagger.core.v3:swagger-annotations:2.2.3` to `2.2.6`
* Updated `org.glassfish.jersey.core:jersey-client:2.36` to `2.37`
* Updated `org.glassfish.jersey.inject:jersey-hk2:2.36` to `2.37`
* Updated `org.glassfish.jersey.media:jersey-media-json-jackson:2.36` to `2.37`
* Updated `org.glassfish.jersey.media:jersey-media-multipart:2.36` to `2.37`

#### Test Dependency Updates

* Updated `nl.jqno.equalsverifier:equalsverifier:3.10.1` to `3.11`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:1.1.2` to `1.2.1`
* Updated `com.exasol:project-keeper-maven-plugin:2.8.0` to `2.9.1`
* Updated `io.github.zlika:reproducible-build-maven-plugin:0.15` to `0.16`
* Updated `org.apache.maven.plugins:maven-deploy-plugin:3.0.0-M1` to `3.0.0`
* Updated `org.apache.maven.plugins:maven-javadoc-plugin:3.4.0` to `3.4.1`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.0.0-M5` to `3.0.0-M7`
* Updated `org.codehaus.mojo:flatten-maven-plugin:1.2.7` to `1.3.0`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.10.0` to `2.13.0`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:exasol-test-setup-abstraction-java:0.3.2` to `1.0.0`
* Updated `com.exasol:extension-manager-client-java:0.1.0` to `0.2.0`
* Updated `com.exasol:test-db-builder-java:3.4.0` to `3.4.1`

#### Test Dependency Updates

* Updated `org.mockito:mockito-junit-jupiter:4.8.0` to `4.8.1`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:1.1.2` to `1.2.1`
* Updated `com.exasol:project-keeper-maven-plugin:2.8.0` to `2.9.1`
* Updated `io.github.zlika:reproducible-build-maven-plugin:0.15` to `0.16`
* Updated `org.apache.maven.plugins:maven-deploy-plugin:3.0.0-M1` to `3.0.0`
* Updated `org.apache.maven.plugins:maven-failsafe-plugin:3.0.0-M5` to `3.0.0-M7`
* Updated `org.apache.maven.plugins:maven-javadoc-plugin:3.4.0` to `3.4.1`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.0.0-M5` to `3.0.0-M7`
* Updated `org.codehaus.mojo:flatten-maven-plugin:1.2.7` to `1.3.0`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.10.0` to `2.13.0`

### Registry

#### Compile Dependency Updates

* Added `source-map-support:^0.5.21`
* Added `aws-cdk-lib:2.50.0`
* Added `constructs:^10.1.148`

#### Development Dependency Updates

* Added `@types/node:18.11.9`
* Added `ts-jest:^29.0.3`
* Added `@types/jest:^29.2.1`
* Added `typescript:~4.8.4`
* Added `@types/prettier:2.7.1`
* Added `jest:^29.2.2`
* Added `aws-cdk:2.50.0`
* Added `ts-node:^10.9.1`

### ParameterValidator

#### Compile Dependency Updates

* Added `@exasol/extension-parameter-validator:0.2.0`

#### Development Dependency Updates

* Added `typescript:4.8.4`
* Added `esbuild:0.15.13`
