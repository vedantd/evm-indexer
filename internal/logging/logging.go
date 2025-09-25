package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	// Pretty console output for dev
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
