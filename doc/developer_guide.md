# Developer Guide

This guide describes how to develop, test and build the Extension Manager.

## Building

To build the binary, run

```sh
go generate ./...
go build -o extension-manager cmd/main.go
```

To run the extension manager, execute

```sh
# Show supported command line arguments
go run cmd/main.go -h
# Start server with custom extension registry
go run cmd/main.go -serverAddress localhost:8080 -extensionRegistryURL /path/to/extensions/
```

After starting the server you can get the OpenApi definition by executing

```sh
curl "http://localhost:8080/openapi.json" -o extension-manager-api.json
```

## Requirement Tracing

You can run requirements tracing by executing:

```sh
./ci/trace-requirements.sh
```

If tracing fails with a `org.xml.sax.SAXParseException` you might need to run `mvn clean` before to delete temporary files like `target/site/jacoco/jacoco.xml`.

## Testing

The extension-manager project contains unit and integration tests that verify
* Loading and executing of JavaScript extensions
* Database interactions
* REST API interface
* Server-side parameter validation using `extension-parameter-validator`

Tests use dummy extensions, no real extensions.

### Non-Parallel Tests

The tests of this project use [`exasol-test-setup-abstraction-server`](https://github.com/exasol/exasol-test-setup-abstraction-server/). There the tests connect to an Exasol database running in a docker container. For performance reasons the test-setup-abstraction reuses that container. This feature is not compatible with running tests in parallel.

Problems would be:

* Name conflicts, e.g. schema names
* Missing isolation, e.g. `EXA_ALL_SCRIPTS` contains objects from other tests
* Issues with the exasol-test-setup-abstraction-server (the download of the server jar is triggered by the first test. The second one tries to use the unfinished jar)

For that reason parallel tests are currently disabled in the CI with `-p 1`.

To run test locally use:

```sh
go test -p 1 ./...
```

To run only tests without the database use:

```sh
go test -p 1 -short ./...
```

Please note that also `-short` tests need `-p 1` because extension integration tests share a directory for building a test extension. Tests will fail randomly without `-p 1`.

## Static Code Analysis

### Go Linter

To install golangci-lint on your machine, follow [these instruction](https://golangci-lint.run/usage/install/#local-installation).

To run the linter, execute

```sh
golangci-lint run
```

File `.golangci.yml` contains configuration like enabled or disabled linters.

### Sonar

Download sonar-scanner as a zip file from [sonarqube.org](https://docs.sonarqube.org/latest/analysis/scan/sonarscanner/) and unpack it.

Run tests to generate code coverage information:

```sh
go test -v -p 1 -count 1 -coverprofile=coverage.out ./...
mvn verify
```

Then run Sonar with the following command in the project root:

```sh
sonar-scanner -Dsonar.organization=exasol -Dsonar.host.url=https://sonarcloud.io -Dsonar.login=$SONAR_TOKEN
```

## Using a Local Extension Interface

To use a local, non-published version of the extension interface for testing EM follow these steps:

1. Build `extension-manager-interface` by running `npm run build`. This is required after each change.
2. Edit [pkg/integrationTesting/extensionForTesting/package.json](./../pkg/integrationTesting/extensionForTesting/package.json) and replace the version of `"@exasol/extension-manager-interface"` with the path to your local clone of [extension-manager-interface](https://github.com/exasol/extension-manager-interface).
3. Edit [pkg/integrationTesting/extensionForTesting/extensionForTestingTemplate.ts](./../pkg/integrationTesting/extensionForTesting/extensionForTestingTemplate.ts) and adapt it to the new API if necessary.

   **Note:** The file contains placeholders that will be replaced during tests. It is not valid TypeScript, so it's normal that the editor complains about the invalid syntax.

Make sure to not commit the modified `package.json`.

## Extension Registry

The extension registry is an HTTPS service that provides a JSON file containing links to all available extensions. The service consists of an S3 Bucket and a CloudFront distribution deployed via AWS Cloud Development Kit (CDK).

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

To deploy the content of the Extension Registry to `test` or `prod` stage, run the following command:

```sh
cd registry-upload
AWS_PROFILE=$profile npm run upload -- --stage=test --no-dry-run
# or
AWS_PROFILE=$profile npm run upload -- --stage=prod --no-dry-run
```

This will upload the JSON file from the `registry-upload/content` folder for the given stage to the S3 bucket and invalidate the CloudFront cache. It will also upload the `testing-extension.js` extension to the `test` stage.

### Upgrade NPM Dependencies

```sh
npx npm-check-updates -u && npm install
```

## Embedding Extension Manager in Other Go Programs

### Embedding the REST API

You can embed the Extension Manager's REST API in other programs that use the [Nightapes/go-rest](https://github.com/Nightapes/go-rest) library:

```go
// Create an instance of `openapi.API`
api := openapi.NewOpenAPI()
// Create a new configuration object
config := ExtensionManagerConfig{ExtensionRegistryURL: "https://<extension-registry>"}
// Add endpoints
err := restAPI.AddPublicEndpoints(api, config)
// Start the server
```

### Embedding the Controller

If you want to directly use the controller:

```go
controller:=extensionController.CreateWithConfig(extensionController.ExtensionManagerConfig{
    ExtensionRegistryURL: "https://example.com/registry.json", 
    BucketFSBasePath: "/buckets/bfsdefault/default/",
    ExtensionSchema: "EXA_EXTENSIONS",
})
var db *sql.DB // create database connection
extensions, err := controller.GetAllExtensions(context.Background(), db)
// ...
```
