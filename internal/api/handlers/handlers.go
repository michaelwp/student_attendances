package handlers

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/michaelwp/student_attendance/internal/config"
	"github.com/michaelwp/student_attendance/internal/repository"
	"github.com/redis/go-redis/v9"
)

type HandlerDependencies struct {
	Repositories *repository.Repositories
	S3Client     *s3.Client
	S3Config     *config.S3Config
	RedisClient  *redis.Client
}

// NewHandlers creates a new instance of all handlers
func NewHandlers(dep *HandlerDependencies) *Handlers {
	return &Handlers{
		Teacher:       NewTeacherHandler(dep.Repositories.Teacher, dep.S3Client, dep.S3Config),
		Class:         NewClassHandler(dep.Repositories.Class),
		Student:       NewStudentHandler(dep.Repositories.Student, dep.S3Client, dep.S3Config),
		Attendance:    NewAttendanceHandler(dep.Repositories.Attendance, dep.Repositories.Student),
		AbsentRequest: NewAbsentRequestHandler(dep.Repositories.AbsentRequest),
		Admin:         NewAdminHandler(dep.Repositories.Admin),
		Auth:          NewAuthHandler(dep.Repositories.Admin, dep.Repositories.Teacher, dep.Repositories.Student, dep.RedisClient),
	}
}
