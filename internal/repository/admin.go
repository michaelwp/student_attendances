package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/michaelwp/student_attendance/internal/models"
)

type adminRepository struct {
	db             *sql.DB
	teacherRepo    TeacherRepository
	studentRepo    StudentRepository
	classRepo      ClassRepository
	attendanceRepo AttendanceRepository
}

// NewAdminRepository creates a new admin repository
func NewAdminRepository(db *sql.DB) AdminRepository {
	return &adminRepository{db: db}
}

// NewAdminRepositoryWithDeps creates a new admin repository with dependencies
func NewAdminRepositoryWithDeps(db *sql.DB, teacherRepo TeacherRepository, studentRepo StudentRepository, classRepo ClassRepository, attendanceRepo AttendanceRepository) AdminRepository {
	return &adminRepository{
		db:             db,
		teacherRepo:    teacherRepo,
		studentRepo:    studentRepo,
		classRepo:      classRepo,
		attendanceRepo: attendanceRepo,
	}
}

func (r *adminRepository) Create(ctx context.Context, admin *models.Admin) error {
	query := `
		INSERT INTO admins (email, password, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		admin.Email,
		admin.Password,
		admin.IsActive,
	).Scan(&admin.ID, &admin.CreatedAt, &admin.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	return nil
}

func (r *adminRepository) GetByID(ctx context.Context, id uint) (*models.Admin, error) {
	query := `
		SELECT id, email, password, last_login, is_active, created_at, updated_at
		FROM admins WHERE id = $1`

	admin := &models.Admin{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&admin.ID,
		&admin.Email,
		&admin.Password,
		&admin.LastLogin,
		&admin.IsActive,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin not found")
		}
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}

	return admin, nil
}

func (r *adminRepository) GetByEmail(ctx context.Context, email string) (*models.Admin, error) {
	query := `
		SELECT id, email, password, last_login, is_active, created_at, updated_at
		FROM admins WHERE email = $1`

	admin := &models.Admin{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&admin.ID,
		&admin.Email,
		&admin.Password,
		&admin.LastLogin,
		&admin.IsActive,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin not found")
		}
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}

	return admin, nil
}

func (r *adminRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Admin, error) {
	query := `
		SELECT id, email, password, last_login, is_active, created_at, updated_at
		FROM admins
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get admins: %w", err)
	}
	defer rows.Close()

	var admins []*models.Admin
	for rows.Next() {
		admin := &models.Admin{}
		err := rows.Scan(
			&admin.ID,
			&admin.Email,
			&admin.Password,
			&admin.LastLogin,
			&admin.IsActive,
			&admin.CreatedAt,
			&admin.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan admin: %w", err)
		}
		admins = append(admins, admin)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return admins, nil
}

func (r *adminRepository) Update(ctx context.Context, admin *models.Admin) error {
	query := `
		UPDATE admins 
		SET email = $2, is_active = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		admin.ID,
		admin.Email,
		admin.IsActive,
	).Scan(&admin.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("admin not found")
		}
		return fmt.Errorf("failed to update admin: %w", err)
	}

	return nil
}

func (r *adminRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM admins WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete admin: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("admin not found")
	}

	return nil
}

func (r *adminRepository) UpdatePassword(ctx context.Context, email string, password string) error {
	query := `
		UPDATE admins 
		SET password = $2, updated_at = NOW()
		WHERE email = $1`

	result, err := r.db.ExecContext(ctx, query, email, password)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("admin not found")
	}

	return nil
}

func (r *adminRepository) UpdateLastLogin(ctx context.Context, id uint, lastLogin time.Time) error {
	query := `
		UPDATE admins 
		SET last_login = $2, updated_at = NOW()
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, lastLogin)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("admin not found")
	}

	return nil
}

func (r *adminRepository) SetActiveStatus(ctx context.Context, id uint, isActive bool) error {
	query := `
		UPDATE admins 
		SET is_active = $2, updated_at = NOW()
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, isActive)
	if err != nil {
		return fmt.Errorf("failed to set active status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("admin not found")
	}

	return nil
}

func (r *adminRepository) GetTotalAdmins(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM admins`

	var total int
	err := r.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total admins: %w", err)
	}

	return total, nil
}

func (r *adminRepository) IsAdminExist(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM admins WHERE email = $1 AND is_active = true)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check admin existence: %w", err)
	}

	return exists, nil
}

func (r *adminRepository) GetPasswordByEmail(ctx context.Context, email string) (string, error) {
	query := `SELECT password FROM admins WHERE email = $1`

	var password string
	err := r.db.QueryRowContext(ctx, query, email).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("admin not found")
		}
		return "", fmt.Errorf("failed to get admin password: %w", err)
	}

	return password, nil
}

func (r *adminRepository) GetStats(ctx context.Context) (*models.AdminStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_admins,
			COUNT(CASE WHEN is_active = true THEN 1 END) as active_admins,
			COUNT(CASE WHEN is_active = false THEN 1 END) as inactive_admins
		FROM admins`

	stats := &models.AdminStats{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalAdmins,
		&stats.ActiveAdmins,
		&stats.InactiveAdmins,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin stats: %w", err)
	}

	return stats, nil
}

func (r *adminRepository) GetDashboardStats(ctx context.Context) (*models.DashboardStats, error) {
	dashboardStats := &models.DashboardStats{}

	// Get admin stats
	adminStats, err := r.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin stats: %w", err)
	}
	dashboardStats.TotalAdmins = adminStats.TotalAdmins
	dashboardStats.ActiveAdmins = adminStats.ActiveAdmins
	dashboardStats.InactiveAdmins = adminStats.InactiveAdmins

	// Get teacher stats (use a direct database query if teacherRepo is nil)
	if r.teacherRepo != nil {
		teacherStats, err := r.teacherRepo.GetStats(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get teacher stats: %w", err)
		}
		dashboardStats.TotalTeachers = teacherStats.TotalTeachers
		dashboardStats.ActiveTeachers = teacherStats.ActiveTeachers
		dashboardStats.InactiveTeachers = teacherStats.InactiveTeachers
	} else {
		// Direct query if teacherRepo is not available
		query := `
			SELECT 
				COUNT(*) as total_teachers,
				COUNT(CASE WHEN is_active = true THEN 1 END) as active_teachers,
				COUNT(CASE WHEN is_active = false THEN 1 END) as inactive_teachers
			FROM teachers`
		err := r.db.QueryRowContext(ctx, query).Scan(
			&dashboardStats.TotalTeachers,
			&dashboardStats.ActiveTeachers,
			&dashboardStats.InactiveTeachers,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get teacher stats: %w", err)
		}
	}

	// Get student stats (use a direct database query if studentRepo is nil)
	if r.studentRepo != nil {
		studentStats, err := r.studentRepo.GetStats(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get student stats: %w", err)
		}
		dashboardStats.TotalStudents = studentStats.TotalStudents
		dashboardStats.ActiveStudents = studentStats.ActiveStudents
		dashboardStats.InactiveStudents = studentStats.InactiveStudents
	} else {
		// Direct query if studentRepo is not available
		query := `
			SELECT 
				COUNT(*) as total_students,
				COUNT(CASE WHEN is_active = true THEN 1 END) as active_students,
				COUNT(CASE WHEN is_active = false THEN 1 END) as inactive_students
			FROM students`
		err := r.db.QueryRowContext(ctx, query).Scan(
			&dashboardStats.TotalStudents,
			&dashboardStats.ActiveStudents,
			&dashboardStats.InactiveStudents,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get student stats: %w", err)
		}
	}

	// Get total classes (use a direct database query if classRepo is nil)
	if r.classRepo != nil {
		totalClasses, err := r.classRepo.GetTotalClasses(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get total classes: %w", err)
		}
		dashboardStats.TotalClasses = totalClasses
	} else {
		// Direct query if classRepo is not available
		query := `SELECT COUNT(*) FROM classes`
		err := r.db.QueryRowContext(ctx, query).Scan(&dashboardStats.TotalClasses)
		if err != nil {
			return nil, fmt.Errorf("failed to get total classes: %w", err)
		}
	}

	// Get today's attendance stats
	todayQuery := `
		SELECT 
			COUNT(*) as total_attendance_today,
			COUNT(CASE WHEN status = 'present' THEN 1 END) as present_today,
			COUNT(CASE WHEN status = 'absent' THEN 1 END) as absent_today,
			COUNT(CASE WHEN status = 'late' THEN 1 END) as late_today
		FROM attendances 
		WHERE DATE(date) = CURRENT_DATE`
	err = r.db.QueryRowContext(ctx, todayQuery).Scan(
		&dashboardStats.TotalAttendanceToday,
		&dashboardStats.PresentToday,
		&dashboardStats.AbsentToday,
		&dashboardStats.LateToday,
	)
	if err != nil {
		// If there's no attendance today, set values to 0
		dashboardStats.TotalAttendanceToday = 0
		dashboardStats.PresentToday = 0
		dashboardStats.AbsentToday = 0
		dashboardStats.LateToday = 0
	}

	return dashboardStats, nil
}
