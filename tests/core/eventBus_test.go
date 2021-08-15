package core

import (
	"context"
	"testing"

	"github.com/jybbang/go-core-architecture/core"
)

func TestEventBus_AddDomainEvent(t *testing.T) {
	type args struct {
		domainEvent core.DomainEventer
	}
	tests := []struct {
		name string
		e    *core.EventBus
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.AddDomainEvent(tt.args.domainEvent)
		})
	}
}

func TestEventBus_PublishDomainEvents(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		e       *core.EventBus
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.PublishDomainEvents(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("EventBus.PublishDomainEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}