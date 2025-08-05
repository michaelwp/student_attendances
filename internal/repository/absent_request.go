package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/michaelwp/student_attendance/internal/models"
)

type absentRequestRepository struct {
	db *sql.DB
}

// NewAbsentRequestRepository creates a new absent request repository
func NewAbsentRequestRepository(db *sql.DB) AbsentRequestRepository {
	return &absentRequestRepository{db: db}
}

func (r *absentRequestRepository) Create(ctx context.Context, request *models.AbsentRequest) error {
	query := `
		INSERT INTO absent_requests (student_id, class_id, request_date, reason, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		request.StudentID,
		request.ClassID,
		request.RequestDate,
		request.Reason,
		request.Status,
	).Scan(&request.ID, &request.CreatedAt, &request.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create absent request: %w", err)
	}

	return nil
}

func (r *absentRequestRepository) GetByID(ctx context.Context, id uint) (*models.AbsentRequest, error) {
	query := `
		SELECT id, student_id, class_id, request_date, reason, status, created_at, updated_at
		FROM absent_requests WHERE id = $1`

	request := &models.AbsentRequest{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&request.ID,
		&request.StudentID,
		&request.ClassID,
		&request.RequestDate,
		&request.Reason,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("absent request not found")
		}
		return nil, fmt.Errorf("failed to get absent request: %w", err)
	}

	return request, nil
}

func (r *absentRequestRepository) GetByStudent(ctx context.Context, studentID string, limit, offset int) ([]*models.AbsentRequest, error) {
	query := `
		SELECT id, student_id, class_id, request_date, reason, status, created_at, updated_at
		FROM absent_requests 
		WHERE student_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get absent requests by student: %w", err)
	}
	defer rows.Close()

	var requests []*models.AbsentRequest
	for rows.Next() {
		request := &models.AbsentRequest{}
		err := rows.Scan(
			&request.ID,
			&request.StudentID,
			&request.ClassID,
			&request.RequestDate,
			&request.Reason,
			&request.Status,
			&request.CreatedAt,
			&request.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan absent request: %w", err)
		}
		requests = append(requests, request)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate absent requests: %w", err)
	}

	return requests, nil
}

func (r *absentRequestRepository) GetByClass(ctx context.Context, classID uint, limit, offset int) ([]*models.AbsentRequest, error) {
	query := `
		SELECT id, student_id, class_id, request_date, reason, status, created_at, updated_at
		FROM absent_requests 
		WHERE class_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, classID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get absent requests by class: %w", err)
	}
	defer rows.Close()

	var requests []*models.AbsentRequest
	for rows.Next() {
		request := &models.AbsentRequest{}
		err := rows.Scan(
			&request.ID,
			&request.StudentID,
			&request.ClassID,
			&request.RequestDate,
			&request.Reason,
			&request.Status,
			&request.CreatedAt,
			&request.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan absent request: %w", err)
		}
		requests = append(requests, request)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate absent requests: %w", err)
	}

	return requests, nil
}

func (r *absentRequestRepository) GetByStatus(ctx context.Context, status models.AbsentRequestStatus, limit, offset int) ([]*models.AbsentRequest, error) {
	query := `
		SELECT id, student_id, class_id, request_date, reason, status, created_at, updated_at
		FROM absent_requests 
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get absent requests by status: %w", err)
	}
	defer rows.Close()

	var requests []*models.AbsentRequest
	for rows.Next() {
		request := &models.AbsentRequest{}
		err := rows.Scan(
			&request.ID,
			&request.StudentID,
			&request.ClassID,
			&request.RequestDate,
			&request.Reason,
			&request.Status,
			&request.CreatedAt,
			&request.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan absent request: %w", err)
		}
		requests = append(requests, request)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate absent requests: %w", err)
	}

	return requests, nil
}

func (r *absentRequestRepository) GetPending(ctx context.Context, limit, offset int) ([]*models.AbsentRequest, error) {
	return r.GetByStatus(ctx, models.AbsentRequestStatusPending, limit, offset)
}

func (r *absentRequestRepository) Update(ctx context.Context, request *models.AbsentRequest) error {
	query := `
		UPDATE absent_requests 
		SET student_id = $2, class_id = $3, request_date = $4, reason = $5, status = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		request.ID,
		request.StudentID,
		request.ClassID,
		request.RequestDate,
		request.Reason,
		request.Status,
	).Scan(&request.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("absent request not found")
		}
		return fmt.Errorf("failed to update absent request: %w", err)
	}

	return nil
}

func (r *absentRequestRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM absent_requests WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete absent request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("absent request not found")
	}

	return nil
}
