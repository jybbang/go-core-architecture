package core

import (
	"context"
	"testing"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/mocks"
)

func TestEventBus_PublishDomainEvents(t *testing.T) {
	core.NewMediatorBuilder().Build()

	e := core.NewEventbusBuilder().MessaingAdapter(mocks.NewMockAdapter()).Build()

	type args struct {
		ctx context.Context
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
			if err := e.PublishDomainEvents(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("EventBus.PublishDomainEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
