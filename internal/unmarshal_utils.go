package internal

import (
	"slices"
)

// UnmarshalStringEnum unmarshals a string value into an enum type and its raw string representation.
func UnmarshalStringEnum[T ~string](data any, enum *T, raw *string, validValues []T) {
	if data == nil {
		return
	}
	if v, ok := data.(string); ok {
		*raw = v
		val := T(v)
		if slices.Contains(validValues, val) {
			*enum = val
			return
		}
		*enum = ""
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
