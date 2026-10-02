// SPDX-License-Identifier: BUSL-1.1

package codec8

import (
	"bytes"
	"errors"
	"maps"
	"testing"
	"time"
)

// specData returns the data field of the spec example: the frame without its
// 8-byte header and 4-byte CRC, which is what Decode receives from ReadFrame.
func specData(t *testing.T) []byte {
	t.Helper()
	frame := mustHex(t, specFrame)
	return frame[8 : len(frame)-4]
}

// withNX replaces the empty NX group and Number of Data 2 at the end of data
// (its last 3 bytes) with one NX element, ID 255, whose length field says
// length while the value is always aa bb cc, then Number of Data 2 again.
func withNX(data []byte, length byte) []byte {
	return append(bytes.Clone(data[:len(data)-3]),
		0x00, 0x01, // NX count: one element
		0x00, 0xFF, // IO ID 255
		0x00, length, // length field
		0xAA, 0xBB, 0xCC, // value: 3 bytes
		0x01, // Number of Data 2
	)
}

func TestDecodeSpecExample(t *testing.T) {
	records, err := Decode(specData(t))
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("got %d records, want 1", len(records))
	}
	r := records[0]

	if want := time.Date(2019, time.June, 10, 11, 36, 32, 0, time.UTC); !r.Timestamp.Equal(want) {
		t.Errorf("Timestamp = %v, want %v", r.Timestamp, want)
	}
	if r.Priority != PriorityHigh {
		t.Errorf("Priority = %d, want %d", r.Priority, PriorityHigh)
	}
	if r.GPS != (GPS{}) {
		t.Errorf("GPS = %+v, want all zero: the example has no fix", r.GPS)
	}
	if r.IO.EventID != 1 {
		t.Errorf("EventID = %d, want 1", r.IO.EventID)
	}
	wantValues := map[uint16]uint64{
		1:  1,         // DIN1, 1-byte group
		17: 29,        // Axis X, 2-byte group
		16: 22949000,  // Total Odometer, 4-byte group
		11: 893700218, // ICCID1, 8-byte group
		14: 500686954, // ICCID2, 8-byte group
	}
	if !maps.Equal(r.IO.Values, wantValues) {
		t.Errorf("Values = %v, want %v", r.IO.Values, wantValues)
	}
	if len(r.IO.RawValues) != 0 {
		t.Errorf("RawValues = %v, want none", r.IO.RawValues)
	}
}

func TestDecodeVariableSizeIO(t *testing.T) {
	records, err := Decode(withNX(specData(t), 0x03))
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if got := records[0].IO.RawValues[255]; !bytes.Equal(got, []byte{0xAA, 0xBB, 0xCC}) {
		t.Errorf("RawValues[255] = % x, want aa bb cc", got)
	}
}

func TestDecodeNoRecords(t *testing.T) {
	records, err := Decode([]byte{codecID8E, 0x00, 0x00})
	if err != nil || len(records) != 0 {
		t.Fatalf("Decode = %v, %v; want no records and no error", records, err)
	}
}

func TestDecodeRejects(t *testing.T) {
	good := specData(t)

	tests := []struct {
		name    string
		data    []byte
		wantErr error
	}{
		{"empty", nil, ErrMalformedData},
		{"Codec 8 instead of 8E", append([]byte{0x08}, good[1:]...), ErrInvalidCodecID},
		{"cut in the middle of a record", good[:40], ErrMalformedData},
		{"NX value shorter than its length", withNX(good, 0x05), ErrMalformedData},
		{"number of data 2 differs", corrupt(good, len(good)-1), ErrMalformedData},
		{"trailing byte", append(bytes.Clone(good), 0x00), ErrMalformedData},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			records, err := Decode(tt.data)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if records != nil {
				t.Errorf("records = %v, want nil on error", records)
			}
		})
	}
}
