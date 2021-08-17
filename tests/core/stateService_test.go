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
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
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

	result := s.Has(ctx, key)

	if result.V != true {
		t.Errorf("TestStateService_Has() ok = %v, expect %v", result.V, true)
	}

	if result.E != nil {
		t.Errorf("TestStateService_Has() err = %v", result.E)
	}
}

func TestStateService_HasNotFoundShouldBeFalseNadNoError(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	result := s.Has(ctx, "zxc")

	if result.V != false {
		t.Errorf("TestStateService_HasNotFoundShouldBeFalseNadNoError() ok = %v, expect %v", result.V, false)
	}

	if result.E != nil {
		t.Errorf("TestStateService_HasNotFoundShouldBeFalseNadNoError() err = %v", result.E)
	}
}

func TestStateService_Get(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
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

	time.Sleep(1 * time.Second)

	result := s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("TestStateService_Get() dest = %v, expect %v", dest, expect)
	}

	if result.E != nil {
		t.Errorf("TestStateService_Get() err = %v", result.E)
	}
}

func TestStateService_GetNotFoundShouldBeError(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	dest := &okCommand{}
	result := s.Get(ctx, "zxc", dest)

	if !errors.Is(result.E, core.ErrNotFound) {
		t.Errorf("TestStateService_GetNotFoundShouldBeError() err = %v, expect %v", result.E, core.ErrNotFound)
	}
}

func TestStateService_Set(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	count := 10000
	key := "qwe"
	expect := &okCommand{
		Expect: 123,
	}
	for i := 0; i < count; i++ {
		go s.Set(ctx, key, expect)
	}

	time.Sleep(1 * time.Second)

	result := s.Set(ctx, key, expect)

	if result.E != nil {
		t.Errorf("TestStateService_Set() err = %v", result.E)
	}

	dest := &okCommand{}
	s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("TestStateService_Set() dest = %v, expect %v", dest, expect)
	}
}

func TestStateService_Delete(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	count := 10000
	key := "qwe"
	expect := okCommand{
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
		t.Errorf("TestStateService_Delete() ok = %v, expect %v", result.V, false)
	}

	if result.E != nil {
		t.Errorf("TestStateService_Delete() err = %v", result.E)
	}
}

func TestStateService_DeleteNotFoundShouldBeNoError(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	key := "qwe"

	result := s.Delete(ctx, key)

	if result.E != nil {
		t.Errorf("TestStateService_DeleteNotFoundShouldBeNoError() err = %v", result.E)
	}
}
