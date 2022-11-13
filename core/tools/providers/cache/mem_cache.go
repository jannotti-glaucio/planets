package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jannotti-glaucio/planets/core/tools/providers/tracer"
)

//MemCache struct for assign ICacheProvider
type MemCache struct {
	Client     *memcache.Client
	Expiration time.Duration
}

var singletonConnection *MemCache = nil

//New create a new instance
func New(host string) (*MemCache, error) {
	if singletonConnection == nil {
		client := memcache.New(host)
		if err := client.Ping(); err != nil {
			return nil, err
		}

		return &MemCache{
			Client: client,
		}, nil
	}
	return singletonConnection, nil
}

//Assign interface
var _ ICacheProvider = &MemCache{}

func (cache *MemCache) WithExpiration(expiration time.Duration) ICacheProvider {
	cache.Expiration = expiration / time.Second
	return cache
}

func (cache *MemCache) Set(ctx context.Context, key string, value interface{}) error {
	identifierTracer := "mem.cache.Set"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: key, Value: value})
	defer span.Finish()

	valueBytes, err := cache.encode(value)
	if err != nil {
		return err
	}

	cache.Client.Set(&memcache.Item{
		Key:        key,
		Value:      valueBytes,
		Expiration: int32(cache.Expiration),
	})
	return nil
}

func (cache *MemCache) Get(ctx context.Context, key string, value interface{}) error {
	identifierTracer := "mem.cache.Get"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: key, Value: value})
	defer span.Finish()

	item, err := cache.Client.Get(key)
	if err != nil {
		return err
	}

	if err := cache.decode(item.Value, value); err != nil {
		return err
	}

	return nil
}

func (cache *MemCache) Delete(ctx context.Context, key string) error {
	identifierTracer := "mem.cache.Delete"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: key, Value: ""})
	defer span.Finish()

	return cache.Client.Delete(key)
}

func (cache *MemCache) Clear(ctx context.Context) error {
	identifierTracer := "mem.cache.Clear"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer)
	defer span.Finish()

	return cache.Client.FlushAll()
}

func (cache *MemCache) encode(v interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (cache *MemCache) decode(data []byte, v interface{}) error {
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)

	return decoder.Decode(v)
}
