# Extension Manager 0.5.9, released 2024-??-??

Code name: Speedup listing extensions

## Summary

This release speeds up listing extensions and installations by caching the extension registry.

Additionally the release fixes a vulnerability in transitive test dependency `io.netty:netty-codec-http:jar:4.1.107.Final` by updating dependencies.

## Features

* #169: Enabled caching for http registry
* #167: Added script for automatically generating extension registry

## Security

* #172: Fix CVE-2024-29025 in transitive test dependency `io.netty:netty-codec-http:jar:4.1.107.Final`
## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `golang:1.20` to `1.21`
* Updated `github.com/exasol/exasol-driver-go:v1.0.4` to `v1.0.6`
* Updated `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v0.3.5` to `v0.3.6`

#### Test Dependency Updates

* Updated `golang.org/x/mod:v0.16.0` to `v0.17.0`

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `com.fasterxml.jackson.core:jackson-annotations:2.16.2` to `2.17.0`
* Updated `com.fasterxml.jackson.core:jackson-core:2.16.2` to `2.17.0`
* Updated `com.fasterxml.jackson.core:jackson-databind:2.16.2` to `2.17.0`
* Updated `io.swagger.core.v3:swagger-annotations:2.2.20` to `2.2.21`

#### Test Dependency Updates

* Updated `nl.jqno.equalsverifier:equalsverifier:3.15.8` to `3.16.1`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:2.0.0` to `2.0.1`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:exasol-test-setup-abstraction-java:2.1.1` to `2.1.2`
* Updated `com.exasol:extension-manager-client-java:0.5.8` to `0.5.10`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:2.0.0` to `2.0.1`

### Registry

#### Compile Dependency Updates

* Updated `aws-cdk-lib:2.132.1` to `2.133.0`

#### Development Dependency Updates

* Updated `@types/node:^20.11.26` to `^20.11.30`
* Updated `typescript:~5.4.2` to `~5.4.3`
* Updated `aws-cdk:2.132.1` to `2.133.0`

### Registry-upload

#### Compile Dependency Updates

* Updated `@aws-sdk/client-cloudfront:^3.451.0` to `^3.535.0`
* Updated `follow-redirects:^1.15.3` to `^1.15.6`
* Updated `@aws-sdk/client-s3:^3.451.0` to `^3.537.0`
* Added `octokit:^3.1.2`
* Updated `@aws-sdk/client-cloudformation:^3.451.0` to `^3.537.0`

#### Development Dependency Updates

* Updated `eslint:^8.53.0` to `^8.57.0`
* Updated `@typescript-eslint/parser:^6.11.0` to `^7.3.1`
* Updated `@types/node:^20.9.1` to `^20.11.30`
* Updated `typescript:~5.2.2` to `~5.4.3`
* Updated `@typescript-eslint/eslint-plugin:^6.11.0` to `^7.3.1`
* Updated `ts-node:^10.9.1` to `^10.9.2`
