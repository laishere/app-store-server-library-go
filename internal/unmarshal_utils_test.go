package internal

import (
	"testing"
)

func TestUnmarshalStringEnum(t *testing.T) {
	type MyEnum string
	var enum MyEnum
	var raw string
	validValues := []MyEnum{"VALUE"}

	// Test valid string
	UnmarshalStringEnum("VALUE", &enum, &raw, validValues)
	if enum != "VALUE" || raw != "VALUE" {
		t.Errorf("Expected VALUE, got enum=%v, raw=%v", enum, raw)
	}

	// Test invalid string
	enum = ""
	raw = ""
	UnmarshalStringEnum("INVALID", &enum, &raw, validValues)
	if enum != "" || raw != "INVALID" {
		t.Errorf("Expected empty enum and raw INVALID, got enum=%v, raw=%v", enum, raw)
	}

	// Test nil
	enum = ""
	raw = ""
	UnmarshalStringEnum(nil, &enum, &raw, validValues)
	if enum != "" || raw != "" {
		t.Errorf("Expected empty, got enum=%v, raw=%v", enum, raw)
	}

	// Test invalid type
	enum = ""
	raw = ""
	UnmarshalStringEnum(123, &enum, &raw, validValues)
	if enum != "" || raw != "" {
		t.Errorf("Expected empty for invalid type, got enum=%v, raw=%v", enum, raw)
	}
}

func TestUnmarshalIntEnum(t *testing.T) {
	type MyEnum int32
	var enum MyEnum
	var raw int32

	// Test valid int (from float64 as json unmarshal does)
	UnmarshalIntEnum(float64(123), &enum, &raw)
	if enum != 123 || raw != 123 {
		t.Errorf("Expected 123, got enum=%v, raw=%v", enum, raw)
	}

	// Test nil
	enum = 0
	raw = 0
	UnmarshalIntEnum(nil, &enum, &raw)
	if enum != 0 || raw != 0 {
		t.Errorf("Expected 0, got enum=%v, raw=%v", enum, raw)
	}

	// Test invalid type
	enum = 0
	raw = 0
	UnmarshalIntEnum("123", &enum, &raw)
	if enum != 0 || raw != 0 {
		t.Errorf("Expected 0 for invalid type, got enum=%v, raw=%v", enum, raw)
	}
}
