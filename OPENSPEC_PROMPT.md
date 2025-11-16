# OpenSpec change prompt

## Context

The packer-plugin-dagger-cli project was scaffolded from a Packer plugin template that included example implementations for builders, datasources, post-processors, and provisioners. The plugin is now focused solely on the `dagger-cli` provisioner, making all scaffolding artifacts unnecessary.

## Goal

Remove all scaffolding-related code, directories, and documentation from the project, leaving only the functional `dagger-cli` provisioner and its supporting infrastructure.

## Scope

**In scope:**
- Remove scaffolding imports from [`main.go`](packer/plugins/packer-plugin-dagger-cli/main.go)
- Remove scaffolding plugin registrations from [`main.go`](packer/plugins/packer-plugin-dagger-cli/main.go)
- Delete entire scaffolding directories:
  - `builder/scaffolding/`
  - `datasource/scaffolding/`
  - `post-processor/scaffolding/`
  - `provisioner/scaffolding/`
- Delete scaffolding documentation files:
  - `docs/builders/builder.mdx`
  - `docs/datasources/datasource.mdx`
  - `docs/post-processors/post-processor.mdx`
  - `docs/provisioners/provisioner.mdx` (NOT `dagger-cli.mdx`)
- Update any references to removed components

**Out of scope:**
- The [`provisioner/dagger-cli/`](packer/plugins/packer-plugin-dagger-cli/provisioner/dagger-cli/) directory and all its contents (MUST be preserved)
- [`docs/provisioners/dagger-cli.mdx`](packer/plugins/packer-plugin-dagger-cli/docs/provisioners/dagger-cli.mdx) (MUST be preserved)
- Build configuration ([`GNUmakefile`](packer/plugins/packer-plugin-dagger-cli/GNUmakefile), [`go.mod`](packer/plugins/packer-plugin-dagger-cli/go.mod), etc.)
- Version management ([`version/`](packer/plugins/packer-plugin-dagger-cli/version/) directory)
- OpenSpec documentation
- [`README.md`](packer/plugins/packer-plugin-dagger-cli/README.md)
- License and repo metadata

## Desired behaviour

After this change:

- [`main.go`](packer/plugins/packer-plugin-dagger-cli/main.go) should only:
  - Import the `dagger-cli` provisioner
  - Register the `dagger-cli` provisioner
  - Import version and plugin SDK packages
- No scaffolding directories exist in the project
- No scaffolding documentation exists in `docs/`
- `go build` completes successfully (no import errors)
- `make dev` target runs successfully
- The plugin's `describe` output shows only the `dagger-cli` provisioner

## Constraints & assumptions

- Assume scaffolding code has no dependencies outside its own directories
- Assume all scaffolding components follow standard naming: "scaffolding", "my-builder", "my-provisioner", "my-post-processor", "my-datasource"
- The plugin should remain buildable and testable after cleanup
- No functional changes to the `dagger-cli` provisioner

## Acceptance criteria

- [ ] [`main.go`](packer/plugins/packer-plugin-dagger-cli/main.go) contains no scaffolding imports
- [ ] [`main.go`](packer/plugins/packer-plugin-dagger-cli/main.go) registers only the `dagger-cli` provisioner
- [ ] All scaffolding directories deleted: `builder/`, `datasource/`, `post-processor/`, `provisioner/scaffolding/`
- [ ] All scaffolding docs deleted except [`docs/provisioners/dagger-cli.mdx`](packer/plugins/packer-plugin-dagger-cli/docs/provisioners/dagger-cli.mdx)
- [ ] `go build` succeeds without errors
- [ ] `make dev` succeeds and produces a working plugin binary
- [ ] Plugin `describe` output contains only `dagger-cli` in provisioners array