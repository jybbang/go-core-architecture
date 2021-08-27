package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/reactivex/rxgo/v2"
	"gopkg.in/jeevatkm/go-model.v1"
)

type cacheProxy struct {
	cache    *cache.Cache
	adapter  stateAdapter
	settings CacheSettings
	ch       chan rxgo.Item
}

type CacheSettings struct {
	ItemExpiration       time.Duration `model:",omitempty"`
	CacheCleanupInterval time.Duration `model:",omitempty"`
	BatchBufferInterval  time.Duration `model:",omitempty"`
	BatchTimeout         time.Duration `model:",omitempty"`
	UseBatch             bool
}

var cacheSync sync.Once

var cacheInstance *cache.Cache

func getCache(settings CacheSettings) *cache.Cache {
	if cacheInstance == nil {
		cacheSync.Do(
			func() {
				cacheInstance = cache.New(settings.ItemExpiration, settings.CacheCleanupInterval)
			})
	}
	return cacheInstance
}

func newCache(state stateAdapter, settings CacheSettings) *cacheProxy {
	s := &CacheSettings{
		ItemExpiration:       time.Duration(5 * time.Minute),
		CacheCleanupInterval: time.Duration(10 * time.Minute),
		BatchBufferInterval:  time.Duration(1 * time.Minute),
		BatchTimeout:         time.Duration(10 * time.Second),
	}

	err := model.Copy(s, settings)
	if err != nil {
		panic(fmt.Errorf("mapping errors occurred: %v", err))
	}

	cache := &cacheProxy{
		cache:    getCache(*s),
		adapter:  state,
		settings: *s,
		ch:       make(chan rxgo.Item, 1),
	}

	if settings.UseBatch {
		observable := rxgo.FromChannel(cache.ch).
			BufferWithTime(rxgo.WithDuration(s.BatchBufferInterval))

		go cache.subscribeBatch(observable)
	}

	return cache
}

func (c *cacheProxy) subscribeBatch(observable rxgo.Observable) {
	ch := observable.Observe()

	for {
		items := <-ch

		batch, ok := items.V.([]interface{})
		if !ok || len(batch) == 0 {
			continue
		}

		kvs := make(map[string]KV)
		for _, v := range batch {
			if kv, ok := v.(KV); ok {
				kvs[kv.K] = kv
			}
		}

		values := make([]KV, 0, len(kvs))
		for _, v := range kvs {
			values = append(values, v)
		}

		timeout, cancel := context.WithTimeout(context.Background(), c.settings.BatchTimeout)
		c.BatchSet(timeout, values)
		cancel()
	}
}

func (c *cacheProxy) Close() {
	c.adapter.Close()
	close(c.ch)
	c.cache.DeleteExpired()
}

func (c *cacheProxy) Has(ctx context.Context, key string) bool {
	_, ok := c.cache.Get(key)
	if !ok {
		return c.adapter.Has(ctx, key)
	}
	return ok
}

func (c *cacheProxy) Get(ctx context.Context, key string, dest interface{}) error {
	value, ok := c.cache.Get(key)
	if !ok {
		err := c.adapter.Get(ctx, key, dest)
		if err == nil {
			c.cache.SetDefault(key, dest)
		}
		return err
	}

	err := model.Copy(dest, value)
	if err != nil {
		return fmt.Errorf("mapping errors occurred: %v", err)
	}

	return nil
}

func (c *cacheProxy) Set(ctx context.Context, key string, value interface{}) error {
	c.cache.SetDefault(key, value)
	if c.settings.UseBatch {
		c.ch <- rxgo.Item{
			V: KV{
				K: key,
				V: value,
			},
		}
	} else {
		c.adapter.Set(ctx, key, value)
	}
	return nil
}

func (c *cacheProxy) Delete(ctx context.Context, key string) error {
	c.cache.Delete(key)
	return c.adapter.Delete(ctx, key)
}

func (c *cacheProxy) BatchSet(ctx context.Context, kvs []KV) error {
	return c.adapter.BatchSet(ctx, kvs)
}
