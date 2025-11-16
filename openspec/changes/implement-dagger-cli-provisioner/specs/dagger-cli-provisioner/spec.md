# Dagger CLI Provisioner Specification

## ADDED Requirements

### Requirement: Module Source Configuration
The provisioner SHALL support specifying Dagger modules via `module` or `file` parameters.

#### Scenario: Remote module via module parameter
- **GIVEN** a Packer template with `module = "github.com/acme/infra//ci?ref=v1.2.3"`
- **WHEN** the provisioner prepares the configuration
- **THEN** the module source SHALL be set to the remote module reference

#### Scenario: Local module via module parameter
- **GIVEN** a Packer template with `module = "./ci"`
- **WHEN** the provisioner prepares the configuration
- **THEN** the module source SHALL be set to the local directory path

#### Scenario: File-based module source
- **GIVEN** a Packer template with `file = "ci/pipeline.py"`
- **WHEN** the provisioner prepares the configuration
- **THEN** the module source SHALL be derived from the file's directory

### Requirement: Module and File Mutual Exclusivity
The provisioner SHALL enforce that `module` and `file` are mutually exclusive configuration parameters.

#### Scenario: Both module and file specified
- **GIVEN** a configuration with both `module` and `file` set
- **WHEN** the Prepare method validates the configuration
- **THEN** it SHALL return error "dagger-cli: module and file are mutually exclusive; specify only one."

#### Scenario: Neither module nor file specified
- **GIVEN** a configuration with neither `module` nor `file` set
- **WHEN** the Prepare method validates the configuration
- **THEN** it SHALL return error "dagger-cli: one of module or file must be provided."

#### Scenario: Only module specified
- **GIVEN** a configuration with only `module` set
- **WHEN** the Prepare method validates the configuration
- **THEN** the configuration SHALL be valid

#### Scenario: Only file specified
- **GIVEN** a configuration with only `file` set
- **WHEN** the Prepare method validates the configuration
- **THEN** the configuration SHALL be valid

### Requirement: Call Parameter Required
The provisioner SHALL require the `call` parameter specifying the Dagger entrypoint.

#### Scenario: Call parameter provided
- **GIVEN** a configuration with `call = "configure_system"`
- **WHEN** the Prepare method validates the configuration
- **THEN** the configuration SHALL be valid

#### Scenario: Call parameter missing
- **GIVEN** a configuration without a `call` parameter
- **WHEN** the Prepare method validates the configuration
- **THEN** it SHALL return error "dagger-cli: call must be provided."

### Requirement: Arguments and Environment Variables
The provisioner SHALL support passing named arguments and environment variables to the Dagger CLI.

#### Scenario: Named arguments as map
- **GIVEN** a configuration with `args = { region = "us-east-1", fast = true, count = 3 }`
- **WHEN** the provisioner constructs the command
- **THEN** it SHALL convert to CLI flags: `--region=us-east-1 --fast --count=3`

#### Scenario: Boolean arguments
- **GIVEN** an argument with value `true`
- **WHEN** the provisioner constructs the command
- **THEN** it SHALL use flag format `--name` without a value

#### Scenario: Boolean false arguments
- **GIVEN** an argument with value `false`
- **WHEN** the provisioner constructs the command
- **THEN** it SHALL omit the flag

#### Scenario: Environment variables overlay
- **GIVEN** a configuration with `env = { ENVIRONMENT = "staging" }`
- **WHEN** the provisioner executes the Dagger CLI
- **THEN** it SHALL merge config env vars onto the process environment

### Requirement: Working Directory Configuration
The provisioner SHALL support overriding the working directory for Dagger CLI execution.

#### Scenario: Custom working directory
- **GIVEN** a configuration with `working_dir = "/custom/path"`
- **WHEN** the provisioner executes the Dagger CLI
- **THEN** it SHALL set the process working directory to "/custom/path"

#### Scenario: Default working directory
- **GIVEN** a configuration without `working_dir` set
- **WHEN** the provisioner executes the Dagger CLI
- **THEN** it SHALL use the Packer default working directory

### Requirement: Log Level Configuration
The provisioner SHALL support configurable log levels for observability.

#### Scenario: Valid log level
- **GIVEN** a configuration with `log_level = "debug"`
- **WHEN** the Prepare method validates the configuration
- **THEN** it SHALL normalize to lowercase and accept the value

#### Scenario: Invalid log level
- **GIVEN** a configuration with `log_level = "invalid"`
- **WHEN** the Prepare method validates the configuration
- **THEN** it SHALL return a validation error with allowed values

#### Scenario: Default log level
- **GIVEN** a configuration without `log_level` set
- **WHEN** the Prepare method validates the configuration
- **THEN** it SHALL default to "info"

#### Scenario: Allowed log levels
- **GIVEN** the provisioner's log level validation
- **WHEN** checking for valid values
- **THEN** it SHALL accept: "trace", "debug", "info", "warn", "error" (case-insensitive)

### Requirement: Retry Configuration
The provisioner SHALL support configurable retry logic with backoff for failed Dagger executions.

#### Scenario: Valid retry count
- **GIVEN** a configuration with `retry_count = 3`
- **WHEN** the Prepare method validates the configuration
- **THEN** the configuration SHALL be valid

#### Scenario: Negative retry count
- **GIVEN** a configuration with `retry_count = -1`
- **WHEN** the Prepare method validates the configuration
- **THEN** it SHALL return a validation error

#### Scenario: Default retry count
- **GIVEN** a configuration without `retry_count` set
- **WHEN** the Prepare method validates the configuration
- **THEN** it SHALL default to 0 (no retries)

#### Scenario: Retry backoff with retries enabled
- **GIVEN** a configuration with `retry_count = 3` and no `retry_backoff`
- **WHEN** the Prepare method validates the configuration
- **THEN** it SHALL default `retry_backoff` to "5s"

#### Scenario: Valid retry backoff duration
- **GIVEN** a configuration with `retry_backoff = "10s"`
- **WHEN** the Prepare method validates the configuration
- **THEN** it SHALL parse the duration successfully

#### Scenario: Invalid retry backoff duration
- **GIVEN** a configuration with `retry_backoff = "invalid"`
- **WHEN** the Prepare method validates the configuration
- **THEN** it SHALL return a parsing error

### Requirement: Cache Mode Configuration
The provisioner SHALL support Dagger caching modes: auto, disabled, and keyed.

#### Scenario: Default cache mode
- **GIVEN** a configuration without `cache_mode` set
- **WHEN** the Prepare method validates the configuration
- **THEN** it SHALL default to "auto"

#### Scenario: Valid cache modes
- **GIVEN** a configuration with `cache_mode` set to one of: "auto", "disabled", "keyed"
- **WHEN** the Prepare method validates the configuration
- **THEN** the configuration SHALL be valid

#### Scenario: Invalid cache mode
- **GIVEN** a configuration with `cache_mode = "invalid"`
- **WHEN** the Prepare method validates the configuration
- **THEN** it SHALL return a validation error with allowed values

#### Scenario: Keyed cache requires cache key
- **GIVEN** a configuration with `cache_mode = "keyed"` and no `cache_key`
- **WHEN** the Prepare method validates the configuration
- **THEN** it SHALL return an error requiring cache_key

#### Scenario: Keyed cache with cache key
- **GIVEN** a configuration with `cache_mode = "keyed"` and `cache_key = "ubuntu-22.04"`
- **WHEN** the Prepare method validates the configuration
- **THEN** it SHALL compute an effective cache key combining cache_key, module/file, call, and args hash

### Requirement: Command Construction
The provisioner SHALL construct Dagger CLI commands following the pattern: `dagger call -m <module> <call> [args...]`

#### Scenario: Basic command with module
- **GIVEN** a configuration with `module = "./ci"` and `call = "build"`
- **WHEN** the provisioner constructs the command
- **THEN** it SHALL produce: `["call", "-m", "./ci", "build"]`

#### Scenario: Command with file-based module
- **GIVEN** a configuration with `file = "ci/pipeline.py"` and `call = "deploy"`
- **WHEN** the provisioner constructs the command
- **THEN** it SHALL produce: `["call", "-m", "ci", "deploy"]`

#### Scenario: Command with arguments
- **GIVEN** a configuration with args `{region: "us-west-2"}`
- **WHEN** the provisioner constructs the command
- **THEN** it SHALL append: `--region=us-west-2`

### Requirement: Retry Execution Logic
The provisioner SHALL implement retry logic with exponential backoff on non-zero exit codes.

#### Scenario: Successful execution on first attempt
- **GIVEN** a configuration with `retry_count = 3`
- **WHEN** the Dagger CLI execution succeeds on first attempt
- **THEN** the provisioner SHALL return success without retrying

#### Scenario: Successful execution after retries
- **GIVEN** a configuration with `retry_count = 3`
- **WHEN** the first 2 attempts fail and the 3rd succeeds
- **THEN** the provisioner SHALL return success after 3 attempts

#### Scenario: All attempts fail
- **GIVEN** a configuration with `retry_count = 2`
- **WHEN** all 3 attempts (1 + 2 retries) fail
- **THEN** the provisioner SHALL return the final error

#### Scenario: Backoff delay calculation
- **GIVEN** a configuration with `retry_backoff = "5s"`
- **WHEN** an attempt fails and needs retry
- **THEN** the provisioner SHALL wait `attempt_number * base_duration` before retrying

### Requirement: Structured Logging
The provisioner SHALL provide structured logging with the `[dagger-cli]` prefix and respect the configured log level.

#### Scenario: Pre-execution summary log
- **GIVEN** a provisioner configuration
- **WHEN** the Provision method starts execution
- **THEN** it SHALL log a summary with: module/file, call, total attempts, log_level, cache_mode, cache_key

#### Scenario: Per-attempt logging
- **GIVEN** an execution attempt
- **WHEN** the attempt starts and finishes
- **THEN** it SHALL log the attempt number, command preview (sanitized), duration, and exit code

#### Scenario: Retry logging
- **GIVEN** a failed attempt with retries remaining
- **WHEN** the attempt fails
- **THEN** it SHALL log: "attempt X/N failed (exit_code=Y). Retrying in Zs…"

#### Scenario: Final failure logging
- **GIVEN** the final attempt fails with no retries remaining
- **WHEN** the attempt completes
- **THEN** it SHALL log: "attempt X/X failed (exit_code=Y). No more retries."

#### Scenario: Log level filtering
- **GIVEN** a log_level of "warn"
- **WHEN** the provisioner logs messages
- **THEN** it SHALL only output warn and error level messages

### Requirement: Cache Mode Environment Variables
The provisioner SHALL configure Dagger caching via environment variables based on cache_mode.

#### Scenario: Auto cache mode
- **GIVEN** a configuration with `cache_mode = "auto"`
- **WHEN** the provisioner executes Dagger CLI
- **THEN** it SHALL NOT set any special cache environment variables

#### Scenario: Disabled cache mode
- **GIVEN** a configuration with `cache_mode = "disabled"`
- **WHEN** the provisioner executes Dagger CLI
- **THEN** it SHALL set environment variable `DAGGER_CACHE_DISABLED=1`

#### Scenario: Keyed cache mode
- **GIVEN** a configuration with `cache_mode = "keyed"` and computed effective cache key
- **WHEN** the provisioner executes Dagger CLI
- **THEN** it SHALL set environment variable `DAGGER_CACHE_KEY=<effective_cache_key>`

### Requirement: Error Messages
The provisioner SHALL provide clear, prefixed error messages for all validation and execution failures.

#### Scenario: Validation error format
- **GIVEN** any configuration validation failure
- **WHEN** an error is returned
- **THEN** the error message SHALL start with "dagger-cli: "

#### Scenario: Execution error context
- **GIVEN** a Dagger CLI execution failure
- **WHEN** the provisioner returns an error
- **THEN** it SHALL include the exit code and attempt number

### Requirement: Packer SDK Integration
The provisioner SHALL integrate with the Packer Plugin SDK following standard provisioner patterns.

#### Scenario: ConfigSpec implementation
- **GIVEN** the provisioner struct
- **WHEN** ConfigSpec() is called
- **THEN** it SHALL return the HCL2 object spec for config binding

#### Scenario: Prepare implementation
- **GIVEN** raw configuration data
- **WHEN** Prepare() is called
- **THEN** it SHALL decode, validate, and store the configuration

#### Scenario: Provision implementation
- **GIVEN** a Packer build context
- **WHEN** Provision() is called with UI and Communicator
- **THEN** it SHALL execute the Dagger CLI and report progress via UI

#### Scenario: Generated HCL2 spec
- **GIVEN** the Config struct with mapstructure tags
- **WHEN** go generate is run
- **THEN** it SHALL generate the HCL2 spec file for config binding