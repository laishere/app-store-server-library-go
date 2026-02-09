package internal

// ParseTimestamp handles both integer and floating point timestamps
func ParseTimestamp(v any) int64 {
	switch val := v.(type) {
	case float64:
		return int64(val) // Truncate toward zero
	case int64:
		return val
	default:
		return 0
	}
}

// UnmarshalStringEnum unmarshals a string value into an enum type and its raw string representation.
func UnmarshalStringEnum[T ~string](data any, enum *T, raw *string) {
	if data == nil {
		return
	}
	if v, ok := data.(string); ok {
		*raw = v
		*enum = T(v)
	}
}

// UnmarshalIntEnum unmarshals an integer value into an enum type and its raw int32 representation.
func UnmarshalIntEnum[T ~int32](data any, enum *T, raw *int32) {
	if data == nil {
		return
	}
	if v, ok := data.(float64); ok {
		*raw = int32(v)
		*enum = T(int32(v))
	}
}

// UnmarshalTimestamp unmarshals a timestamp value (int64 or float64) into an int64.
func UnmarshalTimestamp(data any, ts *int64) {
	if data == nil {
		return
	}
	*ts = ParseTimestamp(data)
}
