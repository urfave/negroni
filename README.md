# Negroni [![GoDoc](https://godoc.org/github.com/codegangsta/negroni?status.png)](http://godoc.org/github.com/codegangsta/negroni)

Negroni is a fancy approach to web middleware in Go. It is tiny, non-intrusive, and encourages use of `net/http` Handlers.

## Getting Started

After installing Go and setting up your [GOPATH](http://golang.org/doc/code.html#GOPATH), create your first `.go` file. We'll call it `server.go`.

~~~ go
package main

import (
  "github.com/codegangsta/negroni"
  "net/http"
)

func main() {
  n := negroni.Classic()
  n.Use(negroni.HandlerFunc(Hello))
  n.Run()
}

func Hello(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  rw.Write([]byte("Hello world"))
}
~~~

Then install the Negroni package (**go 1.1** and greater is required):
~~~
go get github.com/codegangsta/negroni
~~~

Then run your server:
~~~
go run server.go
~~~

You will now have a Go net/http webserver running on `localhost:3000`.

## Handlers
Negroni provides a bidirectional middleware flow. This is done through the `negroni.Handler` interface:

~~~ go
type Handler interface {
  ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
~~~

Middleware should call the next `http.HandlerFunc` in the chain to yield to the next middleware handler. This can be used for great good:

~~~ go
func MyMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // do some stuff before
  next(rw, r)
  // do some stuff after
}
~~~

## Live code reload?
[gin](https://github.com/codegangsta/gin) and [fresh](https://github.com/pilu/fresh) both live reload negroni apps.

## Contributing
Martini is meant to be kept tiny and clean. Most contributions should end up in a repository in the [martini-contrib](https://github.com/martini-contrib) organization. If you do have a contribution for the core of Martini feel free to put up a Pull Request.

## About

Inspired by [express](https://github.com/visionmedia/express) and [sinatra](https://github.com/sinatra/sinatra)

Martini is obsessively designed by none other than the [Code Gangsta](http://codegangsta.io/)
