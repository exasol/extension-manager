# Extension Manager 0.5.10, released 2024-04-25

Code name: Fix reading NULL values from metadata tables

## Summary

This release fixes an error reading `NULL` values from metadata tables:

```
Could not get database installed extension list failed to read metadata tables. Cause: failed to read row of SYS.EXA_ALL_VIRTUAL_SCHEMAS: sql: Scan error on column index 4, name \"ADAPTER_NOTES\": converting NULL to string is unsupported
```

## Bugfixes

* #174: Fixed reading `NULL` values from metadata tables

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `github.com/dop251/goja_nodejs:v0.0.0-20240221231712-27eeffc9c235` to `v0.0.0-20240418154818-2aae10d4cbcf`
* Updated `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v0.3.6` to `v0.3.8`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:exasol-test-setup-abstraction-java:2.1.2` to `2.1.3`
* Updated `com.exasol:extension-manager-client-java:0.5.9` to `0.5.10`

#### Test Dependency Updates

* Updated `com.exasol:udf-debugging-java:0.6.12` to `0.6.13`
* Updated `org.slf4j:slf4j-jdk14:2.0.12` to `2.0.13`
