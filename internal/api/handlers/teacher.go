package handlers

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/michaelwp/student_attendance/internal/config"
	"github.com/michaelwp/student_attendance/pkg"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelwp/student_attendance/internal/models"
	"github.com/michaelwp/student_attendance/internal/repository"
)

type teacherHandler struct {
	teacherRepo repository.TeacherRepository
	s3Config    *config.S3Config
	s3Client    *s3.Client
}

// NewTeacherHandler creates a new teacher handler
func NewTeacherHandler(teacherRepo repository.TeacherRepository, s3Client *s3.Client, s3Config *config.S3Config) TeacherHandler {
	return &teacherHandler{
		teacherRepo: teacherRepo,
		s3Config:    s3Config,
		s3Client:    s3Client,
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
			"translate_key": "error.invalid_request_body",
			"error":         "Invalid request body",
		})
	}

	if err := h.teacherRepo.Create(c.Context(), &teacher); err != nil {
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_teacher_id",
			"error":         "Invalid teacher ID",
		})
	}

	teacher, err := h.teacherRepo.GetByID(c.Context(), uint(id))
	if err != nil {
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.teacher_id_required",
			"error":         "Teacher ID is required",
		})
	}

	teacher, err := h.teacherRepo.GetByTeacherID(c.Context(), teacherID)
	if err != nil {
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_teacher_id",
			"error":         "Invalid teacher ID",
		})
	}

	var teacher models.Teacher
	if err := c.BodyParser(&teacher); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_body",
			"error":         "Invalid request body",
		})
	}

	teacher.ID = uint(id)
	if err := h.teacherRepo.Update(c.Context(), &teacher); err != nil {
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
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_teacher_id",
			"error":         "Invalid teacher ID",
		})
	}

	if err := h.teacherRepo.Delete(c.Context(), uint(id)); err != nil {
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_teacher_id",
			"error":         "Invalid teacher ID",
		})
	}

	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.no_photo_file_provided",
			"error":         "No photo file provided",
		})
	}

	fileContent, err := file.Open()
	if err != nil {
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_teacher_id",
			"error":         "Invalid teacher ID",
		})
	}

	photoPath, err := h.teacherRepo.GetPhotoPath(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.photo_not_found",
			"error":         "Photo not found",
		})
	}

	signedURL, err := h.s3Config.GetSignedURL(h.s3Client, photoPath, time.Hour)
	if err != nil {
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

// ResetPassword godoc
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.teacher.id.required",
			"error":         "Teacher ID is required",
		})
	}

	// Check if a teacher exists
	exist, err := h.teacherRepo.IsTeacherExist(c.Context(), teacherID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.teacher.check.failed",
			"error":         "Failed to check teacher existence",
		})
	}

	if !exist {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.teacher.not.found",
			"error":         "Teacher not found",
		})
	}

	password, err := pkg.GeneratePassword(12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.password.generation.failed",
			"error":         "Failed to generate password",
		})
	}

	err = h.updateCurrentPassword(c.Context(), teacherID, password)
	if err != nil {
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

	exist, err := h.teacherRepo.IsTeacherExist(c.Context(), teacherID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.teacher.check.failed",
			"error":         "Failed to check teacher existence",
		})
	}

	if !exist {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.teacher.not.found",
			"error":         "Teacher not found",
		})
	}

	storedPassword, err := h.teacherRepo.GetPasswordByTeacherID(c.Context(), teacherID)
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

	err = h.updateCurrentPassword(c.Context(), teacherID, request.NewPassword)
	if err != nil {
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
