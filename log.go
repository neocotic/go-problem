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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log/slog"
)

type (
	// LogInfo contains information associated with a Problem that is only relevant for logging purposes.
	LogInfo struct {
		// Level is the LogLevel that has either been explicitly defined during construction or inherited from a Type or
		// another Problem within an err's tree if unwrapped accordingly.
		Level LogLevel
		// Stack is a string representation of the stack trace captured during construction or inherited from another
		// Problem within an err's tree if unwrapped accordingly.
		//
		// Stack is only populated if Generator.StackFlag has FlagLog or either Builder.Stack or WithStack were used and
		// either passed no flags or FlagLog explicitly.
		Stack string
		// UUID is the Universally Unique Identifier generated during construction or inherited from another Problem
		// within an err's tree if unwrapped accordingly.
		//
		// UUID is only populated if Generator.UUIDFlag has FlagLog or either Builder.UUID or WithUUID were used and
		// either passed no flags or FlagLog explicitly.
		UUID string
	}

	// LogLeveler is a function that can be used by a Generator to override the LogLevel derived from a Type (i.e.
	// instead of only Type.LogLevel).
	//
	// It is important to note that if the function returns zero, it will fall back to DefaultLogLevel and not
	// Type.LogLevel.
	LogLeveler func(defType Type) LogLevel

	// Logger is a function used by a Generator to log a message and problem and any additional arguments.
	//
	// The Problem is passed within the last two arguments; its key (Generator.LogArgKey) and value. If
	// Generator.LogArgKey is empty, DefaultLogArgKey is passed.
	Logger func(ctx context.Context, level LogLevel, msg string, args ...any)
)

const (
	// DefaultLogArgKey is the default argument key passed to Logger immediately before the Problem at the end of the
	// arguments, and is used by DefaultGenerator.
	DefaultLogArgKey = "problem"

	// DefaultLogLevel is the LogLevel used when one could not be derived.
	DefaultLogLevel = LogLevelError
	// defaultSlogLevel is the slog.Level used when one could not be derived.
	defaultSlogLevel = slog.LevelError
	// defaultZapLevel is the zapcore.Level used when one could not be derived.
	defaultZapLevel = zapcore.ErrorLevel

	// badZapKey is used when a zap.Logger is passed an unexpected arg key.
	badZapKey = "!BADKEY"
)

// Log logs the given message and Problem along with any additional arguments and context.Background via
// Generator.Logger.
//
// The Problem is passed to Generator.Logger within the last two arguments; its key (Generator.LogArgKey) and value. If
// Generator.LogArgKey is empty, DefaultLogArgKey is used.
//
// If Generator.Logger is nil, DefaultLogger is used to log the message.
func (g *Generator) Log(msg string, prob *Problem, args ...any) {
	g.LogContext(context.Background(), msg, prob, args...)
}

// LogContext logs the given message and Problem along with the context provided and any additional arguments via
// Generator.Logger.
//
// The Problem is passed to Generator.Logger within the last two arguments; its key (Generator.LogArgKey) and value. If
// Generator.LogArgKey is empty, DefaultLogArgKey is used.
//
// If Generator.Logger is nil, DefaultLogger is used to log the message.
func (g *Generator) LogContext(ctx context.Context, msg string, prob *Problem, args ...any) {
	lak := g.LogArgKey
	if lak == "" {
		lak = DefaultLogArgKey
	}
	args = append(args, lak, prob)
	fn := g.Logger
	if fn == nil {
		fn = DefaultLogger()
	}
	fn(ctx, prob.logLevel(), msg, args...)
}

// logLevel checks if Generator.LogLeveler is present and, if so, calls it with the given Type to allow for the LogLevel
// to be overridden, where appropriate. Otherwise, Type.LogLevel is returned.
func (g *Generator) logLevel(defType Type) LogLevel {
	if ll := g.LogLeveler; ll != nil {
		return ll(defType)
	}
	return defType.LogLevel
}

// Log is a convenient shorthand for calling Generator.Log on DefaultGenerator.
func Log(msg string, prob *Problem, args ...any) {
	DefaultGenerator.LogContext(context.Background(), msg, prob, args...)
}

// LogContext is a convenient shorthand for calling Generator.LogContext on the Generator within the given
// context.Context, if any, otherwise DefaultGenerator.
func LogContext(ctx context.Context, msg string, prob *Problem, args ...any) {
	GetGenerator(ctx).LogContext(ctx, msg, prob, args...)
}

var _ slog.LogValuer = (*Problem)(nil)

// LogInfo returns information associated with the Problem that is only relevant for logging purposes.
//
// It is mostly intended to be used when using a custom Logger as a means of providing an alternative to
// Problem.LogValue, which is limited to supporting slog.Logger.
func (p *Problem) LogInfo() LogInfo {
	var info LogInfo
	if p != nil {
		info = p.logInfo
	}
	if info.Level == 0 {
		info.Level = DefaultLogLevel
	}
	return info
}

// LogValue returns a slog.GroupValue representation of the Problem containing attrs for only non-empty fields.
func (p *Problem) LogValue() slog.Value {
	attrs := make([]slog.Attr, 0, 10)
	if p.Code != "" {
		attrs = append(attrs, slog.String("code", string(p.Code)))
	}
	if p.Detail != "" {
		attrs = append(attrs, slog.String("detail", p.Detail))
	}
	if p.err != nil {
		attrs = append(attrs, slog.Any("error", p.err))
	}
	if len(p.Extensions) > 0 {
		attrs = append(attrs, mapLogGroup("extensions", p.Extensions))
	}
	if p.Instance != "" {
		attrs = append(attrs, slog.String("instance", p.Instance))
	}
	if p.logInfo.Stack != "" {
		attrs = append(attrs, slog.String("stack", p.logInfo.Stack))
	}
	if p.Status != 0 {
		attrs = append(attrs, slog.Int("status", p.Status))
	}
	if p.Title != "" {
		attrs = append(attrs, slog.String("title", p.Title))
	}
	if p.Type != "" {
		attrs = append(attrs, slog.String("type", p.Type))
	}
	if p.logInfo.UUID != "" {
		attrs = append(attrs, slog.String("uuid", p.logInfo.UUID))
	}
	return slog.GroupValue(attrs...)
}

// MarshalLogObject appends non-empty fields of the Problem to enc.
func (p *Problem) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if p.Code != "" {
		enc.AddString("code", string(p.Code))
	}
	if p.Detail != "" {
		enc.AddString("detail", p.Detail)
	}
	if p.err != nil {
		enc.AddString("error", p.err.Error())
	}
	if len(p.Extensions) > 0 {
		if err := enc.AddReflected("extensions", p.Extensions); err != nil {
			return err
		}
	}
	if p.Instance != "" {
		enc.AddString("instance", p.Instance)
	}
	if p.logInfo.Stack != "" {
		enc.AddString("stack", p.logInfo.Stack)
	}
	if p.Status != 0 {
		enc.AddInt("status", p.Status)
	}
	if p.Title != "" {
		enc.AddString("title", p.Title)
	}
	if p.Type != "" {
		enc.AddString("type", p.Type)
	}
	if p.logInfo.UUID != "" {
		enc.AddString("uuid", p.logInfo.UUID)
	}
	return nil
}

// logLevel returns the LogLevel recommend to be used to log the Problem.
func (p *Problem) logLevel() LogLevel {
	if p == nil || p.logInfo.Level == 0 {
		return DefaultLogLevel
	}
	return p.logInfo.Level
}

// DefaultLogger returns a Logger that uses slog.Default and is used by DefaultGenerator.
func DefaultLogger() Logger {
	return DefaultLoggerContext(func(_ context.Context, logger *slog.Logger) *slog.Logger {
		return logger
	})
}

// DefaultLoggerContext returns a Logger that uses slog.Default while passing the context to the function provided to
// return the most appropriate slog.Logger.
//
// This can be useful for cases where the context is used to further enrich logs.
func DefaultLoggerContext(handleCtx func(ctx context.Context, logger *slog.Logger) *slog.Logger) Logger {
	return func(ctx context.Context, level LogLevel, msg string, args ...any) {
		handleCtx(ctx, slog.Default()).Log(ctx, level.slogLevel(), msg, args...)
	}
}

// GlobalZapLogger returns a Logger that uses zap.L.
func GlobalZapLogger() Logger {
	return ZapLoggerFromContext(zap.L(), func(_ context.Context, logger *zap.Logger) *zap.Logger {
		return logger
	})
}

// GlobalZapLoggerContext returns a Logger that uses zap.L while passing the context to the function provided to return
// the most appropriate zap.Logger.
//
// This can be useful for cases where the context is used to further enrich logs.
func GlobalZapLoggerContext(handleCtx func(ctx context.Context, logger *zap.Logger) *zap.Logger) Logger {
	return ZapLoggerFromContext(zap.L(), handleCtx)
}

// LoggerFrom returns a Logger that uses the given slog.Logger.
func LoggerFrom(logger *slog.Logger) Logger {
	return LoggerFromContext(logger, func(_ context.Context, _ *slog.Logger) *slog.Logger {
		return logger
	})
}

// LoggerFromContext returns a Logger that uses the given slog.Logger while passing the context to the function provided
// to return the most appropriate slog.Logger.
//
// This can be useful for cases where the context is used to further enrich logs.
func LoggerFromContext(logger *slog.Logger, handleCtx func(ctx context.Context, logger *slog.Logger) *slog.Logger) Logger {
	return func(ctx context.Context, level LogLevel, msg string, args ...any) {
		handleCtx(ctx, logger).Log(ctx, level.slogLevel(), msg, args...)
	}
}

// NoopLogger returns a Logger that does nothing.
func NoopLogger() Logger {
	return func(_ context.Context, _ LogLevel, _ string, _ ...any) {
		// Do nothing
	}
}

// ZapLoggerFrom returns a Logger that uses the given zap.Logger.
func ZapLoggerFrom(logger *zap.Logger) Logger {
	return ZapLoggerFromContext(logger, func(_ context.Context, _ *zap.Logger) *zap.Logger {
		return logger
	})
}

// ZapLoggerFromContext returns a Logger that uses the given zap.Logger while passing the context to the function
// provided to return the most appropriate zap.Logger.
//
// This can be useful for cases where the context is used to further enrich logs.
func ZapLoggerFromContext(logger *zap.Logger, handleCtx func(ctx context.Context, logger *zap.Logger) *zap.Logger) Logger {
	return func(ctx context.Context, level LogLevel, msg string, args ...any) {
		handleCtx(ctx, logger).Log(level.zapLevel(), msg, extractZapFields(args)...)
	}
}

// LogLevel represents a log level and will typically need mapped to one understood by any custom Logger.
//
// It has built-in support for slog.Level when DefaultLogger or LoggerFrom are used.
//
// The zero value is intentionally not mapped in order to represent an undefined value and should be substituted by a
// fallback/default LogLevel.
type LogLevel uint

const (
	// LogLevelDebug represents the DEBUG log level.
	LogLevelDebug LogLevel = iota + 1
	// LogLevelInfo represents the INFO log level.
	LogLevelInfo
	// LogLevelWarn represents the WARN log level.
	LogLevelWarn
	// LogLevelError represents the ERROR log level.
	LogLevelError
)

// slogLevel returns the slog.Level representation of the LogLevel, where possible, otherwise defaultSlogLevel.
func (ll LogLevel) slogLevel() slog.Level {
	switch ll {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	default:
		return defaultSlogLevel
	}
}

// zapLevel returns the zapcore.Level representation of the LogLevel, where possible, otherwise defaultZapLevel.
func (ll LogLevel) zapLevel() zapcore.Level {
	switch ll {
	case LogLevelDebug:
		return zapcore.DebugLevel
	case LogLevelInfo:
		return zapcore.InfoLevel
	case LogLevelWarn:
		return zapcore.WarnLevel
	case LogLevelError:
		return zapcore.ErrorLevel
	default:
		return defaultZapLevel
	}
}

// argsToZapField turns a prefix of the non-empty args slice into a zapcore.Field and returns the unconsumed portion of
// the slice.
//
// If args[0] is a zapcore.Field, it returns it.
// If args[0] is a slog.Attr, it returns it and zapcore.Field contain the slog.Attr key-value pair.
// If args[0] is a string, it treats the first two elements as a key-value pair.
// Otherwise, it treats args[0] as a value with a missing key.
func argsToZapField(args []any) (zapcore.Field, []any) {
	switch k := args[0].(type) {
	case string:
		if len(args) == 1 {
			return zap.String(badZapKey, k), nil
		}
		return zap.Any(k, args[1]), args[2:]
	case slog.Attr:
		return zap.Any(k.Key, k.Value.Any()), args[1:]
	case zap.Field:
		return k, args[1:]
	default:
		return zap.Any(badZapKey, k), args[1:]
	}
}

// extractZapFields consumes all args into a slice of zapcore.Field.
func extractZapFields(args []any) (fields []zapcore.Field) {
	var f zapcore.Field
	for len(args) > 0 {
		f, args = argsToZapField(args)
		fields = append(fields, f)
	}
	return
}

// mapLogGroup returns a slog.Attr with a slog.GroupValue containing all entries within the given map.
func mapLogGroup(key string, m map[string]any) slog.Attr {
	var attrs []any
	for k, v := range m {
		attrs = append(attrs, slog.Any(k, v))
	}
	return slog.Group(key, attrs...)
}
