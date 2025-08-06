package handlers

import "github.com/gofiber/fiber/v2"

// TeacherHandler defines the interface for teacher API operations
type TeacherHandler interface {
	Create(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetByTeacherID(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	UploadPhoto(c *fiber.Ctx) error
	GetPhoto(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
	UpdatePassword(c *fiber.Ctx) error
}

// ClassHandler defines the interface for class API operations
type ClassHandler interface {
	Create(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetByTeacher(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

// StudentHandler defines the interface for student API operations
type StudentHandler interface {
	Create(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetByStudentID(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetByClass(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	UploadPhoto(c *fiber.Ctx) error
	GetPhoto(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
	UpdatePassword(c *fiber.Ctx) error
}

// AttendanceHandler defines the interface for attendance API operations
type AttendanceHandler interface {
	Create(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetByStudent(c *fiber.Ctx) error
	GetByClass(c *fiber.Ctx) error
	GetByDateRange(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

// AbsentRequestHandler defines the interface for absent request API operations
type AbsentRequestHandler interface {
	Create(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetByStudent(c *fiber.Ctx) error
	GetByClass(c *fiber.Ctx) error
	GetPending(c *fiber.Ctx) error
	UpdateStatus(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

// AdminHandler defines the interface for admin API operations
type AdminHandler interface {
	Create(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetByEmail(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	UpdatePassword(c *fiber.Ctx) error
	SetActiveStatus(c *fiber.Ctx) error
}

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

// Handlers aggregates all handler interfaces
type Handlers struct {
	Teacher       TeacherHandler
	Class         ClassHandler
	Student       StudentHandler
	Attendance    AttendanceHandler
	AbsentRequest AbsentRequestHandler
	Admin         AdminHandler
	Auth          AuthHandler
}
