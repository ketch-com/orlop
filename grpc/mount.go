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
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.ketch.com/lib/orlop/service"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

type MountGrpcServerParams struct {
	fx.In

	Name   service.Name
	Server *grpc.Server
}

func MountGrpcServer(params MountGrpcServerParams) fx.Annotated {
	return fx.Annotated{
		Name: "routes",
		Target: func(mux chi.Router) {
			grpcHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
					params.Server.ServeHTTP(w, r)
				} else {
					http.NotFound(w, r)
				}
			})

			mux.Handle(fmt.Sprintf("/%s.{service}/*", params.Name), grpcHandler)
		},
	}
}
