package main

import (
	"fmt"
	"log"
	"net/http"
	"orgguard/db"
	"orgguard/handlers"
)

func main() {
	var (
		addr string = "localhost:"
		port string = "8080"
	)

	// Initialize database
	db.InitDB()
	db.RunMigrations()

	mux := http.NewServeMux()

	// routes
	mux.HandleFunc("/org", handlers.OrganizationHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)

	// 404 for unmatched routes
	mux.HandleFunc("/", handlers.NotFoundHandler)

	// started server
	fmt.Println("Server Listening At ", addr+port)

	// currently mus is not handled
	err := http.ListenAndServe(addr+port, mux)
	if err != nil {
		log.Fatalf("error listening to server : %v", err)
		return
	}

	fmt.Println("Server Exits..")
}
