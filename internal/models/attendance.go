package models

import "time"

type AttendanceStatus string

const (
	AttendanceStatusPresent AttendanceStatus = "present"
	AttendanceStatusAbsent  AttendanceStatus = "absent"
	AttendanceStatusLate    AttendanceStatus = "late"
	AttendanceStatusExcused AttendanceStatus = "excused"
)

type Attendance struct {
	ID          uint             `json:"id" db:"id"`
	StudentID   string           `json:"student_id" db:"student_id"`
	ClassID     uint             `json:"class_id" db:"class_id"`
	Date        time.Time        `json:"date" db:"date"`
	Status      AttendanceStatus `json:"status" db:"status"`
	Description *string          `json:"description" db:"description"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time       `json:"updated_at" db:"updated_at"`
	TimeIn      time.Time        `json:"time_in" db:"time_in"`
	TimeOut     *time.Time       `json:"time_out" db:"time_out"`
	CreatedBy   uint             `json:"created_by" db:"created_by"`
	UpdatedBy   *uint            `json:"updated_by" db:"updated_by"`
	DeletedAt   *time.Time       `json:"deleted_at" db:"deleted_at"`
	DeletedBy   *uint            `json:"deleted_by" db:"deleted_by"`
}

type AttendanceWithStats struct {
	Attendance
	TotalAttendances int `json:"total_attendances" db:"total_attendances"`
	TotalPresent     int `json:"total_present" db:"total_present"`
	TotalAbsent      int `json:"total_absent" db:"total_absent"`
	TotalLate        int `json:"total_late" db:"total_late"`
	TotalExcused     int `json:"total_excused" db:"total_excused"`
}

func (Attendance) TableName() string {
	return "attendances"
}
