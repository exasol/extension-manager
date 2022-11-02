# Developers Guide

## Building

To build the binary, run

```shell
go generate ./...
go build -o extension-manager cmd/main.go
```

To run the extension manager, execute

```shell
go run cmd/main.go
```

After starting the server you can get the OpenApi definition by executing

```shell
curl "http://localhost:8080/openapi.json" -o extension-manager-api.json
```

## Testing

The different components of the project are responsible for testing different
things.

### extension-manager

The extension-manager project contains unit and integration tests that verify
* Loading and executing of JavaScript extensions
* Database interactions
* REST API interface
* Server-side parameter validation using `extension-parameter-validator`
* ...

Tests use dummy extensions, no real extensions.

### Extensions

Extensions are located in the repositories of the virtual schema implementations, e.g. `s3-document-files-virtual-schema`.

Tests for extensions are:
* Verify correct implementation of a specific version of the
  `extension-manager-interface` using the TypeScript compiler
* Unit tests written in TypeScript verify all execution paths of the extension
* Integration tests written in Java use a specific version of the
  `extension-manager` to verify that the extension
  * can be loaded
  * can install a virtual schema and check that it works
  * can update parameters of an existing virtual schema
  * can upgrade a virtual schema created with an older version
  * ...

### Restrictions as Document Virtual Schemas Only Support a Single Version

Document virtual schemas like `s3-document-files-virtual-schema` require a
`SET SCRIPT` that must have a specific name. As this script references a
specific virtual schema JAR archive, it is not possible to install multiple
version of the same virtual schema in the same database `SCHEMA`.

This means that in order to test a new version of a virtual schema, you need
to create a new `SCHEMA` with the required database objects.

### Non-Parallel Tests

The tests of this project use the exasol-test-setup-abstraction-server. There
the tests connect to an Exasol database running in a docker container.  For
performance reasons the test-setup-abstraction reuses that container.  This
feature is not compatible with running tests in parallel.

Problems would be:

* Name conflicts, e.g. schema names
* Missing isolation, e.g. (`EXA_ALL_SCRIPTS`) contains objects from other tests
* Issues with the exasol-test-setup-abstraction-server (the download of the
  server jar is triggered by the first test. The second one tries to use the
  unfinished jar)

For that reason parallel tests are currently disabled in the CI with `-p 1`.

To run test locally use:

```shell
go test -p 1 ./...
```

To run only unit tests use:

```shell
go test -short ./...
```

## Static Code Analysis

### Go Linter

To install golangci-lint on your machine, follow [these instruction](https://golangci-lint.run/usage/install/#local-installation). Then run

```shell
golangci-lint run
```

### Sonar

Download sonar-scanner as a zip file from [sonarqube.org](https://docs.sonarqube.org/latest/analysis/scan/sonarscanner/) and unpack it.

Run tests to generate code coverage information:

```shell
go test -v -p 1 -count 1 -coverprofile=coverage.out ./...
mvn verify
```

Then run Sonar with the following command in the project root:

```shell
sonar-scanner -Dsonar.organization=exasol -Dsonar.host.url=https://sonarcloud.io -Dsonar.login=$SONAR_TOKEN
```

## Using a Local Extension Interface

To use a local, non-published version of the extension interface in integration tests, edit [pkg/integrationTesting/extensionForTesting/package.json](./../pkg/integrationTesting/extensionForTesting/package.json) and replace the version of `"@exasol/extension-manager-interface"` with the path to your local clone of [extension-manager-interface](https://github.com/exasol/extension-manager-interface).

Make sure to not commit the modified `package.json`.

## Extension Registry

The extension registry is an HTTPS service that provides a JSON file containing links to all available extensions. It consists of an S3 Bucket and a CloudFront distribution deployed via CDK.

### Initial Configuration

1. Create file `registry/lib/config.ts` with the following content:
    ```ts
    export const CONFIG = {
        owner: 'your.email@example.com'
    }
    ```
2. Run `npm install`
3. Configure AWS profile and region:
    ```sh
    export AWS_PROFILE=<profile>
    export AWS_REGION=eu-central-1
    ```

### Deploy Changes

Run `npm run cdk diff`. If the output looks good, run `npm run cdk deploy`.

To get the output variables of the deployed stack (e.g. bucket name and CloudFront distribution host name), run the following command:

```sh
aws cloudformation describe-stacks --stack-name ExtensionManagerRegistry --query "Stacks[0].Outputs[].{key:ExportName,value:OutputValue}"
```

### Deploy Registry Content

To deploy the content of the Extension Registry to `test` or `prod` stage, run script `AWS_PROFILE=<profile> upload-registry.sh test|prod`. This will upload one of the json files in the `content` folder to the S3 bucket and invalidate the CloudFront cache.
