# Negroni
[![GoDoc](https://godoc.org/github.com/urfave/negroni?status.svg)](http://godoc.org/github.com/urfave/negroni)
[![Build Status](https://travis-ci.org/urfave/negroni.svg?branch=master)](https://travis-ci.org/urfave/negroni)
[![codebeat](https://codebeat.co/badges/47d320b1-209e-45e8-bd99-9094bc5111e2)](https://codebeat.co/projects/github-com-urfave-negroni)
[![codecov](https://codecov.io/gh/urfave/negroni/branch/master/graph/badge.svg)](https://codecov.io/gh/urfave/negroni)

**Notice:** Esta é uma biblioteca conhecida anteriormente como 
`github.com/codegangsta/negroni` -- Github irá redirecionar automaticamente os pedidos
para este repositório, mas recomendas atualizar suas referências por clareza.

Negroni é uma abordagem idiomática para middleware web em Go. É pequeno, não intrusivo, e incentiva uso da biblioteca `net/http`.

Se gosta da idéia do [Martini](http://github.com/go-martini/martini), mas acha que contém muita mágica, então Negroni é ideal.

Idiomas traduzidos:
* [Deutsch (de_DE)](translations/README_de_de.md)
* [Português Brasileiro (pt_BR)](translations/README_pt_br.md)
* [简体中文 (zh_cn)](translations/README_zh_cn.md)
* [繁體中文 (zh_tw)](translations/README_zh_tw.md)
* [日本語 (ja_JP)](translations/README_ja_JP.md)

## Começando

Depois de instalar Go e definir seu [GOPATH](http://golang.org/doc/code.html#GOPATH), criar seu primeiro arquivo `.go`. Iremos chamá-lo `server.go`.

<!-- { "interrupt": true } -->
``` go
package main

import (
  "fmt"
  "net/http"

  "github.com/urfave/negroni"
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

Depois instale o pacote Negroni (**go 1.1** ou superior)

```
go get github.com/urfave/negroni
```

Depois execute seu servidor:
```
go run server.go
```

Agora terá um servidor web Go net/http rodando em `localhost:3000`.

## Empacotamento
Se você está no Debian, `negroni` também está disponível como  [um pacote](https://packages.debian.org/sid/golang-github-urfave-negroni-dev) que 
você podee instalar via `apt install golang-github-urfave-negroni-dev` (no momento queisto for escrito, ele estava em repositórios do `sid`).

## Negroni é um Framework?
Negroni **não** é a framework. É uma biblioteca que é desenhada para trabalhar diretamente com net/http.

## Roteamento?
Negroni é TSPR(Traga seu próprio Roteamento). A comunidade Go já tem um grande número de roteadores http disponíveis, Negroni tenta rodar bem com todos eles pelo suporte total `net/http`/ Por exemplo, a integração com [Gorilla Mux](http://github.com/gorilla/mux) se parece com isso:

```go
router := mux.NewRouter()
router.HandleFunc("/", HomeHandler)

n := negroni.New(Middleware1, Middleware2)
// Or use a middleware with the Use() function
n.Use(Middleware3)
// router goes last
n.UseHandler(router)

n.Run(":3000")
```

## `negroni.Classic()`
`negroni.Classic()`  fornece alguns middlewares padrão que são úteis para maioria das aplicações:

* `negroni.Recovery` - Panic Recovery Middleware.
* `negroni.Logging` - Request/Response Logging Middleware.
* `negroni.Static` - Static File serving under the "public" directory.

Isso torna muito fácil começar com alguns recursos úteis do Negroni.

## Handlers
Negroni fornece um middleware de fluxo bidirecional. Isso é feito através da interface `negroni.Handler`:

```go
type Handler interface {
  ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
```

Se um middleware não tenha escrito o ResponseWriter, ele deve chamar a próxima `http.HandlerFunc` na cadeia para produzir o próximo handler middleware. Isso pode ser usado muito bem:

```go
func MyMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // do some stuff before
  next(rw, r)
  // do some stuff after
}
```

E pode mapear isso para a cadeia de handler com a função `Use`:

```go
n := negroni.New()
n.Use(negroni.HandlerFunc(MyMiddleware))
```

Você também pode mapear `http.Handler` antigos:

```go
n := negroni.New()

mux := http.NewServeMux()
// map your routes

n.UseHandler(mux)

http.ListenAndServe(":3000", n)
```

## `With()`

Negroni tem uma conveniente função chamada `With`. `With` pega uma ou mais
instâncias `Handler` e retorna um novo `Negroni` com a combinação  dos handlers de receivers e os novos handlers.

```go
// middleware we want to reuse
common := negroni.New()
common.Use(MyMiddleware1)
common.Use(MyMiddleware2)

// `specific` is a new negroni with the handlers from `common` combined with the
// the handlers passed in
specific := common.With(
	SpecificMiddleware1,
	SpecificMiddleware2
)
```

## `Run()`
Negroni tem uma função de conveniência chamada `Run`. `Run` pega um endereço de string idêntico para [http.ListenAndServe](http://golang.org/pkg/net/http#ListenAndServe).

<!-- { "interrupt": true } -->
``` go
package main

import (
  "github.com/urfave/negroni"
)

func main() {
  n := negroni.Classic()
  n.Run(":8080")
}
```
Se nenhum endereço for fornecido, a variável de ambiente `PORT` é usada no lugar.
Se a variável de ambiente `PORT` não for definida, o endereço padrão será usado.
Veja [Run](https://godoc.org/github.com/urfave/negroni#Negroni.Run) para uma descrição completa.

No geral, você irá quere usar métodos `net/http` e passar `negroni` como um `Handler`,
como ele é mais flexível, e.g.:

<!-- { "interrupt": true } -->
``` go
package main

import (
  "fmt"
  "log"
  "net/http"
  "time"

  "github.com/urfave/negroni"
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

## Middleware para Rotas Específicas
Se você tem um grupo de rota com rotas que precisam ser executadas por um middleware específico, pode simplesmente criar uma nova instância de Negroni e usar no seu Manipulador de rota.

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

Se você está usando [Gorilla Mux], aqui é um exemplo usando uma subrota:

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

`With()` pode ser usado para eliminar redundância por middlewares compartilhados através de rotas.

``` go
router := mux.NewRouter()
apiRoutes := mux.NewRouter()
// add api routes here
webRoutes := mux.NewRouter()
// add web routes here

// create common middleware to be shared across routes
common := negroni.New(
	Middleware1,
	Middleware2,
)

// create a new negroni for the api middleware
// using the common middleware as a base
router.PathPrefix("/api").Handler(common.With(
  APIMiddleware1,
  negroni.Wrap(apiRoutes),
))
// create a new negroni for the web middleware
// using the common middleware as a base
router.PathPrefix("/web").Handler(common.With(
  WebMiddleware1,
  negroni.Wrap(webRoutes),
))
```

## Bundled Middleware

### Static

Este middleware servirá arquivos do filesystem. Se os arquivos não existerem,
ele substituirá a requisição para o próximo middleware. Se você quer que as
requisições para arquivos não existentes retorne um `404 File Not Found` para o user
você deve olhar em usando
[http.FileServer](https://golang.org/pkg/net/http/#FileServer) como um handler.

Exemplo:

<!-- { "interrupt": true } -->
``` go
package main

import (
  "fmt"
  "net/http"

  "github.com/urfave/negroni"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!")
  })

  // Example of using a http.FileServer if you want "server-like" rather than "middleware" behavior
  // mux.Handle("/public", http.FileServer(http.Dir("/home/public")))

  n := negroni.New()
  n.Use(negroni.NewStatic(http.Dir("/tmp")))
  n.UseHandler(mux)

  http.ListenAndServe(":3002", n)
}
```

Serviá arquivos do diretório `/tmp` primeiro, mas substituirá chamadas para o
próximo handler se a requisição não encontrar um arquivo no filesystem.

### Recovery

Este middleware pega repostas de `panic` com um código de resposta `500`. Se
algum outro middleware escreveu um código resposta ou conteúdo, este middleware
falhará corretamente enviando um 500 para o cliente, como o cliente já recebeu
o código de resposta HTTP. Adicionalmente um `PanicHandlerFunc` pode ser anexado
para reportar 500 para um serviço de report como um Sentry ou Airbrake.

Exemplo:

<!-- { "interrupt": true } -->
``` go
package main

import (
  "net/http"

  "github.com/urfave/negroni"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    panic("oh no")
  })

  n := negroni.New()
  n.Use(negroni.NewRecovery())
  n.UseHandler(mux)

  http.ListenAndServe(":3003", n)
}
```

Irá retornar um `500 Internal Server Error` para cada requisição. Ele também irá
registrar stack trace bem como um print na stack tarce para quem fez a requisição 
if `PrintStack` estiver setado como `true` (o padrão).

Exemplo com handler de erro:

``` go
package main

import (
  "net/http"

  "github.com/urfave/negroni"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    panic("oh no")
  })

  n := negroni.New()
  recovery := negroni.NewRecovery()
  recovery.PanicHandlerFunc = reportToSentry
  n.Use(recovery)
  n.UseHandler(mux)

  http.ListenAndServe(":3003", n)
}

func reportToSentry(info *negroni.PanicInformation) {
    // write code here to report error to Sentry
}
```


## Logger

Este middleware loga cada entrada de requisição e resposta.

Exemplo:

<!-- { "interrupt": true } -->
``` go
package main

import (
  "fmt"
  "net/http"

  "github.com/urfave/negroni"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!")
  })

  n := negroni.New()
  n.Use(negroni.NewLogger())
  n.UseHandler(mux)

  http.ListenAndServe(":3004", n)
}
```

Irá printar um log simular para:

```
[negroni] 2017-10-04T14:56:25+02:00 | 200 |      378µs | localhost:3004 | GET /
```

em cada request.

Você também pode definir seu próprio formato de log chamando a função `SetFormat`.
O formato é uma string de template como os campos como mencionados na struct `LoggerEntry`.
Então, como exemplo - 

```go
l.SetFormat("[{{.Status}} {{.Duration}}] - {{.Request.UserAgent}}")
```

Mostrará algo como  - `[200 18.263µs] - Go-User-Agent/1.1 `
## Middleware de Terceiros

Aqui está uma lista atual de Middleware Compatíveis com Negroni. Sinta se livre para mandar um PR vinculando seu middleware se construiu um:


| Middleware | Autor | Descrição |
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


## Exemplos
[Alexander Rødseth](https://github.com/xyproto) criou [mooseware](https://github.com/xyproto/mooseware), uma estrutura para escrever um handler middleware Negroni.

[Prasanga Siripala](https://github.com/pjebs) Criou um efetivo estrutura esqueleto para projetos Go/Negroni baseados na web: [Go-Skeleton](https://github.com/pjebs/go-skeleton) 

## Servidor com autoreload?
[gin](https://github.com/codegangsta/gin) e
[fresh](https://github.com/pilu/fresh) são aplicativos para autoreload do Negroni.

## Leitura Essencial para Iniciantes em Go & Negroni
* [Usando um contexto para passar informação de um middleware para o manipulador final](http://elithrar.github.io/article/map-string-interface/)
* [Entendendo middleware](http://mattstauffer.co/blog/laravel-5.0-middleware-replacing-filters)


## Sobre
Negroni é obsessivamente desenhado por ninguém menos que  [Code Gangsta](http://codegangsta.io/)

[Gorilla Mux]: https://github.com/gorilla/mux
[`http.FileSystem`]: https://godoc.org/net/http#FileSystem
