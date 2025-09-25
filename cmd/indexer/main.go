package main

import (
	"os"
	"path/filepath"

	"github.com/vedantd/evm-indexer/internal/config"
	"github.com/vedantd/evm-indexer/internal/logging"
	"github.com/vedantd/evm-indexer/internal/version"

	"github.com/rs/zerolog/log"
)

func main() {
	logging.Init()

	log.Info().
		Str("version", version.Version).
		Msg("starting evm-indexer")

	path := os.Getenv("EVMI_CONFIG")
	if path == "" {
		path = "internal/config/chains.yaml"
	}
	abs, _ := filepath.Abs(path)

	cfg, err := config.LoadFromFile(path)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("path", abs).
			Msg("failed to load config")
	}

	log.Info().
		Int("chains", len(cfg.Chains)).
		Str("path", abs).
		Msg("loaded config")

	for _, c := range cfg.Chains {
		log.Info().
			Str("name", c.Name).
			Uint64("id", c.ChainID).
			Uint64("start", c.StartBlock).
			Int("batch", c.BatchSize).
			Str("mode", c.ReceiptsMode).
			Msg("configured chain")
	}
}
