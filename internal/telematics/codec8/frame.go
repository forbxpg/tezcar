// SPDX-License-Identifier: BUSL-1.1

package codec8

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	maxDataLen = 1280
	minDataLen = 3
	nullByte   = 0x00
	crcLen     = 4
)

func validatePreamble(firstEightBytes [8]byte) error {
	for i := 0; i < 4; i++ {
		if firstEightBytes[i] != nullByte {
			return fmt.Errorf("%w: expected %#x, got %#x", ErrInvalidPreamble, nullByte, firstEightBytes[i])
		}
	}
	return nil
}

func validateDataLength(length uint32) error {
	if length < minDataLen {
		return fmt.Errorf("%w, got %d, min %d", ErrInvalidDataLength, length, minDataLen)
	}
	if length > maxDataLen {
		return fmt.Errorf("%w, got %d, max %d", ErrInvalidDataLength, length, maxDataLen)
	}
	return nil
}

func validateFirstEightBytes(firstEightBytes [8]byte) (uint32, error) {
	if err := validatePreamble(firstEightBytes); err != nil {
		return 0, err
	}
	dataLength := binary.BigEndian.Uint32(firstEightBytes[4:8])
	if err := validateDataLength(dataLength); err != nil {
		return 0, err
	}
	return dataLength, nil
}

func crc16(data []byte) uint16 {
	crc := uint16(0)
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if crc&1 != 0 {
				crc >>= 1
				crc ^= 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}

// ReadFrame reads incoming data from a Codec 8 Extended.
// It validates the frame's preamble, data length, and CRC.
// If the preamble is invalid, the function returns an ErrInvalidPreamble error.
// If the data length is invalid, the function returns an ErrInvalidDataLength error.
// If the CRC is invalid, the function returns an ErrInvalidCRC error.
// If the input reader returns an io.EOF error, the function returns an io.ErrUnexpectedEOF error.
// If the input reader returns any other error, the function returns that error.
// Finally, if the frame is valid, the function returns the frame's data field.
// P.S: You must close the connection if you got errors like:
// - ErrInvalidPreamble
// - ErrInvalidDataLength
// If you got ErrInvalidCRC, you can continue reading the next frame.
// Also wrap up the connection with `bufio.Reader` and make sure
// you have read deadlines set.
func ReadFrame(r io.Reader) ([]byte, error) {
	var firstEightBytes [8]byte
	_, err := io.ReadFull(r, firstEightBytes[:])
	if err != nil {
		return nil, err
	}
	dataLength, err := validateFirstEightBytes(firstEightBytes)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, dataLength+crcLen)
	if _, err := io.ReadFull(r, buf); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, io.ErrUnexpectedEOF
		}
		return nil, err
	}
	data := buf[:dataLength]
	want := binary.BigEndian.Uint32(buf[dataLength:])
	if got := uint32(crc16(data)); got != want {
		return nil, fmt.Errorf("%w: expected %#04x, got %#04x", ErrInvalidCRC, want, got)
	}
	return data, nil
}
