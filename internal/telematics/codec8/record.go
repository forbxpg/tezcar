// SPDX-License-Identifier: BUSL-1.1

package codec8

import (
	"time"
)

// Priority is the priority of a record.
type Priority uint8

const (
	PriorityLow   Priority = iota // 0
	PriorityHigh                  // 1
	PriorityPanic                 // 2
)

// GPS is the GPS data of a record.
type GPS struct {
	Latitude   int32  // in degrees * 10 ^ 7
	Longitude  int32  // in degrees * 10 ^ 7
	Altitude   int16  // in meters
	Angle      uint16 // in degrees
	Speed      uint16 // in km/h
	Satellites uint8  // number of satellites used in the fix
}

// IO is the IO data of a record.
type IO struct {
	EventID   uint16            // event ID
	Values    map[uint16]uint64 // fixed-length values
	RawValues map[uint16][]byte // variable-length values
}

// Record is a single record from the codec8 packet stream.
type Record struct {
	Timestamp time.Time // in UTC
	Priority  Priority  // priority of the record
	GPS       GPS       // GPS data
	IO        IO        // IO data
}
