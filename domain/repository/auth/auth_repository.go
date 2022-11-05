package auth

import (
	entity "github.com/tpp/msf/domain/repository"
	"github.com/tpp/msf/model"
	"github.com/tpp/msf/shared/base"
	"github.com/tpp/msf/shared/context"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	List(ctx context.Context, limit, offset int, sort string, filters base.Filters) ([]*model.Role, error)
	GetUserByEmail(ctx context.Context, email string) (user *model.User, pass string, err error)
	GetUserId(ctx context.Context, email string) (id uint64, err error)
	UpdateTimeLastLogin(ctx context.Context, UserId int64) error
	CountRoleByUserID(ctx context.Context, roleId int64) (int64, error)
	GetRole(ctx context.Context, roleID uint64) (*model.Role, error)
	ListPermission(ctx context.Context) ([]*model.Permission, error)
	GetPermissionByRoleId(ctx context.Context, permissionId int64) (*model.Permission, error)
}

type repo struct {
	base.Repository
}

func (r *repo) CountRoleByUserID(ctx context.Context, roleId int64) (int64, error) {
	var total int64
	err := r.DB(ctx).Table("user_role").Where("status= true").Where("role_id =?", roleId).Count(&total).Error
	if err != nil {
		r.Error(ctx).Err(err).Msg("GetUserError")
		if err == gorm.ErrRecordNotFound {
			return 0, err
		}
		return 0, err
	}
	return total, nil
}

func (r *repo) GetPermissionByRoleId(ctx context.Context, permissionId int64) (*model.Permission, error) {
	var permission *entity.Permission
	err := r.DB(ctx).
		Model(&entity.Permission{}).
		Preload("Roles", " id IN (SELECT role_id FROM permision_role WHERE status = true AND permission_id = ?)", permissionId).
		//Preload("Roles.Permissions").
		//Preload("Org").
		Take(&permission, permissionId).Error

	if err != nil {
		r.Error(ctx).Err(err).Msg("Get Permission Error")
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return permission, nil

}

func (r *repo) ListPermission(ctx context.Context) ([]*model.Permission, error) {
	var permissions []*entity.Permission

	err := r.DB(ctx).Model(&entity.Permission{}).Find(&permissions).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []*model.Permission{}, nil
		}
		r.Error(ctx).Err(err).Msg("List Permission Error")
		return nil, err
	}
	return permissions, nil
}
func (r *repo) GetRole(ctx context.Context, roleID uint64) (*model.Role, error) {
	var role *entity.Role
	err := r.DB(ctx).
		Model(&entity.Role{}).Where("status= true").
		Preload("Permissions").
		//Preload( " id IN (SELECT role_id FROM user_role WHERE status = true AND user_id = ?)", roleID).
		Take(&role, roleID).Error
	if err != nil {
		r.Error(ctx).Err(err).Msg("Get Role ID Error")
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return role, nil
}
func (r *repo) GetUserByEmail(ctx context.Context, email string) (*model.User, string, error) {
	var user *entity.User
	if err := r.DB(ctx).
		Where("email = ?", email).Where("status = ?", true).
		Preload("Org").
		Preload("Roles.Permissions").
		Take(&user).
		Error; err != nil {
		r.Error(ctx).Err(err).Msg("GetUserError")
		if err == gorm.ErrRecordNotFound {
			return nil, "", nil
		}
		return nil, "", err
	}
	return user.User, user.Password, nil
}

// compare email received vs email in database and return email
func (r *repo) GetUserId(ctx context.Context, mail string) (uint64, error) {
	var users *entity.User
	if err := r.DB(ctx).
		Where("email = ?", mail).Where("status = ?", true).
		/*Preload("Org").
		Preload("Roles.Permissions").*/
		Take(&users).
		Error; err != nil {
		r.Error(ctx).Err(err).Msg("GetUserError")
		if err == gorm.ErrRecordNotFound {
			return 0, base.NewApiErrors("email", "This email is not registered")
		}
		return 0, err
	}
	return users.ID, nil
}
func (r *repo) List(ctx context.Context, limit, offset int, sort string, filters base.Filters) ([]*model.Role, error) {
	var roles []*entity.Role
	var total int64
	offset = offset * limit
	err := r.DB(ctx).
		Model(&entity.Role{}).
		Scopes(r.Filter(filters)).
		Count(&total).
		Scopes(r.Paginate(limit, offset, sort)).
		Preload("Permissions").
		Find(&roles).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []*model.Role{}, nil
		}
		r.Error(ctx).Err(err).Msg("List Roles Error")
		return nil, err
	}
	return roles, nil
}

func New() Repository {
	return &repo{base.NewBaseRepository("auth")}
}
func (r *repo) UpdateTimeLastLogin(ctx context.Context, UserId int64) error {
	time := time.Now().Format("2006-01-02 15:04:05")
	if err := r.DB(ctx).Model(&entity.User{}).
		Where("id = ?", UserId).Update("last_login", time).
		Error; err != nil {
		r.Error(ctx).Err(err).Msg("update fail")
		if err == gorm.ErrRecordNotFound {
			return err
		}
		return err
	}
	return nil
}
