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

import "context"

// contextKey is an internal type for managing key/value pairs within a context.Context without conflicting with other
// packages.
type contextKey uint

// contextKeyGenerator is the key associated with a Generator within a context.Context.
const contextKeyGenerator contextKey = 0

// GetGenerator returns the Generator within the given context.Context, otherwise DefaultGenerator.
func GetGenerator(ctx context.Context) *Generator {
	if gen, ok := ctx.Value(contextKeyGenerator).(*Generator); ok && gen != nil {
		return gen
	}
	return DefaultGenerator
}

// UsingGenerator returns a copy of the given parent context.Context containing the Generator provided.
//
// If gen is nil, DefaultGenerator is used.
func UsingGenerator(parent context.Context, gen *Generator) context.Context {
	if gen == nil {
		gen = DefaultGenerator
	}
	return context.WithValue(parent, contextKeyGenerator, gen)
}
