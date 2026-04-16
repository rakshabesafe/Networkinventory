package db

import (
	"context"
	"testing"
)

func TestConnect_FailureNoDB(t *testing.T) {
	// Without a running DB, VerifyConnectivity should fail
	db, err := Connect()
	if err == nil {
		t.Errorf("Expected an error since no neo4j database is running locally")
		db.Close(context.Background())
	}
}
