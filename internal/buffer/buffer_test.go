// Copyright (C) 2024 neocotic
// Copyright (c) 2016 Uber Technologies, Inc.
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

package buffer

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Benchmark_Buffer(b *testing.B) {
	str := strings.Repeat("a", 1024)
	slice := make([]byte, 1024)
	buf := bytes.NewBuffer(slice)
	internal := Get()
	b.Run("ByteSlice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice = append(slice, str...)
			slice = slice[:0]
		}
	})
	b.Run("BytesBuffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf.WriteString(str)
			buf.Reset()
		}
	})
	b.Run("InternalBuffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			internal.AppendString(str)
			internal.Reset()
		}
	})
}

func Test_Buffer_Writes(t *testing.T) {
	buf := Get()
	testCases := map[string]struct {
		fn     func()
		expect string
	}{
		"AppendByte":        {func() { buf.AppendByte('v') }, "v"},
		"AppendString":      {func() { buf.AppendString("foo") }, "foo"},
		"AppendIntPositive": {func() { buf.AppendInt(42) }, "42"},
		"AppendIntNegative": {func() { buf.AppendInt(-42) }, "-42"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			buf.Reset()
			tc.fn()
			assert.Equal(t, tc.expect, buf.String(), "unexpected buffer.String()")
			assert.Equal(t, len(tc.expect), buf.Len(), "unexpected buffer length")
			assert.Equal(t, size, buf.Cap(), "expected buffer capacity to remain constant")
		})
	}
}
