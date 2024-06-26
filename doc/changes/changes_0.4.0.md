# Extension Manager 0.4.0, released 2023-07-03

Code name: Extension Categories

## Summary

This release adds a field for a category to the extension. This allows extensions to specify a category like "virtual schema" or "cloud storage".

Besides its actual implementation EM also provides the ITF &mdash; a framework for integration tests of concrete extensions. The included ITF has the same version as the EM itself.

In the past when a concrete extension added the ITF to its POM file with version _v_ then the ITF ran integration tests using EM's latest version which could differ from _v_.

EM's current release fixes ITF to use the same version of EM as specified in the extension's POM file.

## Features

* #100: Added category field
* #77: Updated integration test framework to use the current EM version for tests

## Bugfixes

* #107: Adapted EM to work with Exasol 8

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `github.com/dop251/goja:v0.0.0-20230621100801-7749907a8a20` to `v0.0.0-20230626124041-ba8a63e79201`
* Updated `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v0.3.2` to `v0.3.3`

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `io.swagger.core.v3:swagger-annotations:2.2.12` to `2.2.14`

#### Test Dependency Updates

* Updated `nl.jqno.equalsverifier:equalsverifier:3.14.2` to `3.14.3`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.3.0` to `0.4.0`

#### Test Dependency Updates

* Added `com.exasol:maven-project-version-getter:1.2.0`

#### Plugin Dependency Updates

* Updated `org.apache.maven.plugins:maven-jar-plugin:2.4` to `3.3.0`

### Registry

#### Compile Dependency Updates

* Updated `aws-cdk-lib:2.72.1` to `2.86.0`
* Updated `constructs:^10.1.300` to `^10.2.64`

#### Development Dependency Updates

* Updated `@types/node:18.15.11` to `20.3.2`
* Updated `@types/jest:^29.5.0` to `^29.5.2`
* Updated `typescript:~5.0.3` to `~5.1.6`
* Updated `@types/prettier:2.7.2` to `2.7.3`
* Updated `aws-cdk:2.72.1` to `2.86.0`
