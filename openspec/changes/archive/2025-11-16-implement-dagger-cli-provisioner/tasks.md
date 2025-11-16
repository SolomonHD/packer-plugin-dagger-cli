# Implementation Tasks

## 1. Project Setup
- [x] 1.1 Create `version/VERSION` file with initial version (0.1.0)
- [x] 1.2 Create `provisioner/dagger-cli/` directory structure
- [ ] 1.3 Remove old scaffolding provisioner references (keeping for backwards compatibility)

## 2. Core Configuration
- [x] 2.1 Define Config struct with all fields (module, file, call, args, env, working_dir, log_level, retry_count, retry_backoff, cache_mode, cache_key)
- [x] 2.2 Implement ConfigSpec() for HCL2 binding
- [x] 2.3 Generate HCL2 spec with go:generate directive
- [x] 2.4 Create config validation unit tests

## 3. Validation Logic (Prepare)
- [x] 3.1 Implement module vs file mutual exclusivity check
- [x] 3.2 Implement call required validation
- [x] 3.3 Implement log_level normalization and validation
- [x] 3.4 Implement retry_count and retry_backoff validation
- [x] 3.5 Implement cache_mode and cache_key validation
- [x] 3.6 Implement effective cache key computation
- [x] 3.7 Add validation unit tests for all error cases

## 4. Command Construction
- [x] 4.1 Create command builder helper function (buildCommand)
- [x] 4.2 Implement module/file to -m flag resolution
- [x] 4.3 Implement args map to CLI flags conversion
- [x] 4.4 Implement environment variable overlay logic
- [x] 4.5 Add cache mode environment variable logic
- [x] 4.6 Add command construction unit tests

## 5. Execution Logic (Provision)
- [x] 5.1 Create command runner interface for testability (CommandRunner)
- [x] 5.2 Implement working directory resolution
- [x] 5.3 Implement retry loop with backoff (linear backoff)
- [x] 5.4 Integrate with Packer UI for logging
- [x] 5.5 Add execution unit tests with fake runners (TestRetryLogic)

## 6. Observability
- [x] 6.1 Create logger helper with [dagger-cli] prefix (logSummary, logAttempt)
- [x] 6.2 Implement pre-execution summary log
- [x] 6.3 Implement per-attempt logging (start/finish/duration)
- [x] 6.4 Implement retry failure logging
- [x] 6.5 Respect log_level in output verbosity (basic implementation)
- [-] 6.6 Sanitize sensitive data in logs (deferred - basic implementation logs safely)

## 7. Testing
- [x] 7.1 Write unit tests for config validation (TestConfigValidation_*)
- [x] 7.2 Write unit tests for command construction (TestCommandConstruction_*)
- [x] 7.3 Write unit tests for retry logic (TestRetryLogic)
- [x] 7.4 Write integration tests with fake command runner (FakeCommandRunner)
- [ ] 7.5 Write Packer acceptance tests (deferred - requires real Dagger setup)
- [ ] 7.6 Test all HCL examples from docs (manual testing recommended)

## 8. Registration
- [x] 8.1 Update main.go to register dagger-cli provisioner
- [x] 8.2 Update go.mod module name (github.com/hashicorp/packer-plugin-dagger-cli)
- [x] 8.3 Update GNUmakefile NAME variable to dagger-cli
- [x] 8.4 Test local build and describe validation (confirmed working)

## 9. Documentation
- [x] 9.1 Create docs/provisioners/dagger-cli.mdx
- [x] 9.2 Add basic usage example (module + call)
- [x] 9.3 Add local module example
- [x] 9.4 Add file-based example
- [x] 9.5 Add observability/retry/caching example
- [x] 9.6 Document all configuration fields
- [x] 9.7 Create example/ templates (dagger-cli-basic.pkr.hcl)
- [x] 9.8 Update README.md with dagger-cli description

## 10. Build & Release
- [x] 10.1 Test local build with version injection from VERSION file
- [x] 10.2 Run describe test to validate plugin metadata (version 0.1.0, provisioners includes "dagger-cli")
- [x] 10.3 Verify binary naming follows conventions (packer-plugin-dagger-cli)
- [x] 10.4 Binary moved to /tmp/packer-plugin-dagger-cli per versioning rules
- [x] 10.5 .goreleaser.yml verified (uses template variables, no changes needed)