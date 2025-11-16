# Implementation Tasks

## 1. Clean up main.go imports
- [x] 1.1 Remove import for `builder/scaffolding`
- [x] 1.2 Remove import for `datasource/scaffolding` (aliased as `scaffoldingData`)
- [x] 1.3 Remove import for `post-processor/scaffolding` (aliased as `scaffoldingPP`)
- [x] 1.4 Remove import for `provisioner/scaffolding` (aliased as `scaffoldingProv`)
- [x] 1.5 Verify `dagger-cli` provisioner import remains
- [x] 1.6 Verify version and SDK imports remain

## 2. Clean up main.go registrations
- [x] 2.1 Remove `pps.RegisterBuilder("my-builder", new(scaffolding.Builder))` line
- [x] 2.2 Remove `pps.RegisterProvisioner("my-provisioner", new(scaffoldingProv.Provisioner))` line
- [x] 2.3 Remove `pps.RegisterPostProcessor("my-post-processor", new(scaffoldingPP.PostProcessor))` line
- [x] 2.4 Remove `pps.RegisterDatasource("my-datasource", new(scaffoldingData.Datasource))` line
- [x] 2.5 Verify `dagger-cli` provisioner registration remains
- [x] 2.6 Verify SetVersion and Run calls remain

## 3. Delete scaffolding directories
- [x] 3.1 Delete `builder/scaffolding/` directory and all contents
- [x] 3.2 Delete `datasource/scaffolding/` directory and all contents
- [x] 3.3 Delete `post-processor/scaffolding/` directory and all contents
- [x] 3.4 Delete `provisioner/scaffolding/` directory and all contents
- [x] 3.5 Verify `provisioner/dagger-cli/` directory remains untouched

## 4. Delete scaffolding documentation
- [x] 4.1 Delete `docs/builders/builder.mdx`
- [x] 4.2 Delete `docs/datasources/datasource.mdx`
- [x] 4.3 Delete `docs/post-processors/post-processor.mdx`
- [x] 4.4 Delete `docs/provisioners/provisioner.mdx` (generic scaffolding doc)
- [x] 4.5 Verify `docs/provisioners/dagger-cli.mdx` remains

## 5. Validation and testing
- [x] 5.1 Run `go build` - must succeed without import errors
- [x] 5.2 Run `make dev` - must build and test successfully
- [x] 5.3 Run binary `describe` - verify output contains only `dagger-cli` in provisioners array
- [x] 5.4 Search codebase for "scaffolding" references - should only find historical references in git
- [x] 5.5 Verify no broken imports or references remain