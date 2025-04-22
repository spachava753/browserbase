package browserbase

import (
	"context"
)

// EnvAPIKeySource implements SecuritySource for ogen and reads the API key from memory.
type EnvAPIKeySource struct {
	APIKey string
}

var _ SecuritySource = (*EnvAPIKeySource)(nil)

// XBBAPIKey fulfills the ogen SecuritySource for the X-BB-API-Key header.
func (s *EnvAPIKeySource) XBBAPIKey(ctx context.Context, operation OperationName) (XBBAPIKey, error) {
	return XBBAPIKey{APIKey: s.APIKey}, nil
}
