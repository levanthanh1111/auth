package auth

import (
	"github.com/tpp/msf/shared/base"
	"github.com/tpp/msf/shared/validator"
)

var v = validator.Get()

type loginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (lr *loginReq) IsValid() error {
	return v.Struct(lr)
}

// request for forgot password
type sendMailReq struct {
	Email string `json:"email" validate:"required,email"`
}
type listParams = base.ListParams

var listAcceptedFilterKeys = []string{"full_name", "email", "org_name", "role_name", "role_id"}

// check for forgot password
func (lr *sendMailReq) IsValid() error {
	return v.Struct(lr)
}
