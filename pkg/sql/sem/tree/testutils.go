// Copyright 2017 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package tree

import (
	"time"

	"github.com/cockroachdb/cockroachdb-parser/pkg/geo"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/types"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/duration"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/timeofday"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/timeutil"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/timeutil/pgdate"
	"github.com/cockroachdb/errors"
)

// presetTypesForTesting is a mapping of qualified names to types that can be mocked out
// for tests to allow the qualified names to be type checked without throwing an error.
var presetTypesForTesting map[string]*types.T

// MockNameTypes populates presetTypesForTesting for a test.
func MockNameTypes(types map[string]*types.T) func() {
	presetTypesForTesting = types
	return func() {
		presetTypesForTesting = nil
	}
}

// SampleDatum is intended to be a more lightweight version of RandDatum for
// when you just need one consistent example of a datum.
func SampleDatum(t *types.T) Datum {
	switch t.Family() {
	case types.BitFamily:
		a, _ := NewDBitArrayFromInt(123, 40)
		return a
	case types.BoolFamily:
		return MakeDBool(true)
	case types.IntFamily:
		return NewDInt(123)
	case types.FloatFamily:
		f := DFloat(123.456)
		return &f
	case types.DecimalFamily:
		d := &DDecimal{}
		// int64(rng.Uint64()) to get negative numbers, too
		d.Decimal.SetFinite(3, 6)
		return d
	case types.StringFamily:
		return NewDString("Carl")
	case types.BytesFamily:
		return NewDBytes("Princess")
	case types.DateFamily:
		return NewDDate(pgdate.MakeCompatibleDateFromDisk(123123))
	case types.TimeFamily:
		return MakeDTime(timeofday.FromInt(789))
	case types.TimeTZFamily:
		return NewDTimeTZFromOffset(timeofday.FromInt(345), 5*60*60 /* OffsetSecs */)
	case types.TimestampFamily:
		return MustMakeDTimestamp(timeutil.Unix(123, 123), time.Second)
	case types.TimestampTZFamily:
		return MustMakeDTimestampTZ(timeutil.Unix(123, 123), time.Second)
	case types.IntervalFamily:
		i, _ := ParseDInterval(duration.IntervalStyle_POSTGRES, "1h1m1s")
		return i
	case types.UuidFamily:
		u, _ := ParseDUuidFromString("3189ad07-52f2-4d60-83e8-4a8347fef718")
		return u
	case types.INetFamily:
		i, _ := ParseDIPAddrFromINetString("127.0.0.1")
		return i
	case types.JsonFamily:
		j, _ := ParseDJSON(`{"a": "b"}`)
		return j
	case types.OidFamily:
		return NewDOidWithType(1009, t)
	case types.PGLSNFamily:
		return NewDPGLSN(0x1000000100)
	case types.RefCursorFamily:
		return NewDRefCursor("Wheezer")
	case types.Box2DFamily:
		b := geo.NewCartesianBoundingBox().AddPoint(1, 2).AddPoint(3, 4)
		return NewDBox2D(*b)
	case types.GeographyFamily:
		return NewDGeography(geo.MustParseGeographyFromEWKB([]byte("\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f")))
	case types.GeometryFamily:
		return NewDGeometry(geo.MustParseGeometryFromEWKB([]byte("\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f")))
	default:
		panic(errors.AssertionFailedf("SampleDatum not implemented for %s", t))
	}
}
