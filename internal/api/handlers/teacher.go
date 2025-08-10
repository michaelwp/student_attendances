package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/michaelwp/student_attendance/internal/config"
	"github.com/michaelwp/student_attendance/pkg"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelwp/student_attendance/internal/models"
	"github.com/michaelwp/student_attendance/internal/repository"
)

type teacherHandler struct {
	teacherRepo       repository.TeacherRepository
	s3Config          *config.S3Config
	s3Client          *s3.Client
	classRepo         repository.ClassRepository
	absentRequestRepo repository.AbsentRequestRepository
}

// NewTeacherHandler creates a new teacher handler
func NewTeacherHandler(
	teacherRepo repository.TeacherRepository,
	s3Client *s3.Client,
	s3Config *config.S3Config,
	classRepo repository.ClassRepository,
	absentRequestRepo repository.AbsentRequestRepository,
) TeacherHandler {
	return &teacherHandler{
		teacherRepo:       teacherRepo,
		s3Config:          s3Config,
		s3Client:          s3Client,
		classRepo:         classRepo,
		absentRequestRepo: absentRequestRepo,
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
		log.Println("error on parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_body",
			"error":         "Invalid request body",
		})
	}

	exist, err := h.teacherRepo.IsTeacherExist(c.Context(), teacher.TeacherID)
	if err != nil {
		log.Println("error on check teacher exist:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_check_teacher_exist",
			"error":         "Failed to check teacher exist",
		})
	}

	if exist {
		log.Println("error on create teacher: teacher already exist")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.teacher_already_exist",
			"error":         "Teacher already exist",
		})
	}

	round, _ := strconv.Atoi(os.Getenv("SALT"))
	hashPassword, err := pkg.HashPassword(teacher.Password, round)
	if err != nil {
		log.Println("error on hash password:", err)
		return err
	}

	teacher.Password = hashPassword

	if err := h.teacherRepo.Create(c.Context(), &teacher); err != nil {
		log.Println("error on create teacher:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_create_teacher",
			"error":         "Failed to create teacher",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"translate_key": "success.teacher_created",
		"message":       "Teacher created successfully",
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
		log.Println("error on parse teacher ID:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_teacher_id",
			"error":         "Invalid teacher ID",
		})
	}

	teacher, err := h.teacherRepo.GetByID(c.Context(), uint(id))
	if err != nil {
		log.Println("error on get teacher by ID:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.teacher_not_found",
			"error":         "Teacher not found",
		})
	}

	teacher.Password = ""
	return c.JSON(fiber.Map{
		"translate_key": "success.teacher_retrieved",
		"message":       "Teacher retrieved successfully",
		"data":          teacher,
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
		log.Println("error on get teacher by teacher ID: Teacher ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.teacher_id_required",
			"error":         "Teacher ID is required",
		})
	}

	teacher, err := h.teacherRepo.GetByTeacherID(c.Context(), teacherID)
	if err != nil {
		log.Println("error on get teacher by teacher ID:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.teacher_not_found",
			"error":         "Teacher not found",
		})
	}

	teacher.Password = ""
	return c.JSON(fiber.Map{
		"translate_key": "success.teacher_retrieved",
		"message":       "Teacher retrieved successfully",
		"data":          teacher,
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
		log.Println("error on get all teachers:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_teachers",
			"error":         "Failed to get teachers",
		})
	}

	for _, teacher := range teachers {
		teacher.Password = ""
	}

	totalTeachers, err := h.teacherRepo.GetTotalTeachers(c.Context())
	if err != nil {
		log.Println("error on get total teachers:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_total_teachers",
			"error":         "Failed to get total teachers",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.teachers_retrieved",
		"message":       "Teachers retrieved successfully",
		"data":          teachers,
		"count":         totalTeachers,
		"limit":         limit,
		"offset":        offset,
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
		log.Println("error on parse teacher ID:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_teacher_id",
			"error":         "Invalid teacher ID",
		})
	}

	var teacher models.Teacher
	if err := c.BodyParser(&teacher); err != nil {
		log.Println("error on parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_body",
			"error":         "Invalid request body",
		})
	}

	teacher.ID = uint(id)
	if err := h.teacherRepo.Update(c.Context(), &teacher); err != nil {
		log.Println("error on update teacher:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_update_teacher",
			"error":         "Failed to update teacher",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.teacher_updated",
		"message":       "Teacher updated successfully",
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
	currentUserID := c.Locals("userID")
	if currentUserID == nil {
		log.Println("error on delete teacher: current user ID is nil")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.unauthorized",
			"error":         "Unauthorized access",
		})
	}

	currentUserIDUint, err := strconv.ParseUint(currentUserID.(string), 10, 64)
	if err != nil {
		log.Println("error converting current user ID to int:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.invalid_current_user_id",
			"error":         "Invalid current user ID",
		})
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Println("error on parse teacher ID:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_teacher_id",
			"error":         "Invalid teacher ID",
		})
	}

	if err := h.teacherRepo.UpdateDeleteInfo(c.Context(), uint(id), uint(currentUserIDUint)); err != nil {
		log.Println("error on update delete info teacher:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_delete_teacher",
			"error":         "Failed to delete teacher",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.teacher_deleted",
		"message":       "Teacher deleted successfully",
	})
}

// UploadPhoto godoc
// @Summary Upload teacher photo
// @Description Upload a teacher's photo to S3 and update the photo path
// @Tags Teachers
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Teacher database ID"
// @Param photo formData file true "Teacher photo"
// @Success 200 {object} map[string]interface{} "Photo uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /teachers/{id}/photo [put]
func (h *teacherHandler) UploadPhoto(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Println("error on parse teacher ID:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_teacher_id",
			"error":         "Invalid teacher ID",
		})
	}

	file, err := c.FormFile("photo")
	if err != nil {
		log.Println("error on get file from form:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.no_photo_file_provided",
			"error":         "No photo file provided",
		})
	}

	fileContent, err := file.Open()
	if err != nil {
		log.Println("error on open file content:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_open_photo_file",
			"error":         "Failed to read photo file",
		})
	}
	defer fileContent.Close()

	buffer, err := io.ReadAll(fileContent)
	if err != nil {
		log.Println("error on get buffer from file content:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_read_file",
			"error":         "Failed to read file",
		})
	}

	filename := "teacher_" + idParam + "_" + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(file.Filename)
	key := fmt.Sprintf("photos/teachers/%d/%s", id, filename)

	if err := h.s3Config.UploadFile(h.s3Client, key, buffer); err != nil {
		log.Println("error on upload file to S3:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_upload_file",
			"error":         "Failed to upload file to S3",
		})
	}

	photoPath := h.s3Config.GetObjectURL(key)
	if err := h.teacherRepo.UpdatePhotoPath(c.Context(), uint(id), key); err != nil {
		log.Println("error on update photo path:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_update_photo_path",
			"error":         "Failed to update photo path",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.photo_uploaded",
		"message":       "Photo uploaded successfully",
		"path":          photoPath,
	})
}

// GetPhoto godoc
// @Summary Get teacher photo
// @Description Get teacher's photo signed URL
// @Tags Teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher database ID"
// @Success 200 {object} map[string]interface{} "Photo URL retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid teacher ID"
// @Failure 404 {object} map[string]interface{} "Photo not found"
// @Router /teachers/{id}/photo [get]
func (h *teacherHandler) GetPhoto(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Println("error on parse teacher ID:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_teacher_id",
			"error":         "Invalid teacher ID",
		})
	}

	photoPath, err := h.teacherRepo.GetPhotoPath(c.Context(), uint(id))
	if err != nil {
		log.Println("error on get photo path:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.photo_not_found",
			"error":         "Photo not found",
		})
	}

	signedURL, err := h.s3Config.GetSignedURL(h.s3Client, photoPath, time.Hour)
	if err != nil {
		log.Println("error on get signed URL:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_generate_signed_url",
			"error":         "Failed to generate signed URL",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.photo_url_retrieved",
		"message":       "Photo URL retrieved successfully",
		"url":           signedURL,
	})
}

// @Summary Reset teacher password
// @Description Reset a teacher's password
// @Tags Teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher ID"
// @Param password body string true "New password"
// @Success 200 {object} map[string]interface{} "Password reset successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /teachers/teacher-id/{teacherId}/reset-password [put]
func (h *teacherHandler) ResetPassword(c *fiber.Ctx) error {
	teacherID := c.Params("teacherId")
	if teacherID == "" {
		log.Println("error on reset password: Teacher ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.teacher.id.required",
			"error":         "Teacher ID is required",
		})
	}

	// Check if a teacher exists
	exist, err := h.teacherRepo.IsTeacherExist(c.Context(), teacherID)
	if err != nil {
		log.Println("error on check teacher existence:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.teacher.check.failed",
			"error":         "Failed to check teacher existence",
		})
	}

	if !exist {
		log.Println("error on reset password: Teacher not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.teacher.not.found",
			"error":         "Teacher not found",
		})
	}

	password, err := pkg.GeneratePassword(12)
	if err != nil {
		log.Println("error on generate password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.password.generation.failed",
			"error":         "Failed to generate password",
		})
	}

	err = h.updateCurrentPassword(c.Context(), teacherID, password)
	if err != nil {
		log.Println("error on update password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.password.update.failed",
			"error":         "Failed to update password",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.password.reset",
		"message":       "Password reset successfully",
		"newPassword":   password,
	})
}

// UpdatePassword godoc
// @Summary Update teacher password
// @Description Update a teacher's password with a new one
// @Tags Teachers
// @Accept json
// @Produce json
// @Param teacherId path string true "Teacher ID"
// @Param request body map[string]string true "Password update request"
// @Success 200 {object} map[string]interface{} "Password updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Teacher not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /teachers/teacher-id/{teacherId}/update-password [put]
func (h *teacherHandler) UpdatePassword(c *fiber.Ctx) error {
	teacherID := c.Params("teacherId")
	if teacherID == "" {
		log.Println("error on update password: Teacher ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.teacher.id.required",
			"error":         "Teacher ID is required",
		})
	}

	var request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("error on parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid.request.body",
			"error":         "Invalid request body",
		})
	}

	if request.NewPassword == "" || request.OldPassword == "" {
		log.Println("error on update password: Password is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.password.required",
			"error":         "Password is required",
		})
	}

	exist, err := h.teacherRepo.IsTeacherExist(c.Context(), teacherID)
	if err != nil {
		log.Println("error on check teacher existence:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.teacher.check.failed",
			"error":         "Failed to check teacher existence",
		})
	}

	if !exist {
		log.Println("error on update password: Teacher not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.teacher.not.found",
			"error":         "Teacher not found",
		})
	}

	storedPassword, err := h.teacherRepo.GetPasswordByTeacherID(c.Context(), teacherID)
	if err != nil {
		log.Println("error on retrieve stored password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.password.retrieval.failed",
			"error":         "Failed to retrieve stored password",
		})
	}

	if err := pkg.ComparePasswords(storedPassword, request.OldPassword); err != nil {
		log.Println("error on compare passwords:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid.old.password",
			"error":         "Invalid old password",
		})
	}

	err = h.updateCurrentPassword(c.Context(), teacherID, request.NewPassword)
	if err != nil {
		log.Println("error on update password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.password.update.failed",
			"error":         "Failed to update password",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.password.updated",
		"message":       "Password updated successfully",
	})
}

// GetStats godoc
// @Summary Get teacher statistics
// @Description Get various statistics about teachers in the system
// @Tags Teachers
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Statistics retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /teachers/stats [get]
func (h *teacherHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.teacherRepo.GetStats(c.Context())
	if err != nil {
		log.Println("error on get teacher stats:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_teacher_stats",
			"error":         "Failed to get teacher statistics",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.stats_retrieved",
		"message":       "Statistics retrieved successfully",
		"data":          stats,
	})
}

func (h *teacherHandler) updateCurrentPassword(ctx context.Context, teacherID string, password string) error {
	round, _ := strconv.Atoi(os.Getenv("SALT"))
	hashPassword, err := pkg.HashPassword(password, round)
	if err != nil {
		log.Println("error on hash password:", err)
		return err
	}

	if err := h.teacherRepo.UpdatePassword(ctx, teacherID, hashPassword); err != nil {
		log.Println("error on update password:", err)
		return err
	}

	return nil
}

// GetProfile godoc
// @Summary Get teacher profile
// @Description Get the profile of the currently authenticated teacher with classes and statistics
// @Tags Teacher Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Profile retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Teacher not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /teacher/profile [get]
func (h *teacherHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		log.Println("error on get profile: invalid user id")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.invalid_user_id",
			"error":         "Invalid user ID",
		})
	}

	userIDStr := userID.(string)
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Println("error on get profile: invalid user id format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_user_id_format",
			"error":         "Invalid user ID format",
		})
	}

	teacherWithClasses, err := h.teacherRepo.GetByIDWithClasses(c.Context(), uint(userIDUint))
	if err != nil {
		log.Println("error on get profile:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.teacher_not_found",
			"error":         "Teacher not found",
		})
	}

	if teacherWithClasses == nil || teacherWithClasses.Teacher == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.teacher_not_found",
			"error":         "Teacher not found",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.profile_retrieved",
		"message":       "Profile retrieved successfully",
		"data":          teacherWithClasses,
	})
}

// UpdateCurrentPassword godoc
// @Summary Update current teacher password
// @Description Update the currently authenticated teacher's password
// @Tags Teacher Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Password update request"
// @Success 200 {object} map[string]interface{} "Password updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Teacher not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /teacher/password [put]
func (h *teacherHandler) UpdateCurrentPassword(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		log.Println("error on update current password: invalid user id")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.invalid_user_id",
			"error":         "Invalid user ID",
		})
	}

	var request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("error on update current password: failed to parse request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid.request.body",
			"error":         "Invalid request body",
		})
	}

	if request.NewPassword == "" || request.OldPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.password.required",
			"error":         "Password is required",
		})
	}

	userIDStr := userID.(string)
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Println("error on update current password: invalid user id format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_user_id_format",
			"error":         "Invalid user ID format",
		})
	}

	teacher, err := h.teacherRepo.GetByID(c.Context(), uint(userIDUint))
	if err != nil || teacher == nil {
		log.Println("error on update current password: teacher not found or repo error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.teacher_not_found",
			"error":         "Teacher not found",
		})
	}

	storedPassword, err := h.teacherRepo.GetPasswordByTeacherID(c.Context(), teacher.TeacherID)
	if err != nil {
		log.Println("error on update current password: failed to retrieve stored password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.password.retrieval.failed",
			"error":         "Failed to retrieve stored password",
		})
	}

	if err := pkg.ComparePasswords(storedPassword, request.OldPassword); err != nil {
		log.Println("error on update current password: failed to compare passwords:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid.old.password",
			"error":         "Invalid old password",
		})
	}

	if err := h.updateCurrentPassword(c.Context(), teacher.TeacherID, request.NewPassword); err != nil {
		log.Println("error on update current password: failed to update current password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.password.update.failed",
			"error":         "Failed to update password",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.password.updated",
		"message":       "Password updated successfully",
	})
}

// GetAbsentRequests godoc
// @Summary Get absent requests for teacher
// @Description Get paginated list of absent requests from students in the teacher's classes
// @Tags Teacher Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of requests to return (max 100)" default(10)
// @Param offset query int false "Number of requests to skip" default(0)
// @Success 200 {object} map[string]interface{} "Absent requests retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Teacher not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /absent-requests/current-teacher [get]
func (h *teacherHandler) GetAbsentRequests(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		log.Println("error on get absent requests: invalid user id")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.invalid_user_id",
			"error":         "Invalid user ID",
		})
	}

	userIDStr := userID.(string)
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Println("error on get absent requests: invalid user id format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_user_id_format",
			"error":         "Invalid user ID format",
		})
	}

	// Get teacher to get teacherID
	teacher, err := h.teacherRepo.GetByID(c.Context(), uint(userIDUint))
	if err != nil || teacher == nil {
		log.Println("error on get absent requests: teacher not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.teacher_not_found",
			"error":         "Teacher not found",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	requests, err := h.absentRequestRepo.GetByTeacher(c.Context(), teacher.TeacherID, limit, offset)
	if err != nil {
		log.Println("error on get absent requests:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_absent_requests",
			"error":         "Failed to get absent requests",
		})
	}

	total, err := h.absentRequestRepo.GetCountByTeacher(c.Context(), teacher.TeacherID)
	if err != nil {
		log.Println("error on get absent requests count:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_absent_requests_count",
			"error":         "Failed to get absent requests count",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.absent_requests_retrieved",
		"message":       "Absent requests retrieved successfully",
		"data":          requests,
		"total":         total,
		"limit":         limit,
		"offset":        offset,
	})
}

// ApproveAbsentRequest godoc
// @Summary Approve absent request
// @Description Approve a student's absent request
// @Tags Teacher Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Absent request ID"
// @Success 200 {object} map[string]interface{} "Request approved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /absent-requests/absent-request-id/{id}/approve [put]
func (h *teacherHandler) ApproveAbsentRequest(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		log.Println("error on approve absent request: invalid user id")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.invalid_user_id",
			"error":         "Invalid user ID",
		})
	}

	userIDStr := userID.(string)
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Println("error on approve absent request: invalid user id format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_user_id_format",
			"error":         "Invalid user ID format",
		})
	}

	idParam := c.Params("id")
	requestID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Println("error on approve absent request: invalid request id:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_id",
			"error":         "Invalid request ID",
		})
	}

	if err := h.absentRequestRepo.Approve(c.Context(), uint(requestID), uint(userIDUint)); err != nil {
		log.Println("error on approve absent request:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_approve_request",
			"error":         "Failed to approve request",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.request_approved",
		"message":       "Request approved successfully",
	})
}

// RejectAbsentRequest godoc
// @Summary Reject absent request
// @Description Reject a student's absent request
// @Tags Teacher Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Absent request ID"
// @Success 200 {object} map[string]interface{} "Request rejected successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /absent-requests/absent-request-id/{id}/reject [put]
func (h *teacherHandler) RejectAbsentRequest(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		log.Println("error on reject absent request: invalid user id")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.invalid_user_id",
			"error":         "Invalid user ID",
		})
	}

	userIDStr := userID.(string)
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Println("error on reject absent request: invalid user id format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_user_id_format",
			"error":         "Invalid user ID format",
		})
	}

	idParam := c.Params("id")
	requestID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Println("error on reject absent request: invalid request id:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_id",
			"error":         "Invalid request ID",
		})
	}

	if err := h.absentRequestRepo.Reject(c.Context(), uint(requestID), uint(userIDUint)); err != nil {
		log.Println("error on reject absent request:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_reject_request",
			"error":         "Failed to reject request",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.request_rejected",
		"message":       "Request rejected successfully",
	})
}
