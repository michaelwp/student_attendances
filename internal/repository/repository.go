package repository

import "database/sql"

// NewRepositories creates a new instance of all repositories
func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Teacher:       NewTeacherRepository(db),
		Class:         NewClassRepository(db),
		Student:       NewStudentRepository(db),
		Attendance:    NewAttendanceRepository(db),
		AbsentRequest: NewAbsentRequestRepository(db),
	}
}