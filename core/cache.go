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

	errs := model.Copy(s, settings)

	if errs != nil {
		panic(fmt.Errorf("mapping errors occurred: %v", errs))
	}

	return &cacheProxy{
		adapter:  state,
		settings: settings,
	}
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

func (c *cacheProxy) IsConnected() bool {
	return c.adapter.IsConnected()
}

func (c *cacheProxy) Connect(ctx context.Context) error {
	c.cache = getCache(c.settings)
	c.ch = make(chan rxgo.Item, 1)

	if c.settings.UseBatch {
		observable := rxgo.FromChannel(c.ch).
			BufferWithTime(rxgo.WithDuration(c.settings.BatchBufferInterval))

		go c.subscribeBatch(observable)
	}

	return c.adapter.Connect(ctx)
}

func (c *cacheProxy) Disconnect() {
	c.adapter.Disconnect()

	c.cache.DeleteExpired()

	close(c.ch)
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
