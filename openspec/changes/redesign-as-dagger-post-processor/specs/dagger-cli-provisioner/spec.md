## REMOVED Requirements

### Requirement: Module Source Configuration
**Reason**: The host-only Dagger CLI provisioner is replaced by the artifact-aware Dagger post-processor.
**Migration**: Use `post-processor "dagger"` for artifact validation or Packer `shell-local` for an arbitrary host-side Dagger call.

### Requirement: Module and File Mutual Exclusivity
**Reason**: The `file` pseudo-module input belongs only to the removed provisioner and does not represent a current Dagger module contract.
**Migration**: Configure a Dagger module reference through the post-processor's `module` field.

### Requirement: Call Parameter Required
**Reason**: The provisioner configuration is removed.
**Migration**: Set `call` on `post-processor "dagger"`.

### Requirement: Arguments and Environment Variables
**Reason**: The provisioner argument and unrestricted environment overlay are removed in favor of post-processor ordinary and explicit secret arguments.
**Migration**: Move ordinary values to `args`, secret-provider references to `secret_args`, and host-only workflows to `shell-local`.

### Requirement: Working Directory Configuration
**Reason**: The provisioner configuration is removed.
**Migration**: Set `working_dir` on the Dagger post-processor when a local module requires it.

### Requirement: Log Level Configuration
**Reason**: The field was validated but never controlled plugin or Dagger output and created a false abstraction over Dagger logging.
**Migration**: Use deterministic default progress output; use Dagger directly for specialized debug sessions.

### Requirement: Retry Configuration
**Reason**: Retry and timeout policy is already supplied by Packer post-processor lifecycle settings and should not be duplicated.
**Migration**: Use Packer's common retry and timeout configuration or implement domain-specific retries in the Dagger function.

### Requirement: Cache Mode Configuration
**Reason**: The modes did not correspond to supported current Dagger cache controls.
**Migration**: Define caching on Dagger functions and operations using Dagger's native cache model.

### Requirement: Command Construction
**Reason**: Provisioner command construction is replaced by artifact-aware post-processor invocation.
**Migration**: Use the Dagger post-processor module, call, argument, secret, and artifact descriptor fields.

### Requirement: Retry Execution Logic
**Reason**: The provisioner retry loop is removed with the provisioner.
**Migration**: Use Packer lifecycle retries or Dagger-module logic.

### Requirement: Structured Logging
**Reason**: Provisioner-specific logging is removed and Dagger output will stream through the Packer UI from the post-processor.
**Migration**: Consume Packer output and optional Dagger reports from `post-processor "dagger"`.

### Requirement: Cache Mode Environment Variables
**Reason**: `DAGGER_CACHE_DISABLED` and `DAGGER_CACHE_KEY` are not a supported plugin-level contract for current Dagger caching.
**Migration**: Remove these settings and use Dagger-native function, layer, and volume caching.

### Requirement: Error Messages
**Reason**: The provisioner error contract is replaced by post-processor errors containing function, artifact, report, and exit-status context.
**Migration**: Update automation to match the new `dagger:` post-processor error prefix and fields.

### Requirement: Packer SDK Integration
**Reason**: The plugin no longer registers a Packer provisioner.
**Migration**: Use the new default Packer post-processor integration.
