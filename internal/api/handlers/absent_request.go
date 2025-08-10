package handlers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelwp/student_attendance/internal/models"
	"github.com/michaelwp/student_attendance/internal/repository"
)

type absentRequestHandler struct {
	absentRequestRepo repository.AbsentRequestRepository
	studentRepo       repository.StudentRepository
}

// NewAbsentRequestHandler creates a new absent request handler
func NewAbsentRequestHandler(
	absentRequestRepo repository.AbsentRequestRepository,
	studentRepo repository.StudentRepository) AbsentRequestHandler {
	return &absentRequestHandler{
		absentRequestRepo: absentRequestRepo,
		studentRepo:       studentRepo,
	}
}

// CreateAbsentRequest godoc
// @Summary Create absent request
// @Description Create a new absence request
// @Tags Absent Requests
// @Accept json
// @Produce json
// @Param request body models.AbsentRequestCreate true "Absent request data"
// @Success 201 {object} map[string]interface{} "Absent request created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /absent-requests [post]
func (h *absentRequestHandler) Create(c *fiber.Ctx) error {
	var requestCreate models.AbsentRequestCreate
	if err := c.BodyParser(&requestCreate); err != nil {
		log.Println("error on create absent request:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_body",
			"error":         "Invalid request body",
		})
	}

	// Convert to AbsentRequest with proper date parsing
	request, err := requestCreate.ToAbsentRequest()
	if err != nil {
		log.Println("error on create absent request: invalid date format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_date_format",
			"error":         "Invalid date format. Use YYYY-MM-DD format.",
		})
	}

	studentID := c.Locals("userID")
	if studentID == nil {
		log.Println("error on create absent request: invalid admin id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_admin_id",
			"error":         "Invalid admin ID",
		})
	}

	studentIDUint, err := strconv.ParseUint(studentID.(string), 10, 32)
	if err != nil {
		log.Println("error on create absent request: invalid student id format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_admin_id",
			"error":         "Invalid admin ID format",
		})
	}

	student, err := h.studentRepo.GetByID(c.Context(), uint(studentIDUint))
	if err != nil {
		log.Println("error on create absent request: failed to get student:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_student",
			"error":         "Failed to get student",
		})
	}

	if student == nil {
		log.Println("error on create absent request: student not found")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.student_not_found",
			"error":         "Student not found",
		})
	}

	request.StudentID = student.StudentID
	request.ClassID = student.ClassesID
	request.Status = models.AbsentRequestStatusPending

	if err := h.absentRequestRepo.Create(c.Context(), request); err != nil {
		log.Println("error on create absent request:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_create_absent_request",
			"error":         "Failed to create absent request",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"translate_key": "success.absent_request_created",
		"message":       "Absent request created successfully",
		"data":          request,
	})
}

// GetAbsentRequestByID godoc
// @Summary Get absent request by ID
// @Description Retrieve a specific absent request by ID
// @Tags Absent Requests
// @Accept json
// @Produce json
// @Param id path int true "Absent request ID"
// @Success 200 {object} map[string]interface{} "Absent request retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid absent request ID"
// @Failure 404 {object} map[string]interface{} "Absent request not found"
// @Router /absent-requests/{id} [get]
func (h *absentRequestHandler) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Println("error on get absent request by id:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_absent_request_id",
			"error":         "Invalid absent request ID",
		})
	}

	request, err := h.absentRequestRepo.GetByID(c.Context(), uint(id))
	if err != nil {
		log.Println("error on get absent request by id:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.absent_request_not_found",
			"error":         "Absent request not found",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.absent_request_retrieved",
		"message":       "Absent request retrieved successfully",
		"data":          request,
	})
}

// GetAbsentRequestsByStudent godoc
// @Summary Get absent requests by student
// @Description Retrieve absent requests for a specific student
// @Tags Absent Requests
// @Accept json
// @Produce json
// @Param studentId path string true "Student ID"
// @Param limit query int false "Number of requests to return (max 100)" default(10)
// @Param offset query int false "Number of requests to skip" default(0)
// @Success 200 {object} map[string]interface{} "Absent requests retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Student ID is required"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /absent-requests/student-id/{studentId} [get]
func (h *absentRequestHandler) GetByStudent(c *fiber.Ctx) error {
	studentID := c.Params("studentId")
	if studentID == "" {
		log.Println("error on get absent requests by student: student ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.student_id_required",
			"error":         "Student ID is required",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	requests, err := h.absentRequestRepo.GetByStudent(c.Context(), studentID, limit, offset)
	if err != nil {
		log.Println("error on get absent requests by student:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_absent_requests",
			"error":         "Failed to get absent requests",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.absent_requests_retrieved",
		"message":       "Absent requests retrieved successfully",
		"data":          requests,
		"count":         len(requests),
		"limit":         limit,
		"offset":        offset,
	})
}

// GetAbsentRequestsByClass godoc
// @Summary Get absent requests by class
// @Description Retrieve absent requests for a specific class
// @Tags Absent Requests
// @Accept json
// @Produce json
// @Param classId path int true "Class ID"
// @Param limit query int false "Number of requests to return (max 100)" default(10)
// @Param offset query int false "Number of requests to skip" default(0)
// @Success 200 {object} map[string]interface{} "Absent requests retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid class ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /absent-requests/class-id/{classId} [get]
func (h *absentRequestHandler) GetByClass(c *fiber.Ctx) error {
	classIDParam := c.Params("classId")
	classID, err := strconv.ParseUint(classIDParam, 10, 32)
	if err != nil {
		log.Println("error on get absent requests by class:", err)
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

	requests, err := h.absentRequestRepo.GetByClass(c.Context(), uint(classID), limit, offset)
	if err != nil {
		log.Println("error on get absent requests by class:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_absent_requests",
			"error":         "Failed to get absent requests",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.absent_requests_retrieved",
		"message":       "Absent requests retrieved successfully",
		"data":          requests,
		"count":         len(requests),
		"limit":         limit,
		"offset":        offset,
	})
}

// GetPendingAbsentRequests godoc
// @Summary Get pending absent requests
// @Description Retrieve all pending absent requests
// @Tags Absent Requests
// @Accept json
// @Produce json
// @Param limit query int false "Number of requests to return (max 100)" default(10)
// @Param offset query int false "Number of requests to skip" default(0)
// @Success 200 {object} map[string]interface{} "Pending absent requests retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /absent-requests/pending [get]
func (h *absentRequestHandler) GetPending(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	requests, err := h.absentRequestRepo.GetPending(c.Context(), limit, offset)
	if err != nil {
		log.Println("error on get pending absent requests:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_pending_absent_requests",
			"error":         "Failed to get pending absent requests",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.pending_absent_requests_retrieved",
		"message":       "Pending absent requests retrieved successfully",
		"data":          requests,
		"count":         len(requests),
		"limit":         limit,
		"offset":        offset,
	})
}

// UpdateAbsentRequestStatus godoc
// @Summary Update absent request status
// @Description Update the status of an absent request (approve/reject)
// @Tags Absent Requests
// @Accept json
// @Produce json
// @Param id path int true "Absent request ID"
// @Param status body object true "Status update with status field"
// @Success 200 {object} map[string]interface{} "Absent request status updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or status value"
// @Failure 404 {object} map[string]interface{} "Absent request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /absent-requests/{id}/status [patch]
func (h *absentRequestHandler) UpdateStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Println("error on update absent request status:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_absent_request_id",
			"error":         "Invalid absent request ID",
		})
	}

	var statusUpdate struct {
		Status models.AbsentRequestStatus `json:"status"`
	}

	if err := c.BodyParser(&statusUpdate); err != nil {
		log.Println("error on update absent request status:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_body",
			"error":         "Invalid request body",
		})
	}

	// Validate status
	if statusUpdate.Status != models.AbsentRequestStatusApproved &&
		statusUpdate.Status != models.AbsentRequestStatusRejected &&
		statusUpdate.Status != models.AbsentRequestStatusPending {
		log.Println("error on update absent request status: invalid status value")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_status_value",
			"error":         "Invalid status value",
		})
	}

	// Get existing request
	request, err := h.absentRequestRepo.GetByID(c.Context(), uint(id))
	if err != nil {
		log.Println("error on update absent request status:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.absent_request_not_found",
			"error":         "Absent request not found",
		})
	}

	// Update status
	request.Status = statusUpdate.Status
	if err := h.absentRequestRepo.Update(c.Context(), request); err != nil {
		log.Println("error on update absent request status:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_update_absent_request_status",
			"error":         "Failed to update absent request status",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.absent_request_status_updated",
		"message":       "Absent request status updated successfully",
		"data":          request,
	})
}

// DeleteAbsentRequest godoc
// @Summary Delete absent request
// @Description Delete an absent request
// @Tags Absent Requests
// @Accept json
// @Produce json
// @Param id path int true "Absent request ID"
// @Success 200 {object} map[string]interface{} "Absent request deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid absent request ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /absent-requests/{id} [delete]
func (h *absentRequestHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Println("error on delete absent request:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_absent_request_id",
			"error":         "Invalid absent request ID",
		})
	}

	studentID := c.Locals("userID")
	if studentID == nil {
		log.Println("error on delete absent requests: invalid student id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_student_id",
			"error":         "Invalid student ID",
		})
	}

	studentIDUint, err := strconv.ParseUint(studentID.(string), 10, 64)
	if err != nil {
		log.Println("error on delete absent request:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_delete_absent_request",
			"error":         "Failed to delete absent request",
		})
	}

	if err := h.absentRequestRepo.UpdateDeleteInfo(c.Context(), uint(id), uint(studentIDUint), uint(studentIDUint)); err != nil {
		log.Println("error on delete absent request:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_delete_absent_request",
			"error":         "Failed to delete absent request",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.absent_request_deleted",
		"message":       "Absent request deleted successfully",
	})
}

// GetByCurrentStudent godoc
// @Summary Get absent requests for current student
// @Description Retrieve absent requests for the currently authenticated student
// @Tags Absent Requests
// @Accept json
// @Produce json
// @Param limit query int false "Number of requests to return (max 100)" default(10)
// @Param offset query int false "Number of requests to skip" default(0)
// @Success 200 {object} map[string]interface{} "Absent requests retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid student ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /absent-requests/current-student [get]
func (h *absentRequestHandler) GetByCurrentStudent(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		log.Println("error on get absent requests: invalid student id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_student_id",
			"error":         "Invalid student ID",
		})
	}

	userIDUint, err := strconv.ParseUint(userID.(string), 10, 64)
	if err != nil {
		log.Println("error on get absent requests: invalid student id")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_absent_requests",
			"error":         "Failed to get absent requests",
		})
	}

	student, err := h.absentRequestRepo.GetByID(c.Context(), uint(userIDUint))
	if err != nil {
		log.Println("error on get absent requests: failed to get student")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_absent_requests",
			"error":         "Failed to get absent requests",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "5"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	students, err := h.absentRequestRepo.GetByStudent(c.Context(), student.StudentID, limit, offset)
	if err != nil {
		log.Println("error on get absent requests:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_absent_requests",
			"error":         "Failed to get absent requests",
		})
	}

	totalAbsentRequests, err := h.absentRequestRepo.GetTotalAbsentRequests(c.Context())

	return c.JSON(fiber.Map{
		"translate_key": "success.absent_requests_retrieved",
		"message":       "Absent requests retrieved successfully",
		"data":          students,
		"total":         totalAbsentRequests,
		"limit":         limit,
		"offset":        offset,
	})
}

// UpdateByCurrentStudent godoc
// @Summary Update absent request by current student
// @Description Update an absent request's date and reason if it belongs to the current student and is pending
// @Tags Absent Requests
// @Accept json
// @Produce json
// @Param id path int true "Absent request ID"
// @Param request body models.AbsentRequestCreate true "Absent request update data"
// @Success 200 {object} map[string]interface{} "Absent request updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Absent request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /absent-requests/absent-request-id/{id} [put]
func (h *absentRequestHandler) UpdateByCurrentStudent(c *fiber.Ctx) error {
	studentIDLocal := c.Locals("userID")
	if studentIDLocal == nil {
		log.Println("error on update absent request: invalid student id")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.invalid_student_id",
			"error":         "Invalid student ID",
		})
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Println("error on update absent request: invalid id param")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_absent_request_id",
			"error":         "Invalid absent request ID",
		})
	}

	// Load current student to get StudentID string
	studentRecordID, _ := strconv.ParseUint(studentIDLocal.(string), 10, 32)
	student, err := h.studentRepo.GetByID(c.Context(), uint(studentRecordID))
	if err != nil || student == nil {
		log.Println("error on update absent request: failed to get student")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.student_not_found",
			"error":         "Student not found",
		})
	}

	existing, err := h.absentRequestRepo.GetByID(c.Context(), uint(id))
	if err != nil || existing == nil {
		log.Println("error on update absent request: not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.absent_request_not_found",
			"error":         "Absent request not found",
		})
	}

	// Ownership and status checks
	if existing.StudentID != student.StudentID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.unauthorized",
			"error":         "You are not allowed to update this request",
		})
	}
	if existing.Status != models.AbsentRequestStatusPending {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.absent_request_not_pending",
			"error":         "Only pending requests can be updated",
		})
	}

	var updateBody models.AbsentRequestCreate
	if err := c.BodyParser(&updateBody); err != nil {
		log.Println("error on update absent request: invalid body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_body",
			"error":         "Invalid request body",
		})
	}
	converted, err := updateBody.ToAbsentRequest()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_date_format",
			"error":         "Invalid date format. Use YYYY-MM-DD format.",
		})
	}

	// Apply updates
	existing.RequestDate = converted.RequestDate
	existing.Reason = converted.Reason

	if err := h.absentRequestRepo.Update(c.Context(), existing); err != nil {
		log.Println("error on update absent request: repo update failed", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_update_absent_request",
			"error":         "Failed to update absent request",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.absent_request_updated",
		"message":       "Absent request updated successfully",
		"data":          existing,
	})
}
