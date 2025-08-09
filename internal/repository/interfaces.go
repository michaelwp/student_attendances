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
	UpdatePassword(ctx context.Context, teacherID string, password string) error
	GetPasswordByTeacherID(ctx context.Context, teacherID string) (string, error)
	IsTeacherExist(ctx context.Context, teacherID string) (bool, error)
	GetStats(ctx context.Context) (*models.TeacherStats, error)
	UpdateDeleteInfo(ctx context.Context, id uint, deletedBy uint) error
}

// ClassRepository defines the interface for class operations
type ClassRepository interface {
	Create(ctx context.Context, class *models.Class) error
	GetByID(ctx context.Context, id uint) (*models.Class, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Class, error)
	GetByTeacher(ctx context.Context, teacherID string) ([]*models.Class, error)
	Update(ctx context.Context, class *models.Class) error
	Delete(ctx context.Context, id uint) error
	GetTotalClasses(ctx context.Context) (int, error)
	UpdateDeleteInfo(ctx context.Context, id uint, deletedBy uint) error
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
	UpdatePassword(ctx context.Context, studentID string, password string) error
	GetPasswordByStudentID(ctx context.Context, studentID string) (string, error)
	IsStudentExist(ctx context.Context, studentID string) (bool, error)
	GetStats(ctx context.Context) (*models.StudentStats, error)
	UpdateDeleteInfo(ctx context.Context, id uint, deletedBy uint) error
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
	UpdateDeleteInfo(ctx context.Context, id uint, deletedBy uint) error
	GetAll(ctx context.Context, limit, offset int) ([]*models.Attendance, error)
	GetCount(ctx context.Context) (int, error)
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

// AdminRepository defines the interface for admin operations
type AdminRepository interface {
	Create(ctx context.Context, admin *models.Admin) error
	GetByID(ctx context.Context, id uint) (*models.Admin, error)
	GetByEmail(ctx context.Context, email string) (*models.Admin, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Admin, error)
	Update(ctx context.Context, admin *models.Admin) error
	Delete(ctx context.Context, id uint) error
	UpdatePassword(ctx context.Context, email string, password string) error
	UpdateLastLogin(ctx context.Context, id uint, lastLogin time.Time) error
	SetActiveStatus(ctx context.Context, id uint, isActive bool) error
	GetTotalAdmins(ctx context.Context) (int, error)
	IsAdminExist(ctx context.Context, email string) (bool, error)
	GetPasswordByEmail(ctx context.Context, email string) (string, error)
	GetStats(ctx context.Context) (*models.AdminStats, error)
	GetDashboardStats(ctx context.Context) (*models.DashboardStats, error)
}

// Repositories aggregates all repository interfaces
type Repositories struct {
	Teacher       TeacherRepository
	Class         ClassRepository
	Student       StudentRepository
	Attendance    AttendanceRepository
	AbsentRequest AbsentRequestRepository
	Admin         AdminRepository
}
