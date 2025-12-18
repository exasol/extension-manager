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

* Updated `com.exasol:extension-manager-client-java:0.5.16` to `0.5.17`
* Updated `org.junit.jupiter:junit-jupiter-api:5.10.2` to `5.14.1`

#### Test Dependency Updates

* Updated `org.junit.jupiter:junit-jupiter-params:5.10.2` to `5.14.1`

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
* Updated `org.apache.maven.plugins:maven-javadoc-plugin:3.11.1` to `3.12.0`
* Updated `org.apache.maven.plugins:maven-resources-plugin:3.3.1` to `3.4.0`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.5.2` to `3.5.4`
* Updated `org.codehaus.mojo:flatten-maven-plugin:1.6.0` to `1.7.3`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.18.0` to `2.20.1`
* Updated `org.jacoco:jacoco-maven-plugin:0.8.12` to `0.8.14`
* Updated `org.sonarsource.scanner.maven:sonar-maven-plugin:5.0.0.4389` to `5.5.0.6356`
* Added `org.sonatype.central:central-publishing-maven-plugin:0.9.0`
* Removed `org.sonatype.plugins:nexus-staging-maven-plugin:1.7.0`
