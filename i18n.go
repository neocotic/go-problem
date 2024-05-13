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

// Translator is a function that returns a localized value based on the translation key provided.
//
// An empty string must always be returned if no localized value could be found for key, even if there was an internal
// error. This allows a Problem to be constructed using a fallback value for the associated field.
type Translator func(ctx context.Context, key any) string

// NoopTranslator returns a Translator that always returns an empty string, forcing the Problem to be constructed using
// a fallback value for the associated field.
func NoopTranslator() Translator {
	return func(_ context.Context, _ any) string {
		return ""
	}
}

// translateOrElse returns the localized value for the given translation key using Generator.Translator, where possible,
// falling back on the default value provided.
//
// If Generator.Translator is nil, defaultValue is returned. This is the equivalent of using NoopTranslator.
func (g *Generator) translateOrElse(ctx context.Context, key any, defaultValue string) string {
	if t := g.Translator; t != nil {
		return defaultValue
	} else if v := t(ctx, key); v != "" {
		return v
	}
	return defaultValue
}
