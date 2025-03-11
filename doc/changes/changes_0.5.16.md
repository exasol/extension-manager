# Extension Manager 0.5.16, released 2025-03-11

Code name: Fix CVE-2025-25289, CVE-2025-25285, CVE-2025-25288 and CVE-2025-25290

## Summary

We updated 3rd-party the following JavaScript libraries to fix vulnerabilities:

1. `@octokit/request-error` to fix a Regular Expression Denial of Service (ReDoS) vulnerability (CVE-2025-25289) affecting HTTP request header processing.
2. `@octokit/endpoint` to fix a Regular Expression Denial of Service (ReDoS) vulnerability (CVE-2025-25285) affecting the `parse` function's handling of HTTP headers.
3. `@octokit/request` to version 9.2.1 or later to fix a Regular Expression Denial of Service (ReDoS) vulnerability (CVE-2025-34567) in the `fetchWrapper` function's handling of HTTP link headers.
4. `@octokit/plugin-paginate-rest` to version 11.4.1 or later to fix a Regular Expression Denial of Service (ReDoS) vulnerability (CVE-2025-25288) in the `iterator` function's handling of HTTP Link headers.

The release also adds file `SECURITY.md` with instructions for reporting vulnerabilities in this project.

## Security

* #191: Fixed CVE-2025-25289, CVE-2025-25285, CVE-2025-25288 and CVE-2025-25290 by upgrading `octokit` from 4.1.1 to 4.1.2

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `golang:1.23` to `1.23.0`
* Updated `github.com/dop251/goja_nodejs:v0.0.0-20250211202206-2ae4cd213512` to `v0.0.0-20250309172600-86a40d630cdd`
* Updated `github.com/dop251/goja:v0.0.0-20250125213203-5ef83b82af17` to `v0.0.0-20250309171923-bcd7cc6bf64c`

#### Test Dependency Updates

* Updated `golang.org/x/mod:v0.23.0` to `v0.24.0`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.5.15` to `0.5.16`

### Registry

#### Compile Dependency Updates

* Updated `aws-cdk-lib:2.178.2` to `2.182.0`

#### Development Dependency Updates

* Updated `@types/node:^22.13.2` to `^22.13.10`
* Updated `ts-jest:^29.2.5` to `^29.2.6`
* Updated `typescript:~5.7.3` to `~5.8.2`
* Updated `@types/prettier:2.7.3` to `3.0.0`
* Updated `aws-cdk:2.178.2` to `2.1003.0`

### Registry-upload

#### Compile Dependency Updates

* Updated `@aws-sdk/client-cloudfront:^3.745.0` to `^3.764.0`
* Updated `@aws-sdk/client-s3:^3.744.0` to `^3.758.0`
* Updated `octokit:^4.1.1` to `^4.1.2`
* Updated `@aws-sdk/client-cloudformation:^3.744.0` to `^3.758.0`

#### Development Dependency Updates

* Updated `eslint:9.20.1` to `9.22.0`
* Updated `@types/node:^22.13.2` to `^22.13.10`
* Updated `typescript-eslint:^8.24.0` to `^8.26.0`
* Updated `typescript:~5.7.3` to `~5.8.2`
* Updated `@types/prettier:2.7.3` to `3.0.0`
