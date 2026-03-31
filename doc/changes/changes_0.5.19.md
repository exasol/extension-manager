# Extension Manager 0.5.19, released 2026-03-31

Code name: Update dependencies on top of 0.5.18

## Summary

This release fixes the following vulnerabilities in dependencies:
* [#118: Pygments has Regular Expression Denial of Service (ReDoS) due to Inefficient Regex for GUID Matching](https://github.com/exasol/extension-manager/security/dependabot/118)
* [#117: Handlebars.js has a Property Access Validation Bypass in container.lookup Low Development](https://github.com/exasol/extension-manager/security/dependabot/117)
* [#116: Handlebars.js has a Prototype Method Access Control Gap via Missing __lookupSetter__ Blocklist Entry](https://github.com/exasol/extension-manager/security/dependabot/116)
* [#115: brace-expansion: Zero-step sequence causes process hang and memory exhaustion](https://github.com/exasol/extension-manager/security/dependabot/115)
* [#114: brace-expansion: Zero-step sequence causes process hang and memory exhaustion](https://github.com/exasol/extension-manager/security/dependabot/114)
* [#113: Handlebars.js has JavaScript Injection in CLI Precompiler via Unescaped Names and Options](https://github.com/exasol/extension-manager/security/dependabot/113)
* [#112: Handlebars.js has JavaScript Injection via AST Type Confusion when passing an object as dynamic partial](https://github.com/exasol/extension-manager/security/dependabot/112)
* [#111: Handlebars.js has Denial of Service via Malformed Decorator Syntax in Template Compilation](https://github.com/exasol/extension-manager/security/dependabot/111)
* [#110: Handlebars.js has JavaScript Injection via AST Type Confusion by tampering @partial-block](https://github.com/exasol/extension-manager/security/dependabot/110)
* [#109: Handlebars.js has JavaScript Injection via AST Type Confusion](https://github.com/exasol/extension-manager/security/dependabot/109)
* [#88: jackson-core: Number Length Constraint Bypass in Async Parser Leads to Potential DoS Condition](https://github.com/exasol/extension-manager/security/dependabot/88)

## Security

* #211: Fix vulnerabilities in dependencies reported by Dependabot

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `golang:1.24.9` to `1.25.0`
* Updated `golang.org/x/mod:v0.31.0` to `v0.34.0`
* Updated `github.com/dop251/goja:v0.0.0-20251201205617-2bb4c724c0f9` to `v0.0.0-20260311135729-065cd970411c`
* Updated `github.com/sirupsen/logrus:v1.9.3` to `v1.9.4`
* Updated `github.com/dop251/goja_nodejs:v0.0.0-20251015164255-5e94316bedaf` to `v0.0.0-20260212111938-1f56ff5bcf14`
* Updated `github.com/exasol/exasol-driver-go:v1.0.15` to `v1.0.16`
* Updated `github.com/go-chi/chi/v5:v5.2.3` to `v5.2.5`

#### Other Dependency Updates

* Removed `toolchain:go1.25.0`

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `com.fasterxml.jackson.core:jackson-annotations:2.20` to `2.21`
* Updated `com.fasterxml.jackson.core:jackson-core:2.20.1` to `2.21.2`
* Updated `com.fasterxml.jackson.core:jackson-databind:2.20.1` to `2.21.2`
* Updated `io.swagger.core.v3:swagger-annotations:2.2.41` to `2.2.45`

#### Test Dependency Updates

* Updated `org.junit.jupiter:junit-jupiter-api:5.14.1` to `5.14.3`
* Updated `org.junit.jupiter:junit-jupiter-params:5.14.1` to `5.14.3`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:2.0.5` to `2.0.6`
* Updated `io.swagger.codegen.v3:swagger-codegen-maven-plugin:3.0.75` to `3.0.78`
* Updated `org.apache.maven.plugins:maven-compiler-plugin:3.14.1` to `3.15.0`
* Updated `org.apache.maven.plugins:maven-source-plugin:3.2.1` to `3.4.0`
* Updated `org.codehaus.mojo:exec-maven-plugin:3.6.2` to `3.6.3`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.20.1` to `2.21.0`
* Updated `org.sonatype.central:central-publishing-maven-plugin:0.9.0` to `0.10.0`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:exasol-test-setup-abstraction-java:2.1.10` to `2.1.11`
* Updated `com.exasol:extension-manager-client-java:0.5.18` to `0.5.19`
* Updated `org.junit.jupiter:junit-jupiter-api:5.14.1` to `5.14.3`

#### Test Dependency Updates

* Updated `org.junit.jupiter:junit-jupiter-params:5.14.1` to `5.14.3`
* Updated `org.mockito:mockito-junit-jupiter:5.21.0` to `5.23.0`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:2.0.5` to `2.0.6`
* Updated `org.apache.maven.plugins:maven-compiler-plugin:3.14.1` to `3.15.0`
* Updated `org.apache.maven.plugins:maven-source-plugin:3.2.1` to `3.4.0`
* Updated `org.codehaus.mojo:versions-maven-plugin:2.20.1` to `2.21.0`
* Updated `org.sonatype.central:central-publishing-maven-plugin:0.9.0` to `0.10.0`

### Registry

#### Compile Dependency Updates

* Updated `aws-cdk-lib:2.233.0` to `2.245.0`
* Updated `constructs:^10.4.4` to `^10.6.0`

#### Development Dependency Updates

* Updated `@types/node:^25.0.3` to `^25.5.0`
* Updated `aws-cdk:2.1100.1` to `2.1115.0`
* Updated `jest:^30.2.0` to `^30.3.0`
* Removed `@types/prettier:3.0.0`

### Registry-upload

#### Compile Dependency Updates

* Updated `@aws-sdk/client-cloudfront:^3.955.0` to `^3.1019.0`
* Updated `@aws-sdk/client-s3:^3.955.0` to `^3.1019.0`
* Updated `@aws-sdk/client-cloudformation:^3.955.0` to `^3.1019.0`

#### Development Dependency Updates

* Updated `eslint:9.39.2` to `10.1.0`
* Added `@eslint/js:10.0.1`
* Updated `@types/node:^25.0.3` to `^25.5.0`
* Updated `typescript-eslint:^8.50.0` to `^8.57.2`
* Removed `@types/prettier:3.0.0`

### ParameterValidator

#### Development Dependency Updates

* Updated `esbuild:0.27.2` to `0.27.4`
