# OpenSpec change prompt

## Context

The [`GNUmakefile`](GNUmakefile) currently provides `make dev` and `make dev-install` targets that build the plugin with a `-dev` prerelease suffix appended to the version.

## Goal

Add `make prod` and `make prod-install` targets that build production-ready binaries without the `-dev` suffix.

## Scope

- In scope:
  - Create `make prod` target similar to `make dev`, but with empty `VersionPrerelease`
  - Create `make prod-install` target similar to `make dev-install`, but with empty `VersionPrerelease`
  - Both targets must follow the same versioning rules as `dev`/`dev-install`
  - `make prod` should run the describe test like `make dev` does
  - `make prod-install` should skip the describe test like `make dev-install` does
  
- Out of scope:
  - Changes to `make dev` or `make dev-install`
  - Changes to version/VERSION file location or format
  - Changes to other Makefile targets

## Desired behaviour

After the change:

- `make prod` builds the plugin with version from `version/VERSION` (no `-dev` suffix)
- `make prod` runs the describe test and validates version reporting
- `make prod` moves binary to `/tmp/${PACKER_PLUGIN_BIN}` after test
- `make prod-install` builds and directly installs the plugin without `-dev` suffix
- Both targets use the same git-derived paths and lowercase conventions as `dev` targets

## Constraints & assumptions

- Assumption: The version in `version/VERSION` is the canonical production version
- Assumption: Production builds should set `VersionPrerelease=""` (empty string) in ldflags
- Constraint: Must follow the packer-versioning-testing rule from `.cursorrules`
- Constraint: Must maintain consistency with existing `dev` target structure

## Acceptance criteria

- [x] `make prod` target exists in [`GNUmakefile`](GNUmakefile)
- [x] `make prod-install` target exists in [`GNUmakefile`](GNUmakefile)
- [x] `make prod` builds with `-X '${PLUGIN_FQN}/version.VersionPrerelease='` (empty string)
- [x] `make prod` runs describe test and validates version matches `version/VERSION` exactly
- [x] `make prod-install` builds with empty VersionPrerelease and installs via `packer plugins install`
- [x] Both targets follow the same patterns as `dev`/`dev-install` (version sourcing, git paths, lowercase)