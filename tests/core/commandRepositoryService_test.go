package core

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/mocks"
)

func Test_commandRepositoryService_Remove(t *testing.T) {
	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel), "testModel").
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

	time.Sleep(1 * time.Second)

	result := r.Remove(ctx, dto)

	if result.E != nil {
		t.Errorf("Test_commandRepositoryService_Remove() err = %v", result.E)
	}

	dto2 := new(testModel)
	result = r.Find(ctx, dto.ID, dto2)

	if !errors.Is(result.E, core.ErrNotFound) {
		t.Errorf("Test_commandRepositoryService_Remove() err = %v, expect %v", result.E, core.ErrNotFound)
	}
}

func Test_commandRepositoryService_RemoveRange(t *testing.T) {
	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel), "testModel").
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

	result := r.RemoveRange(ctx, dtos)

	if result.E != nil {
		t.Errorf("Test_commandRepositoryService_RemoveRange() err = %v", result.E)
	}

	dto2 := new(testModel)
	result = r.Find(ctx, dtos[0].GetID(), dto2)

	if !errors.Is(result.E, core.ErrNotFound) {
		t.Errorf("Test_commandRepositoryService_RemoveRange() err = %v, expect %v", result.E, core.ErrNotFound)
	}
}

func Test_commandRepositoryService_AddRange(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel), "testModel").
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	expect := 10000
	var cntExpect int64 = 0
	rand.Seed(time.Now().UnixNano())
	random := rand.Int()
	var dtos = make([]core.Entitier, 0)
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = random

		dtos = append(dtos, dto)
		cntExpect += 1
	}

	result := r.AddRange(ctx, dtos)

	if result.E != nil {
		t.Errorf("Test_commandRepositoryService_AddRange() err = %v", result.E)
	}

	result = r.CountWithFilter(ctx, "", "")

	if result.V.(int64) != cntExpect {
		t.Errorf("Test_commandRepositoryService_AddRange() cnt = %v, expect %v", result.V, cntExpect)
	}
}

func Test_commandRepositoryService_Update(t *testing.T) {
	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel), "testModel").
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

	time.Sleep(1 * time.Second)

	result := r.Update(ctx, dto)

	if result.E != nil {
		t.Errorf("Test_commandRepositoryService_Update() err = %v", result.E)
	}

	dto2 := new(testModel)
	result = r.Find(ctx, dto.GetID(), dto2)

	if dto2.Expect != 1 {
		t.Errorf("Test_commandRepositoryService_Update() result = %v, expect %v", dto2.Expect, 1)
	}
}

func Test_commandRepositoryService_UpdateRange(t *testing.T) {
	ctx := context.Background()

	mock := mocks.NewMockAdapter()
	r := core.NewRepositoryServiceBuilder(new(testModel), "testModel").
		CommandRepositoryAdapter(mock).
		QueryRepositoryAdapter(mock).
		Create()

	expect := 10000
	var dtos = make([]core.Entitier, 0)
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = i

		dtos = append(dtos, dto)
		r.Add(ctx, dto)
	}

	var cntExpect int64 = 0
	rand.Seed(time.Now().UnixNano())
	random := rand.Int()
	for _, dto := range dtos {
		dto.(*testModel).Expect = random
		cntExpect += 1
	}

	result := r.UpdateRange(ctx, dtos)

	if result.E != nil {
		t.Errorf("Test_commandRepositoryService_UpdateRange() err = %v", result.E)
	}

	result = r.CountWithFilter(ctx, "", "")

	if result.V.(int64) != cntExpect {
		t.Errorf("Test_commandRepositoryService_UpdateRange() cnt = %v, expect %v", result.V, cntExpect)
	}
}
