package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Neo4jDB wraps the Neo4j driver connection to provide a centralized
// interface for database interactions across microservices.
type Neo4jDB struct {
	Driver neo4j.DriverWithContext
}

// Connect initializes and verifies a connection to the Neo4j database.
// It reads connection parameters from environment variables (NEO4J_URI, NEO4J_USER, NEO4J_PASSWORD)
// and implements a retry mechanism to handle database startup delays.
func Connect() (*Neo4jDB, error) {
	uri := os.Getenv("NEO4J_URI")
	if uri == "" {
		uri = "neo4j://localhost:7687"
	}
	user := os.Getenv("NEO4J_USER")
	if user == "" {
		user = "neo4j"
	}
	password := os.Getenv("NEO4J_PASSWORD")
	if password == "" {
		password = "password"
	}

	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	ctx := context.Background()

	// Retry logic
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		err = driver.VerifyConnectivity(ctx)
		if err == nil {
			log.Println("Successfully connected to Neo4j")
			return &Neo4jDB{Driver: driver}, nil
		}
		log.Printf("Failed to verify connectivity (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(3 * time.Second)
	}

	return nil, fmt.Errorf("failed to verify connectivity after %d attempts: %w", maxRetries, err)
}

// Close terminates the connection to the Neo4j database gracefully.
func (db *Neo4jDB) Close(ctx context.Context) error {
	if db.Driver != nil {
		return db.Driver.Close(ctx)
	}
	return nil
}
