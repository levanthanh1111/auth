package users

import (
	"github.com/tpp/msf/model"
	"github.com/tpp/msf/shared/base"
	"github.com/tpp/msf/shared/validator"
)

var v = validator.Get()

type listParams = base.ListParams

var listAcceptedFilterKeys = []string{"full_name", "email", "org_name", "role_name", "role_id"}

type createUserReq struct {
	FullName string  `json:"full_name" schema:"full_name" validate:"required"`
	Email    string  `json:"email" schema:"email" validate:"required,email"`
	Role     []int64 `json:"role"`
}

func (r *createUserReq) IsValid() error {
	return v.Struct(r)
}
func (r *createUserReq) mapUserModel() *model.User {
	return &model.User{
		FullName: r.FullName,
		Email:    r.Email,
		Status:   false,
		IsAdmin:  false,
	}
}

type createRoleReq struct {
	Name model.RoleName `json:"name" schema:"name" validate:"required"`
}

func (r *createRoleReq) IsValid() error {
	return v.Struct(r)
}
func (r *createRoleReq) mapRoleModel() *model.Role {
	return &model.Role{
		Name:   r.Name,
		Status: true,
	}
}

type assignRoleReq struct {
	UserID   uint64  `json:"id" schema:"id" validate:"required"`
	RoleIDs  []int64 `json:"role_ids" schema:"role_id" validate:"required"`
	UserName string  `json:"user_name"`
}

func (r *assignRoleReq) IsValid() error {
	return v.Struct(r)
}

type PassWordReset struct {
	Password string `json:"password" validate:"required,gte=6"`
}
type UpdateName struct {
	FullName string `json:"full_name" validate:"required"`
}
type UpdateActive struct {
	UserId   int  `json:"id"`
	IsActive bool `json:"is_active"`
}

func (p *PassWordReset) IsValid() error {
	return v.Struct(p)
}

type PassWordResetByAdmin struct {
	UserId   int    `json:"id"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (p *PassWordReset) setPassWordReset(password string) *PassWordReset {
	return &PassWordReset{password}
}
