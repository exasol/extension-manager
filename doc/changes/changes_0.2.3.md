# Extension Manager 0.2.3, released 2023-04-??

Code name: Upgrade dependencies on top of 0.2.2

## Summary

This release fixes CVE-2022-41723 in `golang.org/x/net`.

## Bugfixes

* #96: Upgraded dependencies to fix vulnerability

## Dependency Updates

### Extension Manager Java Client

#### Plugin Dependency Updates

* Updated `com.exasol:project-keeper-maven-plugin:2.9.3` to `2.9.6`
* Updated `org.apache.maven.plugins:maven-deploy-plugin:3.0.0` to `3.1.0`
* Updated `org.apache.maven.plugins:maven-enforcer-plugin:3.1.0` to `3.2.1`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.2.2` to `0.2.3`

#### Test Dependency Updates

* Updated `org.mockito:mockito-junit-jupiter:5.1.1` to `5.2.0`

#### Plugin Dependency Updates

* Updated `com.exasol:project-keeper-maven-plugin:2.9.3` to `2.9.6`
* Updated `org.apache.maven.plugins:maven-deploy-plugin:3.0.0` to `3.1.0`
* Updated `org.apache.maven.plugins:maven-enforcer-plugin:3.1.0` to `3.2.1`
