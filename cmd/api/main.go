package main

import (
	"flag"
	"os"

	_ "time/tzdata"

	"github.com/aceberg/ExerciseDiary/internal/api"
)

func main() {
	// DATA_DIR env var sets the data directory; -d flag overrides it.
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "/data/ExerciseDiary"
	}
	flag.StringVar(&dataDir, "d", dataDir, "Path to data directory (overrides DATA_DIR env var)")
	flag.Parse()

	api.Start(dataDir)
}
