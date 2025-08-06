package repository

import (
	"database/sql"
)

// NewRepositories creates a new instance of all repositories
func NewRepositories(db *sql.DB) *Repositories {
	teacherRepo := NewTeacherRepository(db)
	classRepo := NewClassRepository(db)
	studentRepo := NewStudentRepository(db)
	attendanceRepo := NewAttendanceRepository(db)
	absentRequestRepo := NewAbsentRequestRepository(db)
	adminRepo := NewAdminRepositoryWithDeps(db, teacherRepo, studentRepo, classRepo, attendanceRepo)
	
	return &Repositories{
		Teacher:       teacherRepo,
		Class:         classRepo,
		Student:       studentRepo,
		Attendance:    attendanceRepo,
		AbsentRequest: absentRequestRepo,
		Admin:         adminRepo,
	}
}
