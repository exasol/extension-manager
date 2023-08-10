# Extension Manager 0.5.0, released 2023-08-10

Code name: Upgrade Extensions

## Summary

This release supports upgrading installed extensions to their latest version. Extensions must implement the latest `extension-manager-interface` version 0.3.0 to support this.

This release also allows configuring the BucketFS base path where EM expects extension files to be located. EM searches this path recursively, so files are also found in subdirectories.

This release also improves error handling when using extensions not implementing all functions required by EM. EM now returns a helpful error message instead of failing with a `nil`-pointer error.

A common scenario for an extension not implementing a required function is when the extension had been built using an older version of EM's extension interface.

## Features

* #101 Added support for upgrading installed extensions

## Bugfixes

* #105: Ensured that EM can load and use compatible extensions
* #74: Fixed most important linter warnings

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `github.com/dop251/goja:v0.0.0-20230626124041-ba8a63e79201` to `v0.0.0-20230707174833-636fdf960de1`
* Updated `github.com/go-chi/chi/v5:v5.0.8` to `v5.0.10`

#### Test Dependency Updates

* Updated `golang.org/x/mod:v0.11.0` to `v0.12.0`

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `io.swagger.core.v3:swagger-annotations:2.2.14` to `2.2.15`
* Updated `org.glassfish.jersey.core:jersey-client:2.39.1` to `2.40`
* Updated `org.glassfish.jersey.inject:jersey-hk2:2.39.1` to `2.40`
* Updated `org.glassfish.jersey.media:jersey-media-json-jackson:2.39.1` to `2.40`
* Updated `org.glassfish.jersey.media:jersey-media-multipart:2.39.1` to `2.40`

#### Test Dependency Updates

* Updated `nl.jqno.equalsverifier:equalsverifier:3.14.3` to `3.15`
* Updated `org.junit.jupiter:junit-jupiter-api:5.9.3` to `5.10.0`
* Added `org.junit.jupiter:junit-jupiter-params:5.10.0`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:1.2.3` to `1.3.0`
* Updated `com.exasol:project-keeper-maven-plugin:2.9.7` to `2.9.10`
* Updated `org.apache.maven.plugins:maven-gpg-plugin:3.0.1` to `3.1.0`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.0.0` to `3.1.2`
* Updated `org.basepom.maven:duplicate-finder-maven-plugin:1.5.1` to `2.0.1`
* Updated `org.codehaus.mojo:build-helper-maven-plugin:3.3.0` to `3.4.0`
* Updated `org.codehaus.mojo:flatten-maven-plugin:1.4.1` to `1.5.0`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.15.0` to `2.16.0`
* Updated `org.jacoco:jacoco-maven-plugin:0.8.9` to `0.8.10`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:exasol-test-setup-abstraction-java:2.0.1` to `2.0.2`
* Updated `com.exasol:extension-manager-client-java:0.4.0` to `0.5.0`
* Removed `io.netty:netty-handler:4.1.94.Final`
* Updated `org.junit.jupiter:junit-jupiter-api:5.9.3` to `5.10.0`

#### Test Dependency Updates

* Updated `com.exasol:udf-debugging-java:0.6.8` to `0.6.10`
* Updated `org.junit.jupiter:junit-jupiter-params:5.9.3` to `5.10.0`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:1.2.3` to `1.3.0`
* Updated `com.exasol:project-keeper-maven-plugin:2.9.7` to `2.9.10`
* Updated `org.apache.maven.plugins:maven-failsafe-plugin:3.0.0` to `3.1.2`
* Updated `org.apache.maven.plugins:maven-gpg-plugin:3.0.1` to `3.1.0`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.0.0` to `3.1.2`
* Updated `org.basepom.maven:duplicate-finder-maven-plugin:1.5.1` to `2.0.1`
* Updated `org.codehaus.mojo:flatten-maven-plugin:1.4.1` to `1.5.0`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.15.0` to `2.16.0`
* Updated `org.jacoco:jacoco-maven-plugin:0.8.9` to `0.8.10`

### Registry

#### Compile Dependency Updates

* Updated `aws-cdk-lib:2.86.0` to `2.87.0`
* Updated `constructs:^10.2.64` to `^10.2.69`

#### Development Dependency Updates

* Updated `@types/node:20.3.2` to `20.4.2`
* Updated `ts-jest:^29.1.0` to `^29.1.1`
* Updated `@types/jest:^29.5.2` to `^29.5.3`
* Updated `jest:^29.5.0` to `^29.6.1`
* Updated `aws-cdk:2.86.0` to `2.87.0`

### ParameterValidator

#### Compile Dependency Updates

* Updated `@exasol/extension-parameter-validator:0.2.0` to `0.2.1`

#### Development Dependency Updates

* Updated `typescript:5.0.3` to `5.1.6`
* Updated `esbuild:0.17.15` to `0.18.16`
