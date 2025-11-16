# Change: Add Local Build and Testing Workflow

## Why
The current `make dev` target hard-codes `VersionPrerelease=dev` and doesn't follow project versioning rules. Developers can't easily:
- Build with the actual version from `version/VERSION`
- Verify builds with the `describe` command
- Test version reporting before releases
- Follow the packer-versioning-testing.md conventions (single source of truth, git-derived paths, describe testing)

This creates confusion during local development and makes it difficult to verify version behavior before creating releases.

## What Changes
- Add new `build-and-test-workflow` capability defining build/test requirements
- Create or update build script(s) to:
  - Read version from `version/VERSION` as the single source of truth
  - Derive and lowercase git host/organization from `.git`
  - Use environment variables for binary names (e.g., `PACKER_PLUGIN_BIN`)
  - Run `describe` test automatically after build
  - Move successful builds to `/tmp` for cleanup
  - Provide clear success/failure reporting
- Update `GNUmakefile` to align with new workflow or reference build script
- Add inline documentation for usage

## Impact
- Affected specs: **NEW** `build-and-test-workflow`
- Affected code: 
  - `GNUmakefile` (update `dev` target)
  - New build script at repo root (e.g., `build-local.sh` or integrated into Makefile)
  - Documentation (inline comments or README section)
- Breaking: No (existing `make dev` functionality preserved or improved)
- Dependencies: Assumes `.git` directory exists with configured remote