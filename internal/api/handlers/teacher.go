package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelwp/student_attendance/internal/models"
	"github.com/michaelwp/student_attendance/internal/repository"
)

type teacherHandler struct {
	teacherRepo repository.TeacherRepository
}

// NewTeacherHandler creates a new teacher handler
func NewTeacherHandler(teacherRepo repository.TeacherRepository) TeacherHandler {
	return &teacherHandler{
		teacherRepo: teacherRepo,
	}
}

// CreateTeacher godoc
// @Summary Create a new teacher
// @Description Create a new teacher in the system
// @Tags Teachers
// @Accept json
// @Produce json
// @Param teacher body models.Teacher true "Teacher data"
// @Success 201 {object} map[string]interface{} "Teacher created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /teachers [post]
func (h *teacherHandler) Create(c *fiber.Ctx) error {
	var teacher models.Teacher
	if err := c.BodyParser(&teacher); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.teacherRepo.Create(c.Context(), &teacher); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create teacher",
		})
	}

	teacher.Password = ""
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Teacher created successfully",
		"data":    teacher,
	})
}

// GetTeacherByID godoc
// @Summary Get teacher by ID
// @Description Retrieve a specific teacher by their database ID
// @Tags Teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher database ID"
// @Success 200 {object} map[string]interface{} "Teacher retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid teacher ID"
// @Failure 404 {object} map[string]interface{} "Teacher not found"
// @Router /teachers/{id} [get]
func (h *teacherHandler) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid teacher ID",
		})
	}

	teacher, err := h.teacherRepo.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Teacher not found",
		})
	}

	teacher.Password = ""
	return c.JSON(fiber.Map{
		"data": teacher,
	})
}

// GetTeacherByTeacherID godoc
// @Summary Get teacher by teacher ID
// @Description Retrieve a specific teacher by their teacher ID
// @Tags Teachers
// @Accept json
// @Produce json
// @Param teacherId path string true "Teacher ID (not database ID)"
// @Success 200 {object} map[string]interface{} "Teacher retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Teacher ID is required"
// @Failure 404 {object} map[string]interface{} "Teacher not found"
// @Router /teachers/teacher-id/{teacherId} [get]
func (h *teacherHandler) GetByTeacherID(c *fiber.Ctx) error {
	teacherID := c.Params("teacherId")
	if teacherID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Teacher ID is required",
		})
	}

	teacher, err := h.teacherRepo.GetByTeacherID(c.Context(), teacherID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Teacher not found",
		})
	}

	teacher.Password = ""
	return c.JSON(fiber.Map{
		"data": teacher,
	})
}

// GetAllTeachers godoc
// @Summary Get all teachers
// @Description Retrieve a paginated list of all teachers
// @Tags Teachers
// @Accept json
// @Produce json
// @Param limit query int false "Number of teachers to return (max 100)" default(10)
// @Param offset query int false "Number of teachers to skip" default(0)
// @Success 200 {object} map[string]interface{} "List of teachers retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /teachers [get]
func (h *teacherHandler) GetAll(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	teachers, err := h.teacherRepo.GetAll(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get teachers",
		})
	}

	for _, teacher := range teachers {
		teacher.Password = ""
	}

	return c.JSON(fiber.Map{
		"data":   teachers,
		"count":  len(teachers),
		"limit":  limit,
		"offset": offset,
	})
}

// UpdateTeacher godoc
// @Summary Update teacher
// @Description Update a teacher's information
// @Tags Teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher database ID"
// @Param teacher body models.Teacher true "Teacher data"
// @Success 200 {object} map[string]interface{} "Teacher updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /teachers/{id} [put]
func (h *teacherHandler) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid teacher ID",
		})
	}

	var teacher models.Teacher
	if err := c.BodyParser(&teacher); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	teacher.ID = uint(id)
	if err := h.teacherRepo.Update(c.Context(), &teacher); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update teacher",
		})
	}

	teacher.Password = ""
	return c.JSON(fiber.Map{
		"message": "Teacher updated successfully",
		"data":    teacher,
	})
}

// DeleteTeacher godoc
// @Summary Delete teacher
// @Description Delete a teacher from the system
// @Tags Teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher database ID"
// @Success 200 {object} map[string]interface{} "Teacher deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid teacher ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /teachers/{id} [delete]
func (h *teacherHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid teacher ID",
		})
	}

	if err := h.teacherRepo.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete teacher",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Teacher deleted successfully",
	})
}
