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

type studentHandler struct {
	studentRepo repository.StudentRepository
	s3Config    *config.S3Config
	s3Client    *s3.Client
}

// NewStudentHandler creates a new student handler
func NewStudentHandler(
	studentRepo repository.StudentRepository,
	s3Client *s3.Client,
	s3Config *config.S3Config,
) StudentHandler {
	return &studentHandler{
		studentRepo: studentRepo,
		s3Client:    s3Client,
		s3Config:    s3Config,
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
		log.Println("error on parse student body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid.request.body",
			"error":         "Invalid request body",
		})
	}

	exist, err := h.studentRepo.IsStudentExist(c.Context(), student.StudentID)
	if err != nil {
		log.Println("error on check student exist:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_check_student_exist",
			"error":         "Failed to check student exist",
		})
	}

	if exist {
		log.Println("error on create student: student already exist")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.student_already_exist",
			"error":         "Student already exist",
		})
	}

	round, _ := strconv.Atoi(os.Getenv("SALT"))
	hashPassword, err := pkg.HashPassword(student.Password, round)
	if err != nil {
		log.Println("error on hash password:", err)
		return err
	}

	student.Password = hashPassword

	if err := h.studentRepo.Create(c.Context(), &student); err != nil {
		log.Println("error on create student:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.student.creation.failed",
			"error":         "Failed to create student",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"translate_key": "success.student.created.successfully",
		"message":       "Student created successfully",
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
		log.Println("error on parse student id:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid.student.id",
			"error":         "Invalid student ID",
		})
	}

	student, err := h.studentRepo.GetByID(c.Context(), uint(id))
	if err != nil {
		log.Println("error on get student by id:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.student.not.found",
			"error":         "Student not found",
		})
	}

	student.Password = ""
	return c.JSON(fiber.Map{
		"translate_key": "success.student.retrieved.successfully",
		"message":       "Student retrieved successfully",
		"data":          student,
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
		log.Println("error on get student by student id: student ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.student.id.required",
			"error":         "Student ID is required",
		})
	}

	student, err := h.studentRepo.GetByStudentID(c.Context(), studentID)
	if err != nil {
		log.Println("error on get student by student id:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.student.not.found",
			"error":         "Student not found",
		})
	}

	student.Password = ""
	return c.JSON(fiber.Map{
		"translate_key": "success.student.retrieved.successfully",
		"message":       "Student retrieved successfully",
		"data":          student,
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
		log.Println("error on get all students:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "errror.students.retrieval.failed",
			"error":         "Failed to get students",
		})
	}

	for _, student := range students {
		student.Password = ""
	}

	totalStudents, err := h.studentRepo.GetTotalStudents(c.Context())
	if err != nil {
		log.Println("error on get total students:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.total.students.retrieval.failed",
			"error":         "Failed to get total students",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.students.retrieved.successfully",
		"message":       "Students retrieved successfully",
		"data":          students,
		"total":         totalStudents,
		"limit":         limit,
		"offset":        offset,
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
		log.Println("error on get students by class id:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid.class.id",
			"error":         "Invalid class ID",
		})
	}

	students, err := h.studentRepo.GetByClass(c.Context(), uint(classID))
	if err != nil {
		log.Println("error on get students by class id:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.students.retrieval.failed",
			"error":         "Failed to get students",
		})
	}

	for _, student := range students {
		student.Password = ""
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.students.retrieved.successfully",
		"message":       "Students retrieved successfully",
		"data":          students,
		"total":         len(students),
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
		log.Println("error on update student:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid.student.id",
			"error":         "Invalid student ID",
		})
	}

	var student models.Student
	if err := c.BodyParser(&student); err != nil {
		log.Println("error on update student:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid.request.body",
			"error":         "Invalid request body",
		})
	}

	student.ID = uint(id)
	if err := h.studentRepo.Update(c.Context(), &student); err != nil {
		log.Println("error on update student:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.student.update.failed",
			"error":         "Failed to update student",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.student.updated",
		"message":       "Student updated successfully",
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
	currentUserID := c.Locals("userID")
	if currentUserID == nil {
		log.Println("error: unauthorized access, userID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.unauthorized",
			"error":         "Unauthorized access",
		})
	}

	currentUserIDUint, err := strconv.ParseUint(currentUserID.(string), 10, 64)
	if err != nil {
		log.Println("error on get current user id:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.invalid_admin_id",
			"error":         "Invalid admin ID",
		})
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Println("error on delete student:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid.student.id",
			"error":         "Invalid student ID",
		})
	}

	if err := h.studentRepo.UpdateDeleteInfo(c.Context(), uint(id), uint(currentUserIDUint)); err != nil {
		log.Println("error on update delete info student:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.student.deletion.failed",
			"error":         "Failed to delete student",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.student.deleted",
		"message":       "Student deleted successfully",
	})
}

// UploadPhoto godoc
// @Summary Upload student photo
// @Description Upload a photo for a student and update their photo path
// @Tags Students
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Student ID"
// @Param photo formData file true "Student photo"
// @Success 200 {object} map[string]interface{} "Photo uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /students/{id}/photo [put]
func (h *studentHandler) UploadPhoto(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Println("error on upload photo:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid.student.id",
			"error":         "Invalid student ID",
		})
	}

	file, err := c.FormFile("photo")
	if err != nil {
		log.Println("error on upload photo:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "nerror.o.file.uploaded",
			"error":         "No file uploaded",
		})
	}

	fileContent, err := file.Open()
	if err != nil {
		log.Println("error on upload photo:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed.to.open.file",
			"error":         "Failed to read photo file",
		})
	}
	defer fileContent.Close()

	buffer, err := io.ReadAll(fileContent)
	if err != nil {
		log.Println("error on get buffer from file content:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed.to.read.file",
			"error":         "Failed to read file",
		})
	}

	filename := "student_" + idParam + "_" + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(file.Filename)
	key := fmt.Sprintf("photos/students/%d/%s", id, filename)

	if err := h.s3Config.UploadFile(h.s3Client, key, buffer); err != nil {
		log.Println("error on upload file to S3:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed.to.upload.file",
			"error":         "Failed to upload file to S3",
		})
	}

	photoPath := h.s3Config.GetObjectURL(key)
	if err := h.studentRepo.UpdatePhotoPath(c.Context(), uint(id), key); err != nil {
		log.Println("error on update photo path:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.photo.update.failed",
			"error":         "Failed to update photo path",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.photo.uploaded.successfully",
		"message":       "Photo uploaded successfully",
		"photoPath":     photoPath,
	})
}

// GetPhoto godoc
// @Summary Get student photo
// @Description Get the signed URL for a student's photo
// @Tags Students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} map[string]interface{} "Photo URL retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid student ID"
// @Failure 404 {object} map[string]interface{} "Photo not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /students/{id}/photo [get]
func (h *studentHandler) GetPhoto(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Println("error on get photo:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid.student.id",
			"error":         "Invalid student ID",
		})
	}

	photoPath, err := h.studentRepo.GetPhotoPath(c.Context(), uint(id))
	if err != nil {
		log.Println("error on get photo path:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.photo.not.found",
			"error":         "Photo not found",
		})
	}

	expires := time.Minute * 15

	signedURL, err := h.s3Config.GetSignedURL(h.s3Client, photoPath, expires)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed.to.generate.photo.url",
			"error":         "Failed to generate photo URL",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.photo.retrieved",
		"message":       "Photo URL retrieved successfully",
		"photoUrl":      signedURL,
	})
}

// ResetPassword godoc
// @Summary Reset student password
// @Description Reset a student's password
// @Tags Students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param password body string true "New password"
// @Success 200 {object} map[string]interface{} "Password reset successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /students/student-id/{studentId}/reset-password [put]
func (h *studentHandler) ResetPassword(c *fiber.Ctx) error {
	studentID := c.Params("studentId")
	if studentID == "" {
		log.Println("error on reset password: student ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.student.id.required",
			"error":         "Student ID is required",
		})
	}

	// Check if a student exists
	exist, err := h.studentRepo.IsStudentExist(c.Context(), studentID)
	if err != nil {
		log.Println("error on reset password: failed to check student existence:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.student.check.failed",
			"error":         "Failed to check student existence",
		})
	}

	if !exist {
		log.Println("error on reset password: student not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.student.not.found",
			"error":         "Student not found",
		})
	}

	password, err := pkg.GeneratePassword(12)
	if err != nil {
		log.Println("error on reset password: failed to generate password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.password.generation.failed",
			"error":         "Failed to generate password",
		})
	}

	err = h.updateCurrentPassword(c.Context(), studentID, password)
	if err != nil {
		log.Println("error on reset password: failed to update current password:", err)
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
// @Summary Update student password
// @Description Update a student's password with a new one
// @Tags Students
// @Accept json
// @Produce json
// @Param studentId path string true "Student ID"
// @Param request body map[string]string true "Password update request"
// @Success 200 {object} map[string]interface{} "Password updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Student not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /students/student-id/{studentId}/update-password [put]
func (h *studentHandler) UpdatePassword(c *fiber.Ctx) error {
	studentID := c.Params("studentId")
	if studentID == "" {
		log.Println("error on update password: student ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.student.id.required",
			"error":         "Student ID is required",
		})
	}

	var request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("error on update password: failed to parse request body:", err)
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

	exist, err := h.studentRepo.IsStudentExist(c.Context(), studentID)
	if err != nil {
		log.Println("error on update password: failed to check student existence:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.student.check.failed",
			"error":         "Failed to check student existence",
		})
	}

	if !exist {
		log.Println("error on update password: student not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate_key": "error.student.not.found",
			"error":         "Student not found",
		})
	}

	storedPassword, err := h.studentRepo.GetPasswordByStudentID(c.Context(), studentID)
	if err != nil {
		log.Println("error on update password: failed to retrieve stored password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.password.retrieval.failed",
			"error":         "Failed to retrieve stored password",
		})
	}

	if err := pkg.ComparePasswords(storedPassword, request.OldPassword); err != nil {
		log.Println("error on update password: failed to compare passwords:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid.old.password",
			"error":         "Invalid old password",
		})
	}

	err = h.updateCurrentPassword(c.Context(), studentID, request.NewPassword)
	if err != nil {
		log.Println("error on update password: failed to update current password:", err)
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
// @Summary Get student statistics
// @Description Get statistics about students including total active and inactive counts
// @Tags Students
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Statistics retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /students/stats [get]
func (h *studentHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.studentRepo.GetStats(c.Context())
	if err != nil {
		log.Println("error on get stats:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.failed_to_get_student_stats",
			"error":         "Failed to get student statistics",
		})
	}

	return c.JSON(fiber.Map{
		"translate_key": "success.stats_retrieved",
		"message":       "Statistics retrieved successfully",
		"data":          stats,
	})
}

func (h *studentHandler) updateCurrentPassword(ctx context.Context, studentID string, password string) error {
	round, _ := strconv.Atoi(os.Getenv("SALT"))
	hashPassword, err := pkg.HashPassword(password, round)
	if err != nil {
		log.Println("error on hash password:", err)
		return err
	}

	if err := h.studentRepo.UpdatePassword(ctx, studentID, hashPassword); err != nil {
		log.Println("error on update password:", err)
		return err
	}

	return nil
}
