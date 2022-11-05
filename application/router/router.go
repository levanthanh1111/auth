package router

import (
	"net/http"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/tpp/msf/application/handler"
	"github.com/tpp/msf/application/middleware"
)

var Router *chi.Mux

func init() {
	Router = chi.NewRouter()
}

var (
	cors   = middleware.CORS
	auth   = middleware.Auth
	reqID  = middleware.RequestID
	permit = middleware.Permit
)

func SetupHandler(h handler.Handler) {
	Router.Use(cors, chimiddleware.Logger, reqID)
	apidocsHTTPHandler(Router)

	Router.Route("/users", func(r chi.Router) {

		r.With(auth, permit([]string{"VIEW_LIST_USER"})).Get("/", h.Users().List)
		r.With(auth, permit([]string{"VIEW_CURRENT_USER"})).Get("/{id:[0-9]+}", h.Users().Get)
		r.With(auth, permit([]string{""})).Post("/create", h.Users().Create)
		r.With(auth, permit([]string{""})).Post("/create-role", h.Users().CreateRole)
		r.With(auth, permit([]string{""})).Put("/assign-role", h.Users().AssignRole)
		r.With(auth, permit([]string{""})).Put("/is-active", h.Users().UpdateActive)
		r.With(auth).Put("/reset-password", h.Users().UpdatePassWord)
		r.With(auth).Put("/name", h.Users().UpdateName)
		r.With(auth, permit([]string{""})).Put("/admin-reset-password", h.Users().AdminResetPWForUser)
	})

	Router.Route("/contracts", func(r chi.Router) {

		r.With(auth, permit([]string{"VIEW_ALL_CONTRACT_LIST", "VIEW_CONTRACT_LIST"})).Get("/", h.Contracts().List)
		r.With(auth, permit([]string{"VIEW_ALL_CONTRACT_LIST", "VIEW_CONTRACT_LIST"})).Get("/{id:[0-9]+}", h.Contracts().Get)

	})

	Router.Route("/auth", func(r chi.Router) {

		r.Post("/login", h.Auth().Login)
		r.Post("/forgot-password", h.Auth().ForgotPassword)
		r.With(auth, permit([]string{"VIEW_LIST_USER"})).Get("/roles", h.Auth().ListRoles)
		r.With(auth, permit([]string{"VIEW_LIST_USER"})).Get("/roles/{id:[0-9]+}", h.Auth().Get)
		r.With(auth, permit([]string{""})).Get("/permissions", h.Auth().ListPermissions)

	})

}

func apidocsHTTPHandler(route *chi.Mux) {
	route.Handle("/docs", http.StripPrefix("/docs", http.FileServer(http.Dir("docs"))))
}
