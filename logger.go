package negroni

import (
	"log"
	"net/http"
	"time"
)

// Logger returns a middleware handler that logs the request as it goes in and the response as it goes out.
func Logger() Handler {
	return HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next(rw, r)

		res := rw.(ResponseWriter)
		log.Printf("Completed %v %s in %v\n", res.Status(), http.StatusText(res.Status()), time.Since(start))
	})
}
