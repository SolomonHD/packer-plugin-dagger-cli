## ADDED Requirements

### Requirement: Production Build Target
The build system SHALL provide a `prod` target that builds production-ready binaries without prerelease suffixes, runs describe test validation, and cleans up to temporary directory.

#### Scenario: Make prod command
- **WHEN** developer runs `make prod`
- **THEN** the system builds plugin with version from `version/VERSION`
- **AND** sets `VersionPrerelease` to empty string in ldflags
- **AND** runs describe test to validate metadata
- **AND** verifies reported version matches `version/VERSION` exactly (no suffix)
- **AND** moves binary to `/tmp/${PACKER_PLUGIN_BIN}` after test

#### Scenario: Production build follows conventions
- **WHEN** building with `make prod`
- **THEN** version is read from `version/VERSION` file
- **AND** git host and organization are derived from `.git` and lowercased
- **AND** `PACKER_PLUGIN_BIN` environment variable is used for binary name
- **AND** all paths use lowercase conventions

#### Scenario: Production describe test validation
- **WHEN** `make prod` completes build
- **THEN** describe command executes on built binary
- **AND** JSON output is validated
- **AND** version field matches `version/VERSION` without any suffix
- **AND** required fields (version, sdk_version, api_version) are present

#### Scenario: Production build failure handling
- **WHEN** `make prod` encounters build or test failure
- **THEN** clear error message is displayed
- **AND** if binary exists, it is moved to `/tmp/${PACKER_PLUGIN_BIN}`
- **AND** cleanup occurs regardless of success or failure

### Requirement: Production Install Target
The build system SHALL provide a `prod-install` target that builds and installs production-ready plugin binaries directly to Packer plugins directory without prerelease suffixes or describe tests.

#### Scenario: Make prod-install command
- **WHEN** developer runs `make prod-install`
- **THEN** the system builds plugin with version from `version/VERSION`
- **AND** sets `VersionPrerelease` to empty string in ldflags
- **AND** skips describe test (immediate install)
- **AND** installs via `packer plugins install --path`

#### Scenario: Production install follows conventions
- **WHEN** installing with `make prod-install`
- **THEN** version is read from `version/VERSION` file
- **AND** binary name uses `PACKER_PLUGIN_BIN` environment variable
- **AND** git-derived paths are lowercased
- **AND** installation uses correct plugin source path

#### Scenario: Direct production installation
- **WHEN** `make prod-install` is used for quick testing
- **THEN** build and install happen in single step
- **AND** no intermediate describe test is run
- **AND** binary is installed to Packer plugins directory
- **AND** developer can immediately test with Packer templates

### Requirement: Target Documentation
The build system SHALL document the purpose and behavior of production build targets with inline comments.

#### Scenario: Developer reads Makefile
- **WHEN** examining production build targets
- **THEN** comments explain `prod` target builds without prerelease suffix
- **AND** describe test validation for `prod` is documented
- **AND** `prod-install` direct installation behavior is explained
- **AND** usage guidance indicates when to use `prod` vs `prod-install`