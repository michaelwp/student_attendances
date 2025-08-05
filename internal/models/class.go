package models

import "time"

type Class struct {
	ID               uint      `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	HomeroomTeacher  string    `json:"homeroom_teacher" db:"homeroom_teacher"`
	Description      *string   `json:"description" db:"description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

func (Class) TableName() string {
	return "classes"
}