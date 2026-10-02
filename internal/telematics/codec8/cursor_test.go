// SPDX-License-Identifier: BUSL-1.1

package codec8

import (
	"errors"
	"testing"
)

func TestCursorReadsInOrder(t *testing.T) {
	c := cursor{b: []byte{
		0x01,       // u8
		0x02, 0x03, // u16
		0x04, 0x05, 0x06, 0x07, // u32
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, // u64
	}}
	if got := c.u8(); got != 0x01 {
		t.Errorf("u8 = %d; want 0x01", got)
	}
	if got := c.u16(); got != 0x0203 {
		t.Errorf("u16 = %#x, want 0x0203", got)
	}
	if got := c.u32(); got != 0x04050607 {
		t.Errorf("u32 = %#x, want 0x04050607", got)
	}
	if got := c.u64(); got != 0x08090A0B0C0D0E0F {
		t.Errorf("u64 = %#x, want 0x08090a0b0c0d0e0f", got)
	}
	if c.err != nil {
		t.Fatalf("err = %v, want nil", c.err)
	}
	if c.off != len(c.b) {
		t.Errorf("off = %d, want %d: every byte must be consumed", c.off, len(c.b))
	}
}

func TestCursorShortData(t *testing.T) {
	c := cursor{b: []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00}}

	if got := c.u32(); got != 1 {
		t.Fatalf("first u32 = %d, want 1", got)
	}
	if got := c.u32(); got != 0 {
		t.Errorf("second u32 = %d, want 0 when data runs out", got)
	}
	if !errors.Is(c.err, ErrMalformedData) {
		t.Fatalf("err = %v, want ErrMalformedData", c.err)
	}
}

func TestCursorKeepsFirstError(t *testing.T) {
	c := cursor{b: []byte{0x01}}

	c.u16() // fails: 1 byte left, 2 needed
	first := c.err
	if first == nil {
		t.Fatal("u16 on 1 byte: want an error")
	}

	if got := c.u8(); got != 0 {
		t.Errorf("u8 after an error = %d, want 0", got)
	}
	if c.err != first { //nolint:errorlint // identity on purpose: the first error must not be replaced
		t.Errorf("err = %v, want the first error %v", c.err, first)
	}
	if c.off != 0 {
		t.Errorf("off = %d, want 0: nothing may be read after an error", c.off)
	}
}
