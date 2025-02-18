package tracex

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	Traceparent = "traceparent"
	Trace       = "trace"
)

var (
	once   sync.Once
	tracer trace.Tracer
)

// InitTracer init tracer
func InitTracer(traceName string) {
	once.Do(func() {
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))

		otel.SetTracerProvider(sdkTrace.NewTracerProvider())

		tracer = otel.Tracer(traceName)
	})
}

// getTracer returns the singleton instance of the tracer
func getTracer() trace.Tracer {
	InitTracer("world-news")

	return tracer
}

// StartSpan start span
func StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	ctx, span := getTracer().Start(ctx, spanName)

	return ctx, span
}

// InjectTraceparentToResponse inject traceparent to response header
func InjectTraceparentToResponse(ctx context.Context, w http.ResponseWriter) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.HeaderCarrier(w.Header()))
}

// ExtractTraceFromRequest extract trace from request header
func ExtractTraceFromRequest(ctx context.Context, r *http.Request) context.Context {
	propagator := otel.GetTextMapPropagator()

	ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

	// if traceparent is not valid, inject a new traceparent
	if !trace.SpanFromContext(ctx).SpanContext().IsValid() {
		ctx = InjectTraceInContext(ctx)
	}

	return ctx
}

// ExtractTraceparent extract traceparent from context
func ExtractTraceparent(ctx context.Context) string {
	spanCtx := trace.SpanFromContext(InjectTraceInContext(ctx)).SpanContext()

	return fmt.Sprintf("00-%s-%s-%s", spanCtx.TraceID(), spanCtx.SpanID(), spanCtx.TraceFlags())
}

const startTrace = "start-trace"

type startTimeKey struct{}

// InjectTraceInContext ensure trace in context
func InjectTraceInContext(ctx context.Context) context.Context {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		return ctx
	}

	ctx, span = StartSpan(ctx, startTrace)

	defer span.End()

	return context.WithValue(ctx, startTimeKey{}, time.Now())
}

// GetTraceID get trace id
func GetTraceID(ctx context.Context) string {
	spanCtx := trace.SpanFromContext(InjectTraceInContext(ctx)).SpanContext()

	return spanCtx.TraceID().String()
}

// LogField log field
func LogField(ctx context.Context) zap.Field {
	return zap.String(Trace, GetTraceID(ctx))
}

// CalculateDuration calculate duration
func CalculateDuration(ctx context.Context) int64 {
	startTime, ok := ctx.Value(startTimeKey{}).(time.Time)
	if !ok {
		return 0
	}

	return time.Since(startTime).Milliseconds()
}
