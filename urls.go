package urls

import (
	"maps"
	"net/url"
	"slices"
	"strings"
)

const pairsEntryCount = 2

// URL is a type that can produce a URL string, optionally with query parameters.
type URL interface {
	URL(querystring ...url.Values) string
}

// Route is an absolute URL path that may contain {placeholder} segments.
// Call URL() to get the final string, or Bind() to fill the placeholders.
type Route string

func (r Route) appendQuery(queries ...url.Values) string {
	if len(queries) == 0 {
		return string(r)
	}

	var qstring strings.Builder

	for _, q := range queries {
		if encoded := q.Encode(); encoded != "" {
			if qstring.Len() > 0 {
				qstring.WriteString("&")
			}
			qstring.WriteString(encoded)
		}
	}

	if qstring.Len() > 0 {
		return string(r) + "?" + qstring.String()
	}

	return string(r)
}

func (r Route) URL(query ...url.Values) string {
	return r.appendQuery(query...)
}

// Bind replaces {key} placeholders in the route with corresponding values from binding.
// Keys are processed in sorted order for deterministic output.
func (r Route) Bind(binding PathParams) Route {
	pairs := make([]string, 0, len(binding)*pairsEntryCount)
	for _, k := range slices.Sorted(maps.Keys(binding)) {
		v := binding[k]
		pairs = append(pairs, "{"+k+"}", v)
	}
	return Route(strings.NewReplacer(pairs...).Replace(string(r)))
}

// PathParams maps placeholder names to their replacement values for Bind.
type PathParams map[string]string

// TypedRoute is a parameterized route with a typed parameter T.
// Pattern is the route template with {placeholders}, and Build is a function
// that fills those placeholders from a value of type T to produce a concrete Route.
type TypedRoute[T any] struct {
	Pattern Route

	Build func(Route, T) Route
}

// With fills the route placeholders from params and returns the resulting Route.
func (tr TypedRoute[T]) With(params T) Route {
	return Route(tr.Build(tr.Pattern, params))
}
