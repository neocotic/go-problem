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

// Package problem provides support for generating "problem details" in accordance to RFC 9457
// https://datatracker.ietf.org/doc/html/rfc9457, represented as a Problem. A Generator can be used to control a lot of
// the logic applied when generating problems. When not specified, DefaultGenerator, the zero value of Generator, is
// used where appropriate.
//
// While a Problem can be created manually by populating fields, this is not the intended approach will result in many
// of the controls offered by Generator being lost. Instead a Problem is constructed using either Build or New with
// options. Both are the same as the latter uses a Builder under-the-hood, so the only real difference is syntax.
//
// A Problem is an error and, as such, can wrap another error and can be wrapped. Wrap (and Builder.Wrap) can be used
// for wrapped while As and Is can be used for unwrapping.
//
// The package also provides opt-in support for stack trace capturing and UUID generation for problems along with the
// concept of a problem Code.
package problem
