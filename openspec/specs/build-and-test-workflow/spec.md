# build-and-test-workflow Specification

## Purpose
TBD - created by archiving change add-local-build-testing-workflow. Update Purpose after archive.
## Requirements
### Requirement: Single Version Source
The build system SHALL read the plugin version from `version/VERSION` file as the single source of truth for all version-related operations.

#### Scenario: Version file exists and is valid
- **WHEN** building the plugin locally
- **THEN** the build system reads version from `version/VERSION`
- **AND** uses this version in ldflags
- **AND** the built binary reports this version via `describe`

#### Scenario: Version file missing
- **WHEN** `version/VERSION` file does not exist
- **THEN** the build process fails with a clear error message
- **AND** the error indicates the missing file path

### Requirement: Git-Derived Paths
The build system SHALL derive git host and organization from `.git` configuration and normalize them to lowercase for all plugin paths and module references.

#### Scenario: Git remote configured
- **WHEN** building the plugin
- **THEN** the system extracts host and organization from `git remote get-url origin`
- **AND** converts both to lowercase
- **AND** uses lowercase values in module paths and plugin source paths

#### Scenario: Mixed-case git remote
- **WHEN** git remote URL contains mixed case (e.g., `GitHub.com/MyOrg`)
- **THEN** the system normalizes to lowercase (`github.com/myorg`)
- **AND** uses normalized values in build artifacts

### Requirement: Environment Variable for Binary Name
The build system SHALL use an environment variable to specify the plugin binary name, enabling consistent reference throughout the build and test process.

#### Scenario: Binary name environment variable set
- **WHEN** `PACKER_PLUGIN_BIN` environment variable is set to lowercase plugin name
- **THEN** the build outputs to `$PACKER_PLUGIN_BIN`
- **AND** subsequent test commands reference `$PACKER_PLUGIN_BIN`
- **AND** cleanup operations use `$PACKER_PLUGIN_BIN`

#### Scenario: Default binary name
- **WHEN** no binary name environment variable is provided
- **THEN** the system uses `packer-plugin-<name>` as default
- **AND** the plugin name portion is lowercase

### Requirement: Build with Configurable Prerelease
The build system SHALL support building with or without prerelease suffixes via ldflags, enabling both development and release builds.

#### Scenario: Development build with prerelease
- **WHEN** building for local development
- **THEN** the system sets `VersionPrerelease` via ldflags (e.g., `dev`)
- **AND** `describe` output shows version with prerelease suffix

#### Scenario: Clean build without prerelease
- **WHEN** building for release verification
- **THEN** the system sets `VersionPrerelease` to empty string via ldflags
- **AND** `describe` output shows clean version without suffix

### Requirement: Automatic Describe Test
The build system SHALL automatically run the `describe` command on the built binary to verify metadata and version reporting.

#### Scenario: Successful describe test
- **WHEN** the build completes successfully
- **THEN** the system executes `./$PACKER_PLUGIN_BIN describe`
- **AND** validates JSON output is well-formed
- **AND** confirms reported version matches `version/VERSION`
- **AND** verifies required fields are present (version, sdk_version, api_version)

#### Scenario: Describe test failure
- **WHEN** the `describe` command fails or returns invalid JSON
- **THEN** the build process reports the failure clearly
- **AND** includes describe output in error message
- **AND** still proceeds to cleanup step

### Requirement: Binary Cleanup to Temporary Directory
The build system SHALL move the test binary to `/tmp` after describe test completion, preventing repository pollution.

#### Scenario: Successful build and test
- **WHEN** build and describe test complete successfully
- **THEN** the system moves `$PACKER_PLUGIN_BIN` to `/tmp/$PACKER_PLUGIN_BIN`
- **AND** the repo root contains no leftover binaries

#### Scenario: Failed build or test
- **WHEN** build or describe test fails
- **THEN** if binary exists, it is moved to `/tmp/$PACKER_PLUGIN_BIN`
- **AND** cleanup occurs regardless of success/failure

### Requirement: Clear Status Reporting
The build system SHALL provide clear, actionable feedback about build and test results to the developer.

#### Scenario: All steps successful
- **WHEN** build, describe test, and cleanup all succeed
- **THEN** the system reports success status
- **AND** confirms version reported by describe
- **AND** indicates binary location in `/tmp`

#### Scenario: Build failure
- **WHEN** go build fails
- **THEN** the system reports build failure
- **AND** includes relevant error output
- **AND** indicates no binary was produced

#### Scenario: Describe test failure
- **WHEN** build succeeds but describe fails
- **THEN** the system reports describe failure specifically
- **AND** includes describe output or error
- **AND** indicates binary was moved to `/tmp`

### Requirement: Make Target Integration
The build system SHALL integrate with existing Makefile targets, preserving or enhancing the `dev` target functionality.

#### Scenario: Make dev command
- **WHEN** developer runs `make dev`
- **THEN** the system performs version-aware build
- **AND** runs describe test
- **AND** moves binary to `/tmp`
- **AND** provides clear status output

#### Scenario: Backward compatibility
- **WHEN** developer relies on existing `make dev` behavior
- **THEN** the new workflow maintains expected outcomes
- **AND** installation step still works if needed
- **AND** no breaking changes to command interface

### Requirement: Build Documentation
The build system SHALL include inline documentation explaining usage, environment variables, and expected behavior.

#### Scenario: Developer reads build script
- **WHEN** examining build tooling
- **THEN** comments explain version sourcing
- **AND** environment variable usage is documented
- **AND** describe test purpose is clear
- **AND** cleanup behavior is explained

#### Scenario: Makefile help
- **WHEN** developer needs build guidance
- **THEN** Makefile includes comments for key targets
- **AND** critical requirements (git remote, version file) are noted

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

