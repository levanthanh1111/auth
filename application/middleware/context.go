package middleware

import (
	"net/http"

	"github.com/rs/xid"
	"github.com/tpp/msf/shared/context"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.FromBaseContext(r.Context())
		reqID := xid.New().String()
		ctx.WithReqID(reqID)
		w.Header().Set("Request-Id", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
