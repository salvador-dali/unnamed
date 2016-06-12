// Package structs contains structs for all models
package structs

import (
	"time"
)

// ErrorCode stores code of a problem that happened while processing client's request.
// It is sent together with 404 status code. It is up to a client how to present it
type ErrorCode struct {
	Id int `json:"id"`
}

// Id stores information about the id of the element which was just inserted
type Id struct {
	Id int `json:"id"`
}

// Brand stores all information about a Brand model
type Brand struct {
	Id        int        `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Issued_at *time.Time `json:"issued_at,omitempty"`
}
