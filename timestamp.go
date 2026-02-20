package appstore

import (
	"encoding/json"
	"time"
)

// Timestamp represents a UNIX time in milliseconds.
type Timestamp int64

// Time returns the time.Time representation of the timestamp.
func (t Timestamp) Time() time.Time {
	if t == 0 {
		return time.Time{}
	}
	return time.UnixMilli(int64(t))
}

// TimePtr returns the time.Time representation of the timestamp as a pointer.
func (t *Timestamp) TimePtr() *time.Time {
	if t == nil {
		return nil
	}
	t2 := t.Time()
	return &t2
}

// UnixMilli returns the timestamp in milliseconds.
func (t Timestamp) UnixMilli() int64 {
	return int64(t)
}

// IsZero returns true if the timestamp is 0.
func (t Timestamp) IsZero() bool {
	return t == 0
}

// String returns the string representation of the timestamp.
func (t Timestamp) String() string {
	if t == 0 {
		return "N/A"
	}
	return t.Time().Format(time.DateTime)
}

// UnmarshalJSON unmarshals a JSON value into a Timestamp.
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	if val, ok := v.(float64); ok {
		*t = Timestamp(int64(val))
	} else {
		*t = Timestamp(0)
	}
	return nil
}
