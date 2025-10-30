package togglr

import (
	"testing"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext()

	if ctx == nil {
		t.Fatal("NewContext() returned nil")
	}

	if len(ctx) != 0 {
		t.Errorf("Expected empty context, got %d items", len(ctx))
	}
}

func TestRequestContext_WithUserID(t *testing.T) {
	ctx := NewContext()
	userID := "user123"

	result := ctx.WithUserID(userID)

	if result[AttrUserID] != userID {
		t.Errorf("Expected user ID '%s', got %v", userID, result[AttrUserID])
	}

	// Verify it's the same map instance (maps are reference types)
	if len(ctx) != len(result) {
		t.Error("WithUserID should modify and return the same map instance")
	}
}

func TestRequestContext_WithUserEmail(t *testing.T) {
	ctx := NewContext()
	email := "user@example.com"

	result := ctx.WithUserEmail(email)

	if result[AttrUserEmail] != email {
		t.Errorf("Expected user email '%s', got %v", email, result[AttrUserEmail])
	}

	// Verify the original context was modified
	if ctx[AttrUserEmail] != email {
		t.Error("WithUserEmail should modify the original context")
	}
}

func TestRequestContext_WithAnonymous(t *testing.T) {
	tests := []struct {
		name  string
		value bool
	}{
		{"anonymous true", true},
		{"anonymous false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext()
			result := ctx.WithAnonymous(tt.value)

			if result[AttrUserAnonymous] != tt.value {
				t.Errorf("Expected anonymous '%v', got %v", tt.value, result[AttrUserAnonymous])
			}

			if ctx[AttrUserAnonymous] != tt.value {
				t.Error("WithAnonymous should modify the original context")
			}
		})
	}
}

func TestRequestContext_WithCountry(t *testing.T) {
	ctx := NewContext()
	country := "US"

	result := ctx.WithCountry(country)

	if result[AttrCountryCode] != country {
		t.Errorf("Expected country code '%s', got %v", country, result[AttrCountryCode])
	}

	if ctx[AttrCountryCode] != country {
		t.Error("WithCountry should modify the original context")
	}
}

func TestRequestContext_WithRegion(t *testing.T) {
	ctx := NewContext()
	region := "California"

	result := ctx.WithRegion(region)

	if result[AttrRegion] != region {
		t.Errorf("Expected region '%s', got %v", region, result[AttrRegion])
	}

	if ctx[AttrRegion] != region {
		t.Error("WithRegion should modify the original context")
	}
}

func TestRequestContext_WithCity(t *testing.T) {
	ctx := NewContext()
	city := "San Francisco"

	result := ctx.WithCity(city)

	if result[AttrCity] != city {
		t.Errorf("Expected city '%s', got %v", city, result[AttrCity])
	}

	if ctx[AttrCity] != city {
		t.Error("WithCity should modify the original context")
	}
}

func TestRequestContext_WithManufacturer(t *testing.T) {
	ctx := NewContext()
	manufacturer := "Apple"

	result := ctx.WithManufacturer(manufacturer)

	if result[AttrManufacturer] != manufacturer {
		t.Errorf("Expected manufacturer '%s', got %v", manufacturer, result[AttrManufacturer])
	}

	if ctx[AttrManufacturer] != manufacturer {
		t.Error("WithManufacturer should modify the original context")
	}
}

func TestRequestContext_WithDeviceType(t *testing.T) {
	ctx := NewContext()
	deviceType := "mobile"

	result := ctx.WithDeviceType(deviceType)

	if result[AttrDeviceType] != deviceType {
		t.Errorf("Expected device type '%s', got %v", deviceType, result[AttrDeviceType])
	}

	if ctx[AttrDeviceType] != deviceType {
		t.Error("WithDeviceType should modify the original context")
	}
}

func TestRequestContext_WithOS(t *testing.T) {
	ctx := NewContext()
	os := "iOS"

	result := ctx.WithOS(os)

	if result[AttrOS] != os {
		t.Errorf("Expected OS '%s', got %v", os, result[AttrOS])
	}

	if ctx[AttrOS] != os {
		t.Error("WithOS should modify the original context")
	}
}

func TestRequestContext_WithOSVersion(t *testing.T) {
	ctx := NewContext()
	version := "15.0"

	result := ctx.WithOSVersion(version)

	if result[AttrOSVersion] != version {
		t.Errorf("Expected OS version '%s', got %v", version, result[AttrOSVersion])
	}

	if ctx[AttrOSVersion] != version {
		t.Error("WithOSVersion should modify the original context")
	}
}

func TestRequestContext_WithBrowser(t *testing.T) {
	ctx := NewContext()
	browser := "Chrome"

	result := ctx.WithBrowser(browser)

	if result[AttrBrowser] != browser {
		t.Errorf("Expected browser '%s', got %v", browser, result[AttrBrowser])
	}

	if ctx[AttrBrowser] != browser {
		t.Error("WithBrowser should modify the original context")
	}
}

func TestRequestContext_WithBrowserVersion(t *testing.T) {
	ctx := NewContext()
	version := "96.0.4664.110"

	result := ctx.WithBrowserVersion(version)

	if result[AttrBrowserVersion] != version {
		t.Errorf("Expected browser version '%s', got %v", version, result[AttrBrowserVersion])
	}

	if ctx[AttrBrowserVersion] != version {
		t.Error("WithBrowserVersion should modify the original context")
	}
}

func TestRequestContext_WithLanguage(t *testing.T) {
	ctx := NewContext()
	language := "en-US"

	result := ctx.WithLanguage(language)

	if result[AttrLanguage] != language {
		t.Errorf("Expected language '%s', got %v", language, result[AttrLanguage])
	}

	if ctx[AttrLanguage] != language {
		t.Error("WithLanguage should modify the original context")
	}
}

func TestRequestContext_WithConnectionType(t *testing.T) {
	ctx := NewContext()
	connType := "wifi"

	result := ctx.WithConnectionType(connType)

	if result[AttrConnectionType] != connType {
		t.Errorf("Expected connection type '%s', got %v", connType, result[AttrConnectionType])
	}

	if ctx[AttrConnectionType] != connType {
		t.Error("WithConnectionType should modify the original context")
	}
}

func TestRequestContext_WithAge(t *testing.T) {
	ctx := NewContext()
	age := 25

	result := ctx.WithAge(age)

	if result[AttrAge] != age {
		t.Errorf("Expected age '%d', got %v", age, result[AttrAge])
	}

	if ctx[AttrAge] != age {
		t.Error("WithAge should modify the original context")
	}
}

func TestRequestContext_WithGender(t *testing.T) {
	ctx := NewContext()
	gender := "male"

	result := ctx.WithGender(gender)

	if result[AttrGender] != gender {
		t.Errorf("Expected gender '%s', got %v", gender, result[AttrGender])
	}

	if ctx[AttrGender] != gender {
		t.Error("WithGender should modify the original context")
	}
}

func TestRequestContext_WithIP(t *testing.T) {
	ctx := NewContext()
	ip := "192.168.1.1"

	result := ctx.WithIP(ip)

	if result[AttrIP] != ip {
		t.Errorf("Expected IP '%s', got %v", ip, result[AttrIP])
	}

	if ctx[AttrIP] != ip {
		t.Error("WithIP should modify the original context")
	}
}

func TestRequestContext_WithAppVersion(t *testing.T) {
	ctx := NewContext()
	version := "1.2.3"

	result := ctx.WithAppVersion(version)

	if result[AttrAppVersion] != version {
		t.Errorf("Expected app version '%s', got %v", version, result[AttrAppVersion])
	}

	if ctx[AttrAppVersion] != version {
		t.Error("WithAppVersion should modify the original context")
	}
}

func TestRequestContext_WithPlatform(t *testing.T) {
	ctx := NewContext()
	platform := "web"

	result := ctx.WithPlatform(platform)

	if result[AttrPlatform] != platform {
		t.Errorf("Expected platform '%s', got %v", platform, result[AttrPlatform])
	}

	if ctx[AttrPlatform] != platform {
		t.Error("WithPlatform should modify the original context")
	}
}

func TestRequestContext_Set(t *testing.T) {
	ctx := NewContext()
	key := "custom_key"
	value := "custom_value"

	result := ctx.Set(key, value)

	if result[key] != value {
		t.Errorf("Expected custom value '%s', got %v", value, result[key])
	}

	if ctx[key] != value {
		t.Error("Set should modify the original context")
	}
}

func TestRequestContext_SetWithDifferentTypes(t *testing.T) {
	t.Run("string value", func(t *testing.T) {
		ctx := NewContext()
		ctx.Set("str_key", "string_value")
		if ctx["str_key"] != "string_value" {
			t.Errorf("Expected 'string_value', got %v", ctx["str_key"])
		}
	})

	t.Run("int value", func(t *testing.T) {
		ctx := NewContext()
		ctx.Set("int_key", 42)
		if ctx["int_key"] != 42 {
			t.Errorf("Expected 42, got %v", ctx["int_key"])
		}
	})

	t.Run("bool value", func(t *testing.T) {
		ctx := NewContext()
		ctx.Set("bool_key", true)
		if ctx["bool_key"] != true {
			t.Errorf("Expected true, got %v", ctx["bool_key"])
		}
	})

	t.Run("float value", func(t *testing.T) {
		ctx := NewContext()
		ctx.Set("float_key", 3.14)
		if ctx["float_key"] != 3.14 {
			t.Errorf("Expected 3.14, got %v", ctx["float_key"])
		}
	})

	t.Run("slice value", func(t *testing.T) {
		ctx := NewContext()
		slice := []string{"a", "b", "c"}
		ctx.Set("slice_key", slice)
		if ctx["slice_key"] == nil {
			t.Error("Expected slice to be set")
		}
	})

	t.Run("map value", func(t *testing.T) {
		ctx := NewContext()
		m := map[string]int{"count": 10}
		ctx.Set("map_key", m)
		if ctx["map_key"] == nil {
			t.Error("Expected map to be set")
		}
	})
}

func TestRequestContext_MethodChaining(t *testing.T) {
	ctx := NewContext()

	result := ctx.
		WithUserID("user123").
		WithUserEmail("user@example.com").
		WithAnonymous(false).
		WithCountry("US").
		WithRegion("California").
		WithCity("San Francisco").
		WithManufacturer("Apple").
		WithDeviceType("mobile").
		WithOS("iOS").
		WithOSVersion("15.0").
		WithBrowser("Safari").
		WithBrowserVersion("15.0").
		WithLanguage("en-US").
		WithConnectionType("wifi").
		WithAge(25).
		WithGender("male").
		WithIP("192.168.1.1").
		WithAppVersion("1.2.3").
		WithPlatform("ios").
		Set("custom", "value")

	// Verify all values are set correctly
	expectedValues := map[string]any{
		AttrUserID:         "user123",
		AttrUserEmail:      "user@example.com",
		AttrUserAnonymous:  false,
		AttrCountryCode:    "US",
		AttrRegion:         "California",
		AttrCity:           "San Francisco",
		AttrManufacturer:   "Apple",
		AttrDeviceType:     "mobile",
		AttrOS:             "iOS",
		AttrOSVersion:      "15.0",
		AttrBrowser:        "Safari",
		AttrBrowserVersion: "15.0",
		AttrLanguage:       "en-US",
		AttrConnectionType: "wifi",
		AttrAge:            25,
		AttrGender:         "male",
		AttrIP:             "192.168.1.1",
		AttrAppVersion:     "1.2.3",
		AttrPlatform:       "ios",
		"custom":           "value",
	}

	if len(result) != len(expectedValues) {
		t.Errorf("Expected %d items in context, got %d", len(expectedValues), len(result))
	}

	for key, expectedValue := range expectedValues {
		if result[key] != expectedValue {
			t.Errorf("For key '%s': expected '%v', got '%v'", key, expectedValue, result[key])
		}
	}
}

func TestRequestContext_OverwriteValues(t *testing.T) {
	ctx := NewContext()

	// Set initial value
	ctx.WithUserID("user1")
	if ctx[AttrUserID] != "user1" {
		t.Errorf("Expected initial user ID 'user1', got %v", ctx[AttrUserID])
	}

	// Overwrite with new value
	ctx.WithUserID("user2")
	if ctx[AttrUserID] != "user2" {
		t.Errorf("Expected overwritten user ID 'user2', got %v", ctx[AttrUserID])
	}
}

func TestRequestContext_EmptyValues(t *testing.T) {
	ctx := NewContext()

	// Test setting empty strings
	ctx.WithUserID("").
		WithUserEmail("").
		WithCountry("")

	if ctx[AttrUserID] != "" {
		t.Errorf("Expected empty user ID, got %v", ctx[AttrUserID])
	}

	if ctx[AttrUserEmail] != "" {
		t.Errorf("Expected empty user email, got %v", ctx[AttrUserEmail])
	}

	if ctx[AttrCountryCode] != "" {
		t.Errorf("Expected empty country code, got %v", ctx[AttrCountryCode])
	}
}

func TestRequestContext_ZeroValues(t *testing.T) {
	ctx := NewContext()

	// Test setting zero values
	ctx.WithAge(0).
		WithAnonymous(false)

	if ctx[AttrAge] != 0 {
		t.Errorf("Expected age 0, got %v", ctx[AttrAge])
	}

	if ctx[AttrUserAnonymous] != false {
		t.Errorf("Expected anonymous false, got %v", ctx[AttrUserAnonymous])
	}
}

func TestRequestContext_NilValue(t *testing.T) {
	ctx := NewContext()

	ctx.Set("nil_key", nil)

	if ctx["nil_key"] != nil {
		t.Errorf("Expected nil value, got %v", ctx["nil_key"])
	}
}

func TestRequestContext_MultipleContexts(t *testing.T) {
	ctx1 := NewContext().WithUserID("user1")
	ctx2 := NewContext().WithUserID("user2")

	if ctx1[AttrUserID] == ctx2[AttrUserID] {
		t.Error("Different contexts should have independent values")
	}

	if ctx1[AttrUserID] != "user1" {
		t.Errorf("Expected ctx1 user ID 'user1', got %v", ctx1[AttrUserID])
	}

	if ctx2[AttrUserID] != "user2" {
		t.Errorf("Expected ctx2 user ID 'user2', got %v", ctx2[AttrUserID])
	}
}
