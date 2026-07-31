## 1. Canonical Plugin Identity and Build

- [ ] 1.1 Rename the Go module, imports, binary, Make targets, release metadata, and Packer source address from `packer-plugin-dagger-cli`/`dagger-cli` to `packer-plugin-dagger`/`dagger`.
- [ ] 1.2 Register the plugin under the Packer SDK default component name and update generated files for the renamed package paths.
- [ ] 1.3 Make every supported build path embed a valid semantic version and add tests or checks that reject malformed development binaries.

## 2. Retire the Host-Only Provisioner

- [ ] 2.1 Remove the `dagger-cli` provisioner implementation, registration, schema, generated HCL2 spec, and provisioner-specific unit tests.
- [ ] 2.2 Document how existing provisioner users can temporarily invoke Dagger through Packer's shell-local provisioner while migrating.

## 3. Post-Processor Configuration and Registration

- [ ] 3.1 Add and register the default `dagger` post-processor with required module and call settings plus arguments, secret arguments, artifact argument, working directory, report path, and Dagger command settings.
- [ ] 3.2 Validate mutually exclusive or conflicting settings and return actionable configuration diagnostics.
- [ ] 3.3 Regenerate the packer-sdc HCL2 spec after adding the post-processor Config and mapstructure tags.

## 4. Artifact Descriptor and Adapters

- [ ] 4.1 Implement the versioned `packer.dagger/artifact/v1` descriptor and deterministic JSON serialization.
- [ ] 4.2 Write descriptors to private temporary files, pass them as an explicit Dagger `File` argument, and clean them up on success, failure, or cancellation.
- [ ] 4.3 Implement the Amazon EBS artifact adapter with normalized region and AMI ID entries and unit tests for supported artifact encodings.
- [ ] 4.4 Reject unsupported artifact types with a message that identifies the builder ID and supported adapters.

## 5. Dagger Invocation and Secret Handling

- [ ] 5.1 Construct deterministic `dagger call` commands without shell interpolation and support local and remote module references.
- [ ] 5.2 Resolve declared secret inputs through supported environment references, pass them with Dagger secret semantics, and prevent secret values from appearing in logs or errors.
- [ ] 5.3 Stream useful Dagger output through the Packer UI, preserve cancellation, and translate nonzero exits into contextual post-processor errors.
- [ ] 5.4 Add unit tests for command construction, argument collisions, redaction, cancellation, and exit-code handling.

## 6. Artifact Ownership and Reports

- [ ] 6.1 Return the exact input artifact with keep-input-artifact behavior on successful Dagger validation.
- [ ] 6.2 Preserve the input artifact and report the failure when Dagger execution fails.
- [ ] 6.3 Support optional report export to an explicit host path without treating the report as a replacement Packer artifact.
- [ ] 6.4 Add tests for passthrough identity, keep flags, report export, and cleanup behavior.

## 7. Examples and Documentation

- [ ] 7.1 Replace scaffolded examples with a local Dagger fixture that validates a normalized artifact descriptor.
- [ ] 7.2 Add an Amazon EBS example showing Ansible Navigator provisioning followed by Dagger AMI validation.
- [ ] 7.3 Update the README and migration guide to explain the post-processor boundary, security model, supported artifacts, prerequisites, and repository rename.

## 8. Verification and Acceptance

- [ ] 8.1 Run formatting, generation, unit tests, race tests, vet, build, plugin checks, and strict OpenSpec validation.
- [ ] 8.2 Add a gated local Dagger/Packer acceptance test covering descriptor delivery and artifact passthrough.
- [ ] 8.3 Add a gated AWS acceptance test that builds an AMI, validates its metadata and launchability through Dagger, and reliably cleans up temporary resources.
