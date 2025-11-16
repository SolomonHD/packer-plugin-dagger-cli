// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package daggercli

import (
	"context"
	"io"
	"os/exec"
	"testing"
	"time"
)

// FakeCommandRunner is a test runner that simulates command execution
type FakeCommandRunner struct {
	FailCount    int // Number of times to fail before succeeding
	currentFails int // Current failure count
	Called       int // Number of times Run was called
	LastCmd      *exec.Cmd
	ShouldError  bool // Always fail if true
}

func (f *FakeCommandRunner) Run(ctx context.Context, cmd *exec.Cmd) error {
	f.Called++
	f.LastCmd = cmd

	if f.ShouldError {
		return &exec.ExitError{}
	}

	if f.currentFails < f.FailCount {
		f.currentFails++
		return &exec.ExitError{}
	}

	return nil
}

func TestConfigValidation_ModuleFileExclusive(t *testing.T) {
	tests := []struct {
		name        string
		module      string
		file        string
		call        string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "both module and file",
			module:      "github.com/test/mod",
			file:        "./test.py",
			call:        "build",
			shouldError: true,
			errorMsg:    "mutually exclusive",
		},
		{
			name:        "neither module nor file",
			module:      "",
			file:        "",
			call:        "build",
			shouldError: true,
			errorMsg:    "one of module or file must be provided",
		},
		{
			name:        "module only",
			module:      "github.com/test/mod",
			file:        "",
			call:        "build",
			shouldError: false,
		},
		{
			name:        "file only",
			module:      "",
			file:        "./test.py",
			call:        "build",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provisioner{
				runner: &FakeCommandRunner{},
			}

			config := map[string]interface{}{
				"call": tt.call,
			}
			if tt.module != "" {
				config["module"] = tt.module
			}
			if tt.file != "" {
				config["file"] = tt.file
			}

			err := p.Prepare(config)

			if tt.shouldError {
				if err == nil {
					t.Fatalf("expected error containing '%s', got nil", tt.errorMsg)
				}
				if !contains(err.Error(), tt.errorMsg) {
					t.Fatalf("expected error containing '%s', got: %s", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestConfigValidation_CallRequired(t *testing.T) {
	p := &Provisioner{
		runner: &FakeCommandRunner{},
	}

	config := map[string]interface{}{
		"module": "github.com/test/mod",
		// call is missing
	}

	err := p.Prepare(config)
	if err == nil {
		t.Fatal("expected error for missing call, got nil")
	}
	if !contains(err.Error(), "call must be provided") {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestConfigValidation_LogLevel(t *testing.T) {
	tests := []struct {
		name        string
		logLevel    string
		shouldError bool
		expected    string
	}{
		{
			name:        "missing log_level defaults to info",
			logLevel:    "",
			shouldError: false,
			expected:    "info",
		},
		{
			name:        "valid log_level trace",
			logLevel:    "trace",
			shouldError: false,
			expected:    "trace",
		},
		{
			name:        "valid log_level DEBUG (normalized)",
			logLevel:    "DEBUG",
			shouldError: false,
			expected:    "debug",
		},
		{
			name:        "invalid log_level",
			logLevel:    "invalid",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provisioner{
				runner: &FakeCommandRunner{},
			}

			config := map[string]interface{}{
				"module": "github.com/test/mod",
				"call":   "build",
			}
			if tt.logLevel != "" {
				config["log_level"] = tt.logLevel
			}

			err := p.Prepare(config)

			if tt.shouldError {
				if err == nil {
					t.Fatal("expected error for invalid log_level, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if p.config.LogLevel != tt.expected {
					t.Fatalf("expected log_level '%s', got '%s'", tt.expected, p.config.LogLevel)
				}
			}
		})
	}
}

func TestConfigValidation_RetryConfig(t *testing.T) {
	tests := []struct {
		name         string
		retryCount   int
		retryBackoff string
		shouldError  bool
		expectedDur  time.Duration
	}{
		{
			name:         "no retries",
			retryCount:   0,
			retryBackoff: "",
			shouldError:  false,
		},
		{
			name:         "retries with default backoff",
			retryCount:   3,
			retryBackoff: "",
			shouldError:  false,
			expectedDur:  5 * time.Second,
		},
		{
			name:         "retries with custom backoff",
			retryCount:   2,
			retryBackoff: "10s",
			shouldError:  false,
			expectedDur:  10 * time.Second,
		},
		{
			name:         "negative retry count",
			retryCount:   -1,
			retryBackoff: "",
			shouldError:  true,
		},
		{
			name:         "invalid backoff duration",
			retryCount:   1,
			retryBackoff: "invalid",
			shouldError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provisioner{
				runner: &FakeCommandRunner{},
			}

			config := map[string]interface{}{
				"module":      "github.com/test/mod",
				"call":        "build",
				"retry_count": tt.retryCount,
			}
			if tt.retryBackoff != "" {
				config["retry_backoff"] = tt.retryBackoff
			}

			err := p.Prepare(config)

			if tt.shouldError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.retryCount > 0 && p.config.retryDuration != tt.expectedDur {
					t.Fatalf("expected duration %v, got %v", tt.expectedDur, p.config.retryDuration)
				}
			}
		})
	}
}

func TestConfigValidation_CacheMode(t *testing.T) {
	tests := []struct {
		name         string
		cacheMode    string
		cacheKey     string
		shouldError  bool
		expectedMode string
	}{
		{
			name:         "missing cache_mode defaults to auto",
			cacheMode:    "",
			shouldError:  false,
			expectedMode: "auto",
		},
		{
			name:         "valid cache_mode disabled",
			cacheMode:    "disabled",
			shouldError:  false,
			expectedMode: "disabled",
		},
		{
			name:        "keyed mode requires cache_key",
			cacheMode:   "keyed",
			cacheKey:    "",
			shouldError: true,
		},
		{
			name:         "keyed mode with cache_key",
			cacheMode:    "keyed",
			cacheKey:     "my-cache",
			shouldError:  false,
			expectedMode: "keyed",
		},
		{
			name:        "invalid cache_mode",
			cacheMode:   "invalid",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provisioner{
				runner: &FakeCommandRunner{},
			}

			config := map[string]interface{}{
				"module": "github.com/test/mod",
				"call":   "build",
			}
			if tt.cacheMode != "" {
				config["cache_mode"] = tt.cacheMode
			}
			if tt.cacheKey != "" {
				config["cache_key"] = tt.cacheKey
			}

			err := p.Prepare(config)

			if tt.shouldError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if p.config.CacheMode != tt.expectedMode {
					t.Fatalf("expected cache_mode '%s', got '%s'", tt.expectedMode, p.config.CacheMode)
				}
				if tt.cacheMode == "keyed" && p.config.effectiveCacheKey == "" {
					t.Fatal("expected effectiveCacheKey to be computed, got empty string")
				}
			}
		})
	}
}

func TestCommandConstruction_Module(t *testing.T) {
	p := &Provisioner{
		runner: &FakeCommandRunner{},
	}

	config := map[string]interface{}{
		"module": "github.com/test/mod",
		"call":   "build",
	}

	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("unexpected prepare error: %v", err)
	}

	cmd, err := p.buildCommand(context.Background())
	if err != nil {
		t.Fatalf("unexpected buildCommand error: %v", err)
	}

	expectedArgs := []string{"call", "-m", "github.com/test/mod", "build"}
	if !sliceEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestCommandConstruction_File(t *testing.T) {
	p := &Provisioner{
		runner: &FakeCommandRunner{},
	}

	config := map[string]interface{}{
		"file": "ci/pipeline.py",
		"call": "build",
	}

	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("unexpected prepare error: %v", err)
	}

	cmd, err := p.buildCommand(context.Background())
	if err != nil {
		t.Fatalf("unexpected buildCommand error: %v", err)
	}

	// Should use directory of file as module
	expectedArgs := []string{"call", "-m", "ci", "build"}
	if !sliceEqual(cmd.Args[1:], expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args[1:])
	}
}

func TestCommandConstruction_Args(t *testing.T) {
	p := &Provisioner{
		runner: &FakeCommandRunner{},
	}

	config := map[string]interface{}{
		"module": "github.com/test/mod",
		"call":   "build",
		"args": map[string]interface{}{
			"region": "us-east-1",
			"fast":   true,
			"slow":   false,
			"count":  3,
		},
	}

	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("unexpected prepare error: %v", err)
	}

	cmd, err := p.buildCommand(context.Background())
	if err != nil {
		t.Fatalf("unexpected buildCommand error: %v", err)
	}

	args := cmd.Args[1:]

	// Check that expected flags are present (order may vary)
	expectedFlags := []string{"--region=us-east-1", "--fast", "--count=3"}
	for _, flag := range expectedFlags {
		if !sliceContains(args, flag) {
			t.Fatalf("expected flag '%s' in args %v", flag, args)
		}
	}

	// Check that false boolean is NOT present
	if sliceContains(args, "--slow") {
		t.Fatal("false boolean flag should not be present in args")
	}
}

func TestRetryLogic(t *testing.T) {
	t.Run("success on first try", func(t *testing.T) {
		fakeRunner := &FakeCommandRunner{}
		p := &Provisioner{
			runner: fakeRunner,
		}

		config := map[string]interface{}{
			"module": "github.com/test/mod",
			"call":   "build",
		}

		err := p.Prepare(config)
		if err != nil {
			t.Fatalf("unexpected prepare error: %v", err)
		}

		err = p.Provision(context.Background(), &FakeUI{}, nil, nil)
		if err != nil {
			t.Fatalf("unexpected provision error: %v", err)
		}

		if fakeRunner.Called != 1 {
			t.Fatalf("expected 1 call, got %d", fakeRunner.Called)
		}
	})

	t.Run("success on retry", func(t *testing.T) {
		fakeRunner := &FakeCommandRunner{
			FailCount: 2, // Fail twice, succeed on third
		}
		p := &Provisioner{
			runner: fakeRunner,
		}

		config := map[string]interface{}{
			"module":        "github.com/test/mod",
			"call":          "build",
			"retry_count":   3,
			"retry_backoff": "1ms", // Use short backoff for testing
		}

		err := p.Prepare(config)
		if err != nil {
			t.Fatalf("unexpected prepare error: %v", err)
		}

		err = p.Provision(context.Background(), &FakeUI{}, nil, nil)
		if err != nil {
			t.Fatalf("unexpected provision error: %v", err)
		}

		if fakeRunner.Called != 3 {
			t.Fatalf("expected 3 calls, got %d", fakeRunner.Called)
		}
	})

	t.Run("exhaust retries", func(t *testing.T) {
		fakeRunner := &FakeCommandRunner{
			ShouldError: true, // Always fail
		}
		p := &Provisioner{
			runner: fakeRunner,
		}

		config := map[string]interface{}{
			"module":        "github.com/test/mod",
			"call":          "build",
			"retry_count":   2,
			"retry_backoff": "1ms",
		}

		err := p.Prepare(config)
		if err != nil {
			t.Fatalf("unexpected prepare error: %v", err)
		}

		err = p.Provision(context.Background(), &FakeUI{}, nil, nil)
		if err == nil {
			t.Fatal("expected error after exhausting retries, got nil")
		}

		expectedCalls := 1 + 2 // initial + 2 retries
		if fakeRunner.Called != expectedCalls {
			t.Fatalf("expected %d calls, got %d", expectedCalls, fakeRunner.Called)
		}
	})
}

// Helper functions

type FakeUI struct{}

func (f *FakeUI) Say(message string)                                     {}
func (f *FakeUI) Message(message string)                                 {}
func (f *FakeUI) Error(message string)                                   {}
func (f *FakeUI) Machine(category string, args ...string)                {}
func (f *FakeUI) Ask(query string) (string, error)                       { return "", nil }
func (f *FakeUI) Askf(query string, args ...interface{}) (string, error) { return "", nil }
func (f *FakeUI) Sayf(message string, args ...interface{})               {}
func (f *FakeUI) Messagef(message string, args ...interface{})           {}
func (f *FakeUI) Errorf(message string, args ...interface{})             {}
func (f *FakeUI) TrackProgress(src string, currentNum, total int64, stream io.ReadCloser) io.ReadCloser {
	return stream
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
