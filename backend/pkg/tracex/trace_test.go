package tracex

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestExtractTraceparent testing ExtractTraceparent
func TestExtractTraceparent(t *testing.T) {
	InitTracer("test-trace")

	traceparent := ExtractTraceparent(context.Background())

	fmt.Println(traceparent)
}

// TestCalculateDuration testing CalculateDuration
func TestCalculateDuration(t *testing.T) {
	InitTracer("test-trace")

	ctx := InjectTraceInContext(context.Background())

	time.Sleep(time.Second)

	duration := CalculateDuration(ctx)

	fmt.Println(duration)
}
