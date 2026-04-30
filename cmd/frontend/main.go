package main

import (
	"os"

	_ "time/tzdata"

	"github.com/aceberg/ExerciseDiary/internal/store"
	"github.com/aceberg/ExerciseDiary/internal/web"
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// All configuration is read from environment variables:
//
//	PORT      listen port for the web UI  (default: 8080)
//	API_URL   base URL of the API server  (default: http://localhost:8851)
//	API_KEY   X-Api-Key sent to the API   (default: "", no auth)
//	NODE_PATH path to local node_modules  (default: "", use CDN)
func main() {
	port := envOr("PORT", "8080")
	apiURL := envOr("API_URL", "http://localhost:8851")
	apiKey := envOr("API_KEY", "")
	nodePath := envOr("NODE_PATH", "")

	ac := store.NewAPIClient(apiURL, apiKey)
	web.GuiWithStore(ac, ac, port, nodePath)
}
