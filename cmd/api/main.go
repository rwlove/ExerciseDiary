package main

import (
	"flag"

	_ "time/tzdata"

	"github.com/aceberg/ExerciseDiary/internal/api"
)

const (
	defaultDirPath = "/data/ExerciseDiary"
	defaultPort    = "8851"
	defaultAPIKey  = ""
)

func main() {
	dirPtr := flag.String("d", defaultDirPath, "Path to data/config directory")
	portPtr := flag.String("p", defaultPort, "Port to listen on")
	keyPtr := flag.String("k", defaultAPIKey, "API key required on X-Api-Key header (empty = no auth)")
	flag.Parse()

	api.Start(*dirPtr, *portPtr, *keyPtr)
}
