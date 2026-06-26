package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed frontend/dist
var frontendAssets embed.FS

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve embedded frontend assets
	frontendFS, err := fs.Sub(frontendAssets, "frontend/dist")
	if err != nil {
		log.Fatalf("Failed to load frontend assets: %v", err)
	}

	mux := http.NewServeMux()

	// API health endpoint
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "project": "{{PROJECT_NAME}}"})
	})

	// Serve frontend SPA
	mux.Handle("/", http.FileServer(http.FS(frontendFS)))

	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
