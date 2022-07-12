# Introduction

## Acknowledgments

This document's section structure is derived from the "[arc42](https://arc42.org/)" architectural template by Dr. Gernot Starke, Dr. Peter Hruschka.

# Constraints

This section introduces technical system constraints.

# Solution Strategy

## Requirement Overview

Please refer to the [System Requirement Specification](system_requirements.md) for user-level requirements.

# Building Blocks

# Runtime

# Cross-cutting Concerns

# Design Decisions

## Do we Need a Backend?

One option would be to implement everything in the JavaScript client. However, we discarded that option, since it does not allow us to upgrade the installed adapters automatically. An automated job can't run in a browser.

# Quality Scenarios

# Risks
