# Exasol Extension Manager

[![Build Status](https://github.com/exasol/extension-manager/actions/workflows/ci-build.yml/badge.svg)](https://github.com/exasol/extension-manager/actions/workflows/ci-build.yml)

This project contains the Exasol extension manager. It's a tool for installing and managing extensions like Virtual
Schemas.

## REST API doc

The extension-manager exposes a REST API for the frontend.
This API is documented via Swagger. In order to view it, checkout this repo, run `go run ./main/` and open http://localhost:8080/swagger/index.html.


## Additional Information

* [Changelog](doc/changes/changelog.md)
* [Software Design (online, main branch)](https://exasol.github.io/extension-manager/design.html)
* [Software Design (local)](doc/design/design.md)
* [Developers Guide](doc/developers_guide/developers_guide.md)
* [Dependencies](dependencies.md)
