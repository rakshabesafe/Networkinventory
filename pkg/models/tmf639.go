package models

import "time"

// Resource represents a TMF 639 Resource entity
type Resource struct {
	ID                  string    `json:"id"`
	Href                string    `json:"href,omitempty"`
	Category            string    `json:"category,omitempty"`
	Description         string    `json:"description,omitempty"`
	Name                string    `json:"name"`
	ResourceVersion     string    `json:"resourceVersion,omitempty"`
	LifecycleState      string    `json:"lifecycleState,omitempty"` // e.g., "operating", "planned"
	OperationalState    string    `json:"operationalState,omitempty"` // e.g., "enable", "disable"
	UsageState          string    `json:"usageState,omitempty"` // e.g., "idle", "active"
	StartOperatingDate  time.Time `json:"startOperatingDate,omitempty"`
	Type                string    `json:"@type,omitempty"`
	BaseType            string    `json:"@baseType,omitempty"`
	SchemaLocation      string    `json:"@schemaLocation,omitempty"`
}
