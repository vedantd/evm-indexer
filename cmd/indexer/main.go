package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vedantd/evm-indexer/internal/config"
	"github.com/vedantd/evm-indexer/internal/version"
)

func main() {
	fmt.Println("evm-indexer", version.Version)

	// Allow override via EVMI_CONFIG, default to configs/chains.yaml
	path := os.Getenv("EVMI_CONFIG")
	if path == "" {
		path = "internal/config/chains.yaml"
	}
	abs, _ := filepath.Abs(path)

	cfg, err := config.LoadFromFile(path)
	if err != nil {
		fmt.Printf("config error: %v (path=%s)\n", err, abs)
		os.Exit(1)
	}

	fmt.Printf("loaded %d chain(s) from %s\n", len(cfg.Chains), abs)
	for _, c := range cfg.Chains {
		fmt.Printf("- %s (id=%d) start=%d batch=%d mode=%s\n",
			c.Name, c.ChainID, c.StartBlock, c.BatchSize, c.ReceiptsMode)
	}
}
