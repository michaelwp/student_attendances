package handlers

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/michaelwp/student_attendance/internal/config"
	"io"
	"log"
	"path/filepath"
	"strconv"
	"time"

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
func NewStudentHandler(studentRepo repository.StudentRepository, s3Client *s3.Client, s3Config *config.S3Config) StudentHandler {
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate.key": "error.invalid.request.body",
			"error":         "Invalid request body",
		})
	}

	if err := h.studentRepo.Create(c.Context(), &student); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate.key": "error.student.creation.failed",
			"error":         "Failed to create student",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"translate.key": "success.student.created.successfully",
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate.key": "error.invalid.student.id",
			"error":         "Invalid student ID",
		})
	}

	student, err := h.studentRepo.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate.key": "error.student.not.found",
			"error":         "Student not found",
		})
	}

	student.Password = ""
	return c.JSON(fiber.Map{
		"translate.key": "success.student.retrieved.successfully",
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate.key": "error.student.id.required",
			"error":         "Student ID is required",
		})
	}

	student, err := h.studentRepo.GetByStudentID(c.Context(), studentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate.key": "error.student.not.found",
			"error":         "Student not found",
		})
	}

	student.Password = ""
	return c.JSON(fiber.Map{
		"translate.key": "success.student.retrieved.successfully",
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate.key": "errror.students.retrieval.failed",
			"error":         "Failed to get students",
		})
	}

	for _, student := range students {
		student.Password = ""
	}

	totalStudents, err := h.studentRepo.GetTotalStudents(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate.key": "error.total.students.retrieval.failed",
			"error":         "Failed to get total students",
		})
	}

	return c.JSON(fiber.Map{
		"translate.key": "success.students.retrieved.successfully",
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate.key": "error.invalid.class.id",
			"error":         "Invalid class ID",
		})
	}

	students, err := h.studentRepo.GetByClass(c.Context(), uint(classID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate.key": "error.students.retrieval.failed",
			"error":         "Failed to get students",
		})
	}

	for _, student := range students {
		student.Password = ""
	}

	return c.JSON(fiber.Map{
		"translate.key": "success.students.retrieved.successfully",
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate.key": "error.invalid.student.id",
			"error":         "Invalid student ID",
		})
	}

	var student models.Student
	if err := c.BodyParser(&student); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate.key": "error.invalid.request.body",
			"error":         "Invalid request body",
		})
	}

	student.ID = uint(id)
	if err := h.studentRepo.Update(c.Context(), &student); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate.key": "error.student.update.failed",
			"error":         "Failed to update student",
		})
	}

	return c.JSON(fiber.Map{
		"translate.key": "success.student.updated",
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
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate.key": "error.invalid.student.id",
			"error":         "Invalid student ID",
		})
	}

	if err := h.studentRepo.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate.key": "error.student.deletion.failed",
			"error":         "Failed to delete student",
		})
	}

	return c.JSON(fiber.Map{
		"translate.key": "success.student.deleted",
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate.key": "error.invalid.student.id",
			"error":         "Invalid student ID",
		})
	}

	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate.key": "nerror.o.file.uploaded",
			"error":         "No file uploaded",
		})
	}

	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate.key": "error.failed.to.open.file",
			"error":         "Failed to read photo file",
		})
	}
	defer fileContent.Close()

	buffer, err := io.ReadAll(fileContent)
	if err != nil {
		log.Println("error on get buffer from file content:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate.key": "error.failed.to.read.file",
			"error":         "Failed to read file",
		})
	}

	filename := "student_" + idParam + "_" + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(file.Filename)
	key := fmt.Sprintf("photos/students/%d/%s", id, filename)

	if err := h.s3Config.UploadFile(h.s3Client, key, buffer); err != nil {
		log.Println("error on upload file to S3:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate.key": "error.failed.to.upload.file",
			"error":         "Failed to upload file to S3",
		})
	}

	photoPath := h.s3Config.GetObjectURL(key)
	if err := h.studentRepo.UpdatePhotoPath(c.Context(), uint(id), key); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate.key": "error.photo.update.failed",
			"error":         "Failed to update photo path",
		})
	}

	return c.JSON(fiber.Map{
		"translate.key": "success.photo.uploaded.successfully",
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate.key": "error.invalid.student.id",
			"error":         "Invalid student ID",
		})
	}

	photoPath, err := h.studentRepo.GetPhotoPath(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"translate.key": "error.photo.not.found",
			"error":         "Photo not found",
		})
	}

	expires := time.Minute * 15

	signedURL, err := h.s3Config.GetSignedURL(h.s3Client, photoPath, expires)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate.key": "error.failed.to.generate.photo.url",
			"error":         "Failed to generate photo URL",
		})
	}

	return c.JSON(fiber.Map{
		"translate.key": "success.photo.retrieved",
		"message":       "Photo URL retrieved successfully",
		"photoUrl":      signedURL,
	})
}
