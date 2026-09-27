// SPDX-License-Identifier: BUSL-1.1

package codec8

import "errors"

var (
	// ErrInvalidPreamble is returned when the preamble is not four zero bytes.
	ErrInvalidPreamble = errors.New("codec8: invalid preamble. expected four zero bytes")
	// ErrInvalidDataLength is returned when the data length is invalid.
	ErrInvalidDataLength = errors.New("codec8: invalid data length. expected 0x00000000")
	// ErrInvalidCodecID is returned when the codec ID is not 0x8E.
	ErrInvalidCodecID = errors.New("codec8: invalid codec ID. expected 0x8E")
	// ErrInvalidCRC is returned when the CRC is invalid.
	ErrInvalidCRC = errors.New("codec8: invalid CRC. expected 0x0000")
	// ErrMalformedData is returned when the data is malformed.
	ErrMalformedData = errors.New("codec8: malformed data")
)
