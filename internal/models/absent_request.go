package models

import (
	"strings"
	"time"
)

type AbsentRequestStatus string

const (
	AbsentRequestStatusPending  AbsentRequestStatus = "pending"
	AbsentRequestStatusApproved AbsentRequestStatus = "approved"
	AbsentRequestStatusRejected AbsentRequestStatus = "rejected"
)

type AbsentRequest struct {
	ID          uint                `json:"id" db:"id"`
	StudentID   string              `json:"student_id" db:"student_id"`
	ClassID     uint                `json:"class_id" db:"class_id"`
	RequestDate time.Time           `json:"request_date" db:"request_date"`
	Reason      string              `json:"reason" db:"reason"`
	Status      AbsentRequestStatus `json:"status" db:"status"`
	CreatedAt   time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at" db:"updated_at"`
	ApprovedBy  *uint               `json:"approved_by" db:"approved_by"`
	ApprovedAt  *time.Time          `json:"approved_at" db:"approved_at"`
	RejectedBy  *uint               `json:"rejected_by" db:"rejected_by"`
	RejectedAt  *time.Time          `json:"rejected_at" db:"rejected_at"`
	DeletedAt   *time.Time          `json:"deleted_at" db:"deleted_at"`
	DeletedBy   *uint               `json:"deleted_by" db:"deleted_by"`
}

func (AbsentRequest) TableName() string {
	return "absent_requests"
}

// AbsentRequestCreate is used for creating absent requests with date string input
type AbsentRequestCreate struct {
	RequestDate string `json:"request_date"`
	Reason      string `json:"reason"`
}

// ToAbsentRequest converts AbsentRequestCreate to AbsentRequest
func (arc *AbsentRequestCreate) ToAbsentRequest() (*AbsentRequest, error) {
	// Parse date string (YYYY-MM-DD format)
	requestDate, err := time.Parse("2006-01-02", arc.RequestDate)
	if err != nil {
		return nil, err
	}

	return &AbsentRequest{
		RequestDate: requestDate,
		Reason:      strings.TrimSpace(arc.Reason),
		Status:      AbsentRequestStatusPending,
	}, nil
}
