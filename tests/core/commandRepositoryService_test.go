package core

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/mocks"
)

func Test_commandRepositoryService_Remove(t *testing.T) {
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
		go r.Remove(ctx, dto)
	}

	err := r.Remove(ctx, dto)

	if err != nil {
		t.Errorf("Test_commandRepositoryService_Remove() err = %v", err)
	}

	dto2 := new(testModel)
	err = r.Find(ctx, dto.ID, dto2)

	if err != core.ErrNotFound {
		t.Errorf("Test_commandRepositoryService_Remove() err = %v, expect %v", err, core.ErrNotFound)
	}
}

func Test_commandRepositoryService_RemoveRange(t *testing.T) {
	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	ctx := context.Background()
	expect := 10
	var dtos = make([]core.Entitier, 0)
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = i
		dtos = append(dtos, dto)

		r.Add(ctx, dto)
	}

	err := r.RemoveRange(ctx, dtos)

	if err != nil {
		t.Errorf("Test_commandRepositoryService_RemoveRange() err = %v", err)
	}

	dto2 := new(testModel)
	err = r.Find(ctx, dtos[0].GetID(), dto2)

	if err != core.ErrNotFound {
		t.Errorf("Test_commandRepositoryService_RemoveRange() err = %v, expect %v", err, core.ErrNotFound)
	}
}

func Test_commandRepositoryService_AddRange(t *testing.T) {
	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	ctx := context.Background()
	expect := 10000
	sumExpect := 0
	var dtos = make([]core.Entitier, 0)
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = i
		dtos = append(dtos, dto)
		sumExpect += i
	}

	err := r.AddRange(ctx, dtos)

	if err != nil {
		t.Errorf("Test_commandRepositoryService_AddRange() err = %v", err)
	}

	dtos2, err := r.List(ctx)

	sum := 0
	for _, v := range dtos2 {
		sum += v.(*testModel).Expect
	}

	if sum != sumExpect {
		t.Errorf("Test_commandRepositoryService_AddRange() sum = %v, expect %v", sum, sumExpect)
	}

	if err != nil {
		t.Errorf("Test_commandRepositoryService_AddRange() err = %v", err)
	}
}

func Test_commandRepositoryService_Update(t *testing.T) {
	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	ctx := context.Background()
	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 100

	r.Add(ctx, dto)

	dto.Expect = 1
	count := 10000
	for i := 0; i < count; i++ {
		go r.Update(ctx, dto)
	}

	err := r.Update(ctx, dto)

	if err != nil {
		t.Errorf("Test_commandRepositoryService_Update() err = %v", err)
	}

	dto2 := new(testModel)
	r.Find(ctx, dto.GetID(), dto2)

	result := dto2.Expect
	if result != 1 {
		t.Errorf("Test_commandRepositoryService_Update() result = %v, expect %v", result, 1)
	}
}

func Test_commandRepositoryService_UpdateRange(t *testing.T) {
	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	ctx := context.Background()
	expect := 10000
	var dtos = make([]core.Entitier, 0)
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = 0
		dtos = append(dtos, dto)
		r.Add(ctx, dto)
	}

	sumExpect := 0
	for _, dto := range dtos {
		dto.(*testModel).Expect = 1
		sumExpect += 1
	}

	err := r.UpdateRange(ctx, dtos)

	if err != nil {
		t.Errorf("Test_commandRepositoryService_UpdateRange() err = %v", err)
	}

	dtos2, err := r.List(ctx)

	sum := 0
	for _, v := range dtos2 {
		sum += v.(*testModel).Expect
	}

	if sum != sumExpect {
		t.Errorf("Test_commandRepositoryService_UpdateRange() sum = %v, expect %v", sum, sumExpect)
	}
}
