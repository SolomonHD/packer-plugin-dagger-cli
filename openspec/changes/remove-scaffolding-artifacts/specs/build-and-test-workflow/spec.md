# build-and-test-workflow Spec Delta

## ADDED Requirements

### Requirement: Plugin Describe Output Validation
The build system SHALL validate that the plugin's describe output contains only intended components and no scaffolding artifacts.

#### Scenario: Describe output contains only dagger-cli provisioner
- **WHEN** the `describe` command executes successfully
- **THEN** the JSON output's `provisioners` array SHALL contain only `["dagger-cli"]`
- **AND** the `builders` array SHALL be empty or contain no scaffolding entries
- **AND** the `post_processors` array SHALL be empty or contain no scaffolding entries  
- **AND** the `datasources` array SHALL be empty or contain no scaffolding entries

#### Scenario: No scaffolding components in describe output
- **WHEN** validating the describe JSON output
- **THEN** it SHALL NOT contain any references to:
  - "my-builder"
  - "my-provisioner"  
  - "my-post-processor"
  - "my-datasource"
- **AND** the only provisioner listed SHALL be "dagger-cli"

#### Scenario: Describe validation as quality gate
- **WHEN** describe test runs as part of `make dev` or `make prod`
- **THEN** it SHALL verify provisioners array matches expected output
- **AND** warn or fail if unexpected components are present
- **AND** confirm single-purpose plugin architecture