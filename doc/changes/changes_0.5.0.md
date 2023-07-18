# Extension Manager 0.5.0, released 2023-??-??

Code name:

## Summary

This release improves error handling when using extensions not implementing all functions required by EM. EM now returns a helpful error message instead of failing with a `nil`-pointer error.

A common scenario for an extension not implementing a required function is when the extension had been built using an older version of EM's extension interface.

## Bugfixes

* #105: Ensured that EM can load and use compatible extensions
* #74: Fixed most important linter warnings

## Dependency Updates

### Extension Manager Java Client

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:1.2.3` to `1.3.0`
* Updated `com.exasol:project-keeper-maven-plugin:2.9.7` to `2.9.9`
* Updated `org.apache.maven.plugins:maven-gpg-plugin:3.0.1` to `3.1.0`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.0.0` to `3.1.2`
* Updated `org.basepom.maven:duplicate-finder-maven-plugin:1.5.1` to `2.0.1`
* Updated `org.codehaus.mojo:flatten-maven-plugin:1.4.1` to `1.5.0`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.15.0` to `2.16.0`
* Updated `org.jacoco:jacoco-maven-plugin:0.8.9` to `0.8.10`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.4.0` to `0.5.0`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:1.2.3` to `1.3.0`
* Updated `com.exasol:project-keeper-maven-plugin:2.9.7` to `2.9.9`
* Updated `org.apache.maven.plugins:maven-failsafe-plugin:3.0.0` to `3.1.2`
* Updated `org.apache.maven.plugins:maven-gpg-plugin:3.0.1` to `3.1.0`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.0.0` to `3.1.2`
* Updated `org.basepom.maven:duplicate-finder-maven-plugin:1.5.1` to `2.0.1`
* Updated `org.codehaus.mojo:flatten-maven-plugin:1.4.1` to `1.5.0`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.15.0` to `2.16.0`
* Updated `org.jacoco:jacoco-maven-plugin:0.8.9` to `0.8.10`

### Registry

#### Compile Dependency Updates

* Updated `aws-cdk-lib:2.86.0` to `2.87.0`
* Updated `constructs:^10.2.64` to `^10.2.69`

#### Development Dependency Updates

* Updated `@types/node:20.3.2` to `20.4.2`
* Updated `ts-jest:^29.1.0` to `^29.1.1`
* Updated `@types/jest:^29.5.2` to `^29.5.3`
* Updated `jest:^29.5.0` to `^29.6.1`
* Updated `aws-cdk:2.86.0` to `2.87.0`

### ParameterValidator

#### Compile Dependency Updates

* Updated `@exasol/extension-parameter-validator:0.2.0` to `0.2.1`

#### Development Dependency Updates

* Updated `typescript:5.0.3` to `5.1.6`
* Updated `esbuild:0.17.15` to `0.18.13`
