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
	"cmp"
	"errors"
	"fmt"
)

type (
	// Matcher is a function used to conditionally match on a Problem, returning true only if the match is successful.
	//
	// A Matcher is never passed a nil pointer to a Problem.
	Matcher func(p *Problem) bool

	// Operator is used by a Matcher to compare two values of the same type.
	Operator uint8

	// Unwrapper is a function used by Builder.Wrap and Wrap to handle an already wrapped Problem in err's tree.
	//
	// An Unwrapper is effectively responsible for deciding what, if any, information from a wrapped Problem is to be
	// used to construct the new Problem. Any such information will not take precedence over any explicitly defined
	// Problem fields, however, it will take precedence over any information derived from a Definition or its Type.
	Unwrapper func(err error) Problem
)

const (
	// OperatorEquals is used to check if two values of the same type are equal.
	OperatorEquals Operator = iota
	// OperatorNotEquals is used to check if two values of the same type are not equal.
	OperatorNotEquals
	// OperatorGreaterThan is used to check if one value is greater than another value of the same type.
	OperatorGreaterThan
	// OperatorGreaterThanOrEqual is used to check if one value is greater than or equal to another value of the same
	// type.
	OperatorGreaterThanOrEqual
	// OperatorLessThan is used to check if one value is less than another value of the same type.
	OperatorLessThan
	// OperatorLessThanOrEqual is used to check if one value is less than or equal to another value of the same type.
	OperatorLessThanOrEqual
)

// As is a convenient shorthand for calling errors.As with a Problem target, however, it also gracefully handles the
// case where err is nil without a panic.
func As(err error) (*Problem, bool) {
	if err == nil {
		return nil, false
	}
	var p *Problem
	isProblem := errors.As(err, &p)
	return p, isProblem
}

// AsOrElse is a convenient shorthand for calling errors.As with a Problem target, however, it also gracefully handles
// the case where err is nil without a panic.
//
// If no Problem is found in err's tree, defaultProb is returned but false is also returned to indicate that the default
// was used. This can be useful for providing a default Problem if one is not found in err's tree. For example; error
// middleware may use AsOrElse with a catch-all Problem for cases where an error's tree has no Problem to include in the
// response.
func AsOrElse(err error, defaultProb *Problem) (*Problem, bool) {
	if err == nil {
		return defaultProb, false
	}
	var p *Problem
	isProblem := errors.As(err, &p)
	if !isProblem {
		p = defaultProb
	}
	return p, isProblem
}

// AsOrElseGet is a convenient shorthand for calling errors.As with a Problem target, however, it also gracefully
// handles the case where err is nil without a panic.
//
// If no Problem is found in err's tree, defaultProbFunc is used to return a Problem but false is also returned to
// indicate that the default was used. This can be useful for lazily building/constructing a default Problem if one is
// not found in err's tree. For example; error middleware may use AsOrElseGet to lazily construct a catch-all Problem
// for cases where an error's tree has no Problem to include in the response.
func AsOrElseGet(err error, defaultProbFunc func() *Problem) (*Problem, bool) {
	if err == nil {
		return defaultProbFunc(), false
	}
	var p *Problem
	isProblem := errors.As(err, &p)
	if !isProblem {
		p = defaultProbFunc()
	}
	return p, isProblem
}

// AsMatch is a convenient shorthand for calling errors.As with a Problem target, however, it also gracefully handles
// the case where err is nil without a panic.
//
// Additionally, if a Problem is found in err's tree, it must match all matchers provided, otherwise it will be
// unwrapped, and it's tree (excluding itself) will continue to be checked until either a matching Problem is found or
// no Problem is found.
func AsMatch(err error, matchers ...Matcher) (*Problem, bool) {
	if err == nil {
		return nil, false
	}
	var p *Problem
	if !errors.As(err, &p) {
		return nil, false
	}
	if Match(p, matchers...) {
		return p, true
	}
	if p == nil {
		return nil, false
	}
	return AsMatch(p.Unwrap(), matchers...)
}

// AsMatchOrElse is a convenient shorthand for calling errors.As with a Problem target, however, it also gracefully
// handles the case where err is nil without a panic.
//
// Additionally, if a Problem is found in err's tree, it must match all matchers provided, otherwise it will be
// unwrapped, and it's tree (excluding itself) will continue to be checked until either a matching Problem is found or
// no Problem is found.
//
// If no matching Problem is found in err's tree, defaultProb is returned but false is also returned to indicate that
// the default was used. This can be useful for providing a default Problem if one is not found in err's tree. For
// example; error middleware may use AsMatchOrElse with a catch-all Problem for cases where an error's tree has no
// Problem to include in the response.
func AsMatchOrElse(err error, defaultProb *Problem, matchers ...Matcher) *Problem {
	if p, isMatch := AsMatch(err, matchers...); isMatch {
		return p
	}
	return defaultProb
}

// AsMatchOrElseGet is a convenient shorthand for calling errors.As with a Problem target, however, it also gracefully
// handles the case where err is nil without a panic.
//
// Additionally, if a Problem is found in err's tree, it must match all matchers provided, otherwise it will be
// unwrapped, and it's tree (excluding itself) will continue to be checked until either a matching Problem is found or
// no Problem is found.
//
// If no match Problem is found in err's tree, defaultProbFunc is used to return a Problem but false is also returned to
// indicate that the default was used. This can be useful for lazily building/constructing a default Problem if one is
// not found in err's tree. For example; error middleware may use AsMatchOrElseGet to lazily construct a catch-all
// Problem for cases where an error's tree has no Problem to include in the response.
func AsMatchOrElseGet(err error, defaultProbFunc func() *Problem, matchers ...Matcher) *Problem {
	if p, isMatch := AsMatch(err, matchers...); isMatch {
		return p
	}
	return defaultProbFunc()
}

// Is acts as a substitute for errors.Is, returning true if err's tree contains a Problem.
//
// It is effectively a convenient shorthand for calling As where only the boolean return value is returned.
func Is(err error) bool {
	_, isProblem := As(err)
	return isProblem
}

// IsMatch acts as a substitute for errors.Is, returning true if err's tree contains a Problem that matches all matchers
// provided.
//
// It is effectively a convenient shorthand for calling AsMatch where only the boolean return value is returned.
func IsMatch(err error, matchers ...Matcher) bool {
	_, isMatch := AsMatch(err, matchers...)
	return isMatch
}

// HasCode is used to match a Problem based on its Code.
//
// By default, this match is based on whether the values are equal, however, this can be controlled by passing another
// Operator.
func HasCode(code Code, operator ...Operator) Matcher {
	op := operatorOrDefault(operator)
	return func(p *Problem) bool {
		return operate(op, p.Code, code)
	}
}

// HasCodeNS is used to match a Problem based on the NS within its Code using DefaultGenerator.
//
// By default, this match is based on whether the values are equal, however, this can be controlled by passing another
// Operator.
func HasCodeNS(ns NS, operator ...Operator) Matcher {
	return HasCodeNSUsing(DefaultGenerator, ns, operator...)
}

// HasCodeNSUsing is used to match a Problem based on the NS within its Code using the given Generator.
//
// By default, this match is based on whether the values are equal, however, this can be controlled by passing another
// Operator.
func HasCodeNSUsing(gen *Generator, ns NS, operator ...Operator) Matcher {
	c := gen.Coder()
	op := operatorOrDefault(operator)
	return func(p *Problem) bool {
		parsed, err := c.Parse(p.Code)
		return err != nil && operate(op, parsed.NS, ns)
	}
}

// HasCodeValue is used to match a Problem based on the value within its Code using DefaultGenerator.
//
// By default, this match is based on whether the values are equal, however, this can be controlled by passing another
// Operator.
func HasCodeValue(value uint, operator ...Operator) Matcher {
	return HasCodeValueUsing(DefaultGenerator, value, operator...)
}

// HasCodeValueUsing is used to match a Problem based on the value within its Code using the given Generator.
//
// By default, this match is based on whether the values are equal, however, this can be controlled by passing another
// Operator.
func HasCodeValueUsing(gen *Generator, value uint, operator ...Operator) Matcher {
	c := gen.Coder()
	op := operatorOrDefault(operator)
	return func(p *Problem) bool {
		parsed, err := c.Parse(p.Code)
		return err != nil && operate(op, parsed.Value, value)
	}
}

// HasDetail is used to match a Problem based on its detail.
//
// By default, this match is based on whether the values are equal, however, this can be controlled by passing another
// Operator.
func HasDetail(detail string, operator ...Operator) Matcher {
	op := operatorOrDefault(operator)
	return func(p *Problem) bool {
		return operate(op, p.Detail, detail)
	}
}

// HasExtension is used to match a Problem based on whether it contains an extension with the given key.
func HasExtension(key string) Matcher {
	return func(p *Problem) bool {
		_, found := p.Extension(key)
		return found
	}
}

// HasExtensions is used to match a Problem based on whether it contains extensions with the given keys.
func HasExtensions(keys ...string) Matcher {
	return func(p *Problem) bool {
		for _, key := range keys {
			if _, found := p.Extension(key); !found {
				return false
			}
		}
		return true
	}
}

// HasInstance is used to match a Problem based on its instance.
//
// By default, this match is based on whether the values are equal, however, this can be controlled by passing another
// Operator.
func HasInstance(instance string, operator ...Operator) Matcher {
	op := operatorOrDefault(operator)
	return func(p *Problem) bool {
		return operate(op, p.Instance, instance)
	}
}

// HasStack is used to match a Problem based on whether it has a captured stack trace.
func HasStack() Matcher {
	return func(p *Problem) bool {
		return p.Stack != ""
	}
}

// HasStatus is used to match a Problem based on its status.
//
// By default, this match is based on whether the values are equal, however, this can be controlled by passing another
// Operator.
func HasStatus(status int, operator ...Operator) Matcher {
	op := operatorOrDefault(operator)
	return func(p *Problem) bool {
		return operate(op, p.Status, status)
	}
}

// HasTitle is used to match a Problem based on its title.
//
// By default, this match is based on whether the values are equal, however, this can be controlled by passing another
// Operator.
func HasTitle(title string, operator ...Operator) Matcher {
	op := operatorOrDefault(operator)
	return func(p *Problem) bool {
		return operate(op, p.Title, title)
	}
}

// HasType is used to match a Problem based on its type URI.
//
// By default, this match is based on whether the values are equal, however, this can be controlled by passing another
// Operator.
func HasType(typeURI string, operator ...Operator) Matcher {
	op := operatorOrDefault(operator)
	return func(p *Problem) bool {
		return operate(op, p.Type, typeURI)
	}
}

// HasUUID is used to match a Problem based on whether it has a generated UUID.
func HasUUID() Matcher {
	return func(p *Problem) bool {
		return p.UUID != ""
	}
}

// Match returns whether the given Problem matchers all the matchers provided.
//
// If one or more Matcher is provided but prob is nil, false will always be returned as a Matcher assumes prob is not
// nil.
func Match(prob *Problem, matchers ...Matcher) bool {
	if len(matchers) > 0 {
		if prob == nil {
			return false
		}
		for _, m := range matchers {
			if !m(prob) {
				return false
			}
		}
	}
	return true
}

// Or is used to match a Problem on any of the given matchers.
func Or(matchers ...Matcher) Matcher {
	return func(p *Problem) bool {
		for _, m := range matchers {
			if m(p) {
				return true
			}
		}
		return false
	}
}

// FullUnwrapper returns an Unwrapper that extracts all fields from a wrapped Problem in err's tree, if present. These
// fields will not take precedence over any explicitly defined Problem fields, however, it will take precedence over any
// fields derived from a Definition or its Type.
func FullUnwrapper() Unwrapper {
	return unwrapAllFields
}

// NoopUnwrapper returns an Unwrapper that does nothing.
func NoopUnwrapper() Unwrapper {
	return func(_ error) Problem {
		return Problem{}
	}
}

// PropagatedFieldUnwrapper returns an Unwrapper that extracts only fields that are expected to be propagated (e.g.
// captured stack trace, generated "UUID") from a wrapped Problem in err's tree, if present. Any such fields will not
// take precedence over any explicitly defined Problem fields, however, it will take precedence over any fields derived
// from a Definition or its Type.
func PropagatedFieldUnwrapper() Unwrapper {
	return unwrapPropagatedFields
}

// operate returns the result of the given operation.
//
// Panics if op is invalid.
func operate[T cmp.Ordered](op Operator, probValue, otherValue T) bool {
	c := cmp.Compare(probValue, otherValue)
	switch op {
	case OperatorEquals:
		return c == 0
	case OperatorNotEquals:
		return c != 0
	case OperatorGreaterThan:
		return c == 1
	case OperatorGreaterThanOrEqual:
		return c == 1 || c == 0
	case OperatorLessThan:
		return c == -1
	case OperatorLessThanOrEqual:
		return c == -1 || c == 0
	default:
		// Should never happen
		panic(fmt.Errorf("unsupported Operator: %v", op))
	}
}

// operatorOrDefault returns the first Operator if any are given or OperatorEquals if none are given.
func operatorOrDefault(op []Operator) Operator {
	if len(op) > 0 {
		return op[0]
	}
	return OperatorEquals
}

// unwrapAllFields extracts all fields from a wrapped Problem in err's tree, if present. These fields will not take
// precedence over any explicitly defined Problem fields, however, it will take precedence over any fields derived from
// a Definition or its Type.
func unwrapAllFields(err error) Problem {
	if p, isProblem := As(err); isProblem && p != nil {
		return *p
	}
	return Problem{}
}

// unwrapPropagatedFields extracts only fields that are expected to be propagated (e.g. captured stack trace, generated
// "UUID") from a wrapped Problem in err's tree, if present. Any such fields will not take precedence over any
// explicitly defined Problem fields, however, it will take precedence over any fields derived from a Definition or its
// Type.
func unwrapPropagatedFields(err error) Problem {
	if p, isProblem := As(err); isProblem && p != nil {
		return Problem{
			Stack:   p.Stack,
			UUID:    p.UUID,
			logInfo: p.logInfo,
		}
	}
	return Problem{}
}
