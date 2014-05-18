package negroni

import (
	"net/http"
)

type Handler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

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

func Wrap(handler http.Handler) Handler {
	return HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(rw, r)
		next(rw, r)
	})
}

type Negroni struct {
	middleware middleware
	handlers   []Handler
}

func New() *Negroni {
	n := Negroni{
		middleware: middleware{HandlerFunc(voidHandler), &middleware{}},
	}
	return &n
}

func (n *Negroni) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	n.middleware.ServeHTTP(NewResponseWriter(rw), r)
}

func (n *Negroni) Use(handler Handler) {
	n.handlers = append(n.handlers, handler)
	n.middleware = build(0, n.handlers)
}

func (n *Negroni) UseHandler(handler http.Handler) {
	n.Use(Wrap(handler))
}

func build(i int, handlers []Handler) middleware {
	var next middleware

	h := handlers[i]
	if i < len(handlers)-1 {
		next = build(i+1, handlers)
	} else {
		next = middleware{HandlerFunc(voidHandler), &middleware{}}
	}

	return middleware{h, &next}
}

func voidHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// do nothing
}
