# Change: Remove Scaffolding Artifacts

## Why

The plugin was scaffolded from a Packer plugin template that included example implementations for all plugin types (builders, datasources, post-processors, provisioners). The project now focuses solely on the `dagger-cli` provisioner. All scaffolding code serves no purpose and creates maintenance burden and confusion about the plugin's actual capabilities.

## What Changes

- Remove scaffolding imports from [`main.go`](main.go:10-12,14)
- Remove scaffolding plugin registrations from [`main.go`](main.go:22-23,25-26)
- Delete entire scaffolding directories:
  - `builder/scaffolding/` (entire directory)
  - `datasource/scaffolding/` (entire directory)  
  - `post-processor/scaffolding/` (entire directory)
  - `provisioner/scaffolding/` (entire directory)
- Delete scaffolding documentation files:
  - `docs/builders/builder.mdx`
  - `docs/datasources/datasource.mdx`
  - `docs/post-processors/post-processor.mdx`
  - `docs/provisioners/provisioner.mdx` (NOT `dagger-cli.mdx`)

After cleanup, [`main.go`](main.go) will only:
- Import the `dagger-cli` provisioner
- Register the `dagger-cli` provisioner
- Import version and SDK packages

## Impact

- Affected files:
  - [`main.go`](main.go) (import cleanup, registration cleanup)
  - 4 top-level directories removed
  - 4 documentation files removed
  
- No behavioral changes to the `dagger-cli` provisioner
- Build and test systems remain unchanged
- Plugin functionality unchanged (only removes unused code)
- **BREAKING**: None (scaffolding code was never functional/documented for users)

## Validation

- `go build` completes without import errors
- `make dev` succeeds and produces working binary
- `make prod` succeeds
- `describe` output contains only `dagger-cli` in provisioners array
- No references to scaffolding remain in codebase