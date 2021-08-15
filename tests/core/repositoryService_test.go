package core

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/mocks"
)

type testModel struct {
	core.Entity
	Expect int `bson:"expect,omitempty"`
}

func TestRepositoryService_Find(t *testing.T) {
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mocks.NewMockAdapter()).
		QueryRepositoryAdapter(mocks.NewMockAdapter()).
		Build()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	dto2 := new(testModel)
	dto2.ID = uuid.New()
	dto2.Expect = 1234

	ctx := context.Background()

	r.Add(ctx, dto)
	r.Add(ctx, dto2)

	type args struct {
		ctx  context.Context
		dest core.Entitier
		id   uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				ctx:  context.Background(),
				dest: dto2,
				id:   dto.ID,
			},
			wantOk:  true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, err := r.Find(tt.args.ctx, tt.args.dest, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("RepositoryService.Find() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestRepositoryService_Any(t *testing.T) {
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mocks.NewMockAdapter()).
		QueryRepositoryAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, err := r.Any(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.Any() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("RepositoryService.Any() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestRepositoryService_AnyWithFilter(t *testing.T) {
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mocks.NewMockAdapter()).
		QueryRepositoryAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx   context.Context
		query interface{}
		args  interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, err := r.AnyWithFilter(tt.args.ctx, tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.AnyWithFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("RepositoryService.AnyWithFilter() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestRepositoryService_Count(t *testing.T) {
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mocks.NewMockAdapter()).
		QueryRepositoryAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		wantCount int64
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := r.Count(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("RepositoryService.Count() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func TestRepositoryService_CountWithFilter(t *testing.T) {
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mocks.NewMockAdapter()).
		QueryRepositoryAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx   context.Context
		query interface{}
		args  interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantCount int64
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := r.CountWithFilter(tt.args.ctx, tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.CountWithFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("RepositoryService.CountWithFilter() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func TestRepositoryService_List(t *testing.T) {
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mocks.NewMockAdapter()).
		QueryRepositoryAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx  context.Context
		dest []core.Entitier
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.List(tt.args.ctx, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.List() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryService_ListWithFilter(t *testing.T) {
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mocks.NewMockAdapter()).
		QueryRepositoryAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx   context.Context
		dest  []core.Entitier
		query interface{}
		args  interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.ListWithFilter(tt.args.ctx, tt.args.dest, tt.args.query, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.ListWithFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryService_Remove(t *testing.T) {
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mocks.NewMockAdapter()).
		QueryRepositoryAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx    context.Context
		entity core.Entitier
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Remove(tt.args.ctx, tt.args.entity); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryService_RemoveRange(t *testing.T) {
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mocks.NewMockAdapter()).
		QueryRepositoryAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx      context.Context
		entities []core.Entitier
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.RemoveRange(tt.args.ctx, tt.args.entities); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.RemoveRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryService_Add(t *testing.T) {
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mocks.NewMockAdapter()).
		QueryRepositoryAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx    context.Context
		entity core.Entitier
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Add(tt.args.ctx, tt.args.entity); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryService_AddRange(t *testing.T) {
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mocks.NewMockAdapter()).
		QueryRepositoryAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx      context.Context
		entities []core.Entitier
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.AddRange(tt.args.ctx, tt.args.entities); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.AddRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryService_Update(t *testing.T) {
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mocks.NewMockAdapter()).
		QueryRepositoryAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx    context.Context
		entity core.Entitier
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Update(tt.args.ctx, tt.args.entity); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryService_UpdateRange(t *testing.T) {
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mocks.NewMockAdapter()).
		QueryRepositoryAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx      context.Context
		entities []core.Entitier
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.UpdateRange(tt.args.ctx, tt.args.entities); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.UpdateRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
