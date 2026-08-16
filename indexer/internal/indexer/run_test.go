package indexer

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestConfigValidate(t *testing.T) {
	deploymentBlock := uint64(1)
	valid := Config{
		RPC:             "http://localhost:8545",
		Contract:        "0x0000000000000000000000000000000000000001",
		DeploymentBlock: &deploymentBlock,
		Interval:        time.Second,
		BatchSize:       100,
	}
	tests := []struct {
		name        string
		mutate      func(*Config)
		wantErrPart string
	}{
		{name: "valid"},
		{name: "invalid contract", mutate: func(c *Config) { c.Contract = "nope" }, wantErrPart: "-contract"},
		{name: "missing deployment", mutate: func(c *Config) { c.DeploymentBlock = nil }, wantErrPart: "-deployment-block"},
		{name: "invalid interval", mutate: func(c *Config) { c.Interval = 0 }, wantErrPart: "-interval"},
		{name: "invalid batch size", mutate: func(c *Config) { c.BatchSize = 0 }, wantErrPart: "-batch-size"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := valid
			if test.mutate != nil {
				test.mutate(&config)
			}
			err := config.Validate()
			if test.wantErrPart == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.wantErrPart) {
				t.Fatalf("Validate() error = %v, want substring %q", err, test.wantErrPart)
			}
		})
	}
}

func TestRunRequiresStorageBeforeDial(t *testing.T) {
	deploymentBlock := uint64(1)
	err := Run(context.Background(), Config{
		RPC:             "http://not-used.invalid",
		Contract:        "0x0000000000000000000000000000000000000001",
		DeploymentBlock: &deploymentBlock,
		Interval:        time.Second,
		BatchSize:       100,
	}, Dependencies{})
	if err == nil || !strings.Contains(err.Error(), "storage") {
		t.Fatalf("Run() error = %v, want missing storage error", err)
	}
}
