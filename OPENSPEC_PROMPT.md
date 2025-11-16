````md
You are working in a brand-new HashiCorp Packer plugin repository that was created from the official packer-plugin-template. Implement a new provisioner plugin named **"dagger-cli"** that shells out to the Dagger CLI.

## High-level goals

- Provide a `dagger-cli` provisioner that feels idiomatic for Dagger users.
- Primary concepts:
  - `module` → wired to `dagger call -m`
  - `call`   → the Dagger call/entrypoint name
- Support named arguments and environment variables as maps.
- Add first-class observability, retries, and caching controls.
- Keep the implementation simple, with clear validation and logs.

Use the latest Packer Plugin SDK for Go and follow the layout of the packer-plugin-template.

---

## Provisioner name and registration

- Plugin name: `dagger-cli`
- Provisioner type name (HCL block): `dagger-cli`
- Register the provisioner in `main.go` (or equivalent) using the SDK’s provisioner registration pattern.

---

## Configuration schema

Create a `Config` struct for the provisioner. Use these fields and mapstructure tags:

```go
type Config struct {
    // Dagger module / source
    Module      string            `mapstructure:"module"` // module reference for dagger -m
    File        string            `mapstructure:"file"`   // local file path / module dir alternative

    // Dagger call entrypoint
    Call        string            `mapstructure:"call"`

    // Arguments and environment
    Args        map[string]any    `mapstructure:"args"`        // named parameters for the Dagger call
    Env         map[string]string `mapstructure:"env"`         // extra env vars for Dagger CLI
    WorkingDir  string            `mapstructure:"working_dir"` // override process working dir

    // Observability
    LogLevel    string            `mapstructure:"log_level"`   // "trace", "debug", "info", "warn", "error"

    // Retries
    RetryCount   int    `mapstructure:"retry_count"`   // default 0
    RetryBackoff string `mapstructure:"retry_backoff"` // e.g. "5s", "30s", "2m"

    // Caching
    CacheMode string `mapstructure:"cache_mode"` // "auto", "disabled", "keyed"
    CacheKey  string `mapstructure:"cache_key"`
}
````

### Config semantics and validation (Prepare)

Implement `Prepare` to validate and normalize config:

1. **module vs file**

   * `module` and `file` are **mutually exclusive**.
   * If both are set, return an error:

     * `"dagger-cli: module and file are mutually exclusive; specify only one."`
   * It is valid to set **either** `module` **or** `file`.
   * If neither is set, error:

     * `"dagger-cli: one of module or file must be provided."`

2. **call**

   * `Call` is **required**.
   * If empty, return:

     * `"dagger-cli: call must be provided."`

3. **log_level**

   * Normalize to lowercase.
   * Allowed: `"trace"`, `"debug"`, `"info"`, `"warn"`, `"error"`.
   * Default: `"info"` if empty.
   * If invalid, return a validation error.

4. **retry_count / retry_backoff**

   * `RetryCount`:

     * Default 0 (no retries).
     * Must be >= 0; if negative, error.
   * `RetryBackoff`:

     * If `RetryCount > 0` and `RetryBackoff` is empty, default to `"5s"`.
     * Parse `RetryBackoff` with `time.ParseDuration`; if parsing fails, return an error.

5. **cache_mode / cache_key**

   * `CacheMode` default: `"auto"` if empty.
   * Allowed values: `"auto"`, `"disabled"`, `"keyed"`. Otherwise error.
   * When `CacheMode == "keyed"`:

     * Require non-empty `CacheKey`; otherwise error.
     * Precompute an `effectiveCacheKey` string that combines:

       * `CacheKey`
       * `module` or `file`
       * `Call`
       * Optionally a hash of `Args`.
     * You can keep `effectiveCacheKey` as a field on the runtime/provisioner struct.

6. **Args & Env**

   * `Args` is a free-form map; keep it as-is and handle during command building.
   * `Env` is just extra env vars; no special validation needed besides nil handling.

Return any validation errors through the usual Packer `Prepare` mechanism.

---

## Provisioner behavior (Provision)

Implement the core behavior in `Provision`:

1. **Resolve working directory**

   * If `WorkingDir` is set, run the Dagger CLI process using that directory as the working directory.
   * Otherwise, use Packer’s default working directory (the one provided by the SDK).

2. **Compute effective module / source**

   * If `Module` is set: use it directly as the `-m` value.
   * Else: compute a module/source based on `file`:

     * For a simple implementation, you can:

       * Use the directory of `File` as the module (`-m dirOf(file)`).
       * Optionally pass the filename as an argument if your Dagger program needs it.

3. **Command construction**

Build a `[]string` for the Dagger CLI command equivalent to:

```bash
dagger call -m <MODULE> <CALL> [ARGS...]
```

In Go:

```go
cmdArgs := []string{"call", "-m", effectiveModule, effectiveCall}
```

Then:

* Append args from `Args`:

  * For string/number/bool types, convert them into flags:

    * Booleans:

      * If `true`, use `--name` (no value).
      * If `false`, either omit or use `--name=false`; pick one and be consistent.
    * Strings/numbers:

      * Use `--name=value` or `--name`, `value`.
* Env:

  * Start from the existing process environment.
  * Overlay config `Env` onto it.
  * Add any extra env vars needed for logging/caching (e.g., `DAGGER_*`).

4. **Observability / logging**

* Implement a small logger helper that:

  * Normalizes `LogLevel`.
  * Prefixes all messages with `[dagger-cli]`.

* Before the first attempt, log a summary:

  ```text
  [dagger-cli] module=<...> file=<...> call=<...> attempts=<totalAttempts> log_level=<...> cache_mode=<...> cache_key=<...>
  ```

* Before each attempt:

  * Log the attempt number and a **sanitized** command preview:

    ```text
    [dagger-cli] attempt 1/<N> running: dagger call -m <module> <call> [args...]
    ```

  * Log the start time.

* After each attempt:

  * Log end time and duration:

    ```text
    [dagger-cli] attempt 1/<N> finished in 42.3s (exit_code=1)
    ```

* Respect `LogLevel`:

  * `debug`/`trace`: detailed logs, including redacted args/env.
  * `info`: summary only.
  * `warn`/`error`: mainly failures.

5. **Retries**

* Total attempts = `1 + RetryCount`.
* Wrap the Dagger CLI call:

```go
for attempt := 1; attempt <= totalAttempts; attempt++ {
    // run command
    // if success, return nil
    // if failure and attempt < totalAttempts, sleep with backoff
}
```

* Only retry on **non-zero exit codes** from the Dagger command.

* For backoff:

  * Parse `RetryBackoff` into `base`.
  * Use a simple strategy, e.g.:

    ```go
    delay := time.Duration(attempt) * base
    ```

* Log:

  * On non-final failed attempt:

    ```text
    [dagger-cli] attempt 1/<N> failed (exit_code=1). Retrying in 5s…
    ```

  * On final failed attempt:

    ```text
    [dagger-cli] attempt 3/3 failed (exit_code=1). No more retries.
    ```

6. **Caching**

Use `CacheMode` and `effectiveCacheKey` to influence Dagger via env vars:

* `auto`: default Dagger caching; no special env.
* `disabled`: set something like `DAGGER_CACHE_DISABLED=1`.
* `keyed`: set something like `DAGGER_CACHE_KEY=<effectiveCacheKey>`.

Always log cache mode at the start:

```text
[dagger-cli] cache_mode=keyed cache_key=<CacheKey> effective_cache_key=<hash>
```

The exact env var names can be adjusted to match real Dagger CLI conventions; keep them encapsulated in a helper.

---

## Testing

Add unit tests for:

* Config validation in `Prepare`:

  * `module` xor `file`
  * `call` required
  * invalid `log_level`
  * retry / backoff parsing
  * cache mode and key rules
* Retry behavior:

  * Use a fake command runner that fails a configurable number of times.
* Command construction:

  * Confirm that `module` + `call` → `["call", "-m", module, call, ...]`.
  * Confirm args encode correctly.

Abstract the command runner behind an interface so tests do not require the real `dagger` binary.

---

## Documentation

Create a README/docs section for the provisioner with examples:

1. **Basic usage (module + call)**

```hcl
provisioner "dagger-cli" {
  module = "github.com/acme/infra//ci?ref=v1.2.3"
  call   = "configure_system"

  args = {
    region = "us-east-1"
    fast   = true
    count  = 3
  }

  env = {
    ENVIRONMENT = "staging"
  }
}
```

2. **Local module example**

```hcl
provisioner "dagger-cli" {
  module = "./ci"
  call   = "build_image"

  args = {
    region = "us-west-2"
  }
}
```

3. **File-based example**

```hcl
provisioner "dagger-cli" {
  file = "ci/pipeline.py"
  call = "configure_system"

  args = {
    region = "us-east-1"
  }
}
```

4. **Observability, retries, caching example**

```hcl
provisioner "dagger-cli" {
  module = "github.com/acme/infra//ci?ref=v1.2.3"
  call   = "configure_system"

  log_level     = "debug"
  retry_count   = 3
  retry_backoff = "10s"

  cache_mode = "keyed"
  cache_key  = "ubuntu-22.04-base"
}
```

Make sure the implementation, tests, and docs all match these semantics. Use **`call`** and **`module`** consistently as the primary Dagger concepts.

```
```
