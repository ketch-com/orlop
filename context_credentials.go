package orlop

import (
	"context"
	"google.golang.org/grpc/metadata"
)

// ContextCredentials provides credentials to the client based on the context
type ContextCredentials struct{}

// GetRequestMetadata returns authorization metadata
func (j ContextCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok && len(md.Get("Authorization")) > 0 {
		return map[string]string{
			"authorization": md.Get("Authorization")[0],
		}, nil
	}

	return nil, nil
}

// RequireTransportSecurity denotes we require transport security
func (j ContextCredentials) RequireTransportSecurity() bool {
	return true
}
