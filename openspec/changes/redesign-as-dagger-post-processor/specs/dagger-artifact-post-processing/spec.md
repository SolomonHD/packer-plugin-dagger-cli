## ADDED Requirements

### Requirement: Default Dagger Post-Processor

The plugin MUST register a default Packer post-processor addressable as `post-processor "dagger"` and MUST NOT require a component-name suffix.

#### Scenario: Plugin describe output

- **WHEN** the renamed plugin's `describe` command runs
- **THEN** it MUST advertise one default post-processor
- **AND** it MUST advertise no provisioners, builders, or data sources

### Requirement: Dagger Module Invocation

The post-processor SHALL invoke one configured Dagger module function without a shell and SHALL provide deterministic non-interactive execution.

#### Scenario: Valid module call

- **GIVEN** `module`, `call`, and ordinary scalar arguments are configured
- **WHEN** post-processing begins
- **THEN** the plugin SHALL invoke the configured `dagger` command with `call`, the module reference, function name, and arguments
- **AND** argument names SHALL be ordered deterministically
- **AND** Dagger output SHALL stream through the Packer UI

#### Scenario: Required configuration is absent

- **GIVEN** `module` or `call` is empty
- **WHEN** configuration is prepared
- **THEN** the plugin MUST return an actionable validation error before post-processing

#### Scenario: Version-managed Dagger command

- **GIVEN** `dagger_command` is an explicit executable path
- **WHEN** the function is invoked
- **THEN** the plugin SHALL execute that path without shell expansion

### Requirement: Versioned Packer Artifact Descriptor

The post-processor MUST provide the Dagger function a mode-0600 JSON descriptor representing the input Packer artifact through an explicit Dagger `File` argument.

#### Scenario: Generic artifact descriptor

- **GIVEN** any Packer artifact
- **WHEN** its descriptor is generated
- **THEN** the descriptor MUST include schema identifier `packer.dagger/artifact/v1`, builder ID, raw artifact ID, files, and display text
- **AND** the descriptor MUST NOT contain process environment, credentials, or arbitrary builder state
- **AND** it SHALL be passed through the argument named by `artifact_argument`, defaulting to `packer-artifact`

#### Scenario: Reserved argument collision

- **GIVEN** an ordinary or secret argument uses the configured artifact argument name
- **WHEN** configuration is validated
- **THEN** validation MUST fail before Dagger executes

#### Scenario: Descriptor cleanup

- **WHEN** Dagger completes, fails, or is cancelled
- **THEN** the temporary descriptor and its temporary directory MUST be removed

### Requirement: Amazon EBS AMI Descriptor Adapter

The post-processor MUST normalize an `amazon-ebs` artifact into a deterministic region-to-AMI representation while retaining the raw Packer artifact fields.

#### Scenario: Single-region AMI

- **GIVEN** an Amazon EBS artifact identifying one regional AMI
- **WHEN** its descriptor is generated
- **THEN** `attributes.aws.amis` MUST contain that region and AMI ID

#### Scenario: Multi-region AMIs

- **GIVEN** an Amazon EBS artifact identifying AMIs in multiple regions
- **WHEN** its descriptor is generated
- **THEN** `attributes.aws.amis` MUST contain every region and AMI ID sorted by region

#### Scenario: Malformed Amazon artifact metadata

- **GIVEN** an artifact claims the Amazon EBS builder ID but its AMI mapping is malformed
- **WHEN** its descriptor is generated
- **THEN** post-processing MUST fail with the raw artifact ID and an actionable adapter error
- **AND** Dagger MUST NOT execute with an incorrect or partial mapping

### Requirement: Passthrough Artifact Ownership

The post-processor MUST act as a validation gate and preserve the exact input Packer artifact on success.

#### Scenario: Dagger succeeds

- **WHEN** the Dagger function exits successfully
- **THEN** the post-processor MUST return the original artifact instance
- **AND** it MUST force Packer to retain the input artifact
- **AND** it MUST NOT call the artifact's `Destroy` method

#### Scenario: Dagger fails

- **WHEN** the Dagger function exits unsuccessfully
- **THEN** the post-processor MUST return an error containing the actual exit status and function identity
- **AND** it MUST NOT manufacture or return a replacement artifact

### Requirement: Optional Dagger Report Export

The post-processor SHALL optionally export the Dagger function result to an explicit local report path without replacing the Packer artifact.

#### Scenario: Report requested

- **GIVEN** `report_path` is configured
- **WHEN** Dagger succeeds with an exportable result
- **THEN** the result SHALL be written to the requested path
- **AND** the original Packer artifact SHALL remain the post-processor result

#### Scenario: Report export fails

- **WHEN** Dagger cannot export its result to `report_path`
- **THEN** post-processing MUST fail with an error identifying the report path
- **AND** any partial temporary output managed by the plugin MUST be cleaned up

### Requirement: Explicit Secret Arguments

The post-processor MUST distinguish Dagger secret-provider references from ordinary arguments and prevent secret values from entering plugin output.

#### Scenario: Secret reference passed

- **GIVEN** `secret_args` maps a function argument to a supported Dagger secret-provider URI
- **WHEN** Dagger is invoked
- **THEN** the URI SHALL be passed as the named function argument
- **AND** command previews and errors MUST redact its value

#### Scenario: Argument name overlap

- **GIVEN** the same argument name appears in `args` and `secret_args`
- **WHEN** configuration is validated
- **THEN** validation MUST fail before Dagger executes

### Requirement: Cancellation and Cleanup

The post-processor MUST propagate Packer cancellation to Dagger and clean all plugin-managed local resources on every exit path.

#### Scenario: Packer cancels post-processing

- **WHEN** the Packer context is cancelled during Dagger execution
- **THEN** the Dagger process MUST be terminated through its context
- **AND** the post-processor MUST return the context error
- **AND** local descriptor and report temporaries MUST be removed

### Requirement: AWS AMI Validation Acceptance Path

The project MUST provide an opt-in acceptance path proving that the Dagger post-processor receives and validates a real Amazon EBS artifact without leaking credentials or leaving temporary resources.

#### Scenario: Non-destructive AMI metadata validation

- **GIVEN** AWS acceptance testing is explicitly enabled with least-privilege credentials and a test AMI
- **WHEN** the fixture invokes a Dagger validation module
- **THEN** the module MUST receive the correct source region and AMI ID
- **AND** it MUST successfully perform a non-destructive image metadata check
- **AND** the post-processor MUST return the original artifact

#### Scenario: Launch validation cleanup

- **GIVEN** the separately gated launch-test fixture creates an EC2 instance or related resource
- **WHEN** the test succeeds, fails, times out, or is cancelled
- **THEN** every resource created by the fixture MUST be terminated or deleted
- **AND** cleanup failures MUST identify the surviving resource IDs
