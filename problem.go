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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/neocotic/go-optional"
	"strconv"
	"strings"
)

// Extensions is a map that may contain additional information used extend the details of a Problem.
type Extensions map[string]any

var (
	_ xml.Marshaler   = (Extensions)(nil)
	_ xml.Unmarshaler = (Extensions)(nil)
)

// MarshalXML marshals the encoded entries within the map into XML.
//
// This is required in order to allow extensions to be marshaled at the top-level of a Problem. Additionally,
// xml.Unmarshaler does not support marshaling maps by default since XML documents are ordered by design. Unfortunately,
// this function cannot guarantee ordering.
//
// An error is returned if unable to marshal any of the entries or the map contains a key that is either empty or
// reserved (i.e. conflicts with Problem-level fields).
func (es Extensions) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	for k, v := range es {
		err := validationExtensionKey(k)
		if err != nil {
			return err
		}
		err = e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: k}})
		if err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalXML does nothing but is required to ensure that any <extensions> element is not unmarshaled as this cannot
// be supported.
func (es Extensions) UnmarshalXML(_ *xml.Decoder, _ xml.StartElement) error {
	// Does nothing
	return nil
}

type (
	// Problem describes problem details for HTTP APIs based on RFC 9457; https://datatracker.ietf.org/doc/html/rfc9457.
	//
	// While a Problem can be explicitly constructed, it's expected that either a Builder or New (with options) is used
	// to construct a Problem for the greatest level of control and for fallback/default fields to be applied as well as
	// support for wrapping errors. Construction is typically driven by a Generator which, unless defined, will be the
	// DefaultGenerator.
	Problem struct {
		// Code is a unique Code that identifies the specific occurrence of the Problem.
		//
		// When present, consumers can use Code to communicate with the owner of the generator to help debug the problem
		// and the client can use it to handle a specific occurrence of the problem without risk of conflicts with other
		// problems using the same Status, Title, and/or Type.
		//
		// It MUST NOT change from occurrence to occurrence of the problem.
		Code Code `json:"code,omitempty" xml:"code,omitempty"`
		// Detail is a human-readable explanation specific to this occurrence of the Problem.
		//
		// When present, ought to focus on helping the client correct the problem, rather than giving debugging
		// information.
		//
		// Consumers SHOULD NOT parse Detail for information; instead Extensions is a more suitable and less error-prone
		// way to obtain such information.
		Detail string `json:"detail,omitempty" xml:"detail,omitempty"`
		// Extensions may contain additional information used extend the details of the Problem.
		//
		// Clients consuming problem details MUST ignore any such extensions that they don't recognize; this allows
		// problem types to evolve and include additional information in the future.
		//
		// If/when the Problem is marshalled to JSON or XML any such extensions are serialized at the top level.
		// However, such information is lost if an attempt is made to unmarshal the serialized data bac into XML due to
		// limitations with xml.Unmarshaler. JSON data can be unmarshaled without any issues. If Extensions contains a
		// key that is empty or reserved (i.e. conflicts with Problem-level fields), an error will occur when attempting
		// to marshal the Problem to JSON or XML.
		Extensions Extensions `json:"-" xml:"extensions,omitempty"`
		// Instance is a URI reference that identifies the specific occurrence of the Problem.
		//
		// It may or may not yield further information if dereferenced.
		//
		// It may be a relative URI; this means that it must be resolved relative to the document's base URI.
		Instance string `json:"instance,omitempty" xml:"instance,omitempty"`
		// Stack is a string representation of the stack trace captured when the Problem generated.
		//
		// When present, Stack can be used to help debug the problem, however, care should be taken as a stack trace
		// will contain information about the internal system architecture and therefore potentially pose a risk to
		// information leakage.
		//
		// Stack is only populated if Generator.StackFlag has FlagField or either Builder.Stack or WithStack were used
		// and either passed no flags or FlagField explicitly. If FlagField is not present but FlagLog is, the Problem
		// will contain a capture stack trace internally for logging within LogValue, however, Stack will be empty. This
		// can be useful for cases where a stack trace is desired for logging only.
		Stack string `json:"stack,omitempty" xml:"stack,omitempty"`
		// Status is the HTTP status code to be generated by the origin server for this occurrence of the Problem.
		//
		// When present, is only advisory; it conveys the HTTP status code used for the convenience of the consumer.
		// Generators MUST use the same status code in the actual HTTP response, to assure that generic HTTP software
		// that does not understand this format still behaves correctly.
		//
		// Consumers can use Status to determine what the original status code used by the generator was, in cases where
		// it has been changed (e.g. by an intermediary or cache), and when message bodies persist without HTTP
		// information. Generic HTTP software will still use the HTTP status code.
		Status int `json:"status" xml:"status"`
		// Title is a short, human-readable summary of the type of the Problem.
		//
		// It SHOULD NOT change from occurrence to occurrence of the problem, except for purposes of localization (e.g.
		// using proactive content negotiation).
		Title string `json:"title" xml:"title"`
		// Type is a URI reference that identifies the type of the Problem.
		//
		// It is encouraged that, when dereferenced, it provides human-readable documentation for the problem type (e.g.
		// using HTML). DefaultTypeURI should be used in all cases where Type is absent.
		//
		// Consumers MUST use Type as the primary identifier for the problem type; Title is advisory and included only
		// for users who are not aware of the semantics of the URI and do not have the ability to discover them (e.g.
		// offline log analysis). Consumers SHOULD NOT automatically dereference the type URI.
		//
		// It may be a relative URI; this means that it must be resolved relative to the document's base URI.
		Type string `json:"type" xml:"type"`
		// UUID is the Universally Unique Identifier generated along with the Problem using Generator.UUIDGenerator.
		//
		// When present, it MUST be unique so that consumers can use UUID to communicate with the owner of the generator
		// to help debug the problem.
		//
		// UUID is only populated if Generator.UUIDFlag has FlagField or either Builder.UUID or WithUUID were used and
		// either passed no flags or FlagField explicitly. If FlagField is not present but FlagLog is, the Problem will
		// contain a generated "UUID" internally for logging within LogValue, however, UUID will be empty. This can be
		// useful for cases where a "UUID" is desired for logging only.
		UUID string `json:"uuid,omitempty" xml:"uuid,omitempty"`
		// err is the error wrapped within the Problem, where applicable.
		err error
		// logInfo contains the relevant logging information for the Problem.
		logInfo LogInfo
	}

	// jsonProblem is used to allow JSON data to be unmarshaled into a Problem struct without having
	// Problem.UnmarshalJSON invoked, resulting in a stack overflow.
	jsonProblem Problem
)

const (
	// DefaultTitle is the title given to a Problem if one was not explicitly specified or could be derived.
	DefaultTitle = "Unknown Error"

	// nilString is returned as a string representation of a nil Problem.
	nilString = "<nil>"
	// xmlDefaultLocalName is used to detect whenever a Problem is being marshaled to XML without an explicit local
	// name so that it can be replaced with a preferred one.
	xmlDefaultLocalName = "Problem"
	// xmlDefaultSpaceName is used to detect whenever a Problem is being marshaled to XML without an explicit space name
	// so that it can be replaced with a preferred one.
	xmlDefaultSpaceName = ""
	// xmlPreferredLocalName is substituted for xmlDefaultLocalName whenever it is detected while a Problem is being
	// marshaled to XML.
	xmlPreferredLocalName = "problem"
	// xmlPreferredSpaceName is substituted for xmlDefaultSpaceName whenever it is detected while a Problem is being
	// marshaled to XML.
	xmlPreferredSpaceName = "urn:ietf:rfc:9457"
)

var (
	_ error            = (*Problem)(nil)
	_ fmt.Stringer     = (*Problem)(nil)
	_ json.Marshaler   = (*Problem)(nil)
	_ json.Unmarshaler = (*Problem)(nil)
	_ xml.Marshaler    = (*Problem)(nil)
)

// reservedExtensions contains extension keys that are reserved. These are typically the names of serialized fields on a
// Problem and are intended to be used to prevent entries within problem.Extensions overwriting top-level Problem
// fields during marshaling.
var reservedExtensions = map[string]struct{}{
	"code":       {},
	"detail":     {},
	"extensions": {},
	"instance":   {},
	"stack":      {},
	"status":     {},
	"title":      {},
	"type":       {},
	"uuid":       {},
}

// Error returns the most suitable error message for the Problem.
//
// If the Problem wraps another error, the message of that error will be included.
func (p *Problem) Error() string {
	return p.buildString(true)
}

// Extension returns the value of the extension with given key within the Problem, if present.
func (p *Problem) Extension(key string) (value any, found bool) {
	if p != nil {
		value, found = p.Extensions[key]
	}
	return
}

// MarshalJSON marshals the Problem into JSON.
//
// This is required in order to allow Problem.Extensions to be marshaled at the top-level of a Problem. Unfortunately,
// this can only be managed by marshaling the problem details twice so is suboptimal in terms of performance.
//
// An error is returned if unable to marshal the Problem or Problem.Extensions contains a key that is either empty or
// reserved (i.e. conflicts with Problem-level fields).
func (p *Problem) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(*p)
	if err != nil {
		return nil, err
	}
	if len(p.Extensions) == 0 {
		return b, nil
	}
	var m map[string]any
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	for k, v := range p.Extensions {
		err = validationExtensionKey(k)
		if err != nil {
			return nil, err
		}
		m[k] = v
	}
	return json.Marshal(m)
}

// MarshalXML marshals the Problem into XML.
//
// This is required in order for greater control of the local and space names on the xml.StartElement when their default
// values are expected. In such cases, it's preferred to use local and space names that match RFC 9457.
//
// An error is returned if unable to marshal the Problem or Problem.Extensions contains a key that is either empty or
// reserved (i.e. conflicts with Problem-level fields).
func (p *Problem) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if start.Name.Local == xmlDefaultLocalName {
		start.Name.Local = xmlPreferredLocalName
	}
	if start.Name.Space == xmlDefaultSpaceName {
		start.Name.Space = xmlPreferredSpaceName
	}
	return e.EncodeElement(*p, start)
}

// String returns a string representation of the Problem.
func (p *Problem) String() string {
	return p.buildString(false)
}

// UnmarshalJSON unmarshals the JSON data provided into the Problem.
//
// This is required in order to unmarshal any superfluous JSON properties at the top-level into Problem.Extensions.
// Unfortunately, this can only be managed by unmarshaling the JSON twice so is suboptimal in terms of performance.
//
// An error is returned if unable to unmarshal data.
func (p *Problem) UnmarshalJSON(data []byte) error {
	var jp jsonProblem
	if err := json.Unmarshal(data, &jp); err != nil {
		return err
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	for k := range m {
		if _, reserved := reservedExtensions[k]; reserved {
			delete(m, k)
		}
	}
	if len(m) > 0 {
		jp.Extensions = m
	}
	*p = Problem(jp)
	return nil
}

// Unwrap returns the error wrapped by the Problem, if any, otherwise returns nil.
func (p *Problem) Unwrap() error {
	if p == nil {
		return nil
	}
	return p.err
}

// buildString returns a string representation of the Problem while providing control over whether any wrapped error is
// included.
func (p *Problem) buildString(inclErr bool) string {
	if p == nil {
		return nilString
	}
	var sb strings.Builder
	sb.WriteString(strconv.FormatInt(int64(p.Status), 10))
	if p.Title != "" {
		sb.WriteRune(' ')
		sb.WriteString(p.Title)
	}
	if p.Detail != "" {
		sb.WriteString(" - ")
		sb.WriteString(p.Detail)
	}
	if p.Code != "" {
		sb.WriteString(" [")
		sb.WriteString(string(p.Code))
		sb.WriteRune(']')
	}
	if inclErr && p.err != nil {
		sb.WriteString(": ")
		sb.WriteString(p.err.Error())
	}
	return sb.String()
}

// New returns a constructed Problem using context.Background, optionally using the options provided as well.
func (g *Generator) New(opts ...Option) *Problem {
	return g.new(context.Background(), opts, 1)
}

// NewContext returns a constructed Problem using the given context, optionally using the options provided as well.
func (g *Generator) NewContext(ctx context.Context, opts ...Option) *Problem {
	return g.new(ctx, opts, 1)
}

// new constructs a Builder for the Generator and applies the given options to it, allowing control over the number of
// stack frames to be skipped, which is useful for other internal calls.
//
// skipStackFrames is the number of frames before recording the stack trace with zero identifying the caller of build.
func (g *Generator) new(ctx context.Context, opts []Option, skipStackFrames int) *Problem {
	b := &Builder{
		Generator: g,
		ctx:       optional.Of(ctx),
	}
	for _, opt := range opts {
		opt(b)
	}
	return b.build(skipStackFrames + 1)
}

// New is a convenient shorthand for calling Generator.New on DefaultGenerator.
func New(opts ...Option) *Problem {
	return DefaultGenerator.new(context.Background(), opts, 1)
}

// NewContext is a convenient shorthand for calling Generator.NewContext on the Generator within the given
// context.Context, if any, otherwise DefaultGenerator.
func NewContext(ctx context.Context, opts ...Option) *Problem {
	return GetGenerator(ctx).new(ctx, opts, 1)
}
