package model

import "time"

type Profile struct {
	UserID    int64      `json:"id"`
	Nickname  string     `json:"nickname"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (Profile) TableName() string { return "profiles" }
