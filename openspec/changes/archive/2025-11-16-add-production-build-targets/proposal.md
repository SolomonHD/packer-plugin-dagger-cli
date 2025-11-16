# Change: Add Production Build Targets

## Why
The current build system only supports development builds with `-dev` prerelease suffix via `make dev` and `make dev-install`. Production release verification and local production testing require clean builds without prerelease suffixes.

## What Changes
- Add `make prod` target that builds plugin with empty `VersionPrerelease` and runs describe test
- Add `make prod-install` target that builds and installs plugin without prerelease suffix
- Both targets follow the same versioning, git-derived path, and lowercase conventions as existing `dev` targets

## Impact
- Affected specs: `build-and-test-workflow`
- Affected code: [`GNUmakefile`](../../GNUmakefile:1)
- No breaking changes - purely additive functionality
- Enables production release testing and verification workflows