# Extension Manager 0.5.6, released 2023-??-??

Code name: Adapt for JDBC based virtual schema extensions

## Summary

This release adapts Extension Manager integration test framework for JDBC based extensions:
* Skip `upgradeFromPreviousVersion()` when no previous version is available
* Update `createInstanceWithSingleQuote()` to accept any adapter name, not just the hard-coded S3 adapter

## Features

* #161: Adapted shared tests for JDBC based extensions

## Dependency Updates

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.5.5` to `0.5.6`
