<head><link href="oft_spec.css" rel="stylesheet"></head>

# System Requirement Specification &mdash; Extension Manager

## Introduction

The Extension Manager (EM) is a REST service that manages extensions (e.g. virtual schemas) in an Exasol database.

## About This Document

### Target Audience

The target audience are software developers, requirement engineers, software designers. See section ["Stakeholders"](#stakeholders) for more details.

### Goal

The EM main goal is to simplify installing, configuring and updating extensions for database administrators.

## Stakeholders

### Database Administrators

Database Administrators (DBA) use EM through it's REST API or a user interface for managing extensions in their Exasol databases.

## Terms and Abbreviations

The following list gives you an overview of terms and abbreviations commonly used in OFT documents.

* **EM**: Extension Manager
* **UDF** / **User defined function**: Extension point in the Exasol database that allows users to write their own SQL functions, see the [documentation](https://docs.exasol.com/db/latest/database_concepts/udf_scripts.htm) for details
* **Virtual Schema**: Projection of an external data source that can be access like an Exasol database schema.
* **Virtual Schema adapter**: Plug-in for Exasol based on the Virtual Schema API that translates between Exasol and the data source.
* **Extension**: A user managed extension of the Exasol database (e.g. a Virtual Schema, bulk loaders and other in-database integrations) that consists of multiple parts like files in BucketFS, adapter scripts, connections etc.
* **DBA**: [Database Administrator](#database-administrators)
* **DBO**: A Database Object like tables, views, scripts, connections, virtual schemas etc.

## Features

Features are the highest level requirements in this document that describe the main functionality of EM.

### Install Extensions
`feat~install-extension~1`

EM installs extensions. An extension consists of artifacts and database objects.

#### Install Required Artifacts
`feat~install-extension-artifacts~1`

EM installs all artifacts required by an extension from GitHub (Jar files, 3rd party libraries, JDBC drivers etc.).

Note:

In the initial version this feature is omitted to simplify implementation. We assume that all required artifacts are available in BucketFS.

#### Install Database Objects
`feat~install-extension-database-objects~1`

EM installs corresponding database objects for an extension (UDFs, LUA Scripts, credentials, configuration tables etc.).

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
