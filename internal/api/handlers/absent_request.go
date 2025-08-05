package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelwp/student_attendance/internal/models"
	"github.com/michaelwp/student_attendance/internal/repository"
)

type absentRequestHandler struct {
	absentRequestRepo repository.AbsentRequestRepository
}

// NewAbsentRequestHandler creates a new absent request handler
func NewAbsentRequestHandler(absentRequestRepo repository.AbsentRequestRepository) AbsentRequestHandler {
	return &absentRequestHandler{
		absentRequestRepo: absentRequestRepo,
	}
}

// CreateAbsentRequest godoc
// @Summary Create absent request
// @Description Create a new absence request
// @Tags Absent Requests
// @Accept json
// @Produce json
// @Param request body models.AbsentRequest true "Absent request data"
// @Success 201 {object} map[string]interface{} "Absent request created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /absent-requests [post]
func (h *absentRequestHandler) Create(c *fiber.Ctx) error {
	var request models.AbsentRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Default status to pending if not provided
	if request.Status == "" {
		request.Status = models.AbsentRequestStatusPending
	}

	if err := h.absentRequestRepo.Create(c.Context(), &request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create absent request",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Absent request created successfully",
		"data":    request,
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid absent request ID",
		})
	}

	request, err := h.absentRequestRepo.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Absent request not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": request,
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student ID is required",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	requests, err := h.absentRequestRepo.GetByStudent(c.Context(), studentID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get absent requests",
		})
	}

	return c.JSON(fiber.Map{
		"data":   requests,
		"count":  len(requests),
		"limit":  limit,
		"offset": offset,
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid class ID",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	requests, err := h.absentRequestRepo.GetByClass(c.Context(), uint(classID), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get absent requests",
		})
	}

	return c.JSON(fiber.Map{
		"data":   requests,
		"count":  len(requests),
		"limit":  limit,
		"offset": offset,
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get pending absent requests",
		})
	}

	return c.JSON(fiber.Map{
		"data":   requests,
		"count":  len(requests),
		"limit":  limit,
		"offset": offset,
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid absent request ID",
		})
	}

	var statusUpdate struct {
		Status models.AbsentRequestStatus `json:"status"`
	}

	if err := c.BodyParser(&statusUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate status
	if statusUpdate.Status != models.AbsentRequestStatusApproved &&
		statusUpdate.Status != models.AbsentRequestStatusRejected &&
		statusUpdate.Status != models.AbsentRequestStatusPending {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status value",
		})
	}

	// Get existing request
	request, err := h.absentRequestRepo.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Absent request not found",
		})
	}

	// Update status
	request.Status = statusUpdate.Status
	if err := h.absentRequestRepo.Update(c.Context(), request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update absent request status",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Absent request status updated successfully",
		"data":    request,
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid absent request ID",
		})
	}

	if err := h.absentRequestRepo.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete absent request",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Absent request deleted successfully",
	})
}
