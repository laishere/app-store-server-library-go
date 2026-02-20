package appstore

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	assert := assert.New(t)
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
				assert.Error(err, "Timestamp.UnmarshalJSON()")
				return
			}
			assert.NoError(err, "Timestamp.UnmarshalJSON()")
			assert.Equal(tt.want, ts, "Timestamp.UnmarshalJSON()")
		})
	}
}

func TestTimestamp_Time(t *testing.T) {
	assert := assert.New(t)
	ts := Timestamp(1698148900000)
	want := time.UnixMilli(1698148900000)
	assert.Equal(want, ts.Time(), "Timestamp.Time()")

	var zeroTs Timestamp
	assert.Equal(true, zeroTs.Time().IsZero(), "Timestamp.Time() for zero")
}

func TestTimestamp_TimePtr(t *testing.T) {
	assert := assert.New(t)
	ts := Timestamp(1698148900000)
	want := time.UnixMilli(1698148900000)
	assert.Equal(want, *ts.TimePtr(), "Timestamp.TimePtr()")

	nilTs := (*Timestamp)(nil)
	assert.Nil(nilTs.TimePtr(), "Timestamp.TimePtr() for nil")
}

func TestTimestamp_UnixMilli(t *testing.T) {
	assert := assert.New(t)
	ts := Timestamp(1698148900000)
	assert.Equal(int64(1698148900000), ts.UnixMilli(), "Timestamp.UnixMilli()")

	var zeroTs Timestamp
	assert.Equal(int64(0), zeroTs.UnixMilli(), "Timestamp.UnixMilli() for zero")
}

func TestTimestamp_IsZero(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(true, Timestamp(0).IsZero(), "Timestamp(0).IsZero()")
	assert.Equal(false, Timestamp(123).IsZero(), "Timestamp(123).IsZero()")
}

func TestTimestamp_String(t *testing.T) {
	assert := assert.New(t)
	ts := Timestamp(1698148900000)
	assert.Equal(ts.Time().Format(time.DateTime), ts.String(), "Timestamp.String()")

	ts = Timestamp(0)
	assert.Equal("N/A", ts.String(), "Timestamp(0).String()")
}

func TestTimestamp_UnmarshalJSON_DirectError(t *testing.T) {
	assert := assert.New(t)
	var ts Timestamp
	err := ts.UnmarshalJSON([]byte("{"))
	assert.Error(err, "Timestamp.UnmarshalJSON() for invalid json")
}
