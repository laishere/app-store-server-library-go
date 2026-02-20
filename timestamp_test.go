package appstore

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Timestamp
		wantErr bool
	}{
		{
			name:    "integer timestamp",
			input:   "1698148900000",
			want:    Timestamp(1698148900000),
			wantErr: false,
		},
		{
			name:    "float timestamp (Xcode style)",
			input:   "1697680122257.447",
			want:    Timestamp(1697680122257),
			wantErr: false,
		},
		{
			name:    "zero timestamp",
			input:   "0",
			want:    Timestamp(0),
			wantErr: false,
		},
		{
			name:    "null timestamp",
			input:   "null",
			want:    Timestamp(0),
			wantErr: false,
		},
		{
			name:    "invalid timestamp",
			input:   `"invalid"`,
			want:    Timestamp(0),
			wantErr: false,
		},
		{
			name:    "boolean timestamp",
			input:   "true",
			want:    Timestamp(0),
			wantErr: false,
		},
		{
			name:    "array timestamp",
			input:   "[]",
			want:    Timestamp(0),
			wantErr: false,
		},
		{
			name:    "object timestamp",
			input:   "{}",
			want:    Timestamp(0),
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   "{",
			want:    Timestamp(0),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts Timestamp
			err := json.Unmarshal([]byte(tt.input), &ts)
			if tt.wantErr {
				assertError(t, err, "Timestamp.UnmarshalJSON()")
				return
			}
			assertNoError(t, err, "Timestamp.UnmarshalJSON()")
			assertEqual(t, tt.want, ts, "Timestamp.UnmarshalJSON()")
		})
	}
}

func TestTimestamp_Time(t *testing.T) {
	ts := Timestamp(1698148900000)
	want := time.UnixMilli(1698148900000)
	assertEqual(t, want, ts.Time(), "Timestamp.Time()")

	var zeroTs Timestamp
	assertEqual(t, true, zeroTs.Time().IsZero(), "Timestamp.Time() for zero")
}

func TestTimestamp_TimePtr(t *testing.T) {
	ts := Timestamp(1698148900000)
	want := time.UnixMilli(1698148900000)
	assertEqual(t, want, *ts.TimePtr(), "Timestamp.TimePtr()")

	nilTs := (*Timestamp)(nil)
	assertNil(t, nilTs.TimePtr(), "Timestamp.TimePtr() for nil")
}

func TestTimestamp_UnixMilli(t *testing.T) {
	ts := Timestamp(1698148900000)
	assertEqual(t, int64(1698148900000), ts.UnixMilli(), "Timestamp.UnixMilli()")

	var zeroTs Timestamp
	assertEqual(t, int64(0), zeroTs.UnixMilli(), "Timestamp.UnixMilli() for zero")
}

func TestTimestamp_IsZero(t *testing.T) {
	assertEqual(t, true, Timestamp(0).IsZero(), "Timestamp(0).IsZero()")
	assertEqual(t, false, Timestamp(123).IsZero(), "Timestamp(123).IsZero()")
}

func TestTimestamp_String(t *testing.T) {
	ts := Timestamp(1698148900000)
	assertEqual(t, ts.Time().Format(time.DateTime), ts.String(), "Timestamp.String()")

	ts = Timestamp(0)
	assertEqual(t, "N/A", ts.String(), "Timestamp(0).String()")
}

func TestTimestamp_UnmarshalJSON_DirectError(t *testing.T) {
	var ts Timestamp
	err := ts.UnmarshalJSON([]byte("{"))
	assertError(t, err, "Timestamp.UnmarshalJSON() for invalid json")
}
