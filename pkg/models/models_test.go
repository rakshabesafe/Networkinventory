package models

import "testing"

func TestModels(t *testing.T) {
	service := Service{
		ID:   "s1",
		Name: "Test Service",
	}

	if service.ID != "s1" || service.Name != "Test Service" {
		t.Errorf("Service model creation failed")
	}

	resource := Resource{
		ID:   "r1",
		Name: "Test Resource",
	}

	if resource.ID != "r1" || resource.Name != "Test Resource" {
		t.Errorf("Resource model creation failed")
	}
}
