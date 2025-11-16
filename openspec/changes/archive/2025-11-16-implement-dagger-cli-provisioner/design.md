# Design: Dagger CLI Provisioner

## Context
Implementing a Packer provisioner that executes Dagger CLI commands requires careful design of:
- Command execution abstraction for testability
- Retry logic with configurable backoff
- Cache key computation
- Integration with Packer's provisioner framework
- Error handling and observability

**Constraints:**
- Must work with Packer Plugin SDK v0.6.1
- Must not require Dagger daemon or state management
- Must integrate cleanly with Packer's UI system
- Must support both local and remote Dagger modules

**Stakeholders:**
- Packer users wanting Dagger integration
- CI/CD pipeline authors
- Infrastructure-as-code practitioners

## Goals / Non-Goals

**Goals:**
- Simple, testable provisioner implementation
- Clear configuration validation with helpful error messages
- Reliable retry mechanism with exponential backoff
- Observable execution with structured logging
- Support for Dagger's caching mechanisms

**Non-Goals:**
- Managing Dagger daemon lifecycle
- Dagger SDK integration (pure CLI approach)
- Custom Dagger module discovery or management
- Parallel execution of multiple Dagger calls
- Streaming logs from Dagger (rely on Dagger CLI's output)

## Decisions

### Decision 1: Command Execution Interface
**What:** Abstract command execution behind a `CommandRunner` interface for testability.

```go
type CommandRunner interface {
    Run(ctx context.Context, cmd *exec.Cmd) error
}
```

**Why:**
- Enables unit testing without real Dagger CLI
- Allows fake implementations that simulate failures for retry testing
- Standard Go pattern for external process execution

**Alternatives considered:**
- Direct `exec.Command` calls: harder to test
- Mock library approach: adds dependency and complexity

### Decision 2: Configuration Validation Strategy
**What:** Validate all configuration in `Prepare()` method before any execution.

**Why:**
- Fail fast principle - catch errors before provisioning starts
- Follows Packer SDK patterns
- Better user experience with clear validation errors

**Validation order:**
1. Module/file mutual exclusivity
2. Call parameter required check
3. Log level normalization and validation
4. Retry configuration validation
5. Cache mode and key validation
6. Effective cache key computation

### Decision 3: Retry Backoff Implementation
**What:** Implement simple linear backoff: `delay = attempt_number × base_duration`

**Why:**
- Simple to understand and predict
- Configurable via retry_backoff parameter
- Sufficient for transient failures (network, Dagger service)

**Example:**
```
Base: 5s, Retries: 3
Attempt 1: immediate
Attempt 2: wait 5s (1 × 5s)
Attempt 3: wait 10s (2 × 5s)
Attempt 4: wait 15s (3 × 5s)
```

**Alternatives considered:**
- Exponential backoff (2^n): too aggressive for typical provisioning
- Fixed backoff: less flexible
- Jittered backoff: unnecessary complexity

### Decision 4: Effective Cache Key Computation
**What:** Compute cache key from: `hash(cache_key + module/file + call + args_json)`

**Why:**
- Ensures different configurations don't share cache incorrectly
- Includes all inputs that affect Dagger execution
- Uses SHA256 for stable, collision-resistant hashing

**Implementation:**
```go
effectiveCacheKey := fmt.Sprintf("%x", sha256.Sum256([]byte(
    cacheKey + moduleOrFile + call + argsJSON
)))
```

### Decision 5: Args to CLI Flags Conversion
**What:** Convert HCL map args to CLI flags with type-specific handling:
- Boolean true: `--name` (flag only)
- Boolean false: omit
- String/Number: `--name=value`

**Why:**
- Matches Dagger CLI conventions
- Boolean flags are idiomatic in CLI tools
- Eliminates false flag clutter

**Example:**
```hcl
args = {
  region = "us-east-1"
  fast = true
  slow = false
  count = 3
}
```
Becomes: `--region=us-east-1 --fast --count=3`

### Decision 6: Logging Architecture
**What:** Log levels control verbosity; all logs prefixed with `[dagger-cli]`

**Log levels:**
- `trace`: Full command with args, env vars (redacted)
- `debug`: Detailed execution info, full context
- `info`: Summary logs (default)
- `warn`: Retry warnings, non-fatal issues
- `error`: Fatal failures only

**Why:**
- Standard log level hierarchy
- Prefix distinguishes provisioner logs from Packer and Dagger output
- Allows users to control verbosity based on debugging needs

### Decision 7: Cache Mode Environment Variables
**What:** Map cache_mode to Dagger environment variables:
- `auto`: no env vars (Dagger default behavior)
- `disabled`: `DAGGER_CACHE_DISABLED=1`
- `keyed`: `DAGGER_CACHE_KEY=<computed_key>`

**Why:**
- Environment variables are Dagger's standard cache control mechanism
- Declarative approach in HCL config
- No need to construct special CLI flags

**Note:** Verify actual Dagger CLI env var names during implementation.

### Decision 8: File-Based Module Resolution
**What:** When `file` is specified, derive module from file's directory: `-m path/to/dir`

**Why:**
- Dagger modules are directory-based
- Simplifies user config (don't need to specify dir when file gives context)
- Matches user mental model (file belongs to a module)

**Example:**
```hcl
file = "ci/pipeline.py"
```
Results in: `dagger call -m ci <call>`

### Decision 9: Working Directory Handling
**What:** Use Packer's default working directory unless `working_dir` explicitly set.

**Why:**
- Respects Packer's context by default
- Allows override for edge cases (absolute paths, different repo locations)
- Doesn't require users to think about working directories in common cases

### Decision 10: Error Message Format
**What:** All errors prefixed with `"dagger-cli: "` and provide actionable context.

**Examples:**
- `"dagger-cli: module and file are mutually exclusive; specify only one."`
- `"dagger-cli: call must be provided."`
- `"dagger-cli: invalid log_level 'foo'. Allowed: trace, debug, info, warn, error"`

**Why:**
- Clear attribution to this provisioner
- Actionable guidance for users
- Consistent error format

## Risks / Trade-offs

### Risk: Dagger CLI Not in PATH
**Impact:** Provisioner fails immediately when Dagger command not found.

**Mitigation:**
- Clear error message referencing Dagger CLI installation
- Document Dagger CLI as required dependency
- Consider adding Dagger version check in Prepare (future enhancement)

### Risk: Long-Running Dagger Operations
**Impact:** Provisioner may appear hung if Dagger takes a long time.

**Mitigation:**
- Log start time for each attempt
- Users can see Dagger's own output (it's not suppressed)
- Packer's own timeout mechanisms still apply

### Risk: Sensitive Data in Logs
**Impact:** Secrets in args or env vars may leak to logs.

**Mitigation:**
- Sanitize logs by default (don't show full args at info level)
- Only show full command at debug/trace levels
- Document that trace level may expose sensitive data

### Trade-off: CLI vs SDK Approach
**Decision:** Use CLI shelling instead of Dagger SDK

**Pros:**
- Simpler implementation (no Dagger SDK dependency)
- Works with any Dagger CLI version
- No state management required

**Cons:**
- Less control over execution
- Dependent on CLI output format for errors
- Additional process overhead

**Verdict:** CLI approach is correct for a Packer provisioner. SDK would require managing daemon lifecycle and state, which conflicts with Packer's stateless provisioner model.

## Migration Plan

**N/A** - This is a new provisioner; no migration needed.

## Open Questions

1. **Dagger CLI environment variable names:** Need to verify actual env var names for cache control (`DAGGER_CACHE_DISABLED`, `DAGGER_CACHE_KEY`) against Dagger CLI documentation.

2. **Args type handling:** How should we handle complex types (arrays, objects) in args map? Decision: Start with scalar types (string, number, bool) only. Document limitation and consider enhanced support in future versions.

3. **Dagger CLI version compatibility:** Should we enforce minimum Dagger CLI version? Decision: Document tested version but don't enforce in code initially. Add version check if compatibility issues arise.

4. **Communicator usage:** Should we use Packer's Communicator to transfer files? Decision: No - Dagger handles its own file context and mounts. Provisioner only executes CLI commands.

5. **Multiple calls:** Should we support multiple Dagger calls in one provisioner block? Decision: No - keep it simple. Users can add multiple provisioner blocks if needed.