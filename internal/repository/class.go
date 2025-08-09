package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/michaelwp/student_attendance/internal/models"
)

type classRepository struct {
	db *sql.DB
}

// NewClassRepository creates a new class repository
func NewClassRepository(db *sql.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) Create(ctx context.Context, class *models.Class) error {
	query := `
		INSERT INTO classes (name, homeroom_teacher, description, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		class.Name,
		class.HomeroomTeacher,
		class.Description,
	).Scan(&class.ID, &class.CreatedAt, &class.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create class: %w", err)
	}

	return nil
}

func (r *classRepository) GetByID(ctx context.Context, id uint) (*models.Class, error) {
	query := `
		SELECT id, name, homeroom_teacher, description, created_at, updated_at
		FROM classes WHERE id = $1 AND deleted_at IS NULL`

	class := &models.Class{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&class.ID,
		&class.Name,
		&class.HomeroomTeacher,
		&class.Description,
		&class.CreatedAt,
		&class.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("class not found")
		}
		return nil, fmt.Errorf("failed to get class: %w", err)
	}

	return class, nil
}

func (r *classRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Class, error) {
	query := `
		SELECT id, name, homeroom_teacher, description, created_at, updated_at
		FROM classes
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get classes: %w", err)
	}
	defer rows.Close()

	var classes []*models.Class
	for rows.Next() {
		class := &models.Class{}
		err := rows.Scan(
			&class.ID,
			&class.Name,
			&class.HomeroomTeacher,
			&class.Description,
			&class.CreatedAt,
			&class.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan class: %w", err)
		}
		classes = append(classes, class)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate classes: %w", err)
	}

	return classes, nil
}

func (r *classRepository) GetByTeacher(ctx context.Context, teacherID string) ([]*models.Class, error) {
	query := `
		SELECT id, name, homeroom_teacher, description, created_at, updated_at
		FROM classes 
		WHERE homeroom_teacher = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get classes by teacher: %w", err)
	}
	defer rows.Close()

	var classes []*models.Class
	for rows.Next() {
		class := &models.Class{}
		err := rows.Scan(
			&class.ID,
			&class.Name,
			&class.HomeroomTeacher,
			&class.Description,
			&class.CreatedAt,
			&class.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan class: %w", err)
		}
		classes = append(classes, class)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate classes: %w", err)
	}

	return classes, nil
}

func (r *classRepository) Update(ctx context.Context, class *models.Class) error {
	query := `
		UPDATE classes 
		SET name = $2, homeroom_teacher = $3, description = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		class.ID,
		class.Name,
		class.HomeroomTeacher,
		class.Description,
	).Scan(&class.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("class not found")
		}
		return fmt.Errorf("failed to update class: %w", err)
	}

	return nil
}

func (r *classRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM classes WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete class: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("class not found")
	}

	return nil
}

func (r *classRepository) GetTotalClasses(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM classes WHERE deleted_at IS NULL`

	var total int
	err := r.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total classes: %w", err)
	}

	return total, nil
}

func (r *classRepository) UpdateDeleteInfo(ctx context.Context, id uint, deletedBy uint) error {
	query := `
		UPDATE classes 
		SET deleted_at = NOW(), deleted_by = $2
		WHERE id = $1
		RETURNING deleted_at`

	var deletedAt string
	err := r.db.QueryRowContext(ctx, query, id, deletedBy).Scan(&deletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("class not found")
		}
		return fmt.Errorf("failed to update class delete info: %w", err)
	}

	return nil
}
