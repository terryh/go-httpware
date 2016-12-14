package httpware

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
)

// PostgresDB middleware
func PostgresDB(db *sqlx.DB, contextKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var ctx context.Context
			var newReq *http.Request

			if contextKey == "" {
				contextKey = "db"
			}
			ctx = context.WithValue(r.Context(), contextKey, db)
			newReq = r.WithContext(ctx)

			// Process request
			next.ServeHTTP(w, newReq)
		}
		return http.HandlerFunc(fn)
	}
}
