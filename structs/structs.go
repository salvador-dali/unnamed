// Package structs contains structs for all models
package structs

import (
	"time"
)

// ErrorCode stores code of a problem that happened while processing client's request.
// It is sent together with 404 status code. It is up to a client how to present it
type ErrorCode struct {
	Id int `json:"error"`
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

// Tag stores all information about a Tag model
type Tag struct {
	Id          int        `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Issued_at   *time.Time `json:"issued_at,omitempty"`
}

// User stores all information about a User model
type User struct {
	Id            int        `json:"id,omitempty"`
	Nickname      string     `json:"nickname,omitempty"`
	Image         string     `json:"image,omitempty"`
	About         string     `json:"about,omitempty"`
	Expertise     int        `json:"expertise,omitempty"`
	Followers_num int        `json:"followers_num,omitempty"`
	Following_num int        `json:"following_num,omitempty"`
	Purchases_num int        `json:"purchases_num,omitempty"`
	Questions_num int        `json:"questions_num,omitempty"`
	Answers_num   int        `json:"answers_num,omitempty"`
	Issued_at     *time.Time `json:"issued_at,omitempty"`
}

// Purchase stores all information about a Purchase model
type Purchase struct {
	Id          int    `json:"id,omitempty"`
	Image       string `json:"image,omitempty"`
	Description string `json:"description,omitempty"`
	User_id     string `json:"user_id,omitempty"`
	Issued_at   string `json:"issued_at,omitempty"`
	Tags        []int  `json:"tags,omitempty"`
	Brand       int    `json:"brand,omitempty"`
	Likes_num   int    `json:"likes_num,omitempty"`
}
