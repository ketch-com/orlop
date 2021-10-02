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
