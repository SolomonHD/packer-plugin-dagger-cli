# Implementation Tasks

## 1. Build Script Implementation
- [x] 1.1 Create version reading function/logic
  - Read from `version/VERSION` file
  - Handle missing file with clear error
  - Trim whitespace from version string
- [x] 1.2 Implement git path derivation
  - Extract host and org from `git remote get-url origin`
  - Lowercase both values
  - Handle parse failures gracefully
- [x] 1.3 Set up environment variables
  - Define `PACKER_PLUGIN_BIN` with lowercase plugin name
  - Use in build output path
  - Document in script comments
- [x] 1.4 Implement build command
  - Construct ldflags with version from VERSION file
  - Support configurable VersionPrerelease (dev vs empty)
  - Use actual module path from go.mod for ldflags
  - Output to `$PACKER_PLUGIN_BIN`

## 2. Describe Test Integration
- [x] 2.1 Add describe test execution
  - Run `./$PACKER_PLUGIN_BIN describe` after build
  - Capture output and exit code
- [x] 2.2 Implement validation
  - Check JSON is well-formed
  - Verify version matches VERSION file content
  - Confirm required fields present (version, sdk_version, api_version)
- [x] 2.3 Add error handling
  - Report describe failures clearly
  - Include command output in errors
  - Continue to cleanup even on failure

## 3. Cleanup and Reporting
- [x] 3.1 Implement binary cleanup
  - Move `$PACKER_PLUGIN_BIN` to `/tmp/$PACKER_PLUGIN_BIN`
  - Execute cleanup regardless of test outcome
  - Verify repo root is clean
- [x] 3.2 Add status reporting
  - Report build success/failure
  - Report describe test results
  - Report final binary location
  - Include version information in success message

## 4. Makefile Integration
- [x] 4.1 Update `dev` target in GNUmakefile
  - Integrate new build workflow
  - Preserve or improve existing functionality
  - Maintain backward compatibility where possible
- [x] 4.2 Add documentation comments
  - Explain version sourcing
  - Document environment variables
  - Note requirements (git remote, VERSION file)
- [x] 4.3 Consider separate `test-build` target
  - Created `dev-install` target for plugin installation
  - Keep `dev` for testing workflow

## 5. Testing and Validation
- [x] 5.1 Test version reading
  - Verify correct version from VERSION file (0.1.0)
  - Test error handling for missing file
- [x] 5.2 Test git path derivation
  - Tested with git remote URL
  - Verified lowercase conversion (github.com/solomonhd)
- [x] 5.3 Test describe validation
  - Built plugin and ran describe
  - Verified JSON parsing works
  - Confirmed version matching logic (0.1.0-dev)
- [x] 5.4 Test cleanup behavior
  - Verified binary moves to /tmp
  - Checked cleanup occurs after test
- [x] 5.5 Integration test
  - Ran full `make dev` workflow
  - Verified all steps complete correctly
  - Confirmed final state matches requirements

## 6. Documentation
- [x] 6.1 Add inline script comments
  - Explained each major step
  - Documented variables and their purpose
  - Noted assumptions and requirements
- [x] 6.2 Update README if needed
  - Not needed - inline documentation is sufficient
  - Build workflow is self-documenting in Makefile