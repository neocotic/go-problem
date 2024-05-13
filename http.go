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
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
)

// WriteOptions contains options that can be used when writing errors/problems to HTTP responses.
//
// All fields are optional with default behaviour clearly documented.
type WriteOptions struct {
	// ContentType is the content/media type to be used in the HTTP response.
	//
	// The value will be ignored if unsupported or not appropriate for the function called. If empty,
	// Generator.ContentType will be used with a fallback to either ContentTypeJSONUTF8 or a more appropriate
	// content/media type depending on the function called.
	ContentType string
	// LogArgs contains arguments to be passed to Generator.LogContext along with the Problem.
	//
	// If empty, no additional arguments will be passed.
	LogArgs []any
	// LogDisabled is whether the Problem should be logged via Generator.LogContext.
	//
	// By default, the Problem will be logged unless the resolved LogLevel for the Problem is disabled for
	// Generator.Logger.
	LogDisabled bool
	// LogMessage is the message to be passed to Generator.LogContext along with the Problem.
	//
	// If empty, a basic message will be passed.
	LogMessage string
	// Status is the status code to the written to the HTTP response.
	//
	// If less than or equal to zero, Problem.Status will be used with a fallback to http.StatusInternalServerError.
	Status int
}

const (
	// contentTypeHeader is the header representing an HTTP response's content/media type.
	contentTypeHeader = "Content-Type"
	// defaultHTTPLogMessage is the default log message used when writing errors/problems to an HTTP response.
	defaultHTTPLogMessage = "A problem has occurred"
	// defaultHTTPPanicLogMessage is the default log message used when writing an error/problem recovered from a panic
	// to an HTTP response within the Middleware functions.
	defaultHTTPPanicLogMessage = "A panic recovery has occurred"
)

// apply applies the fields from the given WriteOptions, if any and where applicable.
//
// The fields of any WriteOptions found are handled as follows:
//
//   - ContentType is applied if not empty and valid (based on function provided)
//   - LogArgs is applied if not empty
//   - LogDisabled is always applied as only a true value changes anything
//   - LogMessage is applied if not empty
//   - Status is applied if greater than zero
//
// If LogMessage is empty and a non-empty log message is not applied, defaultHTTPLogMessage will be applied.
//
// Panics if ContentType is empty and a non-empty valid content/media type is not applied.
func (wo WriteOptions) apply(opts []WriteOptions, isValidCT func(ct string) bool) WriteOptions {
	if len(opts) > 0 {
		_opts := opts[0]
		if _opts.ContentType != "" && isValidCT(_opts.ContentType) {
			wo.ContentType = _opts.ContentType
		}
		wo.LogDisabled = _opts.LogDisabled
		if len(_opts.LogArgs) > 0 {
			wo.LogArgs = _opts.LogArgs
		}
		if _opts.LogMessage != "" {
			wo.LogMessage = _opts.LogMessage
		}
		if _opts.Status > 0 {
			wo.Status = _opts.Status
		}
	}
	if wo.ContentType == "" {
		// Sanity check - should never happen
		panic(errors.New("missing WriteOptions.ContentType"))
	}
	if wo.LogMessage == "" {
		wo.LogMessage = defaultHTTPLogMessage
	}
	return wo
}

// WriteError writes an HTTP response for a Problem where the Problem is unwrapped from err, where possible, with the
// given function being used to provide a default Problem, relying solely on WriteOptions.ContentType to determine how
// the response is formed, with a graceful fallback to Generator.ContentType and ContentTypeJSONUTF8. WriteOptions can
// also be passed for more granular control.
//
// An error is returned if the Problem fails to be written to w.
func (g *Generator) WriteError(err error, w http.ResponseWriter, req *http.Request, probFunc func(err error) *Problem, opts ...WriteOptions) error {
	prob, isProblem := As(err)
	if !isProblem {
		prob = probFunc(err)
	}
	return g.WriteProblem(prob, w, req, opts...)
}

// WriteErrorJSON writes an HTTP response for a Problem in JSON format where the Problem is unwrapped from err, where
// possible, with the given function being used to provide a default Problem. WriteOptions can also be passed for more
// granular control.
//
// An error is returned if the Problem fails to be written to w.
func (g *Generator) WriteErrorJSON(err error, w http.ResponseWriter, req *http.Request, probFunc func(err error) *Problem, opts ...WriteOptions) error {
	prob, isProblem := As(err)
	if !isProblem {
		prob = probFunc(err)
	}
	return g.WriteProblemJSON(prob, w, req, opts...)
}

// WriteErrorXML writes an HTTP response for a Problem in XML format where the Problem is unwrapped from err, where
// possible, with the given function being used to provide a default Problem. WriteOptions can also be passed for more
// granular control.
//
// An error is returned if the Problem fails to be written to w.
func (g *Generator) WriteErrorXML(err error, w http.ResponseWriter, req *http.Request, probFunc func(err error) *Problem, opts ...WriteOptions) error {
	prob, isProblem := As(err)
	if !isProblem {
		prob = probFunc(err)
	}
	return g.WriteProblemXML(prob, w, req, opts...)
}

// WriteProblem writes an HTTP response for the given Problem, optionally using WriteOptions for more granular control,
// relying solely on WriteOptions.ContentType to determine how the response is formed, with a graceful fallback to
// Generator.ContentType and ContentTypeJSONUTF8.
//
// An error is returned if prob fails to be written to w.
func (g *Generator) WriteProblem(prob *Problem, w http.ResponseWriter, req *http.Request, opts ...WriteOptions) error {
	return g.writeProblem(prob, w, req, WriteOptions{ContentType: g.contentType()}.apply(opts, isValidContentType))
}

// WriteProblemJSON writes an HTTP response for the given Problem in JSON format, optionally using WriteOptions for more
// granular control.
//
// An error is returned if prob fails to be written to w.
func (g *Generator) WriteProblemJSON(prob *Problem, w http.ResponseWriter, req *http.Request, opts ...WriteOptions) error {
	return g.writeProblemJSON(prob, w, req, WriteOptions{ContentType: ContentTypeJSONUTF8}.apply(opts, isValidContentTypeForJSON))
}

// WriteProblemXML writes an HTTP response for the given Problem in XML format, optionally using WriteOptions for more
// granular control.
//
// An error is returned if prob fails to be written to w.
func (g *Generator) WriteProblemXML(prob *Problem, w http.ResponseWriter, req *http.Request, opts ...WriteOptions) error {
	return g.writeProblemXML(prob, w, req, WriteOptions{ContentType: ContentTypeXMLUTF8}.apply(opts, isValidContentTypeForXML))
}

// writeProblem writes an HTTP response for the given Problem using WriteOptions, that are expected to have been
// applied, to determine how the response is formed and whether the Problem is logged.
//
// An error is returned if prob fails to be written to w.
//
// Panics if WriteOptions.ContentType is not recognized.
func (g *Generator) writeProblem(prob *Problem, w http.ResponseWriter, req *http.Request, opts WriteOptions) error {
	switch opts.ContentType {
	case ContentTypeJSON, ContentTypeJSONUTF8:
		return g.writeProblemJSON(prob, w, req, opts)
	case ContentTypeXML, ContentTypeXMLUTF8:
		return g.writeProblemXML(prob, w, req, opts)
	default:
		// Sanity check - should never happen
		panic(fmt.Errorf("unexpected WriteOptions.ContentType applied: %q", opts.ContentType))
	}
}

// writeProblemJSON writes an HTTP response for the given Problem in JSON format using WriteOptions, that are expected
// to have been applied, to determine how the response is formed and whether the Problem is logged.
//
// An error is returned if prob fails to be written to w.
func (g *Generator) writeProblemJSON(prob *Problem, w http.ResponseWriter, req *http.Request, opts WriteOptions) error {
	if !opts.LogDisabled && opts.LogMessage != "" {
		g.LogContext(req.Context(), opts.LogMessage, prob, opts.LogArgs...)
	}

	w.Header().Set(contentTypeHeader, opts.ContentType)
	w.WriteHeader(firstNonZeroValue(opts.Status, prob.Status, http.StatusInternalServerError))

	return json.NewEncoder(w).Encode(prob)
}

// writeProblemXML writes an HTTP response for the given Problem in XML format using WriteOptions, that are expected to
// have been applied, to determine how the response is formed and whether the Problem is logged.
//
// An error is returned if prob fails to be written to w.
func (g *Generator) writeProblemXML(prob *Problem, w http.ResponseWriter, req *http.Request, opts WriteOptions) error {
	if !opts.LogDisabled && opts.LogMessage != "" {
		g.LogContext(req.Context(), opts.LogMessage, prob, opts.LogArgs...)
	}

	w.Header().Set(contentTypeHeader, opts.ContentType)
	w.WriteHeader(firstNonZeroValue(opts.Status, prob.Status, http.StatusInternalServerError))

	return xml.NewEncoder(w).Encode(prob)
}

// Middleware is a convenient shorthand for calling MiddlewareUsing with DefaultGenerator.
func Middleware(probFunc func(err error) *Problem, opts ...WriteOptions) func(http.Handler) http.Handler {
	return MiddlewareUsing(nil, probFunc, opts...)
}

// MiddlewareUsing returns a middleware function that is responsible for populating the HTTP request's context.Context
// with the given Generator (which can be retrieved using GetGenerator) and also provides panic recovery, allowing
// recovered values to be used to form Problem HTTP responses, optionally using WriteOptions for more granular control.
//
// If a value recovered from a panic is not a Problem (which is highly likely), probFunc is called with an error
// representation of that value (if not already an error) to be used to construct a Problem.
func MiddlewareUsing(gen *Generator, probFunc func(err error) *Problem, opts ...WriteOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if gen == nil {
				gen = DefaultGenerator
			}

			req = req.WithContext(UsingGenerator(req.Context(), gen))

			defer func() {
				if r := recover(); r != nil {
					var prob *Problem
					_opts := WriteOptions{
						ContentType: gen.contentType(),
						LogMessage:  defaultHTTPPanicLogMessage,
					}.apply(opts, isValidContentType)
					if err, isErr := r.(error); isErr && err != nil {
						var isProblem bool
						prob, isProblem = As(err)
						if !isProblem {
							prob = probFunc(err)
						}
					} else {
						prob = probFunc(fmt.Errorf("%v", r))
					}
					_ = gen.writeProblem(prob, w, req, _opts)
				}
			}()

			next.ServeHTTP(w, req)
		})
	}
}

// WriteError is a convenient shorthand for calling Generator.WriteError on the Generator within the given HTTP
// request's context.Context, if any, otherwise DefaultGenerator.
func WriteError(err error, w http.ResponseWriter, req *http.Request, fn func(err error) *Problem, opts ...WriteOptions) error {
	return GetGenerator(req.Context()).WriteError(err, w, req, fn, opts...)
}

// WriteErrorJSON is a convenient shorthand for calling Generator.WriteErrorJSON on the Generator within the given HTTP
// request's context.Context, if any, otherwise DefaultGenerator.
func WriteErrorJSON(err error, w http.ResponseWriter, req *http.Request, fn func(err error) *Problem, opts ...WriteOptions) error {
	return GetGenerator(req.Context()).WriteErrorJSON(err, w, req, fn, opts...)
}

// WriteErrorXML is a convenient shorthand for calling Generator.WriteErrorXML on the Generator within the given HTTP
// request's context.Context, if any, otherwise DefaultGenerator.
func WriteErrorXML(err error, w http.ResponseWriter, req *http.Request, fn func(err error) *Problem, opts ...WriteOptions) error {
	return GetGenerator(req.Context()).WriteErrorXML(err, w, req, fn, opts...)
}

// WriteProblem is a convenient shorthand for calling Generator.WriteProblem on the Generator within the given HTTP
// request's context.Context, if any, otherwise DefaultGenerator.
func WriteProblem(prob *Problem, w http.ResponseWriter, req *http.Request, opts ...WriteOptions) error {
	return GetGenerator(req.Context()).WriteProblem(prob, w, req, opts...)
}

// WriteProblemJSON is a convenient shorthand for calling Generator.WriteProblemJSON on the Generator within the given
// HTTP request's context.Context, if any, otherwise DefaultGenerator.
func WriteProblemJSON(prob *Problem, w http.ResponseWriter, req *http.Request, opts ...WriteOptions) error {
	return GetGenerator(req.Context()).WriteProblemJSON(prob, w, req, opts...)
}

// WriteProblemXML is a convenient shorthand for calling Generator.WriteProblemXML on the Generator within the given
// HTTP request's context.Context, if any, otherwise DefaultGenerator.
func WriteProblemXML(prob *Problem, w http.ResponseWriter, req *http.Request, opts ...WriteOptions) error {
	return GetGenerator(req.Context()).WriteProblemXML(prob, w, req, opts...)
}
