// SPDX-License-Identifier: BUSL-1.1
// Package codec8 reads and decodes Teltonika Codec 8 Extended frames.
//
// A device sends frames over one long-lived TCP connection. ReadFrame cuts one
// frame out of that byte stream and checks it; Decode turns the frame's data
// field into records. Only ReadFrame does I/O, so Decode is tested with plain
// byte slices and needs no network.
//
// Documentation of Codec8: https://wiki.teltonika-gps.com/view/Codec

package codec8
