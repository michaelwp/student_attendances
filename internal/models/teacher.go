package models

import "time"

type Teacher struct {
	ID        uint      `json:"id" db:"id"`
	TeacherID string    `json:"teacher_id" db:"teacher_id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	Phone     *string   `json:"phone" db:"phone"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (Teacher) TableName() string {
	return "teachers"
}