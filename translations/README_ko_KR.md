# Negroni

[![GoDoc](https://godoc.org/github.com/urfave/negroni?status.svg)](http://godoc.org/github.com/urfave/negroni)
[![Build Status](https://travis-ci.org/urfave/negroni.svg?branch=master)](https://travis-ci.org/urfave/negroni)
[![codebeat](https://codebeat.co/badges/47d320b1-209e-45e8-bd99-9094bc5111e2)](https://codebeat.co/projects/github-com-urfave-negroni)
[![codecov](https://codecov.io/gh/urfave/negroni/branch/master/graph/badge.svg)](https://codecov.io/gh/urfave/negroni)

**공지:** 이 라이브러리는 아래의 주소로 많이 알려져 왔습니다.
`github.com/codegangsta/negroni` -- Github가 자동으로 이 저장소에 대한 요청을 리다이렉트 시킬 것이지만, 확실한 사용을 위해 참조를 이곳으로 변경하는 것을 추천드립니다.

Negroni는 Go에서 웹 미들웨어로의 자연스러운 접근을 추구합니다. 이것은 작고, 거슬리지 않으며 `net/http` 핸들러의 사용을 지향하는 것을 의미합니다.

만약 당신이 [Martini](https://github.com/go-martini/martini)의 기본적인 컨셉을 원하지만, 이것이 너무 많은 기능을 포함한다고 느껴졌다면 Negroni가 최적의 선택일 것입니다.

## 시작하기

Go 설치와 [GOPATH](http://golang.org/doc/code.html#GOPATH)를 세팅하는 작업을 완료한 뒤, 당신의 첫 `.go` 파일을 생성하세요.
우리는 이를  `server.go` 라고 부를 것입니다.

<!-- { "interrupt": true } -->

```go
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

  n := negroni.Classic() // 기본 미들웨어들을 포함합니다
  n.UseHandler(mux)

  http.ListenAndServe(":3000", n)
}
```

그리고 Negroni 패키지를 설치합니다  (**공지**:  **go 1.1** 이상이 요구됩니다) :

```
go get github.com/urfave/negroni
```

서버를 실행시킵니다:

```
go run server.go
```

이제 `localhost:3000`에서 동작하는 Go `net/http` 웹서버를 가지게 되었습니다.

### 패키징 (*Packaging*)

Debian을 사용중이시라면, `negroni`는 [a package](https://packages.debian.org/sid/golang-github-urfave-negroni-dev)에서도 사용이 가능합니다. `apt install golang-github-urfave-negroni-dev`를 통해서 설치가 가능합니다. (글을 작성할 당시, 이는 `sid` 저장소 안에 있습니다.)

## Negroni는 프레임워크(*Framework*)인가요?

Negroni는 프레임워크가 **아닙니다.** 이는 `net/http`를 직접적으로 이용할 수 있도록 디자인된 미들웨어 중심의 라이브러리입니다.

## 라우팅(*Routing*) ?

Negroni는 *BYOR* (*Bring your own Router*, 나의 라우터를 사용하기)를 지향합니다. Go 커뮤니티에는 이미 좋은 http 라우터들이 존재하기 때문에 Negroni는 그들과 잘 어우러질 수 있도록 `net/http`를 전적으로 지원하고 있습니다. 

예를 들어, 아래는 [Gorilla Mux]를 사용한 예입니다:

```go
router := mux.NewRouter()
router.HandleFunc("/", HomeHandler)

n := negroni.New(Middleware1, Middleware2)
// Use() 함수를 사용해서 미들웨어를 사용할수도 있습니다.
n.Use(Middleware3)
// 라우터의 경우 마지막에 옵니다.
n.UseHandler(router)

http.ListenAndServe(":3001", n)
```

## `negroni.Classic()`

`negroni.Classic()`은 대부분의 어플리케이션에서 유용하게 사용되는 기본적인 미들웨어들을 제공합니다.

- [`negroni.Recovery`](#recovery) - Panic 복구용 미들웨어
- [`negroni.Logger`](#logger) - Request/Response 로깅 미들웨어
- [`negroni.Static`](#static) - "public" 디렉터리 아래의 정적 파일 제공(*serving*)을 위한 미들웨어

이는 Negroni의 유용한 기능들을 사용하기 시작하는데 큰 도움이 되도록 만들어줄 것입니다.

## 핸들러(*Handlers*)

Negroni는 양방향 미들웨어 흐름을 제공합니다. 이는 `negroni.Handler` 인터페이스를 통해 구현합니다.

```go
type Handler interface {
  ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
```

미들웨어가 `ResponseWriter`에 아직 무언가 쓰지 않았다면, 이는 다음 미들웨어에 연결되어있는  `http.HandleFunc`를 호출해야 합니다. 이는 유용하게 사용될 수 있습니다:

```go
func MyMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // next() 처리 이전 작업 수행
  next(rw, r)
  // next() 처리 이후 작업 수행
}
```

이후 `Use` 함수를 통해 핸들러 체인(*handler chain*)에 매핑 시킬 수 있습니다:

```go
n := negroni.New()
n.Use(negroni.HandlerFunc(MyMiddleware))
```

또한, 기존의 `http.Handler`들과도 매핑시킬 수 있습니다:

```go
n := negroni.New()

mux := http.NewServeMux()
// 여기에 라우트들을 매핑하세요

n.UseHandler(mux)

http.ListenAndServe(":3000", n)
```

## `With()`

Negroni는 `With`라고 불리는 편리한 함수를 가지고 있습니다. `With`는 한 개 혹은 그 이상의 `Handler` 인스턴스들을 받아 기존 리시버의 핸들러들과 새로운 핸들러들이 조합된 새로운 `Negroni` 객체를 리턴합니다.

```go
// 재사용을 원하는 미들웨어들
common := negroni.New()
common.Use(MyMiddleware1)
common.Use(MyMiddleware2)

// `specific`은 `common`의 핸들러들과 새로 전달된 핸들러들이 조합된 새로운 `negroni` 객체
specific := common.With(
	SpecificMiddleware1,
	SpecificMiddleware2
)
```

## `Run()`

Negroni는 `Run`이라고 불리는 편리한 함수를 가지고 있습니다. `Run`은 `http.ListenAndServe`와 같이 주소 스트링 값(*addr string*)을 넘겨받습니다.

<!-- { "interrupt": true } -->

```go
package main

import (
  "github.com/urfave/negroni"
)

func main() {
  n := negroni.Classic()
  n.Run(":8080")
}
```

만약 주소 값이 제공되지 않는다면, `PORT` 환경 변수가 대신 사용됩니다. `PORT` 환경 변수 또한 정의되어있지 않다면, 기본 주소(*default address*)가 사용됩니다. 전체 설명을 보시려면 [Run](https://godoc.org/github.com/urfave/negroni#Negroni.Run)을 참고하세요.

일반적으로는, 좀 더 유연한 사용을 위해서 `net/http` 메서드를 사용하여 `negroni` 객체를 핸들러로서 넘기는 것을 선호할 것입니다. 예를 들면:

<!-- { "interrupt": true } -->

```go
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

  n := negroni.Classic() // 기본 미들웨어들을 포함합니다
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

## 라우트 전용 미들웨어(*Route Specific Middleware*)

만약 당신이 라우트들의 라우트 그룹을 가지고 있다면 그런데 그것은 실행되어야하는 미들웨어를 필요로한다.

특정 라우트 그룹만이 사용하는 미들웨어가 있다면, 간단하게 Negroni 인스턴스를 새롭게 생성하여 라우트 핸들러(*route handler*)로서 사용하면 된다.

```go
router := mux.NewRouter()
adminRoutes := mux.NewRouter()
// admin routes를 여기에 추가하세요

// 관리자 미들웨어들을 위한 새로운 negroni 인스턴스를 생성합니다
router.PathPrefix("/admin").Handler(negroni.New(
  Middleware1,
  Middleware2,
  negroni.Wrap(adminRoutes),
))
```

[Gorilla Mux]를 사용하고 있다면, 아래는 서브 라우터(subrouter)를 사용하는 예제입니다:

```go
router := mux.NewRouter()
subRouter := mux.NewRouter().PathPrefix("/subpath").Subrouter().StrictSlash(true)
subRouter.HandleFunc("/", someSubpathHandler) // "/subpath/"
subRouter.HandleFunc("/:id", someSubpathHandler) // "/subpath/:id"

// "/subpath" 는 subRouter와 메인 라우터(main router)의 연결 보장을 위해 반드시 필요합니다
router.PathPrefix("/subpath").Handler(negroni.New(
  Middleware1,
  Middleware2,
  negroni.Wrap(subRouter),
))
```

`With()`는 라우터 간 공유시 발생하는 미들웨어 중복을 방지하기 위해서 사용될 수 있습니다.

```go
router := mux.NewRouter()
apiRoutes := mux.NewRouter()
// api 라우트들을 여기에 추가하세요
webRoutes := mux.NewRouter()
// web 라우트들을 여기에 추가하세요

// 라우터간 공유될 common 미들웨어를 생성합니다
common := negroni.New(
	Middleware1,
	Middleware2,
)

// common 미들웨어를 기반으로 
// api 미들웨어를 위한 새로운 negroni 객체를 생성합니다
router.PathPrefix("/api").Handler(common.With(
  APIMiddleware1,
  negroni.Wrap(apiRoutes),
))
// common 미들웨어를 기반으로 
// web 미들웨어를 위한 새로운 negroni 객체를 생성합니다
router.PathPrefix("/web").Handler(common.With(
  WebMiddleware1,
  negroni.Wrap(webRoutes),
))
```

## 번들 미들웨어(*Bundled Middleware*)

### Static

이 미들웨어는 파일들을 파일 시스템(*filesystem*)으로 제공하는 역할을 수행합니다. 파일이 존재하지 않을 경우, 요청을 다음 미들웨어로 넘깁니다. 존재하지 않는 파일에 대해 `404 File Not Found`를 유저에게 반환하길 원하는 경우 [http.FileServer](https://golang.org/pkg/net/http/#FileServer)를 핸들러로 사용하는 것을 살펴보아야 합니다.

예제:

<!-- { "interrupt": true } -->

```go
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

  // "미들웨어(middleware)"의 역할보다는 "서버와 같은(server-like)" 역할을 수행하기를 원할 때 
  // http.FileServer를 사용한 예제
  // mux.Handle("/public", http.FileServer(http.Dir("/home/public")))  

  n := negroni.New()
  n.Use(negroni.NewStatic(http.Dir("/tmp")))
  n.UseHandler(mux)

  http.ListenAndServe(":3002", n)
}
```

위 코드는 `/tmp` 디렉터리로부터 파일을 제공할 것입니다. 그러나 파일 시스템 상의 파일에 대한 요청이 일치하지 않는 경우 프록시는 다음 핸들러를 호출할 것입니다.

### Recovery

이 미들웨어는 `panic`들을 감지하고 `500` 응답 코드(*response code*)를 반환하는 역할을 수행합니다. 다른 미들웨어가 응답 코드 또는 바디(*body*)를 쓸 경우, 클라이언트는 이미 HTTP 응답 코드를 받았기 때문에 이 미들웨어가 적절한 시점에 `500` 코드를 보내는 것에 실패할 것입니다. 추가적으로 `PanicHandlerFunc`는 Sentry 또는 Aribrake와 같은 에러 보고 서비스에 `500` 코드를 반환하도록 붙을 수 있습니다.

예제:

<!-- { "interrupt": true } -->

```go
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

위 코드는 각 요청에 대해 `500 Internal Server Error` 반환할 것입니다. `PrintStack` 값이 `true` (기본 값)로 설정되어있다면 요청자(*requester*)에게 스택 트레이스(*stack trace*) 값을 출력하는 것처럼 로깅 또한 진행합니다.

`PanicHandlerFunc`를 사용한 예제:

```go
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
    // Sentry에게 에러를 보고하는 코드를 작성하세요
}
```

미들웨어는 `STDOUT`에 기본으로 값들을 출력합니다. 하지만 `SetFormatter()` 함수를 이용해서 출력 프로세스를 커스터마이징 할 수 있습니다. 당연히 `HTMLPanicFormatter`를 사용해서 깔끔한 HTML로도 에러 상황을 보여줄 수 있습니다.

<!-- { "interrupt": true } -->

```go
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
  recovery.Formatter = &negroni.HTMLPanicFormatter{}
  n.Use(recovery)
  n.UseHandler(mux)

  http.ListenAndServe(":3003", n)
}
```

## Logger

이 미들웨어는 서버에 들어오는 요청과 응답들을 기록하는 역할을 수행합니다.

예제:

<!-- { "interrupt": true } -->

```go
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

위 코드는 각 요청에 대해 아래와 같이 출력할 것입니다.

```
[negroni] 2017-10-04T14:56:25+02:00 | 200 |      378µs | localhost:3004 | GET /
```

물론, `SetFormat` 함수를 이용해 사용자만의 로그 포맷(*log format*) 또한 정의할 수 있습니다. 로그 포맷은 `LoggerEntry` 구조체 내부의 필드들로 구성된 템플릿 문자열입니다.

```go
l.SetFormat("[{{.Status}} {{.Duration}}] - {{.Request.UserAgent}}")
```

위 구조는 이와 같이 출력될 것입니다 - `[200 18.263µs] - Go-User-Agent/1.1 `

## Third Party Middleware

아래는 현재(2018.10.10) Negroni와 호환되는 미들웨어들입니다. 당신의 미들웨어가 링크되기를 원한다면 자유롭게 PR을 보내주세요.

| Middleware                                                   | Author                                               | Description                                                  |
| ------------------------------------------------------------ | ---------------------------------------------------- | ------------------------------------------------------------ |
| [authz](https://github.com/casbin/negroni-authz)             | [Yang Luo](https://github.com/hsluoyz)               | ACL, RBAC, ABAC Authorization middlware based on [Casbin](https://github.com/casbin/casbin) |
| [binding](https://github.com/mholt/binding)                  | [Matt Holt](https://github.com/mholt)                | Data binding from HTTP requests into structs                 |
| [cloudwatch](https://github.com/cvillecsteele/negroni-cloudwatch) | [Colin Steele](https://github.com/cvillecsteele)     | AWS cloudwatch metrics middleware                            |
| [cors](https://github.com/rs/cors)                           | [Olivier Poitrey](https://github.com/rs)             | [Cross Origin Resource Sharing](http://www.w3.org/TR/cors/) (CORS) support |
| [csp](https://github.com/awakenetworks/csp)                  | [Awake Networks](https://github.com/awakenetworks)   | [Content Security Policy](https://www.w3.org/TR/CSP2/) (CSP) support |
| [delay](https://github.com/jeffbmartinez/delay)              | [Jeff Martinez](https://github.com/jeffbmartinez)    | Add delays/latency to endpoints. Useful when testing effects of high latency |
| [New Relic Go Agent](https://github.com/yadvendar/negroni-newrelic-go-agent) | [Yadvendar Champawat](https://github.com/yadvendar)  | Official [New Relic Go Agent](https://github.com/newrelic/go-agent) (currently in beta) |
| [gorelic](https://github.com/jingweno/negroni-gorelic)       | [Jingwen Owen Ou](https://github.com/jingweno)       | New Relic agent for Go runtime                               |
| [Graceful](https://github.com/tylerb/graceful)               | [Tyler Bunnell](https://github.com/tylerb)           | Graceful HTTP Shutdown                                       |
| [gzip](https://github.com/phyber/negroni-gzip)               | [phyber](https://github.com/phyber)                  | GZIP response compression                                    |
| [JWT Middleware](https://github.com/auth0/go-jwt-middleware) | [Auth0](https://github.com/auth0)                    | Middleware checks for a JWT on the `Authorization` header on incoming requests and decodes it |
| [JWT Middleware](https://github.com/mfuentesg/go-jwtmiddleware) | [Marcelo Fuentes](https://github.com/mfuentesg)      | JWT middleware for golang                                    |
| [logrus](https://github.com/meatballhat/negroni-logrus)      | [Dan Buch](https://github.com/meatballhat)           | Logrus-based logger                                          |
| [oauth2](https://github.com/goincremental/negroni-oauth2)    | [David Bochenski](https://github.com/bochenski)      | oAuth2 middleware                                            |
| [onthefly](https://github.com/xyproto/onthefly)              | [Alexander Rødseth](https://github.com/xyproto)      | Generate TinySVG, HTML and CSS on the fly                    |
| [permissions2](https://github.com/xyproto/permissions2)      | [Alexander Rødseth](https://github.com/xyproto)      | Cookies, users and permissions                               |
| [prometheus](https://github.com/zbindenren/negroni-prometheus) | [Rene Zbinden](https://github.com/zbindenren)        | Easily create metrics endpoint for the [prometheus](http://prometheus.io) instrumentation tool |
| [prometheus](https://github.com/slok/go-prometheus-middleware) | [Xabier Larrakoetxea](https://github.com/slok)       | [Prometheus](http://prometheus.io) metrics with multiple options that follow standards and try to be measured in a efficent way |
| [render](https://github.com/unrolled/render)                 | [Cory Jacobsen](https://github.com/unrolled)         | Render JSON, XML and HTML templates                          |
| [RestGate](https://github.com/pjebs/restgate)                | [Prasanga Siripala](https://github.com/pjebs)        | Secure authentication for REST API endpoints                 |
| [secure](https://github.com/unrolled/secure)                 | [Cory Jacobsen](https://github.com/unrolled)         | Middleware that implements a few quick security wins         |
| [sessions](https://github.com/goincremental/negroni-sessions) | [David Bochenski](https://github.com/bochenski)      | Session Management                                           |
| [stats](https://github.com/thoas/stats)                      | [Florent Messa](https://github.com/thoas)            | Store information about your web application (response time, etc.) |
| [VanGoH](https://github.com/auroratechnologies/vangoh)       | [Taylor Wrobel](https://github.com/twrobel3)         | Configurable [AWS-Style](http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html) HMAC authentication middleware |
| [xrequestid](https://github.com/pilu/xrequestid)             | [Andrea Franz](https://github.com/pilu)              | Middleware that assigns a random X-Request-Id header to each request |
| [mgo session](https://github.com/joeljames/nigroni-mgo-session) | [Joel James](https://github.com/joeljames)           | Middleware that handles creating and closing mgo sessions per request |
| [digits](https://github.com/bamarni/digits)                  | [Bilal Amarni](https://github.com/bamarni)           | Middleware that handles [Twitter Digits](https://get.digits.com/) authentication |
| [stats](https://github.com/guptachirag/stats)                | [Chirag Gupta](https://github.com/guptachirag/stats) | Middleware that manages qps and latency stats for your endpoints and asynchronously flushes them to influx db |
| [Chaos](https://github.com/falzm/chaos)                      | [Marc Falzon](https://github.com/falzm)              | Middleware for injecting chaotic behavior into application in a programmatic way |

## 예제

[Alexander Rødseth](https://github.com/xyproto)는 Negroni 미들웨어 핸들러를 작성하기 위한 뼈대인 [mooseware](https://github.com/xyproto/mooseware)를 만들었습니다.

[Prasanga Siripala](https://github.com/pjebs)는 웹 기반의 Go/Negroni 프로젝트들을 위한 효율적인 뼈대 구조를 만들었습니다 : [Go-Skeleton](https://github.com/pjebs/go-skeleton) 

## 코드 실시간 새로고침(Live code reload)?

[gin](https://github.com/codegangsta/gin)과 [fresh](https://github.com/pilu/fresh) 모두 negroni 앱의 실시간 새로고침(*live reload*)을 지원합니다.

## Go & Negroni 초심자들이 필수적으로 읽어야하는 자료들

- [Using a Context to pass information from middleware to end handler](http://elithrar.github.io/article/map-string-interface/)
- [Understanding middleware](https://mattstauffer.co/blog/laravel-5.0-middleware-filter-style)

## 추가 정보

Negroni는 [Code Gangsta](https://codegangsta.io/)에 의해 디자인 되었습니다.

[Gorilla Mux]: https://github.com/gorilla/mux
[`http.FileSystem`]: https://godoc.org/net/http#FileSystem