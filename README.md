# urls

A package for constructing URL routes with placeholder substitution and query string building.

## Usage

```go
import "github.com/kron4eg/urls"
```

### Route.Bind — placeholder substitution

```go
route := urls.Route("/orgs/{org}/repos/{repo}")
route = route.Bind(urls.PathParams{"org": "golang", "repo": "go"})
fmt.Println(route) // "/orgs/golang/repos/go"
```

Placeholder keys are processed in sorted order for deterministic output. Missing placeholders are left as-is.

### Route.URL — query string building

```go
route := urls.Route("/search")
fmt.Println(route.URL(url.Values{"q": {"hello world"}, "page": {"2"}}))
// "/search?page=2&q=hello+world"

// Multiple Values can be merged:
fmt.Println(route.URL(
    url.Values{"q": {"hello"}},
    url.Values{"page": {"2"}},
))
// "/search?page=2&q=hello"
```

### TypedRoute — typed parameter binding

```go
type Params struct {
    Org  string
    Repo string
}

repoRoute := urls.TypedRoute[Params]{
    Pattern: urls.Route("/orgs/{org}/repos/{repo}"),
    Build: func(r urls.Route, p Params) urls.Route {
        return r.Bind(urls.PathParams{"org": p.Org, "repo": p.Repo})
    },
}

fmt.Println(repoRoute.With(Params{Org: "golang", Repo: "go"}))
// "/orgs/golang/repos/go"
```

## API

```go
type URL interface {
    URL(querystring ...url.Values) string
}

type Route string
func (r Route) URL(querystring ...url.Values) string
func (r Route) Bind(binding PathParams) Route

type PathParams map[string]string

type TypedRoute[T any] struct {
    Pattern Route
    Build   func(Route, T) Route
}
func (tr TypedRoute[T]) With(params T) Route
```
