package handlers

import (
	"github.com/michaelwp/student_attendance/internal/repository"
)

// NewHandlers creates a new instance of all handlers
func NewHandlers(repos *repository.Repositories) *Handlers {
	return &Handlers{
		Teacher:       NewTeacherHandler(repos.Teacher),
		Class:         NewClassHandler(repos.Class),
		Student:       NewStudentHandler(repos.Student),
		Attendance:    NewAttendanceHandler(repos.Attendance),
		AbsentRequest: NewAbsentRequestHandler(repos.AbsentRequest),
	}
}
