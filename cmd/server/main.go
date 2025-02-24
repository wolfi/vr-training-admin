package main

import (
    "log"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    
    // Serve static files
    fs := http.FileServer(http.Dir("static"))
    mux.Handle("/static/", http.StripPrefix("/static/", fs))
    
    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatal(err)
    }
}
