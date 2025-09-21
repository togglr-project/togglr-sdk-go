package togglr

import (
	"testing"
)

func TestRequestContext(t *testing.T) {
	ctx := NewContext()

	// Test chaining
	ctx = ctx.WithUserID("user123").
		WithCountry("US").
		WithDeviceType("mobile").
		WithOS("iOS")

	// Test values
	if ctx[AttrUserID] != "user123" {
		t.Errorf("Expected user ID 'user123', got %v", ctx[AttrUserID])
	}

	if ctx[AttrCountryCode] != "US" {
		t.Errorf("Expected country 'US', got %v", ctx[AttrCountryCode])
	}

	if ctx[AttrDeviceType] != "mobile" {
		t.Errorf("Expected device type 'mobile', got %v", ctx[AttrDeviceType])
	}

	if ctx[AttrOS] != "iOS" {
		t.Errorf("Expected OS 'iOS', got %v", ctx[AttrOS])
	}
}

func TestRequestContextSet(t *testing.T) {
	ctx := NewContext()

	ctx = ctx.Set("custom_key", "custom_value")

	if ctx["custom_key"] != "custom_value" {
		t.Errorf("Expected custom value 'custom_value', got %v", ctx["custom_key"])
	}
}
