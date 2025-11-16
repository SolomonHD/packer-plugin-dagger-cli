# Implementation Tasks

## 1. Project Setup
- [ ] 1.1 Create `version/VERSION` file with initial version (0.1.0)
- [ ] 1.2 Create `provisioner/dagger-cli/` directory structure
- [ ] 1.3 Remove old scaffolding provisioner references

## 2. Core Configuration
- [ ] 2.1 Define Config struct with all fields (module, file, call, args, env, working_dir, log_level, retry_count, retry_backoff, cache_mode, cache_key)
- [ ] 2.2 Implement ConfigSpec() for HCL2 binding
- [ ] 2.3 Generate HCL2 spec with go:generate directive
- [ ] 2.4 Create config validation unit tests

## 3. Validation Logic (Prepare)
- [ ] 3.1 Implement module vs file mutual exclusivity check
- [ ] 3.2 Implement call required validation
- [ ] 3.3 Implement log_level normalization and validation
- [ ] 3.4 Implement retry_count and retry_backoff validation
- [ ] 3.5 Implement cache_mode and cache_key validation
- [ ] 3.6 Implement effective cache key computation
- [ ] 3.7 Add validation unit tests for all error cases

## 4. Command Construction
- [ ] 4.1 Create command builder helper function
- [ ] 4.2 Implement module/file to -m flag resolution
- [ ] 4.3 Implement args map to CLI flags conversion
- [ ] 4.4 Implement environment variable overlay logic
- [ ] 4.5 Add cache mode environment variable logic
- [ ] 4.6 Add command construction unit tests

## 5. Execution Logic (Provision)
- [ ] 5.1 Create command runner interface for testability
- [ ] 5.2 Implement working directory resolution
- [ ] 5.3 Implement retry loop with backoff
- [ ] 5.4 Integrate with Packer UI for logging
- [ ] 5.5 Add execution unit tests with fake runners

## 6. Observability
- [ ] 6.1 Create logger helper with [dagger-cli] prefix
- [ ] 6.2 Implement pre-execution summary log
- [ ] 6.3 Implement per-attempt logging (start/finish/duration)
- [ ] 6.4 Implement retry failure logging
- [ ] 6.5 Respect log_level in output verbosity
- [ ] 6.6 Sanitize sensitive data in logs

## 7. Testing
- [ ] 7.1 Write unit tests for config validation
- [ ] 7.2 Write unit tests for command construction
- [ ] 7.3 Write unit tests for retry logic
- [ ] 7.4 Write integration tests with fake command runner
- [ ] 7.5 Write Packer acceptance tests
- [ ] 7.6 Test all HCL examples from docs

## 8. Registration
- [ ] 8.1 Update main.go to register dagger-cli provisioner
- [ ] 8.2 Update go.mod module name if needed
- [ ] 8.3 Update GNUmakefile NAME variable to dagger-cli
- [ ] 8.4 Test `make dev` local installation

## 9. Documentation
- [ ] 9.1 Create docs/provisioners/dagger-cli.mdx
- [ ] 9.2 Add basic usage example (module + call)
- [ ] 9.3 Add local module example
- [ ] 9.4 Add file-based example
- [ ] 9.5 Add observability/retry/caching example
- [ ] 9.6 Document all configuration fields
- [ ] 9.7 Create example/ templates
- [ ] 9.8 Update README.md with dagger-cli description

## 10. Build & Release
- [ ] 10.1 Test local build with version injection
- [ ] 10.2 Run describe test to validate plugin metadata
- [ ] 10.3 Verify binary naming follows conventions
- [ ] 10.4 Test plugin installation with `packer plugins install`
- [ ] 10.5 Update .goreleaser.yml if needed for plugin name