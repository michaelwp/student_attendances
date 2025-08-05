package models

import "time"

type AttendanceStatus string

const (
	AttendanceStatusPresent AttendanceStatus = "present"
	AttendanceStatusAbsent  AttendanceStatus = "absent"
	AttendanceStatusLate    AttendanceStatus = "late"
)

type Attendance struct {
	ID          uint             `json:"id" db:"id"`
	StudentID   string           `json:"student_id" db:"student_id"`
	ClassID     uint             `json:"class_id" db:"class_id"`
	Date        time.Time        `json:"date" db:"date"`
	Status      AttendanceStatus `json:"status" db:"status"`
	Description *string          `json:"description" db:"description"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
}

func (Attendance) TableName() string {
	return "attendances"
}