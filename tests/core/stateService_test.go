package core

import (
	"context"
	"testing"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/mocks"
)

func TestStateService_Has(t *testing.T) {
	s := core.NewStateServiceBuilder().
		StateAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx context.Context
		key string
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
			gotOk, err := s.Has(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateService.Has() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("StateService.Has() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestStateService_Get(t *testing.T) {
	s := core.NewStateServiceBuilder().
		StateAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx  context.Context
		key  string
		dest interface{}
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
			gotOk, err := s.Get(tt.args.ctx, tt.args.key, tt.args.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("StateService.Get() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestStateService_Set(t *testing.T) {
	s := core.NewStateServiceBuilder().
		StateAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx   context.Context
		key   string
		value interface{}
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
			if err := s.Set(tt.args.ctx, tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("StateService.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStateService_Delete(t *testing.T) {
	s := core.NewStateServiceBuilder().
		StateAdapter(mocks.NewMockAdapter()).
		Build()

	type args struct {
		ctx context.Context
		key string
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
			if err := s.Delete(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("StateService.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
