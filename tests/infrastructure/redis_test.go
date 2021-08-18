package infrastructure

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/redis"
)

func Test_redisStateService_ConnectionTimeout(t *testing.T) {
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
	expect := testModel{
		Expect: 123,
	}

	result := s.Set(ctx, key, &expect)

	if !errors.Is(result.E, context.DeadlineExceeded) {
		t.Errorf("Test_redisStateService_ConnectionTimeout() err = %v, expect %v", result.E, context.DeadlineExceeded)
	}
}

func Test_redisStateService_Has(t *testing.T) {
	ctx := context.Background()

	redis := redis.NewRedisAdapter(ctx, redis.RedisSettings{
		Host: "localhost:6379",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(redis).
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
		t.Errorf("Test_redisStateService_Has() ok = %v, expect %v", result.V, true)
	}

	if result.E != nil {
		t.Errorf("Test_redisStateService_Has() err = %v", result.E)
	}
}

func Test_redisStateService_Get(t *testing.T) {
	ctx := context.Background()

	redis := redis.NewRedisAdapter(ctx, redis.RedisSettings{
		Host: "localhost:6379",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(redis).
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
		t.Errorf("Test_redisStateService_Get() dest = %v, expect %v", dest, expect)
	}

	if result.E != nil {
		t.Errorf("Test_redisStateService_Get() err = %v", result.E)
	}
}

func Test_redisStateService_GetNotFoundShouldBeError(t *testing.T) {
	ctx := context.Background()

	redis := redis.NewRedisAdapter(ctx, redis.RedisSettings{
		Host: "localhost:6379",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(redis).
		Create()

	dest := &testModel{}
	result := s.Get(ctx, "zxc", dest)

	if !errors.Is(result.E, core.ErrNotFound) {
		t.Errorf("Test_redisStateService_GetNotFoundShouldBeError() err = %v, expect %v", result.E, core.ErrNotFound)
	}
}

func Test_redisStateService_Set(t *testing.T) {
	ctx := context.Background()

	redis := redis.NewRedisAdapter(ctx, redis.RedisSettings{
		Host: "localhost:6379",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(redis).
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
		t.Errorf("Test_redisStateService_Set() err = %v", result.E)
	}

	dest := &testModel{}
	s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("Test_redisStateService_Set() dest = %v, expect %v", dest, expect)
	}
}

func Test_redisStateService_Delete(t *testing.T) {
	ctx := context.Background()

	redis := redis.NewRedisAdapter(ctx, redis.RedisSettings{
		Host: "localhost:6379",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(redis).
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

	if result.E != nil {
		t.Errorf("Test_redisStateService_Delete() err = %v", result.E)
	}

	if result.V != false {
		t.Errorf("Test_redisStateService_Delete() ok = %v, expect %v", result.V, false)
	}
}

func Test_redisStateService_DeleteNotFoundShouldBeNoError(t *testing.T) {
	ctx := context.Background()

	redis := redis.NewRedisAdapter(ctx, redis.RedisSettings{
		Host: "localhost:6379",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(redis).
		Create()

	key := "qwe"

	result := s.Delete(ctx, key)

	if result.E != nil {
		t.Errorf("Test_redisStateService_DeleteNotFoundShouldBeNoError() err = %v", result.E)
	}
}
