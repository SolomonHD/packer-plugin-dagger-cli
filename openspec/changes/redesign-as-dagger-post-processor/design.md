## Context

The repository currently registers one `dagger-cli` provisioner that executes `dagger call` on the Packer host. It ignores both the communicator and `generatedData`, so it neither provisions the guest nor supplies meaningful Packer state to Dagger. Packer already provides `shell-local` for this behavior. The implementation also invents `DAGGER_CACHE_DISABLED` and `DAGGER_CACHE_KEY` controls, validates a log level it does not apply, and has no acceptance test against a real Dagger function.

The differentiated integration point is Packer's post-processing phase. A post-processor receives the completed builder artifact after Ansible Navigator and other provisioners have run. For `amazon-ebs`, the artifact is metadata identifying one or more regional AMIs, not an image file. A Dagger function can consume a normalized descriptor, use explicit AWS credentials and APIs to inspect or boot the AMI, emit a report, and fail the build as a release gate.

## Goals / Non-Goals

**Goals:**

- Rename the project to the durable `packer-plugin-dagger` identity and publish it from `github.com/SolomonHD/dagger`.
- Replace the redundant provisioner with one default `dagger` post-processor.
- Define a stable, versioned artifact descriptor independent of Packer's Go implementation types.
- Preserve the original artifact while allowing Dagger to validate it and export a report.
- Support generic Packer artifacts plus a tested `amazon-ebs` adapter for region-to-AMI metadata.
- Keep Dagger secrets explicit and prevent their values from entering command previews, descriptors, reports, or logs.
- Make cancellation, temporary-file cleanup, and external acceptance-resource cleanup reliable.

**Non-Goals:**

- Replace Ansible, Ansible Navigator, or Molecule for guest configuration and role testing.
- Implement AWS copy, sharing, Inspector, SSM, or promotion logic in Go; those belong in user-selected Dagger modules.
- Return replacement AMI artifacts or automatically destroy the input artifact in the first release.
- Embed or manage the Dagger engine, control Dagger's internal cache, or expose every Dagger global option.
- Preserve compatibility with `provisioner "dagger-cli"` or the unreleased HashiCorp namespace examples.

## Decisions

### Rename around Dagger, not its CLI transport

Rename the repository and Go module to `github.com/SolomonHD/packer-plugin-dagger`, publish the Packer source as `github.com/SolomonHD/dagger`, and build `packer-plugin-dagger`. Register the post-processor with `plugin.DEFAULT_NAME`, yielding `post-processor "dagger"` rather than a duplicated `dagger-cli-dagger-cli` component name.

The initial implementation still invokes the CLI because it already provides module discovery, typed argument conversion, engine version matching, secret-provider references, progress output, cancellation via process context, and result export. The project name remains correct if a later implementation adopts the Go SDK.

Alternative considered: name the plugin `image-gate`. Rejected because users explicitly configure Dagger modules and calls; hiding the required runtime would be misleading and would unnecessarily restrict non-AMI artifacts.

### Pass a versioned descriptor file as an explicit function argument

Before invoking Dagger, create a mode-0600 JSON file with this logical shape:

```json
{
  "schema": "packer.dagger/artifact/v1",
  "builder_id": "mitchellh.amazonebs",
  "id": "us-east-1:ami-0123456789abcdef0",
  "files": [],
  "display": "...",
  "attributes": {
    "aws": {
      "amis": [{"region": "us-east-1", "id": "ami-0123456789abcdef0"}]
    }
  }
}
```

Pass its path to a configurable function argument whose default name is `packer-artifact`. A Dagger module declares that argument as `File`, so Dagger imports it explicitly into the sandbox. The top-level keys are generic; builder-specific adapters populate allowlisted `attributes` only. Never attempt to serialize arbitrary `Artifact.State` values because the interface is not enumerable and state may contain credentials or unsupported Go values.

The descriptor schema is versioned from the first release. Additive fields may appear within v1, while incompatible changes require a new schema identifier.

Alternative considered: flatten artifact values into environment variables. Rejected because it loses structure, creates escaping problems, and makes secret leakage and schema evolution harder to control.

### Make validation passthrough the only initial artifact mode

On Dagger success, return the exact input `packer.Artifact` and force it to be retained. On failure or cancellation, return an error after cleanup. Do not manufacture an AWS artifact, call `Destroy`, or interpret a Dagger result as a replacement artifact.

This safely supports validation, attestation, notification, and user-authored promotion side effects while keeping Packer artifact ownership unambiguous. An optional `report_path` uses Dagger's output export to save a function result locally, but the report is not substituted for the AMI in a post-processor chain.

Alternative considered: return copied AMIs as a replacement artifact. Deferred because a correct artifact must implement builder-specific identity, state, display, and destruction semantics and must handle partial cross-region failure.

### Use a small explicit HCL surface

The post-processor configuration contains:

- required `module` and `call` strings;
- optional scalar `args` for ordinary Dagger function arguments;
- optional `secret_args` mapping argument names to Dagger secret-provider URIs, always redacted;
- optional `artifact_argument` defaulting to `packer-artifact`;
- optional `working_dir`;
- optional `report_path`;
- optional `dagger_command` defaulting to `dagger` for version-managed installations.

Always invoke without a shell and add deterministic non-interactive progress configuration. Reject collisions where `args` or `secret_args` attempts to supply the reserved artifact argument, reject overlapping ordinary and secret argument names, sort argument names for deterministic tests and logs, and validate report destinations before execution.

Do not expose retry, cache, or log-level settings. Packer already supplies post-processor retry/timeout controls, Dagger owns its cache, and Dagger's progress/debug behavior can evolve without a false plugin abstraction.

### Treat AWS behavior as an adapter and acceptance fixture

Recognize the `amazon-ebs` builder ID and parse its documented artifact ID/state into a deterministic `attributes.aws.amis` list. Keep the raw `id` even when parsing fails; fail with an actionable adapter error when an artifact claims to be Amazon EBS but its AMI mapping is malformed.

Unit tests cover descriptor generation from fixed artifact stubs. An opt-in acceptance fixture builds or accepts a disposable test AMI, invokes a minimal Dagger validation module, verifies that the module receives the expected region and AMI ID, and performs a non-destructive `DescribeImages` check. A separately gated launch-test case may create an instance and MUST track and clean every resource it creates.

### Separate repository rename from implementation compatibility

The GitHub repository rename occurs before implementation so new work lands under the final identity. GitHub redirects preserve the old clone URL temporarily, but code, documentation, release assets, module imports, and Packer source are updated together during apply. Because no release or tag exists, no compatibility shim is added.

## Risks / Trade-offs

- [A Dagger module can mutate or delete cloud resources] → Keep the plugin generic, document that module permissions define blast radius, recommend least-privilege validation roles, and make the built-in acceptance fixture non-destructive by default.
- [Descriptor fields could leak builder credentials] → Serialize only allowlisted fields, write mode-0600 temporary files, never include environment or arbitrary state, and test descriptors for known secret-shaped keys.
- [Packer may delete an input artifact after post-processing] → Return the original artifact and force retention in every successful path; unit-test the exact `keep` and `forceOverride` values.
- [Dagger report types vary] → Treat `report_path` as raw Dagger CLI output export and report CLI type errors without guessing conversions.
- [Amazon artifact encoding changes] → Preserve the raw artifact ID, isolate parsing in a tested adapter, and fail clearly rather than emitting an incorrect AMI map.
- [Renaming breaks unreleased local templates] → Provide a direct migration example from `dagger-cli` to `dagger` or `shell-local`; no published compatibility promise exists.

## Migration Plan

1. Rename the GitHub repository and re-clone it at `packer-plugin-dagger` after preserving the OpenSpec planning commit.
2. Update the Go module, imports, release metadata, Make targets, documentation, source address, binary name, and OpenSpec context.
3. Remove the provisioner implementation and register the new default post-processor.
4. Implement descriptor generation, Dagger invocation, report export, redaction, passthrough semantics, and unit tests.
5. Replace scaffolding examples with a minimal local Dagger fixture and an opt-in Amazon EBS validation example.
6. Run unit, race, vet, build, plugin-check, strict OpenSpec, local Dagger, and gated AWS acceptance checks.

Rollback before a release is a Git revert plus a GitHub repository rename back to the former name. Existing local users can invoke the same host-side Dagger call through Packer `shell-local` without the removed provisioner.

## Open Questions

- Confirm the exact Amazon EBS artifact encoding and builder-state keys against the pinned Amazon plugin during implementation, then lock them in fixtures.
- Decide whether the first release permits `cmd://` secret references; disabling them initially would reduce host-command execution risk without limiting environment, file, Vault, 1Password, or AWS secret providers.
