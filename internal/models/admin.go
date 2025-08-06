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

func (Admin) TableName() string {
	return "admins"
}