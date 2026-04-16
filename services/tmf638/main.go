package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"inventory/pkg/db"
)

func main() {
	r := mux.NewRouter()

	// Initialize DB connection
	dbConn, err := db.Connect()
	if err != nil {
		log.Printf("Warning: Failed to connect to DB, running with mock data (or will fail if DB needed): %v", err)
	} else {
		log.Printf("Successfully connected to Neo4j for TMF 638")
		SetDB(dbConn)
	}

	r.HandleFunc("/tmf-api/serviceInventoryManagement/v4/service", CreateServiceHandler).Methods("POST")
	r.HandleFunc("/tmf-api/serviceInventoryManagement/v4/service", ListServicesHandler).Methods("GET")
	r.HandleFunc("/tmf-api/serviceInventoryManagement/v4/service/{id}", GetServiceHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Starting TMF 638 Service Inventory on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
