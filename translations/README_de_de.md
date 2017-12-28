# Negroni [![GoDoc](https://godoc.org/github.com/urfave/negroni?status.svg)](http://godoc.org/github.com/urfave/negroni) [![wercker status](https://app.wercker.com/status/13688a4a94b82d84a0b8d038c4965b61/s "wercker status")](https://app.wercker.com/project/bykey/13688a4a94b82d84a0b8d038c4965b61)

Negroni ist ein Ansatz für eine idiomatische Middleware in Go. Sie ist klein, nicht-intrusiv und unterstützt die Nutzung von `net/http` Handlern.

Wenn Dir die Idee hinter [Martini](http://github.com/go-martini/martini) gefällt, aber Du denkst, es stecke zu viel Magie darin, dann ist Negroni eine passende Alternative.

## Wo fange ich an?

Nachdem Du Go installiert und den [GOPATH](http://golang.org/doc/code.html#GOPATH) eingerichtet hast, erstelle eine `.go`-Datei. Nennen wir sie  `server.go`.

~~~ go
package main

import (
  "github.com/urfave/negroni"
  "net/http"
  "fmt"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Willkommen auf der Homepage!")
  })

  n := negroni.Classic()
  n.UseHandler(mux)
  n.Run(":3000")
}
~~~

Installiere nun das Negroni Package (**go 1.1** und höher werden vorausgesetzt):
~~~
go get github.com/urfave/negroni
~~~

Dann starte Deinen Server:
~~~
go run server.go
~~~

Nun läuft ein `net/http`-Webserver von Go unter `localhost:3000`.

## Hilfe benötigt?
Wenn Du eine Frage hast oder Dir ein bestimmte Funktion wünscht, nutze die [Mailing Liste](https://groups.google.com/forum/#!forum/negroni-users). Issues auf Github werden ausschließlich für Bug Reports und Pull Requests genutzt.

## Ist Negroni ein Framework?
Negroni ist **kein** Framework. Es ist eine Bibliothek, geschaffen, um kompatibel mit `net/http` zu sein.

## Routing?
Negroni ist BYOR (Bring your own Router - Nutze Deinen eigenen Router). Die Go-Community verfügt bereits über eine Vielzahl von großartigen Routern. Negroni versucht möglichst alle zu unterstützen, indem es `net/http` vollständig unterstützt. Beispielsweise sieht eine Implementation mit [Gorilla Mux](http://github.com/gorilla/mux) folgendermaßen aus:

~~~ go
router := mux.NewRouter()
router.HandleFunc("/", HomeHandler)

n := negroni.New(Middleware1, Middleware2)
// Oder nutze eine Middleware mit der Use()-Funktion
n.Use(Middleware3)
// Der Router kommt als letztes
n.UseHandler(router)

n.Run(":3000")
~~~

## `negroni.Classic()`
`negroni.Classic()` stellt einige Standard-Middlewares bereit, die für die meisten Anwendungen von Nutzen ist:

* `negroni.Recovery` - Middleware für Panic Recovery .
* `negroni.Logging` - Anfrage/Rückmeldungs-Logging-Middleware.
* `negroni.Static` - Ausliefern von statischen Dateien unter dem "public" Verzeichnis.

Dies macht es wirklich einfach, mit den nützlichen Funktionen von Negroni zu starten.

## Handlers
Negroni stellt einen bidirektionalen Middleware-Flow bereit. Dies wird durch das `negroni.Handler`-Interface erreicht:

~~~ go
type Handler interface {
  ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
~~~

Wenn eine Middleware nicht bereits den ResponseWriter genutzt hat, sollte sie die nächste `http.HandlerFunc` in der Verkettung von Middlewares aufrufen und diese ausführen. Das kann von großem Nutzen sein:

~~~ go
func MyMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // Mache etwas vor dem Aufruf
  next(rw, r)
  // Mache etwas nach dem Aufruf
}
~~~

Und Du kannst eine Middleware durch die `Use`-Funktion der Verkettung von Middlewares zuordnen.

~~~ go
n := negroni.New()
n.Use(negroni.HandlerFunc(MyMiddleware))
~~~

Stattdessen kannst Du auch herkömmliche `http.Handler` zuordnen:

~~~ go
n := negroni.New()

mux := http.NewServeMux()
// Ordne Deine Routen zu

n.UseHandler(mux)

n.Run(":3000")
~~~

## `Run()`
Negroni hat eine nützliche Funktion namens `Run`. `Run` übernimmt eine Zeichenkette `addr` ähnlich wie [http.ListenAndServe](http://golang.org/pkg/net/http#ListenAndServe).

~~~ go
n := negroni.Classic()
// ...
log.Fatal(http.ListenAndServe(":8080", n))
~~~

## Routenspezifische Middleware
Wenn Du eine Gruppe von Routen hast, welche alle die gleiche Middleware ausführen müssen, kannst Du einfach eine neue Negroni-Instanz erstellen und sie als Route-Handler nutzen:

~~~ go
router := mux.NewRouter()
adminRoutes := mux.NewRouter()
// Füge die Admin-Routen hier hinzu

// Erstelle eine neue Negroni-Instanz für die Admin-Middleware
router.Handle("/admin", negroni.New(
  Middleware1,
  Middleware2,
  negroni.Wrap(adminRoutes),
))
~~~

## Middlewares von Dritten

Hier ist eine aktuelle Liste von Middlewares, die kompatible mit Negroni sind. Tue Dir keinen Zwang an, Dich einzutragen, wenn Du selbst eine Middleware programmiert hast:


| Middleware | Autor | Beschreibung |
| -----------|--------|-------------|
| [authz](https://github.com/casbin/negroni-authz) | [Yang Luo](https://github.com/hsluoyz) | ACL, RBAC, ABAC Authorization middlware based on [Casbin](https://github.com/casbin/casbin) |
| [binding](https://github.com/mholt/binding) | [Matt Holt](https://github.com/mholt) | Data binding from HTTP requests into structs |
| [cloudwatch](https://github.com/cvillecsteele/negroni-cloudwatch) | [Colin Steele](https://github.com/cvillecsteele) | AWS cloudwatch metrics middleware |
| [cors](https://github.com/rs/cors) | [Olivier Poitrey](https://github.com/rs) | [Cross Origin Resource Sharing](http://www.w3.org/TR/cors/) (CORS) support |
| [csp](https://github.com/awakenetworks/csp) | [Awake Networks](https://github.com/awakenetworks) | [Content Security Policy](https://www.w3.org/TR/CSP2/) (CSP) support |
| [delay](https://github.com/jeffbmartinez/delay) | [Jeff Martinez](https://github.com/jeffbmartinez) | Add delays/latency to endpoints. Useful when testing effects of high latency |
| [New Relic Go Agent](https://github.com/yadvendar/negroni-newrelic-go-agent) | [Yadvendar Champawat](https://github.com/yadvendar) | Official [New Relic Go Agent](https://github.com/newrelic/go-agent) (currently in beta)  |
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
| [mgo session](https://github.com/joeljames/nigroni-mgo-session) | [Joel James](https://github.com/joeljames) | Middleware that handles creating and closing mgo sessions per request |
| [digits](https://github.com/bamarni/digits) | [Bilal Amarni](https://github.com/bamarni) | Middleware that handles [Twitter Digits](https://get.digits.com/) authentication |
| [stats](https://github.com/guptachirag/stats) | [Chirag Gupta](https://github.com/guptachirag/stats) | Middleware that manages qps and latency stats for your endpoints and asynchronously flushes them to influx db |

## Beispiele
[Alexander Rødseth](https://github.com/xyproto) programmierte [mooseware](https://github.com/xyproto/mooseware), ein  Grundgerüst zum Erstellen von Negroni Middleware-Handerln.

## Aktualisieren in Echtzeit?
[gin](https://github.com/codegangsta/gin) und [fresh](https://github.com/pilu/fresh)  aktualisieren Deine Negroni-Anwendung automatisch.

## Unverzichbare Informationen für Go- & Negronineulinge

* [Nutze einen Kontext zum Übertragen von Middlewareinformationen an Handler (Englisch)](http://elithrar.github.io/article/map-string-interface/)
* [Middlewares verstehen (Englisch)](http://mattstauffer.co/blog/laravel-5.0-middleware-replacing-filters)

## Über das Projekt

Negroni wurde obsseziv von Niemand gerigeren als dem [Code Gangsta](http://codegangsta.io/) entwickelt.
