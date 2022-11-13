package tracer

import (
	"context"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

//JaegerTracer struct for assign interface
type JaegerTracer struct {
	tracer opentracing.Tracer
	Error  error
	Closer io.Closer
}

var singletonJaegerTracer map[string]*JaegerTracer

//New create a new instance
func New(serviceName string) *JaegerTracer {
	if singletonJaegerTracer[serviceName] != nil {
		return singletonJaegerTracer[serviceName]
	}

	ReporterHost := viper.GetString("jaeger.host")
	ReporterPort := viper.GetString("jaeger.port")

	configuration := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{

			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: ReporterHost + ":" + ReporterPort,
		},
	}
	tracer, closer, err := configuration.NewTracer(
		config.Logger(jaeger.StdLogger),
	)
	opentracing.SetGlobalTracer(tracer)

	return &JaegerTracer{
		tracer: tracer,
		Error:  err,
		Closer: closer,
	}
}

//Token implements the TokenInterface
var _ ITracerProvider = &JaegerTracer{}

func (jaeger *JaegerTracer) ContextWithSpan(ctx context.Context, span Span) context.Context {
	return opentracing.ContextWithSpan(ctx, span)
}

func (jaeger *JaegerTracer) StartSpanWidthContext(ctx context.Context, identifier string, opts ...Options) Span {
	options := jaeger.getOpts(opts...)
	span, _ := opentracing.StartSpanFromContext(ctx, identifier, options...)
	return span
}

func (jaeger *JaegerTracer) StartSpanFromRequest(request *http.Request, identifier string, opts ...Options) Span {
	spanCtx, _ := jaeger.extract(jaeger.tracer, request)
	span := jaeger.tracer.StartSpan(identifier, ext.RPCServerOption(spanCtx))
	return span
}

func (jaeger *JaegerTracer) getOpts(opts ...Options) []opentracing.StartSpanOption {
	options := make([]opentracing.StartSpanOption, len(opts))
	for index, opt := range opts {
		options[index] = opentracing.Tag{Key: opt.Key, Value: opt.Value}
	}
	return options
}

func (jaeger *JaegerTracer) extract(tracer opentracing.Tracer, request *http.Request) (opentracing.SpanContext, error) {
	return tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header))
}
