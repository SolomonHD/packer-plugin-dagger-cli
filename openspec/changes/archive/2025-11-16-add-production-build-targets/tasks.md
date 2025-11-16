# Implementation Tasks

## 1. Add Production Build Target
- [x] 1.1 Create `prod` target in [`GNUmakefile`](../../GNUmakefile:1)
- [x] 1.2 Set `VersionPrerelease` to empty string in ldflags
- [x] 1.3 Include describe test validation (like `dev` target)
- [x] 1.4 Move binary to `/tmp` after test completion
- [x] 1.5 Verify version output matches `version/VERSION` exactly (no suffix)

## 2. Add Production Install Target
- [x] 2.1 Create `prod-install` target in [`GNUmakefile`](../../GNUmakefile:1)
- [x] 2.2 Build with empty `VersionPrerelease` (like `dev-install`)
- [x] 2.3 Skip describe test (like `dev-install`)
- [x] 2.4 Install directly via `packer plugins install`

## 3. Testing & Validation
- [x] 3.1 Test `make prod` builds successfully
- [x] 3.2 Verify `describe` output shows version without `-dev` suffix
- [x] 3.3 Confirm binary moved to `/tmp` after test
- [x] 3.4 Test `make prod-install` installs successfully
- [x] 3.5 Verify installed plugin reports correct version

## 4. Documentation
- [x] 4.1 Add inline comments for `prod` target
- [x] 4.2 Add inline comments for `prod-install` target
- [x] 4.3 Update any README sections referencing build targets