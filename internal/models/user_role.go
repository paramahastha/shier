package models

import "time"

type UserRole struct {
	tableName struct{}  `sql:"user_role"`
	UserID    int       `json:"user_id" sql:",pk"`
	RoleID    int       `json:"role_id" sql:",pk"`
	GrantDate time.Time `json:"grant_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
