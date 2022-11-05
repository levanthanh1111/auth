package repository

import (
	"github.com/tpp/msf/model"
)

type User struct {
	*model.User
	Password string
}
type UserRole struct {
	UserId int64
	RoleId int64
}

func MapUserModels(users []*User) []*model.User {
	if users == nil {
		return nil
	}
	var models = make([]*model.User, len(users))
	for idx, user := range users {
		models[idx] = user.User
	}
	return models
}

type Org = model.Org

type Permission = model.Permission

type Role = model.Role

type Contract = model.Contract

type WithdrawRequest struct {
	*model.WithdrawRequest
	Contract          *Contract
	ProjectContractor *User
}
