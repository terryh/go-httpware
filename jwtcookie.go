package httpware

import (
	"fmt"
	"net/http"
)

// JWTCookie,  A JSON Web Token middleware
// this will parse token from Query String or cookie,
// put jwt token to Authorization token
// which not secure but convenience
func JWTCookie(tokenName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Skip WebSocket
			header := r.Header
			if (header.Get("Upgrade")) == "websocket" {
				// Process request
				next.ServeHTTP(w, r)
				return
			}

			// have not Authorization header
			if header.Get("Authorization") == "" {
				//  try Query String
				if tokenHash := r.URL.Query().Get(tokenName); tokenHash != "" {
					header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenHash))
				} else {
					// try Cookie
					token, err := r.Cookie(tokenName)
					if err == nil {
						header.Set("Authorization", fmt.Sprintf("Bearer %v", token.Value))
					}
				}
			}
			// Process request
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
