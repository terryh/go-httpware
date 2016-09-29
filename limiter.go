package httpware

import (
	"net/http"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/config"
)

// Limiter middleware which use "github.com/didip/tollbooth"
func Limiter(limiterconfig *config.Limiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return tollbooth.LimitHandler(limiterconfig, next)
	}
}
