package repository

import (
	"context"
	"time"

	"github.com/michaelwp/student_attendance/internal/models"
)

// TeacherRepository defines the interface for teacher operations
type TeacherRepository interface {
	Create(ctx context.Context, teacher *models.Teacher) error
	GetByID(ctx context.Context, id uint) (*models.Teacher, error)
	GetByTeacherID(ctx context.Context, teacherID string) (*models.Teacher, error)
	GetByEmail(ctx context.Context, email string) (*models.Teacher, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Teacher, error)
	Update(ctx context.Context, teacher *models.Teacher) error
	Delete(ctx context.Context, id uint) error
	UpdatePhotoPath(ctx context.Context, id uint, photoPath string) error
	GetPhotoPath(ctx context.Context, id uint) (string, error)
	GetTotalTeachers(ctx context.Context) (int, error)
}

// ClassRepository defines the interface for class operations
type ClassRepository interface {
	Create(ctx context.Context, class *models.Class) error
	GetByID(ctx context.Context, id uint) (*models.Class, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Class, error)
	GetByTeacher(ctx context.Context, teacherID string) ([]*models.Class, error)
	Update(ctx context.Context, class *models.Class) error
	Delete(ctx context.Context, id uint) error
}

// StudentRepository defines the interface for student operations
type StudentRepository interface {
	Create(ctx context.Context, student *models.Student) error
	GetByID(ctx context.Context, id uint) (*models.Student, error)
	GetByStudentID(ctx context.Context, studentID string) (*models.Student, error)
	GetByEmail(ctx context.Context, email string) (*models.Student, error)
	GetByClass(ctx context.Context, classID uint) ([]*models.Student, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Student, error)
	Update(ctx context.Context, student *models.Student) error
	Delete(ctx context.Context, id uint) error
	UpdatePhotoPath(ctx context.Context, id uint, photoPath string) error
	GetPhotoPath(ctx context.Context, id uint) (string, error)
	GetTotalStudents(ctx context.Context) (int, error)
}

// AttendanceRepository defines the interface for attendance operations
type AttendanceRepository interface {
	Create(ctx context.Context, attendance *models.Attendance) error
	GetByID(ctx context.Context, id uint) (*models.Attendance, error)
	GetByStudentAndDate(ctx context.Context, studentID string, date time.Time) (*models.Attendance, error)
	GetByStudent(ctx context.Context, studentID string, limit, offset int) ([]*models.Attendance, error)
	GetByClass(ctx context.Context, classID uint, limit, offset int) ([]*models.Attendance, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*models.Attendance, error)
	Update(ctx context.Context, attendance *models.Attendance) error
	Delete(ctx context.Context, id uint) error
}

// AbsentRequestRepository defines the interface for absent request operations
type AbsentRequestRepository interface {
	Create(ctx context.Context, request *models.AbsentRequest) error
	GetByID(ctx context.Context, id uint) (*models.AbsentRequest, error)
	GetByStudent(ctx context.Context, studentID string, limit, offset int) ([]*models.AbsentRequest, error)
	GetByClass(ctx context.Context, classID uint, limit, offset int) ([]*models.AbsentRequest, error)
	GetByStatus(ctx context.Context, status models.AbsentRequestStatus, limit, offset int) ([]*models.AbsentRequest, error)
	GetPending(ctx context.Context, limit, offset int) ([]*models.AbsentRequest, error)
	Update(ctx context.Context, request *models.AbsentRequest) error
	Delete(ctx context.Context, id uint) error
}

// Repositories aggregates all repository interfaces
type Repositories struct {
	Teacher       TeacherRepository
	Class         ClassRepository
	Student       StudentRepository
	Attendance    AttendanceRepository
	AbsentRequest AbsentRequestRepository
}
