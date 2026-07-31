## Why

The current `dagger-cli` provisioner is a host-side wrapper around `dagger call`: it ignores Packer's communicator and artifact data, overlaps with `shell-local`, and exposes cache controls that do not match current Dagger behavior. The useful integration point is after an image is built, where a Dagger function can validate, attest, or promote the actual Packer artifact without competing with Ansible for guest configuration.

## What Changes

- **BREAKING** Rename the repository, Go module, binary, plugin source, documentation identity, and default component from `packer-plugin-dagger-cli` / `dagger-cli` to `packer-plugin-dagger` / `dagger` under `github.com/SolomonHD`.
- **BREAKING** Remove the `dagger-cli` provisioner and its HCL surface, including the unsupported plugin-level cache and log-level abstractions.
- Add a default `dagger` post-processor that invokes a configured Dagger module function after builders and provisioners finish.
- Supply every Dagger function with a normalized, versioned JSON descriptor containing the Packer artifact ID, builder ID, files, string representation, and supported builder state.
- Preserve and return the input artifact by default; use the Dagger call as a cancellable validation/policy gate rather than silently replacing or destroying an AMI.
- Support explicit Dagger module inputs, typed secret references, deterministic non-interactive output, and an optional local report exported from the function result.
- Make AWS `amazon-ebs` AMI validation the first real acceptance path: Ansible Navigator configures the builder, Packer registers the AMI, and Dagger consumes its region-to-AMI metadata for boot, policy, or promotion workflows.
- Correct build/version injection so `make plugin-check` describes a semantic plugin version, and replace remaining scaffolding examples and documentation.

## Capabilities

### New Capabilities

- `dagger-artifact-post-processing`: Defines normalized artifact handoff, Dagger function execution, passthrough semantics, reports, secrets, cancellation, and AMI validation behavior.

### Modified Capabilities

- `dagger-cli-provisioner`: Remove the host-only provisioner and migrate users to the Dagger post-processor or Packer `shell-local`.
- `build-and-test-workflow`: Rename the plugin identity, correct semantic version injection, eliminate scaffolding examples, and add post-processor and AMI acceptance coverage.

## Impact

- Affected code includes `main.go`, `version/`, `GNUmakefile`, `go.mod`, generated HCL2 schemas, the existing provisioner package/tests, a new post-processor package/tests, examples, release metadata, and documentation.
- The GitHub repository becomes `SolomonHD/packer-plugin-dagger`; the Packer source becomes `github.com/SolomonHD/dagger` and the Go module becomes `github.com/SolomonHD/packer-plugin-dagger`.
- Existing `provisioner "dagger-cli"` blocks stop working. Arbitrary host-side calls migrate to `shell-local`; artifact validation migrates to `post-processor "dagger"`.
- Runtime prerequisites remain the Dagger CLI plus a supported Dagger engine runner. AWS acceptance tests additionally require explicitly enabled, least-privilege AWS credentials and must clean up all temporary resources.
