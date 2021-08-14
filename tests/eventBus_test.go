package tests

import (
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
	tests := []struct {
		name    string
		e       *core.EventBus
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.PublishDomainEvents(); (err != nil) != tt.wantErr {
				t.Errorf("EventBus.PublishDomainEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
