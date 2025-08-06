package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michaelwp/student_attendance/internal/models"
	"github.com/michaelwp/student_attendance/internal/repository"
	"github.com/michaelwp/student_attendance/pkg"
	"os"
	"strconv"
)

type adminHandler struct {
	adminRepo repository.AdminRepository
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(adminRepo repository.AdminRepository) AdminHandler {
	return &adminHandler{
		adminRepo: adminRepo,
	}
}

// CreateAdmin godoc
// @Summary Create a new admin
// @Description Create a new admin in the system
// @Tags Admins
// @Accept json
// @Produce json
// @Param admin body models.Admin true "Admin data"
// @Success 201 {object} map[string]interface{} "Admin created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admins [post]
func (h *adminHandler) Create(c *fiber.Ctx) error {
	var admin models.Admin
	if err := c.BodyParser(&admin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_body",
			"error":         "Invalid request body",
		})
	}

	if admin.Email == "" || admin.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.email_and_password_required",
			"error":         "Email and password are required",
		})
	}

	// Hash password before storing
	round, _ := strconv.Atoi(os.Getenv("ROUND"))
	hashedPassword, err := pkg.HashPassword(admin.Password, round)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.password_hashing_failed",
			"error":         "Failed to hash password",
		})
	}
	admin.Password = hashedPassword

	if err := h.adminRepo.Create(c.Context(), &admin); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_create_admin",
			"error":         "Failed to create admin",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"translate_key": "success.admin_created",
		"message":       "Admin created successfully",
	})
}

// GetAdminByID godoc
// @Summary Get admin by ID
// @Description Retrieve a specific admin by their database ID
// @Tags Admins
// @Accept json
// @Produce json
// @Param id path int true "Admin database ID"
// @Success 200 {object} map[string]interface{} "Admin retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid admin ID"
// @Failure 404 {object} map[string]interface{} "Admin not found"
// @Router /admins/{id} [get]
func (h *adminHandler) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_admin_id",
			"error":         "Invalid admin ID",
		})
	}

	admin, err := h.adminRepo.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.admin_not_found",
			"error":         "Admin not found",
		})
	}

	admin.Password = ""
	return c.JSON(fiber.Map{
		"translate_key": "success.admin_retrieved",
		"message":       "Admin retrieved successfully",
		"data":          admin,
	})
}

// GetAdminByEmail godoc
// @Summary Get admin by email
// @Description Retrieve a specific admin by their email
// @Tags Admins
// @Accept json
// @Produce json
// @Param email path string true "Admin email"
// @Success 200 {object} map[string]interface{} "Admin retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Email is required"
// @Failure 404 {object} map[string]interface{} "Admin not found"
// @Router /admins/email/{email} [get]
func (h *adminHandler) GetByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.email_required",
			"error":         "Email is required",
		})
	}

	admin, err := h.adminRepo.GetByEmail(c.Context(), email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.admin_not_found",
			"error":         "Admin not found",
		})
	}

	admin.Password = ""
	return c.JSON(fiber.Map{
		"translate_key": "success.admin_retrieved",
		"message":       "Admin retrieved successfully",
		"data":          admin,
	})
}

// GetAllAdmins godoc
// @Summary Get all admins
// @Description Retrieve a paginated list of all admins
// @Tags Admins
// @Accept json
// @Produce json
// @Param limit query int false "Number of admins to return (max 100)" default(10)
// @Param offset query int false "Number of admins to skip" default(0)
// @Success 200 {object} map[string]interface{} "List of admins retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admins [get]
func (h *adminHandler) GetAll(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	admins, err := h.adminRepo.GetAll(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_admins",
			"error":         "Failed to get admins",
		})
	}

	for _, admin := range admins {
		admin.Password = ""
	}

	totalAdmins, err := h.adminRepo.GetTotalAdmins(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_total_admins",
			"error":         "Failed to get total admins",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.admins_retrieved",
		"message":       "Admins retrieved successfully",
		"data":          admins,
		"total":         totalAdmins,
		"limit":         limit,
		"offset":        offset,
	})
}

// UpdateAdmin godoc
// @Summary Update admin
// @Description Update an admin's information
// @Tags Admins
// @Accept json
// @Produce json
// @Param id path int true "Admin database ID"
// @Param admin body models.Admin true "Admin data"
// @Success 200 {object} map[string]interface{} "Admin updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admins/{id} [put]
func (h *adminHandler) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_admin_id",
			"error":         "Invalid admin ID",
		})
	}

	var admin models.Admin
	if err := c.BodyParser(&admin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_body",
			"error":         "Invalid request body",
		})
	}

	admin.ID = uint(id)
	if err := h.adminRepo.Update(c.Context(), &admin); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_update_admin",
			"error":         "Failed to update admin",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.admin_updated",
		"message":       "Admin updated successfully",
		"data":          admin,
	})
}

// DeleteAdmin godoc
// @Summary Delete admin
// @Description Delete an admin from the system
// @Tags Admins
// @Accept json
// @Produce json
// @Param id path int true "Admin database ID"
// @Success 200 {object} map[string]interface{} "Admin deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid admin ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admins/{id} [delete]
func (h *adminHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_admin_id",
			"error":         "Invalid admin ID",
		})
	}

	if err := h.adminRepo.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_delete_admin",
			"error":         "Failed to delete admin",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.admin_deleted",
		"message":       "Admin deleted successfully",
	})
}

// UpdateAdminPassword godoc
// @Summary Update admin password
// @Description Update an admin's password
// @Tags Admins
// @Accept json
// @Produce json
// @Param id path int true "Admin database ID"
// @Param request body map[string]string true "Password update request"
// @Success 200 {object} map[string]interface{} "Password updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admins/{id}/password [put]
func (h *adminHandler) UpdatePassword(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.email_required",
			"error":         "Email is required",
		})
	}

	var request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_body",
			"error":         "Invalid request body",
		})
	}

	if request.NewPassword == "" || request.OldPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.password.required",
			"error":         "Password is required",
		})
	}

	exist, err := h.adminRepo.IsAdminExist(c.Context(), email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.admin.check.failed",
			"error":         "Failed to check admin existence",
		})
	}

	if !exist {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.admin.not.found",
			"error":         "Admin not found",
		})
	}

	storedPassword, err := h.adminRepo.GetPasswordByEmail(c.Context(), email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.password.retrieval.failed",
			"error":         "Failed to retrieve stored password",
		})
	}

	if err := pkg.ComparePasswords(storedPassword, request.OldPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid.old.password",
			"error":         "Invalid old password",
		})
	}

	// Hash password
	round, _ := strconv.Atoi(os.Getenv("ROUND"))
	hashedPassword, err := pkg.HashPassword(request.NewPassword, round)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.password_hashing_failed",
			"error":         "Failed to hash password",
		})
	}

	if err := h.adminRepo.UpdatePassword(c.Context(), email, hashedPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_update_password",
			"error":         "Failed to update password",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.password_updated",
		"message":       "Password updated successfully",
	})
}

// SetAdminActiveStatus godoc
// @Summary Set admin active status
// @Description Set an admin's active status (activate/deactivate)
// @Tags Admins
// @Accept json
// @Produce json
// @Param id path int true "Admin database ID"
// @Param request body map[string]bool true "Active status request"
// @Success 200 {object} map[string]interface{} "Admin status updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admins/{id}/status [put]
func (h *adminHandler) SetActiveStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_admin_id",
			"error":         "Invalid admin ID",
		})
	}

	var request struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_body",
			"error":         "Invalid request body",
		})
	}

	if err := h.adminRepo.SetActiveStatus(c.Context(), uint(id), request.IsActive); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_update_status",
			"error":         "Failed to update admin status",
		})
	}

	status := "deactivated"
	if request.IsActive {
		status = "activated"
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.admin_status_updated",
		"message":       "Admin " + status + " successfully",
	})
}
