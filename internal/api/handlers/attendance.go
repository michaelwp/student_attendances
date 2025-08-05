package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelwp/student_attendance/internal/models"
	"github.com/michaelwp/student_attendance/internal/repository"
)

type attendanceHandler struct {
	attendanceRepo repository.AttendanceRepository
}

// NewAttendanceHandler creates a new attendance handler
func NewAttendanceHandler(attendanceRepo repository.AttendanceRepository) AttendanceHandler {
	return &attendanceHandler{
		attendanceRepo: attendanceRepo,
	}
}

// CreateAttendance godoc
// @Summary Create attendance record
// @Description Create a new attendance record for a student
// @Tags Attendances
// @Accept json
// @Produce json
// @Param attendance body models.Attendance true "Attendance data"
// @Success 201 {object} map[string]interface{} "Attendance record created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /attendances [post]
func (h *attendanceHandler) Create(c *fiber.Ctx) error {
	var attendance models.Attendance
	if err := c.BodyParser(&attendance); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.attendanceRepo.Create(c.Context(), &attendance); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create attendance",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Attendance created successfully",
		"data":    attendance,
	})
}

// GetAttendanceByID godoc
// @Summary Get attendance by ID
// @Description Retrieve a specific attendance record by ID
// @Tags Attendances
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {object} map[string]interface{} "Attendance record retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid attendance ID"
// @Failure 404 {object} map[string]interface{} "Attendance record not found"
// @Router /attendances/{id} [get]
func (h *attendanceHandler) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid attendance ID",
		})
	}

	attendance, err := h.attendanceRepo.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Attendance not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": attendance,
	})
}

// GetAttendancesByStudent godoc
// @Summary Get attendance by student
// @Description Retrieve attendance records for a specific student
// @Tags Attendances
// @Accept json
// @Produce json
// @Param studentId path string true "Student ID"
// @Param limit query int false "Number of records to return (max 100)" default(10)
// @Param offset query int false "Number of records to skip" default(0)
// @Success 200 {object} map[string]interface{} "Attendance records retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Student ID is required"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /attendances/student-id/{studentId} [get]
func (h *attendanceHandler) GetByStudent(c *fiber.Ctx) error {
	studentID := c.Params("studentId")
	if studentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student ID is required",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	attendances, err := h.attendanceRepo.GetByStudent(c.Context(), studentID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get attendances",
		})
	}

	return c.JSON(fiber.Map{
		"data":   attendances,
		"count":  len(attendances),
		"limit":  limit,
		"offset": offset,
	})
}

// GetAttendancesByClass godoc
// @Summary Get attendance by class
// @Description Retrieve attendance records for a specific class
// @Tags Attendances
// @Accept json
// @Produce json
// @Param classId path int true "Class ID"
// @Param limit query int false "Number of records to return (max 100)" default(10)
// @Param offset query int false "Number of records to skip" default(0)
// @Success 200 {object} map[string]interface{} "Attendance records retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid class ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /attendances/class-id/{classId} [get]
func (h *attendanceHandler) GetByClass(c *fiber.Ctx) error {
	classIDParam := c.Params("classId")
	classID, err := strconv.ParseUint(classIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid class ID",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	attendances, err := h.attendanceRepo.GetByClass(c.Context(), uint(classID), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get attendances",
		})
	}

	return c.JSON(fiber.Map{
		"data":   attendances,
		"count":  len(attendances),
		"limit":  limit,
		"offset": offset,
	})
}

// GetAttendancesByDateRange godoc
// @Summary Get attendance by date range
// @Description Retrieve attendance records within a specific date range
// @Tags Attendances
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)" Format(date)
// @Param end_date query string true "End date (YYYY-MM-DD)" Format(date)
// @Param limit query int false "Number of records to return (max 100)" default(10)
// @Param offset query int false "Number of records to skip" default(0)
// @Success 200 {object} map[string]interface{} "Attendance records retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid date format or missing parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /attendances/date-range [get]
func (h *attendanceHandler) GetByDateRange(c *fiber.Ctx) error {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "start_date and end_date are required",
		})
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid start_date format (use YYYY-MM-DD)",
		})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid end_date format (use YYYY-MM-DD)",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	attendances, err := h.attendanceRepo.GetByDateRange(c.Context(), startDate, endDate, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get attendances",
		})
	}

	return c.JSON(fiber.Map{
		"data":       attendances,
		"count":      len(attendances),
		"limit":      limit,
		"offset":     offset,
		"start_date": startDateStr,
		"end_date":   endDateStr,
	})
}

// UpdateAttendance godoc
// @Summary Update attendance record
// @Description Update an attendance record
// @Tags Attendances
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Param attendance body models.Attendance true "Attendance data"
// @Success 200 {object} map[string]interface{} "Attendance record updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /attendances/{id} [put]
func (h *attendanceHandler) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid attendance ID",
		})
	}

	var attendance models.Attendance
	if err := c.BodyParser(&attendance); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	attendance.ID = uint(id)
	if err := h.attendanceRepo.Update(c.Context(), &attendance); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update attendance",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Attendance updated successfully",
		"data":    attendance,
	})
}

// DeleteAttendance godoc
// @Summary Delete attendance record
// @Description Delete an attendance record
// @Tags Attendances
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {object} map[string]interface{} "Attendance record deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid attendance ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /attendances/{id} [delete]
func (h *attendanceHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid attendance ID",
		})
	}

	if err := h.attendanceRepo.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete attendance",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Attendance deleted successfully",
	})
}
