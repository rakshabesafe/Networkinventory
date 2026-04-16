package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"inventory/pkg/models"
)

func TestGetResourceHandler(t *testing.T) {
	// Setup router to allow path variables
	r := mux.NewRouter()
	r.HandleFunc("/resource/{id}", GetResourceHandler)

	req, err := http.NewRequest("GET", "/resource/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var res models.Resource
	json.NewDecoder(rr.Body).Decode(&res)
	if res.ID != "1" {
		t.Errorf("expected resource ID 1, got %v", res.ID)
	}
}

func TestGetResourceHandler_NotFound(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/resource/{id}", GetResourceHandler)

	req, err := http.NewRequest("GET", "/resource/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
