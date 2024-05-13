// Copyright (C) 2024 neocotic
// Copyright (c) 2023 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package stack

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func Benchmark_Take(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Take(0)
	}
}

func Test_Take(t *testing.T) {
	trace := Take(0)
	lines := strings.Split(trace, "\n")
	require.NotEmpty(t, lines, "expected stacktrace to have at least one frame")
	assert.Contains(
		t,
		lines[0],
		"github.com/neocotic/go-problem/internal/stack.Test_Take",
		"expected stacktrace to start with the test",
	)
}

func Test_Take_WithSkip(t *testing.T) {
	trace := Take(1)
	lines := strings.Split(trace, "\n")
	require.NotEmpty(t, lines, "expected stacktrace to have at least one frame")
	assert.Contains(
		t,
		lines[0],
		"testing.",
		"expected stacktrace to start with the test runner (skipping our own frame)",
	)
}

func Test_Take_WithSkipInnerFunc(t *testing.T) {
	var trace string
	func() {
		trace = Take(2)
	}()
	lines := strings.Split(trace, "\n")
	require.NotEmpty(t, lines, "expected stacktrace to have at least one frame")
	assert.Contains(
		t,
		lines[0],
		"testing.",
		"expected stacktrace to start with the test function (skipping the test function)",
	)
}

func Test_Take_DeepStack(t *testing.T) {
	const (
		N                  = 500
		withStackDepthName = "github.com/neocotic/go-problem/internal/stack.withStackDepth"
	)
	withStackDepth(N, func() {
		trace := Take(0)
		for found := 0; found < N; found++ {
			i := strings.Index(trace, withStackDepthName)
			if i < 0 {
				t.Fatalf(`expected %v occurrences of %q, found %d`,
					N, withStackDepthName, found)
			}
			trace = trace[i+len(withStackDepthName):]
		}
	})
}

func withStackDepth(depth int, f func()) {
	var recurse func(rune) rune
	recurse = func(r rune) rune {
		if r > 0 {
			bytes.Map(recurse, []byte(string([]rune{r - 1})))
		} else {
			f()
		}
		return 0
	}
	recurse(rune(depth))
}
