package main

import (
	"fmt"

	"github.com/vedantd/evm-indexer/internal/version"
)

func main() {
	fmt.Println("evm-indexer", version.Version)
}
