# OpenSpec change prompt

## Context

The packer-plugin-dagger-cli currently fails with "does not support Protobuf" error when users run `packer init`. This indicates the plugin is compiled with an older Packer plugin API version (x4 or earlier) that doesn't support modern gRPC/Protobuf communication.

Packer's latest plugin API version (x5) requires:
- packer-plugin-sdk >= 0.6.2
- gRPC/Protobuf support for plugin communication
- Proper API version declaration in plugin metadata

## Goal

Update the plugin to support Packer's x5 API version, eliminating the "does not support Protobuf" error and ensuring compatibility with modern Packer (>= 1.10.2).

## Scope

**In scope:**
- Update packer-plugin-sdk dependency to latest version supporting x5 API
- Update go.mod and run go mod tidy
- Rebuild and test plugin with x5 API version
- Verify `packer init` works without Protobuf errors
- Update any documentation referencing API versions

**Out of scope:**
- Changing plugin functionality or provisioner behavior
- Modifying existing configuration schema
- Adding new features or capabilities
- Changing version numbering in version/VERSION

## Desired behaviour

After the change:
- `packer init .` successfully downloads and initializes the plugin
- Plugin reports API version as "x5.0" in describe output
- No "does not support Protobuf" errors occur
- Plugin works with Packer >= 1.10.2
- All existing tests pass with new SDK version

## Constraints & assumptions

- Assume packer-plugin-sdk v0.6.2+ supports x5 API (verify latest version)
- Assume no breaking changes in SDK API between v0.6.1 and latest
- Plugin versioning still uses version/VERSION file (no change to versioning strategy)
- Binary naming and build process remain unchanged (follow packer-versioning-testing rules)
- All lowercase conventions for git paths and plugin names remain (follow packer-versioning-testing rules)

## Acceptance criteria

- [ ] go.mod updated with packer-plugin-sdk version supporting x5 API
- [ ] `go mod tidy` executed successfully
- [ ] `make dev` builds plugin without errors
- [ ] `packer-plugin-dagger-cli describe` shows `"api_version": "x5.0"`
- [ ] Test template with `required_plugins` block runs `packer init` without "Protobuf" errors
- [ ] `packer build` successfully executes with updated plugin
- [ ] All existing unit tests pass
- [ ] GNUmakefile build targets still work correctly