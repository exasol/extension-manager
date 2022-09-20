# extension-manager 0.1.0, released 2022-??-??

Code name:

## Summary

## Features

* #1: Added extension interface
* #6: Added backend
* #10: Added REST API
* #13: Added extension discovery
* #17: Added check for preconditions of extensions
* #20: Added swagger doc for REST API
* #14: Added parameters to installations response
* #31: Added parameter validator
* #30: Added optional command line flag for the path to the extension directory
* #37: Added support for more metadata tables, use a specific schema for extensions
* #41: Added step to create the `EXA_EXTENSION` schema before installing an extension
* #45: Added endpoint for installing an extension
* #54: Added database authentication via tokens
* #55: Converted JavaScript ApiError to Go APIError to use correct response status code

## Refactoring

* #2: Moved API to dedicated repo: [exasol/extension-manager-interface](https://github.com/exasol/extension-manager-interface/)
* #23: Moved go code from backend/ folder to project root
* #27: Changed extension registration to use the `global` object
* #66: Make URLs more consistent

## Bugfixes

* #42: Added error handling for exceptions in JavaScript code
* #70: Return status 404 instead of 500 for unknown extension IDs

## Documentation

* #4: Added design

## Dependency Updates

### Compile Dependency Updates

* Added `golang:1.19`
* Added `github.com/stretchr/testify:v1.8.0`
* Added `github.com/dop251/goja:v0.0.0-20220906144433-c4d370b87b45`
* Added `github.com/sirupsen/logrus:v1.9.0`
* Added `github.com/dop251/goja_nodejs:v0.0.0-20220905124449-678b33ca5009`
* Added `github.com/swaggo/http-swagger:v1.3.3`
* Added `github.com/exasol/exasol-driver-go:v0.4.5`
* Added `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v0.2.4`
* Added `github.com/go-chi/chi/v5:v5.0.7`

### Test Dependency Updates

* Added `github.com/Nightapes/go-rest:v0.2.1`
* Added `github.com/kinbiko/jsonassert:v1.1.1`
* Added `github.com/DATA-DOG/go-sqlmock:v1.5.0`
