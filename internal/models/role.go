package models

import "time"

type Role struct {
	ID        int       `json:"id" sql:",pk"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
