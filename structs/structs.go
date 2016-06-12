// Package structs contains structs for all models
package structs

import (
	"time"
)

// Id is a placeholder struct to output id that client asked in unimplemented handler
type Id struct {
	Id int `json:"id"`
}

// Brand is a struct that stores all information about a Brand model
type Brand struct {
	Id        int       `json:"id"`
	Name      string    `json:"name,omitempty"`
	Issued_at time.Time `json:"issued_at,omitempty"`
}
