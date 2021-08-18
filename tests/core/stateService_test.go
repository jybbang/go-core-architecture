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

func Test_stateService_Has(t *testing.T) {
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
		t.Errorf("Test_stateService_Has() ok = %v, expect %v", result.V, true)
	}

	if result.E != nil {
		t.Errorf("Test_stateService_Has() err = %v", result.E)
	}
}

func Test_stateService_HasNotFoundShouldBeFalseNadNoError(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	result := s.Has(ctx, "zxc")

	if result.V != false {
		t.Errorf("Test_stateService_HasNotFoundShouldBeFalseNadNoError() ok = %v, expect %v", result.V, false)
	}

	if result.E != nil {
		t.Errorf("Test_stateService_HasNotFoundShouldBeFalseNadNoError() err = %v", result.E)
	}
}

func Test_stateService_Get(t *testing.T) {
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
		t.Errorf("Test_stateService_Get() dest = %v, expect %v", dest, expect)
	}

	if result.E != nil {
		t.Errorf("Test_stateService_Get() err = %v", result.E)
	}
}

func Test_stateService_GetNotFoundShouldBeError(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	dest := &okCommand{}
	result := s.Get(ctx, "zxc", dest)

	if !errors.Is(result.E, core.ErrNotFound) {
		t.Errorf("Test_stateService_GetNotFoundShouldBeError() err = %v, expect %v", result.E, core.ErrNotFound)
	}
}

func Test_stateService_Set(t *testing.T) {
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
		t.Errorf("Test_stateService_Set() err = %v", result.E)
	}

	dest := &okCommand{}
	s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("Test_stateService_Set() dest = %v, expect %v", dest, expect)
	}
}

func Test_stateService_Delete(t *testing.T) {
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
		t.Errorf("Test_stateService_Delete() ok = %v, expect %v", result.V, false)
	}

	if result.E != nil {
		t.Errorf("Test_stateService_Delete() err = %v", result.E)
	}
}

func Test_stateService_DeleteNotFoundShouldBeNoError(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	s := core.NewStateServiceBuilder().
		StateAdapter(mock).
		Create()

	key := "qwe"

	result := s.Delete(ctx, key)

	if result.E != nil {
		t.Errorf("Test_stateService_DeleteNotFoundShouldBeNoError() err = %v", result.E)
	}
}
