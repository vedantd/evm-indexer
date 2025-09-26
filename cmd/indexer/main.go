package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/vedantd/evm-indexer/internal/config"
	"github.com/vedantd/evm-indexer/internal/ingest/planner"
	"github.com/vedantd/evm-indexer/internal/logging"
	"github.com/vedantd/evm-indexer/internal/version"
)

// define staticHead + method at package level
type staticHead uint64

func (h staticHead) HeadNumber(ctx context.Context) (uint64, error) {
	return uint64(h), nil
}

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

	// Demo: plan a small backfill for the first chain and log first 10 numbers.
	if len(cfg.Chains) > 0 {
		c := cfg.Chains[0]

		out := make(chan uint64, 1000)
		p := &planner.Planner{
			Heads:        staticHead(c.StartBlock + 500), // pretend head is +500
			BatchSize:    100,
			SafetyWindow: 6,
		}
		if err := p.Plan(context.Background(), c.StartBlock, out); err != nil {
			log.Error().Err(err).Msg("planner demo failed")
		} else {
			max := 10
			for i := 0; i < max; i++ {
				select {
				case n := <-out:
					log.Info().Uint64("planned_block", n).Msg("demo")
				default:
					break
				}
			}
		}
	}
}
