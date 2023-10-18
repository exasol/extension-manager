# Extension Manager 0.5.4, released 2023-10-18

Code name: Migrate to Exasol 8

## Summary

Starting with this release, Extension Manager requires Exasol 8. Also integration tests for extension definitions must run with Exasol 8.

To skip tests on non v8 versions, you can use the following new method:

```java
import com.exasol.extensionmanager.itest;
// ...
ExasolVersionCheck.assumeExasolVersion8(exasolTestSetup)
```

## Features

* #152: Migrated tests to Exasol 8

## Dependency Updates

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.5.3` to `0.5.4`
