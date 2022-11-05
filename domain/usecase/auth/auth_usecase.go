package auth

import (
	"github.com/tpp/msf/domain/repository/auth"
	"github.com/tpp/msf/model"
	helper "github.com/tpp/msf/shared/auth"
	"github.com/tpp/msf/shared/base"
	"github.com/tpp/msf/shared/context"
	"github.com/tpp/msf/shared/utils"
)

type Usecase interface {
	ListRoles(ctx context.Context, limit, offset int, sort string, filters base.Filters) ([]*model.Role, error)
	Login(ctx context.Context, email, password string) (user *model.User, accessToken string, err error)
	ForgotPassword(ctx context.Context, email string) (accessTolken string, err error)
	GetRole(ctx context.Context, roleID uint64) (*model.Role, error)
	ListPermission(ctx context.Context) ([]*model.Permission, error)
}

type usecase struct {
	base.Usecase
	repo auth.Repository
}

func (u *usecase) ListPermission(ctx context.Context) ([]*model.Permission, error) {
	permissions, err := u.repo.ListPermission(ctx)
	if err != nil {
		return nil, err
	}
	return permissions, err

}
func (u *usecase) ListRoles(ctx context.Context, limit, offset int, sort string, filters base.Filters) ([]*model.Role, error) {
	roles, err := u.repo.List(ctx, limit, offset, sort, filters)
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		total, err := u.repo.CountRoleByUserID(ctx, int64(role.ID))
		if err != nil {
			return nil, err
		}
		role.SetCount(total)
	}
	return roles, err
}
func (u *usecase) GetRole(ctx context.Context, roleID uint64) (*model.Role, error) {
	var roles *model.Role
	var err error
	roles, err = u.repo.GetRole(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if roles == nil {
		return nil, base.ErrorNotFound
	}
	return roles, err
}

func (u *usecase) Login(ctx context.Context, email, password string) (*model.User, string, error) {
	user, hashedPassword, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		return nil, "", base.NewApiErrors("email", "This email is not registered")
	}

	if !utils.CompareHashAndPassword(hashedPassword, password) {
		return nil, "", base.NewApiErrors("password", "wrong password")
	}

	accessToken, err := helper.GenerateJWTToken(user)
	if err != nil {
		u.Error(ctx).Err(err).Msg("LoginError")
	}
	err = u.repo.UpdateTimeLastLogin(ctx, int64(user.ID))
	if err != nil {
		return user, accessToken, err
	}
	return user, accessToken, nil
}

// forgot password
func (u *usecase) ForgotPassword(ctx context.Context, mail string) (string, error) {
	userId, err := u.repo.GetUserId(ctx, mail)
	if err != nil {
		return "", err
	}
	if userId != 0 {
		accessToken, err := helper.GenerateJWTToken(userId)
		err = u.SendEmail(accessToken, mail, "email_template.html", "Reset Password")
		if err != nil {
			u.Info(ctx).Err(err).Msg("Can not send mail")
			return "", err
		}
	}

	return "Send Mail Successfully", nil
}
func New() Usecase {
	return &usecase{
		Usecase: base.NewBaseUsecase("auth"),
		repo:    auth.New(),
	}
}
