package infrastructure

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/etcd"
)

func Test_etcdStateService_ConnectionTimeout(t *testing.T) {
	timeout := time.Duration(1 * time.Second)
	ctx, c := context.WithTimeout(context.TODO(), timeout)
	defer c()

	etcd := etcd.NewEtcdAdapter(ctx, etcd.EtcdSettings{
		Endpoints: []string{"localhost:2379"},
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(etcd).
		Create()

	time.Sleep(timeout)

	key := "qwe"
	expect := testModel{
		Expect: 123,
	}

	result := s.Set(ctx, key, &expect)

	if !errors.Is(result.E, context.DeadlineExceeded) {
		t.Errorf("Test_etcdStateService_ConnectionTimeout() err = %v, expect %v", result.E, context.DeadlineExceeded)
	}
}

func Test_etcdStateService_Has(t *testing.T) {
	ctx := context.Background()

	etcd := etcd.NewEtcdAdapter(ctx, etcd.EtcdSettings{
		Endpoints: []string{"localhost:2379"},
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(etcd).
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
		t.Errorf("Test_etcdStateService_Has() ok = %v, expect %v", result.V, true)
	}

	if result.E != nil {
		t.Errorf("Test_etcdStateService_Has() err = %v", result.E)
	}
}

func Test_etcdStateService_Get(t *testing.T) {
	ctx := context.Background()

	etcd := etcd.NewEtcdAdapter(ctx, etcd.EtcdSettings{
		Endpoints: []string{"localhost:2379"},
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(etcd).
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
		t.Errorf("Test_etcdStateService_Get() dest = %v, expect %v", dest, expect)
	}

	if result.E != nil {
		t.Errorf("Test_etcdStateService_Get() err = %v", result.E)
	}
}

func Test_etcdStateService_GetNotFoundShouldBeError(t *testing.T) {
	ctx := context.Background()

	etcd := etcd.NewEtcdAdapter(ctx, etcd.EtcdSettings{
		Endpoints: []string{"localhost:2379"},
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(etcd).
		Create()

	dest := &testModel{}
	result := s.Get(ctx, "zxc", dest)

	if !errors.Is(result.E, core.ErrNotFound) {
		t.Errorf("TestStateService_GetNotFoundShouldBeError() err = %v, expect %v", result.E, core.ErrNotFound)
	}
}

func Test_etcdStateService_Set(t *testing.T) {
	ctx := context.Background()

	etcd := etcd.NewEtcdAdapter(ctx, etcd.EtcdSettings{
		Endpoints: []string{"localhost:2379"},
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(etcd).
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
		t.Errorf("Test_etcdStateService_Set() err = %v", result.E)
	}

	dest := &testModel{}
	s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("Test_etcdStateService_Set() dest = %v, expect %v", dest, expect)
	}
}

func Test_etcdStateService_BatchSet(t *testing.T) {
	ctx := context.Background()

	etcd := etcd.NewEtcdAdapter(ctx, etcd.EtcdSettings{
		Endpoints: []string{"localhost:2379"},
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(etcd).
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
		t.Errorf("Test_etcdStateService_BatchSet() err = %v", result.E)
	}

	dest := &testModel{}
	s.Get(ctx, strconv.Itoa(expect), dest)

	if dest.Expect == expect {
		t.Errorf("Test_etcdStateService_BatchSet() dest.Expect = %d, expect %d", dest.Expect, expect)
	}
}

func Test_etcdStateService_Delete(t *testing.T) {
	ctx := context.Background()

	etcd := etcd.NewEtcdAdapter(ctx, etcd.EtcdSettings{
		Endpoints: []string{"localhost:2379"},
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(etcd).
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
		t.Errorf("Test_etcdStateService_Delete() err = %v", result.E)
	}

	if result.V != false {
		t.Errorf("Test_etcdStateService_Delete() ok = %v, expect %v", result.V, false)
	}
}

func Test_etcdStateService_DeleteNotFoundShouldBeNoError(t *testing.T) {
	ctx := context.Background()

	etcd := etcd.NewEtcdAdapter(ctx, etcd.EtcdSettings{
		Endpoints: []string{"localhost:2379"},
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(etcd).
		Create()

	key := "qwe"

	result := s.Delete(ctx, key)

	if result.E != nil {
		t.Errorf("Test_etcdStateService_DeleteNotFoundShouldBeNoError() err = %v", result.E)
	}
}
