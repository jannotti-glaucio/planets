package logger

import "context"

//ILoggerProvider for sign token struct
type ILoggerProvider interface {
	Error(ctx context.Context, v ...interface{})
	Info(ctx context.Context, v ...interface{})
}
