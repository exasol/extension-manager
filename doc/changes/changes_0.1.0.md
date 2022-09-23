# extension-manager 0.1.0, released 2022-??-??

Code name:

## Summary

## Features

* #1: Added extension interface
* #6: Added backend
* #10: Added REST API
* #13: Added extension discovery
* #17: Added check for preconditions of extensions
* #20: Added swagger doc for REST API
* #14: Added parameters to installations response
* #31: Added parameter validator
* #30: Added optional command line flag for the path to the extension directory
* #37: Added support for more metadata tables, use a specific schema for extensions
* #41: Added step to create the `EXA_EXTENSION` schema before installing an extension
* #45: Added endpoint for installing an extension
* #54: Added database authentication via tokens
* #55: Converted JavaScript ApiError to Go APIError to use correct response status code

## Refactoring

* #2: Moved API to dedicated repo: [exasol/extension-manager-interface](https://github.com/exasol/extension-manager-interface/)
* #23: Moved go code from backend/ folder to project root
* #27: Changed extension registration to use the `global` object
* #66: Made URLs more consistent
* #65: Added extension registry

## Bugfixes

* #42: Added error handling for exceptions in JavaScript code
* #70: Return status 404 instead of 500 for unknown extension IDs

## Documentation

* #4: Added design

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Added `golang:1.19`
* Added `github.com/stretchr/testify:v1.8.0`
* Added `github.com/dop251/goja:v0.0.0-20220906144433-c4d370b87b45`
* Added `github.com/sirupsen/logrus:v1.9.0`
* Added `github.com/dop251/goja_nodejs:v0.0.0-20220905124449-678b33ca5009`
* Added `github.com/swaggo/http-swagger:v1.3.3`
* Added `github.com/exasol/exasol-driver-go:v0.4.5`
* Added `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v0.2.4`
* Added `github.com/go-chi/chi/v5:v5.0.7`

#### Test Dependency Updates

* Added `github.com/Nightapes/go-rest:v0.2.1`
* Added `github.com/kinbiko/jsonassert:v1.1.1`
* Added `github.com/DATA-DOG/go-sqlmock:v1.5.0`

### Extension Integration Tests Library

#### Test Dependency Updates

* Added `com.brsanthu:migbase64:2.2`
* Added `com.fasterxml.jackson.core:jackson-annotations:2.13.4`
* Added `com.fasterxml.jackson.core:jackson-core:2.13.4`
* Added `com.fasterxml.jackson.core:jackson-databind:2.13.4`
* Added `io.swagger.core.v3:swagger-annotations:2.2.2`
* Added `org.glassfish.jersey.core:jersey-client:2.36`
* Added `org.glassfish.jersey.inject:jersey-hk2:2.36`
* Added `org.glassfish.jersey.media:jersey-media-json-jackson:2.36`
* Added `org.glassfish.jersey.media:jersey-media-multipart:2.36`

#### Plugin Dependency Updates

* Added `com.exasol:error-code-crawler-maven-plugin:1.1.2`
* Added `com.exasol:project-keeper-maven-plugin:2.8.0`
* Added `io.github.zlika:reproducible-build-maven-plugin:0.15`
* Added `io.swagger.codegen.v3:swagger-codegen-maven-plugin:3.0.34`
* Added `org.apache.maven.plugins:maven-clean-plugin:2.5`
* Added `org.apache.maven.plugins:maven-compiler-plugin:3.10.1`
* Added `org.apache.maven.plugins:maven-deploy-plugin:3.0.0-M1`
* Added `org.apache.maven.plugins:maven-enforcer-plugin:3.1.0`
* Added `org.apache.maven.plugins:maven-gpg-plugin:3.0.1`
* Added `org.apache.maven.plugins:maven-install-plugin:2.4`
* Added `org.apache.maven.plugins:maven-jar-plugin:2.4`
* Added `org.apache.maven.plugins:maven-javadoc-plugin:3.4.0`
* Added `org.apache.maven.plugins:maven-resources-plugin:2.6`
* Added `org.apache.maven.plugins:maven-site-plugin:3.3`
* Added `org.apache.maven.plugins:maven-source-plugin:3.2.1`
* Added `org.apache.maven.plugins:maven-surefire-plugin:3.0.0-M5`
* Added `org.codehaus.mojo:build-helper-maven-plugin:3.3.0`
* Added `org.codehaus.mojo:exec-maven-plugin:3.0.0`
* Added `org.codehaus.mojo:flatten-maven-plugin:1.2.7`
* Added `org.codehaus.mojo:versions-maven-plugin:2.10.0`
* Added `org.jacoco:jacoco-maven-plugin:0.8.8`
* Added `org.sonarsource.scanner.maven:sonar-maven-plugin:3.9.1.2184`
* Added `org.sonatype.ossindex.maven:ossindex-maven-plugin:3.2.0`
* Added `org.sonatype.plugins:nexus-staging-maven-plugin:1.6.13`
