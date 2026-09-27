// SPDX-License-Identifier: BUSL-1.1

package codec8

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	maxDataLen = 1280
	minDataLen = 3
	nullByte   = 0x00
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

// ReadFrame validates the incoming data from codec8 extended frames
// and returns the frame's data field or an error if the data is invalid.
// If the io ends before the frame is complete, an io.ErrUnexpectedEOF error is returned.
// If the packet is complete then io.EOF is returned.
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
	data := make([]byte, dataLength)
	return data, nil // TODO: read the rest of the data from the reader
}
