package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelwp/student_attendance/internal/models"
	"github.com/michaelwp/student_attendance/internal/repository"
)

type studentHandler struct {
	studentRepo repository.StudentRepository
}

// NewStudentHandler creates a new student handler
func NewStudentHandler(studentRepo repository.StudentRepository) StudentHandler {
	return &studentHandler{
		studentRepo: studentRepo,
	}
}

// CreateStudent godoc
// @Summary Create a new student
// @Description Create a new student in the system
// @Tags Students
// @Accept json
// @Produce json
// @Param student body models.Student true "Student data"
// @Success 201 {object} map[string]interface{} "Student created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /students [post]
func (h *studentHandler) Create(c *fiber.Ctx) error {
	var student models.Student
	if err := c.BodyParser(&student); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.studentRepo.Create(c.Context(), &student); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create student",
		})
	}

	student.Password = ""
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Student created successfully",
		"data":    student,
	})
}

// GetStudentByID godoc
// @Summary Get student by ID
// @Description Retrieve a specific student by their database ID
// @Tags Students
// @Accept json
// @Produce json
// @Param id path int true "Student database ID"
// @Success 200 {object} map[string]interface{} "Student retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid student ID"
// @Failure 404 {object} map[string]interface{} "Student not found"
// @Router /students/{id} [get]
func (h *studentHandler) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid student ID",
		})
	}

	student, err := h.studentRepo.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student not found",
		})
	}

	student.Password = ""
	return c.JSON(fiber.Map{
		"data": student,
	})
}

// GetStudentByStudentID godoc
// @Summary Get student by student ID
// @Description Retrieve a specific student by their student ID
// @Tags Students
// @Accept json
// @Produce json
// @Param studentId path string true "Student ID (not database ID)"
// @Success 200 {object} map[string]interface{} "Student retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Student ID is required"
// @Failure 404 {object} map[string]interface{} "Student not found"
// @Router /students/student-id/{studentId} [get]
func (h *studentHandler) GetByStudentID(c *fiber.Ctx) error {
	studentID := c.Params("studentId")
	if studentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student ID is required",
		})
	}

	student, err := h.studentRepo.GetByStudentID(c.Context(), studentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student not found",
		})
	}

	student.Password = ""
	return c.JSON(fiber.Map{
		"data": student,
	})
}

// GetAllStudents godoc
// @Summary Get all students
// @Description Retrieve a paginated list of all students
// @Tags Students
// @Accept json
// @Produce json
// @Param limit query int false "Number of students to return (max 100)" default(10)
// @Param offset query int false "Number of students to skip" default(0)
// @Success 200 {object} map[string]interface{} "List of students retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /students [get]
func (h *studentHandler) GetAll(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	students, err := h.studentRepo.GetAll(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get students",
		})
	}

	for _, student := range students {
		student.Password = ""
	}

	return c.JSON(fiber.Map{
		"data":   students,
		"count":  len(students),
		"limit":  limit,
		"offset": offset,
	})
}

// GetStudentsByClass godoc
// @Summary Get students by class
// @Description Retrieve all students in a specific class
// @Tags Students
// @Accept json
// @Produce json
// @Param classId path int true "Class ID"
// @Success 200 {object} map[string]interface{} "Students retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid class ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /students/class-id/{classId} [get]
func (h *studentHandler) GetByClass(c *fiber.Ctx) error {
	classIDParam := c.Params("classId")
	classID, err := strconv.ParseUint(classIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid class ID",
		})
	}

	students, err := h.studentRepo.GetByClass(c.Context(), uint(classID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get students",
		})
	}

	for _, student := range students {
		student.Password = ""
	}

	return c.JSON(fiber.Map{
		"data":  students,
		"count": len(students),
	})
}

// UpdateStudent godoc
// @Summary Update student
// @Description Update a student's information
// @Tags Students
// @Accept json
// @Produce json
// @Param id path int true "Student database ID"
// @Param student body models.Student true "Student data"
// @Success 200 {object} map[string]interface{} "Student updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /students/{id} [put]
func (h *studentHandler) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid student ID",
		})
	}

	var student models.Student
	if err := c.BodyParser(&student); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	student.ID = uint(id)
	if err := h.studentRepo.Update(c.Context(), &student); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update student",
		})
	}

	student.Password = ""
	return c.JSON(fiber.Map{
		"message": "Student updated successfully",
		"data":    student,
	})
}

// DeleteStudent godoc
// @Summary Delete student
// @Description Delete a student from the system
// @Tags Students
// @Accept json
// @Produce json
// @Param id path int true "Student database ID"
// @Success 200 {object} map[string]interface{} "Student deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid student ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /students/{id} [delete]
func (h *studentHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid student ID",
		})
	}

	if err := h.studentRepo.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete student",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Student deleted successfully",
	})
}
