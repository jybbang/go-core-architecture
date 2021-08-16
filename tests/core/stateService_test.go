package core

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/mocks"
)

func TestStateService_Has(t *testing.T) {
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
		go s.Has(ctx, key)
	}

	time.Sleep(1 * time.Second)

	ok, err := s.Has(ctx, key)

	if ok != true {
		t.Errorf("TestStateService_Has() ok = %v, expect %v", ok, true)
	}

	if err != nil {
		t.Errorf("TestStateService_Has() err = %v", err)
	}
}

func TestStateService_HasNotFoundShouldBeFalseNadNoError(t *testing.T) {
	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	ctx := context.Background()
	ok, err := s.Has(ctx, "zxc")

	if ok != false {
		t.Errorf("TestStateService_HasNotFoundShouldBeFalseNadNoError() ok = %v, expect %v", ok, false)
	}

	if err != nil {
		t.Errorf("TestStateService_HasNotFoundShouldBeFalseNadNoError() err = %v", err)
	}
}

func TestStateService_Get(t *testing.T) {
	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	count := 10000
	key := "qwe"
	expect := &okCommand{
		Expect: 123,
	}
	ctx := context.Background()
	s.Set(ctx, key, expect)

	dest := &okCommand{}
	for i := 0; i < count; i++ {
		go s.Get(ctx, key, dest)
	}

	err := s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("TestStateService_Get() dest = %v, expect %v", dest, expect)
	}

	if err != nil {
		t.Errorf("TestStateService_Get() err = %v", err)
	}
}

func TestStateService_GetNotFoundShouldBeError(t *testing.T) {
	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	ctx := context.Background()

	dest := &okCommand{}
	err := s.Get(ctx, "zxc", dest)

	if !errors.Is(err, core.ErrNotFound) {
		t.Errorf("TestStateService_GetNotFoundShouldBeError() err = %v, expect %v", err, core.ErrNotFound)
	}
}

func TestStateService_Set(t *testing.T) {
	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	ctx := context.Background()

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
		t.Errorf("TestStateService_Set() err = %v", err)
	}

	dest := &okCommand{}
	s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("TestStateService_Set() dest = %v, expect %v", dest, expect)
	}
}

func TestStateService_Delete(t *testing.T) {
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
		t.Errorf("TestStateService_Delete() ok = %v, expect %v", ok, false)
	}

	if err != nil {
		t.Errorf("TestStateService_Delete() err = %v", err)
	}
}

func TestStateService_DeleteNotFoundShouldBeNoError(t *testing.T) {
	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	key := "qwe"
	ctx := context.Background()

	err := s.Delete(ctx, key)

	if err != nil {
		t.Errorf("TestStateService_DeleteNotFoundShouldBeNoError() err = %v", err)
	}
}
