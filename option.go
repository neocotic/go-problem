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

// Option is used to customize the generation of a Problem and/or to override fields derived from a Definition and/or
// Type.
//
// An Option provides an alternative syntax for constructing problems using an underlying Builder as each Option simply
// maps to Builder method.
type Option func(b *Builder)

// FromDefinition customizes a Generator to return a Problem using the given Definition when building a Problem.
//
// The fields of def are treated as defaults when a field is not explicitly defined using another option. FromDefinition
// can conflict with FromType as it effectively assigns to the same underlying field.
func FromDefinition(def Definition) Option {
	return func(b *Builder) {
		b.Definition(def)
	}
}

// FromType customizes a Generator to return a Problem using the given Type when building a Problem.
//
// The fields of defType are treated as defaults when a field is not explicitly defined using another option. FromType
// can conflict with FromDefinition as it effectively assigns to the same underlying field, however, only setting
// Definition.Type.
func FromType(defType Type) Option {
	return func(b *Builder) {
		b.DefinitionType(defType)
	}
}

// WithCode customizes a Generator to return a Problem with the given Code. See Problem.Code for more information.
//
// If code is not empty, it will take precedence over anything provided using FromDefinition or any of the Wrap options.
func WithCode(code Code) Option {
	return func(b *Builder) {
		b.Code(code)
	}
}

// WithDetail customizes a Generator to return a Problem with the given detail. See Problem.Detail for more information.
//
// If detail is not empty, it will take precedence over anything provided using FromDefinition or any of the Wrap
// options. However, if a localized detail is resolved from a translation key using WithDetailKey, that will take
// precedence over detail.
func WithDetail(detail string) Option {
	return func(b *Builder) {
		b.Detail(detail)
	}
}

// WithDetailf customizes a Generator to return a Problem with the formatted detail. See Problem.Detail for more
// information.
//
// If the formatted detail is not empty, it will take precedence over anything provided using FromDefinition or any of
// the Wrap options. However, if a localized detail is resolved from a translation key using WithDetailKey, that will
// take precedence over the formatted detail.
func WithDetailf(format string, args ...any) Option {
	return func(b *Builder) {
		b.Detailf(format, args...)
	}
}

// WithDetailKey customizes a Generator to return a Problem with detail localized using the given translation key. See
// Problem.Detail for more information.
//
// The localized detail will be looked up using Generator.Translator, where possible. If resolved, it will take
// precedence over anything provided using WithDetail, FromDefinition, or any of the Wrap options.
func WithDetailKey(key any) Option {
	return func(b *Builder) {
		b.DetailKey(key)
	}
}

// WithDetailKeyOrElse is a convenient shorthand for using both WithDetailKey and WithDetail.
func WithDetailKeyOrElse(key any, detail string) Option {
	return func(b *Builder) {
		b.DetailKey(key).Detail(detail)
	}
}

// WithExtension customizes a Generator to return a Problem containing the given extension key and value. See
// Problem.Extensions for more information.
//
// When used, it will take precedence over any extensions provided using FromDefinition or any of the Wrap options.
//
// Panics if key is either empty or reserved (i.e. conflicts with Problem-level fields).
//
// WithExtensions may be preferred for providing multiple extensions and does not conflict with usage of WithExtension
// in that neither option will delete/modify extensions unless the key overlaps, in which case the value will be
// overwritten.
func WithExtension(key string, value any) Option {
	return func(b *Builder) {
		b.Extension(key, value)
	}
}

// WithExtensions customizes a Generator to return a Problem with a shallow clone of the given extensions. See
// Problem.Extensions for more information.
//
// If extensions is not empty, it will take precedence over anything provided using FromDefinition or any of the Wrap
// options.
//
// Panics if extensions contains a key that is either empty or reserved (i.e. conflicts with Problem-level fields).
//
// WithExtension may be preferred for providing a single extension and does not conflict with usage of WithExtensions in
// that neither option will delete/modify extensions unless the key overlaps, in which case the value will be
// overwritten.
func WithExtensions(extensions Extensions) Option {
	return func(b *Builder) {
		b.Extensions(extensions)
	}
}

// WithInstance customizes a Generator to return a Problem with the given instance URI reference. See Problem.Instance
// for more information.
//
// An uri.Builder can be used to aid building the URI reference.
//
// If instanceURI is not empty, it will take precedence over anything provided using FromDefinition or any of the Wrap
// options.
func WithInstance(instanceURI string) Option {
	return func(b *Builder) {
		b.Instance(instanceURI)
	}
}

// WithInstancef customizes a Generator to return a Problem with the formatted instance URI reference. See
// Problem.Instance for more information.
//
// If the formatted instance URI reference is not empty, it will take precedence over anything provided using
// FromDefinition or any of the Wrap options.
func WithInstancef(format string, args ...any) Option {
	return func(b *Builder) {
		b.Instancef(format, args...)
	}
}

// WithLogLevel customizes a Generator to return a Problem with the given LogLevel. See Problem.LogLevel for more
// information.
//
// If level is not zero, it will take precedence over anything provided using FromDefinition, FromType, or any of the
// Wrap options.
func WithLogLevel(level LogLevel) Option {
	return func(b *Builder) {
		b.LogLevel(level)
	}
}

// WithStack customizes a Generator to control if/how a captured stack trace is visible on a Problem. See Problem.Stack
// for more information.
//
// By default, Generator.StackFlag is used to control visibility of a stack trace.
//
// If no flags are provided, this is considered equal to passing FlagField and FlagLog. If FlagDisable is given, all
// other flags are ignored. No stack trace is captured if FlagDisable is provided.
//
// If a stack trace needs to be captured and any of the Wrap options are used and a Problem is unwrapped that already
// has a stack trace, its stack trace will be used instead of recapturing a stack trace that could potentially result in
// loss of depth in terms of stack frames.
func WithStack(flags ...Flag) Option {
	return func(b *Builder) {
		b.Stack(flags...)
	}
}

// WithStackFramesSkipped customizes a Generator to skip the given number of additional frames if/when a stack trace is
// captured. See WithStack and Problem.Stack for more information.
//
// If skipped is less than or equal to zero, no additional frames will be skipped.
func WithStackFramesSkipped(skipped int) Option {
	return func(b *Builder) {
		b.StackFramesSkipped(skipped)
	}
}

// WithStatus customizes a Generator to return a Problem with the given status. See Problem.Status for more information.
//
// If status is not zero, it will take precedence over anything provided using FromDefinition, FromType, or any of the
// Wrap options.
func WithStatus(status int) Option {
	return func(b *Builder) {
		b.Status(status)
	}
}

// WithTitle customizes a Generator to return a Problem with the given title. See Problem.Title for more information.
//
// If title is not empty, it will take precedence over anything provided using FromDefinition, FromType, or any of the
// Wrap options. However, if a localized title is resolved from a translation key using WithTitleKey, that will take
// precedence over title.
func WithTitle(title string) Option {
	return func(b *Builder) {
		b.Title(title)
	}
}

// WithTitlef customizes a Generator to return a Problem with the formatted title. See Problem.Title for more
// information.
//
// If the formatted title is not empty, it will take precedence over anything provided using FromDefinition or any of
// the Wrap options. However, if a localized title is resolved from a translation key using WithTitleKey, that will take
// precedence over the formatted title.
func WithTitlef(format string, args ...any) Option {
	return func(b *Builder) {
		b.Titlef(format, args...)
	}
}

// WithTitleKey customizes a Generator to return a Problem with title localized using the given translation key. See
// Problem.Title for more information.
//
// The localized title will be looked up using Generator.Translator, where possible. If resolved, it will take
// precedence over anything provided using WithTitle, FromDefinition, FromType, or any of the Wrap options.
func WithTitleKey(key any) Option {
	return func(b *Builder) {
		b.TitleKey(key)
	}
}

// WithTitleKeyOrElse is a convenient shorthand for using both WithTitleKey and WithTitle.
func WithTitleKeyOrElse(key any, title string) Option {
	return func(b *Builder) {
		b.TitleKey(key).Title(title)
	}
}

// WithType customizes a Generator to return a Problem with the given type URI reference. See Problem.Type for more
// information.
//
// An uri.Builder can be used to aid building the URI reference.
//
// If typeURI is not empty, it will take precedence over anything provided using FromDefinition, FromType, or any of the
// Wrap options.
func WithType(typeURI string) Option {
	return func(b *Builder) {
		b.Type(typeURI)
	}
}

// WithTypef customizes a Generator to return a Problem with the formatted type URI reference. See Problem.Type for more
// information.
//
// If the formatted type URI reference is not empty, it will take precedence over anything provided using
// FromDefinition, FromType, or any of the Wrap options.
func WithTypef(format string, args ...any) Option {
	return func(b *Builder) {
		b.Typef(format, args...)
	}
}

// WithUUID customizes a Generator to control if/how a generated UUID is visible on a Problem. See Problem.UUID for more
// information.
//
// By default, Generator.UUIDFlag is used to control visibility of a UUID.
//
// If no flags are provided, this is considered equal to passing FlagField and FlagLog. If FlagDisable is given, all
// other flags are ignored. No UUID is generated if FlagDisable is provided.
//
// If a UUID needs to be generated and any of the Wrap options are used and a Problem is unwrapped that already has a
// UUID, its UUID will be used instead of generating a UUID to ensure that it is used consistently once generated to
// make tracing any logs, specifically, a lot easier.
func WithUUID(flags ...Flag) Option {
	return func(b *Builder) {
		b.UUID(flags...)
	}
}

// Wrap customizes a Generator to return a Problem wrapping the given error. See Problem.Error and Problem.Unwrap for
// more information.
//
// Additionally, more control can be achieved over the scenario where err's tree contains a Problem by passing an
// Unwrapper; a function responsible for deciding what, if any, information from a wrapped Problem is to be used when
// building a Problem. Any such information will not take precedence over any explicitly defined Problem fields,
// however, it will take precedence over any information derived from a Definition or its Type using FromDefinition
// and/or FromType.
//
// If no Unwrapper is provided, Generator.Unwrapper is used from Builder.Generator if not nil, otherwise from
// DefaultGenerator. If an Unwrapper could still not be resolved, it defaults to PropagatedFieldUnwrapper.
func Wrap(err error, unwrapper ...Unwrapper) Option {
	return func(b *Builder) {
		b.Wrap(err, unwrapper...)
	}
}
