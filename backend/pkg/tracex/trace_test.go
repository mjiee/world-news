package tracex

import (
	"context"
	"fmt"
	"testing"
)

// TestExtractTraceparent testing ExtractTraceparent
func TestExtractTraceparent(t *testing.T) {
	InitTracer("test-trace")

	traceparent := ExtractTraceparent(context.Background())

	fmt.Println(traceparent)
}
