package middleware

import (
	"Lighthouse/internal/database/spec/interfaces"
	"context"
	"net/http"
)

var CtxDbKey = struct{}{}

func ApplyAttachDb(db interfaces.DatabaseConnectorStrategy) ApplyMiddlewareLayer {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), CtxDbKey, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetDbFromContext(ctx context.Context) (interfaces.DatabaseConnectorStrategy, bool) {
	db, ok := ctx.Value(CtxDbKey).(interfaces.DatabaseConnectorStrategy)
	return db, ok
}
