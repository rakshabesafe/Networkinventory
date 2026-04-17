package models

import "time"

// Service represents a TMF 638 Service entity.
// It holds the attributes that define a service in the inventory system,
// such as its state, operational dates, and categorical information.
type Service struct {
	ID                  string    `json:"id"`
	Href                string    `json:"href,omitempty"`
	Category            string    `json:"category,omitempty"`
	Description         string    `json:"description,omitempty"`
	HasStarted          bool      `json:"hasStarted,omitempty"`
	IsServiceEnabled    bool      `json:"isServiceEnabled,omitempty"`
	IsStateful          bool      `json:"isStateful,omitempty"`
	Name                string    `json:"name"`
	ServiceDate         time.Time `json:"serviceDate,omitempty"`
	ServiceType         string    `json:"serviceType,omitempty"`
	StartDate           time.Time `json:"startDate,omitempty"`
	StartMode           string    `json:"startMode,omitempty"`
	State               string    `json:"state,omitempty"` // e.g., "active", "inactive", "terminated"
	Type                string    `json:"@type,omitempty"`
	BaseType            string    `json:"@baseType,omitempty"`
	SchemaLocation      string    `json:"@schemaLocation,omitempty"`
}
