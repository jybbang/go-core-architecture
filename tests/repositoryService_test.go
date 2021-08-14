package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/mocks"
)

type testModel struct {
	core.Entity
	Expect int `bson:"expect,omitempty"`
}

func (e *testModel) CopyWith(src interface{}) bool {
	source, ok := src.(*testModel)
	e.ID = source.ID
	e.CreateUser = source.CreateUser
	e.UpdateUser = source.UpdateUser
	e.CreatedAt = source.CreatedAt
	e.UpdatedAt = source.UpdatedAt
	e.Expect = source.Expect
	return ok
}

func TestRepositoryService_Find(t *testing.T) {
	mock := mocks.NewMockAdapter()
	r := core.GetRepositoryService(new(testModel))
	r.SetCommandRepositoryAdapter(mock).SetQueryRepositoryAdapter(mock)

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	r.Add(dto)

	type args struct {
		dto core.Entitier
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		r       *core.RepositoryService
		args    args
		wantErr bool
	}{
		{
			name: "1",
			r:    r,
			args: args{
				dto: dto,
				id:  dto.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Find(tt.args.dto, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.Find() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryService_Any(t *testing.T) {
	tests := []struct {
		name    string
		r       *core.RepositoryService
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Any()
			if (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.Any() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RepositoryService.Any() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepositoryService_AnyWithFilter(t *testing.T) {
	type args struct {
		query interface{}
		args  interface{}
	}
	tests := []struct {
		name    string
		r       *core.RepositoryService
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.AnyWithFilter(tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.AnyWithFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RepositoryService.AnyWithFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepositoryService_Count(t *testing.T) {
	tests := []struct {
		name    string
		r       *core.RepositoryService
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Count()
			if (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RepositoryService.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepositoryService_CountWithFilter(t *testing.T) {
	type args struct {
		query interface{}
		args  interface{}
	}
	tests := []struct {
		name    string
		r       *core.RepositoryService
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.CountWithFilter(tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.CountWithFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RepositoryService.CountWithFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepositoryService_List(t *testing.T) {
	type args struct {
		dtos []core.Entitier
	}
	tests := []struct {
		name    string
		r       *core.RepositoryService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.List(tt.args.dtos); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.List() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryService_ListWithFilter(t *testing.T) {
	type args struct {
		dtos  []core.Entitier
		query interface{}
		args  interface{}
	}
	tests := []struct {
		name    string
		r       *core.RepositoryService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.ListWithFilter(tt.args.dtos, tt.args.query, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.ListWithFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryService_Remove(t *testing.T) {
	type args struct {
		entity core.Entitier
	}
	tests := []struct {
		name    string
		r       *core.RepositoryService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Remove(tt.args.entity); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryService_RemoveRange(t *testing.T) {
	type args struct {
		entities []core.Entitier
	}
	tests := []struct {
		name    string
		r       *core.RepositoryService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.RemoveRange(tt.args.entities); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.RemoveRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryService_Add(t *testing.T) {
	type args struct {
		entity core.Entitier
	}
	tests := []struct {
		name    string
		r       *core.RepositoryService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Add(tt.args.entity); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryService_AddRange(t *testing.T) {
	type args struct {
		entities []core.Entitier
	}
	tests := []struct {
		name    string
		r       *core.RepositoryService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.AddRange(tt.args.entities); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.AddRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryService_Update(t *testing.T) {
	type args struct {
		entity core.Entitier
	}
	tests := []struct {
		name    string
		r       *core.RepositoryService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Update(tt.args.entity); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryService_UpdateRange(t *testing.T) {
	type args struct {
		entities []core.Entitier
	}
	tests := []struct {
		name    string
		r       *core.RepositoryService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.UpdateRange(tt.args.entities); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.UpdateRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
