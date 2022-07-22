package internal

import "google.golang.org/grpc/status"

type GRPCStatus interface {
	GRPCStatus() *status.Status
}
