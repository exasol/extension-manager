# Extension Manager 0.5.13, released 2024-11-19

Code name: Fix CVE-2024-47535 in `io.netty:netty-common:jar:4.1.108.Final:runtime`

## Summary

This release fixes vulnerability CVE-2024-47535 in `io.netty:netty-common:jar:4.1.108.Final:runtime`.

## Security

* #184: Fixed CVE-2024-47535 in `io.netty:netty-common:jar:4.1.108.Final:runtime`

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `golang:1.21` to `1.22.0`
* Updated `github.com/dop251/goja:v0.0.0-20240610225006-393f6d42497b` to `v0.0.0-20241024094426-79f3a7efcdbd`
* Updated `github.com/dop251/goja_nodejs:v0.0.0-20240418154818-2aae10d4cbcf` to `v0.0.0-20240728170619-29b559befffc`
* Updated `github.com/exasol/exasol-driver-go:v1.0.7` to `v1.0.10`
* Updated `github.com/go-chi/chi/v5:v5.0.12` to `v5.1.0`

#### Test Dependency Updates

* Updated `golang.org/x/mod:v0.18.0` to `v0.22.0`
* Updated `github.com/kinbiko/jsonassert:v1.1.1` to `v1.2.0`

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `com.fasterxml.jackson.core:jackson-annotations:2.17.1` to `2.18.1`
* Updated `com.fasterxml.jackson.core:jackson-core:2.17.1` to `2.18.1`
* Updated `com.fasterxml.jackson.core:jackson-databind:2.17.1` to `2.18.1`
* Updated `io.swagger.core.v3:swagger-annotations:2.2.22` to `2.2.26`
* Updated `org.glassfish.jersey.core:jersey-client:2.41` to `2.45`
* Updated `org.glassfish.jersey.inject:jersey-hk2:2.41` to `2.45`
* Updated `org.glassfish.jersey.media:jersey-media-json-jackson:2.41` to `2.45`
* Updated `org.glassfish.jersey.media:jersey-media-multipart:2.41` to `2.45`

#### Test Dependency Updates

* Updated `nl.jqno.equalsverifier:equalsverifier:3.16.1` to `3.17.3`

#### Plugin Dependency Updates

* Added `com.exasol:quality-summarizer-maven-plugin:0.2.0`
* Updated `io.github.zlika:reproducible-build-maven-plugin:0.16` to `0.17`
* Updated `io.swagger.codegen.v3:swagger-codegen-maven-plugin:3.0.54` to `3.0.64`
* Updated `org.apache.maven.plugins:maven-clean-plugin:2.5` to `3.4.0`
* Updated `org.apache.maven.plugins:maven-gpg-plugin:3.2.4` to `3.2.7`
* Updated `org.apache.maven.plugins:maven-install-plugin:2.4` to `3.1.3`
* Updated `org.apache.maven.plugins:maven-javadoc-plugin:3.7.0` to `3.10.1`
* Updated `org.apache.maven.plugins:maven-resources-plugin:2.6` to `3.3.1`
* Updated `org.apache.maven.plugins:maven-site-plugin:3.3` to `3.9.1`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.2.5` to `3.5.1`
* Updated `org.codehaus.mojo:build-helper-maven-plugin:3.5.0` to `3.6.0`
* Updated `org.codehaus.mojo:exec-maven-plugin:3.2.0` to `3.5.0`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.16.2` to `2.17.1`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:exasol-test-setup-abstraction-java:2.1.4` to `2.1.6`
* Updated `com.exasol:extension-manager-client-java:0.5.12` to `0.5.13`
* Updated `com.exasol:hamcrest-resultset-matcher:1.6.5` to `1.7.0`
* Updated `com.exasol:test-db-builder-java:3.5.4` to `3.6.0`

#### Test Dependency Updates

* Updated `org.mockito:mockito-junit-jupiter:5.12.0` to `5.14.2`
* Updated `org.slf4j:slf4j-jdk14:2.0.13` to `2.0.16`

#### Plugin Dependency Updates

* Added `com.exasol:quality-summarizer-maven-plugin:0.2.0`
* Updated `io.github.zlika:reproducible-build-maven-plugin:0.16` to `0.17`
* Updated `org.apache.maven.plugins:maven-clean-plugin:2.5` to `3.4.0`
* Updated `org.apache.maven.plugins:maven-failsafe-plugin:3.2.5` to `3.5.1`
* Updated `org.apache.maven.plugins:maven-gpg-plugin:3.2.4` to `3.2.7`
* Updated `org.apache.maven.plugins:maven-install-plugin:2.4` to `3.1.3`
* Updated `org.apache.maven.plugins:maven-jar-plugin:3.3.0` to `3.4.2`
* Updated `org.apache.maven.plugins:maven-javadoc-plugin:3.7.0` to `3.10.1`
* Updated `org.apache.maven.plugins:maven-resources-plugin:2.6` to `3.3.1`
* Updated `org.apache.maven.plugins:maven-site-plugin:3.3` to `3.9.1`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.2.5` to `3.5.1`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.16.2` to `2.17.1`

### Registry

#### Compile Dependency Updates

* Updated `aws-cdk-lib:2.146.0` to `2.167.1`
* Updated `constructs:^10.3.0` to `^10.4.2`

#### Development Dependency Updates

* Updated `@types/node:^20.14.2` to `^22.9.0`
* Updated `ts-jest:^29.1.4` to `^29.2.5`
* Updated `@types/jest:^29.5.12` to `^29.5.14`
* Updated `typescript:~5.4.5` to `~5.6.3`
* Updated `aws-cdk:2.146.0` to `2.167.1`

### Registry-upload

#### Compile Dependency Updates

* Updated `@aws-sdk/client-cloudfront:^3.596.0` to `^3.693.0`
* Updated `follow-redirects:^1.15.6` to `^1.15.9`
* Updated `@aws-sdk/client-s3:^3.596.0` to `^3.693.0`
* Updated `@aws-sdk/client-cloudformation:^3.596.0` to `^3.695.0`

#### Development Dependency Updates

* Updated `eslint:^8.57.0` to `9.14.0`
* Updated `@types/node:^20.14.2` to `^22.9.0`
* Added `typescript-eslint:^8.14.0`
* Updated `typescript:~5.4.5` to `~5.6.3`
* Removed `@typescript-eslint/parser:^7.13.0`
* Removed `@typescript-eslint/eslint-plugin:^7.13.0`

### ParameterValidator

#### Development Dependency Updates

* Updated `typescript:5.4.5` to `5.6.3`
* Updated `esbuild:0.21.5` to `0.24.0`
