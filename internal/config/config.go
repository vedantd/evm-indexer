package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Chain represents one blockchain configuration entry.
type Chain struct {
	Name         string `yaml:"name"`
	ChainID      uint64 `yaml:"chain_id"`
	RPCHTTP      string `yaml:"rpc_http"`
	RPCWS        string `yaml:"rpc_ws"`
	StartBlock   uint64 `yaml:"start_block"`
	BatchSize    int    `yaml:"batch_size"`
	ReceiptsMode string `yaml:"receipts_mode"`
}

// AppConfig is the root of our YAML file.
type AppConfig struct {
	Chains []Chain `yaml:"chains"`
}

// LoadFromFile loads a YAML config into AppConfig.
func LoadFromFile(path string) (AppConfig, error) {
	var cfg AppConfig

	b, err := os.ReadFile(path) // read raw bytes from file
	if err != nil {
		return cfg, fmt.Errorf("read config file: %w", err)
	}

	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return cfg, fmt.Errorf("parse yaml: %w", err)
	}

	if len(cfg.Chains) == 0 {
		return cfg, fmt.Errorf("no chains configured")
	}

	return cfg, nil
}
