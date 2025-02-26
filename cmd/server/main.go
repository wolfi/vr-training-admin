package main

import (
	"log"
	"net/http"

	"github.com/saladinomario/vr-training-admin/internal/handlers"
)

// Track registered routes
var registeredRoutes []string

func setupRoutes() *http.ServeMux {
	log.Println("Setting up all application routes...")

	mux := http.NewServeMux()

	// Register dashboard handler
	log.Println("Registering dashboard route")
	mux.HandleFunc("/", handlers.DashboardHandler)
	mux.HandleFunc("/dashboard-content", handlers.DashboardContentHandler)
	// Register scenario routes
	log.Println("Setting up scenario routes")
	handlers.SetupScenarioRoutes(mux)

	// Register avatar routes
	log.Println("Setting up avatar routes")
	handlers.SetupAvatarRoutes(mux)

	// Register observer routes
	log.Println("Setting up observer routes")
	handlers.SetupObserverRoutes(mux)

	// Register settings routes
	log.Println("Setting up settings routes")
	handlers.SetupSettingsRoutes(mux)

	// Register session routes (new)
	log.Println("Setting up session routes")
	handlers.SetupSessionRoutes(mux)

	// Serve static files
	log.Println("Setting up static file server")
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("All routes registered successfully")
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
