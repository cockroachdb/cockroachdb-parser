// Copyright 2014 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

// Package encoding exposes some utilities for encoding data as bytes.
package encoding

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/cockroachdb/apd/v3"
	"github.com/cockroachdb/cockroachdb-parser/pkg/geo/geopb"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/bitarray"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/duration"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/encoding/encodingtype"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/ipaddr"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/protoutil"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/timeofday"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/timetz"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/timeutil"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/uuid"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/redact"
)

const (
	encodedNull = 0x00
	// A marker greater than NULL but lower than any other value.
	// This value is not actually ever present in a stored key, but
	// it's used in keys used as span boundaries for index scans.
	encodedNotNull = 0x01

	floatNaN     = encodedNotNull + 1
	floatNeg     = floatNaN + 1
	floatZero    = floatNeg + 1
	floatPos     = floatZero + 1
	floatNaNDesc = floatPos + 1 // NaN encoded descendingly

	// The gap between floatNaNDesc and bytesMarker was left for
	// compatibility reasons.
	bytesMarker          byte = 0x12
	bytesDescMarker           = bytesMarker + 1
	timeMarker                = bytesDescMarker + 1
	durationBigNegMarker      = timeMarker + 1 // Only used for durations < MinInt64 nanos.
	durationMarker            = durationBigNegMarker + 1
	durationBigPosMarker      = durationMarker + 1 // Only used for durations > MaxInt64 nanos.

	decimalNaN              = durationBigPosMarker + 1 // 24
	decimalNegativeInfinity = decimalNaN + 1
	decimalNegLarge         = decimalNegativeInfinity + 1
	decimalNegMedium        = decimalNegLarge + 11
	decimalNegSmall         = decimalNegMedium + 1
	decimalZero             = decimalNegSmall + 1
	decimalPosSmall         = decimalZero + 1
	decimalPosMedium        = decimalPosSmall + 1
	decimalPosLarge         = decimalPosMedium + 11
	decimalInfinity         = decimalPosLarge + 1
	decimalNaNDesc          = decimalInfinity + 1 // NaN encoded descendingly
	decimalTerminator       = 0x00

	jsonInvertedIndex = decimalNaNDesc + 1
	jsonEmptyArray    = jsonInvertedIndex + 1
	jsonEmptyObject   = jsonEmptyArray + 1

	bitArrayMarker             = jsonEmptyObject + 1
	bitArrayDescMarker         = bitArrayMarker + 1
	bitArrayDataTerminator     = 0x00
	bitArrayDataDescTerminator = 0xff

	timeTZMarker  = bitArrayDescMarker + 1
	geoMarker     = timeTZMarker + 1
	geoDescMarker = geoMarker + 1

	// Markers and terminators for key encoding Datum arrays in sorted order.
	// For the arrayKeyMarker and other types like bytes and bit arrays, it
	// might be unclear why we have a separate marker for the ascending and
	// descending cases. This is necessary because the terminators for these
	// encodings are different depending on the direction the data is encoded
	// in. In order to safely decode a set of bytes without knowing the direction
	// of the encoding, we must store this information in the marker. Otherwise,
	// we would not know what terminator to look for when decoding this format.
	arrayKeyMarker           = geoDescMarker + 1
	arrayKeyDescendingMarker = arrayKeyMarker + 1

	box2DMarker            = arrayKeyDescendingMarker + 1
	geoInvertedIndexMarker = box2DMarker + 1

	emptyArray = geoInvertedIndexMarker + 1
	voidMarker = emptyArray + 1

	// Defining different key markers, for the ascending designation,
	// for handling different JSON values.

	// Postgres currently has a special case (maybe a bug) where the empty JSON
	// Array sorts before all other JSON values. See the bug report:
	// https://www.postgresql.org/message-id/17873-826fdc8bbcace4f1%40postgresql.org
	jsonEmptyArrayKeyMarker = voidMarker + 1
	jsonNullKeyMarker       = jsonEmptyArrayKeyMarker + 1
	jsonStringKeyMarker     = jsonNullKeyMarker + 1
	jsonNumberKeyMarker     = jsonStringKeyMarker + 1
	jsonFalseKeyMarker      = jsonNumberKeyMarker + 1
	jsonTrueKeyMarker       = jsonFalseKeyMarker + 1
	jsonArrayKeyMarker      = jsonTrueKeyMarker + 1
	jsonObjectKeyMarker     = jsonArrayKeyMarker + 1

	arrayKeyTerminator           byte = 0x00
	arrayKeyDescendingTerminator byte = 0xFF
	// We use different null encodings for nulls within key arrays. Doing this
	// allows for the terminator to be less/greater than the null value within
	// arrays. These byte values overlap with encodedNotNull and
	// encodedNotNullDesc, but they can only exist within an encoded array key.
	// Because of the context, they cannot be ambiguous with these other bytes.
	ascendingNullWithinArrayKey  byte = 0x01
	descendingNullWithinArrayKey byte = 0xFE

	// Defining different key markers, for the descending designation,
	// for handling different JSON values.
	jsonEmptyArrayKeyDescendingMarker = jsonObjectKeyMarker + 8
	jsonNullKeyDescendingMarker       = jsonEmptyArrayKeyDescendingMarker - 1
	jsonStringKeyDescendingMarker     = jsonNullKeyDescendingMarker - 1
	jsonNumberKeyDescendingMarker     = jsonStringKeyDescendingMarker - 1
	jsonFalseKeyDescendingMarker      = jsonNumberKeyDescendingMarker - 1
	jsonTrueKeyDescendingMarker       = jsonFalseKeyDescendingMarker - 1
	jsonArrayKeyDescendingMarker      = jsonTrueKeyDescendingMarker - 1
	jsonObjectKeyDescendingMarker     = jsonArrayKeyDescendingMarker - 1

	// Terminators for JSON Key encoding.
	jsonKeyTerminator           byte = 0x00
	jsonKeyDescendingTerminator byte = 0xFF

	// IntMin is chosen such that the range of int tags does not overlap the
	// ascii character set that is frequently used in testing.
	IntMin      = 0x80 // 128
	intMaxWidth = 8
	intZero     = IntMin + intMaxWidth           // 136
	intSmall    = IntMax - intZero - intMaxWidth // 109
	// IntMax is the maximum int tag value.
	IntMax = 0xfd // 253

	// Nulls come last when encoded descendingly.
	// This value is not actually ever present in a stored key, but
	// it's used in keys used as span boundaries for index scans.
	encodedNotNullDesc = 0xfe
	encodedNullDesc    = 0xff

	// offsetSecsToMicros is a constant that allows conversion from seconds
	// to microseconds for offsetSecs type calculations (e.g. for TimeTZ).
	offsetSecsToMicros = 1000000
)

const (
	// EncodedDurationMaxLen is the largest number of bytes used when encoding a
	// Duration.
	EncodedDurationMaxLen = 1 + 3*binary.MaxVarintLen64 // 3 varints are encoded.
	// EncodedTimeTZMaxLen is the largest number of bytes used when encoding a
	// TimeTZ.
	EncodedTimeTZMaxLen = 1 + binary.MaxVarintLen64 + binary.MaxVarintLen32
)

// Direction for ordering results.
type Direction int

// Direction values.
const (
	_ Direction = iota
	Ascending
	Descending
)

const escapeLength = 2

// EncodeUint16Ascending encodes the uint16 value using a big-endian 2 byte
// representation. The bytes are appended to the supplied buffer and
// the final buffer is returned.
func EncodeUint16Ascending(b []byte, v uint16) []byte {
	return append(b, byte(v>>8), byte(v))
}

// PutUint16Ascending encodes the uint16 value using a big-endian 2 byte
// representation at the specified index, lengthening the input slice if
// necessary.
func PutUint16Ascending(b []byte, v uint16, idx int) []byte {
	for len(b) < idx+2 {
		b = append(b, 0)
	}
	b[idx] = byte(v >> 8)
	b[idx+1] = byte(v)
	return b
}

// EncodeUint32Ascending encodes the uint32 value using a big-endian 4 byte
// representation. The bytes are appended to the supplied buffer and
// the final buffer is returned.
func EncodeUint32Ascending(b []byte, v uint32) []byte {
	return append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

// PutUint32Ascending encodes the uint32 value using a big-endian 4 byte
// representation at the specified index, lengthening the input slice if
// necessary.
func PutUint32Ascending(b []byte, v uint32, idx int) []byte {
	for len(b) < idx+4 {
		b = append(b, 0)
	}
	b[idx] = byte(v >> 24)
	b[idx+1] = byte(v >> 16)
	b[idx+2] = byte(v >> 8)
	b[idx+3] = byte(v)
	return b
}

// EncodeUint32Descending encodes the uint32 value so that it sorts in
// reverse order, from largest to smallest.
func EncodeUint32Descending(b []byte, v uint32) []byte {
	return EncodeUint32Ascending(b, ^v)
}

// DecodeUint16Ascending decodes a uint16 from the input buffer, treating
// the input as a big-endian 2 byte uint16 representation. The remainder
// of the input buffer and the decoded uint16 are returned.
func DecodeUint16Ascending(b []byte) ([]byte, uint16, error) {
	if len(b) < 2 {
		return nil, 0, errors.Errorf("insufficient bytes to decode uint16 int value")
	}
	v := binary.BigEndian.Uint16(b)
	return b[2:], v, nil
}

// DecodeUint32Ascending decodes a uint32 from the input buffer, treating
// the input as a big-endian 4 byte uint32 representation. The remainder
// of the input buffer and the decoded uint32 are returned.
func DecodeUint32Ascending(b []byte) ([]byte, uint32, error) {
	if len(b) < 4 {
		return nil, 0, errors.Errorf("insufficient bytes to decode uint32 int value")
	}
	v := binary.BigEndian.Uint32(b)
	return b[4:], v, nil
}

// DecodeUint32Descending decodes a uint32 value which was encoded
// using EncodeUint32Descending.
func DecodeUint32Descending(b []byte) ([]byte, uint32, error) {
	leftover, v, err := DecodeUint32Ascending(b)
	return leftover, ^v, err
}

const uint64AscendingEncodedLength = 8

// EncodeUint64Ascending encodes the uint64 value using a big-endian 8 byte
// representation. The bytes are appended to the supplied buffer and
// the final buffer is returned.
func EncodeUint64Ascending(b []byte, v uint64) []byte {
	return append(b,
		byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
		byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

// EncodeUint64Descending encodes the uint64 value so that it sorts in
// reverse order, from largest to smallest.
func EncodeUint64Descending(b []byte, v uint64) []byte {
	return EncodeUint64Ascending(b, ^v)
}

// DecodeUint64Ascending decodes a uint64 from the input buffer, treating
// the input as a big-endian 8 byte uint64 representation. The remainder
// of the input buffer and the decoded uint64 are returned.
func DecodeUint64Ascending(b []byte) ([]byte, uint64, error) {
	if len(b) < 8 {
		return nil, 0, errors.Errorf("insufficient bytes to decode uint64 int value")
	}
	v := binary.BigEndian.Uint64(b)
	return b[8:], v, nil
}

// DecodeUint64Descending decodes a uint64 value which was encoded
// using EncodeUint64Descending.
func DecodeUint64Descending(b []byte) ([]byte, uint64, error) {
	leftover, v, err := DecodeUint64Ascending(b)
	return leftover, ^v, err
}

// MaxVarintLen is the maximum length of a value encoded using any of:
// - EncodeVarintAscending
// - EncodeVarintDescending
// - EncodeUvarintAscending
// - EncodeUvarintDescending
const MaxVarintLen = 9

// EncodeVarintAscending encodes the int64 value using a variable length
// (length-prefixed) representation. The length is encoded as a single
// byte. If the value to be encoded is negative the length is encoded
// as 8-numBytes. If the value is positive it is encoded as
// 8+numBytes. The encoded bytes are appended to the supplied buffer
// and the final buffer is returned.
func EncodeVarintAscending(b []byte, v int64) []byte {
	if v < 0 {
		switch {
		case v >= -0xff:
			return append(b, IntMin+7, byte(v))
		case v >= -0xffff:
			return append(b, IntMin+6, byte(v>>8), byte(v))
		case v >= -0xffffff:
			return append(b, IntMin+5, byte(v>>16), byte(v>>8), byte(v))
		case v >= -0xffffffff:
			return append(b, IntMin+4, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		case v >= -0xffffffffff:
			return append(b, IntMin+3, byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8),
				byte(v))
		case v >= -0xffffffffffff:
			return append(b, IntMin+2, byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16),
				byte(v>>8), byte(v))
		case v >= -0xffffffffffffff:
			return append(b, IntMin+1, byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24),
				byte(v>>16), byte(v>>8), byte(v))
		default:
			return append(b, IntMin, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
				byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		}
	}
	return EncodeUvarintAscending(b, uint64(v))
}

// EncodeVarintDescending encodes the int64 value so that it sorts in reverse
// order, from largest to smallest.
func EncodeVarintDescending(b []byte, v int64) []byte {
	return EncodeVarintAscending(b, ^v)
}

// getVarintLen returns the encoded length of an encoded varint. Assumes the
// slice has at least one byte.
func getVarintLen(b []byte) (int, error) {
	length := int(b[0]) - intZero
	if length >= 0 {
		if length <= intSmall {
			// just the tag
			return 1, nil
		}
		// tag and length-intSmall bytes
		length = 1 + length - intSmall
	} else {
		// tag and -length bytes
		length = 1 - length
	}

	if length > len(b) {
		return 0, errors.Errorf("varint length %d exceeds slice length %d", length, len(b))
	}
	return length, nil
}

// DecodeVarintAscending decodes a value encoded by EncodeVarintAscending.
func DecodeVarintAscending(b []byte) ([]byte, int64, error) {
	if len(b) == 0 {
		return nil, 0, errors.Errorf("insufficient bytes to decode varint value")
	}
	length := int(b[0]) - intZero
	if length < 0 {
		length = -length
		remB := b[1:]
		if len(remB) < length {
			return nil, 0, errors.Errorf("insufficient bytes to decode varint value: %q", remB)
		}
		var v int64
		// Use the ones-complement of each encoded byte in order to build
		// up a positive number, then take the ones-complement again to
		// arrive at our negative value.
		for _, t := range remB[:length] {
			v = (v << 8) | int64(^t)
		}
		return remB[length:], ^v, nil
	}

	remB, v, err := DecodeUvarintAscending(b)
	if err != nil {
		return remB, 0, err
	}
	if v > math.MaxInt64 {
		return nil, 0, errors.Errorf("varint %d overflows int64", v)
	}
	return remB, int64(v), nil
}

// DecodeVarintDescending decodes a int64 value which was encoded
// using EncodeVarintDescending.
func DecodeVarintDescending(b []byte) ([]byte, int64, error) {
	leftover, v, err := DecodeVarintAscending(b)
	return leftover, ^v, err
}

// EncodeUvarintAscending encodes the uint64 value using a variable length
// (length-prefixed) representation. The length is encoded as a single
// byte indicating the number of encoded bytes (-8) to follow. See
// EncodeVarintAscending for rationale. The encoded bytes are appended to the
// supplied buffer and the final buffer is returned.
func EncodeUvarintAscending(b []byte, v uint64) []byte {
	switch {
	case v <= intSmall:
		return append(b, intZero+byte(v))
	case v <= 0xff:
		return append(b, IntMax-7, byte(v))
	case v <= 0xffff:
		return append(b, IntMax-6, byte(v>>8), byte(v))
	case v <= 0xffffff:
		return append(b, IntMax-5, byte(v>>16), byte(v>>8), byte(v))
	case v <= 0xffffffff:
		return append(b, IntMax-4, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	case v <= 0xffffffffff:
		return append(b, IntMax-3, byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8),
			byte(v))
	case v <= 0xffffffffffff:
		return append(b, IntMax-2, byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16),
			byte(v>>8), byte(v))
	case v <= 0xffffffffffffff:
		return append(b, IntMax-1, byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24),
			byte(v>>16), byte(v>>8), byte(v))
	default:
		return append(b, IntMax, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
			byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
}

// EncodedLengthUvarintAscending returns the length of the variable length
// representation, i.e. the number of bytes appended by EncodeUvarintAscending.
func EncodedLengthUvarintAscending(v uint64) int {
	switch {
	case v <= intSmall:
		return 1
	case v <= 0xff:
		return 2
	case v <= 0xffff:
		return 3
	case v <= 0xffffff:
		return 4
	case v <= 0xffffffff:
		return 5
	case v <= 0xffffffffff:
		return 6
	case v <= 0xffffffffffff:
		return 7
	case v <= 0xffffffffffffff:
		return 8
	default:
		return 9
	}
}

// EncodeUvarintDescending encodes the uint64 value so that it sorts in
// reverse order, from largest to smallest.
func EncodeUvarintDescending(b []byte, v uint64) []byte {
	switch {
	case v == 0:
		return append(b, IntMin+8)
	case v <= 0xff:
		v = ^v
		return append(b, IntMin+7, byte(v))
	case v <= 0xffff:
		v = ^v
		return append(b, IntMin+6, byte(v>>8), byte(v))
	case v <= 0xffffff:
		v = ^v
		return append(b, IntMin+5, byte(v>>16), byte(v>>8), byte(v))
	case v <= 0xffffffff:
		v = ^v
		return append(b, IntMin+4, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	case v <= 0xffffffffff:
		v = ^v
		return append(b, IntMin+3, byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8),
			byte(v))
	case v <= 0xffffffffffff:
		v = ^v
		return append(b, IntMin+2, byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16),
			byte(v>>8), byte(v))
	case v <= 0xffffffffffffff:
		v = ^v
		return append(b, IntMin+1, byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24),
			byte(v>>16), byte(v>>8), byte(v))
	default:
		v = ^v
		return append(b, IntMin, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
			byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
}

// highestByteIndex returns the index (0 to 7) of the highest nonzero byte in v.
func highestByteIndex(v uint64) int {
	l := 0
	if v > 0xffffffff {
		v >>= 32
		l += 4
	}
	if v > 0xffff {
		v >>= 16
		l += 2
	}
	if v > 0xff {
		l++
	}
	return l
}

// EncLenUvarintAscending returns the encoding length for EncodeUvarintAscending
// without actually encoding.
func EncLenUvarintAscending(v uint64) int {
	if v <= intSmall {
		return 1
	}
	return 2 + highestByteIndex(v)
}

// EncLenUvarintDescending returns the encoding length for
// EncodeUvarintDescending without actually encoding.
func EncLenUvarintDescending(v uint64) int {
	if v == 0 {
		return 1
	}
	return 2 + highestByteIndex(v)
}

// GetUvarintLen is similar to DecodeUvarintAscending except that it returns the
// length of the prefix that encodes a uint64 value in bytes without actually
// decoding the value. An error is returned if b does not contain a valid
// encoding of an unsigned int datum.
func GetUvarintLen(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, errors.Errorf("insufficient bytes to decode uvarint value")
	}
	length := int(b[0]) - intZero
	if length <= intSmall {
		return 1, nil
	}
	length -= intSmall
	if length < 0 || length > 8 {
		return 0, errors.Errorf("invalid uvarint length of %d", length)
	} else if len(b) <= length {
		// Note: we use <= for comparison here as opposed to the < in
		// DecodeUvarintAscending because in the latter the first byte for the
		// uvarint is removed as part of decoding. We need to account for the first
		// byte when assessing the size.
		return 0, errors.Errorf("insufficient bytes to decode uvarint value: %q", b)
	}
	return 1 + length, nil
}

// DecodeUvarintAscending decodes a uint64 encoded uint64 from the input
// buffer. The remainder of the input buffer and the decoded uint64
// are returned.
func DecodeUvarintAscending(b []byte) ([]byte, uint64, error) {
	if len(b) == 0 {
		return nil, 0, errors.Errorf("insufficient bytes to decode uvarint value")
	}
	length := int(b[0]) - intZero
	b = b[1:] // skip length byte
	if length <= intSmall {
		return b, uint64(length), nil
	}
	length -= intSmall
	if length < 0 || length > 8 {
		return nil, 0, errors.Errorf("invalid uvarint length of %d", length)
	} else if len(b) < length {
		return nil, 0, errors.Errorf("insufficient bytes to decode uvarint value: %q", b)
	}
	var v uint64
	// It is faster to range over the elements in a slice than to index
	// into the slice on each loop iteration.
	for _, t := range b[:length] {
		v = (v << 8) | uint64(t)
	}
	return b[length:], v, nil
}

// DecodeUvarintDescending decodes a uint64 value which was encoded
// using EncodeUvarintDescending.
func DecodeUvarintDescending(b []byte) ([]byte, uint64, error) {
	if len(b) == 0 {
		return nil, 0, errors.Errorf("insufficient bytes to decode uvarint value")
	}
	length := intZero - int(b[0])
	b = b[1:] // skip length byte
	if length < 0 || length > 8 {
		return nil, 0, errors.Errorf("invalid uvarint length of %d", length)
	} else if len(b) < length {
		return nil, 0, errors.Errorf("insufficient bytes to decode uvarint value: %q", b)
	}
	var x uint64
	for _, t := range b[:length] {
		x = (x << 8) | uint64(^t)
	}
	return b[length:], x, nil
}

const (
	// <term>     -> \x00\x01
	// \x00       -> \x00\xff
	escape                   byte = 0x00
	escapedTerm              byte = 0x01
	escapedJSONObjectKeyTerm byte = 0x02
	escapedJSONArray         byte = 0x03
	escaped00                byte = 0xff
	escapedFF                byte = 0x00
)

type escapes struct {
	escape      byte
	escapedTerm byte
	escaped00   byte
	escapedFF   byte
	marker      byte
}

var (
	ascendingBytesEscapes  = escapes{escape, escapedTerm, escaped00, escapedFF, bytesMarker}
	descendingBytesEscapes = escapes{^escape, ^escapedTerm, ^escaped00, ^escapedFF, bytesDescMarker}

	ascendingGeoEscapes  = escapes{escape, escapedTerm, escaped00, escapedFF, geoMarker}
	descendingGeoEscapes = escapes{^escape, ^escapedTerm, ^escaped00, ^escapedFF, geoDescMarker}
)

// EncodeBytesAscending encodes the []byte value using an escape-based
// encoding. The encoded value is terminated with the sequence
// "\x00\x01" which is guaranteed to not occur elsewhere in the
// encoded value. The encoded bytes are append to the supplied buffer
// and the resulting buffer is returned.
func EncodeBytesAscending(b []byte, data []byte) []byte {
	return encodeBytesAscendingWithTerminatorAndPrefix(b, data, ascendingBytesEscapes.escapedTerm, bytesMarker)
}

// EncodeNextBytesAscending encodes the []byte value with an extra 0x00 byte
// appended before encoding. It's equivalent to
//
//	EncodeBytesAscending(b, append(data, 0x00))
//
// but may avoid an allocation when the data slice does not have additional
// capacity.
func EncodeNextBytesAscending(b []byte, data []byte) []byte {
	b = append(b, bytesMarker)
	return encodeNextBytesAscendingWithTerminator(b, data, ascendingBytesEscapes.escapedTerm)
}

func encodeNextBytesAscendingWithTerminator(b []byte, data []byte, terminator byte) []byte {
	bs := encodeBytesAscendingWithoutTerminatorOrPrefix(b, data)
	return append(bs, escape, escaped00, escape, terminator)
}

// encodeBytesAscendingWithTerminatorAndPrefix encodes the []byte value using an escape-based
// encoding. The encoded value is terminated with the sequence
// "\x00\terminator". The encoded bytes are append to the supplied buffer
// and the resulting buffer is returned. The terminator allows us to pass
// different terminators for things such as JSON key encoding.
func encodeBytesAscendingWithTerminatorAndPrefix(
	b []byte, data []byte, terminator byte, prefix byte,
) []byte {
	b = append(b, prefix)
	return encodeBytesAscendingWithTerminator(b, data, terminator)
}

// encodeBytesAscendingWithTerminator encodes the []byte value using an escape-based
// encoding. The encoded value is terminated with the sequence
// "\x00\terminator". The encoded bytes are append to the supplied buffer
// and the resulting buffer is returned. The terminator allows us to pass
// different terminators for things such as JSON key encoding.
func encodeBytesAscendingWithTerminator(b []byte, data []byte, terminator byte) []byte {
	bs := encodeBytesAscendingWithoutTerminatorOrPrefix(b, data)
	return append(bs, escape, terminator)
}

// encodeBytesAscendingWithoutTerminatorOrPrefix encodes the []byte value using an escape-based
// encoding.
func encodeBytesAscendingWithoutTerminatorOrPrefix(b []byte, data []byte) []byte {
	for {
		// IndexByte is implemented by the go runtime in assembly and is
		// much faster than looping over the bytes in the slice.
		i := bytes.IndexByte(data, escape)
		if i == -1 {
			break
		}
		b = append(b, data[:i]...)
		b = append(b, escape, escaped00)
		data = data[i+1:]
	}
	return append(b, data...)
}

// EncodeBytesDescending encodes the []byte value using an
// escape-based encoding and then inverts (ones complement) the result
// so that it sorts in reverse order, from larger to smaller
// lexicographically.
func EncodeBytesDescending(b []byte, data []byte) []byte {
	n := len(b)
	b = EncodeBytesAscending(b, data)
	b[n] = bytesDescMarker
	onesComplement(b[n+1:])
	return b
}

// EncodeBytesSize returns the size of the []byte value when encoded using
// EncodeBytes{Ascending,Descending}. The function accounts for the encoding
// marker, escaping, and the terminator.
func EncodeBytesSize(data []byte) int {
	// Encoding overhead:
	// +1 for [bytesMarker] prefix
	// +2 for [escape, escapedTerm] suffix
	// +1 for each byte that needs to be escaped
	//
	// NOTE: bytes.Count is implemented by the go runtime in assembly and is
	// much faster than looping over the bytes in the slice, especially when
	// given a single-byte separator.
	return len(data) + 3 + bytes.Count(data, []byte{escape})
}

// EncodeNextBytesSize returns the size of the []byte value when suffixed with a
// zero byte and then encoded using EncodeNextBytes{Ascending,Descending}. The
// function accounts for the encoding marker, escaping, and the terminator.
func EncodeNextBytesSize(data []byte) int {
	// Encoding overhead:
	// +1 for [bytesMarker] prefix
	// +2 for [escape, escapedTerm] suffix
	// +1 for each byte that needs to be escaped
	// +2 for the appended 0x00 byte, plus its escaping byte
	//
	// NOTE: bytes.Count is implemented by the go runtime in assembly and is
	// much faster than looping over the bytes in the slice, especially when
	// given a single-byte separator.
	return len(data) + 5 + bytes.Count(data, []byte{escape})
}

// DecodeBytesAscending decodes a []byte value from the input buffer
// which was encoded using EncodeBytesAscending. The decoded bytes
// are appended to r. The remainder of the input buffer and the
// decoded []byte are returned.
func DecodeBytesAscending(b []byte, r []byte) ([]byte, []byte, error) {
	return decodeBytesInternal(b, r, ascendingBytesEscapes, true /* expectMarker */, false /* deepCopy */)
}

// ValidateDecodeBytesAscending is like DecodeBytesAscending, but discards the
// decoded bytes. The remainder of the input buffer is returned on success.
func ValidateDecodeBytesAscending(b []byte) ([]byte, error) {
	return validateDecodeBytesInternal(b, ascendingBytesEscapes, true /* expectMarker */)
}

// DecodeBytesAscendingDeepCopy is the same as DecodeBytesAscending, but the
// decoded []byte will never alias memory of b.
func DecodeBytesAscendingDeepCopy(b []byte, r []byte) ([]byte, []byte, error) {
	return decodeBytesInternal(b, r, ascendingBytesEscapes, true /* expectMarker */, true /* deepCopy */)
}

// DecodeBytesDescending decodes a []byte value from the input buffer
// which was encoded using EncodeBytesDescending. The decoded bytes
// are appended to r. The remainder of the input buffer and the
// decoded []byte are returned.
//
// Note that this method internally will always perform a deep copy, so there is
// no need to introduce DecodeBytesDescendingDeepCopy to mirror
// DecodeBytesAscendingDeepCopy.
func DecodeBytesDescending(b []byte, r []byte) ([]byte, []byte, error) {
	// Ask for the deep copy to make sure we never get back a sub-slice of `b`,
	// since we're going to modify the contents of the slice.
	b, r, err := decodeBytesInternal(b, r, descendingBytesEscapes, true /* expectMarker */, true /* deepCopy */)
	onesComplement(r)
	return b, r, err
}

// decodeBytesInternal decodes an encoded []byte value from b and appends it to
// r. The remainder of b and the decoded []byte are returned. If deepCopy is
// true, then the decoded []byte will be deep copied from b and there will no
// aliasing of the same memory.
func decodeBytesInternal(
	b []byte, r []byte, e escapes, expectMarker bool, deepCopy bool,
) ([]byte, []byte, error) {
	if expectMarker {
		if len(b) == 0 || b[0] != e.marker {
			return nil, nil, errors.Errorf("did not find marker %#x in buffer %#x", e.marker, b)
		}
		b = b[1:]
	}

	for {
		i := bytes.IndexByte(b, e.escape)
		if i == -1 {
			return nil, nil, errors.Errorf("did not find terminator %#x in buffer %#x", e.escape, b)
		}
		if i+1 >= len(b) {
			return nil, nil, errors.Errorf("malformed escape in buffer %#x", b)
		}
		v := b[i+1]
		if v == e.escapedTerm {
			if r == nil && !deepCopy {
				r = b[:i]
			} else {
				r = append(r, b[:i]...)
			}
			return b[i+2:], r, nil
		}

		if v != e.escaped00 {
			return nil, nil, errors.Errorf("unknown escape sequence: %#x %#x", e.escape, v)
		}

		r = append(r, b[:i]...)
		r = append(r, e.escapedFF)
		b = b[i+2:]
	}
}

// validateDecodeBytesInternal decodes an encoded []byte value from b,
// discarding the decoded value. The remainder of b is returned on success, or a
// non-nil error otherwise.
func validateDecodeBytesInternal(b []byte, e escapes, expectMarker bool) ([]byte, error) {
	if expectMarker {
		if len(b) == 0 || b[0] != e.marker {
			return nil, errors.Errorf("did not find marker %#x in buffer %#x", e.marker, b)
		}
		b = b[1:]
	}

	for {
		i := bytes.IndexByte(b, e.escape)
		if i == -1 {
			return nil, errors.Errorf("did not find terminator %#x in buffer %#x", e.escape, b)
		}
		if i+1 >= len(b) {
			return nil, errors.Errorf("malformed escape in buffer %#x", b)
		}
		v := b[i+1]
		if v == e.escapedTerm {
			return b[i+2:], nil
		}

		if v != e.escaped00 {
			return nil, errors.Errorf("unknown escape sequence: %#x %#x", e.escape, v)
		}
		b = b[i+2:]
	}
}

// getBytesLength finds the length of a bytes encoding.
func getBytesLength(b []byte, e escapes) (int, error) {
	// Skip the tag.
	skipped := 1
	for {
		i := bytes.IndexByte(b[skipped:], e.escape)
		if i == -1 {
			return 0, errors.Errorf("did not find terminator %#x in buffer %#x", e.escape, b)
		}
		if i+1 >= len(b) {
			return 0, errors.Errorf("malformed escape in buffer %#x", b)
		}
		skipped += i + escapeLength
		if b[skipped-1] == e.escapedTerm {
			return skipped, nil
		}
	}
}

// prettyPrintInvertedIndexKey returns a string representation of the path part of a JSON inverted
// index.
func prettyPrintInvertedIndexKey(b []byte) (string, []byte, error) {
	outBytes := ""
	// We're skipping the first byte because it's the JSON tag.
	tempB := b[1:]
	for {
		i := bytes.IndexByte(tempB, escape)

		if i == -1 {
			return "", nil, errors.Errorf("did not find terminator %#x in buffer %#x", escape, b)
		}
		if i+1 >= len(tempB) {
			return "", nil, errors.Errorf("malformed escape in buffer %#x", b)
		}

		switch tempB[i+1] {
		case escapedTerm:
			if len(tempB[:i]) > 0 {
				outBytes = outBytes + strconv.Quote(UnsafeConvertBytesToString(tempB[:i]))
			} else {
				lenOut := len(outBytes)
				if lenOut > 1 && outBytes[lenOut-1] == '/' {
					outBytes = outBytes[:lenOut-1]
				}
			}
			return outBytes, tempB[i+escapeLength:], nil
		case escapedJSONObjectKeyTerm:
			outBytes = outBytes + strconv.Quote(UnsafeConvertBytesToString(tempB[:i])) + "/"
		case escapedJSONArray:
			outBytes = outBytes + "Arr/"
			if i+2 >= len(tempB) {
				// The key ends in an escaped JSON array byte, which is used in
				// spans to scan over non-empty arrays.
				return outBytes, nil, nil
			}
		case escaped00:
			if i+2 >= len(tempB) {
				// The key ends in an escaped NULL byte.
				return outBytes, nil, nil
			}
		default:
			return "", nil, errors.Errorf("malformed escape in buffer %#x", b)

		}

		tempB = tempB[i+escapeLength:]
	}
}

// UnsafeConvertStringToBytes converts a string to a byte array to be used with
// string encoding functions. Note that the output byte array should not be
// modified if the input string is expected to be used again - doing so could
// violate Go semantics.
func UnsafeConvertStringToBytes(s string) []byte {
	// unsafe.StringData output is unspecified for empty string input so always
	// return nil.
	if len(s) == 0 {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// EncodeStringAscending encodes the string value using an escape-based encoding. See
// EncodeBytes for details. The encoded bytes are append to the supplied buffer
// and the resulting buffer is returned.
func EncodeStringAscending(b []byte, s string) []byte {
	return encodeStringAscendingWithTerminatorAndPrefix(b, s, ascendingBytesEscapes.escapedTerm, bytesMarker)
}

// encodeStringAscendingWithTerminatorAndPrefix encodes the string value using an escape-based encoding. See
// EncodeBytes for details. The encoded bytes are append to the supplied buffer
// and the resulting buffer is returned. We can also pass a terminator byte to be used with
// JSON key encoding.
func encodeStringAscendingWithTerminatorAndPrefix(
	b []byte, s string, terminator byte, prefix byte,
) []byte {
	unsafeString := UnsafeConvertStringToBytes(s)
	return encodeBytesAscendingWithTerminatorAndPrefix(b, unsafeString, terminator, prefix)
}

// EncodeJSONKeyStringAscending encodes the JSON key string value with a JSON specific escaped
// terminator. This allows us to encode keys in the same number of bytes as a string,
// while at the same time giving us a sentinel to identify JSON keys. The end parameter is used
// to determine if this is the last key in a a JSON path. If it is we don't add a separator after it.
func EncodeJSONKeyStringAscending(b []byte, s string, end bool) []byte {
	str := UnsafeConvertStringToBytes(s)

	if end {
		return encodeBytesAscendingWithoutTerminatorOrPrefix(b, str)
	}
	return encodeBytesAscendingWithTerminator(b, str, escapedJSONObjectKeyTerm)
}

// EncodeJSONEmptyArray returns a byte array b with a byte to signify an empty JSON array.
func EncodeJSONEmptyArray(b []byte) []byte {
	return append(b, escape, escapedTerm, jsonEmptyArray)
}

// AddJSONPathTerminator adds a json path terminator to a byte array.
func AddJSONPathTerminator(b []byte) []byte {
	return append(b, escape, escapedTerm)
}

// AddJSONPathSeparator adds a json path separator to a byte array.
func AddJSONPathSeparator(b []byte) []byte {
	return append(b, escape, escapedJSONObjectKeyTerm)
}

// EncodeJSONEmptyObject returns a byte array b with a byte to signify an empty JSON object.
func EncodeJSONEmptyObject(b []byte) []byte {
	return append(b, escape, escapedTerm, jsonEmptyObject)
}

// EncodeEmptyArray returns a byte array b with a byte to signify an empty array.
func EncodeEmptyArray(b []byte) []byte {
	return append(b, emptyArray)
}

// EncodeStringDescending is the descending version of EncodeStringAscending.
func EncodeStringDescending(b []byte, s string) []byte {
	arg := UnsafeConvertStringToBytes(s)
	return EncodeBytesDescending(b, arg)
}

// UnsafeConvertBytesToString performs an unsafe conversion from a []byte to a
// string. The returned string will share the underlying memory with the
// []byte which thus allows the string to be mutable through the []byte. We're
// careful to use this method only in situations in which the []byte will not
// be modified.
func UnsafeConvertBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// DecodeUnsafeStringAscending decodes a string value from the input buffer which was
// encoded using EncodeString or EncodeBytes. The r []byte is used as a
// temporary buffer in order to avoid memory allocations. The remainder of the
// input buffer and the decoded string are returned. Note that the returned
// string may share storage with the input buffer.
func DecodeUnsafeStringAscending(b []byte, r []byte) ([]byte, string, error) {
	b, r, err := DecodeBytesAscending(b, r)
	return b, UnsafeConvertBytesToString(r), err
}

// DecodeUnsafeStringAscendingDeepCopy is the same as
// DecodeUnsafeStringAscending but the returned string will never share storage
// with the input buffer.
func DecodeUnsafeStringAscendingDeepCopy(b []byte, r []byte) ([]byte, string, error) {
	b, r, err := DecodeBytesAscendingDeepCopy(b, r)
	return b, UnsafeConvertBytesToString(r), err
}

// DecodeUnsafeStringDescending decodes a string value from the input buffer which
// was encoded using EncodeStringDescending or EncodeBytesDescending. The r
// []byte is used as a temporary buffer in order to avoid memory
// allocations. The remainder of the input buffer and the decoded string are
// returned. Note that the returned string may share storage with the input
// buffer.
func DecodeUnsafeStringDescending(b []byte, r []byte) ([]byte, string, error) {
	b, r, err := DecodeBytesDescending(b, r)
	return b, UnsafeConvertBytesToString(r), err
}

// EncodeNullAscending encodes a NULL value. The encodes bytes are appended to the
// supplied buffer and the final buffer is returned. The encoded value for a
// NULL is guaranteed to not be a prefix for the EncodeVarint, EncodeFloat,
// EncodeBytes and EncodeString encodings.
func EncodeNullAscending(b []byte) []byte {
	return append(b, encodedNull)
}

// EncodeJSONAscending encodes a JSON Type. The encoded bytes are appended to the
// supplied buffer and the final buffer is returned.
func EncodeJSONAscending(b []byte) []byte {
	return append(b, jsonInvertedIndex)
}

// Geo inverted keys are formatted as:
// geoInvertedIndexMarker + EncodeUvarintAscending(cellid) + encoded-bbox
// We don't have a single function to do the whole encoding since a shape is typically
// indexed under multiple cellids, but has a single bbox. So the caller can more
// efficiently
// - append geoInvertedIndex to construct the prefix.
// - encode the bbox once
// - iterate over the cellids and append the encoded cellid to the prefix and then the
//   previously encoded bbox.

// EncodeGeoInvertedAscending appends the geoInvertedIndexMarker.
func EncodeGeoInvertedAscending(b []byte) []byte {
	return append(b, geoInvertedIndexMarker)
}

// Currently only the lowest bit is used to define the encoding kind and the
// remaining 7 bits are unused.
type geoInvertedBBoxEncodingKind byte

const (
	geoInvertedFourFloats geoInvertedBBoxEncodingKind = iota
	geoInvertedTwoFloats
)

// MaxGeoInvertedBBoxLen is the maximum length of the encoded bounding box for
// geo inverted keys.
const MaxGeoInvertedBBoxLen = 1 + 4*uint64AscendingEncodedLength

// EncodeGeoInvertedBBox encodes the bounding box for the geo inverted index.
func EncodeGeoInvertedBBox(b []byte, loX, loY, hiX, hiY float64) []byte {
	encodeTwoFloats := loX == hiX && loY == hiY
	if encodeTwoFloats {
		b = append(b, byte(geoInvertedTwoFloats))
		b = EncodeUntaggedFloatValue(b, loX)
		b = EncodeUntaggedFloatValue(b, loY)
	} else {
		b = append(b, byte(geoInvertedFourFloats))
		b = EncodeUntaggedFloatValue(b, loX)
		b = EncodeUntaggedFloatValue(b, loY)
		b = EncodeUntaggedFloatValue(b, hiX)
		b = EncodeUntaggedFloatValue(b, hiY)
	}
	return b
}

// DecodeGeoInvertedKey decodes the bounding box from the geo inverted key.
// The cellid is skipped in the decoding.
func DecodeGeoInvertedKey(b []byte) (loX, loY, hiX, hiY float64, remaining []byte, err error) {
	// Minimum: 1 byte marker + 1 byte cell length +
	//          1 byte bbox encoding kind + 16 bytes for 2 floats
	if len(b) < 3+2*uint64AscendingEncodedLength {
		return 0, 0, 0, 0, b,
			errors.Errorf("inverted key length %d too small", len(b))
	}
	if b[0] != geoInvertedIndexMarker {
		return 0, 0, 0, 0, b, errors.Errorf("marker is not geoInvertedIndexMarker")
	}
	b = b[1:]
	var cellLen int
	if cellLen, err = getVarintLen(b); err != nil {
		return 0, 0, 0, 0, b, err
	}
	if len(b) < cellLen+17 {
		return 0, 0, 0, 0, b,
			errors.Errorf("insufficient length for encoded bbox in inverted key: %d", len(b)-cellLen)
	}
	encodingKind := geoInvertedBBoxEncodingKind(b[cellLen])
	if encodingKind != geoInvertedTwoFloats && encodingKind != geoInvertedFourFloats {
		return 0, 0, 0, 0, b,
			errors.Errorf("unknown encoding kind for bbox in inverted key: %d", encodingKind)
	}
	b = b[cellLen+1:]
	if b, loX, err = DecodeUntaggedFloatValue(b); err != nil {
		return 0, 0, 0, 0, b, err
	}
	if b, loY, err = DecodeUntaggedFloatValue(b); err != nil {
		return 0, 0, 0, 0, b, err
	}
	if encodingKind == geoInvertedFourFloats {
		if b, hiX, err = DecodeUntaggedFloatValue(b); err != nil {
			return 0, 0, 0, 0, b, err
		}
		if b, hiY, err = DecodeUntaggedFloatValue(b); err != nil {
			return 0, 0, 0, 0, b, err
		}
	} else {
		hiX = loX
		hiY = loY
	}
	return loX, loY, hiX, hiY, b, nil
}

// EncodeNullDescending is the descending equivalent of EncodeNullAscending.
func EncodeNullDescending(b []byte) []byte {
	return append(b, encodedNullDesc)
}

// EncodeNotNullAscending encodes a value that is larger than the NULL marker encoded by
// EncodeNull but less than any encoded value returned by EncodeVarint,
// EncodeFloat, EncodeBytes or EncodeString.
func EncodeNotNullAscending(b []byte) []byte {
	return append(b, encodedNotNull)
}

// EncodeJSONObjectSpanStartAscending encodes the first possible value for JSON
// objects, which is \x00\xff. Non-objects (i.e., scalars and arrays) will
// start with \x00\x01 or \x00\x03 (see AddJSONPathTerminator and
// EncodeArrayAscending), so all objects will be ordered after them.
func EncodeJSONObjectSpanStartAscending(b []byte) []byte {
	return append(b, escape, escaped00)
}

// EncodeArrayAscending encodes a value used to signify membership of an array for JSON objects.
func EncodeArrayAscending(b []byte) []byte {
	return append(b, escape, escapedJSONArray)
}

// EncodeTrueAscending encodes the boolean value true for use with JSON inverted indexes.
func EncodeTrueAscending(b []byte) []byte {
	return append(b, byte(True))
}

// EncodeFalseAscending encodes the boolean value false for use with JSON inverted indexes.
func EncodeFalseAscending(b []byte) []byte {
	return append(b, byte(False))
}

// EncodeNotNullDescending is the descending equivalent of EncodeNotNullAscending.
func EncodeNotNullDescending(b []byte) []byte {
	return append(b, encodedNotNullDesc)
}

// DecodeIfNull decodes a NULL value from the input buffer. If the input buffer
// contains a null at the start of the buffer then it is removed from the
// buffer and true is returned for the second result. Otherwise, the buffer is
// returned unchanged and false is returned for the second result. Since the
// NULL value encoding is guaranteed to never occur as the prefix for the
// EncodeVarint, EncodeFloat, EncodeBytes and EncodeString encodings, it is
// safe to call DecodeIfNull on their encoded values.
// This function handles both ascendingly and descendingly encoded NULLs.
func DecodeIfNull(b []byte) ([]byte, bool) {
	if PeekType(b) == Null {
		return b[1:], true
	}
	return b, false
}

// DecodeIfNotNull decodes a not-NULL value from the input buffer. If the input
// buffer contains a not-NULL marker at the start of the buffer then it is
// removed from the buffer and true is returned for the second
// result. Otherwise, the buffer is returned unchanged and false is returned
// for the second result. Note that the not-NULL marker is identical to the
// empty string encoding, so do not use this routine where it is necessary to
// distinguish not-NULL from the empty string.
// This function handles both ascendingly and descendingly encoded NULLs.
func DecodeIfNotNull(b []byte) ([]byte, bool) {
	if PeekType(b) == NotNull {
		return b[1:], true
	}
	return b, false
}

// EncodeTimeAscending encodes a time value, appends it to the supplied buffer,
// and returns the final buffer. The encoding is guaranteed to be ordered
// Such that if t1.Before(t2) then after EncodeTime(b1, t1), and
// EncodeTime(b2, t2), Compare(b1, b2) < 0. The time zone offset not
// included in the encoding.
func EncodeTimeAscending(b []byte, t time.Time) []byte {
	return encodeTime(b, t.Unix(), int64(t.Nanosecond()))
}

// EncodeTimeDescending is the descending version of EncodeTimeAscending.
func EncodeTimeDescending(b []byte, t time.Time) []byte {
	return encodeTime(b, ^t.Unix(), ^int64(t.Nanosecond()))
}

func encodeTime(b []byte, unix, nanos int64) []byte {
	// Read the unix absolute time. This is the absolute time and is
	// not time zone offset dependent.
	b = append(b, timeMarker)
	b = EncodeVarintAscending(b, unix)
	b = EncodeVarintAscending(b, nanos)
	return b
}

// DecodeTimeAscending decodes a time.Time value which was encoded using
// EncodeTime. The remainder of the input buffer and the decoded
// time.Time are returned.
func DecodeTimeAscending(b []byte) ([]byte, time.Time, error) {
	b, sec, nsec, err := decodeTime(b)
	if err != nil {
		return b, time.Time{}, err
	}
	return b, timeutil.Unix(sec, nsec), nil
}

// DecodeTimeDescending is the descending version of DecodeTimeAscending.
func DecodeTimeDescending(b []byte) ([]byte, time.Time, error) {
	b, sec, nsec, err := decodeTime(b)
	if err != nil {
		return b, time.Time{}, err
	}
	return b, timeutil.Unix(^sec, ^nsec), nil
}

func decodeTime(b []byte) (r []byte, sec int64, nsec int64, err error) {
	if PeekType(b) != Time {
		return nil, 0, 0, errors.Errorf("did not find marker")
	}
	b = b[1:]
	b, sec, err = DecodeVarintAscending(b)
	if err != nil {
		return b, 0, 0, err
	}
	b, nsec, err = DecodeVarintAscending(b)
	if err != nil {
		return b, 0, 0, err
	}
	return b, sec, nsec, nil
}

// EncodeVoidAscendingOrDescending encodes a void (valid for both ascending and descending order).
func EncodeVoidAscendingOrDescending(b []byte) []byte {
	return append(b, voidMarker)
}

// DecodeVoidAscendingOrDescending decodes a void  (valid for both ascending and descending order).
func DecodeVoidAscendingOrDescending(b []byte) ([]byte, error) {
	if PeekType(b) != Void {
		return nil, errors.Errorf("did not find Void marker")
	}
	return b[1:], nil
}

// EncodeBox2DAscending encodes a bounding box in ascending order.
func EncodeBox2DAscending(b []byte, box geopb.BoundingBox) ([]byte, error) {
	b = append(b, box2DMarker)
	b = EncodeFloatAscending(b, box.LoX)
	b = EncodeFloatAscending(b, box.HiX)
	b = EncodeFloatAscending(b, box.LoY)
	b = EncodeFloatAscending(b, box.HiY)
	return b, nil
}

// EncodeBox2DDescending encodes a bounding box in descending order.
func EncodeBox2DDescending(b []byte, box geopb.BoundingBox) ([]byte, error) {
	b = append(b, box2DMarker)
	b = EncodeFloatDescending(b, box.LoX)
	b = EncodeFloatDescending(b, box.HiX)
	b = EncodeFloatDescending(b, box.LoY)
	b = EncodeFloatDescending(b, box.HiY)
	return b, nil
}

// DecodeBox2DAscending decodes a box2D object in ascending order.
func DecodeBox2DAscending(b []byte) ([]byte, geopb.BoundingBox, error) {
	box := geopb.BoundingBox{}
	if PeekType(b) != Box2D {
		return nil, box, errors.Errorf("did not find Box2D marker")
	}

	b = b[1:]
	var err error
	b, box.LoX, err = DecodeFloatAscending(b)
	if err != nil {
		return nil, box, err
	}
	b, box.HiX, err = DecodeFloatAscending(b)
	if err != nil {
		return nil, box, err
	}
	b, box.LoY, err = DecodeFloatAscending(b)
	if err != nil {
		return nil, box, err
	}
	b, box.HiY, err = DecodeFloatAscending(b)
	if err != nil {
		return nil, box, err
	}
	return b, box, nil
}

// DecodeBox2DDescending decodes a box2D object in descending order.
func DecodeBox2DDescending(b []byte) ([]byte, geopb.BoundingBox, error) {
	box := geopb.BoundingBox{}
	if PeekType(b) != Box2D {
		return nil, box, errors.Errorf("did not find Box2D marker")
	}

	b = b[1:]
	var err error
	b, box.LoX, err = DecodeFloatDescending(b)
	if err != nil {
		return nil, box, err
	}
	b, box.HiX, err = DecodeFloatDescending(b)
	if err != nil {
		return nil, box, err
	}
	b, box.LoY, err = DecodeFloatDescending(b)
	if err != nil {
		return nil, box, err
	}
	b, box.HiY, err = DecodeFloatDescending(b)
	if err != nil {
		return nil, box, err
	}
	return b, box, nil
}

// EncodeGeoAscending encodes a geopb.SpatialObject value in ascending order and
// returns the new buffer.
// It is sorted by the given curve index, followed by the bytes of the spatial object.
func EncodeGeoAscending(b []byte, curveIndex uint64, so *geopb.SpatialObject) ([]byte, error) {
	b = append(b, geoMarker)
	b = EncodeUint64Ascending(b, curveIndex)

	data, err := protoutil.Marshal(so)
	if err != nil {
		return nil, err
	}
	b = encodeBytesAscendingWithTerminator(b, data, ascendingGeoEscapes.escapedTerm)
	return b, nil
}

// EncodeGeoDescending encodes a geopb.SpatialObject value in descending order and
// returns the new buffer.
// It is sorted by the given curve index, followed by the bytes of the spatial object.
func EncodeGeoDescending(b []byte, curveIndex uint64, so *geopb.SpatialObject) ([]byte, error) {
	b = append(b, geoDescMarker)
	b = EncodeUint64Descending(b, curveIndex)

	data, err := protoutil.Marshal(so)
	if err != nil {
		return nil, err
	}
	n := len(b)
	b = encodeBytesAscendingWithTerminator(b, data, ascendingGeoEscapes.escapedTerm)
	onesComplement(b[n:])
	return b, nil
}

// DecodeGeoAscending decodes a geopb.SpatialObject value that was encoded
// in ascending order back into a geopb.SpatialObject. The so parameter
// must already be empty/reset.
func DecodeGeoAscending(b []byte, so *geopb.SpatialObject) ([]byte, error) {
	if PeekType(b) != Geo {
		return nil, errors.Errorf("did not find Geo marker")
	}
	b = b[1:]
	var err error
	b, _, err = DecodeUint64Ascending(b)
	if err != nil {
		return nil, err
	}

	var pbBytes []byte
	b, pbBytes, err = decodeBytesInternal(b, pbBytes, ascendingGeoEscapes, false /* expectMarker */, false /* deepCopy */)
	if err != nil {
		return b, err
	}
	// Not using protoutil.Unmarshal since the call to so.Reset() will waste the
	// pre-allocated EWKB.
	err = so.Unmarshal(pbBytes)
	return b, err
}

// DecodeGeoDescending decodes a geopb.SpatialObject value that was encoded
// in descending order back into a geopb.SpatialObject. The so parameter
// must already be empty/reset.
func DecodeGeoDescending(b []byte, so *geopb.SpatialObject) ([]byte, error) {
	if PeekType(b) != GeoDesc {
		return nil, errors.Errorf("did not find Geo marker")
	}
	b = b[1:]
	var err error
	b, _, err = DecodeUint64Descending(b)
	if err != nil {
		return nil, err
	}

	var pbBytes []byte
	b, pbBytes, err = decodeBytesInternal(b, pbBytes, descendingGeoEscapes, false /* expectMarker */, false /* deepCopy */)
	if err != nil {
		return b, err
	}
	onesComplement(pbBytes)
	// Not using protoutil.Unmarshal since the call to so.Reset() will waste the
	// pre-allocated EWKB.
	err = so.Unmarshal(pbBytes)
	return b, err
}

// EncodeTimeTZAscending encodes a timetz.TimeTZ value and appends it to
// the supplied buffer and returns the final buffer.
// The encoding is guaranteed to be ordered such that if t1.Before(t2)
// then after encodeTimeTZ(b1, t1) and encodeTimeTZ(b2, t2),
// Compare(b1, b2) < 0.
// The time zone offset is included in the encoding.
func EncodeTimeTZAscending(b []byte, t timetz.TimeTZ) []byte {
	// Do not use TimeOfDay's add function, as it loses 24:00:00 encoding.
	return encodeTimeTZ(b, int64(t.TimeOfDay)+int64(t.OffsetSecs)*offsetSecsToMicros, t.OffsetSecs)
}

// EncodeTimeTZDescending is the descending version of EncodeTimeTZAscending.
func EncodeTimeTZDescending(b []byte, t timetz.TimeTZ) []byte {
	// Do not use TimeOfDay's add function, as it loses 24:00:00 encoding.
	return encodeTimeTZ(b, ^(int64(t.TimeOfDay) + int64(t.OffsetSecs)*offsetSecsToMicros), ^t.OffsetSecs)
}

func encodeTimeTZ(b []byte, unixMicros int64, offsetSecs int32) []byte {
	b = append(b, timeTZMarker)
	b = EncodeVarintAscending(b, unixMicros)
	b = EncodeVarintAscending(b, int64(offsetSecs))
	return b
}

// DecodeTimeTZAscending decodes a timetz.TimeTZ value which was encoded
// using encodeTimeTZ. The remainder of the input buffer and the decoded
// timetz.TimeTZ are returned.
func DecodeTimeTZAscending(b []byte) ([]byte, timetz.TimeTZ, error) {
	b, unixMicros, offsetSecs, err := decodeTimeTZ(b)
	if err != nil {
		return nil, timetz.TimeTZ{}, err
	}
	// Do not use timeofday.FromInt, as it loses 24:00:00 encoding.
	return b, timetz.TimeTZ{
		TimeOfDay:  timeofday.TimeOfDay(unixMicros - int64(offsetSecs)*offsetSecsToMicros),
		OffsetSecs: offsetSecs,
	}, nil
}

// DecodeTimeTZDescending is the descending version of DecodeTimeTZAscending.
func DecodeTimeTZDescending(b []byte) ([]byte, timetz.TimeTZ, error) {
	b, unixMicros, offsetSecs, err := decodeTimeTZ(b)
	if err != nil {
		return nil, timetz.TimeTZ{}, err
	}
	// Do not use timeofday.FromInt, as it loses 24:00:00 encoding.
	return b, timetz.TimeTZ{
		TimeOfDay:  timeofday.TimeOfDay(^unixMicros - int64(^offsetSecs)*offsetSecsToMicros),
		OffsetSecs: ^offsetSecs,
	}, nil
}

func decodeTimeTZ(b []byte) ([]byte, int64, int32, error) {
	if PeekType(b) != TimeTZ {
		return nil, 0, 0, errors.Errorf("did not find marker")
	}
	b = b[1:]
	var err error
	var unixMicros int64
	b, unixMicros, err = DecodeVarintAscending(b)
	if err != nil {
		return nil, 0, 0, err
	}
	var offsetSecs int64
	b, offsetSecs, err = DecodeVarintAscending(b)
	if err != nil {
		return nil, 0, 0, err
	}
	return b, unixMicros, int32(offsetSecs), nil
}

// EncodeDurationAscending encodes a duration.Duration value, appends it to the
// supplied buffer, and returns the final buffer. The encoding is guaranteed to
// be ordered such that if t1.Compare(t2) < 0 (or = 0 or > 0) then bytes.Compare
// will order them the same way after encoding.
func EncodeDurationAscending(b []byte, d duration.Duration) ([]byte, error) {
	sortNanos, months, days, err := d.Encode()
	if err != nil {
		// TODO(dan): Handle this using d.EncodeBigInt() and the
		// durationBigNeg/durationBigPos markers.
		return b, err
	}
	b = append(b, durationMarker)
	b = EncodeVarintAscending(b, sortNanos)
	b = EncodeVarintAscending(b, months)
	b = EncodeVarintAscending(b, days)
	return b, nil
}

// EncodeDurationDescending is the descending version of EncodeDurationAscending.
func EncodeDurationDescending(b []byte, d duration.Duration) ([]byte, error) {
	sortNanos, months, days, err := d.Encode()
	if err != nil {
		// TODO(dan): Handle this using d.EncodeBigInt() and the
		// durationBigNeg/durationBigPos markers.
		return b, err
	}
	b = append(b, durationMarker)
	b = EncodeVarintDescending(b, sortNanos)
	b = EncodeVarintDescending(b, months)
	b = EncodeVarintDescending(b, days)
	return b, nil
}

// DecodeDurationAscending decodes a duration.Duration value which was encoded
// using EncodeDurationAscending. The remainder of the input buffer and the
// decoded duration.Duration are returned.
func DecodeDurationAscending(b []byte) ([]byte, duration.Duration, error) {
	if PeekType(b) != Duration {
		return nil, duration.Duration{}, errors.Errorf("did not find marker %x", b)
	}
	b = b[1:]
	b, sortNanos, err := DecodeVarintAscending(b)
	if err != nil {
		return b, duration.Duration{}, err
	}
	b, months, err := DecodeVarintAscending(b)
	if err != nil {
		return b, duration.Duration{}, err
	}
	b, days, err := DecodeVarintAscending(b)
	if err != nil {
		return b, duration.Duration{}, err
	}
	d, err := duration.Decode(sortNanos, months, days)
	if err != nil {
		return b, duration.Duration{}, err
	}
	return b, d, nil
}

// DecodeDurationDescending is the descending version of DecodeDurationAscending.
func DecodeDurationDescending(b []byte) ([]byte, duration.Duration, error) {
	if PeekType(b) != Duration {
		return nil, duration.Duration{}, errors.Errorf("did not find marker")
	}
	b = b[1:]
	b, sortNanos, err := DecodeVarintDescending(b)
	if err != nil {
		return b, duration.Duration{}, err
	}
	b, months, err := DecodeVarintDescending(b)
	if err != nil {
		return b, duration.Duration{}, err
	}
	b, days, err := DecodeVarintDescending(b)
	if err != nil {
		return b, duration.Duration{}, err
	}
	d, err := duration.Decode(sortNanos, months, days)
	if err != nil {
		return b, duration.Duration{}, err
	}
	return b, d, nil
}

// EncodeBitArrayAscending encodes a bitarray.BitArray value, appends it to the
// supplied buffer, and returns the final buffer. The encoding is guaranteed to
// be ordered such that if t1.Compare(t2) < 0 (or = 0 or > 0) then bytes.Compare
// will order them the same way after encoding.
//
// The encoding uses varint encoding for each word of the backing
// array. This is a trade-off. The alternative is to encode the entire
// backing word array as a byte array, using byte array encoding and escaped
// special bytes (via  `encodeBytesAscendingWithoutTerminatorOrPrefix`).
// There are two arguments against this alternative:
//   - the bytes must be encoded big endian, but the most common architectures
//     running CockroachDB are little-endian, so the bytes would need
//     to be reordered prior to encoding.
//   - when decoding or skipping over a value, the decoding/sizing loop
//     would need to look at every byte of the encoding to find the
//     terminator.
//
// In contrast, the chosen encoding using varints is endianness-agnostic
// and enables fast decoding/skipping thanks ot the tag bytes.
func EncodeBitArrayAscending(b []byte, d bitarray.BitArray) []byte {
	b = append(b, bitArrayMarker)
	words, lastBitsUsed := d.EncodingParts()
	for _, w := range words {
		b = EncodeUvarintAscending(b, w)
	}
	b = append(b, bitArrayDataTerminator)
	b = EncodeUvarintAscending(b, lastBitsUsed)
	return b
}

// EncodeBitArrayDescending is the descending version of EncodeBitArrayAscending.
func EncodeBitArrayDescending(b []byte, d bitarray.BitArray) []byte {
	b = append(b, bitArrayDescMarker)
	words, lastBitsUsed := d.EncodingParts()
	for _, w := range words {
		b = EncodeUvarintDescending(b, w)
	}
	b = append(b, bitArrayDataDescTerminator)
	b = EncodeUvarintDescending(b, lastBitsUsed)
	return b
}

// DecodeBitArrayAscending decodes a bit array which was encoded using
// EncodeBitArrayAscending. The remainder of the input buffer and the
// decoded bit array are returned.
func DecodeBitArrayAscending(b []byte) ([]byte, bitarray.BitArray, error) {
	if PeekType(b) != BitArray {
		return nil, bitarray.BitArray{}, errors.Errorf("did not find marker %x", b)
	}
	b = b[1:]

	// First compute the length.
	numWords, _, err := getBitArrayWordsLen(b, bitArrayDataTerminator)
	if err != nil {
		return b, bitarray.BitArray{}, err
	}
	// Decode the words.
	words := make([]uint64, numWords)
	for i := range words {
		b, words[i], err = DecodeUvarintAscending(b)
		if err != nil {
			return b, bitarray.BitArray{}, err
		}
	}
	// Decode the final part.
	if len(b) == 0 || b[0] != bitArrayDataTerminator {
		return b, bitarray.BitArray{}, errBitArrayTerminatorMissing
	}
	b = b[1:]
	b, lastVal, err := DecodeUvarintAscending(b)
	if err != nil {
		return b, bitarray.BitArray{}, err
	}
	ba, err := bitarray.FromEncodingParts(words, lastVal)
	return b, ba, err
}

var errBitArrayTerminatorMissing = errors.New("cannot find bit array data terminator")

// getBitArrayWordsLen returns the number of bit array words in the
// encoded bytes and the size in bytes of the encoded word array
// (excluding the terminator byte).
func getBitArrayWordsLen(b []byte, term byte) (int, int, error) {
	bSearch := b
	numWords := 0
	sz := 0
	for {
		if len(bSearch) == 0 {
			return 0, 0, errors.Errorf("slice too short for bit array (%d)", len(b))
		}
		if bSearch[0] == term {
			break
		}
		vLen, err := getVarintLen(bSearch)
		if err != nil {
			return 0, 0, err
		}
		bSearch = bSearch[vLen:]
		numWords++
		sz += vLen
	}
	return numWords, sz, nil
}

// DecodeBitArrayDescending is the descending version of DecodeBitArrayAscending.
func DecodeBitArrayDescending(b []byte) ([]byte, bitarray.BitArray, error) {
	if PeekType(b) != BitArrayDesc {
		return nil, bitarray.BitArray{}, errors.Errorf("did not find marker %x", b)
	}
	b = b[1:]

	// First compute the length.
	numWords, _, err := getBitArrayWordsLen(b, bitArrayDataDescTerminator)
	if err != nil {
		return b, bitarray.BitArray{}, err
	}
	// Decode the words.
	words := make([]uint64, numWords)
	for i := range words {
		b, words[i], err = DecodeUvarintDescending(b)
		if err != nil {
			return b, bitarray.BitArray{}, err
		}
	}
	// Decode the final part.
	if len(b) == 0 || b[0] != bitArrayDataDescTerminator {
		return b, bitarray.BitArray{}, errBitArrayTerminatorMissing
	}
	b = b[1:]
	b, lastVal, err := DecodeUvarintDescending(b)
	if err != nil {
		return b, bitarray.BitArray{}, err
	}
	ba, err := bitarray.FromEncodingParts(words, lastVal)
	return b, ba, err
}

// Type represents the type of a value encoded by
// Encode{Null,NotNull,Varint,Uvarint,Float,Bytes}.
//
//go:generate stringer -type=Type
type Type encodingtype.T

// Type values.
// TODO(dan, arjun): Make this into a proto enum.
// The 'Type' annotations are necessary for producing stringer-generated values.
const (
	Unknown   Type = 0
	Null      Type = 1
	NotNull   Type = 2
	Int       Type = 3
	Float     Type = 4
	Decimal   Type = 5
	Bytes     Type = 6
	BytesDesc Type = 7 // Bytes encoded descendingly
	Time      Type = 8
	Duration  Type = 9
	True      Type = 10
	False     Type = 11
	UUID      Type = 12
	Array     Type = 13
	IPAddr    Type = 14
	// SentinelType is used for bit manipulation to check if the encoded type
	// value requires more than 4 bits, and thus will be encoded in two bytes. It
	// is not used as a type value, and thus intentionally overlaps with the
	// subsequent type value. The 'Type' annotation is intentionally omitted here.
	SentinelType        = 15
	JSON           Type = 15
	Tuple          Type = 16
	BitArray       Type = 17
	BitArrayDesc   Type = 18 // BitArray encoded descendingly
	TimeTZ         Type = 19
	Geo            Type = 20
	GeoDesc        Type = 21
	ArrayKeyAsc    Type = 22 // Array key encoding
	ArrayKeyDesc   Type = 23 // Array key encoded descendingly
	Box2D          Type = 24
	Void           Type = 25
	TSQuery        Type = 26
	TSVector       Type = 27
	JSONNull       Type = 28
	JSONNullDesc   Type = 29
	JSONString     Type = 30
	JSONStringDesc Type = 31
	JSONNumber     Type = 32
	JSONNumberDesc Type = 33
	JSONFalse      Type = 34
	JSONFalseDesc  Type = 35
	JSONTrue       Type = 36
	JSONTrueDesc   Type = 37
	JSONArray      Type = 38
	JSONArrayDesc  Type = 39
	JSONObject     Type = 40
	JSONObjectDesc Type = 41
	// Special case
	JsonEmptyArray     Type = 42
	JsonEmptyArrayDesc Type = 43
	PGVector           Type = 44
)

// typMap maps an encoded type byte to a decoded Type. It's got 256 slots, one
// for every possible byte value.
var typMap [256]Type

func init() {
	buf := []byte{0}
	for i := range typMap {
		buf[0] = byte(i)
		typMap[i] = slowPeekType(buf)
	}
}

// PeekType peeks at the type of the value encoded at the start of b.
func PeekType(b []byte) Type {
	if len(b) >= 1 {
		return typMap[b[0]]
	}
	return Unknown
}

// slowPeekType is the old implementation of PeekType. It's used to generate
// the lookup table for PeekType.
func slowPeekType(b []byte) Type {
	if len(b) >= 1 {
		m := b[0]
		switch {
		case m == encodedNull, m == encodedNullDesc:
			return Null
		case m == encodedNotNull, m == encodedNotNullDesc:
			return NotNull
		case m == arrayKeyMarker:
			return ArrayKeyAsc
		case m == arrayKeyDescendingMarker:
			return ArrayKeyDesc
		case m == jsonNullKeyMarker:
			return JSONNull
		case m == jsonNullKeyDescendingMarker:
			return JSONNullDesc
		case m == jsonStringKeyMarker:
			return JSONString
		case m == jsonStringKeyDescendingMarker:
			return JSONStringDesc
		case m == jsonNumberKeyMarker:
			return JSONNumber
		case m == jsonNumberKeyDescendingMarker:
			return JSONNumberDesc
		case m == jsonFalseKeyMarker:
			return JSONFalse
		case m == jsonFalseKeyDescendingMarker:
			return JSONFalseDesc
		case m == jsonTrueKeyMarker:
			return JSONTrue
		case m == jsonTrueKeyDescendingMarker:
			return JSONTrueDesc
		case m == jsonArrayKeyMarker:
			return JSONArray
		case m == jsonArrayKeyDescendingMarker:
			return JSONArrayDesc
		case m == jsonEmptyArrayKeyMarker:
			return JsonEmptyArray
		case m == jsonEmptyArrayKeyDescendingMarker:
			return JsonEmptyArrayDesc
		case m == jsonObjectKeyMarker:
			return JSONObject
		case m == jsonObjectKeyDescendingMarker:
			return JSONObjectDesc
		case m == bytesMarker:
			return Bytes
		case m == bytesDescMarker:
			return BytesDesc
		case m == bitArrayMarker:
			return BitArray
		case m == bitArrayDescMarker:
			return BitArrayDesc
		case m == timeMarker:
			return Time
		case m == timeTZMarker:
			return TimeTZ
		case m == geoMarker:
			return Geo
		case m == box2DMarker:
			return Box2D
		case m == geoDescMarker:
			return GeoDesc
		case m == byte(Array):
			return Array
		case m == byte(True):
			return True
		case m == byte(False):
			return False
		case m == durationBigNegMarker, m == durationMarker, m == durationBigPosMarker:
			return Duration
		case m >= IntMin && m <= IntMax:
			return Int
		case m >= floatNaN && m <= floatNaNDesc:
			return Float
		case m >= decimalNaN && m <= decimalNaNDesc:
			return Decimal
		case m == voidMarker:
			return Void
		}
	}
	return Unknown
}

// GetMultiVarintLen find the length of <num> encoded varints that follow a
// 1-byte tag.
func GetMultiVarintLen(b []byte, num int) (int, error) {
	p := 1
	for i := 0; i < num && p < len(b); i++ {
		len, err := getVarintLen(b[p:])
		if err != nil {
			return 0, err
		}
		p += len
	}
	return p, nil
}

// getMultiNonsortingVarintLen finds the length of <num> encoded nonsorting varints.
func getMultiNonsortingVarintLen(b []byte, num int) (int, error) {
	p := 0
	for i := 0; i < num && p < len(b); i++ {
		_, len, _, err := DecodeNonsortingStdlibVarint(b[p:])
		if err != nil {
			return 0, err
		}
		p += len
	}
	return p, nil
}

func getArrayOrJSONLength(
	buf []byte, dir Direction, keyDoneFn func(buf []byte, dir Direction) bool,
) (int, error) {
	result := 0

	for {
		if len(buf) == 0 {
			return 0, errors.AssertionFailedf("invalid encoding (unterminated)")
		}
		if keyDoneFn(buf, dir) {
			// Increment to include the terminator byte.
			result++
			break
		}
		next, err := PeekLength(buf)
		if err != nil {
			return 0, err
		}
		// Shift buf over by the encoded data amount.
		buf = buf[next:]
		result += next
	}
	return result, nil
}

// peekBox2DLength peeks to look at the length of a box2d encoding.
func peekBox2DLength(b []byte) (int, error) {
	length := 0
	curr := b
	for i := 0; i < 4; i++ {
		if len(curr) == 0 {
			return 0, errors.Newf("slice too short for box2d")
		}
		switch curr[0] {
		case floatNaN, floatNaNDesc, floatZero:
			length++
			curr = curr[1:]
		case floatNeg, floatPos:
			length += 9
			curr = curr[9:]
		default:
			return 0, errors.Newf("unexpected marker for box2d: %x", curr[0])
		}
	}
	return length, nil
}

// PeekLength returns the length of the encoded value at the start of b.  Note:
// if this function succeeds, it's not a guarantee that decoding the value will
// succeed. PeekLength is meant to be used on key encoded data only.
func PeekLength(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, errors.Errorf("empty slice")
	}
	m := b[0]
	switch m {
	case encodedNull, encodedNullDesc, encodedNotNull, encodedNotNullDesc,
		floatNaN, floatNaNDesc, floatZero, decimalZero, byte(True), byte(False),
		emptyArray, voidMarker, jsonNullKeyMarker, jsonNullKeyDescendingMarker,
		jsonFalseKeyMarker, jsonFalseKeyDescendingMarker, jsonTrueKeyMarker,
		jsonTrueKeyDescendingMarker:
		// ascendingNullWithinArrayKey and descendingNullWithinArrayKey also
		// contain the same byte values as encodedNotNull and encodedNotNullDesc
		// respectively, but they cannot be included explicitly in the case
		// statement.
		return 1, nil
	case bitArrayMarker, bitArrayDescMarker:
		terminator := byte(bitArrayDataTerminator)
		if m == bitArrayDescMarker {
			terminator = bitArrayDataDescTerminator
		}
		_, n, err := getBitArrayWordsLen(b[1:], terminator)
		if err != nil {
			return 1 + n, err
		}
		m, err := getVarintLen(b[n+2:])
		if err != nil {
			return 1 + n + m + 1, err
		}
		return 1 + n + m + 1, nil
	case jsonStringKeyMarker, jsonStringKeyDescendingMarker,
		jsonNumberKeyMarker, jsonNumberKeyDescendingMarker:
		dir := Ascending
		if (m == jsonStringKeyDescendingMarker) ||
			(m == jsonNumberKeyDescendingMarker) {
			dir = Descending
		}
		length, err := getArrayOrJSONLength(b[1:], dir, IsJSONKeyDone)
		return 1 + length, err
	case jsonArrayKeyMarker, jsonArrayKeyDescendingMarker,
		jsonObjectKeyMarker, jsonObjectKeyDescendingMarker,
		jsonEmptyArrayKeyMarker, jsonEmptyArrayKeyDescendingMarker:
		dir := Ascending
		if (m == jsonArrayKeyDescendingMarker) ||
			(m == jsonObjectKeyDescendingMarker) ||
			(m == jsonEmptyArrayKeyDescendingMarker) {
			dir = Descending
		}
		// removing the starter tag
		b = b[1:]

		// Getting the number of elements present
		// in the container.
		numberElems, err := getVarintLen(b)
		if err != nil {
			return -1, errors.AssertionFailedf("failed to get the number of elements" +
				"in the container")
		}
		length, err := getArrayOrJSONLength(b[numberElems:], dir, IsJSONKeyDone)
		return 1 + numberElems + length, err
	case arrayKeyMarker, arrayKeyDescendingMarker:
		dir := Ascending
		if m == arrayKeyDescendingMarker {
			dir = Descending
		}
		length, err := getArrayOrJSONLength(b[1:], dir, IsArrayKeyDone)
		return 1 + length, err
	case bytesMarker:
		return getBytesLength(b, ascendingBytesEscapes)
	case box2DMarker:
		if len(b) == 0 {
			return 0, errors.Newf("slice too short for box2d")
		}
		length, err := peekBox2DLength(b[1:])
		if err != nil {
			return 0, err
		}
		return 1 + length, nil
	case geoInvertedIndexMarker:
		return getGeoInvertedIndexKeyLength(b)
	case geoMarker:
		// Expect to reserve at least 8 bytes for int64.
		if len(b) < 8 {
			return 0, errors.Errorf("slice too short for spatial object (%d)", len(b))
		}
		ret, err := getBytesLength(b[8:], ascendingGeoEscapes)
		if err != nil {
			return 0, err
		}
		return 8 + ret, nil
	case jsonInvertedIndex:
		return getJSONInvertedIndexKeyLength(b)
	case bytesDescMarker:
		return getBytesLength(b, descendingBytesEscapes)
	case geoDescMarker:
		// Expect to reserve at least 8 bytes for int64.
		if len(b) < 8 {
			return 0, errors.Errorf("slice too short for spatial object (%d)", len(b))
		}
		ret, err := getBytesLength(b[8:], descendingGeoEscapes)
		if err != nil {
			return 0, err
		}
		return 8 + ret, nil
	case timeMarker, timeTZMarker:
		return GetMultiVarintLen(b, 2)
	case durationBigNegMarker, durationMarker, durationBigPosMarker:
		return GetMultiVarintLen(b, 3)
	case floatNeg, floatPos:
		// the marker is followed by 8 bytes
		if len(b) < 9 {
			return 0, errors.Errorf("slice too short for float (%d)", len(b))
		}
		return 9, nil
	}
	if m >= IntMin && m <= IntMax {
		return getVarintLen(b)
	}
	if m >= decimalNaN && m <= decimalNaNDesc {
		return getDecimalLen(b)
	}
	return 0, errors.Errorf("unknown tag %d", m)
}

// PrettyPrintValue returns the string representation of all contiguous
// decodable values in the provided byte slice, separated by a provided
// separator.
// The directions each value is encoded may be provided. If valDirs is nil,
// all values are decoded and printed with the default direction (ascending).
func PrettyPrintValue(buf *redact.StringBuilder, valDirs []Direction, b []byte, sep string) {
	safeSep := redact.SafeString(sep)
	allDecoded := prettyPrintValueImpl(buf, valDirs, b, safeSep)
	if allDecoded {
		return
	}
	// If we failed to decoded everything above, assume the key was the result of a
	// `PrefixEnd()`. Attempt to undo PrefixEnd & retry the process, otherwise return
	// what we were able to decode.
	if undoPrefixEnd, ok := UndoPrefixEnd(b); ok {
		// When we UndoPrefixEnd, we may have lost a tail of 0xFFs. Try to add
		// enough of them to get something decoded. This is best-effort, we have to stop
		// somewhere.
		cap := 20
		if len(valDirs) > len(b) {
			cap = len(valDirs) - len(b)
		}
		for i := 0; i < cap; i++ {
			if allDecoded := prettyPrintValueImpl(buf, valDirs, undoPrefixEnd, safeSep); allDecoded {
				buf.Reset()
				buf.Print(sep + "PrefixEnd")
				return
			}
			undoPrefixEnd = append(undoPrefixEnd, 0xFF)
		}
	}
}

// PrettyPrintValuesWithTypes returns a slice containing each contiguous decodable value
// in the provided byte slice along with a slice containing the type of each value.
// The directions each value is encoded may be provided. If valDirs is nil,
// all values are decoded and printed with the default direction (ascending).
func PrettyPrintValuesWithTypes(valDirs []Direction, b []byte) (vals []string, types []Type) {
	vals1, types1, allDecoded := prettyPrintValuesWithTypesImpl(valDirs, b)
	if allDecoded {
		return vals1, types1
	}
	// If we failed to decoded everything above, assume the key was the result of a
	// `PrefixEnd()`. Attempt to undo PrefixEnd & retry the process, otherwise return
	// what we were able to decode.
	if undoPrefixEnd, ok := UndoPrefixEnd(b); ok {
		// When we UndoPrefixEnd, we may have lost a tail of 0xFFs. Try to add
		// enough of them to get something decoded. This is best-effort, we have to stop
		// somewhere.
		cap := 20
		if len(valDirs) > len(b) {
			cap = len(valDirs) - len(b)
		}
		for i := 0; i < cap; i++ {
			if vals2, types2, allDecoded := prettyPrintValuesWithTypesImpl(valDirs, undoPrefixEnd); allDecoded {
				vals2 = append(vals2, "PrefixEnd")
				types2 = append(types2, Bytes)
				return vals2, types2
			}
			undoPrefixEnd = append(undoPrefixEnd, 0xFF)
		}
	}
	return vals1, types1
}

func prettyPrintValuesWithTypesImpl(
	valDirs []Direction, b []byte,
) (vals []string, types []Type, allDecoded bool) {
	allDecoded = true
	for len(b) > 0 {
		var valDir Direction
		if len(valDirs) > 0 {
			valDir = valDirs[0]
			valDirs = valDirs[1:]
		}

		bb, s, err := prettyPrintFirstValue(valDir, b)
		if err != nil {
			// If we fail to decode, mark as unknown and attempt
			// to continue - it's possible we can still decode the
			// remainder of the key bytes.
			allDecoded = false
			vals = append(vals, "???")
			types = append(types, Unknown)
		} else {
			vals = append(vals, s)
			types = append(types, PeekType(b))
		}
		b = bb
	}
	return vals, types, allDecoded
}

func prettyPrintValueImpl(
	buf *redact.StringBuilder, valDirs []Direction, b []byte, sep redact.SafeString,
) bool {
	allDecoded := true
	for len(b) > 0 {
		// If there are more values than encoding directions specified,
		// valDir will contain the 0 value of Direction.
		// prettyPrintFirstValue will then use the default encoding
		// direction per each value type.
		var valDir Direction
		if len(valDirs) > 0 {
			valDir = valDirs[0]
			valDirs = valDirs[1:]
		}

		bb, s, err := prettyPrintFirstValue(valDir, b)
		if err != nil {
			// If we fail to decode, mark as unknown and attempt
			// to continue - it's possible we can still decode the
			// remainder of the key bytes.
			allDecoded = false
			// Mark the separator as safe.
			buf.Print(sep)
			buf.SafeString("???")
		} else {
			buf.Print(sep)
			buf.Print(redact.Safe(s))
		}
		b = bb
	}
	return allDecoded
}

// prettyPrintFirstValue returns a string representation of the first decodable
// value in the provided byte slice, along with the remaining byte slice
// after decoding.
//
// Ascending will be the default direction (when dir is the 0 value) for all
// values.
func prettyPrintFirstValue(dir Direction, b []byte) ([]byte, string, error) {
	var err error
	switch typ := PeekType(b); typ {
	case Null:
		b, _ = DecodeIfNull(b)
		return b, "NULL", nil
	case True:
		return b[1:], "True", nil
	case False:
		return b[1:], "False", nil
	case Array:
		return b[1:], "Arr", nil
	case ArrayKeyAsc, ArrayKeyDesc:
		encDir := Ascending
		if typ == ArrayKeyDesc {
			encDir = Descending
		}
		var build strings.Builder
		buf, err := ValidateAndConsumeArrayKeyMarker(b, encDir)
		if err != nil {
			return nil, "", err
		}
		build.WriteString("ARRAY[")
		first := true
		// Use the array key decoding logic, but instead of calling out
		// to keyside.Decode, just make a recursive call.
		for {
			if len(buf) == 0 {
				return nil, "", errors.AssertionFailedf("invalid array (unterminated)")
			}
			if IsArrayKeyDone(buf, encDir) {
				buf = buf[1:]
				break
			}
			var next string
			if IsNextByteArrayEncodedNull(buf, dir) {
				next = "NULL"
				buf = buf[1:]
			} else {
				buf, next, err = prettyPrintFirstValue(dir, buf)
				if err != nil {
					return nil, "", err
				}
			}
			if !first {
				build.WriteString(",")
			}
			build.WriteString(next)
			first = false
		}
		build.WriteString("]")
		return buf, build.String(), nil
	case NotNull:
		b, _ = DecodeIfNotNull(b)
		return b, "!NULL", nil
	case Int:
		var i int64
		if dir == Descending {
			b, i, err = DecodeVarintDescending(b)
		} else {
			b, i, err = DecodeVarintAscending(b)
		}
		if err != nil {
			return b, "", err
		}
		return b, strconv.FormatInt(i, 10), nil
	case Float:
		var f float64
		if dir == Descending {
			b, f, err = DecodeFloatDescending(b)
		} else {
			b, f, err = DecodeFloatAscending(b)
		}
		if err != nil {
			return b, "", err
		}
		return b, strconv.FormatFloat(f, 'g', -1, 64), nil
	case Decimal:
		var d apd.Decimal
		if dir == Descending {
			b, d, err = DecodeDecimalDescending(b, nil)
		} else {
			b, d, err = DecodeDecimalAscending(b, nil)
		}
		if err != nil {
			return b, "", err
		}
		return b, d.String(), nil
	case BitArray:
		if dir == Descending {
			return b, "", errors.Errorf("descending bit column dir but ascending bit array encoding")
		}
		var d bitarray.BitArray
		b, d, err = DecodeBitArrayAscending(b)
		return b, "B" + d.String(), err
	case BitArrayDesc:
		if dir == Ascending {
			return b, "", errors.Errorf("ascending bit column dir but descending bit array encoding")
		}
		var d bitarray.BitArray
		b, d, err = DecodeBitArrayDescending(b)
		return b, "B" + d.String(), err
	case Bytes:
		if dir == Descending {
			return b, "", errors.Errorf("descending bytes column dir but ascending bytes encoding")
		}
		var s string
		b, s, err = DecodeUnsafeStringAscending(b, nil)
		if err != nil {
			return b, "", err
		}
		return b, strconv.Quote(s), nil
	case BytesDesc:
		if dir == Ascending {
			return b, "", errors.Errorf("ascending bytes column dir but descending bytes encoding")
		}

		var s string
		b, s, err = DecodeUnsafeStringDescending(b, nil)
		if err != nil {
			return b, "", err
		}
		return b, strconv.Quote(s), nil
	case Time:
		var t time.Time
		if dir == Descending {
			b, t, err = DecodeTimeDescending(b)
		} else {
			b, t, err = DecodeTimeAscending(b)
		}
		if err != nil {
			return b, "", err
		}
		return b, t.UTC().Format(time.RFC3339Nano), nil
	case TimeTZ:
		var t timetz.TimeTZ
		if dir == Descending {
			b, t, err = DecodeTimeTZDescending(b)
		} else {
			b, t, err = DecodeTimeTZAscending(b)
		}
		if err != nil {
			return b, "", err
		}
		return b, t.String(), nil
	case Duration:
		var d duration.Duration
		if dir == Descending {
			b, d, err = DecodeDurationDescending(b)
		} else {
			b, d, err = DecodeDurationAscending(b)
		}
		if err != nil {
			return b, "", err
		}
		return b, d.StringNanos(), nil
	default:
		if len(b) >= 1 {
			switch b[0] {
			case jsonInvertedIndex:
				var str string
				str, b, err = prettyPrintInvertedIndexKey(b)
				if err != nil {
					return b, "", err
				}
				if str == "" {
					return prettyPrintFirstValue(dir, b)
				}
				return b, str, nil
			case jsonEmptyArray:
				return b[1:], "[]", nil
			case jsonEmptyObject:
				return b[1:], "{}", nil
			case emptyArray:
				return b[1:], "[]", nil
			}
		}
		// This shouldn't ever happen, but if it does, return an empty slice.
		return nil, strconv.Quote(string(b)), nil
	}
}

// UndoPrefixEnd is a partial inverse for roachpb.Key.PrefixEnd.
//
// In general, we can't undo PrefixEnd because it is lossy; we don't know how
// many FFs were stripped from the original key. For example:
//   - key:            01 02 03 FF FF
//   - PrefixEnd:      01 02 04
//   - UndoPrefixEnd:  01 02 03
//
// Some keys are not possible results of PrefixEnd; in particular, PrefixEnd
// keys never end in 00. If an impossible key is passed, the second return value
// is false.
//
// Specifically, calling UndoPrefixEnd will reverse the effects of calling a
// PrefixEnd on a byte sequence, except when the byte sequence represents a
// maximal prefix (i.e., 0xff...). This is because PrefixEnd is a lossy
// operation: PrefixEnd(0xff) returns 0xff rather than wrapping around to the
// minimal prefix 0x00. For consistency, UndoPrefixEnd is also lossy:
// UndoPrefixEnd(0x00) returns 0x00 rather than wrapping around to the maximal
// prefix 0xff.
//
// Formally:
//
//	PrefixEnd(UndoPrefixEnd(p)) = p for all non-minimal prefixes p
//	UndoPrefixEnd(PrefixEnd(p)) = p for all non-maximal prefixes p
//
// A minimal prefix is any prefix that consists only of one or more 0x00 bytes;
// analogously, a maximal prefix is any prefix that consists only of one or more
// 0xff bytes.
//
// UndoPrefixEnd is implemented here to avoid a circular dependency on roachpb,
// but arguably belongs in a byte-manipulation utility package.
func UndoPrefixEnd(b []byte) (_ []byte, ok bool) {
	if len(b) == 0 || b[len(b)-1] == 0 {
		// Not a possible result of PrefixEnd.
		return nil, false
	}
	out := append([]byte(nil), b...)
	out[len(out)-1]--
	return out, true
}

// MaxNonsortingVarintLen is the maximum length of an EncodeNonsortingVarint
// encoded value.
const MaxNonsortingVarintLen = binary.MaxVarintLen64

// EncodeNonsortingStdlibVarint encodes an int value using encoding/binary, appends it
// to the supplied buffer, and returns the final buffer.
func EncodeNonsortingStdlibVarint(appendTo []byte, x int64) []byte {
	// Fixed size array to allocate this on the stack.
	var scratch [binary.MaxVarintLen64]byte
	i := binary.PutVarint(scratch[:binary.MaxVarintLen64], x)
	return append(appendTo, scratch[:i]...)
}

// DecodeNonsortingStdlibVarint decodes a value encoded by EncodeNonsortingVarint. It
// returns the length of the encoded varint and value.
func DecodeNonsortingStdlibVarint(b []byte) (remaining []byte, length int, value int64, err error) {
	value, length = binary.Varint(b)
	if length <= 0 {
		return nil, 0, 0, fmt.Errorf("int64 varint decoding failed: %d", length)
	}
	return b[length:], length, value, nil
}

// MaxNonsortingUvarintLen is the maximum length of an EncodeNonsortingUvarint
// encoded value.
const MaxNonsortingUvarintLen = 10

// EncodeNonsortingUvarint encodes a uint64, appends it to the supplied buffer,
// and returns the final buffer. The encoding used is similar to
// encoding/binary, but with the most significant bits first:
//   - Unsigned integers are serialized 7 bits at a time, starting with the
//     most significant bits.
//   - The most significant bit (msb) in each output byte indicates if there
//     is a continuation byte (msb = 1).
func EncodeNonsortingUvarint(appendTo []byte, x uint64) []byte {
	switch {
	case x < (1 << 7):
		return append(appendTo, byte(x))
	case x < (1 << 14):
		return append(appendTo, 0x80|byte(x>>7), 0x7f&byte(x))
	case x < (1 << 21):
		return append(appendTo, 0x80|byte(x>>14), 0x80|byte(x>>7), 0x7f&byte(x))
	case x < (1 << 28):
		return append(appendTo, 0x80|byte(x>>21), 0x80|byte(x>>14), 0x80|byte(x>>7), 0x7f&byte(x))
	case x < (1 << 35):
		return append(appendTo, 0x80|byte(x>>28), 0x80|byte(x>>21), 0x80|byte(x>>14), 0x80|byte(x>>7), 0x7f&byte(x))
	case x < (1 << 42):
		return append(appendTo, 0x80|byte(x>>35), 0x80|byte(x>>28), 0x80|byte(x>>21), 0x80|byte(x>>14), 0x80|byte(x>>7), 0x7f&byte(x))
	case x < (1 << 49):
		return append(appendTo, 0x80|byte(x>>42), 0x80|byte(x>>35), 0x80|byte(x>>28), 0x80|byte(x>>21), 0x80|byte(x>>14), 0x80|byte(x>>7), 0x7f&byte(x))
	case x < (1 << 56):
		return append(appendTo, 0x80|byte(x>>49), 0x80|byte(x>>42), 0x80|byte(x>>35), 0x80|byte(x>>28), 0x80|byte(x>>21), 0x80|byte(x>>14), 0x80|byte(x>>7), 0x7f&byte(x))
	case x < (1 << 63):
		return append(appendTo, 0x80|byte(x>>56), 0x80|byte(x>>49), 0x80|byte(x>>42), 0x80|byte(x>>35), 0x80|byte(x>>28), 0x80|byte(x>>21), 0x80|byte(x>>14), 0x80|byte(x>>7), 0x7f&byte(x))
	default:
		return append(appendTo, 0x80|byte(x>>63), 0x80|byte(x>>56), 0x80|byte(x>>49), 0x80|byte(x>>42), 0x80|byte(x>>35), 0x80|byte(x>>28), 0x80|byte(x>>21), 0x80|byte(x>>14), 0x80|byte(x>>7), 0x7f&byte(x))
	}
}

// DecodeNonsortingUvarint decodes a value encoded by EncodeNonsortingUvarint. It
// returns the length of the encoded varint and value.
func DecodeNonsortingUvarint(buf []byte) (remaining []byte, length int, value uint64, err error) {
	// TODO(dan): Handle overflow.
	for i, b := range buf {
		value += uint64(b & 0x7f)
		if b < 0x80 {
			return buf[i+1:], i + 1, value, nil
		}
		value <<= 7
	}
	return buf, 0, 0, nil
}

// DecodeNonsortingStdlibUvarint decodes a value encoded with binary.PutUvarint. It
// returns the length of the encoded varint and value.
func DecodeNonsortingStdlibUvarint(
	buf []byte,
) (remaining []byte, length int, value uint64, err error) {
	i, n := binary.Uvarint(buf)
	if n <= 0 {
		return buf, 0, 0, errors.New("buffer too small")
	}
	return buf[n:], n, i, nil
}

// PeekLengthNonsortingUvarint returns the length of the value that starts at
// the beginning of buf and was encoded by EncodeNonsortingUvarint.
func PeekLengthNonsortingUvarint(buf []byte) int {
	for i, b := range buf {
		if b&0x80 == 0 {
			return i + 1
		}
	}
	return 0
}

// NoColumnID is a sentinel for the EncodeFooValue methods representing an
// invalid column id.
const NoColumnID uint32 = 0

// EncodeValueTag encodes the prefix that is used by each of the EncodeFooValue
// methods.
//
// The prefix uses varints to encode a column id and type, packing them into a
// single byte when they're small (colIDDelta < 8 and typ < 15). This works by
// shifting the colIDDelta "left" by 4 and putting any type less than 15 in the low
// bytes. The result is uvarint encoded and fits in one byte if the original
// column id fit in 3 bits. If it doesn't fit in one byte, the most significant
// bits spill to the "left", leaving the type bits always at the very "right".
//
// If the type is > 15, the reserved sentinel of 15 is placed in the type bits
// and a uvarint follows with the type value. This means that there are always
// one or two uvarints.
//
// Together, this means the everything but the last byte of the first uvarint
// can be dropped if the column id isn't needed.
func EncodeValueTag(appendTo []byte, colIDDelta uint32, typ Type) []byte {
	if typ >= SentinelType {
		appendTo = EncodeNonsortingUvarint(appendTo, uint64(colIDDelta)<<4|uint64(SentinelType))
		return EncodeNonsortingUvarint(appendTo, uint64(typ))
	}
	if colIDDelta == NoColumnID {
		// TODO(dan): EncodeValueTag is not inlined by the compiler. Copying this
		// special case into one of the EncodeFooValue functions speeds it up by
		// ~4ns.
		return append(appendTo, byte(typ))
	}
	return EncodeNonsortingUvarint(appendTo, uint64(colIDDelta)<<4|uint64(typ))
}

// EncodeNullValue encodes a null value, appends it to the supplied buffer, and
// returns the final buffer.
func EncodeNullValue(appendTo []byte, colIDDelta uint32) []byte {
	return EncodeValueTag(appendTo, colIDDelta, Null)
}

// EncodeNotNullValue encodes a not null value, appends it to the supplied
// buffer, and returns the final buffer.
func EncodeNotNullValue(appendTo []byte, colIDDelta uint32) []byte {
	return EncodeValueTag(appendTo, colIDDelta, NotNull)
}

// EncodeBoolValue encodes a bool value, appends it to the supplied buffer, and
// returns the final buffer.
func EncodeBoolValue(appendTo []byte, colIDDelta uint32, b bool) []byte {
	if b {
		return EncodeValueTag(appendTo, colIDDelta, True)
	}
	return EncodeValueTag(appendTo, colIDDelta, False)
}

// EncodeIntValue encodes an int value with its value tag, appends it to the
// supplied buffer, and returns the final buffer.
func EncodeIntValue(appendTo []byte, colIDDelta uint32, i int64) []byte {
	appendTo = EncodeValueTag(appendTo, colIDDelta, Int)
	return EncodeUntaggedIntValue(appendTo, i)
}

// EncodeUntaggedIntValue encodes an int value, appends it to the supplied buffer, and
// returns the final buffer.
func EncodeUntaggedIntValue(appendTo []byte, i int64) []byte {
	return EncodeNonsortingStdlibVarint(appendTo, i)
}

const floatValueEncodedLength = uint64AscendingEncodedLength

// EncodeFloatValue encodes a float value with its value tag, appends it to the
// supplied buffer, and returns the final buffer.
func EncodeFloatValue(appendTo []byte, colIDDelta uint32, f float64) []byte {
	appendTo = EncodeValueTag(appendTo, colIDDelta, Float)
	return EncodeUntaggedFloatValue(appendTo, f)
}

// EncodeUntaggedFloatValue encodes a float value, appends it to the supplied buffer,
// and returns the final buffer.
func EncodeUntaggedFloatValue(appendTo []byte, f float64) []byte {
	return EncodeUint64Ascending(appendTo, math.Float64bits(f))
}

// EncodeUntaggedFloat32Value encodes a float32 value, appends it to the supplied buffer,
// and returns the final buffer.
func EncodeUntaggedFloat32Value(appendTo []byte, f float32) []byte {
	return EncodeUint32Ascending(appendTo, math.Float32bits(f))
}

// EncodeBytesValue encodes a byte array value with its value tag, appends it to
// the supplied buffer, and returns the final buffer.
func EncodeBytesValue(appendTo []byte, colIDDelta uint32, data []byte) []byte {
	appendTo = EncodeValueTag(appendTo, colIDDelta, Bytes)
	return EncodeUntaggedBytesValue(appendTo, data)
}

// EncodeUntaggedBytesValue encodes a byte array value, appends it to the supplied
// buffer, and returns the final buffer.
func EncodeUntaggedBytesValue(appendTo []byte, data []byte) []byte {
	appendTo = EncodeNonsortingUvarint(appendTo, uint64(len(data)))
	return append(appendTo, data...)
}

// EncodeArrayValue encodes a byte array value with its value tag, appends it to
// the supplied buffer, and returns the final buffer.
func EncodeArrayValue(appendTo []byte, colIDDelta uint32, data []byte) []byte {
	appendTo = EncodeValueTag(appendTo, colIDDelta, Array)
	return EncodeUntaggedBytesValue(appendTo, data)
}

// EncodeTimeValue encodes a time.Time value with its value tag, appends it to
// the supplied buffer, and returns the final buffer.
func EncodeTimeValue(appendTo []byte, colIDDelta uint32, t time.Time) []byte {
	appendTo = EncodeValueTag(appendTo, colIDDelta, Time)
	return EncodeUntaggedTimeValue(appendTo, t)
}

// EncodeUntaggedTimeValue encodes a time.Time value, appends it to the supplied buffer,
// and returns the final buffer.
func EncodeUntaggedTimeValue(appendTo []byte, t time.Time) []byte {
	appendTo = EncodeNonsortingStdlibVarint(appendTo, t.Unix())
	return EncodeNonsortingStdlibVarint(appendTo, int64(t.Nanosecond()))
}

// EncodeTimeTZValue encodes a timetz.TimeTZ value with its value tag, appends it to
// the supplied buffer, and returns the final buffer.
func EncodeTimeTZValue(appendTo []byte, colIDDelta uint32, t timetz.TimeTZ) []byte {
	appendTo = EncodeValueTag(appendTo, colIDDelta, TimeTZ)
	return EncodeUntaggedTimeTZValue(appendTo, t)
}

// EncodeUntaggedTimeTZValue encodes a time.Time value, appends it to the supplied buffer,
// and returns the final buffer.
func EncodeUntaggedTimeTZValue(appendTo []byte, t timetz.TimeTZ) []byte {
	appendTo = EncodeNonsortingStdlibVarint(appendTo, int64(t.TimeOfDay))
	return EncodeNonsortingStdlibVarint(appendTo, int64(t.OffsetSecs))
}

// EncodeVoidValue encodes a void with its value tag, appends it to
// the supplied buffer and returns the final buffer.
func EncodeVoidValue(appendTo []byte, colIDDelta uint32) []byte {
	return EncodeValueTag(appendTo, colIDDelta, Void)
}

// EncodeBox2DValue encodes a geopb.BoundingBox with its value tag, appends it to
// the supplied buffer and returns the final buffer.
func EncodeBox2DValue(appendTo []byte, colIDDelta uint32, b geopb.BoundingBox) ([]byte, error) {
	appendTo = EncodeValueTag(appendTo, colIDDelta, Box2D)
	return EncodeUntaggedBox2DValue(appendTo, b)
}

// EncodeUntaggedBox2DValue encodes a geopb.BoundingBox value, appends it to the supplied buffer,
// and returns the final buffer.
func EncodeUntaggedBox2DValue(appendTo []byte, b geopb.BoundingBox) ([]byte, error) {
	appendTo = EncodeFloatAscending(appendTo, b.LoX)
	appendTo = EncodeFloatAscending(appendTo, b.HiX)
	appendTo = EncodeFloatAscending(appendTo, b.LoY)
	appendTo = EncodeFloatAscending(appendTo, b.HiY)
	return appendTo, nil
}

// EncodeGeoValue encodes a geopb.SpatialObject value with its value tag, appends it to
// the supplied buffer, and returns the final buffer.
func EncodeGeoValue(appendTo []byte, colIDDelta uint32, so *geopb.SpatialObject) ([]byte, error) {
	appendTo = EncodeValueTag(appendTo, colIDDelta, Geo)
	return EncodeUntaggedGeoValue(appendTo, so)
}

// EncodeUntaggedGeoValue encodes a geopb.SpatialObject value, appends it to the supplied buffer,
// and returns the final buffer.
func EncodeUntaggedGeoValue(appendTo []byte, so *geopb.SpatialObject) ([]byte, error) {
	bytes, err := protoutil.Marshal(so)
	if err != nil {
		return nil, err
	}
	return EncodeUntaggedBytesValue(appendTo, bytes), nil
}

// EncodeDecimalValue encodes an apd.Decimal value with its value tag, appends
// it to the supplied buffer, and returns the final buffer.
func EncodeDecimalValue(appendTo []byte, colIDDelta uint32, d *apd.Decimal) []byte {
	appendTo = EncodeValueTag(appendTo, colIDDelta, Decimal)
	return EncodeUntaggedDecimalValue(appendTo, d)
}

// EncodeUntaggedDecimalValue encodes an apd.Decimal value, appends it to the supplied
// buffer, and returns the final buffer.
func EncodeUntaggedDecimalValue(appendTo []byte, d *apd.Decimal) []byte {
	// To avoid the allocation, leave space for the varint, encode the decimal,
	// encode the varint, and shift the encoded decimal to the end of the
	// varint.
	varintPos := len(appendTo)
	// Manually append 10 (binary.MaxVarintLen64) 0s to avoid the allocation.
	appendTo = append(appendTo, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	decOffset := len(appendTo)
	appendTo = EncodeNonsortingDecimal(appendTo, d)
	decLen := len(appendTo) - decOffset
	varintLen := binary.PutUvarint(appendTo[varintPos:decOffset], uint64(decLen))
	copy(appendTo[varintPos+varintLen:varintPos+varintLen+decLen], appendTo[decOffset:decOffset+decLen])
	return appendTo[:varintPos+varintLen+decLen]
}

// EncodeDurationValue encodes a duration.Duration value with its value tag,
// appends it to the supplied buffer, and returns the final buffer.
func EncodeDurationValue(appendTo []byte, colIDDelta uint32, d duration.Duration) []byte {
	appendTo = EncodeValueTag(appendTo, colIDDelta, Duration)
	return EncodeUntaggedDurationValue(appendTo, d)
}

// EncodeUntaggedDurationValue encodes a duration.Duration value, appends it to the
// supplied buffer, and returns the final buffer.
func EncodeUntaggedDurationValue(appendTo []byte, d duration.Duration) []byte {
	appendTo = EncodeNonsortingStdlibVarint(appendTo, d.Months)
	appendTo = EncodeNonsortingStdlibVarint(appendTo, d.Days)
	return EncodeNonsortingStdlibVarint(appendTo, d.Nanos())
}

// EncodeBitArrayValue encodes a bit array value with its value tag,
// appends it to the supplied buffer, and returns the final buffer.
func EncodeBitArrayValue(appendTo []byte, colIDDelta uint32, d bitarray.BitArray) []byte {
	appendTo = EncodeValueTag(appendTo, colIDDelta, BitArray)
	return EncodeUntaggedBitArrayValue(appendTo, d)
}

// EncodeUntaggedBitArrayValue encodes a bit array value, appends it to the
// supplied buffer, and returns the final buffer.
func EncodeUntaggedBitArrayValue(appendTo []byte, d bitarray.BitArray) []byte {
	bitLen := d.BitLen()
	words, _ := d.EncodingParts()

	appendTo = EncodeNonsortingUvarint(appendTo, uint64(bitLen))
	for _, w := range words {
		appendTo = EncodeUint64Ascending(appendTo, w)
	}
	return appendTo
}

// EncodeUUIDValue encodes a uuid.UUID value with its value tag, appends it to
// the supplied buffer, and returns the final buffer.
func EncodeUUIDValue(appendTo []byte, colIDDelta uint32, u uuid.UUID) []byte {
	appendTo = EncodeValueTag(appendTo, colIDDelta, UUID)
	return EncodeUntaggedUUIDValue(appendTo, u)
}

// EncodeUntaggedUUIDValue encodes a uuid.UUID value, appends it to the supplied buffer,
// and returns the final buffer.
func EncodeUntaggedUUIDValue(appendTo []byte, u uuid.UUID) []byte {
	return append(appendTo, u.GetBytes()...)
}

// EncodeIPAddrValue encodes a ipaddr.IPAddr value with its value tag, appends
// it to the supplied buffer, and returns the final buffer.
func EncodeIPAddrValue(appendTo []byte, colIDDelta uint32, u ipaddr.IPAddr) []byte {
	appendTo = EncodeValueTag(appendTo, colIDDelta, IPAddr)
	return EncodeUntaggedIPAddrValue(appendTo, u)
}

// EncodeUntaggedIPAddrValue encodes a ipaddr.IPAddr value, appends it to the
// supplied buffer, and returns the final buffer.
func EncodeUntaggedIPAddrValue(appendTo []byte, u ipaddr.IPAddr) []byte {
	return u.ToBuffer(appendTo)
}

// EncodeJSONValue encodes an already-byte-encoded JSON value with no value tag
// but with a length prefix, appends it to the supplied buffer, and returns the
// final buffer.
func EncodeJSONValue(appendTo []byte, colIDDelta uint32, data []byte) []byte {
	appendTo = EncodeValueTag(appendTo, colIDDelta, JSON)
	return EncodeUntaggedBytesValue(appendTo, data)
}

// EncodeTSQueryValue encodes an already-byte-encoded TSQuery value with no
// value tag but with a length prefix, appends it to the supplied buffer, and
// returns the final buffer.
func EncodeTSQueryValue(appendTo []byte, colIDDelta uint32, data []byte) []byte {
	appendTo = EncodeValueTag(appendTo, colIDDelta, TSQuery)
	return EncodeUntaggedBytesValue(appendTo, data)
}

// EncodeTSVectorValue encodes an already-byte-encoded TSVector value with no
// value tag but with a length prefix, appends it to the supplied buffer, and
// returns the final buffer.
func EncodeTSVectorValue(appendTo []byte, colIDDelta uint32, data []byte) []byte {
	appendTo = EncodeValueTag(appendTo, colIDDelta, TSVector)
	return EncodeUntaggedBytesValue(appendTo, data)
}

// EncodePGVectorValue encodes an already-byte-encoded PGVector value with no
// value tag but with a length prefix, appends it to the supplied buffer, and
// returns the final buffer.
func EncodePGVectorValue(appendTo []byte, colIDDelta uint32, data []byte) []byte {
	appendTo = EncodeValueTag(appendTo, colIDDelta, PGVector)
	return EncodeUntaggedBytesValue(appendTo, data)
}

// DecodeValueTag decodes a value encoded by EncodeValueTag, used as a prefix in
// each of the other EncodeFooValue methods.
//
// The tag is structured such that the encoded column id can be dropped from the
// front by removing the first `typeOffset` bytes. DecodeValueTag,
// PeekValueLength and each of the DecodeFooValue methods will still work as
// expected with `b[typeOffset:]`. (Except, obviously, the column id is no
// longer encoded so if this suffix is passed back to DecodeValueTag, the
// returned colIDDelta should be discarded.)
//
// Concretely:
//
//	b := ...
//	typeOffset, _, colIDDelta, typ, err := DecodeValueTag(b)
//	_, _, _, typ, err := DecodeValueTag(b[typeOffset:])
//
// will return the same typ and err and
//
//	DecodeFooValue(b)
//	DecodeFooValue(b[typeOffset:])
//
// will return the same thing. PeekValueLength works as expected with either of
// `b` or `b[typeOffset:]`.
func DecodeValueTag(
	b []byte,
) (typeOffset int, dataOffset int, colIDDelta uint32, typ Type, err error) {
	// TODO(dan): This can be made faster by special casing the single byte
	// version and skipping the column id extraction when it's not needed.
	if len(b) == 0 {
		return 0, 0, 0, Unknown, fmt.Errorf("empty array")
	}
	var n int
	var tag uint64
	b, n, tag, err = DecodeNonsortingUvarint(b)
	if err != nil {
		return 0, 0, 0, Unknown, err
	}
	colIDDelta = uint32(tag >> 4)

	typ = Type(tag & 0xf)
	typeOffset = n - 1
	dataOffset = n
	if typ == SentinelType {
		_, n, tag, err = DecodeNonsortingUvarint(b)
		if err != nil {
			return 0, 0, 0, Unknown, err
		}
		typ = Type(tag)
		dataOffset += n
	}
	return typeOffset, dataOffset, colIDDelta, typ, nil
}

// DecodeBoolValue decodes a value encoded by EncodeBoolValue.
func DecodeBoolValue(buf []byte) (remaining []byte, b bool, err error) {
	_, dataOffset, _, typ, err := DecodeValueTag(buf)
	if err != nil {
		return buf, false, err
	}
	buf = buf[dataOffset:]
	switch typ {
	case True:
		return buf, true, nil
	case False:
		return buf, false, nil
	default:
		return buf, false, fmt.Errorf("value type is not %s or %s: %s", True, False, typ)
	}
}

// DecodeIntValue decodes a value encoded by EncodeIntValue.
func DecodeIntValue(b []byte) (remaining []byte, i int64, err error) {
	b, err = decodeValueTypeAssert(b, Int)
	if err != nil {
		return b, 0, err
	}
	return DecodeUntaggedIntValue(b)
}

// DecodeUntaggedIntValue decodes a value encoded by EncodeUntaggedIntValue.
func DecodeUntaggedIntValue(b []byte) (remaining []byte, i int64, err error) {
	b, _, i, err = DecodeNonsortingStdlibVarint(b)
	return b, i, err
}

// DecodeFloatValue decodes a value encoded by EncodeFloatValue.
func DecodeFloatValue(b []byte) (remaining []byte, f float64, err error) {
	b, err = decodeValueTypeAssert(b, Float)
	if err != nil {
		return b, 0, err
	}
	return DecodeUntaggedFloatValue(b)
}

// DecodeUntaggedFloatValue decodes a value encoded by EncodeUntaggedFloatValue.
func DecodeUntaggedFloatValue(b []byte) (remaining []byte, f float64, err error) {
	if len(b) < 8 {
		return b, 0, fmt.Errorf("float64 value should be exactly 8 bytes: %d", len(b))
	}
	var i uint64
	b, i, err = DecodeUint64Ascending(b)
	return b, math.Float64frombits(i), err
}

// DecodeUntaggedFloat32Value decodes a value encoded by EncodeUntaggedFloat32Value.
func DecodeUntaggedFloat32Value(b []byte) (remaining []byte, f float32, err error) {
	if len(b) < 4 {
		return b, 0, fmt.Errorf("float32 value should be exactly 4 bytes: %d", len(b))
	}
	var i uint32
	b, i, err = DecodeUint32Ascending(b)
	return b, math.Float32frombits(i), err
}

// DecodeBytesValue decodes a value encoded by EncodeBytesValue.
func DecodeBytesValue(b []byte) (remaining []byte, data []byte, err error) {
	b, err = decodeValueTypeAssert(b, Bytes)
	if err != nil {
		return b, nil, err
	}
	return DecodeUntaggedBytesValue(b)
}

// DecodeUntaggedBytesValue decodes a value encoded by EncodeUntaggedBytesValue.
func DecodeUntaggedBytesValue(b []byte) (remaining, data []byte, err error) {
	var i uint64
	b, _, i, err = DecodeNonsortingUvarint(b)
	if err != nil {
		return b, nil, err
	}
	return b[int(i):], b[:int(i)], nil
}

// DecodeTimeValue decodes a value encoded by EncodeTimeValue.
func DecodeTimeValue(b []byte) (remaining []byte, t time.Time, err error) {
	b, err = decodeValueTypeAssert(b, Time)
	if err != nil {
		return b, time.Time{}, err
	}
	return DecodeUntaggedTimeValue(b)
}

// DecodeUntaggedTimeValue decodes a value encoded by EncodeUntaggedTimeValue.
func DecodeUntaggedTimeValue(b []byte) (remaining []byte, t time.Time, err error) {
	var sec, nsec int64
	b, _, sec, err = DecodeNonsortingStdlibVarint(b)
	if err != nil {
		return b, time.Time{}, err
	}
	b, _, nsec, err = DecodeNonsortingStdlibVarint(b)
	if err != nil {
		return b, time.Time{}, err
	}
	return b, timeutil.Unix(sec, nsec), nil
}

// DecodeTimeTZValue decodes a value encoded by EncodeTimeTZValue.
func DecodeTimeTZValue(b []byte) (remaining []byte, t timetz.TimeTZ, err error) {
	b, err = decodeValueTypeAssert(b, TimeTZ)
	if err != nil {
		return b, timetz.TimeTZ{}, err
	}
	return DecodeUntaggedTimeTZValue(b)
}

// DecodeUntaggedTimeTZValue decodes a value encoded by EncodeUntaggedTimeTZValue.
func DecodeUntaggedTimeTZValue(b []byte) (remaining []byte, t timetz.TimeTZ, err error) {
	var timeOfDayMicros int64
	b, _, timeOfDayMicros, err = DecodeNonsortingStdlibVarint(b)
	if err != nil {
		return b, timetz.TimeTZ{}, err
	}
	var offsetSecs int64
	b, _, offsetSecs, err = DecodeNonsortingStdlibVarint(b)
	if err != nil {
		return b, timetz.TimeTZ{}, err
	}
	// Do not use timeofday.FromInt as it truncates 24:00 into 00:00.
	return b, timetz.MakeTimeTZ(timeofday.TimeOfDay(timeOfDayMicros), int32(offsetSecs)), nil
}

// DecodeDecimalValue decodes a value encoded by EncodeDecimalValue.
func DecodeDecimalValue(b []byte) (remaining []byte, d apd.Decimal, err error) {
	b, err = decodeValueTypeAssert(b, Decimal)
	if err != nil {
		return b, apd.Decimal{}, err
	}
	return DecodeUntaggedDecimalValue(b)
}

// DecodeUntaggedBox2DValue decodes a value encoded by EncodeUntaggedBox2DValue.
func DecodeUntaggedBox2DValue(b []byte) (remaining []byte, box geopb.BoundingBox, err error) {
	box = geopb.BoundingBox{}
	remaining = b

	remaining, box.LoX, err = DecodeFloatAscending(remaining)
	if err != nil {
		return b, box, err
	}
	remaining, box.HiX, err = DecodeFloatAscending(remaining)
	if err != nil {
		return b, box, err
	}
	remaining, box.LoY, err = DecodeFloatAscending(remaining)
	if err != nil {
		return b, box, err
	}
	remaining, box.HiY, err = DecodeFloatAscending(remaining)
	if err != nil {
		return b, box, err
	}
	return remaining, box, err
}

// DecodeUntaggedGeoValue decodes a value encoded by EncodeUntaggedGeoValue into
// the provided geopb.SpatialObject reference. The so parameter must already be
// empty/reset.
func DecodeUntaggedGeoValue(b []byte, so *geopb.SpatialObject) (remaining []byte, err error) {
	var data []byte
	remaining, data, err = DecodeUntaggedBytesValue(b)
	if err != nil {
		return b, err
	}
	// Not using protoutil.Unmarshal since the call to so.Reset() will waste the
	// pre-allocated EWKB.
	err = so.Unmarshal(data)
	return remaining, err
}

// DecodeUntaggedDecimalValue decodes a value encoded by EncodeUntaggedDecimalValue.
func DecodeUntaggedDecimalValue(b []byte) (remaining []byte, d apd.Decimal, err error) {
	var i uint64
	b, _, i, err = DecodeNonsortingStdlibUvarint(b)
	if err != nil {
		return b, apd.Decimal{}, err
	}
	d, err = DecodeNonsortingDecimal(b[:int(i)], nil)
	return b[int(i):], d, err
}

// DecodeIntoUntaggedDecimalValue is like DecodeUntaggedDecimalValue except it
// writes the new Decimal into the input apd.Decimal pointer, which must be
// non-nil.
func DecodeIntoUntaggedDecimalValue(d *apd.Decimal, b []byte) (remaining []byte, err error) {
	var i uint64
	b, _, i, err = DecodeNonsortingStdlibUvarint(b)
	if err != nil {
		return b, err
	}
	err = DecodeIntoNonsortingDecimal(d, b[:int(i)], nil)
	return b[int(i):], err
}

// DecodeDurationValue decodes a value encoded by EncodeUntaggedDurationValue.
func DecodeDurationValue(b []byte) (remaining []byte, d duration.Duration, err error) {
	b, err = decodeValueTypeAssert(b, Duration)
	if err != nil {
		return b, duration.Duration{}, err
	}
	return DecodeUntaggedDurationValue(b)
}

// DecodeUntaggedDurationValue decodes a value encoded by EncodeUntaggedDurationValue.
func DecodeUntaggedDurationValue(b []byte) (remaining []byte, d duration.Duration, err error) {
	var months, days, nanos int64
	b, _, months, err = DecodeNonsortingStdlibVarint(b)
	if err != nil {
		return b, duration.Duration{}, err
	}
	b, _, days, err = DecodeNonsortingStdlibVarint(b)
	if err != nil {
		return b, duration.Duration{}, err
	}
	b, _, nanos, err = DecodeNonsortingStdlibVarint(b)
	if err != nil {
		return b, duration.Duration{}, err
	}
	return b, duration.DecodeDuration(months, days, nanos), nil
}

// DecodeBitArrayValue decodes a value encoded by EncodeUntaggedBitArrayValue.
func DecodeBitArrayValue(b []byte) (remaining []byte, d bitarray.BitArray, err error) {
	b, err = decodeValueTypeAssert(b, BitArray)
	if err != nil {
		return b, bitarray.BitArray{}, err
	}
	return DecodeUntaggedBitArrayValue(b)
}

// DecodeUntaggedBitArrayValue decodes a value encoded by EncodeUntaggedBitArrayValue.
func DecodeUntaggedBitArrayValue(b []byte) (remaining []byte, d bitarray.BitArray, err error) {
	var bitLen uint64
	b, _, bitLen, err = DecodeNonsortingUvarint(b)
	if err != nil {
		return b, bitarray.BitArray{}, err
	}
	words, lastBitsUsed := bitarray.EncodingPartsForBitLen(uint(bitLen))
	for i := range words {
		var val uint64
		b, val, err = DecodeUint64Ascending(b)
		if err != nil {
			return b, bitarray.BitArray{}, err
		}
		words[i] = val
	}
	ba, err := bitarray.FromEncodingParts(words, lastBitsUsed)
	return b, ba, err
}

const uuidValueEncodedLength = 16

var _ [uuidValueEncodedLength]byte = uuid.UUID{} // Assert that uuid.UUID is length 16.

// DecodeUUIDValue decodes a value encoded by EncodeUUIDValue.
func DecodeUUIDValue(b []byte) (remaining []byte, u uuid.UUID, err error) {
	b, err = decodeValueTypeAssert(b, UUID)
	if err != nil {
		return b, u, err
	}
	return DecodeUntaggedUUIDValue(b)
}

// DecodeUntaggedUUIDValue decodes a value encoded by EncodeUntaggedUUIDValue.
func DecodeUntaggedUUIDValue(b []byte) (remaining []byte, u uuid.UUID, err error) {
	if len(b) < uuidValueEncodedLength {
		return b, uuid.UUID{}, errors.Errorf("invalid uuid length of %d", len(b))
	}
	u, err = uuid.FromBytes(b[:uuidValueEncodedLength])
	if err != nil {
		return b, uuid.UUID{}, err
	}
	return b[uuidValueEncodedLength:], u, nil
}

// DecodeIPAddrValue decodes a value encoded by EncodeIPAddrValue.
func DecodeIPAddrValue(b []byte) (remaining []byte, u ipaddr.IPAddr, err error) {
	b, err = decodeValueTypeAssert(b, IPAddr)
	if err != nil {
		return b, u, err
	}
	return DecodeUntaggedIPAddrValue(b)
}

// DecodeUntaggedIPAddrValue decodes a value encoded by EncodeUntaggedIPAddrValue.
func DecodeUntaggedIPAddrValue(b []byte) (remaining []byte, u ipaddr.IPAddr, err error) {
	remaining, err = u.FromBuffer(b)
	return remaining, u, err
}

func decodeValueTypeAssert(b []byte, expected Type) ([]byte, error) {
	_, dataOffset, _, typ, err := DecodeValueTag(b)
	if err != nil {
		return b, err
	}
	b = b[dataOffset:]
	if typ != expected {
		return b, errors.Errorf("value type is not %s: %s", expected, typ)
	}
	return b, nil
}

// PeekValueLength returns the length of the encoded value at the start of b.
// Note: If this function succeeds, it's not a guarantee that decoding the value
// will succeed.
//
// `b` can point either at beginning of the "full tag" with the column id, or it
// can point to the beginning of the type part of the tag, as indicated by the
// `typeOffset` returned by this or DecodeValueTag.
//
// The length returned is the full length of the encoded value, including the
// entire tag.
func PeekValueLength(b []byte) (typeOffset int, length int, err error) {
	if len(b) == 0 {
		return 0, 0, nil
	}
	var dataOffset int
	var typ Type
	typeOffset, dataOffset, _, typ, err = DecodeValueTag(b)
	if err != nil {
		return 0, 0, err
	}
	length, err = PeekValueLengthWithOffsetsAndType(b, dataOffset, typ)
	return typeOffset, length, err
}

// PeekValueLengthWithOffsetsAndType is the same as PeekValueLength, except it
// expects a dataOffset and typ value from a previous call to DecodeValueTag
// on its input byte slice. Use this if you've already called DecodeValueTag
// on the input for another reason, to avoid it getting called twice.
func PeekValueLengthWithOffsetsAndType(b []byte, dataOffset int, typ Type) (length int, err error) {
	b = b[dataOffset:]
	switch typ {
	case Null:
		return dataOffset, nil
	case Void:
		return dataOffset, nil
	case True, False:
		return dataOffset, nil
	case Int:
		_, n, _, err := DecodeNonsortingStdlibVarint(b)
		return dataOffset + n, err
	case Float:
		return dataOffset + floatValueEncodedLength, nil
	case Bytes, Array, JSON, Geo, TSVector, TSQuery, PGVector:
		_, n, i, err := DecodeNonsortingUvarint(b)
		return dataOffset + n + int(i), err
	case Box2D:
		length, err := peekBox2DLength(b)
		if err != nil {
			return 0, err
		}
		return dataOffset + length, nil
	case BitArray:
		_, n, bitLen, err := DecodeNonsortingUvarint(b)
		if err != nil {
			return 0, err
		}
		numWords, _ := bitarray.SizesForBitLen(uint(bitLen))
		return dataOffset + n + int(numWords)*8, err
	case Tuple:
		rem, l, numTuples, err := DecodeNonsortingUvarint(b)
		if err != nil {
			return 0, errors.Wrapf(err, "cannot decode tuple header: ")
		}
		for i := 0; i < int(numTuples); i++ {
			_, entryLen, err := PeekValueLength(rem)
			if err != nil {
				return 0, errors.Wrapf(err, "cannot peek tuple entry %d", i)
			}
			l += entryLen
			rem = rem[entryLen:]
		}
		return dataOffset + l, nil
	case Decimal:
		_, n, i, err := DecodeNonsortingStdlibUvarint(b)
		return dataOffset + n + int(i), err
	case Time, TimeTZ:
		n, err := getMultiNonsortingVarintLen(b, 2)
		return dataOffset + n, err
	case Duration:
		n, err := getMultiNonsortingVarintLen(b, 3)
		return dataOffset + n, err
	case UUID:
		return dataOffset + uuidValueEncodedLength, err
	case IPAddr:
		family := ipaddr.IPFamily(b[0])
		if family == ipaddr.IPv4family {
			return dataOffset + ipaddr.IPv4size, err
		} else if family == ipaddr.IPv6family {
			return dataOffset + ipaddr.IPv6size, err
		}
		return 0, errors.Errorf("got invalid INET IP family: %d", family)
	default:
		return 0, errors.Errorf("unknown type %s", typ)
	}
}

// PrintableBytes returns true iff the given byte array is a valid
// UTF-8 sequence and it is printable.
func PrintableBytes(b []byte) bool {
	return len(bytes.TrimLeftFunc(b, isValidAndPrintableRune)) == 0
}

func isValidAndPrintableRune(r rune) bool {
	return r != utf8.RuneError && unicode.IsPrint(r)
}

// PrettyPrintJSONValueEncoded returns a string representation of the encoded
// JSON object. It is injected from util/json to avoid an import cycle.
var PrettyPrintJSONValueEncoded func([]byte) (string, error)

var prettyPrintJSONValueEncodedNilErr = errors.New("PrettyPrintJSONValueEncoded is not injected")

// PrettyPrintArrayValueEncoded returns a string representation of the encoded
// array object if possible. It is injected from rowenc/valueside to avoid an
// import cycle.
var PrettyPrintArrayValueEncoded func([]byte) (string, error)

var prettyPrintArrayValueEncodedNilErr = errors.New("PrettyPrintArrayValueEncoded is not injected")

// PrettyPrintTupleValueEncoded returns a string representation of the encoded
// tuple object if possible. It is injected from rowenc/valueside to avoid an
// import cycle.
var PrettyPrintTupleValueEncoded func([]byte) ([]byte, string, error)

var prettyPrintTupleValueEncodedNilErr = errors.New("PrettyPrintTupleValueEncoded is not injected")

// PrettyPrintValueEncoded returns a string representation of the first
// decodable value in the provided byte slice, along with the remaining byte
// slice after decoding.
func PrettyPrintValueEncoded(b []byte) ([]byte, string, error) {
	_, dataOffset, _, typ, err := DecodeValueTag(b)
	if err != nil {
		return b, "", err
	}
	switch typ {
	case Null:
		b = b[dataOffset:]
		return b, "NULL", nil
	case True:
		b = b[dataOffset:]
		return b, "true", nil
	case False:
		b = b[dataOffset:]
		return b, "false", nil
	case Int:
		var i int64
		b, i, err = DecodeIntValue(b)
		if err != nil {
			return b, "", err
		}
		return b, strconv.FormatInt(i, 10), nil
	case Float:
		var f float64
		b, f, err = DecodeFloatValue(b)
		if err != nil {
			return b, "", err
		}
		return b, strconv.FormatFloat(f, 'g', -1, 64), nil
	case Decimal:
		var d apd.Decimal
		b, d, err = DecodeDecimalValue(b)
		if err != nil {
			return b, "", err
		}
		return b, d.String(), nil
	case Bytes:
		var data []byte
		b, data, err = DecodeBytesValue(b)
		if err != nil {
			return b, "", err
		}
		if PrintableBytes(data) {
			return b, string(data), nil
		}
		// The following code extends hex.EncodeToString().
		dst := make([]byte, 2+hex.EncodedLen(len(data)))
		dst[0], dst[1] = '0', 'x'
		hex.Encode(dst[2:], data)
		return b, string(dst), nil
	case Time:
		var t time.Time
		b, t, err = DecodeTimeValue(b)
		if err != nil {
			return b, "", err
		}
		return b, t.UTC().Format(time.RFC3339Nano), nil
	case TimeTZ:
		var t timetz.TimeTZ
		b, t, err = DecodeTimeTZValue(b)
		if err != nil {
			return b, "", err
		}
		return b, t.String(), nil
	case Duration:
		var d duration.Duration
		b, d, err = DecodeDurationValue(b)
		if err != nil {
			return b, "", err
		}
		return b, d.StringNanos(), nil
	case BitArray:
		var d bitarray.BitArray
		b, d, err = DecodeBitArrayValue(b)
		if err != nil {
			return b, "", err
		}
		return b, "B" + d.String(), nil
	case UUID:
		var u uuid.UUID
		b, u, err = DecodeUUIDValue(b)
		if err != nil {
			return b, "", err
		}
		return b, u.String(), nil
	case IPAddr:
		var ipAddr ipaddr.IPAddr
		b, ipAddr, err = DecodeIPAddrValue(b)
		if err != nil {
			return b, "", err
		}
		return b, ipAddr.String(), nil
	case JSON:
		b = b[dataOffset:]
		var data []byte
		b, data, err = DecodeUntaggedBytesValue(b)
		if err != nil {
			return b, "", err
		}
		if PrettyPrintJSONValueEncoded == nil {
			return b, "", prettyPrintJSONValueEncodedNilErr
		}
		var s string
		s, err = PrettyPrintJSONValueEncoded(data)
		return b, s, err
	case Array:
		b = b[dataOffset:]
		var data []byte
		b, data, err = DecodeUntaggedBytesValue(b)
		if err != nil {
			return b, "", err
		}
		if PrettyPrintArrayValueEncoded == nil {
			return b, "", prettyPrintArrayValueEncodedNilErr
		}
		var s string
		s, err = PrettyPrintArrayValueEncoded(data)
		return b, s, err
	case Tuple:
		b = b[dataOffset:]
		if PrettyPrintTupleValueEncoded == nil {
			return b, "", prettyPrintTupleValueEncodedNilErr
		}
		var s string
		b, s, err = PrettyPrintTupleValueEncoded(b)
		return b, s, err
	default:
		return b, "", errors.Errorf("unknown type %s", typ)
	}
}

// getInvertedIndexKeyLength finds the length of an inverted index key
// encoded as a byte array.
func getInvertedIndexKeyLength(b []byte) (int, error) {
	skipped := 0
	for {
		i := bytes.IndexByte(b[skipped:], escape)
		if i == -1 {
			return 0, errors.Errorf("malformed inverted index key in buffer %#x", b)
		}
		skipped += i + escapeLength
		switch b[skipped-1] {
		case escapedTerm, jsonEmptyObject, jsonEmptyArray:
			return skipped, nil
		}
	}
}

// getJSONInvertedIndexKeyLength returns the length of encoded JSON inverted index
// key at the start of b.
func getJSONInvertedIndexKeyLength(buf []byte) (int, error) {
	len, err := getInvertedIndexKeyLength(buf)
	if err != nil {
		return 0, err
	}

	switch buf[len] {
	case jsonEmptyArray, jsonEmptyObject:
		return len + 1, nil

	default:
		valLen, err := PeekLength(buf[len:])
		if err != nil {
			return 0, err
		}

		return len + valLen, nil
	}
}

func getGeoInvertedIndexKeyLength(buf []byte) (int, error) {
	// Minimum: 1 byte marker + 1 byte cell length +
	//          1 byte bbox encoding kind + 16 bytes for 2 floats
	if len(buf) < 3+2*uint64AscendingEncodedLength {
		return 0, errors.Errorf("buf length %d too small", len(buf))
	}
	var cellLen int
	var err error
	if cellLen, err = getVarintLen(buf[1:]); err != nil {
		return 0, err
	}
	encodingKind := geoInvertedBBoxEncodingKind(buf[1+cellLen])
	floatsLen := 4 * uint64AscendingEncodedLength
	if encodingKind == geoInvertedTwoFloats {
		floatsLen = 2 * uint64AscendingEncodedLength
	}
	return 1 + cellLen + 1 + floatsLen, nil
}

// EncodeJSONNullKeyMarker adds a JSON Null key encoding marker
// to buf and returns the new buffer.
func EncodeJSONNullKeyMarker(buf []byte, dir Direction) []byte {
	switch dir {
	case Ascending:
		return append(buf, jsonNullKeyMarker)
	case Descending:
		return append(buf, jsonNullKeyDescendingMarker)
	default:
		panic("invalid direction")
	}
}

// EncodeJSONStringKeyMarker adds a JSON String key encoding marker
// to buf and returns the new buffer.
func EncodeJSONStringKeyMarker(buf []byte, dir Direction) []byte {
	switch dir {
	case Ascending:
		return append(buf, jsonStringKeyMarker)
	case Descending:
		return append(buf, jsonStringKeyDescendingMarker)
	default:
		panic("invalid direction")
	}
}

// EncodeJSONNumberKeyMarker adds a JSON Number key encoding marker
// to buf and returns the new buffer.
func EncodeJSONNumberKeyMarker(buf []byte, dir Direction) []byte {
	switch dir {
	case Ascending:
		return append(buf, jsonNumberKeyMarker)
	case Descending:
		return append(buf, jsonNumberKeyDescendingMarker)
	default:
		panic("invalid direction")
	}
}

// EncodeJSONFalseKeyMarker adds a JSON False key encoding marker
// to buf and returns the new buffer.
func EncodeJSONFalseKeyMarker(buf []byte, dir Direction) []byte {
	switch dir {
	case Ascending:
		return append(buf, jsonFalseKeyMarker)
	case Descending:
		return append(buf, jsonFalseKeyDescendingMarker)
	default:
		panic("invalid direction")
	}
}

// EncodeJSONTrueKeyMarker adds a JSON True key encoding marker
// to buf and returns the new buffer.
func EncodeJSONTrueKeyMarker(buf []byte, dir Direction) []byte {
	switch dir {
	case Ascending:
		return append(buf, jsonTrueKeyMarker)
	case Descending:
		return append(buf, jsonTrueKeyDescendingMarker)
	default:
		panic("invalid direction")
	}
}

// EncodeJSONArrayKeyMarker adds a JSON Array key encoding marker
// to buf and returns the new buffer.
func EncodeJSONArrayKeyMarker(buf []byte, dir Direction, arrayLength int64) []byte {
	switch dir {
	case Ascending:
		if arrayLength == 0 {
			return append(buf, jsonEmptyArrayKeyMarker)
		}
		return append(buf, jsonArrayKeyMarker)
	case Descending:
		if arrayLength == 0 {
			return append(buf, jsonEmptyArrayKeyDescendingMarker)
		}
		return append(buf, jsonArrayKeyDescendingMarker)
	default:
		panic("invalid direction")
	}
}

// EncodeJSONKeyTerminator adds a JSON Key terminator
// to buf and returns the buffer.
func EncodeJSONKeyTerminator(buf []byte, dir Direction) []byte {
	switch dir {
	case Ascending:
		return append(buf, jsonKeyTerminator)
	case Descending:
		return append(buf, jsonKeyDescendingTerminator)
	default:
		panic("invalid direction")
	}
}

// EncodeJSONObjectKeyMarker adds a JSON Object key encoding marker
// to buf and returns the new buffer.
func EncodeJSONObjectKeyMarker(buf []byte, dir Direction) []byte {
	switch dir {
	case Ascending:
		return append(buf, jsonObjectKeyMarker)
	case Descending:
		return append(buf, jsonObjectKeyDescendingMarker)
	default:
		panic("invalid direction")
	}
}

func EncodeJSONValueLength(buf []byte, dir Direction, v int64) []byte {
	switch dir {
	case Ascending:
		return EncodeVarintAscending(buf, v)
	case Descending:
		return EncodeVarintDescending(buf, v)
	default:
		panic("invalid direction")
	}
}

func DecodeJSONValueLength(buf []byte, dir Direction) ([]byte, int64, error) {
	var v int64
	var err error
	switch dir {
	case Ascending:
		buf, v, err = DecodeVarintAscending(buf)
		return buf, v, err
	case Descending:
		buf, v, err = DecodeVarintDescending(buf)
		return buf, v, err
	default:
		panic("invalid direction")
	}
}

// EncodeArrayKeyMarker adds the array key encoding marker to buf and
// returns the new buffer.
func EncodeArrayKeyMarker(buf []byte, dir Direction) []byte {
	switch dir {
	case Ascending:
		return append(buf, arrayKeyMarker)
	case Descending:
		return append(buf, arrayKeyDescendingMarker)
	default:
		panic("invalid direction")
	}
}

// EncodeArrayKeyTerminator adds the array key terminator to buf and
// returns the new buffer.
func EncodeArrayKeyTerminator(buf []byte, dir Direction) []byte {
	switch dir {
	case Ascending:
		return append(buf, arrayKeyTerminator)
	case Descending:
		return append(buf, arrayKeyDescendingTerminator)
	default:
		panic("invalid direction")
	}
}

// EncodeNullWithinArrayKey encodes NULL within a key encoded array.
func EncodeNullWithinArrayKey(buf []byte, dir Direction) []byte {
	switch dir {
	case Ascending:
		return append(buf, ascendingNullWithinArrayKey)
	case Descending:
		return append(buf, descendingNullWithinArrayKey)
	default:
		panic("invalid direction")
	}
}

// IsNextByteArrayEncodedNull returns if the first byte in the input
// is the NULL encoded byte within an array key.
func IsNextByteArrayEncodedNull(buf []byte, dir Direction) bool {
	expected := ascendingNullWithinArrayKey
	if dir == Descending {
		expected = descendingNullWithinArrayKey
	}
	return buf[0] == expected
}

// ValidateAndConsumeJSONKeyMarker checks that the marker at the front
// of buf is valid/invalid for a given JSON value for the given direction.
// If the JSON marker is valid, the marker is consumed and the remaining
// bytes in the array are returned.
func ValidateAndConsumeJSONKeyMarker(buf []byte, dir Direction) ([]byte, Type, error) {
	typ := PeekType(buf)
	switch dir {
	case Descending:
		switch typ {
		case JSONNullDesc, JSONNumberDesc, JSONStringDesc, JSONFalseDesc,
			JSONTrueDesc, JSONArrayDesc, JSONObjectDesc, JsonEmptyArrayDesc:
			return buf[1:], typ, nil
		default:
			return nil, Unknown, errors.Newf("invalid type found %s", typ)
		}
	case Ascending:
		switch typ {
		case JSONNull, JSONNumber, JSONString, JSONFalse, JSONTrue, JSONArray,
			JSONObject, JsonEmptyArray:
			return buf[1:], typ, nil
		default:
			return nil, Unknown, errors.Newf("invalid type found %s", typ)
		}
	default:
		return nil, Unknown, errors.Newf("invalid direction %s", dir)
	}
}

// ValidateAndConsumeArrayKeyMarker checks that the marker at the front
// of buf is valid for an array of the given direction, and consumes it
// if so. It returns an error if the tag is invalid.
func ValidateAndConsumeArrayKeyMarker(buf []byte, dir Direction) ([]byte, error) {
	typ := PeekType(buf)
	expected := ArrayKeyAsc
	if dir == Descending {
		expected = ArrayKeyDesc
	}
	if typ != expected {
		return nil, errors.Newf("invalid type found %s", typ)
	}
	return buf[1:], nil
}

// IsArrayKeyDone returns if the first byte in the input is the array
// terminator for the input direction.
func IsArrayKeyDone(buf []byte, dir Direction) bool {
	expected := arrayKeyTerminator
	if dir == Descending {
		expected = arrayKeyDescendingTerminator
	}
	return buf[0] == expected
}

// isJSONKeyDone returns if the first byte in the input is the JSON
// terminator for the input direction.
func IsJSONKeyDone(buf []byte, dir Direction) bool {
	expected := jsonKeyTerminator
	if dir == Descending {
		expected = jsonKeyDescendingTerminator
	}
	return buf[0] == expected
}

// BytesNext returns the next possible byte slice, using the extra capacity
// of the provided slice if possible, and if not, appending an \x00.
func BytesNext(b []byte) []byte {
	if cap(b) > len(b) {
		bNext := b[:len(b)+1]
		if bNext[len(bNext)-1] == 0 {
			return bNext
		}
	}
	// TODO(spencer): Do we need to enforce KeyMaxLength here?
	// Switched to "make and copy" pattern in #4963 for performance.
	bn := make([]byte, len(b)+1)
	copy(bn, b)
	bn[len(bn)-1] = 0
	return bn
}

// BytesPrevish returns a previous byte slice in lexicographical ordering. It is
// impossible in general to find the exact previous byte slice, because it has
// an infinite number of 0xff bytes at the end, so this returns the nearest
// previous slice right-padded with 0xff up to length bytes. It may reuse the
// given slice when possible.
func BytesPrevish(b []byte, length int) []byte {
	bLen := len(b)
	// An empty slice has no previous slice.
	if bLen == 0 {
		return b
	}
	// If the last byte is 0, just remove it.
	if b[bLen-1] == 0 {
		return b[:bLen-1]
	}
	// Otherwise, decrement the last byte and right-pad with 0xff.
	if bLen > length {
		length = bLen
	}
	buf := make([]byte, length)
	copy(buf, b)
	buf[bLen-1]--
	copy(buf[bLen:], bytes.Repeat([]byte{0xff}, length-bLen))
	return buf
}

// unsafeWrapper is implementation of SafeFormatter. This is used to mark
// arguments as unsafe for redaction. This would make sure that redact.Unsafe() is implementing SafeFormatter interface
// without affecting invocations.
// TODO(aa-joshi): This is a temporary solution to mark arguments as unsafe. We should move/update this into cockroachdb/redact package.
type unsafeWrapper struct {
	a any
}

func (uw unsafeWrapper) SafeFormat(w redact.SafePrinter, _ rune) {
	w.Print(redact.Unsafe(uw.a))
}

func Unsafe(args any) any {
	return unsafeWrapper{a: args}
}
