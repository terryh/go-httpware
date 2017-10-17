package httpware

import (
	"net/http"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
)

// Limiter middleware which use "github.com/didip/tollbooth"
func Limiter(limiterconfig *limiter.Limiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return tollbooth.LimitHandler(limiterconfig, next)
	}
}
