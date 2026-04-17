package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"inventory/pkg/db"
)

// main is the entry point for the TMF 639 Resource Inventory microservice.
// It initializes the database connection, sets up the HTTP router with TMF-compliant endpoints,
// and starts the HTTP server.
func main() {
	r := mux.NewRouter()

	// Initialize DB connection
	dbConn, err := db.Connect()
	if err != nil {
		log.Printf("Warning: Failed to connect to DB, running with mock data (or will fail if DB needed): %v", err)
	} else {
		log.Printf("Successfully connected to Neo4j for TMF 639")
		// Inject DB dependency into handlers
		SetDB(dbConn)
	}

	r.HandleFunc("/tmf-api/resourceInventoryManagement/v4/resource", CreateResourceHandler).Methods("POST")
	r.HandleFunc("/tmf-api/resourceInventoryManagement/v4/resource", ListResourcesHandler).Methods("GET")
	r.HandleFunc("/tmf-api/resourceInventoryManagement/v4/resource/{id}", GetResourceHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Starting TMF 639 Resource Inventory on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
