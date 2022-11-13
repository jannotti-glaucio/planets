package tracer

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go"
)

type Options struct {
	Key   string
	Value interface{}
}

type Span interface {
	opentracing.Span
}

//ITracerProvider for sign token struct
type ITracerProvider interface {
	ContextWithSpan(ctx context.Context, span Span) context.Context
	StartSpanWidthContext(ctx context.Context, identifier string, opts ...Options) Span
	StartSpanFromRequest(request *http.Request, identifier string, opts ...Options) Span
}
