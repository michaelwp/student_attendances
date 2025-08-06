package models

type UserType int

const (
	UserTypeAdmin UserType = iota
	UserTypeStudent
	UserTypeTeacher
)

func (u UserType) String() string {
	return []string{"admin", "student", "teacher"}[u]
}
