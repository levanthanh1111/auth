package handler

import (
	"github.com/tpp/msf/application/handler/auth"
	"github.com/tpp/msf/application/handler/contracts"
	"github.com/tpp/msf/application/handler/users"
)

type Handler interface {
	Users() users.Handler
	Contracts() contracts.Handler
	Auth() auth.Handler
}

type handler struct {
	user     users.Handler
	contract contracts.Handler
	auth     auth.Handler
}

func (h *handler) Users() users.Handler {
	return h.user
}
func (h *handler) Contracts() contracts.Handler {
	return h.contract
}
func (h *handler) Auth() auth.Handler {
	return h.auth
}

func New() Handler {
	return &handler{
		user:     users.New(),
		contract: contracts.New(),
		auth:     auth.New(),
	}
}
