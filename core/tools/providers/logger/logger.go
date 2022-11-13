package logger

import (
	"context"

	log "github.com/amoghe/distillog"
	"github.com/jannotti-glaucio/planets/core/tools/providers/tracer"
)

//Logger struct for assign interface
type Logger struct{}

//New create a new instance
func New() *Logger {
	return &Logger{}
}

//Token implements the TokenInterface
var _ ILoggerProvider = &Logger{}

//Error ...
func (l *Logger) Error(ctx context.Context, v ...interface{}) {
	identifierTracer := "logger.error"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer, Value: v})
	defer span.Finish()

	log.Errorln(v...)
}

//Info ...
func (l *Logger) Info(ctx context.Context, v ...interface{}) {
	identifierTracer := "logger.info"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer, Value: v})
	defer span.Finish()

	log.Infoln(v...)
}
