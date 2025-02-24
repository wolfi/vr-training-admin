// cmd/server/main.go
package main

import (
	"log"
	"net/http"

	"github.com/saladinomario/vr-training-admin/internal/handlers"
)

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Register dashboard handler
	mux.HandleFunc("/", handlers.DashboardHandler)

	// Register scenario routes
	handlers.SetupScenarioRoutes(mux)

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}

func main() {
	mux := setupRoutes()

	log.Println("Server starting on :8080")
	log.Println("Visit http://localhost:8080 to view the application")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
