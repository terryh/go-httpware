package httpware

import (
	"net/http"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
)

// Mongo middleware which put *mgo.Session at golang 1.7 context in name 'contextKey'
// dbsession *mgo.Sessoion
// contextKey string, if contextKey == "", will use "db" for the contextKey
// there are many to wrap db connection, global variable or custom struct
// wrap with db connect
func Mongo(dbsession *mgo.Session, contextKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sCopy := dbsession.Copy()
			defer sCopy.Close()

			var ctx context.Context
			var newReq *http.Request

			if contextKey == "" {
				ctx = context.WithValue(r.Context(), "db", sCopy)
				newReq = r.WithContext(ctx)
			} else {
				ctx = context.WithValue(r.Context(), contextKey, sCopy)
				newReq = r.WithContext(ctx)
			}
			// Process request
			next.ServeHTTP(w, newReq)
		}
		return http.HandlerFunc(fn)
	}
}
