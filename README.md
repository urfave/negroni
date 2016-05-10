# Negroni [![GoDoc](https://godoc.org/github.com/codegangsta/negroni?status.svg)](http://godoc.org/github.com/codegangsta/negroni) [![wercker status](https://app.wercker.com/status/13688a4a94b82d84a0b8d038c4965b61/s "wercker status")](https://app.wercker.com/project/bykey/13688a4a94b82d84a0b8d038c4965b61) [![codebeat](https://codebeat.co/badges/47d320b1-209e-45e8-bd99-9094bc5111e2)](https://codebeat.co/projects/github-com-codegangsta-negroni)

Negroni is an idiomatic approach to web middleware in Go. It is tiny,
non-intrusive, and encourages use of `net/http` Handlers.

If you like the idea of [Martini](https://github.com/go-martini/martini), but
you think it contains too much magic, then Negroni is a great fit.

Language Translations:
* [German (de_DE)](translations/README_de_de.md)
* [Português Brasileiro (pt_BR)](translations/README_pt_br.md)
* [简体中文 (zh_cn)](translations/README_zh_cn.md)
* [繁體中文 (zh_tw)](translations/README_zh_tw.md)

## Getting Started

After installing Go and setting up your
[GOPATH](http://golang.org/doc/code.html#GOPATH), create your first `.go` file.
We'll call it `server.go`.

<!-- { "interrupt": true } -->
``` go
package main

import (
  "fmt"
  "net/http"

  "github.com/codegangsta/negroni"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!")
  })

  n := negroni.Classic() // Includes some default middlewares
  n.UseHandler(mux)

  http.ListenAndServe(":3000", n)
}
```

Then install the Negroni package (&gt;= **go 1.1** is required):

```
go get github.com/codegangsta/negroni
```

Then run your server:

```
go run server.go
```

You will now have a Go `net/http` webserver running on `localhost:3000`.

## Is Negroni a Framework?

Negroni is **not** a framework. It is a middleware-focused library that is
designed to work directly with net/http.

## Routing?

Negroni is BYOR (Bring your own Router). The Go community already has a number
of great http routers available, and Negroni tries to play well with all of them
by fully supporting `net/http`. For instance, integrating with [Gorilla Mux]
looks like so:

``` go
router := mux.NewRouter()
router.HandleFunc("/", HomeHandler)

n := negroni.New(Middleware1, Middleware2)
// Or use a middleware with the Use() function
n.Use(Middleware3)
// router goes last
n.UseHandler(router)

http.ListenAndServe(":3000", n)
```

## `negroni.Classic()`

`negroni.Classic()` provides some default middleware that is useful for most
applications:

* [`negroni.Recovery`](#negronirecovery) - Panic Recovery Middleware.
* [`negroni.Logger`](#negronilogger) - Request/Response Logger Middleware.
* [`negroni.Static`](#negronistatic) - Static File serving under the "public"
  directory.

This makes it really easy to get started with some useful features from Negroni.

### `negroni.Recovery`

The default recovery middleware, returning status `500` and logging `PANIC: `
for any panic recovered when calling `next()`.  All struct members are
configurable:

* `Logger` - defaults to a new `log.Logger` prefixed `[negroni]` writing to
  `os.Stdout` with no flags set.
* `PrintStack` - default `true`, causing the error to be written to the
  `http.ResponseWriter`
* `StackSize` - default 8k, used to make a byte slice receiver passed to
  `runtime.Stack` when retrieving the stack for logging during recovery
* `StackAll` - default `false`, passed to `runtime.Stack` when retrieving the
  stack for logging during recovery

### `negroni.Logger`

The default logging middleware, writing `[negroni] Started {method} {path}` and
`[negroni] Completed {status} {text-status} in {time}` to `os.Stdout` for every
request.  There are no configurable struct members.

### `negroni.Static`

The default static file serving middleware, meant as a best-effort file serving
handler for `GET` or `HEAD` requests matching paths within a single
[`http.FileSystem`].  All struct members are configurable:

* `Dir` - required argument to `NewStatic` of type [`http.FileSystem`], defaulting
  to `http.Dir("public")` when constructed via `negroni.Classic()`
* `Prefix` - default `""` optional prefix which is stripped (if non-empty) prior
  to file existence check.
* `IndexFile` - default `"index.html"` file which will be served if the request
  path matches an existing directory.

**NOTE**: Nonexistent files will *not* result in a `404` status.  If this
behavior is desired, please consider using
[`http.FileServer`](https://godoc.org/net/http#FileServer) directly.

## Handlers

Negroni provides a bidirectional middleware flow. This is done through the
`negroni.Handler` interface:

``` go
type Handler interface {
  ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
```

If a middleware hasn't already written to the `ResponseWriter`, it should call
the next `http.HandlerFunc` in the chain to yield to the next middleware
handler.  This can be used for great good:

``` go
func MyMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // do some stuff before
  next(rw, r)
  // do some stuff after
}
```

And you can map it to the handler chain with the `Use` function:

``` go
n := negroni.New()
n.Use(negroni.HandlerFunc(MyMiddleware))
```

You can also map plain old `http.Handler`s:

``` go
n := negroni.New()

mux := http.NewServeMux()
// map your routes

n.UseHandler(mux)

http.ListenAndServe(":3000", n)
```

## `Run()`

Negroni has a convenience function called `Run`. `Run` takes an addr string
identical to [`http.ListenAndServe`](https://godoc.org/net/http#ListenAndServe).

<!-- { "interrupt": true } -->
``` go
package main

import (
  "github.com/codegangsta/negroni"
)

func main() {
  n := negroni.Classic()
  n.Run(":8080")
}
```

In general, you will want to use `net/http` methods and pass `negroni` as a
`Handler`, as this is more flexible, e.g.:

<!-- { "interrupt": true } -->
``` go
package main

import (
  "fmt"
  "log"
  "net/http"
  "time"

  "github.com/codegangsta/negroni"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!")
  })

  n := negroni.Classic() // Includes some default middlewares
  n.UseHandler(mux)

  s := &http.Server{
    Addr:           ":8080",
    Handler:        n,
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }
  log.Fatal(s.ListenAndServe())
}
```

## Route Specific Middleware

If you have a route group of routes that need specific middleware to be
executed, you can simply create a new Negroni instance and use it as your route
handler.

``` go
router := mux.NewRouter()
adminRoutes := mux.NewRouter()
// add admin routes here

// Create a new negroni for the admin middleware
router.PathPrefix("/admin").Handler(negroni.New(
  Middleware1,
  Middleware2,
  negroni.Wrap(adminRoutes),
))
```

If you are using [Gorilla Mux], here is an example using a subrouter:

``` go
router := mux.NewRouter()
subRouter := mux.NewRouter().PathPrefix("/subpath").Subrouter().StrictSlash(true)
subRouter.HandleFunc("/", someSubpathHandler) // "/subpath/"
subRouter.HandleFunc("/:id", someSubpathHandler) // "/subpath/:id"

// "/subpath" is necessary to ensure the subRouter and main router linkup
router.PathPrefix("/subpath").Handler(negroni.New(
  Middleware1,
  Middleware2,
  negroni.Wrap(subRouter),
))
```

## Third Party Middleware

Here is a current list of Negroni compatible middlware. Feel free to put up a PR
linking your middleware if you have built one:

| Middleware | Author | Description |
| -----------|--------|-------------|
| [binding](https://github.com/mholt/binding) | [Matt Holt](https://github.com/mholt) | Data binding from HTTP requests into structs |
| [cors](https://github.com/rs/cors) | [Olivier Poitrey](https://github.com/rs) | [Cross Origin Resource Sharing](http://www.w3.org/TR/cors/) (CORS) support |
| [delay](https://github.com/jeffbmartinez/delay) | [Jeff Martinez](https://github.com/jeffbmartinez) | Add delays/latency to endpoints. Useful when testing effects of high latency |
| [gorelic](https://github.com/jingweno/negroni-gorelic) | [Jingwen Owen Ou](https://github.com/jingweno) | New Relic agent for Go runtime |
| [Graceful](https://github.com/tylerb/graceful) | [Tyler Bunnell](https://github.com/tylerb) | Graceful HTTP Shutdown |
| [gzip](https://github.com/phyber/negroni-gzip) | [phyber](https://github.com/phyber) | GZIP response compression |
| [JWT Middleware](https://github.com/auth0/go-jwt-middleware) | [Auth0](https://github.com/auth0) | Middleware checks for a JWT on the `Authorization` header on incoming requests and decodes it|
| [logrus](https://github.com/meatballhat/negroni-logrus) | [Dan Buch](https://github.com/meatballhat) | Logrus-based logger |
| [oauth2](https://github.com/goincremental/negroni-oauth2) | [David Bochenski](https://github.com/bochenski) | oAuth2 middleware |
| [onthefly](https://github.com/xyproto/onthefly) | [Alexander Rødseth](https://github.com/xyproto) | Generate TinySVG, HTML and CSS on the fly |
| [permissions2](https://github.com/xyproto/permissions2) | [Alexander Rødseth](https://github.com/xyproto) | Cookies, users and permissions |
| [prometheus](https://github.com/zbindenren/negroni-prometheus) | [Rene Zbinden](https://github.com/zbindenren) | Easily create metrics endpoint for the [prometheus](http://prometheus.io) instrumentation tool |
| [render](https://github.com/unrolled/render) | [Cory Jacobsen](https://github.com/unrolled) | Render JSON, XML and HTML templates |
| [RestGate](https://github.com/pjebs/restgate) | [Prasanga Siripala](https://github.com/pjebs) | Secure authentication for REST API endpoints |
| [secure](https://github.com/unrolled/secure) | [Cory Jacobsen](https://github.com/unrolled) | Middleware that implements a few quick security wins |
| [sessions](https://github.com/goincremental/negroni-sessions) | [David Bochenski](https://github.com/bochenski) | Session Management |
| [stats](https://github.com/thoas/stats) | [Florent Messa](https://github.com/thoas) | Store information about your web application (response time, etc.) |
| [VanGoH](https://github.com/auroratechnologies/vangoh) | [Taylor Wrobel](https://github.com/twrobel3) | Configurable [AWS-Style](http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html) HMAC authentication middleware |
| [xrequestid](https://github.com/pilu/xrequestid) | [Andrea Franz](https://github.com/pilu) | Middleware that assigns a random X-Request-Id header to each request |

## Examples

[Alexander Rødseth](https://github.com/xyproto) created
[mooseware](https://github.com/xyproto/mooseware), a skeleton for writing a
Negroni middleware handler.

## Live code reload?

[gin](https://github.com/codegangsta/gin) and
[fresh](https://github.com/pilu/fresh) both live reload negroni apps.

## Essential Reading for Beginners of Go & Negroni

* [Using a Context to pass information from middleware to end handler](http://elithrar.github.io/article/map-string-interface/)
* [Understanding middleware](https://mattstauffer.co/blog/laravel-5.0-middleware-filter-style)

## About

Negroni is obsessively designed by none other than the [Code
Gangsta](https://codegangsta.io/)

[Gorilla Mux]: https://github.com/gorilla/mux
[`http.FileSystem`]: https://godoc.org/net/http#FileSystem
