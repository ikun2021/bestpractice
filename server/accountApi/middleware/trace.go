package middleware

import (
	"github.com/gin-gonic/gin"
	"github/lunxun9527/bestpractice/pkg/xtrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	traceName = "jaeger-id"
)

func Trace(conf xtrace.JaegerTraceConf) gin.HandlerFunc {
	tracer := otel.Tracer(traceName)
	propagator := otel.GetTextMapPropagator()
	return func(c *gin.Context) {
		spanName := c.Request.URL.Path
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		spanCtx, span := tracer.Start(
			ctx,
			spanName,
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(
				conf.Name, spanName, c.Request)...),
		)
		defer span.End()
		propagator.Inject(spanCtx, propagation.HeaderCarrier(c.Request.Header))
		c.Request.WithContext(spanCtx)

		c.Next()
		span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(c.Writer.Status())...)
		span.SetStatus(semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(
			c.Writer.Status(), oteltrace.SpanKindServer))
	}
}
