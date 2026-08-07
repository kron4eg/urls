package urls

import (
	"net/url"
	"strconv"
	"testing"
)

func TestRoute_URL(t *testing.T) {
	tests := []struct {
		name    string
		route   Route
		queries []url.Values
		want    string
	}{
		{
			name:    "no query params",
			route:   Route("/posts"),
			queries: nil,
			want:    "/posts",
		},
		{
			name:    "single query param",
			route:   Route("/posts"),
			queries: []url.Values{{"page": {"2"}}},
			want:    "/posts?page=2",
		},
		{
			name:    "multiple query params in one Values",
			route:   Route("/search"),
			queries: []url.Values{{"q": {"hello"}, "sort": {"asc"}}},
			want:    "/search?q=hello&sort=asc",
		},
		{
			name:    "multiple Values merged",
			route:   Route("/search"),
			queries: []url.Values{{"q": {"hello"}}, {"sort": {"asc"}}},
			want:    "/search?q=hello&sort=asc",
		},
		{
			name:  "same key across multiple Values appends",
			route: Route("/filter"),
			queries: []url.Values{
				{"tag": {"go"}},
				{"tag": {"templ"}},
			},
			want: "/filter?tag=go&tag=templ",
		},
		{
			name:    "value with special chars is URL-encoded",
			route:   Route("/search"),
			queries: []url.Values{{"q": {"hello world & more"}}},
			want:    "/search?q=hello+world+%26+more",
		},
		{
			name:    "empty Values skips query string",
			route:   Route("/posts"),
			queries: []url.Values{{}},
			want:    "/posts",
		},
		{
			name:    "multiple empty Values skips query string",
			route:   Route("/posts"),
			queries: []url.Values{{}, {}},
			want:    "/posts",
		},
		{
			name:    "root route with query",
			route:   Route("/"),
			queries: []url.Values{{"utm_source": {"twitter"}}},
			want:    "/?utm_source=twitter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.route.URL(tt.queries...)
			if got != tt.want {
				t.Errorf("URL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRoute_Bind(t *testing.T) {
	tests := []struct {
		name    string
		route   Route
		binding PathParams
		want    string
	}{
		{
			name:    "single placeholder",
			route:   Route("/posts/{id}"),
			binding: PathParams{"id": "42"},
			want:    "/posts/42",
		},
		{
			name:    "multiple placeholders",
			route:   Route("/orgs/{org}/repos/{repo}"),
			binding: PathParams{"org": "golang", "repo": "go"},
			want:    "/orgs/golang/repos/go",
		},
		{
			name:    "no placeholders returns route unchanged",
			route:   Route("/dashboard"),
			binding: PathParams{},
			want:    "/dashboard",
		},
		{
			name:    "empty binding with placeholders leaves placeholders",
			route:   Route("/posts/{id}"),
			binding: PathParams{},
			want:    "/posts/{id}",
		},
		{
			name:    "placeholder not in binding left as-is",
			route:   Route("/posts/{id}/comments/{cid}"),
			binding: PathParams{"id": "42"},
			want:    "/posts/42/comments/{cid}",
		},
		{
			name:    "binding key not in pattern has no effect",
			route:   Route("/posts/{id}"),
			binding: PathParams{"id": "42", "extra": "ignored"},
			want:    "/posts/42",
		},
		{
			name:    "sorted key order produces consistent results",
			route:   Route("/{b}/{a}"),
			binding: PathParams{"a": "first", "b": "second"},
			want:    "/second/first",
		},
		{
			name:    "replacement value contains special characters",
			route:   Route("/users/{name}"),
			binding: PathParams{"name": "john.doe"},
			want:    "/users/john.doe",
		},
		{
			name:    "repeated same placeholder, last sorted wins due to Replacer",
			route:   Route("/{a}/{a}"),
			binding: PathParams{"a": "final"},
			want:    "/final/final",
		},
		{
			name:    "nil binding",
			route:   Route("/posts/{id}"),
			binding: nil,
			want:    "/posts/{id}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.route.Bind(tt.binding)
			if string(got) != tt.want {
				t.Errorf("Bind() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRoute_Bind_isDeterministic(t *testing.T) {
	route := Route("/{a}/{b}/{c}")
	binding := PathParams{"c": "3", "b": "2", "a": "1"}

	for range 100 {
		got := route.Bind(binding)
		if string(got) != "/1/2/3" {
			t.Fatalf("Bind() not deterministic: got %q, want %q", got, "/1/2/3")
		}
	}
}

func TestTypedRoute_With(t *testing.T) {
	t.Run("int param", func(t *testing.T) {
		tr := TypedRoute[int]{
			Pattern: Route("/posts/{id}"),
			Build: func(r Route, p int) Route {
				return r.Bind(PathParams{"id": strconv.Itoa(p)})
			},
		}
		got := tr.With(42)
		if string(got) != "/posts/42" {
			t.Errorf("With() = %q, want %q", got, "/posts/42")
		}
	})

	t.Run("struct param", func(t *testing.T) {
		type params struct {
			Org  string
			Repo string
		}
		tr := TypedRoute[params]{
			Pattern: Route("/orgs/{org}/repos/{repo}"),
			Build: func(r Route, p params) Route {
				return r.Bind(PathParams{"org": p.Org, "repo": p.Repo})
			},
		}
		got := tr.With(params{Org: "golang", Repo: "go"})
		if string(got) != "/orgs/golang/repos/go" {
			t.Errorf("With() = %q, want %q", got, "/orgs/golang/repos/go")
		}
	})
}

func TestRoute_implementsURL(t *testing.T) {
	var _ URL = Route("/test")
}
