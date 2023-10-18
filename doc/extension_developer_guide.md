# Extension Developer Guide

This guide describes how to create and test an extension definition for the Extension Manager.

## Extensions

Definition and implementation of each extension are located in a common repository, e.g. `s3-document-files-virtual-schema`.

The tests for an extensions usually will include
* Using the TypeScript compiler to verify correct implementation of a specific version of the `extension-manager-interface`
* Unit tests written in TypeScript to verify all execution paths of the extension's implementation.
* Integration tests written in Java using a specific version of the `extension-manager` to verify that the extension
  * can be loaded
  * can install a virtual schema and check that it works
  * can update parameters of an existing virtual schema
  * can upgrade a virtual schema created with an older version
  * ...

### Restrictions as Document-based Virtual Schemas Only Support a Single Version

Document-based virtual schemas like `s3-document-files-virtual-schema` require a `SET SCRIPT` that must have a specific name. As this script references a specific virtual schema JAR archive, it is not possible to install multiple version of the same virtual schema in the same database schema.

This means that in order to test a new version of a virtual schema, you need to create a new database schema with the required database objects.

## Extension Manager Interface

Extension definitions are written in TypeScript and compiled to a single JavaScript file. They implement the [extension-manager-interface](https://github.com/exasol/extension-manager-interface/). See [testing-extension](../extension-manager-integration-test-java/testing-extension) for an example including build scripts.

## Extension Integration Test Framework for Java

The Extension Integration Test Framework for Java (EITFJ) allows writing integration tests for extensions and their extension definitions.

### Preconditions

We assume your extension definition project is located in folder `$EXTENSION`. `$EXTENSION_ID` is the filename of your JavaScript extension definition.

The project in `$EXTENSION` must fulfill the following preconditions:
* NPM modules are already installed to `node_modules` before running integration tests.
* `package.json` is configured to build the extension definition with `npm run build`.
* The build process writes the JavaScript file to `$EXTENSION/dist/$EXTENSION_ID`

If your extension definition uses a different build process you can create a custom `ExtensionBuilder`.

### Using EITFJ

The EITFJ library is published to [Maven Central](https://central.sonatype.com/artifact/com.exasol/extension-manager-integration-test-java), so you can add it to your project as follows:

```xml
<dependency>
    <groupId>com.exasol</groupId>
    <artifactId>extension-manager-integration-test-java</artifactId>
    <version>$VERSION</version>
    <scope>test</scope>
</dependency>
```

See [`ExampleIT.java`](../extension-manager-integration-test-java/src/test/java/com/exasol/extensionmanager/ExampleIT.java) for an example of how to use EITFJ in your integration tests. Adapt the following constants depending to your own extension definition:

* `EXTENSION_SOURCE_DIR`: relative path to the directory containing the extension definition sources (`$EXTENSION`)
* `EXTENSION_ID`: file name of the built JavaScript file (`$EXTENSION_ID`)

Depending on the requirements of your extension you might also need to upload the adapter JAR or a JDBC driver to BucketFS in `@BeforeAll`.

#### Preconditions for Using EITFJ

EM only works with Exasol DB version 8 or later and does not support 7.1. `ExtensionManagerSetup.create()` verifies the correct DB version by executing a query against the given Exasol test setup.

When you test your project also with Exasol version 7.1, you can use the following code to skip the extension integration tests for version 7.1:

```java
import com.exasol.extensionmanager.itest;
// ...
exasolTestSetup = ...
ExasolVersionCheck.assumeExasolVersion8(exasolTestSetup);
setup = ExtensionManagerSetup.create(exasolTestSetup, /* ... */);
```

#### Features of Class `ExtensionManagerSetup`

Class `ExtensionManagerSetup` offers the following useful features:

* `ExtensionManagerSetup.create()` downloads and starts the EM REST interface, builds your extension definition and adds it to EM's extension registry.
* Call `setup.client()` to get a client for the EM's REST interface. It allows you to install your extension, create a new instance etc.
* Call `setup.client().assertRequestFails()` to verify that a REST call fails with an expected status code and error message. This allows testing that your extension definitions throws an expected error.
* Call `setup.previousVersionManager()` to prepare a previous version of your extension. This is useful for testing the upgrade process.
* Call `setup.exasolMetadata()` to verify that expected database objects like `SCRIPT`, `CONNECTION` or `VIRTUAL SCHEMA` were created.
* Call `setup.addVirtualSchemaToCleanupQueue()` and `setup.addConnectionToCleanupQueue()` to delete a `CONNECTION` or `VIRTUAL SCHEMA` after a test.

### EITFJ Configuration

EITFJ works without additional configuration. During development you can however create file `extension-test.properties` to simplify local testing. We recommend adding this file to `.gitignore` to avoid accidentally committing it.

`extension-test.properties` supports the following optional settings:

* `localExtensionManager`: Path to a local clone of the `extension-manager` repository. This allows testing against a local version of extension manager that was not yet released. By default EITFJ will install extension manager using `go install`.
* `buildExtension`: Set this to `false` in order to skip building the extension definition before the tests. Use this to speedup tests when the extension definition is not modified.
* `buildExtensionManager`: Set this to `false` to skip building/installing the extension manager binary. Use this to speedup tests when extension manager is not modified.
* `extensionManagerVersion`: Version of EM to use during tests. By default EITFJ uses the same version as the version defined in `pom.xml` for `extension-manager-integration-test-java`. Changing this is not recommended.
