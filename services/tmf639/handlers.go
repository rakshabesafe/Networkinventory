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

func SetDB(d *db.Neo4jDB) {
	neo4jDB = d
}

var (
	mockResources = map[string]models.Resource{
		"1": {ID: "1", Name: "Router-1", LifecycleState: "operating"},
	}
	mockMutex sync.RWMutex
)

func CreateResourceHandler(w http.ResponseWriter, r *http.Request) {
	var resource models.Resource
	if err := json.NewDecoder(r.Body).Decode(&resource); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if resource.ID == "" {
		resource.ID = fmt.Sprintf("gen-%d", time.Now().UnixNano())
	}

	if neo4jDB != nil {
		ctx := context.Background()
		_, err := neo4j.ExecuteQuery(ctx, neo4jDB.Driver,
			"MERGE (r:Resource {id: $id}) ON CREATE SET r.name = $name, r.lifecycleState = $lifecycleState ON MATCH SET r.name = $name, r.lifecycleState = $lifecycleState RETURN r",
			map[string]any{
				"id":             resource.ID,
				"name":           resource.Name,
				"lifecycleState": resource.LifecycleState,
			}, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			http.Error(w, "Failed to create resource in DB: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		mockMutex.Lock()
		mockResources[resource.ID] = resource
		mockMutex.Unlock()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resource)
}

func ListResourcesHandler(w http.ResponseWriter, r *http.Request) {
	var resources []models.Resource

	if neo4jDB != nil {
		ctx := context.Background()
		result, err := neo4j.ExecuteQuery(ctx, neo4jDB.Driver,
			"MATCH (r:Resource) RETURN r.id AS id, r.name AS name, r.lifecycleState AS lifecycleState",
			nil, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			http.Error(w, "Failed to list resources: "+err.Error(), http.StatusInternalServerError)
			return
		}

		for _, record := range result.Records {
			id, _ := record.Get("id")
			name, _ := record.Get("name")
			state, _ := record.Get("lifecycleState")
			resources = append(resources, models.Resource{
				ID:             id.(string),
				Name:           name.(string),
				LifecycleState: state.(string),
			})
		}
	} else {
		mockMutex.RLock()
		for _, res := range mockResources {
			resources = append(resources, res)
		}
		mockMutex.RUnlock()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resources)
}

func GetResourceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if neo4jDB != nil {
		ctx := context.Background()
		result, err := neo4j.ExecuteQuery(ctx, neo4jDB.Driver,
			"MATCH (r:Resource {id: $id}) RETURN r.id AS id, r.name AS name, r.lifecycleState AS lifecycleState",
			map[string]any{"id": id}, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			http.Error(w, "Failed to get resource: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if len(result.Records) == 0 {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}

		record := result.Records[0]
		rID, _ := record.Get("id")
		name, _ := record.Get("name")
		state, _ := record.Get("lifecycleState")

		resource := models.Resource{
			ID:             rID.(string),
			Name:           name.(string),
			LifecycleState: state.(string),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resource)
		return
	}

	mockMutex.RLock()
	resource, ok := mockResources[id]
	mockMutex.RUnlock()

	if !ok {
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}
