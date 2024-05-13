// Copyright (C) 2024 neocotic
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package uri provides support for constructing URI references to be used within problems.
package uri

import (
	"fmt"
	"github.com/neocotic/go-optional"
	"maps"
	"net/url"
	"slices"
	"strings"
)

// Builder is used to construct a URI reference.
type Builder struct {
	// base is the base URL to be used when building a URI reference. See Builder.Base for more information.
	base *url.URL
	// fragment is the fragment to be used when building a URI reference. See Builder.Fragment for more information.
	fragment string
	// path is the path to be used when building a URI reference. See Builder.Path for more information.
	path string
	// pathValues contains any path values to be used when building a URI reference. See Builder.PathValue for more
	// information.
	pathValues map[string]string
	// queries contains any query parameters to be used when building a URI reference. See Builder.Query for more
	// information.
	queries url.Values
	// trailingSlash controls whether a trailing slash is enforced when building a URI reference. See
	// Builder.TrailingSlash for more information.
	trailingSlash optional.Optional[bool]
}

var _ fmt.Stringer = (*Builder)(nil)

// AddQueries adds the entries from the given map to the query parameters to be used when building a URI reference.
//
// Builder.Queries can be used to replace the values of any previously set query parameter whose key is present in
// queries, if needed.
func (b *Builder) AddQueries(queries url.Values) *Builder {
	if len(queries) == 0 {
		return b
	}
	if b.queries == nil {
		b.queries = make(url.Values)
	}
	for key, values := range queries {
		for _, value := range values {
			b.queries.Add(key, value)
		}
	}
	return b
}

// AddQuery adds the given key and value to a query parameter with the given key and value to be used when building a
// URI reference.
//
// value is typically a string, otherwise it is passed to fmt.Sprint and its resulting string is used instead.
//
// Builder.Query can be used to replace the value of any previously set query parameter with the same key, if needed.
func (b *Builder) AddQuery(key string, value any) *Builder {
	if b.queries == nil {
		b.queries = make(url.Values)
	}
	switch v := value.(type) {
	case string:
		b.queries.Add(key, v)
	default:
		b.queries.Add(key, fmt.Sprint(v))
	}
	return b
}

// AddQueryf adds the given key and formatted value to a query parameter with the given key and value to be used when
// building a URI reference.
//
// Builder.Queryf can be used to replace the value of any previously set query parameter with the same key, if needed.
func (b *Builder) AddQueryf(key, format string, args ...any) *Builder {
	if b.queries == nil {
		b.queries = make(url.Values)
	}
	b.queries.Add(key, fmt.Sprintf(format, args...))
	return b
}

// Base sets the given base to be used when building a URI reference.
//
// Any query or fragment segments within base are ignored. These should be specified using the relevant Builder methods.
// Its path is retained, however, can be overridden and/or expanded using the Builder.Path and Builder.PathValue
// methods.
//
// base is typically a string or url.URL (or pointers to them), otherwise it is passed to fmt.Sprint and its resulting
// string is used instead.
//
// base is parsed as an url.URL and passed to Builder.BaseURL, however, any error returned by url.Parse is ignored
// therefore it is recommended to parse base as a url.URL and then pass it directly to Builder.BaseURL so that any error
// can be handled.
func (b *Builder) Base(base any) *Builder {
	var baseURL *url.URL
	switch v := base.(type) {
	case string:
		baseURL, _ = url.Parse(v)
	case *string:
		if v != nil {
			baseURL, _ = url.Parse(*v)
		}
	case url.URL:
		baseURL = &v
	case *url.URL:
		baseURL = v
	default:
		baseURL, _ = url.Parse(fmt.Sprint(v))
	}
	return b.BaseURL(baseURL)
}

// Basef sets the given formatted base to be used when building a URI reference.
//
// Any query or fragment segments within the formatted base are ignored. These should be specified using the relevant
// Builder methods. Its path is retained, however, can be overridden and/or expanded using the Builder.Path and
// Builder.PathValue methods.
//
// The formatted base is parsed as an url.URL and passed to Builder.BaseURL, however, any error returned by url.Parse is
// ignored therefore it is recommended to format and then parse the base as a url.URL and then pass it directly to
// Builder.BaseURL so that any error can be handled.
func (b *Builder) Basef(format string, args ...any) *Builder {
	baseURL, _ := url.Parse(fmt.Sprintf(format, args...))
	return b.BaseURL(baseURL)
}

// BaseURL sets the given base to be used when building a URI reference.
//
// Any query or fragment segments within base are ignored. These should be specified using the relevant Builder methods.
// Its path is retained, however, can be overridden and/or expanded using the Builder.Path and Builder.PathValue
// methods.
func (b *Builder) BaseURL(base *url.URL) *Builder {
	if base == nil {
		b.base = nil
		return b
	}
	clone := *base
	clone.Fragment = ""
	clone.RawFragment = ""
	clone.RawQuery = ""
	b.base = &clone
	return b
}

// Clone returns a clone of the Builder.
func (b *Builder) Clone() *Builder {
	if b == nil {
		return nil
	}
	clone := *b
	clone.pathValues = maps.Clone(b.pathValues)
	// Deep clone required for url.Values
	var query url.Values
	if b.queries != nil {
		query = make(url.Values, len(b.queries))
		for key, values := range b.queries {
			query[key] = slices.Clone(values)
		}
	}
	clone.queries = query
	return &clone
}

// Fragment sets the given fragment to be used when a building a URI reference.
//
// fragment is typically a string, otherwise it is passed to fmt.Sprint and its resulting string is used instead.
func (b *Builder) Fragment(fragment any) *Builder {
	switch v := fragment.(type) {
	case string:
		b.fragment = v
	default:
		b.fragment = fmt.Sprint(v)
	}
	return b
}

// Fragmentf sets the given formatted fragment to be used when a building a URI reference.
func (b *Builder) Fragmentf(format string, args ...any) *Builder {
	b.fragment = fmt.Sprintf(format, args...)
	return b
}

// Path sets the given path to be used when a building a URI reference.
//
// path is typically a string, otherwise it is passed to fmt.Sprint and its resulting string is used instead.
//
// When not empty, path will take precedence over that of the base, where specified.
func (b *Builder) Path(path any) *Builder {
	switch v := path.(type) {
	case string:
		b.path = v
	default:
		b.path = fmt.Sprint(v)
	}
	return b
}

// Pathf sets the given formatted path to be used when a building a URI reference.
//
// When not empty, the formatted will take precedence over that of the base, where specified.
func (b *Builder) Pathf(format string, args ...any) *Builder {
	b.path = fmt.Sprintf(format, args...)
	return b
}

// PathValue sets a path value with the given name and value to be used when building a URI reference.
//
// value is typically a string, otherwise it is passed to fmt.Sprint and its resulting string is used instead.
//
// Any path value given is replaced in the path of the URI reference when built. The name of the path value is looked up
// as with a colon (:) prefix and replaced with the value after being passed through url.PathEscape.
func (b *Builder) PathValue(name string, value any) *Builder {
	if b.pathValues == nil {
		b.pathValues = make(map[string]string)
	}
	switch v := value.(type) {
	case string:
		b.pathValues[name] = v
	default:
		b.pathValues[name] = fmt.Sprint(v)
	}
	return b
}

// PathValuef sets a path value with the given name and formatted value to be used when building a URI reference.
//
// Any path value given is replaced in the path of the URI reference when built. The name of the path value is looked up
// as with a colon (:) prefix and replaced with the value after being passed through url.PathEscape.
func (b *Builder) PathValuef(name, format string, args ...any) *Builder {
	if b.pathValues == nil {
		b.pathValues = make(map[string]string)
	}
	b.pathValues[name] = fmt.Sprintf(format, args...)
	return b
}

// PathValues sets the path values with the entries within the given map to be used when building a URI reference.
//
// Any path value given is replaced in the path of the URI reference when built. The name of the path value is looked up
// as with a colon (:) prefix and replaced with the value after being passed through url.PathEscape.
func (b *Builder) PathValues(pathValues map[string]string) *Builder {
	if len(pathValues) == 0 {
		b.pathValues = nil
		return b
	}
	if b.pathValues == nil {
		b.pathValues = make(map[string]string)
	}
	for name, value := range pathValues {
		b.pathValues[name] = value
	}
	return b
}

// Queries sets the query parameters with the entries within the given map to be used when building a URI reference.
//
// While queries is a multi-value map, Queries will replace the values of any previously set query parameter whose key
// is present in queries. Builder.AddQueries can be used to add additional values to the same keys, if needed.
func (b *Builder) Queries(queries url.Values) *Builder {
	if len(queries) == 0 {
		b.queries = nil
		return b
	}
	if b.queries == nil {
		b.queries = make(url.Values)
	}
	for key, values := range queries {
		b.queries[key] = values
	}
	return b
}

// Query sets a query parameter with the given key and value to be used when building a URI reference.
//
// value is typically a string, otherwise it is passed to fmt.Sprint and its resulting string is used instead.
//
// Query will replace the value of any previously set query parameter with the same key. Builder.AddQuery can be used to
// add additional values to the same key, resulting in a multi-value query parameter, if needed.
func (b *Builder) Query(key string, value any) *Builder {
	if b.queries == nil {
		b.queries = make(url.Values)
	}
	switch v := value.(type) {
	case string:
		b.queries.Set(key, v)
	default:
		b.queries.Set(key, fmt.Sprint(v))
	}
	return b
}

// Queryf sets a query parameter with the given key and formatted value to be used when building a URI reference.
//
// Queryf will replace the value of any previously set query parameter with the same key. Builder.AddQueryf can be used
// to add additional formatted values to the same key, resulting in a multi-value query parameter, if needed.
func (b *Builder) Queryf(key, format string, args ...any) *Builder {
	if b.queries == nil {
		b.queries = make(url.Values)
	}
	b.queries.Set(key, fmt.Sprintf(format, args...))
	return b
}

// Reset clears all information used to build a URI reference.
func (b *Builder) Reset() *Builder {
	b.base = nil
	b.fragment = ""
	b.path = ""
	b.pathValues = nil
	b.queries = nil
	b.trailingSlash = optional.Empty[bool]()
	return b
}

// String constructs a URI reference and returns its string representation.
func (b *Builder) String() string {
	return b.URL().String()
}

// TrailingSlash sets whether a trailing slash is enforced at the end of the path when building a URI reference.
//
// By default, the presence of such a trailing slash is entirely optional and not added or removed when constructed.
//
// If nothing is passed to TrailingSlash, this is considered equal to passing true. If true, a trailing slash will be
// enforced at the end of the path. If false, any trailing slash at the end of the pth will be removed.
func (b *Builder) TrailingSlash(trailingSlash ...bool) *Builder {
	if len(trailingSlash) > 0 {
		b.trailingSlash = optional.Of(trailingSlash[0])
	} else {
		b.trailingSlash = optional.Of(true)
	}
	return b
}

// URL constructs a URI reference and returns its url.URL representation.
func (b *Builder) URL() *url.URL {
	var u *url.URL
	if b.base != nil {
		clone := *b.base
		u = &clone
	} else {
		u = &url.URL{}
	}
	if b.path != "" {
		u.Path = b.path
	}
	if !u.IsAbs() && u.Path == "" || u.Path[0] != '/' {
		u.Path = "/" + u.Path
	}
	u.Fragment = b.fragment
	if len(b.pathValues) > 0 {
		var oldNew []string
		for name, value := range b.pathValues {
			oldNew = append(oldNew, ":"+name, url.PathEscape(value))
		}
		u.Path = strings.NewReplacer(oldNew...).Replace(u.Path)
	}
	if trailingSlash, present := b.trailingSlash.Get(); present {
		pl := len(u.Path)
		hasTrailingSlash := u.Path[pl-1] == '/'
		switch {
		case trailingSlash && !hasTrailingSlash:
			u.Path += "/"
		case !trailingSlash && hasTrailingSlash:
			u.Path = u.Path[:pl-1]
		}
	}
	u.RawQuery = b.queries.Encode()
	return u
}

// Build returns a Builder.
func Build() Builder {
	return Builder{}
}
