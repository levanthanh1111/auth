package model

import (
	"time"
)

// User model
type User struct {
	ID         uint64    `json:"id" mapstructure:"id"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	IsAdmin    bool      `json:"is_admin"`
	Roles      []*Role   `json:"roles" gorm:"many2many:user_role;"`
	OrgID      uint64    `json:"-"`
	Org        *Org      `json:"orgs"`
	Status     bool      `json:"status"`
	Last_login time.Time `json:"last_Login"`
	CreateAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdateAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type ListUserRes struct {
	ID         uint64    `json:"id" mapstructure:"id"`
	FullName   string    `json:"full_name"`
	Roles      []*Role   `json:"roles" `
	Email      string    `json:"email"`
	Status     bool      `json:"status"`
	Last_login time.Time `json:"last_Login"`
	CreateAt   time.Time `json:"created_at"`
	UpdateAt   time.Time `json:"updated_at"`
}

func (r *User) MapUserModel() *ListUserRes {
	return &ListUserRes{
		ID:         r.ID,
		FullName:   r.FullName,
		Roles:      r.Roles,
		Email:      r.Email,
		Status:     r.Status,
		Last_login: r.Last_login,
		CreateAt:   r.CreateAt,
		UpdateAt:   r.UpdateAt,
	}
}
