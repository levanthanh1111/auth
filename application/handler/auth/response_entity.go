package auth

import "github.com/tpp/msf/model"

type userRes struct {
	AccessToken string      `json:"access_token"`
	Data        *model.User `json:"data"`
}

// response email
type emailRes struct {
	Message string `json:"message"`
}
