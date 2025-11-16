// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package daggercli

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

// Config holds the provisioner configuration
type Config struct {
	// Dagger module / source
	Module string `mapstructure:"module"`
	File   string `mapstructure:"file"`

	// Dagger call entrypoint
	Call string `mapstructure:"call" required:"true"`

	// Arguments and environment
	Args       map[string]interface{} `mapstructure:"args"`
	Env        map[string]string      `mapstructure:"env"`
	WorkingDir string                 `mapstructure:"working_dir"`

	// Observability
	LogLevel string `mapstructure:"log_level"`

	// Retries
	RetryCount   int    `mapstructure:"retry_count"`
	RetryBackoff string `mapstructure:"retry_backoff"`

	// Caching
	CacheMode string `mapstructure:"cache_mode"`
	CacheKey  string `mapstructure:"cache_key"`

	ctx interpolate.Context

	// Computed fields (not exposed in HCL)
	effectiveCacheKey string
	retryDuration     time.Duration
}

// Provisioner implements the Packer provisioner interface
type Provisioner struct {
	config Config
	runner CommandRunner
}

// CommandRunner abstracts command execution for testability
type CommandRunner interface {
	Run(ctx context.Context, cmd *exec.Cmd) error
}

// DefaultCommandRunner executes commands using os/exec
type DefaultCommandRunner struct{}

func (r *DefaultCommandRunner) Run(ctx context.Context, cmd *exec.Cmd) error {
	return cmd.Run()
}

// ConfigSpec returns the HCL2 spec for the config
func (p *Provisioner) ConfigSpec() hcldec.ObjectSpec {
	return p.config.FlatMapstructure().HCL2Spec()
}

// Prepare validates and normalizes the configuration
func (p *Provisioner) Prepare(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		PluginType:         "packer.provisioner.dagger-cli",
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return err
	}

	// Initialize default runner if not set (for testing)
	if p.runner == nil {
		p.runner = &DefaultCommandRunner{}
	}

	// Validation 1: module vs file mutual exclusivity
	if p.config.Module != "" && p.config.File != "" {
		return fmt.Errorf("dagger-cli: module and file are mutually exclusive; specify only one")
	}
	if p.config.Module == "" && p.config.File == "" {
		return fmt.Errorf("dagger-cli: one of module or file must be provided")
	}

	// Validation 2: call is required
	if p.config.Call == "" {
		return fmt.Errorf("dagger-cli: call must be provided")
	}

	// Validation 3: log_level normalization and validation
	if p.config.LogLevel == "" {
		p.config.LogLevel = "info"
	} else {
		p.config.LogLevel = strings.ToLower(p.config.LogLevel)
		validLevels := map[string]bool{
			"trace": true,
			"debug": true,
			"info":  true,
			"warn":  true,
			"error": true,
		}
		if !validLevels[p.config.LogLevel] {
			return fmt.Errorf("dagger-cli: invalid log_level '%s'. Allowed: trace, debug, info, warn, error", p.config.LogLevel)
		}
	}

	// Validation 4: retry_count and retry_backoff
	if p.config.RetryCount < 0 {
		return fmt.Errorf("dagger-cli: retry_count must be >= 0")
	}
	if p.config.RetryCount > 0 {
		if p.config.RetryBackoff == "" {
			p.config.RetryBackoff = "5s"
		}
		duration, err := time.ParseDuration(p.config.RetryBackoff)
		if err != nil {
			return fmt.Errorf("dagger-cli: invalid retry_backoff '%s': %v", p.config.RetryBackoff, err)
		}
		p.config.retryDuration = duration
	}

	// Validation 5: cache_mode and cache_key
	if p.config.CacheMode == "" {
		p.config.CacheMode = "auto"
	} else {
		p.config.CacheMode = strings.ToLower(p.config.CacheMode)
		validModes := map[string]bool{
			"auto":     true,
			"disabled": true,
			"keyed":    true,
		}
		if !validModes[p.config.CacheMode] {
			return fmt.Errorf("dagger-cli: invalid cache_mode '%s'. Allowed: auto, disabled, keyed", p.config.CacheMode)
		}
	}

	if p.config.CacheMode == "keyed" && p.config.CacheKey == "" {
		return fmt.Errorf("dagger-cli: cache_key is required when cache_mode is 'keyed'")
	}

	// Validation 6: compute effective cache key
	if p.config.CacheMode == "keyed" {
		moduleOrFile := p.config.Module
		if moduleOrFile == "" {
			moduleOrFile = p.config.File
		}

		// Serialize args to JSON for stable hashing
		argsJSON := ""
		if len(p.config.Args) > 0 {
			argsBytes, err := json.Marshal(p.config.Args)
			if err != nil {
				return fmt.Errorf("dagger-cli: failed to serialize args for cache key: %v", err)
			}
			argsJSON = string(argsBytes)
		}

		// Compute SHA256 hash of cache key components
		hashInput := p.config.CacheKey + moduleOrFile + p.config.Call + argsJSON
		hash := sha256.Sum256([]byte(hashInput))
		p.config.effectiveCacheKey = fmt.Sprintf("%x", hash)
	}

	return nil
}

// Provision executes the Dagger CLI command
func (p *Provisioner) Provision(ctx context.Context, ui packer.Ui, _ packer.Communicator, generatedData map[string]interface{}) error {
	// Log initial summary
	p.logSummary(ui)

	totalAttempts := 1 + p.config.RetryCount

	for attempt := 1; attempt <= totalAttempts; attempt++ {
		startTime := time.Now()

		// Log attempt start
		p.logAttempt(ui, attempt, totalAttempts, "starting")

		// Build command
		cmd, err := p.buildCommand(ctx)
		if err != nil {
			return fmt.Errorf("dagger-cli: failed to build command: %v", err)
		}

		// Execute command
		err = p.runner.Run(ctx, cmd)
		duration := time.Since(startTime)

		if err == nil {
			// Success
			p.logAttempt(ui, attempt, totalAttempts, fmt.Sprintf("finished in %.1fs (success)", duration.Seconds()))
			return nil
		}

		// Failed attempt
		exitCode := "unknown"
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = fmt.Sprintf("%d", exitErr.ExitCode())
		}

		p.logAttempt(ui, attempt, totalAttempts, fmt.Sprintf("finished in %.1fs (exit_code=%s)", duration.Seconds(), exitCode))

		// Check if we should retry
		if attempt < totalAttempts {
			// Calculate backoff delay (linear: attempt * base)
			delay := time.Duration(attempt) * p.config.retryDuration
			ui.Say(fmt.Sprintf("[dagger-cli] attempt %d/%d failed (exit_code=%s). Retrying in %s...",
				attempt, totalAttempts, exitCode, delay))
			time.Sleep(delay)
		} else {
			// Final attempt failed
			ui.Error(fmt.Sprintf("[dagger-cli] attempt %d/%d failed (exit_code=%s). No more retries.",
				attempt, totalAttempts, exitCode))
			return fmt.Errorf("dagger-cli: command failed after %d attempts: %v", totalAttempts, err)
		}
	}

	return fmt.Errorf("dagger-cli: unexpected end of retry loop")
}

// buildCommand constructs the Dagger CLI command
func (p *Provisioner) buildCommand(ctx context.Context) (*exec.Cmd, error) {
	// Build base command args
	args := []string{"call"}

	// Add module flag
	if p.config.Module != "" {
		args = append(args, "-m", p.config.Module)
	} else {
		// Derive module from file directory
		moduleDir := filepath.Dir(p.config.File)
		args = append(args, "-m", moduleDir)
	}

	// Add call name
	args = append(args, p.config.Call)

	// Convert args map to CLI flags
	for name, value := range p.config.Args {
		switch v := value.(type) {
		case bool:
			if v {
				// True boolean: add flag without value
				args = append(args, "--"+name)
			}
			// False boolean: omit entirely
		case string:
			args = append(args, "--"+name+"="+v)
		case int, int64, float64:
			args = append(args, "--"+name+"="+fmt.Sprintf("%v", v))
		default:
			return nil, fmt.Errorf("dagger-cli: unsupported type for arg '%s': %T", name, v)
		}
	}

	// Create command
	cmd := exec.CommandContext(ctx, "dagger", args...)

	// Set working directory
	if p.config.WorkingDir != "" {
		cmd.Dir = p.config.WorkingDir
	}

	// Set environment variables
	cmd.Env = os.Environ()
	for key, value := range p.config.Env {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	// Add cache mode environment variables
	switch p.config.CacheMode {
	case "disabled":
		cmd.Env = append(cmd.Env, "DAGGER_CACHE_DISABLED=1")
	case "keyed":
		cmd.Env = append(cmd.Env, "DAGGER_CACHE_KEY="+p.config.effectiveCacheKey)
	}

	// Connect stdout/stderr to UI
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, nil
}

// logSummary logs the initial configuration summary
func (p *Provisioner) logSummary(ui packer.Ui) {
	moduleOrFile := p.config.Module
	if moduleOrFile == "" {
		moduleOrFile = p.config.File
	}

	ui.Say(fmt.Sprintf("[dagger-cli] module=%s call=%s attempts=%d log_level=%s cache_mode=%s",
		moduleOrFile, p.config.Call, 1+p.config.RetryCount, p.config.LogLevel, p.config.CacheMode))

	if p.config.CacheMode == "keyed" {
		ui.Say(fmt.Sprintf("[dagger-cli] cache_key=%s effective_cache_key=%s",
			p.config.CacheKey, p.config.effectiveCacheKey))
	}
}

// logAttempt logs attempt status
func (p *Provisioner) logAttempt(ui packer.Ui, attempt, total int, status string) {
	ui.Say(fmt.Sprintf("[dagger-cli] attempt %d/%d %s", attempt, total, status))
}
