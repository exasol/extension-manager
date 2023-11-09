# Extension Manager 0.5.6, released 2023-11-09

Code name: Adapt for JDBC based virtual schema extensions

## Summary

This release adapts Extension Manager integration test framework for JDBC based extensions:
* Skip `upgradeFromPreviousVersion()` when no previous version is available
* Update `createInstanceWithSingleQuote()` to accept any adapter name, not just the hard-coded S3 adapter

The release also allows extensions to ignore the file size of required BucketFS files by setting a negative `fileSize`. This is useful for JDBC drivers where we don't know the version and file size beforehand and want to allow arbitrary versions and file sizes.

## Features

* #161: Adapted shared tests for JDBC based extensions

## Dependency Updates

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.5.5` to `0.5.6`
