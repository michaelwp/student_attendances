package models

import "time"

type Admin struct {
	ID        uint       `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"-" db:"password"`
	LastLogin *time.Time `json:"last_login" db:"last_login"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type AdminStats struct {
	TotalAdmins    int `json:"total_admins" db:"total_admins"`
	ActiveAdmins   int `json:"active_admins" db:"active_admins"`
	InactiveAdmins int `json:"inactive_admins" db:"inactive_admins"`
}

func (Admin) TableName() string {
	return "admins"
}
