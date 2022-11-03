# Extension Manager 0.2.0, released 2022-11-02

Code name: Add extension registry

## Summary

This release adds a CDK stack for deploying the infrastructure of the Extension Registry to AWS.

We also moved all Go sources to the `pkg` directory. Projects that use this library will need to adapt the imports by replacing `"github.com/exasol/extension-manager/*"` with `"github.com/exasol/extension-manager/pkg/*"`.

## Features

* #80: Added prefix to log messages from JS extensions
* #82: Added infrastructure for extension registry

## Refactoring

* #86: Moved Go sources to `pkg` directory

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `github.com/dop251/goja:v0.0.0-20220906144433-c4d370b87b45` to `v0.0.0-20221019153710-09250e0eba20`
* Updated `github.com/dop251/goja_nodejs:v0.0.0-20220905124449-678b33ca5009` to `v0.0.0-20221009164102-3aa5028e57f6`
* Updated `github.com/exasol/exasol-driver-go:v0.4.5` to `v0.4.6`

### Extension Manager Java Client

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:1.1.2` to `1.2.1`
* Updated `com.exasol:project-keeper-maven-plugin:2.8.0` to `2.9.0`
* Updated `io.github.zlika:reproducible-build-maven-plugin:0.15` to `0.16`
* Updated `org.apache.maven.plugins:maven-deploy-plugin:3.0.0-M1` to `3.0.0`
* Updated `org.apache.maven.plugins:maven-javadoc-plugin:3.4.0` to `3.4.1`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.0.0-M5` to `3.0.0-M7`
* Updated `org.codehaus.mojo:flatten-maven-plugin:1.2.7` to `1.3.0`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.10.0` to `2.13.0`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.1.0` to `0.2.0`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:1.1.2` to `1.2.1`
* Updated `com.exasol:project-keeper-maven-plugin:2.8.0` to `2.9.0`
* Updated `io.github.zlika:reproducible-build-maven-plugin:0.15` to `0.16`
* Updated `org.apache.maven.plugins:maven-deploy-plugin:3.0.0-M1` to `3.0.0`
* Updated `org.apache.maven.plugins:maven-failsafe-plugin:3.0.0-M5` to `3.0.0-M7`
* Updated `org.apache.maven.plugins:maven-javadoc-plugin:3.4.0` to `3.4.1`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.0.0-M5` to `3.0.0-M7`
* Updated `org.codehaus.mojo:flatten-maven-plugin:1.2.7` to `1.3.0`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.10.0` to `2.13.0`
