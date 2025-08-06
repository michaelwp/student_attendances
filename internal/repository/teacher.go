package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/michaelwp/student_attendance/internal/models"
)

type teacherRepository struct {
	db *sql.DB
}

// NewTeacherRepository creates a new teacher repository
func NewTeacherRepository(db *sql.DB) TeacherRepository {
	return &teacherRepository{db: db}
}

func (r *teacherRepository) Create(ctx context.Context, teacher *models.Teacher) error {
	query := `
		INSERT INTO teachers (teacher_id, first_name, last_name, email, phone, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		teacher.TeacherID,
		teacher.FirstName,
		teacher.LastName,
		teacher.Email,
		teacher.Phone,
		teacher.Password,
	).Scan(&teacher.ID, &teacher.CreatedAt, &teacher.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create teacher: %w", err)
	}

	return nil
}

func (r *teacherRepository) GetByID(ctx context.Context, id uint) (*models.Teacher, error) {
	query := `
		SELECT id, teacher_id, first_name, last_name, email, phone, password, created_at, updated_at
		FROM teachers WHERE id = $1`

	teacher := &models.Teacher{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&teacher.ID,
		&teacher.TeacherID,
		&teacher.FirstName,
		&teacher.LastName,
		&teacher.Email,
		&teacher.Phone,
		&teacher.Password,
		&teacher.CreatedAt,
		&teacher.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("teacher not found")
		}
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}

	return teacher, nil
}

func (r *teacherRepository) GetByTeacherID(ctx context.Context, teacherID string) (*models.Teacher, error) {
	query := `
		SELECT id, teacher_id, first_name, last_name, email, phone, password, created_at, updated_at
		FROM teachers WHERE teacher_id = $1`

	teacher := &models.Teacher{}
	err := r.db.QueryRowContext(ctx, query, teacherID).Scan(
		&teacher.ID,
		&teacher.TeacherID,
		&teacher.FirstName,
		&teacher.LastName,
		&teacher.Email,
		&teacher.Phone,
		&teacher.Password,
		&teacher.CreatedAt,
		&teacher.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("teacher not found")
		}
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}

	return teacher, nil
}

func (r *teacherRepository) GetByEmail(ctx context.Context, email string) (*models.Teacher, error) {
	query := `
		SELECT id, teacher_id, first_name, last_name, email, phone, password, created_at, updated_at
		FROM teachers WHERE email = $1`

	teacher := &models.Teacher{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&teacher.ID,
		&teacher.TeacherID,
		&teacher.FirstName,
		&teacher.LastName,
		&teacher.Email,
		&teacher.Phone,
		&teacher.Password,
		&teacher.CreatedAt,
		&teacher.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("teacher not found")
		}
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}

	return teacher, nil
}

func (r *teacherRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Teacher, error) {
	query := `
		SELECT id, teacher_id, first_name, last_name, email, phone, password, created_at, updated_at
		FROM teachers
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get teachers: %w", err)
	}
	defer rows.Close()

	var teachers []*models.Teacher
	for rows.Next() {
		teacher := &models.Teacher{}
		err := rows.Scan(
			&teacher.ID,
			&teacher.TeacherID,
			&teacher.FirstName,
			&teacher.LastName,
			&teacher.Email,
			&teacher.Phone,
			&teacher.Password,
			&teacher.CreatedAt,
			&teacher.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan teacher: %w", err)
		}
		teachers = append(teachers, teacher)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate teachers: %w", err)
	}

	return teachers, nil
}

func (r *teacherRepository) Update(ctx context.Context, teacher *models.Teacher) error {
	query := `
		UPDATE teachers 
		SET teacher_id = $2, first_name = $3, last_name = $4, email = $5, phone = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		teacher.ID,
		teacher.TeacherID,
		teacher.FirstName,
		teacher.LastName,
		teacher.Email,
		teacher.Phone,
	).Scan(&teacher.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("teacher not found")
		}
		return fmt.Errorf("failed to update teacher: %w", err)
	}

	return nil
}

func (r *teacherRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM teachers WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete teacher: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("teacher not found")
	}

	return nil
}

func (r *teacherRepository) UpdatePhotoPath(ctx context.Context, id uint, photoPath string) error {
	query := `
		UPDATE teachers 
		SET photo_path = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	var updatedAt string
	err := r.db.QueryRowContext(ctx, query, id, photoPath).Scan(&updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("teacher not found")
		}
		return fmt.Errorf("failed to update teacher photo path: %w", err)
	}

	return nil
}

func (r *teacherRepository) GetPhotoPath(ctx context.Context, id uint) (string, error) {
	query := `SELECT photo_path FROM teachers WHERE id = $1`

	var photoPath string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&photoPath)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("teacher not found")
		}
		return "", fmt.Errorf("failed to get teacher photo path: %w", err)
	}

	return photoPath, nil
}

func (r *teacherRepository) GetTotalTeachers(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM teachers`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total teachers: %w", err)
	}

	return count, nil
}

func (r *teacherRepository) UpdatePassword(ctx context.Context, teacherID string, password string) error {
	query := `
		UPDATE teachers
		SET password = $2, updated_at = NOW()
		WHERE teacher_id = $1
		RETURNING updated_at`

	var updatedAt string
	err := r.db.QueryRowContext(ctx, query, teacherID, password).Scan(&updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("teacher not found")
		}
		return fmt.Errorf("failed to update teacher password: %w", err)
	}

	return nil
}

func (r *teacherRepository) GetPasswordByTeacherID(ctx context.Context, teacherID string) (string, error) {
	query := `SELECT password FROM teachers WHERE teacher_id = $1`

	var password string
	err := r.db.QueryRowContext(ctx, query, teacherID).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("teacher not found")
		}
		return "", fmt.Errorf("failed to get teacher password: %w", err)
	}

	return password, nil
}

func (r *teacherRepository) IsTeacherExist(ctx context.Context, teacherID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM teachers WHERE teacher_id = $1 AND is_active = true)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, teacherID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check teacher existence: %w", err)
	}

	return exists, nil
}

func (r *teacherRepository) GetStats(ctx context.Context) (*models.TeacherStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_teachers,
			COUNT(CASE WHEN is_active = true THEN 1 END) as active_teachers,
			COUNT(CASE WHEN is_active = false THEN 1 END) as inactive_teachers
		FROM teachers`

	stats := &models.TeacherStats{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalTeachers,
		&stats.ActiveTeachers,
		&stats.InactiveTeachers,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get teachers stats: %w", err)
	}

	return stats, nil
}
