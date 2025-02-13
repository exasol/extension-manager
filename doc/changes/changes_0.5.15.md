# Extension Manager 0.5.15, released 2025-02-13

Code name: Fix CWE-346, CWE-346, CVE-2025-23206, CVE-2025-25193, CVE-2025-24970, CVE-2025-25193 and CVE-2025-24970

## Summary

This release fixes the following vulnerabilities: CWE-346, CWE-346, CVE-2025-23206, CVE-2025-25193, CVE-2025-24970, CVE-2025-25193 and CVE-2025-24970.

**Note:** This release upgrades to Go 1.23.

## Security

* #189: Fixed vulnerabilities by upgrading dependencies

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `golang:1.22.0` to `1.23`
* Updated `github.com/dop251/goja_nodejs:v0.0.0-20240728170619-29b559befffc` to `v0.0.0-20250211202206-2ae4cd213512`
* Updated `github.com/exasol/exasol-driver-go:v1.0.10` to `v1.0.12`
* Updated `github.com/go-chi/chi/v5:v5.2.0` to `v5.2.1`
* Updated `github.com/dop251/goja:v0.0.0-20241024094426-79f3a7efcdbd` to `v0.0.0-20250125213203-5ef83b82af17`
* Updated `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v0.3.10` to `v0.3.11`

#### Test Dependency Updates

* Updated `golang.org/x/mod:v0.22.0` to `v0.23.0`

#### Other Dependency Updates

* Added `toolchain:go1.23.6`

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `io.swagger.core.v3:swagger-annotations:2.2.27` to `2.2.28`

#### Test Dependency Updates

* Updated `nl.jqno.equalsverifier:equalsverifier:3.18.1` to `3.19`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:exasol-test-setup-abstraction-java:2.1.6` to `2.1.7`
* Updated `com.exasol:extension-manager-client-java:0.5.14` to `0.5.15`

#### Test Dependency Updates

* Updated `com.exasol:udf-debugging-java:0.6.14` to `0.6.15`

### Registry

#### Compile Dependency Updates

* Updated `aws-cdk-lib:2.175.1` to `2.178.2`

#### Development Dependency Updates

* Updated `@types/node:^22.10.6` to `^22.13.2`
* Updated `aws-cdk:2.175.1` to `2.178.2`

### Registry-upload

#### Compile Dependency Updates

* Updated `@aws-sdk/client-cloudfront:^3.726.1` to `^3.745.0`
* Updated `@aws-sdk/client-s3:^3.726.1` to `^3.744.0`
* Updated `octokit:^4.1.0` to `^4.1.1`
* Updated `@aws-sdk/client-cloudformation:^3.726.1` to `^3.744.0`

#### Development Dependency Updates

* Updated `eslint:9.18.0` to `9.20.1`
* Updated `@types/node:^22.10.6` to `^22.13.2`
* Updated `typescript-eslint:^8.20.0` to `^8.24.0`

### ParameterValidator

#### Compile Dependency Updates

* Updated `@exasol/extension-parameter-validator:0.3.0` to `0.3.1`

#### Development Dependency Updates

* Updated `esbuild:0.24.2` to `0.25.0`
