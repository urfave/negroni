package negroni

import (
	"log"
	"net/http"
	"os"
)

// Handler handler is an interface that objects can implement to be registered to serve as middleware
// in the Negroni middleware stack.
// ServeHTTP should yield to the next middleware in the chain by invoking the next http.HandlerFunc
// passed in.
type Handler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// HandlerFunc is an adapter to allow the use of ordinary functions as Negroni handlers.
// If f is a function with the appropriate signature, HandlerFunc(f) is a Handler object that calls f.
type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h(rw, r, next)
}

type middleware struct {
	handler Handler
	next    *middleware
}

func (h middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	res := rw.(ResponseWriter)
	if !res.Written() {
		h.handler.ServeHTTP(rw, r, h.next.ServeHTTP)
	}
}

// Wrap converts a http.Handler into a negroni.Handler so it can be used as a Negroni
// middleware. The next http.HandlerFunc is automatically called after the Handler
// is executed.
func Wrap(handler http.Handler) Handler {
	return HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(rw, r)
		next(rw, r)
	})
}

// Negroni is a stack of Middleware Handlers that can be invoked as an http.Handler.
// Negroni middleware is evaluated in the order that they are added to the stack using
// the Use and UseHandler methods.
type Negroni struct {
	middleware middleware
	handlers   []Handler
}

// New returns a new Negroni instance with no middleware preconfigured.
func New() *Negroni {
	return &Negroni{
		middleware: middleware{HandlerFunc(notFoundHandler), &middleware{}},
	}
}

// Classic returns a new Negroni instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging
// Static - Static File Serving
func Classic() *Negroni {
	n := New()
	n.Use(NewRecovery())
	n.Use(NewLogger())
	n.Use(NewStatic("public"))
	return n
}

func (n *Negroni) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	n.middleware.ServeHTTP(NewResponseWriter(rw), r)
}

// Use adds a Handler onto the middleware stack. Handlers are invoked in the order they are added to a Negroni.
func (n *Negroni) Use(handler Handler) {
	n.handlers = append(n.handlers, handler)
	n.middleware = build(0, n.handlers)
}

// UseHandler adds a http.Handler onto the middleware stack. Handlers are invoked in the order they are added to a Negroni.
func (n *Negroni) UseHandler(handler http.Handler) {
	n.Use(Wrap(handler))
}

// Run is a convenience function that runs the negroni stack as a http server. Run looks for the PORT and HOST
// environment variables for configuration. Otherwise Run will run the server at :3000
func (n *Negroni) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	host := os.Getenv("HOST")

	l := log.New(os.Stdout, "[negroni] ", 0)
	l.Printf("listening on %s:%s\n", host, port)
	l.Fatalln(http.ListenAndServe(host+":"+port, n))

}

func build(i int, handlers []Handler) middleware {
	var next middleware

	h := handlers[i]
	if i < len(handlers)-1 {
		next = build(i+1, handlers)
	} else {
		next = middleware{HandlerFunc(notFoundHandler), &middleware{}}
	}

	return middleware{h, &next}
}

func notFoundHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	http.NotFound(rw, r)
}
