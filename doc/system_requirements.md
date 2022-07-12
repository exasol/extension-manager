# System Requirement Specification &mdash; Extension Manager

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

### Install Extensions
`feat~install-extension~1`

EM allows the user to install extensions.
An extension consists of artifacts and database objects.

Needs: req


#### Install Required Artifacts
`req~install-extension-artifacts~1`

EM installs all artifacts required by an extension from GitHub (Jar files, 3rd
party libraries, JDBC drivers etc.).

Note:

In the initial version the developers assume that all required artifacts are
available in BucketFS to simplify implementation.

Covers:

* `feat~install-extension~1`

Needs impl, utest, itest

#### Install Database Objects
`req~install-extension-database-objects~1`

EM installs corresponding database objects for an extension (UDFs, LUA
Scripts, credentials, configuration tables etc.).

Covers:

* `feat~install-extension~1`

Needs: req

### Configure an Extension
`feat~configure-extension~1`

EM configures an extension (e.g. setup a Virtual Schema source system).

Needs: req

### Uninstall an Extension
`feat~uninstall-extension~1`

EM uninstalls an extension including it's artifacts, configuration and database objects.

Needs: req

### Upgrade to Latest Version
`feat~upgrade-extension~1`

EM upgrades an extension to it's latest version.

Needs: req

### List Installed Extensions
`feat~list-extension~1`

EM lists installed extensions.

Needs: req

## Functional Requirements

## Non-functional Requirements

### UI Languages
`feat~ui-languages~1`

The current SaaS implementation only supports English as language in the user
interface. To avoid complexity EM currently only supports English language in
the user interface, too.  This avoids additional efforts for UI translation
until this is required.

## Constraints

### EM Uses a Dedicated Schema
`const~use-reserved-schema~1`

EM manages extensions only in a reserved, Exasol-controlled schema called `EXA_EXTENSION`.

Rationale:

This makes it clear to DBAs that objects in this schema are managed and should not be modified by hand.

Needs: dsn

### EM works with Exasol SaaS
`const~works-with-saas~1`

EM works in an Exasol SaaS environment.

Needs: dsn

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
