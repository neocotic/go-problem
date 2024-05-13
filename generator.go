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

// Generator is responsible for generating a Problem. Its zero value (DefaultGenerator) is usable.
type Generator struct {
	// CodeNSValidator is the NSValidator used to perform additional validation on a NS used within a Code constructed
	// and/or parsed by a Coder.
	//
	// When nil, a Code may contain any rune except Generator.CodeSeparator within its NS so long as it is not empty.
	// Even when not nil, the NS of a Code must not contain Generator.CodeSeparator and cannot be empty.
	//
	// For example;
	//
	//	g := &Generator{CodeNSValidator: ComposeNSValidator(
	//		LenNSValidator(1, 4),
	//		UnicodeNSValidator(unicode.IsUpper, unicode.ToUpper),
	//	)}
	//	c := g.Coder()
	//	c.Validate("USER-404")   // nil
	//	c.Validate("USERS-404")  // ErrCode
	//	c.Validate("user-404")   // ErrCode
	CodeNSValidator NSValidator
	// CodeSeparator is the rune used to separate the NS and value within a Code constructed and/or parsed by a Coder.
	//
	// If zero or less, DefaultCodeSeparator will be used. Otherwise, it must be a printable rune otherwise a Coder will
	// always return an ErrCode when attempting to construct and/or parse a Code.
	//
	// For example;
	//
	//	g := &Generator{CodeSeparator: '.'}
	//	c := g.Coder("USER")
	//	c.MustBuild(404)  // "USER.404"
	CodeSeparator rune
	// CodeValueLen is the number of digits to be included in the value of a Code constructed and/or parsed by a Coder.
	//
	// If zero or less, a Code may contain any number of digits within its value so long as there's at least one.
	// Otherwise, a value cannot contain more digits than CodeValueLen and any value containing fewer digits will be
	// right-padded with zero.
	//
	// For example;
	//
	//	g := &Generator{CodeValueLen: 8}
	//	c := g.Coder("USER")
	//	c.MustBuild(404)  // "USER.40400000"
	CodeValueLen int
	// ContentType is the value used to populate the Content-Type header when Generator.WriteError or
	// Generator.WriteProblem are called without a WriteOptions.ContentType being passed. This also applies to the
	// Middleware functions as they call Generator.WriteError internally.
	//
	// If empty, ContentTypeJSONUTF8 will be used.
	ContentType string
	// LogArgKey is the key passed along with a Problem within the last two arguments to Generator.Logger.
	//
	// If empty, DefaultLogArgKey will be passed.
	//
	// This allows a somewhat more granular level of control without needing to provide a custom Logger.
	LogArgKey string
	// LogLeveler is the problem.LogLeveler used to override the LogLevel derived from a Type (i.e. instead of only
	// Type.LogLevel).
	//
	// If nil, Type.LogLevel will be used, with a fallback to DefaultLogLevel.
	//
	// For example;
	//
	//	leveler := func(defType Type) LogLevel {
	//		switch {
	//		case defType.Status == 0:
	//			return defType.LogLevel
	//		case def.Type.Status > 500:
	//			return LogLevelError
	//		default:
	//			return LogLevelWarn
	//		}
	//	}
	//	g := &Generator{LogLeveler: leveler}
	LogLeveler LogLeveler
	// Logger is the problem.Logger used by Generator.Log and Generator.LogContext to log a message along with any
	// arguments (incl. the Problem).
	//
	// If nil, DefaultLogger will be used.
	//
	// For example;
	//
	//	logger := slog.NewLogLogger(slog.NewJSONHandler(os.Stderr, nil), slog.LevelDebug)
	//	g := &Generator{Logger: LoggerFrom(logger)}
	Logger Logger
	// StackFlag provides control over the capturing of a stack trace and its visibility on a Problem.
	//
	// StackFlag is the default Flag. If Builder.Stack or WithStack are used, but no flags are provided, this is
	// considered equal to passing FlagField and FlagLog. This would mean that the UUID will be generated and the fully
	// visible on the Problem both in terms of field and within the logs. If FlagDisable is ever passed, all other flags
	// are ignored and the stack trace is not captured (or inherited) and will not be visible on the Problem.
	//
	// For example;
	//
	//	g := &Generator{StackFlag: FlagDisable}          // Stack trace is not captured or inherited
	//	g := &Generator{StackFlag: FlagField}            // Stack trace accessible via Problem.Stack
	//	g := &Generator{StackFlag: FlagLog}              // Stack trace visible only in logs
	//	g := &Generator{StackFlag: FlagField | FlagLog}  // Stack trace accessible via Problem.Stack and visible in logs
	StackFlag Flag
	// Translator is the problem.Translator used to provide localized values for translation keys, where possible, when
	// constructing a Problem.
	//
	// If nil, NoopTranslator will be used, which will always return an empty string, forcing the Problem to be
	// constructed using a fallback value for the associated field.
	//
	// For example;
	//
	//	translator := func(ctx context.Context, key any) string {
	//		t, ok := ctx.Value(contextKeyUT).(ut.Translator)
	//		if !ok {
	//			var err error
	//			if t, err = newUT(ctx); err != nil {
	//				return ""
	//			}
	//		}
	//		if v, err := t.T(key); err != nil {
	//			return ""
	//		} else {
	//			return v
	//		}
	//	}
	//	g := &Generator{Translator: translator}
	Translator Translator
	// Typer is the problem.Typer used to override the type URI reference derived from a Type (i.e. instead of only
	// Type.URI).
	//
	// If nil, Type.URI will be used, with a fallback to DefaultTypeURI.
	//
	// For example;
	//
	//	typer := func(defType Type) string {
	//		if defType.Status == http.StatusNotFound {
	//			return "https://datatracker.ietf.org/doc/html/rfc9110#name-404-not-found"
	//		}
	//		return defType.URI
	//	}
	//	g := &Generator{Typer: typer}
	Typer Typer
	// Unwrapper is the problem.Unwrapper used by Builder.Wrap and Wrap handle an already wrapped Problem in an error's
	// tree.
	//
	// If nil, PropagatedFieldUnwrapper will be used.
	//
	// For example;
	//
	//	g := &Generator{Unwrapper: FullUnwrapper()}
	Unwrapper Unwrapper
	// UUIDFlag provides control over the generation of an Universally Unique Identifier and its visibility on a
	// Problem.
	//
	// UUIDFlag is the default Flag. If Builder.UUID or WithUUID are used, but no flags are provided, this is considered
	// equal to passing FlagField and FlagLog. This would mean that the UUID will be generated and the fully visible on
	// the Problem both in terms of field and within the logs. If FlagDisable is ever passed, all other flags are
	// ignored and the UUID is not generated (or inherited) and will not be visible on the Problem.
	//
	// For example;
	//
	//	g := &Generator{UUIDFlag: FlagDisable}          // UUID not generated or inherited
	//	g := &Generator{UUIDFlag: FlagField}            // UUID accessible via Problem.UUID
	//	g := &Generator{UUIDFlag: FlagLog}              // UUID visible only in logs
	//	g := &Generator{UUIDFlag: FlagField | FlagLog}  // UUID accessible via Problem.UUID and visible in logs
	UUIDFlag Flag
	// UUIDGenerator returns the problem.UUIDGenerator used to generate a Universally Unique Identifier when
	// constructing a Problem.
	//
	// If nil, V4UUIDGenerator will be used to generate a UUID. A UUID is only generated for a Problem if needed. That
	// is; if Builder.UUID or WithUUID are used and not passed FlagDisable or, if not used, Generator.UUIDFlag contains
	// FlagField and/or FlagLog.
	//
	// For example;
	//
	//	nanoidGenerator := func(ng nanoid.generator, err error) UUIDGenerator {
	//		if err != nil {
	//			panic(err)
	//		}
	//		return func(_ context.Context) string {
	//			return ng()
	//		}
	//	}
	//	g := &Generator{UUIDGenerator: nanoidGenerator(nanoid.Canonic())}
	UUIDGenerator UUIDGenerator
}

// DefaultGenerator is the default Generator used when none is given to some top-level functions and structs.
//
// While relatively unopinionated, it is designed to work out-of-the-box with the most commonly desired behaviour having
// the following characteristics:
//
//   - Stack traces are not captured and UUIDs are not generated by default (see Generator.StackFlag and
//     Generator.UUIDFlag respectively for more information)
//   - Any UUID that is generated (e.g. via Builder.UUID or WithUUID) is a (V4) UUID (see Generator.UUIDGenerator for
//     more information)
//   - Any stack trace, UUID, or LogLevel of a Problem found in the tree of an error passed to Builder.Wrap or Wrap is
//     unwrapped and treated as defaults for the generated Problem by default (see Generator.Unwrapper for more
//     information)
//   - Any translation keys are ignored (see Generator.Translator for more information)
//   - Any Code constructed and/or parsed can have any non-empty NS and value and are separated by DefaultCodeSeparator
//     (see Generator.CodeNSValidator, Generator.CodeValueLen, and Generator.CodeSeparator respectively for more
//     information)
//   - Any message that is logged (e.g. via Generator.Log or Generator.LogContext) is done so using slog.Default with
//     DefaultLogArgKey passed as the key along with a Problem within the last two arguments (see Generator.Logger and
//     Generator.LogArgKey respectively for more information)
//   - The LogLevel derived from a Type is always Type.LogLevel (see Generator.LogLeveler for more information)
var DefaultGenerator = &Generator{}
