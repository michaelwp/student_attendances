package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelwp/student_attendance/internal/models"
	"github.com/michaelwp/student_attendance/internal/repository"
	"github.com/michaelwp/student_attendance/pkg"
	"golang.org/x/crypto/bcrypt"
)

type attendanceHandler struct {
	attendanceRepo repository.AttendanceRepository
	studentRepo    repository.StudentRepository
}

// NewAttendanceHandler creates a new attendance handler
func NewAttendanceHandler(attendanceRepo repository.AttendanceRepository, studentRepo repository.StudentRepository) AttendanceHandler {
	return &attendanceHandler{
		attendanceRepo: attendanceRepo,
		studentRepo:    studentRepo,
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
		log.Println("Error parsing attendance body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_body",
			"error":         "Invalid request body",
		})
	}

	if err := h.attendanceRepo.Create(c.Context(), &attendance); err != nil {
		log.Println("Error creating attendance:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_create_attendance",
			"error":         "Failed to create attendance",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"translate_key": "success.attendance_created",
		"message":       "Attendance created successfully",
		"data":          attendance,
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
		log.Println("Error parsing attendance ID:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_attendance_id",
			"error":         "Invalid attendance ID",
		})
	}

	attendance, err := h.attendanceRepo.GetByID(c.Context(), uint(id))
	if err != nil {
		log.Println("Error getting attendance:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.attendance_not_found",
			"error":         "Attendance not found",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.attendance_retrieved",
		"message":       "Attendance retrieved successfully",
		"data":          attendance,
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
		log.Println("error parsing student ID:")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_student_id",
			"error":         "Student ID is required",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	attendances, err := h.attendanceRepo.GetByStudent(c.Context(), studentID, limit, offset)
	if err != nil {
		log.Println("Error getting attendances:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_attendances",
			"error":         "Failed to get attendances",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.attendances_retrieved",
		"message":       "Attendance records retrieved successfully",
		"data":          attendances,
		"count":         len(attendances),
		"limit":         limit,
		"offset":        offset,
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
		log.Println("Error parsing class ID:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_class_id",
			"error":         "Invalid class ID",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	attendances, err := h.attendanceRepo.GetByClass(c.Context(), uint(classID), limit, offset)
	if err != nil {
		log.Println("Error getting attendances:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_attendances",
			"error":         "Failed to get attendances",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.attendances_retrieved",
		"message":       "Attendance records retrieved successfully",
		"data":          attendances,
		"count":         len(attendances),
		"limit":         limit,
		"offset":        offset,
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
		log.Println("error parsing date range:")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_date_range",
			"error":         "start_date and end_date are required",
		})
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		log.Println("error parsing date range:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_date_range",
			"error":         "Invalid start_date format (use YYYY-MM-DD)",
		})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		log.Println("error parsing date range:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_date_range",
			"error":         "Invalid end_date format (use YYYY-MM-DD)",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	attendances, err := h.attendanceRepo.GetByDateRange(c.Context(), startDate, endDate, limit, offset)
	if err != nil {
		log.Println("Error getting attendances:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_attendances",
			"error":         "Failed to get attendances",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.attendances_retrieved",
		"message":       "Attendance records retrieved successfully",
		"data":          attendances,
		"count":         len(attendances),
		"limit":         limit,
		"offset":        offset,
		"start_date":    startDateStr,
		"end_date":      endDateStr,
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
		log.Println("Error parsing attendance ID:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_attendance_id",
			"error":         "Invalid attendance ID",
		})
	}

	var attendance models.Attendance
	if err := c.BodyParser(&attendance); err != nil {
		log.Println("Error parsing attendance:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_body",
			"error":         "Invalid request body",
		})
	}

	adminID := c.Locals("userID")
	if adminID == nil {
		log.Println("error: unauthorized access, userID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.unauthorized",
			"error":         "Unauthorized access",
		})
	}

	adminIDUint64, err := strconv.ParseUint(adminID.(string), 10, 64)
	if err != nil {
		log.Println("error on get current user id:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.invalid_admin_id",
			"error":         "Invalid admin ID",
		})
	}

	adminIDUint := uint(adminIDUint64)
	attendance.UpdatedBy = &adminIDUint

	attendance.ID = uint(id)
	if err := h.attendanceRepo.Update(c.Context(), &attendance); err != nil {
		log.Println("Error updating attendance:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_update_attendance",
			"error":         "Failed to update attendance",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.attendance_updated",
		"message":       "Attendance updated successfully",
		"data":          attendance,
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
		log.Println("Error parsing attendance ID:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_attendance_id",
			"error":         "Invalid attendance ID",
		})
	}

	adminID := c.Locals("userID")
	if adminID == nil {
		log.Println("error: unauthorized access, userID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.unauthorized",
			"error":         "Unauthorized access",
		})
	}

	adminIDUint, err := strconv.ParseUint(adminID.(string), 10, 64)
	if err != nil {
		log.Println("error on get current user id:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.invalid_admin_id",
			"error":         "Invalid admin ID",
		})
	}

	if err := h.attendanceRepo.UpdateDeleteInfo(c.Context(), uint(id), uint(adminIDUint)); err != nil {
		log.Println("Error deleting attendance:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_delete_attendance",
			"error":         "Failed to delete attendance",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.attendance_deleted",
		"message":       "Attendance deleted successfully",
	})
}

// MarkAttendance godoc
// @Summary Mark student attendance (Public endpoint)
// @Description Allow students to mark their own attendance using student ID and password
// @Tags Public
// @Accept json
// @Produce json
// @Param request body object{student_id=string,password=string} true "Student credentials"
// @Success 200 {object} object{student_name=string,message=string} "Attendance marked successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or missing parameters"
// @Failure 401 {object} map[string]interface{} "Invalid student credentials"
// @Failure 409 {object} map[string]interface{} "Attendance already marked today"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /attendance/mark [post]
func (h *attendanceHandler) MarkAttendance(c *fiber.Ctx) error {
	var request struct {
		StudentID string `json:"student_id"`
		Password  string `json:"password"`
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("Error parsing request:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_body",
			"error":         "Invalid request body",
		})
	}

	if request.StudentID == "" || request.Password == "" {
		log.Println("error parsing request body:")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.missing_credentials",
			"error":         "Student ID and password are required",
		})
	}

	// Get student by ID
	student, err := h.studentRepo.GetByStudentID(c.Context(), request.StudentID)
	if err != nil {
		log.Println("Error getting student:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.invalid_credentials",
			"error":         "Invalid student credentials",
		})
	}

	if student == nil {
		log.Println("Student not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.invalid_credentials",
			"error":         "Invalid student credentials",
		})
	}

	// Verify password
	hashedPassword, err := h.studentRepo.GetPasswordByStudentID(c.Context(), request.StudentID)
	if err != nil {
		log.Println("Error getting student password:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.invalid_credentials",
			"error":         "Invalid student credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(request.Password)); err != nil {
		log.Println("Error verifying password:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.invalid_credentials",
			"error":         "Invalid student credentials",
		})
	}

	// Check if the student is active
	if !student.IsActive {
		log.Println("Student account is inactive")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.account_inactive",
			"error":         "Student account is inactive",
		})
	}

	// Check if attendance already marked today
	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	existingAttendance, err := h.attendanceRepo.GetByStudentAndDate(c.Context(), request.StudentID, todayStart)
	if err != nil {
		log.Println("Error getting attendance:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_attendance",
			"error":         "Failed to get attendance",
		})
	}

	if existingAttendance != nil {
		log.Println("Attendance already marked today")
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"translate_key": "error.attendance_already_marked",
			"error":         "Attendance already marked for today",
		})
	}

	// Create an attendance record
	attendance := models.Attendance{
		StudentID:   request.StudentID,
		ClassID:     student.ClassesID,
		Date:        todayStart,
		Status:      "present",
		Description: pkg.StringPtr("Self-marked attendance"),
		CreatedBy:   student.ID,
	}

	if err := h.attendanceRepo.Create(c.Context(), &attendance); err != nil {
		log.Println("Error creating attendance:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_mark_attendance",
			"error":         "Failed to mark attendance",
		})
	}

	studentName := student.FirstName + " " + student.LastName

	return c.JSON(fiber.Map{
		"translate_key": "attendance.marked_successfully",
		"student_name":  studentName,
		"message":       "Attendance marked successfully",
	})
}

// GetAll godoc
// @Summary Get all attendance records
// @Description Retrieve all attendance records with pagination
// @Tags Attendances
// @Accept json
// @Produce json
// @Param limit query int false "Number of records to return (max 100)" default(10)
// @Param offset query int false "Number of records to skip" default(0)
// @Success 200 {object} map[string]interface{} "Attendance records retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /attendances [get]
func (h *attendanceHandler) GetAll(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	attendances, err := h.attendanceRepo.GetAll(c.Context(), limit, offset)
	if err != nil {
		log.Println("Error getting attendances:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_attendances",
			"error":         "Failed to get attendances",
		})
	}

	total, err := h.attendanceRepo.GetCount(c.Context())
	if err != nil {
		log.Println("Error getting attendance count:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_attendances",
			"error":         "Failed to get attendance count",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.attendances_retrieved",
		"message":       "Attendance records retrieved successfully",
		"data":          attendances,
		"total":         total,
		"limit":         limit,
		"offset":        offset,
	})
}
