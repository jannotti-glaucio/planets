package cache

import (
	"context"
	"time"
)

//ICacheProvider interface of cache
type ICacheProvider interface {
	WithExpiration(time.Duration) ICacheProvider
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}
