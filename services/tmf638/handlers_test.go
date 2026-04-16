package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"inventory/pkg/models"
)

func TestGetServiceHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/service/{id}", GetServiceHandler)

	req, err := http.NewRequest("GET", "/service/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var srv models.Service
	json.NewDecoder(rr.Body).Decode(&srv)
	if srv.ID != "1" {
		t.Errorf("expected service ID 1, got %v", srv.ID)
	}
}
