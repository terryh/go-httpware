# go-httpware
some http middleware just depend on http.Handler from standard library
please use Golang 1.7 or later version which include context in standard library
    
# Example to use middleware

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/NYTimes/gziphandler"
	"github.com/didip/tollbooth"
	"github.com/go-zoo/bone"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/terryh/go-httpware"
	"github.com/unrolled/secure"
)

func MyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Context().Value("db").(*mgo.Session))
	fmt.Fprintf(w, "Hello %v\n", bone.GetValue(r, "id"))
}

func timeoutHandler(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, 1*time.Second, "timed out")
}

func main() {

	// prepare db
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	// setup secure middleware
	var secureMiddleware *secure.Secure
	secureMiddleware = secure.New(secure.Options{
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		IsDevelopment:      true,
		//SSLRedirect:        true,
	})

	common := alice.New(
		httpware.SimpleLogger,
		httpware.Recovery,

		secureMiddleware.Handler,
		gziphandler.GzipHandler,
		cors.New(cors.Options{AllowedHeaders: []string{"*"}, AllowCredentials: true}).Handler,

		httpware.Limiter(tollbooth.NewLimiter(5, time.Minute)),
		httpware.Mongo(session, "db"),
		httpware.JWTCookie("token"),
		//httpware.JWTAuth("secret", "HS256", "token"),
	)

	mux := bone.New()

	mux.Get("/hello/:id", common.Then(http.HandlerFunc(MyHandler)))
	mux.Get("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":3000", mux)
}
```
