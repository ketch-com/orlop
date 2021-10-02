// Copyright (c) 2021 Ketch Kloud, Inc.
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

package grpc

import (
	"context"
	"google.golang.org/grpc/metadata"
)

// SharedContextCredentials provides context-based or token-based credentials to the client
type SharedContextCredentials struct {
	tokenProvider func(ctx context.Context) string
}

// GetRequestMetadata returns authorization metadata
func (j SharedContextCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	token := j.tokenProvider(ctx)

	if md, ok := metadata.FromIncomingContext(ctx); ok && len(md.Get("Authorization")) > 0 {
		token = md.Get("Authorization")[0]
	}

	return map[string]string{
		"authorization": token,
	}, nil
}

// RequireTransportSecurity denotes we require transport security
func (j SharedContextCredentials) RequireTransportSecurity() bool {
	return true
}
