package main

import (
	"flag"

	_ "time/tzdata"

	"github.com/aceberg/ExerciseDiary/internal/store"
	"github.com/aceberg/ExerciseDiary/internal/web"
)

const (
	defaultPort    = "8080"
	defaultAPIURL  = "http://localhost:8851"
	defaultAPIKey  = ""
	defaultDirPath = "/data/ExerciseDiary"
	defaultNode    = ""
)

func main() {
	portPtr := flag.String("p", defaultPort, "Port for the frontend to listen on")
	apiURLPtr := flag.String("a", defaultAPIURL, "Base URL of the ExerciseDiary API server")
	keyPtr := flag.String("k", defaultAPIKey, "API key sent to the API server (X-Api-Key header)")
	dirPtr := flag.String("d", defaultDirPath, "Path to data/config directory (used for static assets)")
	nodePtr := flag.String("n", defaultNode, "Path to node_modules (empty = use CDN)")
	flag.Parse()

	ac := store.NewAPIClient(*apiURLPtr, *keyPtr)
	web.GuiWithStore(ac, ac, *portPtr, *dirPtr, *nodePtr)
}
