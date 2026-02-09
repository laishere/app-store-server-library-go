package internal

import "testing"

func TestParseTimestamp(t *testing.T) {
	if got := ParseTimestamp(float64(123.456)); got != 123 {
		t.Errorf("ParseTimestamp(123.456) = %v; want 123", got)
	}
	if got := ParseTimestamp(float64(123.999)); got != 123 {
		t.Errorf("ParseTimestamp(123.999) = %v; want 123", got)
	}
	if got := ParseTimestamp(int64(123)); got != 123 {
		t.Errorf("ParseTimestamp(int64(123)) = %v; want 123", got)
	}
	if got := ParseTimestamp("invalid"); got != 0 {
		t.Errorf("ParseTimestamp(\"invalid\") = %v; want 0", got)
	}
}
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

func TestUnmarshalTimestamp(t *testing.T) {
	var ts int64

	// Test valid timestamp
	UnmarshalTimestamp(float64(123.456), &ts)
	if ts != 123 {
		t.Errorf("Expected 123, got %v", ts)
	}

	// Test nil
	ts = 999
	UnmarshalTimestamp(nil, &ts)
	if ts != 999 {
		t.Errorf("Expected 999 (unchanged), got %v", ts)
	}
}
