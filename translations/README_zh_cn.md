# Negroni
[![GoDoc](https://godoc.org/github.com/urfave/negroni?status.svg)](http://godoc.org/github.com/urfave/negroni)
[![Build Status](https://travis-ci.org/urfave/negroni.svg?branch=master)](https://travis-ci.org/urfave/negroni)
[![codebeat](https://codebeat.co/badges/47d320b1-209e-45e8-bd99-9094bc5111e2)](https://codebeat.co/projects/github-com-urfave-negroni)
[![codecov](https://codecov.io/gh/urfave/negroni/branch/master/graph/badge.svg)](https://codecov.io/gh/urfave/negroni)

**注意:** 本函式库原来自于
`github.com/codegangsta/negroni` -- Github会自动将连线转到本连结, 但我们建议你更新一下参照.

在Go语言里，Negroni 是一个很地道的 web 中间件，它是微型，非嵌入式，并鼓励使用原生 `net/http` 处理器的库。

如果你用过并喜欢 [Martini](http://github.com/go-martini/martini) 框架，但又不想框架中有太多魔幻性的特征，那 Negroni 就是你的菜了，相信它非常适合你。

语言翻译:
* [German (de_DE)](translations/README_de_de.md)
* [Português Brasileiro (pt_BR)](translations/README_pt_br.md)
* [简体中文 (zh_cn)](translations/README_zh_cn.md)
* [繁體中文 (zh_tw)](translations/README_zh_tw.md)
* [日本語 (ja_JP)](translations/README_ja_JP.md)

## 入门指导

当安装了 Go 语言并设置好了 [GOPATH](http://golang.org/doc/code.html#GOPATH) 后，新建你第一个`.go` 文件，我们叫它 `server.go` 吧。

``` go
package main

import (
  "github.com/urfave/negroni"
  "net/http"
  "fmt"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!")
  })

  n := negroni.Classic()
  n.UseHandler(mux)
  n.Run(":3000")
}
```

然后安装 Negroni 包（它依赖 **Go 1.1** 或更高的版本）：
```
go get github.com/urfave/negroni
```

然后运行刚建好的 server.go 文件:
```
go run server.go
```

这时一个 Go `net/http` Web 服务器就跑在 `localhost:3000` 上，使用浏览器打开 `localhost:3000` 可以看到输出结果。

### 打包
如果`negroni`在Debian环境下是个[套件](https://packages.debian.org/sid/golang-github-urfave-negroni-dev), 可直接
执行`apt install golang-github-urfave-negroni-dev`安装(这在`sid`仓库中).


## Negroni 是一个框架吗？
Negroni **不**是一个框架，它是为了方便使用 `net/http` 而设计的一个库而已。

## 路由呢？
Negroni 没有带路由功能，使用 Negroni 时，需要找一个适合你的路由。不过好在 Go 社区里已经有相当多可用的路由，Negroni 更喜欢和那些完全支持 `net/http` 库的路由组合使用，比如，结合 [Gorilla Mux](http://github.com/gorilla/mux) 使用像这样：

``` go
router := mux.NewRouter()
router.HandleFunc("/", HomeHandler)

n := negroni.New(Middleware1, Middleware2)
// Or use a middleware with the Use() function
n.Use(Middleware3)
// router goes last
n.UseHandler(router)

n.Run(":3000")
```

## `negroni.Classic()` 经典实例
`negroni.Classic()` 提供一些默认的中间件，这些中间件在多数应用都很有用。

* `negroni.Recovery` - 异常（恐慌）恢复中间件
* `negroni.Logging` - 请求 / 响应 log 日志中间件
* `negroni.Static` - 静态文件处理中间件，默认目录在 "public" 下.

`negroni.Classic()` 让你一开始就非常容易上手 Negroni ，并使用它那些通用的功能。

## Handlers (处理器)
Negroni 提供双向的中间件机制，这个特征很棒，都是得益于 `negroni.Handler` 这个接口。

``` go
type Handler interface {
  ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
```

如果一个中间件没有写入 ResponseWriter 响应，它会在中间件链里调用下一个 `http.HandlerFunc` 执行下去， 它可以这么优雅的使用。如下：

``` go
func MyMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // do some stuff before
  next(rw, r)
  // do some stuff after
}
```

你也可以用 `Use` 函数把这些 `http.Handler` 处理器引进到处理器链上来：

``` go
n := negroni.New()
n.Use(negroni.HandlerFunc(MyMiddleware))
```

你还可以使用 `http.Handler`(s) 把 `http.Handler` 处理器引进来。

``` go
n := negroni.New()

mux := http.NewServeMux()
// map your routes

n.UseHandler(mux)

n.Run(":3000")
```

## `Run()`
尼格龙尼有一个很好用的函数`Run`, `Run`接收addr字串辨识[http.ListenAndServe](http://golang.org/pkg/net/http#ListenAndServe).

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

未提供路径情况下会使用系统环境变数`PORT`, 若未定义该系统环境变数则会用预设路径, 请见[Run](https://godoc.org/github.com/urfave/negroni#Negroni.Run)细看说明.

一般来说, 你会希望使用 `net/http` 方法, 并且将尼格龙尼当作处理器传入, 这相对起来弹性比较大, 例如:

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

  n := negroni.Classic() // 導入一些預設中介器
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

## 特定路由中间件
如果你需要群组路由功能，需要借助特定的路由中间件完成，做法很简单，只需建立一个新 Negroni 实例，传人路由处理器里即可。

``` go
router := mux.NewRouter()
adminRoutes := mux.NewRouter()
// add admin routes here

// Create a new negroni for the admin middleware
router.Handle("/admin", negroni.New(
  Middleware1,
  Middleware2,
  negroni.Wrap(adminRoutes),
))
```

如果你使用 [Gorilla Mux](https://github.com/gorilla/mux), 下方是一个使用 subrounter 的例子:

``` go
router := mux.NewRouter()
subRouter := mux.NewRouter().PathPrefix("/subpath").Subrouter().StrictSlash(true)
subRouter.HandleFunc("/", someSubpathHandler) // "/subpath/"
subRouter.HandleFunc("/:id", someSubpathHandler) // "/subpath/:id"

// "/subpath" 是用来保证subRouter与主要路由连结的必要参数
router.PathPrefix("/subpath").Handler(negroni.New(
  Middleware1,
  Middleware2,
  negroni.Wrap(subRouter),
))
```

`With()` 可被用来降低在跨路由分享时多余的中介器.

``` go
router := mux.NewRouter()
apiRoutes := mux.NewRouter()
// 在此新增API路由
webRoutes := mux.NewRouter()
// 在此新增Web路由

// 建立通用中介器来跨路由分享
common := negroni.New(
  Middleware1,
  Middleware2,
)

// 为API中介器建立新的negroni
// 使用通用中介器作底
router.PathPrefix("/api").Handler(common.With(
  APIMiddleware1,
  negroni.Wrap(apiRoutes),
))
// 为Web中介器建立新的negroni
// 使用通用中介器作底
router.PathPrefix("/web").Handler(common.With(
  WebMiddleware1,
  negroni.Wrap(webRoutes),
))
```

## 内建中介器

### 静态

本中介器会在档案系统上服务档案. 若档案不存在, 会将流量导(proxy)到下个中介器.
如果你想要返回`404 File Not Found`给档案不存在的请求, 请使用[http.FileServer](https://golang.org/pkg/net/http/#FileServer)
作为处理器.

范例:

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

  // http.FileServer的使用范例, 若你预期要"像伺服器"而非"中介器"的行为
  // mux.Handle("/public", http.FileServer(http.Dir("/home/public")))

  n := negroni.New()
  n.Use(negroni.NewStatic(http.Dir("/tmp")))
  n.UseHandler(mux)

  http.ListenAndServe(":3002", n)
}
```

从`/tmp`目录开始服务档案 但如果请求的档案在档案系统中不符合, 代理会
呼叫下个处理器.

### 恢复

本中介器接收`panic`跟错误代码`500`的回应. 如果其他任何中介器写了回应
的HTTP代码或内容的话, 中介器会无法顺利地传送500给用户端, 因为用户端
已经收到HTTP的回应代码. 另外, 可以挂载`ErrorHandlerFunc`来回报500
的错误到错误回报系统, 如: Sentry或Airbrake.

范例:

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


将回传`500 Internal Server Error`到每个结果. 也会把结果纪录到堆叠追踪,
`PrintStack`设成`true`(预设值)的话也会印到注册者.

加错误处理器的范例:

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
  recovery.ErrorHandlerFunc = reportToSentry
  n.Use(recovery)
  n.UseHandler(mux)

  http.ListenAndServe(":3003", n)
}

func reportToSentry(error interface{}) {
    // 在这写些程式回报错误给Sentry
}
```


## Logger

本中介器纪录各个进入的请求与回应.

范例:

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

在每个请求印的纪录会看起来像:

```
[negroni] Started GET /
[negroni] Completed 200 OK in 145.446µs
```

## 第三方中间件

以下的兼容 Negroni 的中间件列表，如果你也有兼容 Negroni 的中间件，可以提交到这个列表来交换链接，我们很乐意做这样有益的事情。


|    中间件    |    作者    |    描述     |
| -------------|------------|-------------|
| [authz](https://github.com/casbin/negroni) | [Yang Luo](https://github.com/hsluoyz) | 支持ACL, RBAC, ABAC的权限管理中间件，基于[Casbin](https://github.com/casbin/casbin) |
| [binding](https://github.com/mholt/binding) | [Matt Holt](https://github.com/mholt) | HTTP 请求数据注入到 structs 实体|
| [cloudwatch](https://github.com/cvillecsteele/negroni-cloudwatch) | [Colin Steele](https://github.com/cvillecsteele) | AWS CloudWatch 矩阵的中间件 |
| [cors](https://github.com/rs/cors) | [Olivier Poitrey](https://github.com/rs) | [Cross Origin Resource Sharing](http://www.w3.org/TR/cors/) (CORS) support |
| [csp](https://github.com/awakenetworks/csp) | [Awake Networks](https://github.com/awakenetworks) | 基于[Content Security Policy](https://www.w3.org/TR/CSP2/)(CSP) |
| [delay](https://github.com/jeffbmartinez/delay) | [Jeff Martinez](https://github.com/jeffbmartinez) | 为endpoints增加延迟时间. 在测试严重网路延迟的效应时好用 |
| [New Relic Go Agent](https://github.com/yadvendar/negroni-newrelic-go-agent) | [Yadvendar Champawat](https://github.com/yadvendar) | 官网 [New Relic Go Agent](https://github.com/newrelic/go-agent) (目前正在测试阶段)  |
| [gorelic](https://github.com/jingweno/negroni-gorelic) | [Jingwen Owen Ou](https://github.com/jingweno) | New Relic agent for Go runtime |
| [Graceful](https://github.com/stretchr/graceful) | [Tyler Bunnell](https://github.com/tylerb) | 优雅关闭 HTTP 的中间件 |
| [gzip](https://github.com/phyber/negroni-gzip) | [phyber](https://github.com/phyber) | 响应流 GZIP 压缩 |
| [JWT Middleware](https://github.com/auth0/go-jwt-middleware) | [Auth0](https://github.com/auth0) | Middleware checks for a JWT on the `Authorization` header on incoming requests and decodes it|
| [logrus](https://github.com/meatballhat/negroni-logrus) | [Dan Buch](https://github.com/meatballhat) | 基于 Logrus-based logger 日志 |
| [oauth2](https://github.com/goincremental/negroni-oauth2) | [David Bochenski](https://github.com/bochenski) | oAuth2 中间件 |
| [onthefly](https://github.com/xyproto/onthefly) | [Alexander Rødseth](https://github.com/xyproto) | 快速生成 TinySVG， HTML and CSS 中间件 |
| [permissions2](https://github.com/xyproto/permissions2) | [Alexander Rødseth](https://github.com/xyproto) | Cookies， 用户和权限 |
| [prometheus](https://github.com/zbindenren/negroni-prometheus) | [Rene Zbinden](https://github.com/zbindenren) | 简易建立矩阵端点给[prometheus](http://prometheus.io)建构工具 |
| [render](https://github.com/unrolled/render) | [Cory Jacobsen](https://github.com/unrolled) | 渲染 JSON, XML and HTML 中间件 |
| [RestGate](https://github.com/pjebs/restgate) | [Prasanga Siripala](https://github.com/pjebs) | REST API 接口的安全认证 |
| [secure](https://github.com/unrolled/secure) | [Cory Jacobsen](https://github.com/unrolled) | Middleware that implements a few quick security wins |
| [sessions](https://github.com/goincremental/negroni-sessions) | [David Bochenski](https://github.com/bochenski) | Session 会话管理 |
| [stats](https://github.com/thoas/stats) | [Florent Messa](https://github.com/thoas) | 检测 web 应用当前运行状态信息 （响应时间等等。） |
| [VanGoH](https://github.com/auroratechnologies/vangoh) | [Taylor Wrobel](https://github.com/twrobel3) | Configurable [AWS-Style](http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html) 基于 HMAC 鉴权认证的中间件 |
| [xrequestid](https://github.com/pilu/xrequestid) | [Andrea Franz](https://github.com/pilu) | 给每个请求指定一个随机 X-Request-Id 头的中间件 |
| [mgo session](https://github.com/joeljames/nigroni-mgo-session) | [Joel James](https://github.com/joeljames) | 处理在每个请求建立与关闭mgo sessions |
| [digits](https://github.com/bamarni/digits) | [Bilal Amarni](https://github.com/bamarni) | 处理[Twitter Digits](https://get.digits.com/)的认证 |

## 范例
[Alexander Rødseth](https://github.com/xyproto) 创建的 [mooseware](https://github.com/xyproto/mooseware) 是一个写兼容 Negroni 中间件的处理器骨架的范例。

## 即时编译
[gin](https://github.com/codegangsta/gin) 和 [fresh](https://github.com/pilu/fresh) 这两个应用是即时编译的 Negroni 工具，推荐用户开发的时候使用。

## Go & Negroni 初学者必读推荐

* [在中间件中使用上下文把消息传递给后端处理器](http://elithrar.github.io/article/map-string-interface/)
* [了解中间件](http://mattstauffer.co/blog/laravel-5.0-middleware-replacing-filters)

## 关于

尼格龙尼正是[Code Gangsta](https://codegangsta.io/)的执着设计.

[Gorilla Mux]: https://github.com/gorilla/mux
[`http.FileSystem`]: https://godoc.org/net/http#FileSystem
