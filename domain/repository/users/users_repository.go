//go:generate mockery --name=Repository
package users

import (
	entity "github.com/tpp/msf/domain/repository"
	"github.com/tpp/msf/model"
	"github.com/tpp/msf/shared/base"
	"github.com/tpp/msf/shared/context"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	List(ctx context.Context, limit, offset int, sort string, filters base.Filters) ([]*model.User, int64, error)
	Create(ctx context.Context, user *entity.User) error
	CreateRole(ctx context.Context, role *entity.Role) error
	Get(ctx context.Context, userID uint64) (*model.User, error)
	GetOrg(ctx context.Context, orgID uint64) (*model.Org, error)
	GetRole(ctx context.Context, roleID uint64) (*model.Role, error)
	UpdateUserRole(ctx context.Context, userID, roleID uint64) error
	UpdatePassWord(ctx context.Context, userID int64, passWord string) error
	UpdateName(ctx context.Context, name string, userID int64) error
	UpdateActive(ctx context.Context, userID int64, active bool) error
	AssignMultipleRole(ctx context.Context, userRoles []entity.UserRole) error
	GetRolesByUserID(ctx context.Context, userID int64) ([]*entity.UserRole, error)
	DeleteRolesByUserID(ctx context.Context, roleIds []int64, userId int64) error
	ActiveRolesByUserID(ctx context.Context, roleIds []int64, userId int64) error
	GetRolesByUserIDtoUpdate(ctx context.Context, userID int64) ([]*entity.UserRole, error)
	TimeCreateUser(ctx context.Context, UserId int64) error
	TimeCreateRole(ctx context.Context, RoleId int64) error
}

type repo struct {
	base.Repository
}

func (r *repo) List(ctx context.Context, limit, offset int, sort string, filters base.Filters) ([]*model.User, int64, error) {
	var users []*entity.User
	var total int64
	offset = offset * limit

	err := r.DB(ctx).
		Model(&entity.User{}).
		Scopes(r.Filter(filters)).
		Count(&total).
		Scopes(r.Paginate(limit, offset, sort)).
		Preload("Roles.Permissions").
		Preload("Org").
		Find(&users).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		r.Error(ctx).Err(err).Msg("ListUserError")
		return []*model.User{}, 0, err
	}
	return entity.MapUserModels(users), total, nil
}

func (r *repo) Create(ctx context.Context, user *entity.User) error {
	if err := r.DB(ctx).Create(user).Error; err != nil {
		r.Error(ctx).Err(err).Msg("StoreUserError")
		return err
	}

	return nil
}

func (r *repo) CreateRole(ctx context.Context, role *model.Role) error {
	var roles *entity.Role
	err := r.DB(ctx).
		Where("name = ?", role.Name).
		Take(&roles).
		Error
	if err == nil {
		return base.NewApiErrors("name", "name already existed")
	} else {

		errr := r.DB(ctx).Select("Name", "CreateAt", "UpdateAt").Create(&role).Error
		if errr != nil {
			r.Error(ctx).Err(errr).Msg("StoreUserError")
			return errr
		}
	}
	return nil
}

func (r *repo) Get(ctx context.Context, userID uint64) (*model.User, error) {
	var user *entity.User
	err := r.DB(ctx).
		Model(&entity.User{}).
		Preload("Roles", " id IN (SELECT role_id FROM user_role WHERE status = true AND user_id = ?)", userID).
		Preload("Roles.Permissions").
		Preload("Org").
		Take(&user, userID).Error

	if err != nil {
		r.Error(ctx).Err(err).Msg("GetUserError")
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return user.User, nil
}

func (r *repo) GetOrg(ctx context.Context, orgID uint64) (*model.Org, error) {
	var org *entity.Org
	if err := r.DB(ctx).Preload("Role").Take(&org, orgID).Error; err != nil {
		r.Error(ctx).Err(err).Msg("GetOrgError")
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return org, nil
}

func (r *repo) GetRole(ctx context.Context, roleID uint64) (*model.Role, error) {
	var role *entity.Role
	if err := r.DB(ctx).Take(&role, roleID).Error; err != nil {
		r.Error(ctx).Err(err).Msg("GetRoleError")
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return role, nil
}

func (r *repo) UpdateUserRole(ctx context.Context, userID, roleID uint64) error {
	if err := r.DB(ctx).Table("user_role").Where("user_id = ?", userID).Update("role_id", roleID).Error; err != nil {
		r.Error(ctx).Err(err).Msg("GetRoleError")
		return err
	}
	return nil
}
func (r *repo) GetRolesByUserID(ctx context.Context, userID int64) ([]*entity.UserRole, error) {
	var roles []*entity.UserRole
	if err := r.DB(ctx).Table("user_role").Where("user_id = ?", userID).Where("status = ?", true).Find(&roles).Error; err != nil {
		r.Error(ctx).Err(err).Msg("GetRoleError")
		return nil, err
	}
	return roles, nil
}
func (r *repo) GetRolesByUserIDtoUpdate(ctx context.Context, userID int64) ([]*entity.UserRole, error) {
	var roles []*entity.UserRole
	if err := r.DB(ctx).Table("user_role").Where("user_id = ?", userID).Find(&roles).Error; err != nil {
		r.Error(ctx).Err(err).Msg("GetRoleError")
		return nil, err
	}
	return roles, nil
}
func (r *repo) UpdatePassWord(ctx context.Context, userID int64, passWord string) error {

	if err := r.DB(ctx).Table("users").Where("id = ?", userID).Update("password", passWord).Update("status", true).Error; err != nil {
		r.Error(ctx).Err(err).Msg("update password error")
		return err
	}
	return nil
}
func (r *repo) UpdateName(ctx context.Context, name string, userID int64) error {

	if err := r.DB(ctx).Table("users").Where("id = ?", userID).Update("full_name", name).Error; err != nil {
		r.Error(ctx).Err(err).Msg("update name error")
		return err
	}
	return nil
}
func (r *repo) UpdateActive(ctx context.Context, userID int64, active bool) error {

	if err := r.DB(ctx).Table("users").Where("id = ?", userID).Update("status", active).Error; err != nil {
		r.Error(ctx).Err(err).Msg("update active error")
		return err
	}
	return nil
}
func (r *repo) AssignMultipleRole(ctx context.Context, userRoles []entity.UserRole) error {

	if err := r.DB(ctx).Table("user_role").Create(userRoles).Error; err != nil {
		r.Error(ctx).Err(err).Msg("Create user_role fail")
		return err
	}
	return nil
}

func (r *repo) DeleteRolesByUserID(ctx context.Context, roleIds []int64, userId int64) error {

	if err := r.DB(ctx).Table("user_role").Where("user_id = ?", userId).Where(" role_id IN ? ", roleIds).Update("status", false).Error; err != nil {
		r.Error(ctx).Err(err).Msg("update active error")
		return err
	}
	return nil
}

func (r *repo) ActiveRolesByUserID(ctx context.Context, roleIds []int64, userId int64) error {

	if err := r.DB(ctx).Table("user_role").Where("user_id = ?", userId).Where(" role_id IN ? ", roleIds).Update("status", true).Error; err != nil {
		r.Error(ctx).Err(err).Msg("update active error")
		return err
	}
	return nil
}
func (r *repo) TimeCreateUser(ctx context.Context, UserId int64) error {
	time := time.Now().Format("2006-01-02 15:04:05")
	if err := r.DB(ctx).Model(&entity.User{}).
		Where("id = ?", UserId).Update("created_at", time).
		Error; err != nil {
		r.Error(ctx).Err(err).Msg("update fail")
		if err == gorm.ErrRecordNotFound {
			return err
		}
		return err
	}
	return nil
}
func (r *repo) TimeCreateRole(ctx context.Context, RoleId int64) error {
	time := time.Now().Format("2006-01-02 15:04:05")
	if err := r.DB(ctx).Model(&entity.Role{}).
		Where("id = ?", RoleId).Update("created_at", time).
		Error; err != nil {
		r.Error(ctx).Err(err).Msg("update fail")
		if err == gorm.ErrRecordNotFound {
			return err
		}
		return err
	}
	return nil
}
func New() Repository {
	baseRepo := base.NewBaseRepository("users")

	orgsJoin := func(d *gorm.DB) *gorm.DB { return d.Joins(`INNER JOIN "orgs" ON "org_id" = "orgs"."id"`) }

	baseRepo.RegisterActualKey("role_name", "roles.name")
	baseRepo.RegisterActualKey("org_name", "orgs.name")

	baseRepo.RegisterJoinScope("org_name", orgsJoin)
	baseRepo.RegisterJoinScope("role_name", orgsJoin, func(d *gorm.DB) *gorm.DB { return d.Joins(`INNER JOIN "roles" ON "orgs"."type" = "roles"."id"`) })

	return &repo{baseRepo}
}
