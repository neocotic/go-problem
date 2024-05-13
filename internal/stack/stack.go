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

// Package stack provides support for capturing stack traces efficiently.
//
// It is a simplification of the zap package's internal implementation.
package stack

import (
	"github.com/neocotic/go-problem/internal/buffer"
	"runtime"
	"sync"
)

// Stack is a captured stack trace.
//
// A Stack is intended to be constructed via a sync.Pool using Capture.
type Stack struct {
	frames  *runtime.Frames
	pcs     []uintptr
	storage []uintptr
}

// Free releases resources associated with the Stack and returns it to a sync.Pool.
//
// Callers must not retain references to the Stack after calling Free.
func (s *Stack) Free() {
	s.frames = nil
	s.frames = nil
	pool.Put(s)
}

// Len returns the number of frames in the Stack, however, it never changes if Next is called.
func (s *Stack) Len() int {
	return len(s.pcs)
}

// Next returns the next frame in the Stack and a boolean indication of whether there are even more.
func (s *Stack) Next() (frame runtime.Frame, more bool) {
	return s.frames.Next()
}

var pool = &sync.Pool{
	New: func() interface{} {
		return &Stack{
			storage: make([]uintptr, 64),
		}
	},
}

// Capture captures a stack trace, skipping the given number of frames.
//
// When skip is zero, this will identify the caller of Capture.
//
// The caller must call Stack.Free on the returned Stack after using it to ensure that it's resources are released, and
// it's returned to a sync.Pool.
func Capture(skip int) *Stack {
	stack := pool.Get().(*Stack)
	stack.pcs = stack.storage

	fc := runtime.Callers(skip+2, stack.pcs)
	pcs := stack.pcs
	for fc == len(pcs) {
		pcs = make([]uintptr, len(pcs)*2)
		fc = runtime.Callers(skip+2, pcs)
	}

	stack.storage = pcs
	stack.pcs = pcs[:fc]
	stack.frames = runtime.CallersFrames(stack.pcs)
	return stack
}

// Take captures the current stack trace and returns its string representation.
//
// skip is the number of frames before recording the stack trace with zero identifying the caller of Take.
func Take(skip int) string {
	stack := Capture(skip + 1)
	defer stack.Free()

	buf := buffer.Get()
	defer buf.Free()

	fmt := NewFormatter(buf)
	fmt.FormatStack(stack)
	return buf.String()
}

// Formatter is responsible for formatting a stack trace into a readable string representation.
type Formatter struct {
	buf      *buffer.Buffer
	nonEmpty bool
}

// NewFormatter returns a new Formatter.
func NewFormatter(buf *buffer.Buffer) Formatter {
	return Formatter{buf: buf}
}

// FormatFrame formats the given frame, appending it to the buffer.
func (f *Formatter) FormatFrame(frame runtime.Frame) {
	if f.nonEmpty {
		f.buf.AppendByte('\n')
	}
	f.nonEmpty = true
	f.buf.AppendString(frame.Function)
	f.buf.AppendByte('\n')
	f.buf.AppendByte('\t')
	f.buf.AppendString(frame.File)
	f.buf.AppendByte(':')
	f.buf.AppendInt(int64(frame.Line))
}

// FormatStack formats all remaining frames in the given Stack (excl. final runtime.main/runtime.goexit frame),
// appending them to the buffer.
func (f *Formatter) FormatStack(stack *Stack) {
	for frame, more := stack.Next(); more; frame, more = stack.Next() {
		f.FormatFrame(frame)
	}
}
