package models

import "time"

type User struct {
	ID        int       `json:"id" sql:",pk"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Roles     []Role    `json:"roles" pg:"many2many:user_role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
