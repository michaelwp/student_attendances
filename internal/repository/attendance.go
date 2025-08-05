package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/michaelwp/student_attendance/internal/models"
)

type attendanceRepository struct {
	db *sql.DB
}

// NewAttendanceRepository creates a new attendance repository
func NewAttendanceRepository(db *sql.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) Create(ctx context.Context, attendance *models.Attendance) error {
	query := `
		INSERT INTO attendances (student_id, class_id, date, status, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		attendance.StudentID,
		attendance.ClassID,
		attendance.Date,
		attendance.Status,
		attendance.Description,
	).Scan(&attendance.ID, &attendance.CreatedAt, &attendance.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create attendance: %w", err)
	}

	return nil
}

func (r *attendanceRepository) GetByID(ctx context.Context, id uint) (*models.Attendance, error) {
	query := `
		SELECT id, student_id, class_id, date, status, description, created_at, updated_at
		FROM attendances WHERE id = $1`

	attendance := &models.Attendance{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&attendance.ID,
		&attendance.StudentID,
		&attendance.ClassID,
		&attendance.Date,
		&attendance.Status,
		&attendance.Description,
		&attendance.CreatedAt,
		&attendance.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("attendance not found")
		}
		return nil, fmt.Errorf("failed to get attendance: %w", err)
	}

	return attendance, nil
}

func (r *attendanceRepository) GetByStudentAndDate(ctx context.Context, studentID string, date time.Time) (*models.Attendance, error) {
	query := `
		SELECT id, student_id, class_id, date, status, description, created_at, updated_at
		FROM attendances 
		WHERE student_id = $1 AND DATE(date) = DATE($2)`

	attendance := &models.Attendance{}
	err := r.db.QueryRowContext(ctx, query, studentID, date).Scan(
		&attendance.ID,
		&attendance.StudentID,
		&attendance.ClassID,
		&attendance.Date,
		&attendance.Status,
		&attendance.Description,
		&attendance.CreatedAt,
		&attendance.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("attendance not found")
		}
		return nil, fmt.Errorf("failed to get attendance: %w", err)
	}

	return attendance, nil
}

func (r *attendanceRepository) GetByStudent(ctx context.Context, studentID string, limit, offset int) ([]*models.Attendance, error) {
	query := `
		SELECT id, student_id, class_id, date, status, description, created_at, updated_at
		FROM attendances 
		WHERE student_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances by student: %w", err)
	}
	defer rows.Close()

	var attendances []*models.Attendance
	for rows.Next() {
		attendance := &models.Attendance{}
		err := rows.Scan(
			&attendance.ID,
			&attendance.StudentID,
			&attendance.ClassID,
			&attendance.Date,
			&attendance.Status,
			&attendance.Description,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attendance: %w", err)
		}
		attendances = append(attendances, attendance)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate attendances: %w", err)
	}

	return attendances, nil
}

func (r *attendanceRepository) GetByClass(ctx context.Context, classID uint, limit, offset int) ([]*models.Attendance, error) {
	query := `
		SELECT id, student_id, class_id, date, status, description, created_at, updated_at
		FROM attendances 
		WHERE class_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, classID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances by class: %w", err)
	}
	defer rows.Close()

	var attendances []*models.Attendance
	for rows.Next() {
		attendance := &models.Attendance{}
		err := rows.Scan(
			&attendance.ID,
			&attendance.StudentID,
			&attendance.ClassID,
			&attendance.Date,
			&attendance.Status,
			&attendance.Description,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attendance: %w", err)
		}
		attendances = append(attendances, attendance)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate attendances: %w", err)
	}

	return attendances, nil
}

func (r *attendanceRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*models.Attendance, error) {
	query := `
		SELECT id, student_id, class_id, date, status, description, created_at, updated_at
		FROM attendances 
		WHERE DATE(date) >= DATE($1) AND DATE(date) <= DATE($2)
		ORDER BY date DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances by date range: %w", err)
	}
	defer rows.Close()

	var attendances []*models.Attendance
	for rows.Next() {
		attendance := &models.Attendance{}
		err := rows.Scan(
			&attendance.ID,
			&attendance.StudentID,
			&attendance.ClassID,
			&attendance.Date,
			&attendance.Status,
			&attendance.Description,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attendance: %w", err)
		}
		attendances = append(attendances, attendance)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate attendances: %w", err)
	}

	return attendances, nil
}

func (r *attendanceRepository) Update(ctx context.Context, attendance *models.Attendance) error {
	query := `
		UPDATE attendances 
		SET student_id = $2, class_id = $3, date = $4, status = $5, description = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		attendance.ID,
		attendance.StudentID,
		attendance.ClassID,
		attendance.Date,
		attendance.Status,
		attendance.Description,
	).Scan(&attendance.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("attendance not found")
		}
		return fmt.Errorf("failed to update attendance: %w", err)
	}

	return nil
}

func (r *attendanceRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM attendances WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete attendance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("attendance not found")
	}

	return nil
}