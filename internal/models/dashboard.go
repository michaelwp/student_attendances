package models

// DashboardStats represents comprehensive dashboard statistics
type DashboardStats struct {
	// Admin Stats
	TotalAdmins    int `json:"total_admins" db:"total_admins"`
	ActiveAdmins   int `json:"active_admins" db:"active_admins"`
	InactiveAdmins int `json:"inactive_admins" db:"inactive_admins"`

	// Teacher Stats
	TotalTeachers    int `json:"total_teachers" db:"total_teachers"`
	ActiveTeachers   int `json:"active_teachers" db:"active_teachers"`
	InactiveTeachers int `json:"inactive_teachers" db:"inactive_teachers"`

	// Student Stats
	TotalStudents    int `json:"total_students" db:"total_students"`
	ActiveStudents   int `json:"active_students" db:"active_students"`
	InactiveStudents int `json:"inactive_students" db:"inactive_students"`

	// Class Stats
	TotalClasses int `json:"total_classes" db:"total_classes"`

	// Today's Attendance (if available)
	TotalAttendanceToday int `json:"total_attendance_today,omitempty"`
	PresentToday         int `json:"present_today,omitempty"`
	AbsentToday          int `json:"absent_today,omitempty"`
	LateToday            int `json:"late_today,omitempty"`
}