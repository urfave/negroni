// Package negroni is a lightweight library that extends the
// builtin net/http package to enable the use of middlewares.
package negroni

import (
	"log"
	"net/http"
	"os"
)

// Node (linked-list)
type middleware struct {
	handler Handler
	next    *middleware
}

// Queue (linked-list)
type Negroni struct {
	first *middleware
	last  *middleware
}

// New registers optional middlewares that implement
// http.Handler or negroni.Handler and returns a new
// Negroni.
func New(handlers ...interface{}) *Negroni {
	c := &Negroni{emptyMiddleware(), nil}
	for _, handler := range handlers {
		c.Use(handler)
	}
	return c
}

// Classic returns a new Negroni instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging
// Static - Static File Serving
func Classic() *Negroni {
	return New(NewRecovery(), NewLogger(), NewStatic(http.Dir("public")))
}

// Use registers a handler that implements http.Handler
// or negroni.Handler.
func (c *Negroni) Use(handler interface{}) {
	switch handler.(type) {
	case Handler:
		c.load(handler.(Handler))
	case http.Handler:
		c.load(wrap(handler.(http.Handler)))
	}
}

// Run takes a network address and calls http.ListenAndServe.
func (c *Negroni) Run(addr string) {
	l := log.New(os.Stdout, "[negroni] ", 0)
	l.Printf("listening on %s", addr)
	l.Fatal(http.ListenAndServe(addr, c))
}

// ServeHTTP is implemented by Negroni so it can be called by
// the net/http package in order to do its own thing.
func (c *Negroni) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.first.mServeHTTP(NewResponseWriter(w), r)
}

// mServeHTTP is recursively called after the current middleware's
// handler is processed.
func (m *middleware) mServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(w, r, m.next.mServeHTTP)
}

// Handler is an interface that can be implemented by middlewares.
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, http.HandlerFunc)
}

// HandlerFunc accepts a function with the following parameters to
// create a HandlerFunc that can be used to implement the negroni.Handler
// interface.
type HandlerFunc func(http.ResponseWriter, *http.Request, http.HandlerFunc)

// ServeHTTP is an implementation of negroni.Handler's interface
func (fn HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fn(w, r, next)
}

// For a negroni to not be empty, it must have at least one
// middleware whose next isn't nil
func (c *Negroni) isEmpty() bool {
	return c.first == nil || c.last == nil
}

// load registers a middleware that implements http.Handler or
// negroni.Handler to Negroni.
func (c *Negroni) load(handler Handler) {
	middleware := &middleware{handler, emptyMiddleware()}
	if c.isEmpty() {
		c.first = middleware
		c.last = c.first
	} else {
		oldlast := c.last
		c.last = middleware
		oldlast.next = c.last
	}
}

// wrap takes an http.Handler and transforms it to a negroni.HandlerFunc
// object.
func wrap(handler http.Handler) Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(w, r)
		next(w, r)
	})
}

// emptyMiddleware is always the last middleware to be processed.
func emptyMiddleware() *middleware {
	return &middleware{HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}), nil}
}
