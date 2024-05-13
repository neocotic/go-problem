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
	"github.com/google/uuid"
	"io"
)

// UUIDGenerator is a function used by a Generator to generate a Universally Unique Identifier.
type UUIDGenerator func(ctx context.Context) string

// fallbackUUID is used by each built-in UUIDGenerator when an error occurs while trying to generate a UUID.
const fallbackUUID = "00000000-0000-0000-0000-000000000000"

// V4UUIDGenerator returns a UUIDGenerator that generates a (V4) UUID and is used by DefaultGenerator.
//
// The strength of the generated UUIDs are based on the strength of the crypto/rand package.
//
// Uses the randomness pool if it was enabled with uuid.EnableRandPool.
func V4UUIDGenerator() UUIDGenerator {
	return func(_ context.Context) string {
		return handleUUID(uuid.NewRandom())
	}
}

// V4UUIDGeneratorFromReader returns a UUIDGenerator that generates a (V4) UUID based on bytes read from the given
// reader.
func V4UUIDGeneratorFromReader(reader io.Reader) UUIDGenerator {
	return func(_ context.Context) string {
		return handleUUID(uuid.NewRandomFromReader(reader))
	}
}

// uuid returns a generated UUID using UUIDGenerator, where possible.
//
// If UUIDGenerator is nil, a (V4) UUID is generated and returned.
func (g *Generator) uuid(ctx context.Context) string {
	fn := g.UUIDGenerator
	if fn == nil {
		return handleUUID(uuid.NewRandom())
	}
	return fn(ctx)
}

// handleUUID returns fallbackUUID if err is not nil, otherwise the string representation of the given UUID is returned.
func handleUUID(_uuid uuid.UUID, err error) string {
	if err != nil {
		return fallbackUUID
	} else {
		return _uuid.String()
	}
}
