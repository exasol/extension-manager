# Extension Manager 0.5.11, released 2024-05-06

Code name: Improve error message for creating duplicate instances

## Summary

This release adds a new test to the shared integration test class `AbstractVirtualSchemaExtensionIT` that verifies the error message in case the user tries to create a virtual schema that already exists. It also allows overriding the parameter name of the virtual schema. This is required in case an extensions is not based on the base virtual schema extension.

## Bugfix

* #177: Improve error message for creating duplicate instances

## Dependency Updates

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `com.fasterxml.jackson.core:jackson-annotations:2.17.0` to `2.17.1`
* Updated `com.fasterxml.jackson.core:jackson-core:2.17.0` to `2.17.1`
* Updated `com.fasterxml.jackson.core:jackson-databind:2.17.0` to `2.17.1`

#### Test Dependency Updates

* Updated `org.junit.jupiter:junit-jupiter-api:5.10.1` to `5.10.2`
* Updated `org.junit.jupiter:junit-jupiter-params:5.10.1` to `5.10.2`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.5.10` to `0.5.11`
* Updated `org.junit.jupiter:junit-jupiter-api:5.10.1` to `5.10.2`

#### Test Dependency Updates

* Updated `org.junit.jupiter:junit-jupiter-params:5.10.1` to `5.10.2`
