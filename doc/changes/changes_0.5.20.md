# Extension Manager 0.5.20, released 2026-06-11

Code name: Fix vulnerabilities in Go and NPM dependencies

## Summary

This release fixes vulnerabilities in the Go toolchain and in JavaScript dependencies used by the registry and upload tooling.

Fixed vulnerabilities:
* Go standard library, [`net/textproto`](https://pkg.go.dev/vuln/GO-2026-5039): arbitrary inputs are included in errors without any escaping.
* Go standard library, [`crypto/x509`](https://pkg.go.dev/vuln/GO-2026-5037): inefficient candidate hostname parsing.
* Go standard library, [`mime`](https://pkg.go.dev/vuln/GO-2026-5038): quadratic complexity in `WordDecoder.DecodeHeader`.
* Go module, [`golang.org/x/crypto`](https://pkg.go.dev/vuln/?q=golang.org%2Fx%2Fcrypto): SSH vulnerabilities including deadlocks on unexpected responses, memory leaks when rejecting channels, server panics during `CheckHostKey`/`Authenticate`, bypass of certificate restrictions, byte arithmetic underflow and panic, and SSH agent constraint handling issues.
* Go module, [`golang.org/x/net`](https://pkg.go.dev/vuln/?q=golang.org%2Fx%2Fnet): infinite loop in HTTP/2 transport when given a bad `SETTINGS_MAX_FRAME_SIZE`.
* npm, [`brace-expansion`](https://github.com/advisories/GHSA-f886-m6hf-6m8v) and [`brace-expansion`](https://github.com/advisories/GHSA-jxxr-4gwj-5jf2): zero-step sequences can cause process hangs and memory exhaustion; large numeric ranges can defeat documented `max` DoS protection.
* npm, [`fast-uri`](https://github.com/advisories/GHSA-q3j6-qgpj-74h6) and [`fast-uri`](https://github.com/advisories/GHSA-v39h-62p7-jpjc): percent-encoded dot segments can enable path traversal; percent-encoded authority delimiters can cause host confusion.
* npm, [`fast-xml-builder`](https://github.com/advisories/GHSA-5wm8-gmm8-39j9): attribute values with unwanted quotes can bypass malicious or unwanted attributes.
* npm, [`fast-xml-parser`](https://github.com/advisories/GHSA-gh4j-gqv2-49f6): XML comments and CDATA can be injected via unescaped delimiters.
* npm, [`follow-redirects`](https://github.com/advisories/GHSA-r4q5-vmmm-2653): custom authentication headers can leak to cross-domain redirect targets.

## Security

* #219: Fix vulnerabilities in Go and NPM dependencies

## Dependency Updates

### Extension-manager

#### Compile Dependency Updates

* Updated `golang:1.25.0` to `1.26.0`
* Updated `golang.org/x/mod:v0.34.0` to `v0.37.0`
* Updated `github.com/dop251/goja:v0.0.0-20260311135729-065cd970411c` to `v0.0.0-20260607120635-348e6bea910d`
* Updated `github.com/exasol/exasol-driver-go:v1.0.16` to `v1.0.17`
* Updated `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v1.0.0` to `v1.0.1`
* Updated `github.com/go-chi/chi/v5:v5.2.5` to `v5.3.0`

### Extension Manager Java Client

#### Compile Dependency Updates

* Updated `com.fasterxml.jackson.core:jackson-annotations:2.21` to `2.22`
* Updated `com.fasterxml.jackson.core:jackson-core:2.21.2` to `2.22.0`
* Updated `com.fasterxml.jackson.core:jackson-databind:2.21.2` to `2.22.0`
* Updated `io.swagger.core.v3:swagger-annotations:2.2.45` to `2.2.50`

#### Test Dependency Updates

* Updated `org.junit.jupiter:junit-jupiter-api:5.14.3` to `5.14.4`
* Updated `org.junit.jupiter:junit-jupiter-params:5.14.3` to `5.14.4`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:2.0.6` to `2.0.7`
* Updated `io.github.git-commit-id:git-commit-id-maven-plugin:9.0.2` to `10.0.0`
* Updated `org.apache.maven.plugins:maven-resources-plugin:3.4.0` to `3.5.0`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.5.4` to `3.5.5`

### Extension Integration Tests Library

#### Compile Dependency Updates

* Updated `com.exasol:extension-manager-client-java:0.5.19` to `0.5.20`
* Updated `com.exasol:test-db-builder-java:3.6.4` to `4.0.0`
* Updated `org.junit.jupiter:junit-jupiter-api:5.14.3` to `5.14.4`

#### Test Dependency Updates

* Updated `org.junit.jupiter:junit-jupiter-params:5.14.3` to `5.14.4`
* Updated `org.slf4j:slf4j-jdk14:2.0.17` to `2.0.18`

#### Plugin Dependency Updates

* Updated `com.exasol:error-code-crawler-maven-plugin:2.0.6` to `2.0.7`
* Updated `io.github.git-commit-id:git-commit-id-maven-plugin:9.0.2` to `10.0.0`
* Updated `org.apache.maven.plugins:maven-failsafe-plugin:3.5.4` to `3.5.5`
* Updated `org.apache.maven.plugins:maven-resources-plugin:3.4.0` to `3.5.0`
* Updated `org.apache.maven.plugins:maven-surefire-plugin:3.5.4` to `3.5.5`

### Registry

#### Compile Dependency Updates

* Updated `aws-cdk-lib:2.245.0` to `2.258.1`

#### Development Dependency Updates

* Updated `@types/node:^25.5.0` to `^25.9.2`
* Updated `ts-jest:^29.4.6` to `^29.4.11`
* Updated `typescript:~5.9.3` to `~6.0.3`
* Updated `aws-cdk:2.1115.0` to `2.1126.0`
* Updated `jest:^30.3.0` to `^30.4.2`

### Registry-upload

#### Compile Dependency Updates

* Updated `@aws-sdk/client-cloudfront:^3.1019.0` to `^3.1065.0`
* Updated `@aws-sdk/client-s3:^3.1019.0` to `^3.1065.0`
* Updated `@aws-sdk/client-cloudformation:^3.1019.0` to `^3.1065.0`
* Removed `follow-redirects:^1.15.11`

#### Development Dependency Updates

* Updated `eslint:10.1.0` to `10.4.1`
* Updated `@types/node:^25.5.0` to `^25.9.2`
* Updated `typescript-eslint:^8.57.2` to `^8.61.0`
* Updated `typescript:~5.9.3` to `~6.0.3`
* Removed `@types/follow-redirects:^1.14.4`
* Removed `ts-node:^10.9.2`

### ParameterValidator

#### Compile Dependency Updates

* Updated `@exasol/extension-parameter-validator:0.3.1` to `0.3.2`

#### Development Dependency Updates

* Updated `typescript:5.9.3` to `6.0.3`
* Updated `esbuild:0.27.4` to `0.28.0`
