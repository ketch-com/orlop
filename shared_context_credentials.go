package orlop

import (
	"context"
	"google.golang.org/grpc/metadata"
)

// SharedContextCredentials provides context-based or token-based credentials to the client
type SharedContextCredentials struct {
	token string
}

// GetRequestMetadata returns authorization metadata
func (j SharedContextCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	token := j.token

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
