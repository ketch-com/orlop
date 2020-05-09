package orlop

import (
	"context"
	"fmt"
	"testing"
)

func testRunner(ctx context.Context, cfg *EmbeddedConfig) error {
	return fmt.Errorf("%s", "Oops!")
}

func TestRun(t *testing.T) {
	Run("wheelhouse", testRunner, &EmbeddedConfig{})
}
