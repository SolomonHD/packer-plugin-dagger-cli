## ADDED Requirements

### Requirement: Canonical Dagger Plugin Identity

The project MUST consistently use the `packer-plugin-dagger` repository, module, binary, and Packer source identity under the `SolomonHD` GitHub namespace.

#### Scenario: Repository and module metadata

- **WHEN** project identity is inspected
- **THEN** the repository SHALL be named `packer-plugin-dagger`
- **AND** the Go module SHALL be `github.com/SolomonHD/packer-plugin-dagger`
- **AND** the Packer source SHALL be `github.com/SolomonHD/dagger`
- **AND** documentation, examples, release configuration, and generated metadata MUST NOT claim the plugin is maintained by HashiCorp

#### Scenario: Binary identity

- **WHEN** the plugin is built for release
- **THEN** its binary name MUST begin with `packer-plugin-dagger_`
- **AND** its `describe` output MUST advertise only the intended Dagger post-processor

### Requirement: Semantic Version for Every Executable Build

Every documented build and plugin-check path MUST inject or default to a valid semantic plugin version before package initialization.

#### Scenario: Plugin compatibility check

- **WHEN** `make plugin-check` builds and describes the plugin
- **THEN** plugin initialization MUST NOT panic with `Malformed version: dev`
- **AND** the compatibility check MUST complete using the version from `version/VERSION`

#### Scenario: Direct developer build

- **WHEN** a developer runs a documented direct build command
- **THEN** the resulting executable MUST start and describe itself with a valid semantic development version

### Requirement: Real Dagger Post-Processor Fixtures

The repository MUST replace scaffolding templates with runnable examples and tests for the Dagger post-processor lifecycle.

#### Scenario: Local fixture

- **WHEN** the local example is initialized and validated with its documented prerequisites
- **THEN** it MUST reference `github.com/SolomonHD/dagger`
- **AND** it MUST invoke `post-processor "dagger"`
- **AND** it MUST use a checked-in minimal Dagger module that validates a synthetic artifact descriptor
- **AND** it MUST contain no scaffolding builder, data source, post-processor, or documentation references

#### Scenario: Amazon EBS example

- **WHEN** a user reads the AWS example
- **THEN** it SHALL show Ansible Navigator configuring an Amazon EBS builder before Dagger validates the registered AMI
- **AND** it SHALL document required AWS permissions, expected cost-bearing resources, cleanup, and opt-in execution
