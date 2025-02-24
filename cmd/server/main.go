// cmd/server/main.go
package main

import (
	"log"
	"net/http"

	"github.com/saladinomario/vr-training-admin/internal/handlers"
)

// Track registered routes
var registeredRoutes []string

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Register dashboard handler
	mux.HandleFunc("/", handlers.DashboardHandler)
	registeredRoutes = append(registeredRoutes, "/ -> DashboardHandler")

	// Register scenario routes
	handlers.SetupScenarioRoutes(mux)
	registeredRoutes = append(registeredRoutes, "/scenarios -> SetupScenarioRoutes")

	// Register avatar routes
	handlers.SetupAvatarRoutes(mux)
	registeredRoutes = append(registeredRoutes, "/avatar-lab -> SetupAvatarRoutes")

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	registeredRoutes = append(registeredRoutes, "/static/ -> Static Files")

	return mux
}

// Print dynamically registered routes
func printRegisteredRoutes() {
	log.Println("=== Registered Routes ===")
	for _, route := range registeredRoutes {
		log.Println(route)
	}
	log.Println("==========================")
}

func main() {
	mux := setupRoutes()
	printRegisteredRoutes()

	log.Println("Server starting on :8080")
	log.Println("Visit http://localhost:8080 to view the application")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
