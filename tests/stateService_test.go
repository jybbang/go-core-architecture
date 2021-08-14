package tests

import (
	"testing"

	"github.com/jybbang/go-core-architecture/core"
)

func TestStateService_Has(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		s       *core.StateService
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Has(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateService.Has() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StateService.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateService_Get(t *testing.T) {
	type args struct {
		key  string
		dest core.Entitier
	}
	tests := []struct {
		name    string
		s       *core.StateService
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Get(tt.args.key, tt.args.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StateService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateService_Set(t *testing.T) {
	type args struct {
		key  string
		item interface{}
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
			if err := tt.s.Set(tt.args.key, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("StateService.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
