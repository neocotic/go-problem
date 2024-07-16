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

package problem

import (
	"context"
	"github.com/neocotic/go-optional"
)

type (
	// Type represents a reusable problem type that may contain default values that can be used when generating a Problem
	// from a specific type.
	//
	// A Type, optionally used in combination with a Definition, provides an alternative approach to generating problems
	// where they can be used to dictate all information populated within a Problem and/or combined with options to
	// provide more granular control and overrides.
	Type struct {
		// LogLevel is the default LogLevel to be assigned to a Problem generated from the Type. See Problem.LogLevel for
		// more information.
		//
		// The LogLevel may be overridden by Generator.LogLeveler.
		//
		// If LogLevel is zero, the default used is DefaultLogLevel.
		LogLevel LogLevel `json:"logLevel" xml:"logLevel" yaml:"logLevel"`
		// Status is the default status to be assigned to a Problem generated from the Type. See Problem.Status for more
		// information.
		//
		// If Status is zero, the default used is http.StatusInternalServerError.
		Status int `json:"status" xml:"status" yaml:"status"`
		// Title is the default title to be assigned to a Problem generated from the Type. See Problem.Title for more
		// information.
		//
		// If Title is empty, the default used is DefaultTitle.
		//
		// If TitleKey is not empty and can be resolved it will take precedence over Title.
		Title string `json:"title" xml:"title" yaml:"title"`
		// TitleKey is the translation key of the default title to be assigned to a Problem generated from the Type. See
		// Problem.Title for more information.
		//
		// The localized title will be looked up using Generator.Translator, where possible. If resolved, it will take
		// precedence over Title.
		//
		// If TitleKey is empty it, it is ignored.
		TitleKey any `json:"titleKey" xml:"titleKey" yaml:"titleKey"`
		// URI is the default type URI to be assigned to a Problem generated from the Type. See Problem.Type for more
		// information.
		//
		// An uri.Builder can be used to aid building the URI reference.
		//
		// If URI is empty, the default used is DefaultTypeURI.
		URI string `json:"uri" xml:"uri" yaml:"uri"`
	}

	// Typer is a function that can be used by a Generator to override the type URI reference derived from a Type (i.e.
	// instead of only Type.URI).
	//
	// An uri.Builder can be used to aid building the URI reference.
	//
	// It is important to note that if the function returns an empty string, it will fall back to DefaultTypeURI and not
	// Type.URI.
	Typer func(defType Type) string
)

const (
	// ContentTypeJSON is the recommended content/media type to represent a problem in JSON format.
	ContentTypeJSON = "application/problem+json"
	// ContentTypeJSONUTF8 is the recommended content/media type to represent a problem in JSON format with UTF-8
	// encoding.
	ContentTypeJSONUTF8 = ContentTypeJSON + "; charset=utf-8"
	// ContentTypeXML is the recommended content/media type to represent a Problem in XML format.
	ContentTypeXML = "application/problem+xml"
	// ContentTypeXMLUTF8 is the recommended content/media type to represent a problem in XML format with UTF-8
	// encoding.
	ContentTypeXMLUTF8 = ContentTypeXML + "; charset=utf-8"

	// DefaultTypeURI is the default problem type URI, indicating that a problem has no additional semantics beyond that
	// its status.
	//
	// Typically, when used, the problem title SHOULD be the same as the recommended HTTP status text for that code
	// (e.g. "Not Found" for 404).
	DefaultTypeURI = "about:blank"
)

// Build is a convenient shorthand for calling Generator.Build on DefaultGenerator with the Type already passed to
// Builder.DefinitionType.
func (t Type) Build() *Builder {
	return &Builder{
		Generator: DefaultGenerator,
		ctx:       optional.Of(context.Background()),
		def:       Definition{Type: t},
	}
}

// BuildContext is a convenient shorthand for calling Generator.BuildContext on the Generator within the given
// context.Context, if any, otherwise DefaultGenerator, with the Type already passed to Builder.DefinitionType.
func (t Type) BuildContext(ctx context.Context) *Builder {
	return &Builder{
		Generator: GetGenerator(ctx),
		ctx:       optional.Of(ctx),
		def:       Definition{Type: t},
	}
}

// BuildContextUsing is a convenient shorthand for calling Generator.BuildContext with the Type already passed to
// Builder.DefinitionType.
func (t Type) BuildContextUsing(ctx context.Context, gen *Generator) *Builder {
	return &Builder{
		Generator: gen,
		ctx:       optional.Of(ctx),
		def:       Definition{Type: t},
	}
}

// BuildUsing is a convenient shorthand for calling Generator.Build with the Type already passed to
// Builder.DefinitionType.
func (t Type) BuildUsing(gen *Generator) *Builder {
	return &Builder{
		Generator: gen,
		ctx:       optional.Of(context.Background()),
		def:       Definition{Type: t},
	}
}

// New is a convenient shorthand for calling Generator.New on DefaultGenerator, including FromType with the Type along
// with any specified options.
func (t Type) New(opts ...Option) *Problem {
	opts = append([]Option{FromType(t)}, opts...)
	return DefaultGenerator.new(context.Background(), opts, 1)
}

// NewContext is a convenient shorthand for calling Generator.NewContext on the Generator within the given
// context.Context, if any, otherwise DefaultGenerator, including FromType with the Type along with any specified
// options.
func (t Type) NewContext(ctx context.Context, opts ...Option) *Problem {
	opts = append([]Option{FromType(t)}, opts...)
	return GetGenerator(ctx).new(ctx, opts, 1)
}

// NewContextUsing is an alternative for calling Generator.NewContext, including FromType with the Type along with any
// specified options.
func (t Type) NewContextUsing(ctx context.Context, gen *Generator, opts ...Option) *Problem {
	opts = append([]Option{FromType(t)}, opts...)
	return gen.new(ctx, opts, 1)
}

// NewUsing is an alternative for calling Generator.New, including FromType with the Type along with any specified
// options.
func (t Type) NewUsing(gen *Generator, opts ...Option) *Problem {
	opts = append([]Option{FromType(t)}, opts...)
	return gen.new(context.Background(), opts, 1)
}

// contentType returns Generator.ContentType if not empty and valid, otherwise ContentTypeJSONUTF8.
func (g *Generator) contentType() string {
	if g.ContentType != "" && isValidContentType(g.ContentType) {
		return g.ContentType
	}
	return ContentTypeJSONUTF8
}

// typeURI checks if Generator.Typer is present and, if so, calls it with the given Type to allow for the type URI
// reference to be overridden, where appropriate. Otherwise, Type.URI is returned.
func (g *Generator) typeURI(defType Type) string {
	if t := g.Typer; t != nil {
		return t(defType)
	}
	return defType.URI
}

// isValidContentType returns whether the given content-type is valid when representing a Problem in any supported form.
func isValidContentType(ct string) bool {
	switch ct {
	case ContentTypeJSON, ContentTypeJSONUTF8, ContentTypeXML, ContentTypeXMLUTF8:
		return true
	default:
		return false
	}
}

// isValidContentTypeForJSON returns whether the given content-type is valid when representing a Problem in its JSON
// form.
func isValidContentTypeForJSON(ct string) bool {
	switch ct {
	case ContentTypeJSON, ContentTypeJSONUTF8:
		return true
	default:
		return false
	}
}

// isValidContentTypeForXML returns whether the given content-type is valid when representing a Problem in its XML form.
func isValidContentTypeForXML(ct string) bool {
	switch ct {
	case ContentTypeXML, ContentTypeXMLUTF8:
		return true
	default:
		return false
	}
}
