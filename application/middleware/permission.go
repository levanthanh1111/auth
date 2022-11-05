package middleware

import (
	"net/http"

	"github.com/tpp/msf/model"
	"github.com/tpp/msf/shared/context"
)

func Permit(permissions []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := context.FromBaseContext(r.Context())
			user := ctx.User()
			if user == nil {
				responseUnauthorizedError(w)
			} else if !hasPermisson(user, permissions) {
				responseForbiddenError(w)
			} else {
				next.ServeHTTP(w, r)
			}

		})
	}
}

func responseUnauthorizedError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func responseForbiddenError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}

func hasPermisson(user *model.User, permissions []string) bool {
	if user.IsAdmin {
		return true
	}

	userPermissions := getUserPermissons(user)
	for _, p := range permissions {
		if isElementExist(userPermissions, p) {
			return true
		}
	}
	return false
}

func getUserPermissons(user *model.User) []string {
	var userPermissions []*model.Permission
	if len(user.Roles) > 0 {
		userPermissions = user.Roles[0].Permissions
	}

	var permissonReturn []string
	for _, p := range userPermissions {
		permissonReturn = append(permissonReturn, p.Name)
	}
	return permissonReturn
}

func isElementExist(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
