# Extension Manager 0.5.5, released 2023-11-06

Code name: Fix misleading error

## Summary

This release fixes a misleading error in the integration test framework. When `PreviousVersionManager.prepareBucketFsFile()` is called with an invalid URL it fails with `NoSuchFileException` file exception that hides the actual exception about the invalid URL. The release also improves other error messages in Extension Manager and fixes a bug that didn't let you uninstall extensions that don't support instances (e.g. because they only require scripts).

The release also adds a base class that simplifies writing integration tests for extensions.

## Features

* #159: Extracted common code for extension integration tests

## Bugfixes

* #156: Fixed misleading error in `PreviousVersionManager.prepareBucketFsFile()`
* #155: Fixed uninstalling extensions that don't support instances

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `github.com/dop251/goja:v0.0.0-20231014103939-873a1496dc8e` to `v0.0.0-20231027120936-b396bb4c349d`
* Updated `github.com/dop251/goja_nodejs:v0.0.0-20230914102007-198ba9a8b098` to `v0.0.0-20231022114343-5c1f9037c9ab`
* Updated `github.com/exasol/exasol-driver-go:v1.0.3` to `v1.0.4`

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `io.swagger.core.v3:swagger-annotations:2.2.17` to `2.2.18`
* Updated `org.glassfish.jersey.core:jersey-client:2.40` to `2.41`
* Updated `org.glassfish.jersey.inject:jersey-hk2:2.40` to `2.41`
* Updated `org.glassfish.jersey.media:jersey-media-json-jackson:2.40` to `2.41`
* Updated `org.glassfish.jersey.media:jersey-media-multipart:2.40` to `2.41`

#### Test Dependency Updates

* Updated `nl.jqno.equalsverifier:equalsverifier:3.15.2` to `3.15.3`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:1.3.0` to `1.3.1`
* Updated `com.exasol:project-keeper-maven-plugin:2.9.12` to `2.9.15`
* Updated `org.apache.maven.plugins:maven-enforcer-plugin:3.4.0` to `3.4.1`
* Updated `org.apache.maven.plugins:maven-javadoc-plugin:3.5.0` to `3.6.0`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.16.0` to `2.16.1`
* Updated `org.jacoco:jacoco-maven-plugin:0.8.10` to `0.8.11`
* Updated `org.sonarsource.scanner.maven:sonar-maven-plugin:3.9.1.2184` to `3.10.0.2594`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.5.4` to `0.5.5`
* Updated `com.exasol:hamcrest-resultset-matcher:1.6.1` to `1.6.2`

#### Test Dependency Updates

* Updated `org.mockito:mockito-junit-jupiter:5.6.0` to `5.7.0`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:1.3.0` to `1.3.1`
* Updated `com.exasol:project-keeper-maven-plugin:2.9.12` to `2.9.15`
* Updated `org.apache.maven.plugins:maven-enforcer-plugin:3.4.0` to `3.4.1`
* Updated `org.apache.maven.plugins:maven-javadoc-plugin:3.5.0` to `3.6.0`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.16.0` to `2.16.1`
* Updated `org.jacoco:jacoco-maven-plugin:0.8.10` to `0.8.11`
* Updated `org.sonarsource.scanner.maven:sonar-maven-plugin:3.9.1.2184` to `3.10.0.2594`

### Registry

#### Compile Dependency Updates

* Updated `aws-cdk-lib:2.101.1` to `2.104.0`

#### Development Dependency Updates

* Updated `@types/node:^20.8.6` to `^20.8.10`
* Updated `@types/jest:^29.5.5` to `^29.5.7`
* Updated `aws-cdk:2.101.1` to `2.104.0`

### Registry-upload

#### Compile Dependency Updates

* Updated `@aws-sdk/client-cloudfront:^3.429.0` to `^3.441.0`
* Updated `@aws-sdk/client-s3:^3.429.0` to `^3.441.0`
* Updated `@aws-sdk/client-cloudformation:^3.429.0` to `^3.441.0`

#### Development Dependency Updates

* Updated `eslint:^8.51.0` to `^8.52.0`
* Updated `@types/follow-redirects:^1.14.2` to `^1.14.3`
* Updated `@typescript-eslint/parser:^6.8.0` to `^6.9.1`
* Updated `@types/node:^20.8.6` to `^20.8.10`
* Updated `@typescript-eslint/eslint-plugin:^6.8.0` to `^6.9.1`

### ParameterValidator

#### Development Dependency Updates

* Updated `esbuild:0.19.2` to `0.19.5`
