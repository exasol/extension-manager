# Extension Manager 0.5.17, released 2025-??-??

Code name:

## Summary

## Features

* ISSUE_NUMBER: description

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `golang:1.23.0` to `1.24.9`
* Updated `golang.org/x/mod:v0.24.0` to `v0.31.0`
* Updated `github.com/dop251/goja_nodejs:v0.0.0-20250309172600-86a40d630cdd` to `v0.0.0-20251015164255-5e94316bedaf`
* Updated `github.com/exasol/exasol-driver-go:v1.0.12` to `v1.0.15`
* Updated `github.com/go-chi/chi/v5:v5.2.1` to `v5.2.3`
* Updated `github.com/stretchr/testify:v1.10.0` to `v1.11.1`
* Updated `github.com/dop251/goja:v0.0.0-20250309171923-bcd7cc6bf64c` to `v0.0.0-20251201205617-2bb4c724c0f9`
* Updated `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v0.3.11` to `v1.0.0`

#### Other Dependency Updates

* Updated `toolchain:go1.23.6` to `go1.25.0`

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `com.fasterxml.jackson.core:jackson-annotations:2.18.2` to `2.20`
* Updated `com.fasterxml.jackson.core:jackson-core:2.18.2` to `2.20.1`
* Updated `com.fasterxml.jackson.core:jackson-databind:2.18.2` to `2.20.1`
* Updated `io.swagger.core.v3:swagger-annotations:2.2.28` to `2.2.41`
* Updated `org.glassfish.jersey.core:jersey-client:2.45` to `2.47`
* Updated `org.glassfish.jersey.inject:jersey-hk2:2.45` to `2.47`
* Updated `org.glassfish.jersey.media:jersey-media-json-jackson:2.45` to `2.47`
* Updated `org.glassfish.jersey.media:jersey-media-multipart:2.45` to `2.47`

#### Test Dependency Updates

* Updated `nl.jqno.equalsverifier:equalsverifier:3.19` to `3.19.4`
* Updated `org.junit.jupiter:junit-jupiter-api:5.10.2` to `5.14.1`
* Updated `org.junit.jupiter:junit-jupiter-params:5.10.2` to `5.14.1`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:2.0.3` to `2.0.5`
* Updated `com.exasol:quality-summarizer-maven-plugin:0.2.0` to `0.2.1`
* Added `io.github.git-commit-id:git-commit-id-maven-plugin:9.0.2`
* Removed `io.github.zlika:reproducible-build-maven-plugin:0.17`
* Updated `io.swagger.codegen.v3:swagger-codegen-maven-plugin:3.0.64` to `3.0.75`
* Added `org.apache.maven.plugins:maven-artifact-plugin:3.6.1`
* Updated `org.apache.maven.plugins:maven-clean-plugin:3.4.0` to `3.5.0`
* Updated `org.apache.maven.plugins:maven-compiler-plugin:3.13.0` to `3.14.1`
* Updated `org.apache.maven.plugins:maven-deploy-plugin:3.1.3` to `3.1.4`
* Updated `org.apache.maven.plugins:maven-enforcer-plugin:3.5.0` to `3.6.2`
* Updated `org.apache.maven.plugins:maven-gpg-plugin:3.2.7` to `3.2.8`
* Updated `org.apache.maven.plugins:maven-install-plugin:3.1.3` to `3.1.4`
* Updated `org.apache.maven.plugins:maven-javadoc-plugin:3.11.1` to `3.12.0`
* Updated `org.apache.maven.plugins:maven-resources-plugin:3.3.1` to `3.4.0`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.5.2` to `3.5.4`
* Updated `org.codehaus.mojo:build-helper-maven-plugin:3.6.0` to `3.6.1`
* Updated `org.codehaus.mojo:exec-maven-plugin:3.5.0` to `3.6.2`
* Updated `org.codehaus.mojo:flatten-maven-plugin:1.6.0` to `1.7.3`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.18.0` to `2.20.1`
* Updated `org.jacoco:jacoco-maven-plugin:0.8.12` to `0.8.14`
* Updated `org.sonarsource.scanner.maven:sonar-maven-plugin:5.0.0.4389` to `5.5.0.6356`
* Added `org.sonatype.central:central-publishing-maven-plugin:0.9.0`
* Removed `org.sonatype.plugins:nexus-staging-maven-plugin:1.7.0`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:exasol-test-setup-abstraction-java:2.1.7` to `2.1.10`
* Updated `com.exasol:extension-manager-client-java:0.5.16` to `0.5.17`
* Updated `com.exasol:hamcrest-resultset-matcher:1.7.0` to `1.7.2`
* Updated `com.exasol:test-db-builder-java:3.6.0` to `3.6.4`
* Updated `org.junit.jupiter:junit-jupiter-api:5.10.2` to `5.14.1`

#### Test Dependency Updates

* Updated `com.exasol:maven-project-version-getter:1.2.1` to `1.2.2`
* Updated `com.exasol:udf-debugging-java:0.6.15` to `0.6.18`
* Updated `org.junit.jupiter:junit-jupiter-params:5.10.2` to `5.14.1`
* Updated `org.mockito:mockito-junit-jupiter:5.15.2` to `5.21.0`
* Updated `org.slf4j:slf4j-jdk14:2.0.16` to `2.0.17`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:2.0.3` to `2.0.5`
* Updated `com.exasol:quality-summarizer-maven-plugin:0.2.0` to `0.2.1`
* Added `io.github.git-commit-id:git-commit-id-maven-plugin:9.0.2`
* Removed `io.github.zlika:reproducible-build-maven-plugin:0.17`
* Added `org.apache.maven.plugins:maven-artifact-plugin:3.6.1`
* Updated `org.apache.maven.plugins:maven-clean-plugin:3.4.0` to `3.5.0`
* Updated `org.apache.maven.plugins:maven-compiler-plugin:3.13.0` to `3.14.1`
* Updated `org.apache.maven.plugins:maven-deploy-plugin:3.1.3` to `3.1.4`
* Updated `org.apache.maven.plugins:maven-enforcer-plugin:3.5.0` to `3.6.2`
* Updated `org.apache.maven.plugins:maven-failsafe-plugin:3.5.2` to `3.5.4`
* Updated `org.apache.maven.plugins:maven-gpg-plugin:3.2.7` to `3.2.8`
* Updated `org.apache.maven.plugins:maven-install-plugin:3.1.3` to `3.1.4`
* Updated `org.apache.maven.plugins:maven-jar-plugin:3.4.2` to `3.5.0`
* Updated `org.apache.maven.plugins:maven-javadoc-plugin:3.11.1` to `3.12.0`
* Updated `org.apache.maven.plugins:maven-resources-plugin:3.3.1` to `3.4.0`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.5.2` to `3.5.4`
* Updated `org.codehaus.mojo:flatten-maven-plugin:1.6.0` to `1.7.3`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.18.0` to `2.20.1`
* Updated `org.jacoco:jacoco-maven-plugin:0.8.12` to `0.8.14`
* Updated `org.sonarsource.scanner.maven:sonar-maven-plugin:5.0.0.4389` to `5.5.0.6356`
* Added `org.sonatype.central:central-publishing-maven-plugin:0.9.0`
* Removed `org.sonatype.plugins:nexus-staging-maven-plugin:1.7.0`

### Registry

#### Compile Dependency Updates

* Updated `aws-cdk-lib:2.182.0` to `2.233.0`
* Updated `constructs:^10.4.2` to `^10.4.4`

#### Development Dependency Updates

* Updated `@types/node:^22.13.10` to `^25.0.3`
* Updated `ts-jest:^29.2.6` to `^29.4.6`
* Updated `@types/jest:^29.5.14` to `^30.0.0`
* Updated `typescript:~5.8.2` to `~5.9.3`
* Updated `aws-cdk:2.1003.0` to `2.1100.1`
* Updated `jest:^29.7.0` to `^30.2.0`

### Registry-upload

#### Compile Dependency Updates

* Updated `@aws-sdk/client-cloudfront:^3.764.0` to `^3.955.0`
* Updated `follow-redirects:^1.15.9` to `^1.15.11`
* Updated `@aws-sdk/client-s3:^3.758.0` to `^3.955.0`
* Updated `octokit:^4.1.2` to `^5.0.5`
* Updated `@aws-sdk/client-cloudformation:^3.758.0` to `^3.955.0`

#### Development Dependency Updates

* Updated `eslint:9.22.0` to `9.39.2`
* Updated `@types/node:^22.13.10` to `^25.0.3`
* Updated `typescript-eslint:^8.26.0` to `^8.50.0`
* Updated `typescript:~5.8.2` to `~5.9.3`

### ParameterValidator

#### Development Dependency Updates

* Updated `typescript:5.7.3` to `5.9.3`
* Updated `esbuild:0.25.0` to `0.27.2`
