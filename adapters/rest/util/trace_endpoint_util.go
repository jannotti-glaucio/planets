package util

import (
	"context"

	"github.com/jannotti-glaucio/planets/core/tools/providers/tracer"
	"github.com/labstack/echo/v4"
)

func TraceRestEndpoint(echoCtx echo.Context, identificator string) (tracer.Span, context.Context) {
	jeagerTracer := tracer.New("rest")
	span := jeagerTracer.StartSpanFromRequest(echoCtx.Request(), identificator)
	ctx := jeagerTracer.ContextWithSpan(echoCtx.Request().Context(), span)
	return span, ctx
}
