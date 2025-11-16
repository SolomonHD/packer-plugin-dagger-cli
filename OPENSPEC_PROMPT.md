# OpenSpec change prompt

## Context

This Packer plugin currently has versioning issues that make local testing difficult:

- The [`GNUmakefile`](GNUmakefile:15) hard-codes `VersionPrerelease=dev` in the `dev` target
- There's no clear way to build with the actual version from [`version/VERSION`](version/VERSION)
- The `describe` command (which shows plugin metadata) may fail or show incorrect version info
- The existing build process doesn't follow the project's versioning rules consistently

The plugin must follow these versioning rules (from `~/.kilocode/rules/packer-versioning-testing.md`):
- Single source of truth: `version/VERSION` file
- No hard-coded versions in Go code
- Derive git host/org from `.git`, then lowercase them
- Use environment variables for plugin binary names
- Run `describe` test after build, then move binary to `/tmp`

## Goal

Make it simple and obvious how to:
1. Build the plugin for local testing with proper versioning
2. Run the `describe` test to verify the build
3. Understand what version will be reported by the plugin

The solution should feel natural for someone familiar with Go projects or Makefiles.

## Scope

In scope:
- Create or modify build script(s) at the repo root to handle versioning correctly
- Ensure builds read from `version/VERSION` 
- Support both dev builds (with prerelease suffix) and clean builds (for releases)
- Include the `describe` test workflow
- Update `GNUmakefile` if needed to align with the new approach
- Provide clear usage instructions (inline comments or brief docs)

Out of scope:
- Modifying Go source files (version package is already set up correctly)
- Changing the release/CI process
- Altering the plugin's runtime behavior

## Desired behaviour

After implementation:

- Running `make dev` (or similar) should:
  - Read version from `version/VERSION`
  - Build with appropriate ldflags (including VersionPrerelease)
  - Use lowercase git host/org/plugin-name
  - Automatically run `./<binary> describe` to verify
  - Move the binary to `/tmp` after verification
  - Report success/failure clearly

- The `describe` output should show the correct version from `version/VERSION`

- It should be obvious from error messages if version/VERSION is missing or git remotes can't be parsed

## Constraints & assumptions

- Assume the `version/VERSION` file exists and contains a valid semantic version (e.g., `0.1.0`)
- Assume `.git` directory exists and has a remote configured
- The git remote URL can be parsed to extract host and organization
- Prefer bash/make for scripting (already used in this project)
- The solution should work on Unix-like systems (Linux, macOS)

## Acceptance criteria

- [ ] A build command exists that reads version from `version/VERSION` and builds correctly
- [ ] The built binary's `describe` output matches the version in `version/VERSION`
- [ ] After build+test, the binary is moved to `/tmp`, not left in the repo root
- [ ] Build failures or describe failures produce clear error messages
- [ ] The approach is documented (in the script itself or a config file comment)