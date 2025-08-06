package models

import "time"

type Student struct {
	ID        uint      `json:"id" db:"id"`
	StudentID string    `json:"student_id" db:"student_id"`
	ClassesID uint      `json:"classes_id" db:"classes_id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	Phone     *string   `json:"phone" db:"phone"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	IsActive  bool      `json:"is_active" db:"is_active"`
}

type StudentStats struct {
	TotalStudents    int `json:"total_students" db:"total_students"`
	ActiveStudents   int `json:"active_students" db:"active_students"`
	InactiveStudents int `json:"inactive_students" db:"inactive_students"`
}

func (Student) TableName() string {
	return "students"
}
