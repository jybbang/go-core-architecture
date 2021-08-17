package core

import (
	"context"
	"errors"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/mocks"
)

func Test_queryRepositoryService_Find(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel), "testModel").
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	r.Add(ctx, dto)

	dto2 := new(testModel)
	count := 10000
	for i := 0; i < count; i++ {
		go r.Find(ctx, dto.ID, dto2)
	}

	time.Sleep(1 * time.Second)

	result := r.Find(ctx, dto.ID, dto2)

	if !reflect.DeepEqual(dto2.ID, dto.ID) || !reflect.DeepEqual(dto2.Expect, dto.Expect) {
		t.Errorf("Test_queryRepositoryService_Find() result = %v, expect %v", dto2, dto)
	}

	if result.E != nil {
		t.Errorf("Test_queryRepositoryService_Find() err = %v", result.E)
	}
}

func Test_queryRepositoryService_FindnotFoundShouldBeError(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel), "testModel").
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	dto2 := new(testModel)
	result := r.Find(ctx, dto.ID, dto2)

	if !errors.Is(result.E, core.ErrNotFound) {
		t.Errorf("TestStateService_GetNotFoundShouldBeError() err = %v, expect %v", result.E, core.ErrNotFound)
	}
}

func Test_queryRepositoryService_Any(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel), "testModel").
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	r.Add(ctx, dto)

	count := 10000
	for i := 0; i < count; i++ {
		go r.Any(ctx)
	}

	time.Sleep(1 * time.Second)

	result := r.Any(ctx)

	if result.V != true {
		t.Errorf("Test_queryRepositoryService_Any() ok = %v, expect %v", result.V, true)
	}

	if result.E != nil {
		t.Errorf("Test_queryRepositoryService_Any() err = %v", result.E)
	}
}

func Test_queryRepositoryService_AnyWithFilter(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel), "testModel").
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	r.Add(ctx, dto)

	result := r.AnyWithFilter(ctx, "", "")

	if result.V != true {
		t.Errorf("Test_queryRepositoryService_AnyWithFilter() ok = %v, expect %v", result.V, true)
	}

	if result.E != nil {
		t.Errorf("Test_queryRepositoryService_AnyWithFilter() err = %v", result.E)
	}
}

func Test_queryRepositoryService_Count(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel), "testModel").
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	expect := 10000
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = i

		r.Add(ctx, dto)
	}

	count := 10000
	for i := 0; i < count; i++ {
		go r.Count(ctx)
	}

	time.Sleep(1 * time.Second)

	result := r.Count(ctx)

	if result.V.(int64) < int64(expect) {
		t.Errorf("Test_queryRepositoryService_Count() result = %v, expect %v", result.V, expect)
	}

	if result.E != nil {
		t.Errorf("Test_queryRepositoryService_Count() err = %v", result.E)
	}
}

func Test_queryRepositoryService_CountWithFilter(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel), "testModel").
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	expect := 10000
	rand.Seed(time.Now().UnixNano())
	random := rand.Int()
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = random

		r.Add(ctx, dto)
	}

	result := r.CountWithFilter(ctx, "", "")

	if result.V.(int64) != int64(expect) {
		t.Errorf("Test_queryRepositoryService_CountWithFilter() result = %v, expect %v", result.V, expect)
	}

	if result.E != nil {
		t.Errorf("Test_queryRepositoryService_CountWithFilter() err = %v", result.E)
	}
}

func Test_queryRepositoryService_List(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel), "testModel").
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	expect := 100
	cntExpect := 0
	rand.Seed(time.Now().UnixNano())
	random := rand.Int()
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = random

		r.Add(ctx, dto)
		cntExpect += 1
	}

	var dest = make([]*testModel, 0)
	result := r.List(ctx, &dest)

	if result.E != nil {
		t.Errorf("Test_queryRepositoryService_List() err = %v", result.E)
	}

	cnt := len(dest)

	if cnt < cntExpect {
		t.Errorf("Test_queryRepositoryService_List() cnt = %v, expect %v", cnt, cntExpect)
	}
}

func Test_queryRepositoryService_ListWithFilter(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel), "testModel").
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	expect := 10000
	cntExpect := 0
	rand.Seed(time.Now().UnixNano())
	random := rand.Int()
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = random

		r.Add(ctx, dto)
		cntExpect += 1
	}

	var dest = make([]*testModel, 0)
	result := r.ListWithFilter(ctx, "", "", &dest)

	cnt := len(dest)

	if cnt != cntExpect {
		t.Errorf("Test_queryRepositoryService_ListWithFilter() cnt = %v, expect %v", cnt, cntExpect)
	}

	if result.E != nil {
		t.Errorf("Test_queryRepositoryService_ListWithFilter() err = %v", result.E)
	}
}
