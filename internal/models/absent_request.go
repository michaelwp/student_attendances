package models

import "time"

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
}

func (AbsentRequest) TableName() string {
	return "absent_requests"
}