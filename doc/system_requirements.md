 System Requirement Specification &mdash; Extension Manager

## Introduction

The Extension Manager (EM) is a REST service that manages extensions
(e.g. virtual schemas) in an Exasol database.

## About This Document

### Target Audience

The target audience are software developers, requirement engineers, software
designers. See section ["Stakeholders"](#stakeholders) for more details.

### Goal

The EM main goal is to simplify installing, configuring and updating
extensions for database administrators.

## Stakeholders

### Database Administrators

Database Administrators (DBA) use EM through it's REST API or a user interface
for managing extensions in their Exasol databases.

## Terms and Abbreviations

The following list gives you an overview of terms and abbreviations commonly
used in the requirements specification.

* **EM**: Extension Manager
* **UDF** / **User defined function**: Extension point in the Exasol database
  that allows users to write their own SQL functions, see
  [UDF documentation](https://docs.exasol.com/db/latest/database_concepts/udf_scripts.htm)
  for details
* **Virtual Schema**: Projection of an external data source that can be
  accessed like an Exasol database schema.
* **Virtual Schema adapter**: Plug-in based on the Virtual Schema
  API that translates between Exasol and the original data source.
* **Extension**: A user managed extension of the Exasol database (e.g. a
  Virtual Schema, bulk loaders and other in-database integration).
  An extension might consist of multiple parts e.g. files in BucketFS, adapter
  scripts, connections.
* **DBA**: [Database Administrator](#database-administrators) (role)
* **DBO**: A Database Object, e.g. a table, view, script, connection, virtual
  schema.

## Features

Features are the highest level requirements in this document that describe the main functionality of EM.

### List Extensions
`feat~list-extensions~1`

EM lists extensions.

Needs: dsn

<!--
list only valid, complete extensions
based on extension definitions

-->


### Install Extensions
`feat~install-extension~1`

EM allows the user to install extensions.

Needs: req

### Configure an Extension
`feat~configure-extension~1`

EM allows the user to configure an extension, e.g. in order to set up a
Virtual Schema source system.

Needs: req

### Updating an Extension
`feat~update-extension~1`

EM allows users to install a new version of an extension that was already
installed in an older version.

Needs: req

### Uninstall an Extension
`feat~uninstall-extension~1`

EM allows the user to uninstall an extension.

Needs: req

<!-- ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ -->
## Functional Requirements

### Extensions
`req~extension~1`

An extension might consist of JDBC driver, artifacts, configuration and
database objects. Dependening on the nature of the extension not all artifacts
might be required.

((see github project extension-manager-interface
for parameter conditions))

Covers:
* [`feat~install-extension~1`](#configure-an-extension)
* [`feat~configure-extension~1`](#updating-an-extension)
* [`feat~uninstall-extension~1`](#extensions)

Needs: dsn

### Installation

#### Install Required Artifacts
`req~install-extension-artifacts~1`

EM installs all artifacts required by an extension from GitHub (Jar files, 3rd
party libraries, JDBC drivers etc.).

Covers:

* [`feat~install-extension~1`](#configure-an-extension)

Needs: dsn

#### Install Database Objects
`req~install-extension-database-objects~1`

<!--

See https://github.com/exasol/s3-document-files-virtual-schema/blob/main/doc/user_guide/user_guide.md

1. CREATE SCHEMA ADAPTER;
2. create adapter script
   - in database schema "ADAPTER"
   - name ADAPTER.S3_FILES_ADAPTER
   - reference to jar in bucketfs
3. create UDF function = "set script"
   - in database schema "ADAPTER"
   - name ADAPTER.IMPORT_FROM_S3_DOCUMENT_FILES
   - columns
     - DATA_LOADER
     - SCHEMA_MAPPING_REQUEST
     - CONNECTION_NAME
   - reference to jar file in bucketfs
4. create connection
   - parameters depending on VS variant
   - incl. credentials
5. create virtual schema
   - ADAPTER.S3_FILES_ADAPTER WITH
   - reference to connection
   - MAPPING in bucket fs

-->

EM installs corresponding database objects for an extension (UDFs, LUA
Scripts, credentials, configuration tables etc.).

(( is UDF really a Database object?
Are there lua scripts, indeed?
lua is inline in extension definition either adapter script or set set
))

Covers:

* [`feat~install-extension~1`](#configure-an-extension)

Needs: dsn

### Update extension
`req~update-extension~1`

When updating an extension that is already installed in an older version, EM
checks if any parameter definition has changed. If there were breaking
changes, EM cannot perform the update automatically and aborts the
installation with an appropriate error message.

(( how to define "breaking" changes? ))

Rationale: The only option would be to add update scripts that define how to
convert the parameters from one version to another. However, that is currently
out of scope.


Covers:
* [`feat~update-extension~1`](#uninstall-an-extension)

Needs: dsn

### Uninstalling extensions
`req~uninstall-extension~1`

Covers:
* [`feat~uninstall-extension~1`](#extensions)

Needs: dsn

### Define Configuration Parameters
`req~define-configuration-parameters~1`

EM allows extensions to define a set of parameters. Each extension might have
different parameters.

Covers:
* [`feat~configure-extension~1`](#updating-an-extension)

Needs: dsn

#### Parameter Types
`req~parameter-types~1`

EM supports the following types for configuration parameters
* strings
* select a single option from a given list of available values
* mandatory parameters
* optional parameters
* conditional parameters, i.e. parameters depending on other parameter's values

Rationale: EM's UI can then present all relevant parameters to the user and
allow the user to assign a value to each parameter, e.g.  enter credentials,
select values from option lists.

Covers
* [`feat~configure-extension~1`](#updating-an-extension)

Needs: dsn

#### Validation of parameter values
`req~validate-parameter-values~1`

EM validates parameter values selected oder entered by the user.

Rationale:
* improve user experience
* detect errors as soon as possible
* ensure security

Covers:
* [`feat~configure-extension~1`](#updating-an-extension)

Needs: dsn

## Non-functional Requirements

### UI Languages
`req~ui-languages~1`

The current SaaS implementation only supports English as language in the user
interface. To avoid complexity EM currently only supports English language in
the user interface, too.  This avoids additional efforts for UI translation
until this is required.

Covers:
* [`feat~configure-extension~1`](#updating-an-extension)
* [`feat~install-extension~1`](#configure-an-extension)
* [`feat~uninstall-extension~1`](#extensions)

## Constraints

### EM Uses a Dedicated Schema
`const~use-reserved-schema~1`

EM manages extensions only in a reserved, Exasol-controlled schema called `EXA_EXTENSION`.

Rationale:

This makes it clear to DBAs that objects in this schema are managed and should not be modified by hand.

Needs: impl, utest, itest

### EM works with Exasol SaaS
`const~works-with-saas~1`

EM works in an Exasol SaaS environment.

Needs: impl, utest, itest

## Out-of-Scope

### Downgrade to an Older Version

EM does not need to support downgrading an installed extension to an older version.

## Potential Future Enhancements

The following features are out of scope for now, but the architecture of the
extension manager must be prepared to support them.

### On Premise Support

The extension manager is also useful for customer using Exasol on premise.

### Automatic Installation of Required Files

Currently the extension manager expects that required files like virtual
schema JARs or JDBC drivers are already available in BucketFS. A future
version might download and install these files automatically or update them to
the latest version.

### Automatic Updates of Installed Extensions

When new versions of a virtual schema become available that potentially fix
security issues, it would be helpful to automatically install the new version
and update the virtual schemas during a maintenance window.

### Install Older Versions

If an older version of the VS JAR is available on BucketFS we could allow the
user to choose the version they want to use.

Currently the extension manager will always use the Jar matching the extension's version.

### Use Custom Adapter Scripts

In case the user wants to use an existing `ADAPTER SCRIPT` or explicitly wants
to create a new one, the extension manager allows them to choose between these
options.

### Usage via a Automation Tools

Currently the extension manager is only meant to be used by the SaaS UI. In
the future we might allow using the extension manager also via automation
tools like Terraform that automatically configure virtual schemas.
