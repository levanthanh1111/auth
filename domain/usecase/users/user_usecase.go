package users

import (
	"github.com/rs/xid"
	entity "github.com/tpp/msf/domain/repository"
	helper "github.com/tpp/msf/shared/auth"
	"github.com/tpp/msf/shared/utils"
	"gorm.io/gorm"
	"strings"

	"github.com/tpp/msf/domain/repository/users"
	"github.com/tpp/msf/model"
	"github.com/tpp/msf/shared/base"
	"github.com/tpp/msf/shared/context"
)

type Usecase interface {
	ListUsers(ctx context.Context, limit, offset int, sort string, filters base.Filters) ([]*model.User, int64, error)
	GetUser(ctx context.Context, userID uint64) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	CreateRole(ctx context.Context, user *model.Role) error
	AssignRole(ctx context.Context, userID int64, roleIDs []int64) error
	UpdatePassWord(ctx context.Context, userID int64, passWord string) error
	AssignMultipleRole(ctx context.Context, userID int64, roleIds []int64) error
	UpdateName(ctx context.Context, name string, userID int64) error
	UpdateActive(ctx context.Context, userID int64, isActive bool) error
}

type usecase struct {
	base.Usecase
	userRepo users.Repository
}

func (u *usecase) ListUsers(ctx context.Context, limit, offset int, sort string, filters base.Filters) ([]*model.User, int64, error) {
	return u.userRepo.List(ctx, limit, offset, sort, filters)
}

func (u *usecase) GetUser(ctx context.Context, userID uint64) (*model.User, error) {
	var user *model.User
	var err error
	user, err = u.userRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, base.ErrorNotFound
	}
	return user, nil
}

func (u *usecase) CreateUser(ctx context.Context, user *model.User) error {

	autoPassword := xid.New().String()[:6]
	hashPassword := utils.HashAndSaltPassword(autoPassword)

	user1 := entity.User{User: user, Password: hashPassword}
	err := u.userRepo.Create(ctx, &user1)
	if err != nil {
		if strings.Contains(err.Error(), "duplicated") || strings.Contains(err.Error(), "email") {

			return base.NewApiErrors("email", "email already existed")
		}
	} else {
		accessToken, err := helper.GenerateJWTToken(user1.ID)
		err = u.SendEmail(accessToken, user1.Email, "active_template.html", "Active your account here")
		if err != nil {
			u.Info(ctx).Err(err).Msg("Can not send mail")
		}
		err = u.userRepo.TimeCreateUser(ctx, int64(user1.ID))
		if err != nil {
			return err
		}
	}
	return nil
}
func (u *usecase) CreateRole(ctx context.Context, role *model.Role) error {
	err := u.userRepo.CreateRole(ctx, role)
	if err != nil {
		u.Info(ctx).Err(err).Msg("Can not create role")
		return err
	}
	err = u.userRepo.TimeCreateRole(ctx, int64(role.ID))
	if err != nil {
		return err
	}
	return nil
}

func (u *usecase) AssignRole(ctx context.Context, userID int64, roleIDs []int64) error {
	var roleIDsForInsert []int64
	var roleIdsForDelete []int64
	var roleIDsForDeactivate []int64
	roleUsersExisted, error := u.userRepo.GetRolesByUserIDtoUpdate(ctx, userID)
	if error != gorm.ErrRecordNotFound && error != nil {
		return error
	}
	if error == gorm.ErrRecordNotFound {
		roleIDsForInsert = roleIDs
	} else {
		var roleIDsExist []int64
		for _, roleUserExisted := range roleUsersExisted {
			roleIDsExist = append(roleIDsExist, roleUserExisted.RoleId)
		}
		roleIDsForInsert, _ = getElementNotExistedOtherArrAndTheSame(roleIDs, roleIDsExist)
		roleIdsForDelete, roleIDsForDeactivate = getElementNotExistedOtherArrAndTheSame(roleIDsExist, roleIDs)
	}
	if len(roleIDsForInsert) > 0 {
		for _, roleID := range roleIDsForInsert {
			role, err := u.userRepo.GetRole(ctx, uint64(roleID))
			if err != nil {
				return err
			}
			if role == nil {
				return base.NewApiErrors("role_id", "this role does not existed")
			}
		}
		user, err := u.userRepo.Get(ctx, uint64(userID))
		if err != nil {
			return err
		}
		if user == nil {
			return base.NewApiErrors("user_id", "this user does not existed")
		}

		var userRoles []entity.UserRole
		for _, role := range roleIDsForInsert {
			userRoles = append(userRoles, entity.UserRole{
				UserId: userID, RoleId: role,
			})
		}
		err = u.userRepo.AssignMultipleRole(ctx, userRoles)
	}
	if len(roleIdsForDelete) > 0 {
		err1 := u.userRepo.DeleteRolesByUserID(ctx, roleIdsForDelete, userID)
		if err1 != nil {
			return err1
		}
	}
	if len(roleIDsForDeactivate) > 0 {
		err2 := u.userRepo.ActiveRolesByUserID(ctx, roleIDsForDeactivate, userID)
		if err2 != nil {
			return err2
		}
	}

	return nil
}
func (u *usecase) UpdatePassWord(ctx context.Context, userID int64, passWord string) error {
	_, err := u.userRepo.Get(ctx, uint64(userID))
	if err != nil {
		return err
	}
	err = u.userRepo.UpdatePassWord(ctx, userID, utils.HashAndSaltPassword(passWord))
	if err != nil {
		return err
	}
	return nil
}

func (u *usecase) AssignMultipleRole(ctx context.Context, userID int64, roleIds []int64) error {
	var userRoles []entity.UserRole
	for _, role := range roleIds {
		userRoles = append(userRoles, entity.UserRole{
			UserId: userID, RoleId: role,
		})
	}
	return u.userRepo.AssignMultipleRole(ctx, userRoles)
}

func (u *usecase) UpdateName(ctx context.Context, name string, userID int64) error {
	return u.userRepo.UpdateName(ctx, name, userID)
}
func (u *usecase) UpdateActive(ctx context.Context, userID int64, isActive bool) error {
	return u.userRepo.UpdateActive(ctx, userID, isActive)
}
func New() Usecase {
	return &usecase{
		Usecase:  base.NewBaseUsecase("users"),
		userRepo: users.New(),
	}
}
func getElementNotExistedOtherArrAndTheSame(source []int64, dest []int64) ([]int64, []int64) {
	var isExisted = false
	var diff []int64
	var same []int64
	for _, e1 := range source {
		for _, e2 := range dest {
			if int64(e1) == int64(e2) {
				same = append(same, e1)
				isExisted = true
				break
			}
		}
		if !isExisted {
			diff = append(diff, e1)
		}
		isExisted = false
	}
	return diff, same
}
