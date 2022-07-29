# Exasol Extension Manager

[![Build Status](https://github.com/exasol/extension-manager/actions/workflows/ci-build.yml/badge.svg)](https://github.com/exasol/extension-manager/actions/workflows/ci-build.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/exasol/extension-manager.svg)](https://pkg.go.dev/github.com/exasol/extension-manager)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=com.exasol%3Aextension-manager&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=com.exasol%3Aextension-manager)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=com.exasol%3Aextension-manager&metric=coverage)](https://sonarcloud.io/summary/new_code?id=com.exasol%3Aextension-manager)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=com.exasol%3Aextension-manager&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=com.exasol%3Aextension-manager)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=com.exasol%3Aextension-manager&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=com.exasol%3Aextension-manager)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=com.exasol%3Aextension-manager&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=com.exasol%3Aextension-manager)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=com.exasol%3Aextension-manager&metric=bugs)](https://sonarcloud.io/summary/new_code?id=com.exasol%3Aextension-manager)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=com.exasol%3Aextension-manager&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=com.exasol%3Aextension-manager)

This project contains the Exasol extension manager. It's a tool for installing and managing extensions like Virtual Schemas.

## REST API Documentation

The extension-manager exposes a REST API for the frontend.
<!-- markdown-link-check-disable-next-line -->
This API is documented via Swagger. In order to view it, checkout this repo, run `go run ./...` and open `http://localhost:8080/swagger/index.html`.

You can find the API definition here:
* [swagger.json](./generatedApiDocs/swagger.json)
* [swagger.yaml](./generatedApiDocs/swagger.yaml)

## Additional Information

* [Changelog](doc/changes/changelog.md)
* [Software Design (online, main branch)](https://exasol.github.io/extension-manager/design.html)
* [Software Design (local)](doc/design.md)
* [Developers Guide](doc/developers_guide.md)
* [Dependencies](dependencies.md)

## Related Projects

* [extension-manager](https://github.com/exasol/extension-manager): Extension manager backend written in Go (this repo)
* [extension-manager-interface](https://github.com/exasol/extension-manager-interface/): Extension interface defined in TypeScript, published to npm as [@exasol/extension-manager-interface](https://www.npmjs.com/package/@exasol/extension-manager-interface)
* [extension-parameter-validator](https://github.com/exasol/extension-parameter-validator): Validator for extension parameters written in TypeScript, published to npm as [@exasol/extension-parameter-validator](https://www.npmjs.com/package/@exasol/extension-parameter-validator)
* Virtual Schemas providing extensions
  * [s3-document-files-virtual-schema](https://github.com/exasol/s3-document-files-virtual-schema/): Work in progress, see [PR#84](https://github.com/exasol/s3-document-files-virtual-schema/pull/84)
