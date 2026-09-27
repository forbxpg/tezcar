// SPDX-License-Identifier: BUSL-1.1

package codec8

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"testing"
	"testing/iotest"
)

// specFrame is the Codec 8 Extended example from the Teltonika wiki:
// one record, data field length 0x4A, CRC 0x2994.
const specFrame = "000000000000004A8E010000016B412CEE000100000000000000000000000000000000" +
	"010005000100010100010011001D00010010015E2C880002000B000000003544C87A" +
	"000E000000001DD7E06A00000100002994"

// mustHex decodes a hex string or fails the test.
func mustHex(t testing.TB, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("bad hex in test: %v", err)
	}
	return b
}

// corrupt returns a copy of b with byte i flipped, leaving b untouched.
func corrupt(b []byte, i int) []byte {
	c := bytes.Clone(b)
	c[i] ^= 0xFF
	return c
}

func TestReadFrame(t *testing.T) {
	frame := mustHex(t, specFrame)

	tests := []struct {
		name    string
		input   io.Reader
		wantErr error // nil: the frame must be read
	}{
		{"spec example", bytes.NewReader(frame), nil},
		{"arrives one byte at a time", iotest.OneByteReader(bytes.NewReader(frame)), nil},
		{"stream ends between frames", bytes.NewReader(nil), io.EOF},
		{"cut inside the header", bytes.NewReader(frame[:5]), io.ErrUnexpectedEOF},
		{"cut right after the header", bytes.NewReader(frame[:8]), io.ErrUnexpectedEOF},
		{"cut inside the CRC", bytes.NewReader(frame[:len(frame)-1]), io.ErrUnexpectedEOF},
		{"non-zero preamble", bytes.NewReader(corrupt(frame, 0)), ErrInvalidPreamble},
		{"length above the cap", bytes.NewReader(mustHex(t, "00000000FFFFFFFF")), ErrInvalidDataLength},
		{"length below the minimum", bytes.NewReader(mustHex(t, "0000000000000002")), ErrInvalidDataLength},
		{"corrupted data byte", bytes.NewReader(corrupt(frame, 20)), ErrInvalidCRC},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := ReadFrame(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if len(data) != 0x4A || data[0] != 0x8E {
				t.Errorf("data: len %d, first byte %#x; want len 74, first byte 0x8e", len(data), data[0])
			}
		})
	}
}

func TestReadFrameConsecutive(t *testing.T) {
	frame := mustHex(t, specFrame)
	stream := bytes.NewReader(append(bytes.Clone(frame), frame...))

	for i := range 2 {
		if _, err := ReadFrame(stream); err != nil {
			t.Fatalf("frame %d: %v", i+1, err)
		}
	}
	if _, err := ReadFrame(stream); !errors.Is(err, io.EOF) {
		t.Fatalf("after two frames: err = %v, want io.EOF", err)
	}
}

func TestReadFrameStaysAlignedAfterBadCRC(t *testing.T) {
	frame := mustHex(t, specFrame)
	stream := bytes.NewReader(append(corrupt(frame, 20), frame...))

	if _, err := ReadFrame(stream); !errors.Is(err, ErrInvalidCRC) {
		t.Fatalf("first frame: err = %v, want ErrInvalidCRC", err)
	}
	if _, err := ReadFrame(stream); err != nil {
		t.Fatalf("frame after a bad CRC must still be readable: %v", err)
	}
}

func TestCRC16(t *testing.T) {
	// "123456789" is the standard check input; 0xBB3D is the published
	// check value for CRC-16/ARC (https://reveng.sourceforge.io/crc-catalogue/16.htm).
	if got := crc16([]byte("123456789")); got != 0xBB3D {
		t.Fatalf("crc16 = %#04x, want 0xbb3d", got)
	}
}
