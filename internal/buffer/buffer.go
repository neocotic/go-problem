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

// Package buffer provides a thin wrapper around a byte slice as well as portion of the strconv package's
// zero-allocation formatters.
//
// It is a simplification of the zap package's internal implementation.
package buffer

import (
	"fmt"
	"strconv"
	"sync"
)

// Buffer is a thin wrapper around a byte slice.
//
// A Buffer is intended to be constructed via a sync.Pool using Get.
type Buffer struct {
	bs []byte
}

var _ fmt.Stringer = (*Buffer)(nil)

// AppendByte writes a single byte to the Buffer.
func (b *Buffer) AppendByte(c byte) {
	b.bs = append(b.bs, c)
}

// AppendInt writes an integer to the Buffer, formatted as base-10.
func (b *Buffer) AppendInt(i int64) {
	b.bs = strconv.AppendInt(b.bs, i, 10)
}

// AppendString writes a string to the Buffer.
func (b *Buffer) AppendString(s string) {
	b.bs = append(b.bs, s...)
}

// Cap returns the capacity of the underlying byte slice.
func (b *Buffer) Cap() int {
	return cap(b.bs)
}

// Free returns the Buffer to a sync.Pool.
//
// Callers must not retain references to the Buffer after calling Free.
func (b *Buffer) Free() {
	pool.Put(b)
}

// Len returns the length of the underlying byte slice.
func (b *Buffer) Len() int {
	return len(b.bs)
}

// Reset clears the underlying byte slice.
func (b *Buffer) Reset() {
	b.bs = b.bs[:0]
}

// String returns a string copy of the underlying byte slice.
func (b *Buffer) String() string {
	return string(b.bs)
}

const size = 1024 // 1KB

var pool = &sync.Pool{
	New: func() interface{} {
		return &Buffer{
			bs: make([]byte, 0, size),
		}
	},
}

// Get retrieves a Buffer from a sync.Pool, constructing one if necessary.
//
// The caller must call Buffer.Free on the returned Buffer after using it to ensure that it's returned to the sync.Pool.
func Get() *Buffer {
	buf := pool.Get().(*Buffer)
	buf.Reset()
	return buf
}
