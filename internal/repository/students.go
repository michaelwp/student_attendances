package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/michaelwp/student_attendance/internal/models"
)

type studentRepository struct {
	db *sql.DB
}

// NewStudentRepository creates a new student repository
func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Create(ctx context.Context, student *models.Student) error {
	query := `
		INSERT INTO students (student_id, classes_id, first_name, last_name, email, phone, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		student.StudentID,
		student.ClassesID,
		student.FirstName,
		student.LastName,
		student.Email,
		student.Phone,
		student.Password,
	).Scan(&student.ID, &student.CreatedAt, &student.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create student: %w", err)
	}

	return nil
}

func (r *studentRepository) GetByID(ctx context.Context, id uint) (*models.Student, error) {
	query := `
		SELECT id, student_id, classes_id, first_name, last_name, email, phone, password, created_at, updated_at, is_active
		FROM students WHERE id = $1 AND deleted_at IS NULL`

	student := &models.Student{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&student.ID,
		&student.StudentID,
		&student.ClassesID,
		&student.FirstName,
		&student.LastName,
		&student.Email,
		&student.Phone,
		&student.Password,
		&student.CreatedAt,
		&student.UpdatedAt,
		&student.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student not found")
		}
		return nil, fmt.Errorf("failed to get student: %w", err)
	}

	return student, nil
}

func (r *studentRepository) GetByStudentID(ctx context.Context, studentID string) (*models.Student, error) {
	query := `
		SELECT id, student_id, classes_id, first_name, last_name, email, phone, password, created_at, updated_at, is_active
		FROM students WHERE student_id = $1 AND deleted_at IS NULL`

	student := &models.Student{}
	err := r.db.QueryRowContext(ctx, query, studentID).Scan(
		&student.ID,
		&student.StudentID,
		&student.ClassesID,
		&student.FirstName,
		&student.LastName,
		&student.Email,
		&student.Phone,
		&student.Password,
		&student.CreatedAt,
		&student.UpdatedAt,
		&student.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student not found")
		}
		return nil, fmt.Errorf("failed to get student: %w", err)
	}

	return student, nil
}

func (r *studentRepository) GetByEmail(ctx context.Context, email string) (*models.Student, error) {
	query := `
		SELECT id, student_id, classes_id, first_name, last_name, email, phone, password, created_at, updated_at, is_active
		FROM students WHERE email = $1 AND deleted_at IS NULL`

	student := &models.Student{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&student.ID,
		&student.StudentID,
		&student.ClassesID,
		&student.FirstName,
		&student.LastName,
		&student.Email,
		&student.Phone,
		&student.Password,
		&student.CreatedAt,
		&student.UpdatedAt,
		&student.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student not found")
		}
		return nil, fmt.Errorf("failed to get student: %w", err)
	}

	return student, nil
}

func (r *studentRepository) GetByClass(ctx context.Context, classID uint) ([]*models.Student, error) {
	query := `
		SELECT id, student_id, classes_id, first_name, last_name, email, phone, password, created_at, updated_at, is_active
		FROM students 
		WHERE classes_id = $1 AND deleted_at IS NULL
		ORDER BY first_name, last_name`

	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, fmt.Errorf("failed to get students by class: %w", err)
	}
	defer rows.Close()

	var students []*models.Student
	for rows.Next() {
		student := &models.Student{}
		err := rows.Scan(
			&student.ID,
			&student.StudentID,
			&student.ClassesID,
			&student.FirstName,
			&student.LastName,
			&student.Email,
			&student.Phone,
			&student.Password,
			&student.CreatedAt,
			&student.UpdatedAt,
			&student.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan student: %w", err)
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate students: %w", err)
	}

	return students, nil
}

func (r *studentRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Student, error) {
	query := `
		SELECT id, student_id, classes_id, first_name, last_name, email, phone, password, created_at, updated_at, is_active
		FROM students
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get students: %w", err)
	}
	defer rows.Close()

	var students []*models.Student
	for rows.Next() {
		student := &models.Student{}
		err := rows.Scan(
			&student.ID,
			&student.StudentID,
			&student.ClassesID,
			&student.FirstName,
			&student.LastName,
			&student.Email,
			&student.Phone,
			&student.Password,
			&student.CreatedAt,
			&student.UpdatedAt,
			&student.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan student: %w", err)
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate students: %w", err)
	}

	return students, nil
}

func (r *studentRepository) Update(ctx context.Context, student *models.Student) error {
	query := `
		UPDATE students 
		SET student_id = $2, classes_id = $3, first_name = $4, last_name = $5, email = $6, phone = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		student.ID,
		student.StudentID,
		student.ClassesID,
		student.FirstName,
		student.LastName,
		student.Email,
		student.Phone,
	).Scan(&student.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("student not found")
		}
		return fmt.Errorf("failed to update student: %w", err)
	}

	return nil
}

func (r *studentRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM students WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("student not found")
	}

	return nil
}

func (r *studentRepository) UpdatePhotoPath(ctx context.Context, id uint, photoPath string) error {
	query := `
		UPDATE students 
		SET photo_path = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	var updatedAt string
	err := r.db.QueryRowContext(ctx, query, id, photoPath).Scan(&updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("student not found")
		}
		return fmt.Errorf("failed to update student photo path: %w", err)
	}

	return nil
}

func (r *studentRepository) GetPhotoPath(ctx context.Context, id uint) (string, error) {
	query := `SELECT photo_path FROM students WHERE id = $1`

	var photoPath string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&photoPath)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("student not found")
		}
		return "", fmt.Errorf("failed to get student photo path: %w", err)
	}

	return photoPath, nil
}

func (r *studentRepository) GetTotalStudents(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM students WHERE deleted_at IS NULL`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total students: %w", err)
	}

	return count, nil
}

func (r *studentRepository) UpdatePassword(ctx context.Context, studentID string, password string) error {
	query := `
		UPDATE students 
		SET password = $2, updated_at = NOW()
		WHERE student_id = $1
		RETURNING updated_at`

	var updatedAt string
	err := r.db.QueryRowContext(ctx, query, studentID, password).Scan(&updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("student not found")
		}
		return fmt.Errorf("failed to update student password: %w", err)
	}

	return nil
}

func (r *studentRepository) GetPasswordByStudentID(ctx context.Context, studentID string) (string, error) {
	query := `SELECT password FROM students WHERE student_id = $1`

	var password string
	err := r.db.QueryRowContext(ctx, query, studentID).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("student not found")
		}
		return "", fmt.Errorf("failed to get student password: %w", err)
	}

	return password, nil
}

func (r *studentRepository) IsStudentExist(ctx context.Context, studentID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM students WHERE student_id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, studentID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check student existence: %w", err)
	}

	return exists, nil
}

func (r *studentRepository) GetStats(ctx context.Context) (*models.StudentStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_students,
			COUNT(CASE WHEN is_active = true THEN 1 END) as active_students,
			COUNT(CASE WHEN is_active = false THEN 1 END) as inactive_students
		FROM students
		WHERE deleted_at IS NULL
		`

	stats := &models.StudentStats{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalStudents,
		&stats.ActiveStudents,
		&stats.InactiveStudents,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get students stats: %w", err)
	}

	return stats, nil
}

func (r *studentRepository) UpdateDeleteInfo(ctx context.Context, id uint, deletedBy uint) error {
	query := `
		UPDATE students
		SET deleted_at = NOW(), deleted_by = $2
		WHERE id = $1
		RETURNING deleted_at`

	var deletedAt string
	err := r.db.QueryRowContext(ctx, query, id, deletedBy).Scan(&deletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("student not found")
		}
		return fmt.Errorf("failed to update student delete info: %w", err)
	}

	return nil
}
