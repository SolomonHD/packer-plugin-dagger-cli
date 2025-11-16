# Change: Implement dagger-cli Provisioner Plugin

## Why
Packer users need a way to leverage Dagger's container-based pipeline capabilities during image provisioning. This enables:
- Reuse of existing Dagger modules for system configuration
- Container-based reproducibility
- Modern CI/CD tooling integration in image builds
- Simplified cross-platform provisioning workflows

## What Changes
- **NEW**: Complete `dagger-cli` provisioner implementation
- Replace scaffolding provisioner with production dagger-cli provisioner
- Add configuration schema with module/file sources, call entrypoint, args, env vars
- Implement retry logic with exponential backoff
- Add structured logging with configurable levels (trace/debug/info/warn/error)
- Support Dagger caching modes (auto/disabled/keyed)
- Comprehensive config validation in Prepare phase
- Unit tests with fake command runners
- Integration tests with Packer acceptance test framework
- Documentation with multiple usage examples

## Impact
**Affected specs:**
- `dagger-cli-provisioner` (NEW capability)

**Affected code:**
- `provisioner/scaffolding/` → replaced with `provisioner/dagger-cli/`
- `main.go` → update provisioner registration
- `docs/provisioners/` → add provisioner documentation
- `example/` → add example templates
- Test files throughout provisioner directory

**No breaking changes** - This is a new provisioner; existing code unchanged.