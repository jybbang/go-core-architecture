package infrastructure

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/etcd"
	"github.com/jybbang/go-core-architecture/infrastructure/mocks"
)

func TestEtcdStateService_ConnectionTimeout(t *testing.T) {
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
	expect := okCommand{
		Expect: 123,
	}
	err := s.Set(ctx, key, &expect)

	if err != context.DeadlineExceeded {
		t.Errorf("TestEtcdStateService_ConnectionTimeout() err = %v, expect %v", err, context.DeadlineExceeded)
	}
}

func TestEtcdStateService_Has(t *testing.T) {
	ctx := context.Background()

	etcd := etcd.NewEtcdAdapter(ctx, etcd.EtcdSettings{
		Endpoints: []string{"localhost:2379"},
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(etcd).
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
		t.Errorf("TestEtcdStateService_Has() ok = %v, expect %v", ok, true)
	}

	if err != nil {
		t.Errorf("TestEtcdStateService_Has() err = %v", err)
	}
}

func TestEtcdStateService_HasNotFoundShouldBeFalseNadNoError(t *testing.T) {
	ctx := context.Background()

	etcd := etcd.NewEtcdAdapter(ctx, etcd.EtcdSettings{
		Endpoints: []string{"localhost:2379"},
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(etcd).
		Create()

	ok, err := s.Has(ctx, "zxc")

	if ok != false {
		t.Errorf("TestEtcdStateService_HasNotFoundShouldBeFalseNadNoError() ok = %v, expect %v", ok, false)
	}

	if err != nil {
		t.Errorf("TestEtcdStateService_HasNotFoundShouldBeFalseNadNoError() err = %v", err)
	}
}

func TestEtcdStateService_Get(t *testing.T) {
	ctx := context.Background()

	etcd := etcd.NewEtcdAdapter(ctx, etcd.EtcdSettings{
		Endpoints: []string{"localhost:2379"},
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(etcd).
		Create()

	count := 10000
	key := "qwe"
	expect := &okCommand{
		Expect: 123,
	}
	s.Set(ctx, key, expect)

	dest := &okCommand{
		Expect: 123,
	}
	for i := 0; i < count; i++ {
		go s.Get(ctx, key, dest)
	}

	err := s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("TestEtcdStateService_Get() dest = %v, expect %v", dest, expect)
	}

	if err != nil {
		t.Errorf("TestEtcdStateService_Get() err = %v", err)
	}
}

func TestEtcdStateService_GetNotFoundShouldBeError(t *testing.T) {
	ctx := context.Background()

	etcd := etcd.NewEtcdAdapter(ctx, etcd.EtcdSettings{
		Endpoints: []string{"localhost:2379"},
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(etcd).
		Create()

	dest := &okCommand{
		Expect: 123,
	}

	err := s.Get(ctx, "zxc", dest)

	if err != core.ErrNotFound {
		t.Errorf("TestEtcdStateService_GetNotFoundShouldBeError() err = %v, expect %v", err, core.ErrNotFound)
	}
}

func TestEtcdStateService_Set(t *testing.T) {
	ctx := context.Background()

	etcd := etcd.NewEtcdAdapter(ctx, etcd.EtcdSettings{
		Endpoints: []string{"localhost:2379"},
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(etcd).
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
		t.Errorf("TestEtcdStateService_Set() err = %v", err)
	}

	dest := &okCommand{}
	s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("TestEtcdStateService_Set() dest = %v, expect %v", dest, expect)
	}
}

func TestEtcdStateService_Delete(t *testing.T) {
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
	s.Set(ctx, key, &expect)

	for i := 0; i < count; i++ {
		go s.Delete(ctx, key)
	}
	s.Delete(ctx, key)

	time.Sleep(1 * time.Second)

	ok, err := s.Has(ctx, key)

	if ok != false {
		t.Errorf("TestEtcdStateService_Delete() ok = %v, expect %v", ok, false)
	}

	if err != nil {
		t.Errorf("TestEtcdStateService_Delete() err = %v", err)
	}
}

func TestEtcdStateService_DeleteNotFoundShouldBeNoError(t *testing.T) {
	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	key := "qwe"
	ctx := context.Background()

	err := s.Delete(ctx, key)

	if err != nil {
		t.Errorf("TestEtcdStateService_DeleteNotFoundShouldBeNoError() err = %v", err)
	}
}
