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

// Definition represents a reusable definition of problem occurrence that may contain default values that can be used
// when generating a Problem from a specific definition.
//
// A Definition provides an alternative approach to generating problems where they can be used to dictate all
// information populated within a Problem and/or combined with options to provide more granular control and overrides.
type Definition struct {
	// Code is the default code to be assigned to a Problem generated from the Definition. See Problem.Code for more
	// information.
	//
	// If Code is empty, no default is used.
	Code Code `json:"code" xml:"code" yaml:"code"`
	// Detail is the default detail to be assigned to a Problem generated from the Definition. See Problem.Detail for
	// more information.
	//
	// If Detail is empty, no default is used.
	//
	// If DetailKey is not empty and can be resolved it will take precedence over Detail.
	Detail string `json:"detail" xml:"detail" yaml:"detail"`
	// DetailKey is the translation key of the default detail to be assigned to a Problem generated from the Definition.
	// See Problem.Detail for more information.
	//
	// The localized detail will be looked up using Generator.Translator, where possible. If resolved, it will take
	// precedence over Detail.
	//
	// If DetailKey is empty it, it is ignored.
	DetailKey any `json:"detailKey" xml:"detailKey" yaml:"detailKey"`
	// Extensions is the default extensions to be assigned to a Problem generated from the Definition. See
	// Problem.Extensions for more information.
	//
	// If Extensions is nil, no default is used.
	Extensions map[string]any `json:"extensions" xml:"extensions" yaml:"extensions"`
	// Instance is the default instance URI to be assigned to a Problem generated from the Definition. See
	// Problem.Instance for more information.
	//
	// An uri.Builder can be used to aid building the URI reference.
	//
	// If Instance is empty, no default is used.
	Instance string `json:"instance" xml:"instance" yaml:"instance"`
	// Type contains fields defining the type of Problem generated from the Definition, typically containing additional
	// default values.
	Type Type `json:"type" xml:"type" yaml:"type"`
}

// Build is a convenient shorthand for calling Generator.Build on DefaultGenerator with the Definition already passed to
// Builder.Definition.
func (d Definition) Build() *Builder {
	return &Builder{
		Generator: DefaultGenerator,
		ctx:       optional.Of(context.Background()),
		def:       d,
	}
}

// BuildContext is a convenient shorthand for calling Generator.BuildContext on the Generator within the given
// context.Context, if any, otherwise DefaultGenerator, with the Definition already passed to Builder.Definition.
func (d Definition) BuildContext(ctx context.Context) *Builder {
	return &Builder{
		Generator: GetGenerator(ctx),
		ctx:       optional.Of(ctx),
		def:       d,
	}
}

// BuildContextUsing is a convenient shorthand for calling Generator.BuildContext with the Definition already passed to
// Builder.Definition.
func (d Definition) BuildContextUsing(ctx context.Context, gen *Generator) *Builder {
	return &Builder{
		Generator: gen,
		ctx:       optional.Of(ctx),
		def:       d,
	}
}

// BuildUsing is a convenient shorthand for calling Generator.Build with the Definition already passed to
// Builder.Definition.
func (d Definition) BuildUsing(gen *Generator) *Builder {
	return &Builder{
		Generator: gen,
		ctx:       optional.Of(context.Background()),
		def:       d,
	}
}

// New is a convenient shorthand for calling Generator.New on DefaultGenerator, including FromDefinition with the
// Definition along with any specified options.
func (d Definition) New(opts ...Option) *Problem {
	opts = append([]Option{FromDefinition(d)}, opts...)
	return DefaultGenerator.new(context.Background(), opts, 1)
}

// NewContext is a convenient shorthand for calling Generator.NewContext on the Generator within the given
// context.Context, if any, otherwise DefaultGenerator, including FromDefinition with the Definition along with any
// specified options.
func (d Definition) NewContext(ctx context.Context, opts ...Option) *Problem {
	opts = append([]Option{FromDefinition(d)}, opts...)
	return GetGenerator(ctx).new(ctx, opts, 1)
}

// NewContextUsing is an alternative for calling Generator.NewContext, including FromDefinition with the Definition
// along with any specified options.
func (d Definition) NewContextUsing(ctx context.Context, gen *Generator, opts ...Option) *Problem {
	opts = append([]Option{FromDefinition(d)}, opts...)
	return gen.new(ctx, opts, 1)
}

// NewUsing is an alternative for calling Generator.New, including FromDefinition with the Definition along with any
// specified options.
func (d Definition) NewUsing(gen *Generator, opts ...Option) *Problem {
	opts = append([]Option{FromDefinition(d)}, opts...)
	return gen.new(context.Background(), opts, 1)
}
