package models

import "time"

type Teacher struct {
	ID        uint       `json:"id" db:"id"`
	TeacherID string     `json:"teacher_id" db:"teacher_id"`
	FirstName string     `json:"first_name" db:"first_name"`
	LastName  string     `json:"last_name" db:"last_name"`
	Email     string     `json:"email" db:"email"`
	Phone     *string    `json:"phone" db:"phone"`
	Password  string     `json:"-" db:"password"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	DeletedBy *uint      `json:"deleted_by" db:"deleted_by"`
}

type TeacherStats struct {
	TotalTeachers    int `json:"total_teachers" db:"total_teachers"`
	ActiveTeachers   int `json:"active_teachers" db:"active_teachers"`
	InactiveTeachers int `json:"inactive_teachers" db:"inactive_teachers"`
}

func (Teacher) TableName() string {
	return "teachers"
}
