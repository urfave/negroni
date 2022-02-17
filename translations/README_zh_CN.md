# Negroni
[![GoDoc](https://godoc.org/github.com/urfave/negroni?status.svg)](http://godoc.org/github.com/urfave/negroni)
[![Build Status](https://travis-ci.org/urfave/negroni.svg?branch=master)](https://travis-ci.org/urfave/negroni)
[![codebeat](https://codebeat.co/badges/47d320b1-209e-45e8-bd99-9094bc5111e2)](https://codebeat.co/projects/github-com-urfave-negroni)
[![codecov](https://codecov.io/gh/urfave/negroni/branch/master/graph/badge.svg)](https://codecov.io/gh/urfave/negroni)

**注意:** 本函数库原本位于 `github.com/codegangsta/negroni` -- Github 会自动将请求重定向到当前地址, 但我们建议你更新一下引用路径。

在 Go 语言里，Negroni 是一个很地道的 Web 中间件，它轻量且非侵入，鼓励使用原生 `net/http` 处理器（Handlers）。

如果你喜欢用 [Martini](http://github.com/go-martini/martini) ，但又觉得它太魔幻，那么 Negroni 就是个很好的选择了。

各国语言翻译:
* [德语 (de_DE)](./README_de_de.md)
* [葡萄牙语 (pt_BR)](./README_pt_br.md)
* [简体中文 (zh_CN)](./README_zh_CN.md)
* [繁體中文 (zh_tw)](./README_zh_tw.md)
* [日本語 (ja_JP)](./README_ja_JP.md)
* [法语 (fr_FR)](./README_fr_FR.md)
* [韩语 (ko_KR)](./README_ko_KR.md)

## 入门指导

在安装了 Go 语言并设置好了 [GOPATH](http://golang.org/doc/code.html#GOPATH) 后，新建第一个 `.go` 文件，命名为 `server.go`。

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

然后安装 Negroni 包（**注意**：要求 **Go 1.1** 或更高的版本的 Go 语言环境）：

```
go get github.com/urfave/negroni
```

最后运行刚建好的 server.go 文件:

```
go run server.go
```

这时一个 Go `net/http` Web 服务器会跑在 `localhost:3000` 上，使用浏览器打开 `localhost:3000` 可看到输出的结果。

### 第三方包

如果你使用 Debian 系统，你可以执行 `apt install golang-github-urfave-negroni-dev` 来安装 `negroni`。 [包地址](https://packages.debian.org/sid/golang-github-urfave-negroni-dev) (写该文档时，它是在 `sid` 仓库中).


## Negroni 是一个框架吗？

Negroni **不**是一个框架，它是为了方便使用 `net/http` 而设计的一个库而已。

## 路由呢？

Negroni 没有带路由功能，因此使用 Negroni 时，需要有一个适合你的路由。不过好在 Go 社区里已经有相当多好用的路由，Negroni 和那些完全支持 `net/http` 库的路由搭配使用更佳，比如搭配 [Gorilla Mux](http://github.com/gorilla/mux) 路由，可以这样使用：

``` go
router := mux.NewRouter()
router.HandleFunc("/", HomeHandler)

n := negroni.New(Middleware1, Middleware2)
// Or use a middleware with the Use() function
n.Use(Middleware3)
// router goes last
n.UseHandler(router)

http.ListenAndServe(":3001", n)
```

## `negroni.Classic()` 经典的实例

`negroni.Classic()` 提供一些默认的中间件，这些中间件在多数应用都很有用。

* `negroni.Recovery` - 异常（恐慌）恢复中间件
* `negroni.Logging` - 请求 / 响应的日志中间件
* `negroni.Static` - 静态文件处理中间件，默认目录为 "public"

`negroni.Classic()` 使你可以轻松上手 Negroni 那些有用的特性。

## Handlers (处理器)

Negroni 提供双向的中间件机制，通过 `negroni.Handler` 这个接口来实现：

``` go
type Handler interface {
  ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
```

在中间件尚未写入 `ResponseWriter` 时，它会调用链上的下一个 `http.HandlerFunc` 以执行下一个中间件处理器。以下代码就是优雅的使用方式：

``` go
func MyMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // do some stuff before
  next(rw, r)
  // do some stuff after
}
```

你也可以用 `Use` 函数把它映射到处理器链上：

``` go
n := negroni.New()
n.Use(negroni.HandlerFunc(MyMiddleware))
```

你也可以映射 `http.Handler`。

``` go
n := negroni.New()

mux := http.NewServeMux()
// map your routes

n.UseHandler(mux)

http.ListenAndServe(":3000", n)
```

## `With()`

Negroni 有一个便利的函数叫 `With`。 `With` 函数可以把一个或多个 `Handler` 实例和接收者处理器集合组合成新的处理器集合，并返回新的 `Negroni` 实例对象。

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

Negroni 有一个便利的函数叫 `Run`。 `Run` 接收 addr 地址字符串 [http.ListenAndServe](http://golang.org/pkg/net/http#ListenAndServe).

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

未提供地址的情况下， `PORT` 系统环境变量会被使用。 若未定义该系统环境变量，默认的地址则会被使用。请参考 [Run](https://godoc.org/github.com/urfave/negroni#Negroni.Run) 的详情说明。

一般来说，使用 `net/http` 方法并将 Negroni 当作处理器传入，这样更灵活，例如:

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

## 特定路由（分组路由）

如果你有一组需要执行特定中间件的路由，你可以新建一个 Negroni 实例，然后把它当作你的路由处理器即可。

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

如果你使用 [Gorilla Mux](https://github.com/gorilla/mux)，下面是一个使用 subrouter 的例子:

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

`With()` 可被用来减少中间件在跨路由共享时的的冗余.

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

## 内置中间件

### 静态文件处理

该中间件通过文件系统来代理（处理）静态文件。 一旦文件不存在，请求代理会转到下个中间件。如果你想要处理不存在的文件返回 `404 File Not Found` HTTP 错误，请参考使用 [http.FileServer](https://golang.org/pkg/net/http/#FileServer) 作为处理器吧。

例子：

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

中间件首先在 `/tmp` 找文件，一旦在文件系统中找不到匹配的文件，代理将调用下一个处理器。


### 恢复

该中间件捕捉 `panic` 并返回错误代码为 `500` 的响应。如果其他中间件写了响应代码或 Body 内容的话，该中间件会无法顺利地传送 500 错误给客户端，因为客户端已经收到 HTTP 响应代码。另外，可以附上 `PanicHandlerFunc` 来报 500 错误给错误报告系统，如 Sentry 或 Airbrake。

例子：

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

它将输出 `500 Internal Server Error` 给每个请求。如果 `PrintStack` 设成 `true` （默认值）的话，它也会把错误日志写入请求方追踪堆栈。

加错误处理器的例子：

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
    // 在这写些程式回报错误给 Sentry
}
```

默认情况下，这个中间件会简要输出日志信息到 STDOUT 上。当然你也可以通过 `SetFormatter()` 函数自定义输出的日志。

当发生崩溃时，同样你也可以通过 `HTMLPanicFormatter` 来显示美化的 HTML 输出结果。

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
  recovery.Formatter = &negroni.HTMLPanicFormatter{}
  n.Use(recovery)
  n.UseHandler(mux)

  http.ListenAndServe(":3003", n)
}
```

## Logger

该中间件负责打印各个请求和响应日志。

例子：

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

每个请求打印日志将如下所示：

```
[negroni] 2017-10-04T14:56:25+02:00 | 200 |      378µs | localhost:3004 | GET /
```

你也可以调用 `SetFormat` 函数来自定义日志的格式。格式是一个预定义了字段的 `LoggerEntry` 结构体。如下所示：

```go
l.SetFormat("[{{.Status}} {{.Duration}}] - {{.Request.UserAgent}}")
```

会输出：`[200 18.263µs] - Go-User-Agent/1.1 `。

## 第三方兼容中间件

以下是兼容 Negroni 的中间件列表，如果你也有兼容 Negroni 的中间件，如果想提交自己的中间件，建议你附上 PR 链接。

| 中间件                                                                       | 作者                                                 | 描述                                                                                                                        |
| ---------------------------------------------------------------------------- | ---------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| [authz](https://github.com/casbin/negroni-authz)                             | [Yang Luo](https://github.com/hsluoyz)               | 支持ACL, RBAC, ABAC的权限管理中间件，基于[Casbin](https://github.com/casbin/casbin)                                         |
| [binding](https://github.com/mholt/binding)                                  | [Matt Holt](https://github.com/mholt)                | HTTP 请求数据注入到 structs 实体                                                                                            |
| [cloudwatch](https://github.com/cvillecsteele/negroni-cloudwatch)            | [Colin Steele](https://github.com/cvillecsteele)     | AWS CloudWatch 矩阵的中间件                                                                                                 |
| [cors](https://github.com/rs/cors)                                           | [Olivier Poitrey](https://github.com/rs)             | [Cross Origin Resource Sharing](http://www.w3.org/TR/cors/) (CORS) support                                                  |
| [csp](https://github.com/awakenetworks/csp)                                  | [Awake Networks](https://github.com/awakenetworks)   | 基于[Content Security Policy](https://www.w3.org/TR/CSP2/)(CSP)                                                             |
| [delay](https://github.com/jeffbmartinez/delay)                              | [Jeff Martinez](https://github.com/jeffbmartinez)    | 为endpoints增加延迟时间. 在测试严重网路延迟的效应时好用                                                                     |
| [New Relic Go Agent](https://github.com/yadvendar/negroni-newrelic-go-agent) | [Yadvendar Champawat](https://github.com/yadvendar)  | 官网 [New Relic Go Agent](https://github.com/newrelic/go-agent) (目前正在测试阶段)                                          |
| [gorelic](https://github.com/jingweno/negroni-gorelic)                       | [Jingwen Owen Ou](https://github.com/jingweno)       | New Relic agent for Go runtime                                                                                              |
| [Graceful](https://github.com/stretchr/graceful)                             | [Tyler Bunnell](https://github.com/tylerb)           | 优雅关闭 HTTP 的中间件                                                                                                      |
| [gzip](https://github.com/phyber/negroni-gzip)                               | [phyber](https://github.com/phyber)                  | 响应流 GZIP 压缩                                                                                                            |
| [JWT Middleware](https://github.com/auth0/go-jwt-middleware)                 | [Auth0](https://github.com/auth0)                    | Middleware checks for a JWT on the `Authorization` header on incoming requests and decodes it                               |
| [logrus](https://github.com/meatballhat/negroni-logrus)                      | [Dan Buch](https://github.com/meatballhat)           | 基于 Logrus-based logger 日志                                                                                               |
| [oauth2](https://github.com/goincremental/negroni-oauth2)                    | [David Bochenski](https://github.com/bochenski)      | oAuth2 中间件                                                                                                               |
| [onthefly](https://github.com/xyproto/onthefly)                              | [Alexander Rødseth](https://github.com/xyproto)      | 快速生成 TinySVG， HTML and CSS 中间件                                                                                      |
| [permissions2](https://github.com/xyproto/permissions2)                      | [Alexander Rødseth](https://github.com/xyproto)      | Cookies， 用户和权限                                                                                                        |
| [prometheus](https://github.com/zbindenren/negroni-prometheus)               | [Rene Zbinden](https://github.com/zbindenren)        | 简易建立矩阵端点给[prometheus](http://prometheus.io)建构工具                                                                |
| [render](https://github.com/unrolled/render)                                 | [Cory Jacobsen](https://github.com/unrolled)         | 渲染 JSON, XML and HTML 中间件                                                                                              |
| [RestGate](https://github.com/pjebs/restgate)                                | [Prasanga Siripala](https://github.com/pjebs)        | REST API 接口的安全认证                                                                                                     |
| [secure](https://github.com/unrolled/secure)                                 | [Cory Jacobsen](https://github.com/unrolled)         | Middleware that implements a few quick security wins                                                                        |
| [sessions](https://github.com/goincremental/negroni-sessions)                | [David Bochenski](https://github.com/bochenski)      | Session 会话管理                                                                                                            |
| [stats](https://github.com/thoas/stats)                                      | [Florent Messa](https://github.com/thoas)            | 检测 web 应用当前运行状态信息 （响应时间等等。）                                                                            |
| [VanGoH](https://github.com/auroratechnologies/vangoh)                       | [Taylor Wrobel](https://github.com/twrobel3)         | Configurable [AWS-Style](http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html) 基于 HMAC 鉴权认证的中间件 |
| [xrequestid](https://github.com/pilu/xrequestid)                             | [Andrea Franz](https://github.com/pilu)              | 给每个请求指定一个随机 X-Request-Id 头的中间件                                                                              |
| [mgo session](https://github.com/joeljames/nigroni-mgo-session)              | [Joel James](https://github.com/joeljames)           | 处理在每个请求建立与关闭 mgo sessions                                                                                       |
| [digits](https://github.com/bamarni/digits)                                  | [Bilal Amarni](https://github.com/bamarni)           | 处理 [Twitter Digits](https://get.digits.com/) 的认证                                                                       |
| [stats](https://github.com/guptachirag/stats)                                | [Chirag Gupta](https://github.com/guptachirag/stats) | endpoints用的管理QPS与延迟状态的中间件非同步地将状态刷入InfluxDB                                                            |
| [Chaos](https://github.com/falzm/chaos)                                      | [Marc Falzon](https://github.com/falzm)              | 以编程方式在应用程式中插入无序行为的中间件                                                                                  |

## 范例

[Alexander Rødseth](https://github.com/xyproto) 创建的 [mooseware](https://github.com/xyproto/mooseware) 是一个编写兼容 Negroni 中间件处理器的骨架。

[Prasanga Siripala](https://github.com/pjebs) 创建的 [Go-Skeleton](https://github.com/pjebs/go-skeleton) 是一个高效编写基于 web 的 Go/Negroni 项目的骨架。

## 即时编译

[gin](https://github.com/codegangsta/gin) 和 [fresh](https://github.com/pilu/fresh) 这两个应用是即时编译的 Negroni 工具。

## Go & Negroni 初学者必读推荐

* [使用上下文把消息从中间件传递给后端处理器](http://elithrar.github.io/article/map-string-interface/)
* [理解中间件](http://mattstauffer.co/blog/laravel-5.0-middleware-replacing-filters)

## 关于 Negroni

 Negroni 原由 [Code Gangsta](https://codegangsta.io/) 主导设计开发。
