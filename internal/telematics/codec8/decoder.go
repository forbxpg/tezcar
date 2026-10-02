// SPDX-License-Identifier: BUSL-1.1

package codec8

import (
	"bytes"
	"fmt"
	"time"
)

const codecID8E = 0x8E

func decodeGPS(c *cursor) GPS {
	return GPS{
		Longitude:  int32(c.u32()), //nolint:gosec // two's complement by spec: the top bit is the sign
		Latitude:   int32(c.u32()), //nolint:gosec // two's complement by spec: the top bit is the sign
		Altitude:   int16(c.u16()), //nolint:gosec // signed metres: below sea level is negative
		Angle:      c.u16(),
		Satellites: c.u8(),
		Speed:      c.u16(),
	}
}

func decodeIO(c *cursor) IO {
	eventID := c.u16()
	c.u16()
	values := make(map[uint16]uint64)

	for _, size := range []uint8{1, 2, 4, 8} {
		n := c.u16()
		for range n {
			key := c.u16()
			var value uint64
			switch size {
			case 1:
				value = uint64(c.u8())
			case 2:
				value = uint64(c.u16())
			case 4:
				value = uint64(c.u32())
			case 8:
				value = c.u64()
			}
			values[key] = value
		}
	}
	nx := c.u16()
	rawValues := make(map[uint16][]byte, nx)
	for range nx {
		id := c.u16()
		length := c.u16()
		rawValues[id] = bytes.Clone(c.take(int(length)))
	}
	return IO{
		EventID:   eventID,
		Values:    values,
		RawValues: rawValues,
	}
}

func decodeRecord(c *cursor) Record {
	return Record{
		Timestamp: time.UnixMilli(int64(c.u64())).UTC(), //nolint:gosec // milliseconds overflow int64 only after ~292 million years
		Priority:  Priority(c.u8()),
		GPS:       decodeGPS(c),
		IO:        decodeIO(c),
	}
}

// Decode turns the data field returned by ReadFrame into records.
//
// It returns ErrInvalidCodecID if the codec is not Codec 8 Extended, and
// ErrMalformedData if the data ends early, the two record counts differ, or
// bytes are left over. On error it returns no records.
func Decode(data []byte) ([]Record, error) {
	c := cursor{b: data}
	codecID := c.u8()
	dataLen := c.u8()

	// Number of Data 1
	if c.err == nil && codecID != codecID8E {
		return nil, fmt.Errorf("%w: expected codecID=%d, got %d", ErrInvalidCodecID, codecID8E, codecID)
	}
	records := make([]Record, 0, dataLen)
	for range dataLen {
		records = append(records, decodeRecord(&c))
	}

	// Number of Data 2
	dataLenTwo := c.u8()
	if c.err == nil && dataLenTwo != dataLen {
		return nil, fmt.Errorf("%w: number of data 1 is %d, number of data 2 is %d", ErrMalformedData, dataLen, dataLenTwo)
	}
	if c.err == nil && c.off != len(data) {
		return nil, fmt.Errorf("%w: %d bytes left over", ErrMalformedData, len(data)-c.off)
	}
	if c.err != nil {
		return nil, c.err
	}
	return records, nil
}
