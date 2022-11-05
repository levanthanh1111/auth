package middleware

import (
	"github.com/tpp/msf/model"
	"net/http"
	"strings"

	"github.com/tpp/msf/config"
	"github.com/tpp/msf/domain/usecase/users"
	"github.com/tpp/msf/external-adapter/db"
	"github.com/tpp/msf/shared/auth"
	"github.com/tpp/msf/shared/context"
	"github.com/tpp/msf/shared/log"
)

func Auth(next http.Handler) http.Handler {
	var authUsecase = users.New()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := context.FromBaseContext(r.Context())
		logger := log.Logger.With().Str("req_id", ctx.ReqID()).Logger()

		jwtHeader := config.GetConfig[string]("auth.header")
		tokenSchema := config.GetConfig[string]("auth.token_schema")
		if tokenSchema == "" {
			tokenSchema = "Bearer"
		}
		tokenString := r.Header.Get(jwtHeader)

		// Check token string with format
		if !strings.HasPrefix(strings.ToLower(tokenString), strings.ToLower(tokenSchema)+" ") {
			logger.Error().Str("error", "wrong token chema format").Msg("AuthError")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString = strings.Split(tokenString, " ")[1]
		ID, err := auth.ParseJWTToken(tokenString)
		if err != nil {
			logger.Error().Err(err).Msg("AuthError")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var dbc = db.GetDBInstance()
		var user *model.User
		if user, err = authUsecase.GetUser(context.Background().WithDBTx(dbc), uint64(ID)); err != nil || ID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx.WithUser(user)))
	})
}
