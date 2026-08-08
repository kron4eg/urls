# urls

A package for type safe URL routes.

## API docs

[![Go Reference](https://pkg.go.dev/badge/github.com/kron4eg/urls.svg)](https://pkg.go.dev/github.com/kron4eg/urls)

## Usage

```go
-- routes.go --
import "github.com/kron4eg/urls"

var (
	Index      = urls.Route("/")
	Posts      = urls.Route("/posts")
	PostDetail = urls.TypedRoute[PostDetailParams]{
		Pattern: urls.Route("/posts/{id}"),
		Build: func(r urls.Route, p PostDetailParams) urls.Route {
			return r.Bind(urls.PathParams{"id": strconv.Itoa(p.ID)})
		},
	}
)

// PostDetailParams are the typed parameters for the PostDetail route.
type PostDetailParams struct {
	ID int
}

-- server.go --
func newMux() *http.ServerMux {
	mux := http.NewServeMux()
	mux.HandleFunc(Index.String(), http.HandlerFunc(Index))
	mux.HandleFunc(Posts.String(), http.HandlerFunc(Posts))
	mux.HandleFunc(PostDetail.Pattern.String(), http.HandlerFunc(h.PostDetail))
}

-- handlers.go --
func RedirectToIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, Index.URL(), http.StatusFound)
}

func Posts(w http.ResponseWriter, r *http.Request) {
	// look Ma, no string literals in revers URLs
	typesafePost42URL := PostDetail.With(PostDetailParams{ID: 42}).URL()
	http.Redirect(w, r, typesafePost42URL, http.StatusFound)
}
```

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
