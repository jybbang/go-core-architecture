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

	err := r.Find(ctx, dto.ID, dto2)

	if !reflect.DeepEqual(dto2, dto) {
		t.Errorf("Test_queryRepositoryService_Find() result = %v, expect %v", dto2, dto)
	}

	if err != nil {
		t.Errorf("Test_queryRepositoryService_Find() err = %v", err)
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
	err := r.Find(ctx, dto.ID, dto2)

	if !errors.Is(err, core.ErrNotFound) {
		t.Errorf("TestStateService_GetNotFoundShouldBeError() err = %v, expect %v", err, core.ErrNotFound)
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

	ok, err := r.Any(ctx)

	if ok != true {
		t.Errorf("Test_queryRepositoryService_Any() ok = %v, expect %v", ok, true)
	}

	if err != nil {
		t.Errorf("Test_queryRepositoryService_Any() err = %v", err)
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

	ok, err := r.AnyWithFilter(ctx, "", "")

	if ok != true {
		t.Errorf("Test_queryRepositoryService_AnyWithFilter() ok = %v, expect %v", ok, true)
	}

	if err != nil {
		t.Errorf("Test_queryRepositoryService_AnyWithFilter() err = %v", err)
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

	result, err := r.Count(ctx)

	if result < int64(expect) {
		t.Errorf("Test_queryRepositoryService_Count() result = %v, expect %v", result, expect)
	}

	if err != nil {
		t.Errorf("Test_queryRepositoryService_Count() err = %v", err)
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

	result, err := r.CountWithFilter(ctx, "", "")

	if result != int64(expect) {
		t.Errorf("Test_queryRepositoryService_CountWithFilter() result = %v, expect %v", result, expect)
	}

	if err != nil {
		t.Errorf("Test_queryRepositoryService_CountWithFilter() err = %v", err)
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
	err := r.List(ctx, &dest)

	if err != nil {
		t.Errorf("Test_queryRepositoryService_List() err = %v", err)
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
	err := r.ListWithFilter(ctx, "", "", &dest)

	cnt := len(dest)

	if cnt != cntExpect {
		t.Errorf("Test_queryRepositoryService_ListWithFilter() cnt = %v, expect %v", cnt, cntExpect)
	}

	if err != nil {
		t.Errorf("Test_queryRepositoryService_ListWithFilter() err = %v", err)
	}
}
