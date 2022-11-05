package model

import "time"

// RoleName type def
type RoleName string

// Role model
type Role struct {
	ID          uint64        `json:"id"`
	Name        RoleName      `json:"name"`
	Permissions []*Permission `json:"permissions" gorm:"many2many:permision_role;"`
	CreateAt    time.Time     `json:"created_at" gorm:"column:created_at"`
	UpdateAt    time.Time     `json:"updated_at" gorm:"column:updated_at"`
	Status      bool          `json:"status"`
	TotalUsers  int64         `json:"total_users"`
}

type PermissionNew struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type ListRoleRes struct {
	ID          uint64           `json:"id"`
	Name        string           `json:"name"`
	Permissions []*PermissionNew `json:"permissions"  gorm:"many2many:permision_role;"`
	CreateAt    time.Time        `json:"created_at" `
	UpdateAt    time.Time        `json:"updated_at" `
	Status      bool             `json:"status"`
	TotalUsers  int64            `json:"total_users"`
}

func (r *Role) MapRoleListModel() *ListRoleRes {
	var permissions []*PermissionNew
	for _, p := range r.Permissions {
		var p1 = p.ConvertPermissionNew()
		permissions = append(permissions, &p1)
	}
	return &ListRoleRes{
		ID:          r.ID,
		Name:        string(r.Name),
		Permissions: permissions,
		CreateAt:    r.CreateAt,
		UpdateAt:    r.UpdateAt,
		Status:      r.Status,
		TotalUsers:  r.TotalUsers,
	}
}

func (r *Role) SetCount(total int64) {
	r.TotalUsers = total
}

const (
	PlannerRoleID           uint64 = 1
	ProjectContractorRoleID uint64 = 2
	SupplyVendorRoleID      uint64 = 3
)
