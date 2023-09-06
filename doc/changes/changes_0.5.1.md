# Extension Manager 0.5.1, released 2023-??-??

Code name: Support `select` parameter type

## Summary

This release contains the following changes:
* It adds support for the `select` parameter type.
* In Integration test class `PreviousExtensionVersion` the builder property `adapterFileName` is now optional. This is useful for extensions that don't need an adapter file like Lua based virtual schemas.
* Command line flag `-addCauseToInternalServerError` of the standalone server now allows adding the root cause error message to internal server errors (status 500).

    This is helpful during debugging because the error message contains the actual root cause and you don't need to check the log for errors. The Java integration test framework `extension-manager-integration-test-java` enables this flag automatically.

    ⚠️Warning⚠️: Do not enable this flag in production environments as this might leak internal information.

## Features

* #132: Added support for the `select` parameter type
* #134: Made adapter file optional for `PreviousExtensionVersion`
* #131: Added command line flag `-addCauseToInternalServerError`

## Documentation

* #129: Improved description of deployment process

## Dependency Updates

### Extension Manager Java Client

#### Test Dependency Updates

* Updated `nl.jqno.equalsverifier:equalsverifier:3.15` to `3.15.1`

#### Plugin Dependency Updates

* Updated `com.exasol:project-keeper-maven-plugin:2.9.10` to `2.9.11`
* Updated `org.apache.maven.plugins:maven-enforcer-plugin:3.3.0` to `3.4.0`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:exasol-test-setup-abstraction-java:2.0.2` to `2.0.3`
* Updated `com.exasol:extension-manager-client-java:0.5.0` to `0.5.1`
* Updated `com.exasol:test-db-builder-java:3.4.2` to `3.5.0`

#### Test Dependency Updates

* Updated `org.mockito:mockito-junit-jupiter:5.4.0` to `5.5.0`

#### Plugin Dependency Updates

* Updated `com.exasol:project-keeper-maven-plugin:2.9.10` to `2.9.11`
* Updated `org.apache.maven.plugins:maven-enforcer-plugin:3.3.0` to `3.4.0`

### ParameterValidator

#### Compile Dependency Updates

* Updated `@exasol/extension-parameter-validator:0.2.1` to `0.3.0`

#### Development Dependency Updates

* Updated `typescript:5.1.6` to `5.2.2`
* Updated `esbuild:0.18.16` to `0.19.2`
