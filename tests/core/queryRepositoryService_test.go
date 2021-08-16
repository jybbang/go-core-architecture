package core

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/mocks"
)

func Test_queryRepositoryService_Find(t *testing.T) {
	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	ctx := context.Background()
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
	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	ctx := context.Background()

	dto2 := new(testModel)
	err := r.Find(ctx, dto.ID, dto2)

	if err != core.ErrNotFound {
		t.Errorf("TestStateService_GetNotFoundShouldBeError() err = %v, expect %v", err, core.ErrNotFound)
	}
}

func Test_queryRepositoryService_Any(t *testing.T) {
	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	ctx := context.Background()
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
	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	ctx := context.Background()
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
	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	ctx := context.Background()
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
	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	ctx := context.Background()
	expect := 10000
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = i

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
	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	ctx := context.Background()
	expect := 10000
	sumExpect := 0
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = i

		r.Add(ctx, dto)
		sumExpect += i
	}

	dtos2, err := r.List(ctx)

	if err != nil {
		t.Errorf("Test_queryRepositoryService_List() err = %v", err)
	}

	sum := 0
	for _, v := range dtos2 {
		sum += v.(*testModel).Expect
	}

	if sum != sumExpect {
		t.Errorf("Test_queryRepositoryService_List() sum = %v, expect %v", sum, sumExpect)
	}
}

func Test_queryRepositoryService_ListWithFilter(t *testing.T) {
	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	ctx := context.Background()
	expect := 10000
	sumExpect := 0
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = i

		r.Add(ctx, dto)
		sumExpect += i
	}

	dtos2, err := r.ListWithFilter(ctx, "", "")

	sum := 0
	for _, v := range dtos2 {
		sum += v.(*testModel).Expect
	}

	if sum != sumExpect {
		t.Errorf("Test_queryRepositoryService_ListWithFilter() sum = %v, expect %v", sum, sumExpect)
	}

	if err != nil {
		t.Errorf("Test_queryRepositoryService_ListWithFilter() err = %v", err)
	}
}
