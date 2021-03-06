package core

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/mocks"
)

func Test_cache_Has(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		UseCache(core.CacheSettings{ItemExpiration: 10 * time.Second}).
		Create()

	count := 10000
	key := "qwe"
	expect := testModel{
		Expect: 123,
	}

	s.Set(ctx, key, &expect)

	for i := 0; i < count; i++ {
		go s.Has(ctx, key)
	}

	time.Sleep(1 * time.Second)

	result := s.Has(ctx, key)

	if result.V != true {
		t.Errorf("Test_cache_Has() ok = %v, expect %v", result.V, true)
	}

	if result.E != nil {
		t.Errorf("Test_cache_Has() err = %v", result.E)
	}
}

func Test_cache_HasShouldBeOkAfterCacheExpired(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		UseCache(core.CacheSettings{ItemExpiration: 1 * time.Second}).
		Create()

	count := 10000
	key := "qwe"
	expect := testModel{
		Expect: 123,
	}

	s.Set(ctx, key, &expect)

	for i := 0; i < count; i++ {
		go s.Has(ctx, key)
	}

	time.Sleep(2 * time.Second)

	result := s.Has(ctx, key)

	if result.V != true {
		t.Errorf("Test_cache_HasShouldBeOkAfterCacheExpired() ok = %v, expect %v", result.V, true)
	}

	if result.E != nil {
		t.Errorf("Test_cache_HasShouldBeOkAfterCacheExpired() err = %v", result.E)
	}
}

func Test_cache_Get(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		UseCache(core.CacheSettings{ItemExpiration: 10 * time.Second}).
		Create()

	count := 10000
	key := "qwe"
	expect := &testModel{
		Expect: 123,
	}

	s.Set(ctx, key, expect)

	dest := &testModel{}
	for i := 0; i < count; i++ {
		go s.Get(ctx, key, dest)
	}

	time.Sleep(1 * time.Second)

	result := s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("Test_cache_Get() dest = %v, expect %v", dest, expect)
	}

	if result.E != nil {
		t.Errorf("Test_cache_Get() err = %v", result.E)
	}
}

func Test_cache_GetShouldBeOkAfterCacheExpired(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		UseCache(core.CacheSettings{ItemExpiration: 1 * time.Second}).
		Create()

	count := 10000
	key := "qwe"
	expect := &testModel{
		Expect: 123,
	}

	s.Set(ctx, key, expect)

	dest := &testModel{}
	for i := 0; i < count; i++ {
		go s.Get(ctx, key, dest)
	}

	time.Sleep(2 * time.Second)

	result := s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("Test_cache_GetShouldBeOkAfterCacheExpired() dest = %v, expect %v", dest, expect)
	}

	if result.E != nil {
		t.Errorf("Test_cache_GetShouldBeOkAfterCacheExpired() err = %v", result.E)
	}
}

func Test_cache_GetNotFoundShouldBeError(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		UseCache(core.CacheSettings{ItemExpiration: 10 * time.Second}).
		Create()

	dest := &testModel{}
	result := s.Get(ctx, "zxc", dest)

	if !errors.Is(result.E, core.ErrNotFound) {
		t.Errorf("Test_cache_GetNotFoundShouldBeError() err = %v, expect %v", result.E, core.ErrNotFound)
	}
}

func Test_cache_Set(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		UseCache(core.CacheSettings{ItemExpiration: 10 * time.Second}).
		Create()

	count := 10000
	key := "qwe"
	expect := &testModel{
		Expect: 123,
	}
	for i := 0; i < count; i++ {
		go s.Set(ctx, key, expect)
	}

	time.Sleep(1 * time.Second)

	result := s.Set(ctx, key, expect)

	if result.E != nil {
		t.Errorf("Test_cache_Set() err = %v", result.E)
	}

	dest := &testModel{}
	s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("Test_cache_Set() dest = %v, expect %v", dest, expect)
	}
}

func Test_cache_SetWithUseBatch(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		UseCache(core.CacheSettings{
			ItemExpiration:      1 * time.Second,
			UseBatch:            true,
			BatchBufferInterval: 1 * time.Second}).
		Create()

	expect := 10000
	for i := 0; i < expect; i++ {
		go s.Set(ctx, strconv.Itoa(i), &testModel{
			Expect: i,
		})
	}

	time.Sleep(2 * time.Second)

	dest := &testModel{}
	s.Get(ctx, strconv.Itoa(expect), dest)

	if dest.Expect == expect {
		t.Errorf("Test_cache_SetUseBatch() dest.Expect = %d, expect %d", dest.Expect, expect)
	}
}

func Test_cache_BatchSet(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		UseCache(core.CacheSettings{ItemExpiration: 10 * time.Second}).
		Create()

	expect := 10000

	kvs := make([]core.KV, 0)
	for i := 0; i < expect; i++ {
		kvs = append(
			kvs,
			core.KV{
				K: strconv.Itoa(i),
				V: &testModel{
					Expect: i,
				}})
	}

	result := s.BatchSet(ctx, kvs)

	if result.E != nil {
		t.Errorf("Test_cache_BatchSet() err = %v", result.E)
	}

	dest := &testModel{}
	s.Get(ctx, strconv.Itoa(expect), dest)

	if dest.Expect == expect {
		t.Errorf("Test_cache_BatchSet() dest.Expect = %d, expect %d", dest.Expect, expect)
	}
}

func Test_cache_Delete(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		UseCache(core.CacheSettings{ItemExpiration: 10 * time.Second}).
		Create()

	count := 10000
	key := "qwe"
	expect := testModel{
		Expect: 123,
	}

	s.Set(ctx, key, &expect)
	for i := 0; i < count; i++ {
		go s.Delete(ctx, key)
	}

	time.Sleep(1 * time.Second)

	s.Delete(ctx, key)

	result := s.Has(ctx, key)

	if result.V != false {
		t.Errorf("Test_cache_Delete() ok = %v, expect %v", result.V, false)
	}

	if result.E != nil {
		t.Errorf("Test_cache_Delete() err = %v", result.E)
	}
}

func Test_cache_DeleteNotFoundShouldBeNoError(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		UseCache(core.CacheSettings{ItemExpiration: 10 * time.Second}).
		Create()

	key := "qwe"

	result := s.Delete(ctx, key)

	if result.E != nil {
		t.Errorf("Test_cache_DeleteNotFoundShouldBeNoError() err = %v", result.E)
	}
}
