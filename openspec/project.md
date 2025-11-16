# Project Context

## Purpose
HashiCorp Packer plugin that provides a `dagger-cli` provisioner for executing Dagger pipelines during image builds. This plugin shells out to the Dagger CLI to orchestrate build steps, enabling infrastructure-as-code pipeline integration within Packer workflows.

## Tech Stack
- Go 1.23.2
- Packer Plugin SDK v0.6.1
- Dagger CLI (external dependency)
- HCL2 for configuration

## Project Conventions

### Code Style
- Follow standard Go formatting (gofmt)
- Use `mapstructure` tags for HCL config binding
- Prefix all log messages with `[dagger-cli]`
- Use lowercase for all plugin paths, module names, and Git organizations
- Error messages follow pattern: `"dagger-cli: <description>"`

### Architecture Patterns
- Single provisioner plugin following Packer Plugin SDK patterns
- Command execution abstracted behind interfaces for testability
- Configuration validation in `Prepare()` method
- Execution logic in `Provision()` method
- Retry logic with configurable backoff
- Structured logging with configurable levels

### Testing Strategy
- Unit tests for config validation
- Unit tests for command construction
- Integration tests using fake command runners
- Acceptance tests with actual Packer builds
- Test files: `*_test.go` alongside implementation

### Git Workflow
- Follow HashiCorp's MPL-2.0 license
- Use semantic versioning via version/VERSION file
- Build binaries with naming: `packer-plugin-<name>_v<version>_<api>_<os>_<arch>`

## Domain Context
- **Dagger**: Container-based CI/CD pipeline engine with CLI interface
- **Provisioners**: Packer components that install/configure software on machines
- **Module**: Dagger's unit of code/pipeline definition (can be remote GitHub repo or local path)
- **Call**: Dagger's function/entrypoint to execute
- **Args**: Named parameters passed to Dagger functions
- **Caching**: Dagger's built-in layer caching system

## Important Constraints
- Plugin must work with existing Packer workflows (builders, communicators)
- Must not require Dagger daemon or state management (stateless execution)
- Configuration must be HCL2-compatible for Packer templates
- Must support both local and remote Dagger modules
- Log output must integrate with Packer's UI system

## External Dependencies
- Dagger CLI must be installed and available in PATH
- Git (if using remote modules)
- Network access (for remote modules and container pulls)
