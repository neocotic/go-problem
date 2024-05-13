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
	"errors"
	"fmt"
	"github.com/neocotic/go-optional"
	"github.com/neocotic/go-problem/internal/stack"
	"maps"
	"net/http"
)

// Flag provides control over the generation of specific data and its visibility on their respective fields on a
// Problem.
//
// A Generator will a field corresponding to each support Flag use case that contains the default Flag.
//
// If a Flag-supporting method on a Builder is called or a Flag-supporting Option is used, but no flags are provided,
// this is considered equal to passing FlagField and FlagLog. This would mean that the corresponding data will be
// generated and the fully visible on the Problem both in terms of field and within the logs. If FlagDisable is ever
// passed, all other flags are ignored and the corresponding data is not generated (or inherited) and will not be
// visible on the Problem.
type Flag uint8

// FlagDisable disables generation (or inheritance) of the corresponding data and hides it on the Problem.
const FlagDisable Flag = 0
const (
	// FlagField triggers generation (or inheritance) of the corresponding data and populates the respective exported
	// field on the Problem.
	//
	// Effectively, if the Problem is accessed either directly or deserialized, the data will be accessible, however,
	// FlagField alone will not result in the data being present any time the Problem is logged. FlagLog is required for
	// that.
	FlagField Flag = 1 << iota
	// FlagLog triggers generation (or inheritance) of the corresponding data and populates the respective unexported
	// field on the Problem.
	//
	// Effectively, any time the Problem is logged the data will be present, however, if the Problem is accessed either
	// directly or deserialized the data will be inaccessible unless used in combination with FlagField.
	FlagLog
)

// Builder is used to construct a Problem using methods to define fields and/or override fields derived from a
// Definition and/or Type.
type Builder struct {
	// Generator is the Generator to be used when building a Problem.
	//
	// If Generator is nil, DefaultGenerator will be used.
	Generator *Generator
	// code is the explicitly defined Code to be used. See Builder.Code for more information.
	code Code
	// ctx is the context to be used when building a Problem.
	ctx optional.Optional[context.Context]
	// def is the Definition whose fields are to be treated as defaults when a field is not explicitly defined. See
	// Builder.Definition and Builder.DefinitionType for more information.
	def Definition
	// detail is the explicitly defined detail to be used. See Builder.Detail for more information.
	detail string
	// detailKey is the explicitly defined translation key to be used to resolve a localized detail. See
	// Builder.DetailKey for more information.
	detailKey any
	// err is the explicitly defined error to be wrapped. See Builder.Wrap for more information.
	err error
	// extensions is a shallow clone of the explicitly defined extensions to be used. See Builder.Extension and
	// Builder.Extensions for more information.
	extensions map[string]any
	// instanceURI is the explicitly defined instance URI reference to be used. See Builder.Instance for more
	// information.
	instanceURI string
	// logLevel is the explicitly defined LogLevel to be used. See Builder.LogLevel for more information.
	logLevel LogLevel
	// problem contains any fields unwrapped from err using an Unwrapper. See Builder.Wrap for more information.
	problem Problem
	// stack is the captured stack trace to be used. See Builder.Stack for more information.
	//
	// stack is captured lazily and priority is given to any existing stack contained within problem. getStack must be
	// used to access the stack trace.
	stack string
	// stackFlag contains the stack trace flags to be used. See Builder.Stack for more information.
	stackFlag optional.Optional[Flag]
	// stackFramesSkipped contains the number of additional stack frames to be skipped. See Builder.StackFramesSkipped
	// for more information.
	stackFramesSkipped int
	// status is the explicitly defined status to be used. See Builder.Status for more information.
	status int
	// title is the explicitly defined title to be used. See Builder.Title for more information.
	title string
	// titleKey is the explicitly defined translation key to be used to resolve a localized title. See Builder.TitleKey
	// for more information.
	titleKey any
	// typeURI is the explicitly defined type URI reference to be used. See Builder.Type for more information.
	typeURI string
	// uuid is the generated "UUID" to be used. See Builder.UUID for more information.
	//
	// uuid is generated lazily and priority is given to any existing "UUID" contained within problem. getUUID must be
	// used to access the "UUID".
	uuid string
	// uuidFlag contains the "UUID" flags to be used. See Builder.UUID for more information.
	uuidFlag optional.Optional[Flag]
}

var _ fmt.Stringer = (*Builder)(nil)

var (
	// errExtensionKeyEmpty is returned if an empty extension key is encountered.
	errExtensionKeyEmpty = errors.New("extension key cannot be empty")
	// errExtensionKeyReserved is returned if a reserved extension key is encountered.
	errExtensionKeyReserved = errors.New("extension key is reserved")
)

// Clone returns a clone of the Builder.
func (b *Builder) Clone() *Builder {
	if b == nil {
		return nil
	}
	clone := *b
	// Shallow clone will have to do since extensions could contain any type of values
	clone.extensions = maps.Clone(b.extensions)
	return &clone
}

// Code sets the given Code to be used when building a Problem. See Problem.Code for more information.
//
// If code is not empty, it will take precedence over anything provided using Builder.Definition or Builder.Wrap.
func (b *Builder) Code(code Code) *Builder {
	b.code = code
	return b
}

// Definition sets the given Definition to be used when building a Problem.
//
// The fields of def are treated as defaults when a field is not explicitly defined. This method can conflict with
// Builder.DefinitionType as it effectively assigns to the same underlying field.
func (b *Builder) Definition(def Definition) *Builder {
	b.def = def
	return b
}

// DefinitionType sets the given Type to be used when building a Problem.
//
// The fields of defType are treated as defaults when a field is not explicitly defined. This method can conflict with
// Builder.Definition as it effectively assigns to the same underlying field, however, only setting Definition.Type.
func (b *Builder) DefinitionType(defType Type) *Builder {
	b.def.Type = defType
	return b
}

// Detail sets the given detail to be used when building a Problem. See Problem.Detail for more information.
//
// If detail is not empty, it will take precedence over anything provided using Builder.Definition or Builder.Wrap.
// However, if a localized detail is resolved from a translation key using Builder.DetailKey, that will take precedence
// over detail.
func (b *Builder) Detail(detail string) *Builder {
	b.detail = detail
	return b
}

// Detailf sets the given formatted detail to be used when building a Problem. See Problem.Detail for more information.
//
// If the formatted detail is not empty, it will take precedence over anything provided using Builder.Definition or
// Builder.Wrap. However, if a localized detail is resolved from a translation key using Builder.DetailKey, that will
// take precedence over the formatted detail.
func (b *Builder) Detailf(format string, args ...any) *Builder {
	b.detail = fmt.Sprintf(format, args...)
	return b
}

// DetailKey sets the translation key for to be used to localize the detail when building a Problem. See Problem.Detail
// for more information.
//
// The localized detail will be looked up using Generator.Translator, where possible. If resolved, it will take
// precedence over anything provided using Builder.Detail, Builder.Definition, or Builder.Wrap.
func (b *Builder) DetailKey(key any) *Builder {
	b.detailKey = key
	return b
}

// Extension appends the given extension key and value to that used when building a Problem. See Problem.Extensions for
// more information.
//
// When used, it will take precedence over any extensions provided using Builder.Definition or Builder.Wrap.
//
// Panics if key is either empty or reserved (i.e. conflicts with Problem-level fields).
//
// Builder.Extensions may be preferred for providing multiple extensions and does not conflict with usage of Extension
// in that neither method will delete/modify extensions unless the key overlaps, in which case the value will be
// overwritten.
func (b *Builder) Extension(key string, value any) *Builder {
	if b.extensions == nil {
		b.extensions = make(Extensions)
	}
	if err := validationExtensionKey(key); err != nil {
		panic(err)
	}
	b.extensions[key] = value
	return b
}

// Extensions sets a shallow clone of the given extensions to be used when building a Problem. See Problem.Extensions
// for more information.
//
// If extensions is not empty, it will take precedence over anything provided using Builder.Definition or Builder.Wrap.
//
// Panics if extensions contains a key that is either empty or reserved (i.e. conflicts with Problem-level fields).
//
// Builder.Extension may be preferred for providing a single extension and does not conflict with usage of Extensions in
// that neither method will delete/modify extensions unless the key overlaps, in which case the value will be
// overwritten.
func (b *Builder) Extensions(extensions Extensions) *Builder {
	l := len(extensions)
	if l == 0 {
		b.extensions = nil
		return b
	}
	if b.extensions == nil {
		b.extensions = make(Extensions, l)
	}
	for k, v := range extensions {
		if err := validationExtensionKey(k); err != nil {
			panic(err)
		}
		b.extensions[k] = v
	}
	return b
}

// Instance sets the instance URI reference to be used when building a Problem. See Problem.Instance for more
// information.
//
// An uri.Builder can be used to aid building the URI reference.
//
// If instanceURI is not empty, it will take precedence over anything provided using Builder.Definition or Builder.Wrap.
func (b *Builder) Instance(instanceURI string) *Builder {
	b.instanceURI = instanceURI
	return b
}

// Instancef sets the formatted instance URI reference to be used when building a Problem. See Problem.Instance for more
// information.
//
// If the formatted instance URI reference is not empty, it will take precedence over anything provided using
// Builder.Definition or Builder.Wrap.
func (b *Builder) Instancef(format string, args ...any) *Builder {
	b.instanceURI = fmt.Sprintf(format, args...)
	return b
}

// LogLevel sets the LogLevel to be used when building a Problem. See Problem.LogLevel for more information.
//
// If level is not zero, it will take precedence over anything provided using Builder.Definition,
// Builder.DefinitionType, or Builder.Wrap.
func (b *Builder) LogLevel(level LogLevel) *Builder {
	b.logLevel = level
	return b
}

// Problem returns a constructed Problem.
func (b *Builder) Problem() *Problem {
	return b.build(1)
}

// Reset clears all information used to build a Problem.
func (b *Builder) Reset() *Builder {
	// Retain Generator and ctx
	b.code = ""
	b.def = Definition{}
	b.detail = ""
	b.detailKey = nil
	b.err = nil
	b.extensions = nil
	b.instanceURI = ""
	b.logLevel = 0
	b.problem = Problem{}
	b.stack = ""
	b.stackFlag = optional.Empty[Flag]()
	b.stackFramesSkipped = 0
	b.status = 0
	b.title = ""
	b.titleKey = nil
	b.typeURI = ""
	b.uuid = ""
	b.uuidFlag = optional.Empty[Flag]()
	return b
}

// Stack sets the flags to be used to control if/how a captured stack trace is visible when building a Problem. See
// Problem.Stack for more information.
//
// By default, Generator.StackFlag is used to control visibility of a stack trace.
//
// If no flags are provided, this is considered equal to passing FlagField and FlagLog. If FlagDisable is given, all
// other flags are ignored. No stack trace is captured if FlagDisable is provided.
//
// If a stack trace needs to be captured and Builder.Wrap is used and a Problem is unwrapped that already has a stack
// trace, its stack trace will be used instead of recapturing a stack trace that could potentially result in loss of
// depth in terms of stack frames.
func (b *Builder) Stack(flags ...Flag) *Builder {
	b.stackFlag = resolveFlag(flags)
	return b
}

// StackFramesSkipped sets the number of additional frames to be skipped if/when a stack trace is captured. See
// Builder.Stack and Problem.Stack for more information.
//
// If skipped is less than or equal to zero, no additional frames will be skipped.
func (b *Builder) StackFramesSkipped(skipped int) *Builder {
	b.stackFramesSkipped = skipped
	return b
}

// Status sets the given status to be used when building a Problem. See Problem.Status for more information.
//
// If status is not zero, it will take precedence over anything provided using Builder.Definition,
// Builder.DefinitionType, or Builder.Wrap.
func (b *Builder) Status(status int) *Builder {
	b.status = status
	return b
}

// String constructs a Problem and returns its string representation.
func (b *Builder) String() string {
	return b.build(1).String()
}

// Title sets the given title to be used when building a Problem. See Problem.Title for more information.
//
// If title is not empty, it will take precedence over anything provided using Builder.Definition or Builder.Wrap.
// However, if a localized title is resolved from a translation key using Builder.TitleKey, that will take precedence
// over title.
func (b *Builder) Title(title string) *Builder {
	b.title = title
	return b
}

// Titlef sets the given formatted title to be used when building a Problem. See Problem.Title for more information.
//
// If the formatted title is not empty, it will take precedence over anything provided using Builder.Definition or
// Builder.Wrap. However, if a localized title is resolved from a translation key using Builder.TitleKey, that will take
// precedence over the formatted title.
func (b *Builder) Titlef(format string, args ...any) *Builder {
	b.title = fmt.Sprintf(format, args...)
	return b
}

// TitleKey sets the translation key for to be used to localize the title when building a Problem. See Problem.Title for
// more information.
//
// The localized title will be looked up using Generator.Translator, where possible. If resolved, it will take
// precedence over anything provided using Builder.Title, Builder.Definition, or Builder.Wrap.
func (b *Builder) TitleKey(key any) *Builder {
	b.titleKey = key
	return b
}

// Type sets the type URI reference to be used when building a Problem. See Problem.Type for more information.
//
// An uri.Builder can be used to aid building the URI reference.
//
// If typeURI is not empty, it will take precedence over anything provided using Builder.Definition,
// Builder.DefinitionType, or Builder.Wrap.
func (b *Builder) Type(typeURI string) *Builder {
	b.typeURI = typeURI
	return b
}

// Typef sets the formatted type URI reference to be used when building a Problem. See Problem.Type for more
// information.
//
// If the formatted type URI reference is not empty, it will take precedence over anything provided using
// Builder.Definition, Builder.DefinitionType, or Builder.Wrap.
func (b *Builder) Typef(format string, args ...any) *Builder {
	b.typeURI = fmt.Sprintf(format, args...)
	return b
}

// UUID sets the flags to be used to control if/how a generated "UUID" is visible when building a Problem. See
// Problem.UUID for more information.
//
// By default, Generator.UUIDFlag is used to control visibility of a "UUID".
//
// If no flags are provided, this is considered equal to passing FlagField and FlagLog. If FlagDisable is given, all
// other flags are ignored. No "UUID" is generated if FlagDisable is provided.
//
// If a "UUID" needs to be generated and Builder.Wrap is used and a Problem is unwrapped that already has a "UUID", its
// "UUID" will be used instead of generating a "UUID" to ensure that it is used consistently once generated to make
// tracing any logs, specifically, a lot easier.
func (b *Builder) UUID(flags ...Flag) *Builder {
	b.uuidFlag = resolveFlag(flags)
	return b
}

// Wrap sets the error to be wrapped when building a Problem. See Problem.Error and Problem.Unwrap for more information.
//
// Additionally, more control can be achieved over the scenario where err's tree contains a Problem by passing an
// Unwrapper; a function responsible for deciding what, if any, information from a wrapped Problem is to be used when
// building a Problem. Any such information will not take precedence over any explicitly defined Problem fields,
// however, it will take precedence over any information derived from a Definition or its Type using Builder.Definition
// and/or Builder.DefinitionType.
//
// If no Unwrapper is provided, Generator.Unwrapper is used from Builder.Generator if not nil, otherwise from
// DefaultGenerator. If an Unwrapper could still not be resolved, it defaults to PropagatedFieldUnwrapper.
func (b *Builder) Wrap(err error, unwrapper ...Unwrapper) *Builder {
	var _unwrapper Unwrapper
	if len(unwrapper) > 0 {
		_unwrapper = unwrapper[0]
	} else if g := b.Generator; g != nil {
		_unwrapper = g.Unwrapper
	} else {
		_unwrapper = DefaultGenerator.Unwrapper
	}
	if _unwrapper == nil {
		_unwrapper = unwrapPropagatedFields
	}
	b.err = err
	b.problem = _unwrapper(err)
	return b
}

// build effectively does the heavy lifting for Builder.Problem but allows control over the number of stack frames to be
// skipped, which is useful for other internal calls.
//
// skipStackFrames is the number of frames before recording the stack trace with zero identifying the caller of build.
func (b *Builder) build(skipStackFrames int) *Problem {
	skipStackFrames++
	ctx := b.ctx.OrElseGet(context.Background)
	g := b.Generator
	if g == nil {
		g = GetGenerator(ctx)
	}
	return &Problem{
		Code:       b.buildCode(),
		Detail:     b.buildDetail(ctx, g),
		Extensions: b.buildExtensions(),
		Instance:   b.buildInstance(),
		Stack:      b.buildStack(g, skipStackFrames),
		Status:     b.buildStatus(),
		Title:      b.buildTitle(ctx, g),
		Type:       b.buildType(g),
		UUID:       b.buildUUID(ctx, g),
		err:        b.err,
		logInfo:    b.buildLogInfo(ctx, g, skipStackFrames),
	}
}

// buildCode returns the most suitable Code for building a Problem.
func (b *Builder) buildCode() Code {
	return firstNonZeroValue(b.code, b.problem.Code, b.def.Code)
}

// buildDetail returns the most suitable detail for building a Problem.
func (b *Builder) buildDetail(ctx context.Context, gen *Generator) string {
	var v string
	if v = gen.translateOrElse(ctx, b.detailKey, b.detail); v != "" {
		return v
	}
	if v = b.problem.Detail; v != "" {
		return v
	}
	return gen.translateOrElse(ctx, b.def.DetailKey, b.def.Detail)
}

// buildExtensions returns a shallow clone of the most suitable extensions for building a Problem.
func (b *Builder) buildExtensions() map[string]any {
	return maps.Clone(firstNonNilMap(b.extensions, b.problem.Extensions, b.def.Extensions))
}

// buildInstance returns the most suitable instance URI reference for building a Problem.
func (b *Builder) buildInstance() string {
	return firstNonZeroValue(b.instanceURI, b.problem.Instance, b.def.Instance)
}

// buildLogInfo returns the most suitable log information for building a Problem.
//
// The stack trace or UUID in the returned logInfo will be empty if stackFlag or uuidFlag do not contain FlagLog
// respectively.
//
// skipStackFrames is the number of frames before recording the stack trace with zero identifying the caller of
// buildLogInfo.
func (b *Builder) buildLogInfo(ctx context.Context, gen *Generator, skipStackFrames int) (info LogInfo) {
	info.Level = firstNonZeroValue(b.logLevel, b.problem.logInfo.Level, gen.logLevel(b.def.Type))
	if checkFlag(b.stackFlag.OrElse(gen.StackFlag), FlagLog) {
		info.Stack = b.getStack(skipStackFrames + 1)
	}
	if checkFlag(b.uuidFlag.OrElse(gen.UUIDFlag), FlagLog) {
		info.UUID = b.getUUID(ctx, gen)
	}
	return
}

// buildStack returns the most suitable stack trace for building a Problem.
//
// An empty string is returned if stackFlag does not contain FlagField.
//
// skipStackFrames is the number of frames before recording the stack trace with zero identifying the caller of
// buildStack.
func (b *Builder) buildStack(gen *Generator, skipStackFrames int) string {
	if checkFlag(b.stackFlag.OrElse(gen.StackFlag), FlagField) {
		return b.getStack(skipStackFrames + 1)
	}
	return ""
}

// buildStatus returns the most suitable status for building a Problem. 500 is returned if no suitable status could be
// derived.
func (b *Builder) buildStatus() int {
	return firstNonZeroValue(b.status, b.problem.Status, b.def.Type.Status, http.StatusInternalServerError)
}

// buildTitle returns the most suitable title for building a Problem.
func (b *Builder) buildTitle(ctx context.Context, gen *Generator) string {
	var v string
	if v = gen.translateOrElse(ctx, b.titleKey, b.title); v != "" {
		return v
	}
	if v = b.problem.Title; v != "" {
		return v
	}
	if v = gen.translateOrElse(ctx, b.def.Type.TitleKey, b.def.Type.Title); v != "" {
		return v
	}
	return DefaultTitle
}

// buildType returns the most suitable type URI reference for building a Problem. DefaultTypeURI is returned if no
// suitable type URI reference could be derived.
func (b *Builder) buildType(gen *Generator) string {
	return firstNonZeroValue(b.typeURI, b.problem.Type, gen.typeURI(b.def.Type), DefaultTypeURI)
}

// buildUUID returns the most suitable "UUID" for building a Problem.
//
// An empty string is returned if uuidFlag does not contain FlagField.
func (b *Builder) buildUUID(ctx context.Context, gen *Generator) string {
	if checkFlag(b.uuidFlag.OrElse(gen.UUIDFlag), FlagField) {
		return b.getUUID(ctx, gen)
	}
	return ""
}

// getStack returns a lazily captured stack trace to be used for building a Problem. Priority is given to any existing
// stack contained within problem.
//
// skip is the number of frames before recording the stack trace with zero identifying the caller of getStack.
func (b *Builder) getStack(skip int) string {
	if b.stack != "" {
		return b.stack
	}
	switch {
	case b.problem.Stack != "":
		b.stack = b.problem.Stack
	case b.problem.logInfo.Stack != "":
		b.stack = b.problem.logInfo.Stack
	default:
		if b.stackFramesSkipped > 0 {
			skip += b.stackFramesSkipped
		}
		b.stack = stack.Take(skip + 1)
	}
	return b.stack
}

// getUUID returns a lazily generated "UUID" to be used for building a Problem. Priority is given to any existing uuid
// contained within problem.
func (b *Builder) getUUID(ctx context.Context, gen *Generator) string {
	if b.uuid != "" {
		return b.uuid
	}
	switch {
	case b.problem.UUID != "":
		b.uuid = b.problem.UUID
	case b.problem.logInfo.UUID != "":
		b.uuid = b.problem.logInfo.UUID
	default:
		b.uuid = gen.uuid(ctx)
	}
	return b.uuid
}

// Build returns a Builder for the Generator with context.Background which can be used to construct problems.
func (g *Generator) Build() Builder {
	return Builder{
		Generator: g,
		ctx:       optional.Of(context.Background()),
	}
}

// BuildContext returns a Builder for the Generator with the given context which can be used to construct problems.
func (g *Generator) BuildContext(ctx context.Context) Builder {
	return Builder{
		Generator: g,
		ctx:       optional.Of(ctx),
	}
}

// Build is a convenient shorthand for calling Generator.Build on DefaultGenerator.
func Build() Builder {
	return Builder{
		Generator: DefaultGenerator,
		ctx:       optional.Of(context.Background()),
	}
}

// BuildContext is a convenient shorthand for calling Generator.BuildContext on the Generator within the given
// context.Context, if any, otherwise DefaultGenerator.
func BuildContext(ctx context.Context) Builder {
	return Builder{
		Generator: GetGenerator(ctx),
		ctx:       optional.Of(ctx),
	}
}

// checkFlag returns whether the given Flag contains the other provided.
func checkFlag(flag, other Flag) bool {
	return flag&other == other
}

// firstNonNilMap returns the first non-nil map from those provided.
func firstNonNilMap[K comparable, V any](maps ...map[K]V) map[K]V {
	for _, m := range maps {
		if m != nil {
			return m
		}
	}
	return nil
}

// firstNonZeroValue returns the first non-zero value from those provided.
func firstNonZeroValue[T comparable](values ...T) T {
	var zero T
	for _, v := range values {
		if v != zero {
			return v
		}
	}
	return zero
}

// resolveFlag returns an optional Flag based on the given flags.
//
// If flags is empty, this is considered equal to passing FlagField and FlagLog. If FlagDisable is given, all other
// flags are ignored.
func resolveFlag(flags []Flag) optional.Optional[Flag] {
	var res Flag
	if len(flags) > 0 {
		for _, f := range flags {
			if f == FlagDisable {
				res = f
				break
			} else {
				res |= f
			}
		}
	} else {
		res = FlagField | FlagLog
	}
	return optional.Of(res)
}

// validationExtensionKey returns an error if the extension key provided is either empty or reserved.
func validationExtensionKey(key string) error {
	if key == "" {
		return errExtensionKeyEmpty
	}
	if _, reserved := reservedExtensions[key]; reserved {
		return fmt.Errorf("%w: %q", errExtensionKeyReserved, key)
	}
	return nil
}
