# Extension Manager 0.5.8, released 2024-03-12

Code name: Fix vulnerabilities CVE-2024-25710 and CVE-2024-26308 in compile dependency `org.apache.commons:commons-compress` of the integration test framework

## Summary

This release fixed vulnerabilities CVE-2024-25710 and CVE-2024-26308 in compile dependency `org.apache.commons:commons-compress` of the integration test framework.

## Security

* #165: Fixed vulnerabilities CVE-2024-25710 and CVE-2024-26308 in compile dependency `org.apache.commons:commons-compress`

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `github.com/stretchr/testify:v1.8.4` to `v1.9.0`
* Updated `github.com/dop251/goja:v0.0.0-20231027120936-b396bb4c349d` to `v0.0.0-20240220182346-e401ed450204`
* Updated `github.com/dop251/goja_nodejs:v0.0.0-20231022114343-5c1f9037c9ab` to `v0.0.0-20240221231712-27eeffc9c235`
* Updated `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v0.3.4` to `v0.3.5`
* Updated `github.com/go-chi/chi/v5:v5.0.10` to `v5.0.12`

#### Test Dependency Updates

* Updated `golang.org/x/mod:v0.14.0` to `v0.16.0`
* Updated `github.com/DATA-DOG/go-sqlmock:v1.5.0` to `v1.5.2`

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `com.fasterxml.jackson.core:jackson-annotations:2.16.0` to `2.16.2`
* Updated `com.fasterxml.jackson.core:jackson-core:2.16.0` to `2.16.2`
* Updated `com.fasterxml.jackson.core:jackson-databind:2.16.0` to `2.16.2`
* Updated `io.swagger.core.v3:swagger-annotations:2.2.19` to `2.2.20`

#### Test Dependency Updates

* Updated `nl.jqno.equalsverifier:equalsverifier:3.15.3` to `3.15.8`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:1.3.1` to `2.0.0`
* Removed `com.exasol:project-keeper-maven-plugin:2.9.16`
* Updated `io.swagger.codegen.v3:swagger-codegen-maven-plugin:3.0.47` to `3.0.54`
* Updated `org.apache.maven.plugins:maven-compiler-plugin:3.11.0` to `3.12.1`
* Updated `org.apache.maven.plugins:maven-javadoc-plugin:3.6.2` to `3.6.3`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.2.2` to `3.2.5`
* Added `org.apache.maven.plugins:maven-toolchains-plugin:3.1.0`
* Updated `org.codehaus.mojo:build-helper-maven-plugin:3.4.0` to `3.5.0`
* Updated `org.codehaus.mojo:exec-maven-plugin:3.1.0` to `3.2.0`
* Updated `org.codehaus.mojo:flatten-maven-plugin:1.5.0` to `1.6.0`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.16.1` to `2.16.2`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:exasol-test-setup-abstraction-java:2.1.0` to `2.1.1`
* Updated `com.exasol:extension-manager-client-java:0.5.7` to `0.5.8`
* Updated `com.exasol:hamcrest-resultset-matcher:1.6.2` to `1.6.5`
* Updated `com.exasol:test-db-builder-java:3.5.2` to `3.5.4`

#### Test Dependency Updates

* Updated `com.exasol:udf-debugging-java:0.6.11` to `0.6.12`
* Updated `org.mockito:mockito-junit-jupiter:5.7.0` to `5.11.0`
* Updated `org.slf4j:slf4j-jdk14:2.0.9` to `2.0.12`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:1.3.1` to `2.0.0`
* Removed `com.exasol:project-keeper-maven-plugin:2.9.16`
* Updated `org.apache.maven.plugins:maven-compiler-plugin:3.11.0` to `3.12.1`
* Updated `org.apache.maven.plugins:maven-failsafe-plugin:3.2.2` to `3.2.5`
* Updated `org.apache.maven.plugins:maven-javadoc-plugin:3.6.2` to `3.6.3`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.2.2` to `3.2.5`
* Added `org.apache.maven.plugins:maven-toolchains-plugin:3.1.0`
* Updated `org.codehaus.mojo:flatten-maven-plugin:1.5.0` to `1.6.0`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.16.1` to `2.16.2`

### Registry

#### Compile Dependency Updates

* Updated `aws-cdk-lib:2.110.0` to `2.132.1`

#### Development Dependency Updates

* Updated `@types/node:^20.9.1` to `^20.11.26`
* Updated `ts-jest:^29.1.1` to `^29.1.2`
* Updated `@types/jest:^29.5.8` to `^29.5.12`
* Updated `typescript:~5.2.2` to `~5.4.2`
* Updated `aws-cdk:2.110.0` to `2.132.1`
* Updated `ts-node:^10.9.1` to `^10.9.2`

### ParameterValidator

#### Development Dependency Updates

* Updated `typescript:5.2.2` to `5.4.2`
* Updated `esbuild:0.19.5` to `0.20.1`
