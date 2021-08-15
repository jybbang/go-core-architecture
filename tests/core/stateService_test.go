package core

import (
	"context"
	"reflect"
	"testing"

	"github.com/jybbang/go-core-architecture/core"
)

func TestStateService_Setup(t *testing.T) {
	tests := []struct {
		name string
		s    *core.StateService
		want *core.StateService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Initialize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StateService.Setup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateService_SetStateAdapter(t *testing.T) {
	type args struct {
		adapter core.StateAdapter
	}
	tests := []struct {
		name string
		s    *core.StateService
		args args
		want *core.StateService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.SetStateAdapter(tt.args.adapter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StateService.SetStateAdapter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateService_Has(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		s       *core.StateService
		args    args
		wantOk  bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, err := tt.s.Has(tt.args.ctx, tt.args.key)
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
	type args struct {
		ctx  context.Context
		key  string
		dest interface{}
	}
	tests := []struct {
		name    string
		s       *core.StateService
		args    args
		wantOk  bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, err := tt.s.Get(tt.args.ctx, tt.args.key, tt.args.dest)
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
	type args struct {
		ctx   context.Context
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		s       *core.StateService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Set(tt.args.ctx, tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("StateService.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
