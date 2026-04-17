package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"inventory/pkg/db"
	"inventory/pkg/models"
)

var neo4jDB *db.Neo4jDB

// SetDB injects the Neo4j database instance into the handler context.
func SetDB(d *db.Neo4jDB) {
	neo4jDB = d
}

var (
	mockServices = map[string]models.Service{
		"1": {ID: "1", Name: "Internet-Access-Service", State: "active"},
	}
	mockMutex sync.RWMutex
)

// CreateServiceHandler handles HTTP POST requests to create a new Service entity.
// It parses the JSON payload, generates an ID if missing, and persists the entity
// to the Neo4j database using a Cypher MERGE query. If the DB is unavailable,
// it falls back to an in-memory thread-safe map.
func CreateServiceHandler(w http.ResponseWriter, r *http.Request) {
	var service models.Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if service.ID == "" {
		service.ID = fmt.Sprintf("gen-%d", time.Now().UnixNano())
	}

	if neo4jDB != nil {
		ctx := context.Background()
		_, err := neo4j.ExecuteQuery(ctx, neo4jDB.Driver,
			"MERGE (s:Service {id: $id}) ON CREATE SET s.name = $name, s.state = $state ON MATCH SET s.name = $name, s.state = $state RETURN s",
			map[string]any{
				"id":    service.ID,
				"name":  service.Name,
				"state": service.State,
			}, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			http.Error(w, "Failed to create service in DB: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		mockMutex.Lock()
		mockServices[service.ID] = service
		mockMutex.Unlock()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(service)
}

// ListServicesHandler handles HTTP GET requests to retrieve a list of all Service entities.
// It queries the Neo4j database to fetch all Service nodes. If the DB is unavailable,
// it returns the items stored in the in-memory fallback map.
func ListServicesHandler(w http.ResponseWriter, r *http.Request) {
	var services []models.Service

	if neo4jDB != nil {
		ctx := context.Background()
		result, err := neo4j.ExecuteQuery(ctx, neo4jDB.Driver,
			"MATCH (s:Service) RETURN s.id AS id, s.name AS name, s.state AS state",
			nil, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			http.Error(w, "Failed to list services: "+err.Error(), http.StatusInternalServerError)
			return
		}

		for _, record := range result.Records {
			id, _ := record.Get("id")
			name, _ := record.Get("name")
			state, _ := record.Get("state")
			services = append(services, models.Service{
				ID:    id.(string),
				Name:  name.(string),
				State: state.(string),
			})
		}
	} else {
		mockMutex.RLock()
		for _, srv := range mockServices {
			services = append(services, srv)
		}
		mockMutex.RUnlock()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

// GetServiceHandler handles HTTP GET requests to retrieve a specific Service entity by its ID.
// It extracts the ID from the URL path variables and queries the Neo4j database.
// If the DB is unavailable, it looks up the item in the in-memory fallback map.
func GetServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if neo4jDB != nil {
		ctx := context.Background()
		result, err := neo4j.ExecuteQuery(ctx, neo4jDB.Driver,
			"MATCH (s:Service {id: $id}) RETURN s.id AS id, s.name AS name, s.state AS state",
			map[string]any{"id": id}, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			http.Error(w, "Failed to get service: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if len(result.Records) == 0 {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}

		record := result.Records[0]
		sID, _ := record.Get("id")
		name, _ := record.Get("name")
		state, _ := record.Get("state")

		service := models.Service{
			ID:    sID.(string),
			Name:  name.(string),
			State: state.(string),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(service)
		return
	}

	mockMutex.RLock()
	service, ok := mockServices[id]
	mockMutex.RUnlock()

	if !ok {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(service)
}
