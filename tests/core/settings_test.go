package core

import (
	"testing"
	"time"

	"github.com/jybbang/go-core-architecture/core"
)

func Test_settings_toGobreakerSettings(t *testing.T) {
	expect := 10
	expectName := "test"

	settings := &core.CircuitBreakerSettings{
		AllowedRequestInHalfOpen: expect,
	}

	cb := settings.ToGobreakerSettings(expectName)

	if cb.MaxRequests != uint32(expect) {
		t.Errorf("Test_CircuitBreakerSettings_toGobreakerSettings() MaxRequests = %v, expect %v", cb.MaxRequests, expect)
	}
	if cb.Name != expectName {
		t.Errorf("Test_CircuitBreakerSettings_toGobreakerSettings() Name = %v, expect %v", cb.Name, expectName)
	}
	if cb.Timeout != time.Duration(60*time.Second) {
		t.Errorf("Test_CircuitBreakerSettings_toGobreakerSettings() Timeout = %v, expect %v", cb.Timeout, 60)
	}
	if cb.Interval != time.Duration(60*time.Second) {
		t.Errorf("Test_CircuitBreakerSettings_toGobreakerSettings() Interval = %v, expect %v", cb.Interval, 60)
	}
}
