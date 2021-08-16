package infrastructure

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/mocks"
	"github.com/jybbang/go-core-architecture/infrastructure/redis"
)

func TestRedisStateService_ConnectionTimeout(t *testing.T) {
	timeout := time.Duration(1 * time.Second)
	ctx, c := context.WithTimeout(context.TODO(), timeout)
	defer c()

	redis := redis.NewRedisAdapter(ctx, redis.RedisSettings{
		Host: "localhost:6379",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(redis).
		Create()

	time.Sleep(timeout)

	key := "qwe"
	expect := okCommand{
		Expect: 123,
	}
	err := s.Set(ctx, key, &expect)

	if err != context.DeadlineExceeded {
		t.Errorf("TestRedisStateService_ConnectionTimeout() err = %v, expect %v", err, context.DeadlineExceeded)
	}
}

func TestRedisStateService_Has(t *testing.T) {
	ctx := context.Background()

	redis := redis.NewRedisAdapter(ctx, redis.RedisSettings{
		Host: "localhost:6379",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(redis).
		Create()

	count := 10000
	key := "qwe"
	expect := okCommand{
		Expect: 123,
	}
	s.Set(ctx, key, &expect)

	for i := 0; i < count; i++ {
		go s.Has(ctx, key)
	}

	time.Sleep(1 * time.Second)

	ok, err := s.Has(ctx, key)

	if ok != true {
		t.Errorf("TestRedisStateService_Has() ok = %v, expect %v", ok, true)
	}

	if err != nil {
		t.Errorf("TestRedisStateService_Has() err = %v", err)
	}
}

func TestRedisStateService_HasNotFoundShouldBeFalseNadNoError(t *testing.T) {
	ctx := context.Background()

	redis := redis.NewRedisAdapter(ctx, redis.RedisSettings{
		Host: "localhost:6379",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(redis).
		Create()

	ok, err := s.Has(ctx, "zxc")

	if ok != false {
		t.Errorf("TestRedisStateService_HasNotFoundShouldBeFalseNadNoError() ok = %v, expect %v", ok, false)
	}

	if err != nil {
		t.Errorf("TestRedisStateService_HasNotFoundShouldBeFalseNadNoError() err = %v", err)
	}
}

func TestRedisStateService_Get(t *testing.T) {
	ctx := context.Background()

	redis := redis.NewRedisAdapter(ctx, redis.RedisSettings{
		Host: "localhost:6379",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(redis).
		Create()

	count := 10000
	key := "qwe"
	expect := &okCommand{
		Expect: 123,
	}
	s.Set(ctx, key, expect)

	dest := &okCommand{}
	for i := 0; i < count; i++ {
		go s.Get(ctx, key, dest)
	}

	err := s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("TestRedisStateService_Get() dest = %v, expect %v", dest, expect)
	}

	if err != nil {
		t.Errorf("TestRedisStateService_Get() err = %v", err)
	}
}

func TestRedisStateService_GetNotFoundShouldBeError(t *testing.T) {
	ctx := context.Background()

	redis := redis.NewRedisAdapter(ctx, redis.RedisSettings{
		Host: "localhost:6379",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(redis).
		Create()

	dest := &okCommand{}
	err := s.Get(ctx, "zxc", dest)

	if err != core.ErrNotFound {
		t.Errorf("TestRedisStateService_GetNotFoundShouldBeError() err = %v, expect %v", err, core.ErrNotFound)
	}
}

func TestRedisStateService_Set(t *testing.T) {
	ctx := context.Background()

	redis := redis.NewRedisAdapter(ctx, redis.RedisSettings{
		Host: "localhost:6379",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(redis).
		Create()

	count := 10000
	key := "qwe"
	expect := &okCommand{
		Expect: 123,
	}
	for i := 0; i < count; i++ {
		go s.Set(ctx, key, expect)
	}

	err := s.Set(ctx, key, expect)

	if err != nil {
		t.Errorf("TestRedisStateService_Set() err = %v", err)
	}

	dest := &okCommand{}
	s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("TestRedisStateService_Set() dest = %v, expect %v", dest, expect)
	}
}

func TestRedisStateService_Delete(t *testing.T) {
	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	count := 10000
	key := "qwe"
	expect := okCommand{
		Expect: 123,
	}
	ctx := context.Background()

	for i := 0; i < count; i++ {
		go s.Set(ctx, key, &expect)
		go s.Delete(ctx, key)
	}

	time.Sleep(1 * time.Second)

	s.Delete(ctx, key)

	ok, err := s.Has(ctx, key)

	if ok != false {
		t.Errorf("TestRedisStateService_Delete() ok = %v, expect %v", ok, false)
	}

	if err != nil {
		t.Errorf("TestRedisStateService_Delete() err = %v", err)
	}
}

func TestRedisStateService_DeleteNotFoundShouldBeNoError(t *testing.T) {
	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	key := "qwe"
	ctx := context.Background()

	err := s.Delete(ctx, key)

	if err != nil {
		t.Errorf("TestRedisStateService_DeleteNotFoundShouldBeNoError() err = %v", err)
	}
}
