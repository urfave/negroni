# Negroni
[![GoDoc](https://godoc.org/github.com/urfave/negroni?status.svg)](http://godoc.org/github.com/urfave/negroni)
[![Build Status](https://travis-ci.org/urfave/negroni.svg?branch=master)](https://travis-ci.org/urfave/negroni)
[![codebeat](https://codebeat.co/badges/47d320b1-209e-45e8-bd99-9094bc5111e2)](https://codebeat.co/projects/github-com-urfave-negroni)
[![codecov](https://codecov.io/gh/urfave/negroni/branch/master/graph/badge.svg)](https://codecov.io/gh/urfave/negroni)

**Note:** Ce projet était initiallement connu comme
`github.com/codegangsta/negroni` -- Github redirigera automatiquement les requêtes vers ce dépôt.
Nous vous recommandons néanmoins d'utiliser la référence vers ce nouveau dépôt pour plus de clarté.

Negroni approche la question de la création de *middleware* de manière pragmatique.
La librairie se veut légère, non intrusive et encourage l'utilisation des *Handlers* de
la librairie standard `net/http`.

Si vous appréciez le projet [Martini](https://github.com/go-martini/martini) et estimez
qu'une certaine magie s'en dégage, Negroni sera sans doute plus approprié.

## Démarrer avec Negroni

Une fois Go installé et votre variable [GOPATH](http://golang.org/doc/code.html#GOPATH) à jour,
créez votre premier fichier `.go` et nommez le `server.go`.

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

  n := negroni.Classic() // Inclue les "middlewares" par défaut.
  n.UseHandler(mux)

  http.ListenAndServe(":3000", n)
}
```

Installez au préalable le paquet Negroni (**NOTE**: Une version de Go &gt;= **go 1.1** est nécessaire):

```
go get github.com/urfave/negroni
```

Démarrez le serveur:

```
go run server.go
```

Vous avez dès à présent un serveur web Go basé sur `net/http` disponible à l'adresse `localhost:3000`.

### Paquets

Si vous utilisez Debian, `negroni` est aussi disponible en tant que [paquet]
(https://packages.debian.org/sid/golang-github-urfave-negroni-dev).
La commande `apt install golang-github-urfave-negroni-dev` vous permettra de l'installer
(À ce jour, vous les trouverez dans les dépôts `sid`).

## Negroni est-il un *framework* ?

Negroni **n'est pas** un *framework*. Considérez le comme une librairie centrée sur
l'utilisation de *middleware* développés pour fonctionner directement avec la librairie `net/http`.

## Redirection (*Routing*) ?

Negroni est *BYOR* (*Bring your own Router*, Apporter votre propre routeur).
La communauté Go offre un nombre importants de routeur et Negroni met tout en oeuvre
pour fonctionner avec chacun d'entre eux en assurant un support complet de la librairie `net/http`.
Par exemple, une utilisation avec [Gorilla Mux] se présente sous la forme:

``` go
router := mux.NewRouter()
router.HandleFunc("/", HomeHandler)

n := negroni.New(Middleware1, Middleware2)
// on peut utiliser également la fonction Use() pour ajouter un "middleware"
n.Use(Middleware3)
// le routeur se trouve toujours en dernier.
n.UseHandler(router)

http.ListenAndServe(":3001", n)
```

## `negroni.Classic()`

L'instance `negroni.Classic()` propose par défaut trois middlewares qui seront utiles à la plupart
des applications:

* [`negroni.Recovery`](#recovery) - Récupère des appels à `panic`.
* [`negroni.Logger`](#logger) - Journalise les requêtes et les réponses.
* [`negroni.Static`](#static) - Sers les fichiers statiques présent dans le dossier "public".

Elle offre un démarrage aisé sans recourir à la configuration pour utiliser quelques une des fonctions les plus utiles de Negroni.

## *Handlers*

Negroni offre un flux bidirectionnel via les *middlewares* (aller-retour entre la requête et la réponse). Tout repose sur l'interface `negroni.Handler` :

 ``` go
type Handler interface {
  ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
```

Si un *middleware* n'a pas écrit au `ResponseWriter`, il doit faire appel au prochain `http.Handlerfunc` de la chaîne pour que le prochain *middleware* soit appelé:

``` go
func MyMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // faire quelque chose avant
  next(rw, r)
  // faire quelque chose après
}
```

Vous pouvez insérer votre *middleware* dans la chaîne en utilisant la fonction `Use`:

``` go
n := negroni.New()
n.Use(negroni.HandlerFunc(MyMiddleware))
```

Vous pour également utiliser un `http.Handler` classique:

``` go
n := negroni.New()

mux := http.NewServeMux()
// définissez vos routes

n.UseHandler(mux)

http.ListenAndServe(":3000", n)
```

## `With()`

La méthode `With()` vous permet de regrouper un ou plusieurs `Handler` au sein
d'une nouvelle instance `Negroni`. Cette dernière est la combinaison des `Handlers` de l'ancienne et de la nouvelle instance.

```go
// "middleware" à réutiliser.
common := negroni.New()
common.Use(MyMiddleware1)
common.Use(MyMiddleware2)

// `specific` devient une nouvelle instance avec les "handlers" provenant de `common` ainsi
// que ceux passés en paramètres.
specific := common.With(
	SpecificMiddleware1,
	SpecificMiddleware2
)
```

## `Run()`

Negroni peut-être démarrer en utilisant la méthode `Run()`. Cette dernière
prend en paramètre l'adresse du serveur, à l'instar de la méthode [`http.ListenAndServe`](https://godoc.org/net/http#ListenAndServe).

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
Si aucune adresse n'est renseignée, la variable d'environnement `PORT` est utilisée.
Si cette dernière n'est pas définie, l'adresse par défaut est utilisée.
Pour une description détaillée, veuillez-vous référer à la documentation de la méthode
[Run]((https://godoc.org/github.com/urfave/negroni#Negroni.Run).

De manière générale, vous voudrez vous servir de la librairie `net/http` et utiliser `negroni`
comme un simple `Handler` pour plus de flexibilité.

Par exemple:

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

  n := negroni.Classic() // Inclue les "middlewares" par défaut
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

## Redirection spécifique

Si un ensemble de routes nécessite l'appel à des *middleware* spécifiques,
vous pouvez simplement créer une nouvelle instance Negroni et l'utiliser comme
`Handler` pour cet ensemble.

``` go
router := mux.NewRouter()
adminRoutes := mux.NewRouter()
// ajout des routes relatives à l'administration

// Création d'une nouvelle instance pour le "middleware" admin.
router.PathPrefix("/admin").Handler(negroni.New(
  Middleware1,
  Middleware2,
  negroni.Wrap(adminRoutes),
))
```

Si vous utilisez [Gorilla Mux], vous pourriez utiliser un *subrouter*:

``` go
router := mux.NewRouter()
subRouter := mux.NewRouter().PathPrefix("/subpath").Subrouter().StrictSlash(true)
subRouter.HandleFunc("/", someSubpathHandler) // "/subpath/"
subRouter.HandleFunc("/:id", someSubpathHandler) // "/subpath/:id"

// "/subpath" est nécessaire pour assurer la cohésion entre le `subrouter` et le routeur principal
router.PathPrefix("/subpath").Handler(negroni.New(
  Middleware1,
  Middleware2,
  negroni.Wrap(subRouter),
))
```

La méthode `With()` peut aider à réduire la duplication des *middlewares* partagés par
plusieurs routes.

``` go
router := mux.NewRouter()
apiRoutes := mux.NewRouter()
// ajout des routes api ici
webRoutes := mux.NewRouter()
// ajout des routes web ici

// création un "middleware" commun pour faciliter le partage
common := negroni.New(
	Middleware1,
	Middleware2,
)

// création d'une nouvelle instance pour le "middleware"
// api en utilisant le "middleware" commun comme base.
router.PathPrefix("/api").Handler(common.With(
  APIMiddleware1,
  negroni.Wrap(apiRoutes),
))
// création d'une nouvelle instance pour le "middleware"
// web en utilisant le "middleware" commun comme base.
router.PathPrefix("/web").Handler(common.With(
  WebMiddleware1,
  negroni.Wrap(webRoutes),
))
```

## *Middlewares* fournis

### Static

Ce *middleware* va servir les fichiers présents sur le système de fichiers.
Si un fichier n'existe pas, il transmet la requête au *middleware* suivant.
Si vous souhaitez retourner le message `404 File Not Found` pour les fichiers non existants,
vous pouvez utiliser la fonction [http.FileServer](https://golang.org/pkg/net/http/#FileServer)
comme `Handler`.

Exemple:

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

  // Exemple d'usage de la fonction http.FileServer pour avoir un comportement similaire à un
  // serveur HTTP "standard" plutôt que le comportement "middleware"
  // mux.Handle("/public", http.FileServer(http.Dir("/home/public")))

  n := negroni.New()
  n.Use(negroni.NewStatic(http.Dir("/tmp")))
  n.UseHandler(mux)

  http.ListenAndServe(":3002", n)
}
```

Ce programme servira les fichiers depuis le dossier `/tmp` en premier lieu.
Si le fichier n'est pas trouvé, il transmet la requête au *middleware* suivant.

### Recupération (*Recovery*)

Ce *middleware* capture les appels à `panic` et renvoie une réponse `500` à
la requête correspondante. Si un autre *middleware* a déjà renvoyé une réponse (vide ou non),
le renvoie de la réponse `500` au client échouera, le client en ayant déjà obtenu une.

Il est possible d'adjoindre au *middleware* une fonction de type `PanicHandlerFunc`
pour collecter les erreurs `500` et les transmettre à un service de rapport d'erreur
tels Sentry ou Airbrake.

Exemple:

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

Ce programme renverra une erreur `500 Internal Server Error` à chaque requête reçue.
Il transmettra à son *logger* associé la trace de la pile d'exécution et affichera cette même trace sur la sortie standard si la valeur `PrintStack` est mise à `true`. (valeur par défaut)

Exemple avec l'utilisation d'une `PanicHandlerFunc`:

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
    // code envoyant le rapport d'erreur à Sentry
}
```

## Journalisation (*Logger*)

Ce *middleware* va *logger* toutes les requêtes et réponses.

Exemple:

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

Ce programme affichera un *log* similaire à celui-ci pour chaque requête/réponse:

```
[negroni] 2017-10-04T14:56:25+02:00 | 200 |      378µs | localhost:3004 | GET /
```

Il est possible de modifier le format par défaut en utilisant la fonction `SetFormat`.
Le format est `template` dont les champs associés sont les propriétés de l'objet `LoggerEntry`.

Par exemple:

```go
l.SetFormat("[{{.Status}} {{.Duration}}] - {{.Request.UserAgent}}")
```

Ce format proposera un affichage similaire à: `[200 18.263µs] - Go-User-Agent/1.1 `

## *Middlewares* tiers

Vous trouverez ici une liste de *middlewares* compatibles avec Negroni.
N'hésitez pas à créer une PR pour renseigner un middleware de votre cru:

| Middleware | Author | Description |
| -----------|--------|-------------|
| [authz](https://github.com/casbin/negroni-authz) | [Yang Luo](https://github.com/hsluoyz) | Un *middleware* de gestion d'accès ACL, RBAC, ABAC basé sur [Casbin](https://github.com/casbin/casbin) |
| [binding](https://github.com/mholt/binding) | [Matt Holt](https://github.com/mholt) | Associez facilement les données des requêtes HTTP vers des structures Go |
| [cloudwatch](https://github.com/cvillecsteele/negroni-cloudwatch) | [Colin Steele](https://github.com/cvillecsteele) | *Middleware* pour utiliser les métriques AWS cloudwatch |
| [cors](https://github.com/rs/cors) | [Olivier Poitrey](https://github.com/rs) | Support [Cross Origin Resource Sharing](http://www.w3.org/TR/cors/) (CORS) |
| [csp](https://github.com/awakenetworks/csp) | [Awake Networks](https://github.com/awakenetworks) | Support [Content Security Policy](https://www.w3.org/TR/CSP2/) (CSP) |
| [delay](https://github.com/jeffbmartinez/delay) | [Jeff Martinez](https://github.com/jeffbmartinez) | Ajouter des délais de réponse sur les routes. Utile pour tester les effets de la latence. |
| [New Relic Go Agent](https://github.com/yadvendar/negroni-newrelic-go-agent) | [Yadvendar Champawat](https://github.com/yadvendar) | [Agent New Relic Go](https://github.com/newrelic/go-agent) officiel |
| [gorelic](https://github.com/jingweno/negroni-gorelic) | [Jingwen Owen Ou](https://github.com/jingweno) | Agent New Relic agent pour le runtime Go |
| [Graceful](https://github.com/tylerb/graceful) | [Tyler Bunnell](https://github.com/tylerb) | Graceful HTTP Shutdown |
| [gzip](https://github.com/phyber/negroni-gzip) | [phyber](https://github.com/phyber) | Compression GZIP des réponses |
| [JWT Middleware](https://github.com/auth0/go-jwt-middleware) | [Auth0](https://github.com/auth0) | Middleware vérifiant la présence d'un JWT dans le *header* `Authorization` et le décode |
| [logrus](https://github.com/meatballhat/negroni-logrus) | [Dan Buch](https://github.com/meatballhat) | *Logger* basé sur Logrus |
| [oauth2](https://github.com/goincremental/negroni-oauth2) | [David Bochenski](https://github.com/bochenski) | Middleware oAuth2 |
| [onthefly](https://github.com/xyproto/onthefly) | [Alexander Rødseth](https://github.com/xyproto) | Générer des éléments TinySVG, HTML et CSS à la volée |
| [permissions2](https://github.com/xyproto/permissions2) | [Alexander Rødseth](https://github.com/xyproto) | Cookies, utilisateurs et permissions |
| [prometheus](https://github.com/zbindenren/negroni-prometheus) | [Rene Zbinden](https://github.com/zbindenren) | Créer des métriques facilement avec l'outil [prometheus](http://prometheus.io) |
| [render](https://github.com/unrolled/render) | [Cory Jacobsen](https://github.com/unrolled) | Rendre des templates JSON, XML et HTML |
| [RestGate](https://github.com/pjebs/restgate) | [Prasanga Siripala](https://github.com/pjebs) | Authentification sécurisée pour les APIs REST |
| [secure](https://github.com/unrolled/secure) | [Cory Jacobsen](https://github.com/unrolled) | Middleware implémentant des basiques de sécurité |
| [sessions](https://github.com/goincremental/negroni-sessions) | [David Bochenski](https://github.com/bochenski) | Gestions des sessions |
| [stats](https://github.com/thoas/stats) | [Florent Messa](https://github.com/thoas) | Stockez des informations à propos de votre application web (temps de réponse, etc.) |
| [VanGoH](https://github.com/auroratechnologies/vangoh) | [Taylor Wrobel](https://github.com/twrobel3) | *Middleware* d'authentification HMAC configurable basée sur [AWS-Style](http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html) |
| [xrequestid](https://github.com/pilu/xrequestid) | [Andrea Franz](https://github.com/pilu) | Un *middleware* qui assigne un *header* `X-Request-Id` à chaque requête |
| [mgo session](https://github.com/joeljames/nigroni-mgo-session) | [Joel James](https://github.com/joeljames) | Un *middleware* qui gère les sessions mgo pour chaque requête (ouverture, fermeture) |
| [digits](https://github.com/bamarni/digits) | [Bilal Amarni](https://github.com/bamarni) | Un *middleware* qui gère l'authentification via [Twitter Digits](https://get.digits.com/) |
| [stats](https://github.com/guptachirag/stats) | [Chirag Gupta](https://github.com/guptachirag/stats) |
Middleware qui gère les statistiques qps et latence pour vos points de terminaison et les envoie de manière asynchrone à influx db |

## Exemples

[Alexander Rødseth](https://github.com/xyproto) a créé
[mooseware](https://github.com/xyproto/mooseware), un squelette pour écrire un *middleware* Negroni.

[Prasanga Siripala](https://github.com/pjebs) a créé un squelette pour les applications web basées sur Go et Negroni: [Go-Skeleton](https://github.com/pjebs/go-skeleton)

## Rechargement automatique du code ?

[gin](https://github.com/codegangsta/gin) et
[fresh](https://github.com/pilu/fresh) permettent tous deux de recharger les applications Negroni
suite à une modification opérée dans le code.

## Lectures pour les débutants avec Go et Negroni

* [Using a Context to pass information from middleware to end handler](http://elithrar.github.io/article/map-string-interface/)
* [Understanding middleware](https://mattstauffer.co/blog/laravel-5.0-middleware-filter-style)

## À propos

Negroni est obsessivement développé par nulle autre personne que [Code
Gangsta](https://codegangsta.io/)

[Gorilla Mux]: https://github.com/gorilla/mux
[`http.FileSystem`]: https://godoc.org/net/http#FileSystem
