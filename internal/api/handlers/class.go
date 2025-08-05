package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelwp/student_attendance/internal/models"
	"github.com/michaelwp/student_attendance/internal/repository"
)

type classHandler struct {
	classRepo repository.ClassRepository
}

// NewClassHandler creates a new class handler
func NewClassHandler(classRepo repository.ClassRepository) ClassHandler {
	return &classHandler{
		classRepo: classRepo,
	}
}

// CreateClass godoc
// @Summary Create a new class
// @Description Create a new class in the system
// @Tags Classes
// @Accept json
// @Produce json
// @Param class body models.Class true "Class data"
// @Success 201 {object} map[string]interface{} "Class created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /classes [post]
func (h *classHandler) Create(c *fiber.Ctx) error {
	var class models.Class
	if err := c.BodyParser(&class); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.classRepo.Create(c.Context(), &class); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create class",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Class created successfully",
		"data":    class,
	})
}

// GetClassByID godoc
// @Summary Get class by ID
// @Description Retrieve a specific class by ID
// @Tags Classes
// @Accept json
// @Produce json
// @Param id path int true "Class ID"
// @Success 200 {object} map[string]interface{} "Class retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid class ID"
// @Failure 404 {object} map[string]interface{} "Class not found"
// @Router /classes/{id} [get]
func (h *classHandler) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid class ID",
		})
	}

	class, err := h.classRepo.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Class not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": class,
	})
}

// GetAllClasses godoc
// @Summary Get all classes
// @Description Retrieve a paginated list of all classes
// @Tags Classes
// @Accept json
// @Produce json
// @Param limit query int false "Number of classes to return (max 100)" default(10)
// @Param offset query int false "Number of classes to skip" default(0)
// @Success 200 {object} map[string]interface{} "List of classes retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /classes [get]
func (h *classHandler) GetAll(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	classes, err := h.classRepo.GetAll(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get classes",
		})
	}

	return c.JSON(fiber.Map{
		"data":   classes,
		"count":  len(classes),
		"limit":  limit,
		"offset": offset,
	})
}

// GetClassesByTeacher godoc
// @Summary Get classes by teacher
// @Description Retrieve all classes assigned to a specific teacher
// @Tags Classes
// @Accept json
// @Produce json
// @Param teacherId path string true "Teacher ID"
// @Success 200 {object} map[string]interface{} "Classes retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Teacher ID is required"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /classes/teacher-id/{teacherId} [get]
func (h *classHandler) GetByTeacher(c *fiber.Ctx) error {
	teacherID := c.Params("teacherId")
	if teacherID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Teacher ID is required",
		})
	}

	classes, err := h.classRepo.GetByTeacher(c.Context(), teacherID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get classes",
		})
	}

	return c.JSON(fiber.Map{
		"data":  classes,
		"count": len(classes),
	})
}

// UpdateClass godoc
// @Summary Update class
// @Description Update a class's information
// @Tags Classes
// @Accept json
// @Produce json
// @Param id path int true "Class ID"
// @Param class body models.Class true "Class data"
// @Success 200 {object} map[string]interface{} "Class updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /classes/{id} [put]
func (h *classHandler) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid class ID",
		})
	}

	var class models.Class
	if err := c.BodyParser(&class); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	class.ID = uint(id)
	if err := h.classRepo.Update(c.Context(), &class); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update class",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Class updated successfully",
		"data":    class,
	})
}

// DeleteClass godoc
// @Summary Delete class
// @Description Delete a class from the system
// @Tags Classes
// @Accept json
// @Produce json
// @Param id path int true "Class ID"
// @Success 200 {object} map[string]interface{} "Class deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid class ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /classes/{id} [delete]
func (h *classHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid class ID",
		})
	}

	if err := h.classRepo.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete class",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Class deleted successfully",
	})
}
