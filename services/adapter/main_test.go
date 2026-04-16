package main

import (
	"testing"
)

func TestDiscoverResources(t *testing.T) {
	// A simple test to ensure it runs without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("DiscoverResources panicked: %v", r)
		}
	}()

	DiscoverResources()
}
