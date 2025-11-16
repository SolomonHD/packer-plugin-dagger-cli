# Implementation Tasks

## 1. Build Script Implementation
- [ ] 1.1 Create version reading function/logic
  - Read from `version/VERSION` file
  - Handle missing file with clear error
  - Trim whitespace from version string
- [ ] 1.2 Implement git path derivation
  - Extract host and org from `git remote get-url origin`
  - Lowercase both values
  - Handle parse failures gracefully
- [ ] 1.3 Set up environment variables
  - Define `PACKER_PLUGIN_BIN` with lowercase plugin name
  - Use in build output path
  - Document in script comments
- [ ] 1.4 Implement build command
  - Construct ldflags with version from VERSION file
  - Support configurable VersionPrerelease (dev vs empty)
  - Use lowercase git-derived module path
  - Output to `$PACKER_PLUGIN_BIN`

## 2. Describe Test Integration
- [ ] 2.1 Add describe test execution
  - Run `./$PACKER_PLUGIN_BIN describe` after build
  - Capture output and exit code
- [ ] 2.2 Implement validation
  - Check JSON is well-formed
  - Verify version matches VERSION file content
  - Confirm required fields present (version, sdk_version, api_version)
- [ ] 2.3 Add error handling
  - Report describe failures clearly
  - Include command output in errors
  - Continue to cleanup even on failure

## 3. Cleanup and Reporting
- [ ] 3.1 Implement binary cleanup
  - Move `$PACKER_PLUGIN_BIN` to `/tmp/$PACKER_PLUGIN_BIN`
  - Execute cleanup regardless of test outcome
  - Verify repo root is clean
- [ ] 3.2 Add status reporting
  - Report build success/failure
  - Report describe test results
  - Report final binary location
  - Include version information in success message

## 4. Makefile Integration
- [ ] 4.1 Update `dev` target in GNUmakefile
  - Integrate new build workflow
  - Preserve or improve existing functionality
  - Maintain backward compatibility where possible
- [ ] 4.2 Add documentation comments
  - Explain version sourcing
  - Document environment variables
  - Note requirements (git remote, VERSION file)
- [ ] 4.3 Consider separate `test-build` target
  - Optional: create distinct target for version-aware testing
  - Keep `dev` for typical development workflow

## 5. Testing and Validation
- [ ] 5.1 Test version reading
  - Verify correct version from VERSION file
  - Test error handling for missing file
- [ ] 5.2 Test git path derivation
  - Test with various git remote formats
  - Verify lowercase conversion
- [ ] 5.3 Test describe validation
  - Build plugin and run describe
  - Verify JSON parsing works
  - Confirm version matching logic
- [ ] 5.4 Test cleanup behavior
  - Verify binary moves to /tmp
  - Check cleanup on both success and failure
- [ ] 5.5 Integration test
  - Run full `make dev` workflow
  - Verify all steps complete correctly
  - Check final state matches requirements

## 6. Documentation
- [ ] 6.1 Add inline script comments
  - Explain each major step
  - Document variables and their purpose
  - Note any assumptions or requirements
- [ ] 6.2 Update README if needed
  - Add build/test workflow section if helpful
  - Document environment variables
  - Explain version behavior