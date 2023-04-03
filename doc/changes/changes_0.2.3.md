# Extension Manager 0.2.3, released 2023-04-??

Code name: Upgrade dependencies on top of 0.2.2

## Summary

This release fixes CVE-2022-41723 in `golang.org/x/net`.

## Bugfixes

* #96: Upgraded dependencies to fix vulnerability

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `github.com/stretchr/testify:v1.8.1` to `v1.8.2`
* Updated `github.com/dop251/goja:v0.0.0-20230203172422-5460598cfa32` to `v0.0.0-20230402114112-623f9dda9079`
* Updated `github.com/dop251/goja_nodejs:v0.0.0-20230207183254-2229640ea097` to `v0.0.0-20230322100729-2550c7b6c124`
* Updated `github.com/swaggo/http-swagger:v1.3.3` to `v1.3.4`
* Updated `github.com/exasol/exasol-driver-go:v0.4.6` to `v0.4.7`

#### Test Dependency Updates

* Updated `github.com/Nightapes/go-rest:v0.3.1` to `v0.3.2`

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `io.swagger.core.v3:swagger-annotations:2.2.8` to `2.2.9`

#### Test Dependency Updates

* Updated `nl.jqno.equalsverifier:equalsverifier:3.13.1` to `3.14.1`

#### Plugin Dependency Updates

* Updated `com.exasol:project-keeper-maven-plugin:2.9.3` to `2.9.6`
* Updated `io.swagger.codegen.v3:swagger-codegen-maven-plugin:3.0.39` to `3.0.41`
* Updated `org.apache.maven.plugins:maven-deploy-plugin:3.0.0` to `3.1.0`
* Updated `org.apache.maven.plugins:maven-enforcer-plugin:3.1.0` to `3.2.1`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.2.2` to `0.2.3`
* Updated `com.exasol:hamcrest-resultset-matcher:1.5.2` to `1.5.3`

#### Test Dependency Updates

* Updated `org.mockito:mockito-junit-jupiter:5.1.1` to `5.2.0`

#### Plugin Dependency Updates

* Updated `com.exasol:project-keeper-maven-plugin:2.9.3` to `2.9.6`
* Updated `org.apache.maven.plugins:maven-deploy-plugin:3.0.0` to `3.1.0`
* Updated `org.apache.maven.plugins:maven-enforcer-plugin:3.1.0` to `3.2.1`
