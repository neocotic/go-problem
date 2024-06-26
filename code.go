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
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type (
	// Code is intended to be a unique string that identifies a specific occurrence of a Problem.
	//
	// Consumers can use a Code to communicate with the owner of the generator to help debug a problem and the client
	// can use it to handle a specific occurrence of a problem without risk of conflicts with other problems using the
	// same status, title, and/or type URI reference.
	//
	// A well-formed Code is expected to be formed of a namespace (NS) and a numeric value, separated using
	// Generator.CodeSeparator.
	//
	// While Code is not strictly associated with problems, it is common for such a concept to be used in combination
	// and so is provided as an optional extra.
	//
	// Since Code is just a string, there's nothing to prevent a Code being explicitly declared, for example;
	//
	//	code := Code("AUTH-400")
	//
	// However, this is discouraged in favour of using Coder to construct consistent codes that conform to validation
	// defined on a Generator. That said; it really just depends on how important Code consistency is to the generator.
	Code string

	// Coder is used to construct and/or parse a Code.
	//
	// Since certain fields of Generator are used by a Coder to construct/parse a Code, the same Generator should be
	// used for both operations to ensure consistency and validity of each Code.
	Coder struct {
		// Generator is the Generator to be used when building/parsing a Code.
		//
		// If Generator is nil, DefaultGenerator will be used.
		Generator *Generator
		// NS is the namespace to be used when building/parsing a Code. It is only required when building a Code but,
		// when present when parsing a Code, it also validates that the parsed Code was constructed using the same NS.
		NS NS
	}

	// NS represents a namespace within the problem generator (e.g. web application). It can be used to distinguish
	// problems generated by different logical locations within the code base. This allows generic and possibly
	// conflicting Code values to be reused across multiple namespaces while not inhibiting debugging.
	//
	// For example; a problem generated within an authentication service may use a NS of "AUTH" while an user service
	// may use a NS of "USER".
	NS string

	// NSValidator is a function used validate a given NS when a Coder is constructing or parsing a Code.
	//
	// Such validation can be useful to ensure each NS used meets a standard.
	//
	// The error returned should not be an ErrCode as a Coder will wrap it in an ErrCode.
	NSValidator func(ns NS) error

	// ParsedCode contains any information parsed from a Code, including the original Code.
	ParsedCode struct {
		// Code is the Code which was parsed.
		Code Code
		// NS is the namespace found within the parsed Code.
		NS NS
		// Value is the value found within the parsed Code.
		Value uint
	}
)

// DefaultCodeSeparator is the default rune used to separate the NS and value of a Code and is used by DefaultGenerator.
const DefaultCodeSeparator rune = '-'

// ErrCode is returned when a Code cannot be constructed or parsed.
var ErrCode = errors.New("invalid problem code")

// Build returns a constructed Code using the given value, if able.
//
// Coder.NS is required as it's also used during the construction and is separated from value using
// Generator.CodeSeparator.
//
// An ErrCode is returned only in the following cases:
//   - Generator.CodeSeparator is a non-printable rune
//   - Coder.ValidateNS rejects Coder.NS
//   - Coder.ValidateValue rejects value
func (c Coder) Build(value uint) (Code, error) {
	g := c.Generator
	if g == nil {
		g = DefaultGenerator
	}

	sep, err := g.codeSeparator()
	if err != nil {
		return "", err
	}

	suffix := strconv.FormatUint(uint64(value), 10)
	if err = g.validateCodeValue(suffix); err != nil {
		return "", err
	}
	if vl := g.CodeValueLen; vl > 0 {
		for len(suffix) < vl {
			suffix += "0"
		}
	}

	if err = g.validateCodeNS(c.NS, sep); err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString(string(c.NS))
	sb.WriteRune(sep)
	sb.WriteString(suffix)
	return Code(sb.String()), nil
}

// MustBuild is a convenient shorthand for calling Coder.Build that panics if it returns an error.
func (c Coder) MustBuild(value uint) Code {
	if code, err := c.Build(value); err != nil {
		panic(err)
	} else {
		return code
	}
}

// MustParse is a convenient shorthand for calling Coder.Parse that panics if it returns an error.
func (c Coder) MustParse(code Code) ParsedCode {
	if parsed, err := c.Parse(code); err != nil {
		panic(err)
	} else {
		return parsed
	}
}

// MustValidate is a convenient shorthand for calling Coder.Validate that panics if it returns an error.
func (c Coder) MustValidate(code Code) {
	if err := c.Validate(code); err != nil {
		panic(err)
	}
}

// Parse parses the given Code, if able.
//
// Coder.NS is optional but, where present, will result in additional validation on the parsed NS, asserting that they
// are equal. Generator.CodeSeparator is used during parsing to separate the NS and value of the Code.
//
// The returned ParsedCode will contain as much information parsed for code as possible, even though it may be invalid.
// Only when the returned error is nil can it be assumed that ParsedCode contains only valid information.
//
// An ErrCode is returned only in the following cases:
//   - Generator.CodeSeparator is a non-printable rune
//   - Coder.ValidateNS rejects the parsed NS
//   - Coder.ValidateValue rejects the parsed value, or it cannot be parsed as an uint
func (c Coder) Parse(code Code) (ParsedCode, error) {
	g := c.Generator
	if g == nil {
		g = DefaultGenerator
	}

	pc := ParsedCode{Code: code}
	sep, err := g.codeSeparator()
	if err != nil {
		return pc, err
	}

	var (
		inVal    bool
		nsb, vsb strings.Builder
	)

	for _, r := range code {
		if inVal {
			vsb.WriteRune(r)
		} else if r == sep {
			inVal = true
		} else {
			nsb.WriteRune(r)
		}
	}

	if !inVal {
		return pc, fmt.Errorf("%w: Generator.CodeSeparator %q not found: %q", ErrCode, sep, code)
	}

	pc.NS = NS(nsb.String())
	if err = g.validateCodeNS(pc.NS, sep); err != nil {
		return pc, err
	}

	if c.NS != "" && c.NS != pc.NS {
		return pc, fmt.Errorf("%w: NS parsed is unexpected (want %q, got %q): %q", ErrCode, c.NS, pc.NS, code)
	}

	valStr := vsb.String()
	if err = g.validateCodeValue(valStr); err != nil {
		return pc, err
	}

	val, err := strconv.ParseUint(valStr, 10, 0)
	pc.Value = uint(val)
	if err != nil {
		return pc, fmt.Errorf("%w: value cannot be parsed: %q: %w", ErrCode, code, err)
	}

	return pc, nil
}

// Validate validates the given Code and returns an ErrCode if invalid.
//
// It is effectively a convenient shorthand for calling Coder.Parse where only the error is returned.
//
// An ErrCode is returned only in the following cases:
//   - Generator.CodeSeparator is a non-printable rune
//   - Coder.ValidateNS rejects the parsed NS
//   - Coder.ValidateValue rejects the parsed value, or it cannot be parsed as an uint
func (c Coder) Validate(code Code) error {
	_, err := c.Parse(code)
	return err
}

// ValidateNS validates the given NS and returns an ErrCode if invalid.
//
// An ErrCode is returned only in the following cases:
//   - Generator.CodeSeparator is a non-printable rune
//   - ns is empty
//   - ns contains Generator.CodeSeparator
//   - Generator.CodeNSValidator rejects ns, if not nil
func (c Coder) ValidateNS(ns NS) error {
	g := c.Generator
	if g == nil {
		g = DefaultGenerator
	}
	sep, err := g.codeSeparator()
	if err != nil {
		return err
	}
	return g.validateCodeNS(ns, sep)
}

// ValidateValue validates the given value and returns an ErrCode if invalid.
//
// An ErrCode is returned only if number of digits in the string representation of is greater than
// Generator.CodeValueLen, if greater than zero.
func (c Coder) ValidateValue(value uint) error {
	g := c.Generator
	if g == nil {
		g = DefaultGenerator
	}
	s := strconv.FormatUint(uint64(value), 10)
	return g.validateCodeValue(s)
}

// BuildCode is a convenient shorthand for calling Coder.Build on a Coder using DefaultGenerator and optionally a given
// NS.
func BuildCode(value uint, ns ...NS) (Code, error) {
	return DefaultGenerator.Coder(ns...).Build(value)
}

// MustBuildCode is a convenient shorthand for calling Coder.MustBuild on a Coder using DefaultGenerator and optionally
// a given NS.
func MustBuildCode(value uint, ns ...NS) Code {
	return DefaultGenerator.Coder(ns...).MustBuild(value)
}

// MustParseCode is a convenient shorthand for calling Coder.MustParse on a Coder using DefaultGenerator and optionally
// a given NS.
func MustParseCode(code Code, ns ...NS) ParsedCode {
	return DefaultGenerator.Coder(ns...).MustParse(code)
}

// MustValidateCode is a convenient shorthand for calling Coder.MustValidate on a Coder using DefaultGenerator and
// optionally a given NS.
func MustValidateCode(code Code, ns ...NS) {
	DefaultGenerator.Coder(ns...).MustValidate(code)
}

// ParseCode is a convenient shorthand for calling Coder.Parse on a Coder using DefaultGenerator and optionally a given
// NS.
func ParseCode(code Code, ns ...NS) (ParsedCode, error) {
	return DefaultGenerator.Coder(ns...).Parse(code)
}

// ValidateCode is a convenient shorthand for calling Coder.Validate on a Coder using DefaultGenerator and optionally a
// given NS.
func ValidateCode(code Code, ns ...NS) error {
	return DefaultGenerator.Coder(ns...).Validate(code)
}

// ComposeNSValidator returns a NSValidator composed of each of the given validators.
//
// For example;
//
//	ComposeNSValidator(LenNSValidator(4, 8), UnicodeNSValidator(unicode.IsUpper, unicode.ToUpper))
func ComposeNSValidator(validators ...NSValidator) NSValidator {
	return func(ns NS) error {
		for _, validator := range validators {
			if err := validator(ns); err != nil {
				return err
			}
		}
		return nil
	}
}

// LenNSValidator returns a NSValidator that asserts that a NS contains at least the minimum and, optionally, at most
// the maximum number of characters. Otherwise, an error is returned.
//
// Since a Coder validates that a NS is not empty by default, min must be at least one and, if max is provided, must be
// greater than or equal to min.
//
// For example;
//
//	LenNSValidator(4)
//	LenNSValidator(4, 8)
func LenNSValidator(min int, max ...int) NSValidator {
	var _max int
	if len(max) > 0 {
		_max = max[0]
	} else {
		_max = math.MaxInt
	}
	return func(ns NS) error {
		if min < 1 {
			return fmt.Errorf("LenNSValidator min is less than min (want 1, got %v)", min)
		} else if _max < min {
			return fmt.Errorf("LenNSValidator max is less than min (want %v, got %v)", min, _max)
		} else if l := len(ns); l < min {
			return fmt.Errorf("NS contains too few characters (want %v, got %v): %q", min, l, ns)
		} else if l > _max {
			return fmt.Errorf("NS contains too many characters (want %v, got %v): %q", _max, l, ns)
		}
		return nil
	}
}

// RegexpNSValidator returns a NSValidator that asserts that a NS matches the given regular expression. Otherwise, an
// error is returned.
//
// If expr fails to compile into a regexp.Regexp, an error is always returned.
//
// For example;
//
//	RegexpNSValidator(`[a-zA-Z]+`)
func RegexpNSValidator(expr string) NSValidator {
	r, err := regexp.Compile(expr)
	return func(ns NS) error {
		if err != nil {
			return fmt.Errorf("RegexpNSValidator expr could not be compiled: %q: %w", expr, err)
		} else if !r.MatchString(string(ns)) {
			return fmt.Errorf("NS does not match regexp (want %q): %q", expr, ns)
		}
		return nil
	}
}

// UnicodeNSValidator returns a NSValidator that asserts that a NS contains only unicode runes that meet the given
// predicate. Otherwise, an error is returned containing the desired rune using mapper.
//
// For example;
//
//	UnicodeNSValidator(unicode.IsLower, unicode.ToLower)
//	UnicodeNSValidator(unicode.IsUpper, unicode.ToUpper)
func UnicodeNSValidator(predicate func(r rune) bool, mapper func(r rune) rune) NSValidator {
	return func(ns NS) error {
		for i, r := range ns {
			if !predicate(r) {
				return fmt.Errorf("NS contains invalid character at index %v (want %q, got %q): %q", i, mapper(r), r, ns)
			}
		}
		return nil
	}
}

// Coder returns a Coder for the Generator, optionally with an NS.
//
// An NS is only required when using Coder to build a Code but, when present when parsing a Code, it also validates that
// the parsed Code was constructed using the same NS.
func (g *Generator) Coder(ns ...NS) Coder {
	var _ns NS
	if len(ns) > 0 {
		_ns = ns[0]
	}
	return Coder{
		Generator: g,
		NS:        _ns,
	}
}

// codeSeparator returns the rune to be used to separate the NS and value of a Code.
//
// If Generator.CodeSeparator is less than or equal to zero, DefaultCodeSeparator is returned, If
// Generator.CodeSeparator is a printable rune, it is returned. Otherwise, an ErrCode is returned.
func (g *Generator) codeSeparator() (rune, error) {
	if sep := g.CodeSeparator; sep <= 0 {
		return DefaultCodeSeparator, nil
	} else if unicode.IsPrint(sep) {
		return sep, nil
	} else {
		return sep, fmt.Errorf("%w: Generator.CodeSeparator is not printable: %q", ErrCode, sep)
	}
}

// validateCodeNS validates the given NS and returns an ErrCode if invalid.
func (g *Generator) validateCodeNS(ns NS, sep rune) error {
	if ns == "" {
		return fmt.Errorf("%w: NS is empty", ErrCode)
	}
	if strings.ContainsRune(string(ns), sep) {
		return fmt.Errorf("%w: NS contains Generator.CodeSeparator: %q", ErrCode, ns)
	}
	if v := g.CodeNSValidator; v != nil {
		if err := v(ns); err != nil {
			return fmt.Errorf("%w: %w", ErrCode, err)
		}
	}
	return nil
}

// validateCodeValue validates the given string representation of a value and returns an ErrCode if invalid.
func (g *Generator) validateCodeValue(value string) error {
	if value == "" {
		return fmt.Errorf("%w: value is empty", ErrCode)
	}
	if vl := g.CodeValueLen; vl > 0 {
		if l := len(value); l > vl {
			return fmt.Errorf("%w: value contains too many characters (want %v, got %v): %q", ErrCode, vl, l, value)
		}
	}
	return nil
}
