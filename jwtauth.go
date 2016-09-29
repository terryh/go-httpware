package httpware

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"github.com/dgrijalva/jwt-go"
)

const (
	bearer_ = "Bearer "
)

// Auth http middleware for JWT authentication
// Auth(secret,signMethod, contextKey string)
// secret string ....
// signingMethod string name for sign method
// contextKey string set value at golang 1.7 http context
func JWTAuth(secret, signingMethod, contextKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var auth, authhead string
			// header we want
			authhead = r.Header.Get("Authorization")

			if strings.HasPrefix(authhead, bearer_) {
				auth = strings.Split(authhead, bearer_)[1]
			} else {
				http.Error(w, "invalid jwt header", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(auth, func(t *jwt.Token) (interface{}, error) {
				// Check the signing method
				if t.Method.Alg() != signingMethod {
					return nil, fmt.Errorf("unexpected jwt signing method %s", t.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if contextKey != "" {
				// set to golang 1.7 context
				// do have context Key name
				ctx := context.WithValue(r.Context(), contextKey, token)
				newReq := r.WithContext(ctx)
				// Process request
				next.ServeHTTP(w, newReq)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
