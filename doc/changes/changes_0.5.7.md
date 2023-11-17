# Extension Manager 0.5.7, released 2023-11-27

Code name: Fix base integration test

## Summary

This release fixes a hard coded project name in `AbstractScriptExtensionIT`, which made the test base unusable for other extensions.
## Bugfix

* #163: Fixed hard coded project name in `AbstractScriptExtensionIT`

## Dependency Updates

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `com.fasterxml.jackson.core:jackson-annotations:2.15.3` to `2.16.0`
* Updated `com.fasterxml.jackson.core:jackson-core:2.15.3` to `2.16.0`
* Updated `com.fasterxml.jackson.core:jackson-databind:2.15.3` to `2.16.0`
* Updated `io.swagger.core.v3:swagger-annotations:2.2.18` to `2.2.19`

#### Plugin Dependency Updates

* Updated `com.exasol:project-keeper-maven-plugin:2.9.15` to `2.9.16`
* Updated `org.apache.maven.plugins:maven-javadoc-plugin:3.6.0` to `3.6.2`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.1.2` to `3.2.2`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:exasol-test-setup-abstraction-java:2.0.4` to `2.1.0`
* Updated `com.exasol:extension-manager-client-java:0.5.6` to `0.5.7`
* Updated `com.exasol:test-db-builder-java:3.5.1` to `3.5.2`

#### Plugin Dependency Updates

* Updated `com.exasol:project-keeper-maven-plugin:2.9.15` to `2.9.16`
* Updated `org.apache.maven.plugins:maven-failsafe-plugin:3.1.2` to `3.2.2`
* Updated `org.apache.maven.plugins:maven-javadoc-plugin:3.6.0` to `3.6.2`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.1.2` to `3.2.2`

### Registry

#### Compile Dependency Updates

* Updated `aws-cdk-lib:2.104.0` to `2.110.0`

#### Development Dependency Updates

* Updated `@types/node:^20.8.10` to `^20.9.1`
* Updated `@types/jest:^29.5.7` to `^29.5.8`
* Updated `aws-cdk:2.104.0` to `2.110.0`

### Registry-upload

#### Compile Dependency Updates

* Updated `@aws-sdk/client-cloudfront:^3.441.0` to `^3.451.0`
* Updated `@aws-sdk/client-s3:^3.441.0` to `^3.451.0`
* Updated `@aws-sdk/client-cloudformation:^3.441.0` to `^3.451.0`

#### Development Dependency Updates

* Updated `eslint:^8.52.0` to `^8.53.0`
* Updated `@types/follow-redirects:^1.14.3` to `^1.14.4`
* Updated `@typescript-eslint/parser:^6.9.1` to `^6.11.0`
* Updated `@types/node:^20.8.10` to `^20.9.1`
* Updated `@typescript-eslint/eslint-plugin:^6.9.1` to `^6.11.0`
