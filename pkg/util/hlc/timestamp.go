// Copyright 2014 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package hlc

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"
	"unsafe"

	"github.com/cockroachdb/cockroachdb-parser/pkg/util/timeutil"
	"github.com/cockroachdb/errors"
)

// Timestamp constant values.
var (
	// MaxTimestamp is the max value allowed for Timestamp.
	MaxTimestamp = Timestamp{WallTime: math.MaxInt64, Logical: math.MaxInt32}
	// MinTimestamp is the min value allowed for Timestamp.
	MinTimestamp = Timestamp{WallTime: 0, Logical: 1}
	// MaxClockTimestamp is the max value allowed for ClockTimestamp.
	MaxClockTimestamp = ClockTimestamp{WallTime: math.MaxInt64, Logical: math.MaxInt32}
	// MinClockTimestamp is the min value allowed for ClockTimestamp.
	MinClockTimestamp = ClockTimestamp{WallTime: 0, Logical: 1}
)

// EqOrdering returns whether the receiver sorts equally to the parameter.
//
// This method is split from tests of structural equality (Equal and the equals
// operator) because it does not consider differences in flags and only
// considers whether the walltime and logical time differ between the
// timestamps.
//
// TODO(nvanbenschoten): remove this method and its uses. It's now equivalent to
// a direct struct comparison.
func (t Timestamp) EqOrdering(s Timestamp) bool {
	return t.WallTime == s.WallTime && t.Logical == s.Logical
}

// Less returns whether the receiver is less than the parameter.
func (t Timestamp) Less(s Timestamp) bool {
	return t.WallTime < s.WallTime || (t.WallTime == s.WallTime && t.Logical < s.Logical)
}

// LessEq returns whether the receiver is less than or equal to the parameter.
func (t Timestamp) LessEq(s Timestamp) bool {
	return t.WallTime < s.WallTime || (t.WallTime == s.WallTime && t.Logical <= s.Logical)
}

// After returns whether the provided timestamp is after our timestamp.
// This matches the behavior of time.After.
func (t Timestamp) After(s Timestamp) bool {
	return !t.LessEq(s)
}

// Compare returns -1 if this timestamp is lesser than the given timestamp, 1 if
// it is greater, and 0 if they are equal.
func (t Timestamp) Compare(s Timestamp) int {
	if t.WallTime > s.WallTime {
		return 1
	} else if t.WallTime < s.WallTime {
		return -1
	} else if t.Logical > s.Logical {
		return 1
	} else if t.Logical < s.Logical {
		return -1
	} else {
		return 0
	}
}

// String implements the fmt.Stringer interface.
// Outputs: seconds.nanos,logical
func (t Timestamp) String() string {
	// The following code was originally written as
	//   fmt.Sprintf("%d.%09d,%d", t.WallTime/1e9, t.WallTime%1e9, t.Logical).
	// The main problem with the original code was that it would put
	// a negative sign in the middle (after the decimal point) if
	// the value happened to be negative.
	buf := make([]byte, 0, 21)

	w := t.WallTime
	if w == 0 {
		buf = append(buf, '0', ',')
	} else {
		if w < 0 {
			w = -w
			buf = append(buf, '-')
		}

		s, ns := uint64(w/1e9), uint64(w%1e9)
		buf = strconv.AppendUint(buf, s, 10)

		prev := len(buf)
		buf = append(buf, '.', '0', '0', '0', '0', '0', '0', '0', '0', '0', ',')
		// The following lines expand to:
		// 		for i := 0; i < 9; i++ {
		// 			buf[len(buf)-2-i] = byte('0' + ns%10)
		// 			ns = ns / 10
		// 		}
		zeroBuf := buf[prev+1 : len(buf)-1]
		zeroBuf[8] = byte('0' + ns%10)
		ns = ns / 10
		zeroBuf[7] = byte('0' + ns%10)
		ns = ns / 10
		zeroBuf[6] = byte('0' + ns%10)
		ns = ns / 10
		zeroBuf[5] = byte('0' + ns%10)
		ns = ns / 10
		zeroBuf[4] = byte('0' + ns%10)
		ns = ns / 10
		zeroBuf[3] = byte('0' + ns%10)
		ns = ns / 10
		zeroBuf[2] = byte('0' + ns%10)
		ns = ns / 10
		zeroBuf[1] = byte('0' + ns%10)
		ns = ns / 10
		zeroBuf[0] = byte('0' + ns%10)
	}
	buf = strconv.AppendInt(buf, int64(t.Logical), 10)
	return *(*string)(unsafe.Pointer(&buf))
}

// SafeValue implements the redact.SafeValue interface.
func (Timestamp) SafeValue() {}

var (
	timestampRegexp = regexp.MustCompile(
		`^(?P<sign>-)?(?P<secs>\d{1,19})(?:\.(?P<nanos>\d{1,20}))?(?:,(?P<logical>-?\d{1,10}))?$`)
	signSubexp    = 1
	secsSubexp    = 2
	nanosSubexp   = 3
	logicalSubexp = 4
)

// ParseTimestamp attempts to parse the string generated from
// Timestamp.String().
func ParseTimestamp(str string) (_ Timestamp, err error) {
	matches := timestampRegexp.FindStringSubmatch(str)
	if matches == nil {
		return Timestamp{}, errors.Errorf("failed to parse %q as Timestamp", str)
	}
	defer func() {
		if err != nil {
			err = errors.Wrapf(err, "failed to parse %q as Timestamp", str)
		}
	}()
	seconds, err := strconv.ParseInt(matches[secsSubexp], 10, 64)
	if err != nil {
		return Timestamp{}, err
	}
	var nanos int64
	if nanosMatch := matches[nanosSubexp]; nanosMatch != "" {
		nanos, err = strconv.ParseInt(nanosMatch, 10, 64)
		if err != nil {
			return Timestamp{}, err
		}
	}
	wallTime := seconds*time.Second.Nanoseconds() + nanos
	if matches[signSubexp] == "-" {
		wallTime *= -1
	}
	var logical int64
	if logicalMatch := matches[logicalSubexp]; logicalMatch != "" {
		logical, err = strconv.ParseInt(logicalMatch, 10, 32)
		if err != nil {
			return Timestamp{}, err
		}
	}
	t := Timestamp{
		WallTime: wallTime,
		Logical:  int32(logical),
	}
	return t, nil
}

// AsOfSystemTime returns a string to be used in an AS OF SYSTEM TIME query.
func (t Timestamp) AsOfSystemTime() string {
	return fmt.Sprintf("%d.%010d", t.WallTime, t.Logical)
}

// IsEmpty returns true if t is an empty Timestamp.
// gcassert:inline
func (t Timestamp) IsEmpty() bool {
	// TODO(radu): consider making timestamps with zero wall time and non-zero
	// logical time illegal. Then we can check just the wall time.
	return t == Timestamp{}
}

// IsSet returns true if t is not an empty Timestamp.
// gcassert:inline
func (t Timestamp) IsSet() bool {
	return !t.IsEmpty()
}

// AddDuration adds a given duration to this Timestamp. Normally if you want to
// bump your clock to the higher of two timestamps, use Forward, however this
// method is here to create a hlc.Timestamp in the future (or past).
func (t Timestamp) AddDuration(duration time.Duration) Timestamp {
	return t.Add(duration.Nanoseconds(), 0)
}

// Add returns a timestamp with the WallTime and Logical components increased.
// wallTime is expressed in nanos.
func (t Timestamp) Add(wallTime int64, logical int32) Timestamp {
	return Timestamp{
		WallTime: t.WallTime + wallTime,
		Logical:  t.Logical + logical,
	}
}

// Clone return a new timestamp that has the same contents as the receiver.
func (t Timestamp) Clone() *Timestamp {
	return &t
}

// Next returns the timestamp with the next later timestamp.
func (t Timestamp) Next() Timestamp {
	if t.Logical == math.MaxInt32 {
		if t.WallTime == math.MaxInt64 {
			panic("cannot take the next value to a max timestamp")
		}
		return Timestamp{
			WallTime: t.WallTime + 1,
		}
	}
	return Timestamp{
		WallTime: t.WallTime,
		Logical:  t.Logical + 1,
	}
}

// Prev returns the next earliest timestamp.
func (t Timestamp) Prev() Timestamp {
	if t.Logical > 0 {
		return Timestamp{
			WallTime: t.WallTime,
			Logical:  t.Logical - 1,
		}
	} else if t.WallTime > 0 {
		return Timestamp{
			WallTime: t.WallTime - 1,
			Logical:  math.MaxInt32,
		}
	}
	panic("cannot take the previous value to a zero timestamp")
}

// FloorPrev returns a timestamp earlier than the current timestamp. If it
// can subtract a logical tick without wrapping around, it does so. Otherwise
// it subtracts a nanosecond from the walltime.
func (t Timestamp) FloorPrev() Timestamp {
	if t.Logical > 0 {
		return Timestamp{
			WallTime: t.WallTime,
			Logical:  t.Logical - 1,
		}
	} else if t.WallTime > 0 {
		return Timestamp{
			WallTime: t.WallTime - 1,
			Logical:  0,
		}
	}
	panic("cannot take the previous value to a zero timestamp")
}

// WallNext adds 1 to the WallTime and resets Logical.
func (t Timestamp) WallNext() Timestamp {
	return Timestamp{
		WallTime: t.WallTime + 1,
		Logical:  0,
	}
}

// WallPrev subtracts 1 from the WallTime and resets Logical.
func (t Timestamp) WallPrev() Timestamp {
	return Timestamp{
		WallTime: t.WallTime - 1,
		Logical:  0,
	}
}

// Forward replaces the receiver with the argument, if that moves it forwards in
// time. Returns true if the timestamp was adjusted to a larger time and false
// otherwise.
func (t *Timestamp) Forward(s Timestamp) bool {
	if t.Less(s) {
		*t = s
		return true
	}
	return false
}

// Backward replaces the receiver with the argument, if that moves it backwards
// in time.
func (t *Timestamp) Backward(s Timestamp) {
	if s.Less(*t) {
		*t = s
	}
}

// GoTime converts the timestamp to a time.Time.
func (t Timestamp) GoTime() time.Time {
	return timeutil.Unix(0, t.WallTime)
}

// ToLegacyTimestamp converts a Timestamp to a LegacyTimestamp.
func (t Timestamp) ToLegacyTimestamp() LegacyTimestamp { return LegacyTimestamp(t) }

// ToTimestamp converts a LegacyTimestamp to a Timestamp.
func (t LegacyTimestamp) ToTimestamp() Timestamp { return Timestamp(t) }

// Less returns whether the receiver is less than the parameter.
func (t LegacyTimestamp) Less(s LegacyTimestamp) bool {
	return t.ToTimestamp().Less(s.ToTimestamp())
}

// String implements the fmt.Stringer interface.
func (t LegacyTimestamp) String() string {
	return t.ToTimestamp().String()
}

// SafeValue implements the redact.SafeValue interface.
func (LegacyTimestamp) SafeValue() {}

// ClockTimestamp is a Timestamp with the added capability of being able to
// update a peer's HLC clock. It possesses this capability because the clock
// timestamp itself is guaranteed to have come from an HLC clock somewhere in
// the system. As such, a clock timestamp is a promise that some node in the
// system has a clock with a reading equal to or above its value.
type ClockTimestamp Timestamp

// UnsafeToClockTimestamp converts a Timestamp to a ClockTimestamp, regardless
// of whether the timestamp actually came from a clock. The method should only
// be used in tests.
func (t Timestamp) UnsafeToClockTimestamp() ClockTimestamp { return ClockTimestamp(t) }

// ToTimestamp upcasts a ClockTimestamp into a Timestamp.
func (t ClockTimestamp) ToTimestamp() Timestamp { return Timestamp(t) }

// Less returns whether the receiver is less than the parameter.
func (t ClockTimestamp) Less(s ClockTimestamp) bool { return Timestamp(t).Less(Timestamp(s)) }

// LessEq returns whether the receiver is less than or equal to the parameter.
func (t ClockTimestamp) LessEq(s ClockTimestamp) bool { return Timestamp(t).LessEq(Timestamp(s)) }

// Ignore unused warnings.
var _ = ClockTimestamp.LessEq

// String implements the fmt.Stringer interface.
func (t ClockTimestamp) String() string { return t.ToTimestamp().String() }

// SafeValue implements the redact.SafeValue interface.
func (t ClockTimestamp) SafeValue() {}

// IsEmpty returns true if t is an empty ClockTimestamp.
func (t ClockTimestamp) IsEmpty() bool { return Timestamp(t).IsEmpty() }

// Forward is like Timestamp.Forward, but for ClockTimestamps.
func (t *ClockTimestamp) Forward(s ClockTimestamp) bool { return (*Timestamp)(t).Forward(Timestamp(s)) }

// Backward is like Timestamp.Backward, but for ClockTimestamps.
func (t *ClockTimestamp) Backward(s ClockTimestamp) { (*Timestamp)(t).Backward(Timestamp(s)) }

// BackwardWithTimestamp is like Backward, but with a Timestamp parameter.
func (t *ClockTimestamp) BackwardWithTimestamp(s Timestamp) { (*Timestamp)(t).Backward(s) }

// Reset implements the protoutil.Message interface.
func (t *ClockTimestamp) Reset() { (*Timestamp)(t).Reset() }

// ProtoMessage implements the protoutil.Message interface.
func (t *ClockTimestamp) ProtoMessage() {}

// MarshalTo implements the protoutil.Message interface.
func (t *ClockTimestamp) MarshalTo(data []byte) (int, error) { return (*Timestamp)(t).MarshalTo(data) }

// MarshalToSizedBuffer implements the protoutil.Message interface.
func (t *ClockTimestamp) MarshalToSizedBuffer(data []byte) (int, error) {
	return (*Timestamp)(t).MarshalToSizedBuffer(data)
}

// Unmarshal implements the protoutil.Message interface.
func (t *ClockTimestamp) Unmarshal(data []byte) error { return (*Timestamp)(t).Unmarshal(data) }

// Size implements the protoutil.Message interface.
func (t *ClockTimestamp) Size() int { return (*Timestamp)(t).Size() }

// Equal is needed for the gogoproto.equal option.
func (t *ClockTimestamp) Equal(that interface{}) bool {
	switch v := that.(type) {
	case nil:
		return t == nil
	case ClockTimestamp:
		return (*Timestamp)(t).Equal((Timestamp)(v))
	case *ClockTimestamp:
		return (*Timestamp)(t).Equal((*Timestamp)(v))
	default:
		return false
	}
}

// NewPopulatedClockTimestamp is needed for the gogoproto.populate option.
func NewPopulatedClockTimestamp(r randyTimestamp, easy bool) *ClockTimestamp {
	return (*ClockTimestamp)(NewPopulatedTimestamp(r, easy))
}

// Ignore unused warnings. The function is called, but by generated functions
// that themselves are unused.
var _ = NewPopulatedClockTimestamp
