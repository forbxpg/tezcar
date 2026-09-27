// SPDX-License-Identifier: BUSL-1.1

package codec8

import (
	"encoding/binary"
	"fmt"
)

// cursor is a helper type that allows reading bytes from a byte slice with bounds checking.
type cursor struct {
	b   []byte
	off int
	err error
}

func (c *cursor) take(n int) []byte {
	if c.err != nil {
		return nil
	}
	if n < 0 {
		c.err = fmt.Errorf("%w: negative length", ErrMalformedData)
		return nil
	}
	remaining := len(c.b) - c.off
	if remaining < n {
		c.err = fmt.Errorf("%w: not enough bytes. want %d, have %d", ErrMalformedData, n, remaining)
		return nil
	}
	v := c.b[c.off : c.off+n : c.off+n] // cap == len, so an append by the caller can't overwrite the bytes that follow
	c.off += n
	return v
}

func (c *cursor) u64() uint64 {
	if c.err != nil {
		return 0
	}
	arr := c.take(8)
	if arr == nil {
		return 0
	}
	return binary.BigEndian.Uint64(arr)
}

func (c *cursor) u32() uint32 {
	if c.err != nil {
		return 0
	}
	arr := c.take(4)
	if arr == nil {
		return 0
	}
	return binary.BigEndian.Uint32(arr)
}

func (c *cursor) u16() uint16 {
	if c.err != nil {
		return 0
	}
	arr := c.take(2)
	if arr == nil {
		return 0
	}
	return binary.BigEndian.Uint16(arr)
}

func (c *cursor) u8() uint8 {
	if c.err != nil {
		return 0
	}
	arr := c.take(1)
	if arr == nil {
		return 0
	}
	return arr[0]
}
