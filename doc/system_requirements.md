<head><link href="oft_spec.css" rel="stylesheet"></head>

# System Requirement Specification &mdash; Extension Manager

## Introduction

The Extension Manager (EM) is a REST service that manages extensions (e.g. virtual schemas) in an Exasol database.

## About This Document

### Target Audience

The target audience are software developers, requirement engineers, software designers. See section ["Stakeholders"](#stakeholders) for more details.

### Goal

The EM main goal is to simplify managing extensions for customers.

## Stakeholders

### Customers

Customers use EM through it's REST API or a user interface for managing extensions in their Exasol databases.

## Terms and Abbreviations

The following list gives you an overview of terms and abbreviations commonly used in OFT documents.

* **EM**: Extension Manager
* **UDF** / **User defined function**: Extension point in the Exasol database that allows users to write their own SQL functions, see the [documentation](https://docs.exasol.com/db/latest/database_concepts/udf_scripts.htm) for details
* **Virtual Schema**: Projection of an external data source that can be access like an Exasol database schema.
* **Virtual Schema adapter**: Plug-in for Exasol based on the Virtual Schema API that translates between Exasol and the data source.
* **Extension**: A user managed extension of the Exasol database (e.g. a Virtual Schema) that consists of multiple parts like files in BucketFS, adapter scripts, connections etc.

## Features

Features are the highest level requirements in this document that describe the main functionality of EM.

## Functional Requirements

## Non-functional Requirements

### Constraints
