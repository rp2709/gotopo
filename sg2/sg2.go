// Package sg2 computes the topocentric geometric position of the Sun using
// the Solar Geometry 2 (SG2) algorithm of Blanc and Wald (2012).
//
// Geographic points use the conventional GIS axis order: X is longitude and
// Y is latitude, both in decimal degrees. If the point has a Z coordinate it
// is interpreted as ellipsoidal height in metres; otherwise height is zero.
//
// SunPosition returns azimuth clockwise from true north in [0, 360) degrees
// and geometric altitude in degrees (positive above the horizon). Altitude is
// not corrected for atmospheric refraction, matching the core SG2 output
// (topocentric gamma_S0). The supplied time is interpreted as an absolute
// instant and converted to UTC; UTC approximates UT1 for this implementation.
//
// SG2's published accuracy is approximately 0.003 degrees over 1980–2030.
// The direct trigonometric series is evaluated without the optional
// precomputed daily interpolation tables.
package sg2

import (
	"fmt"
	"math"
	"time"

	"github.com/twpayne/go-geom"
)

// SunPosition returns the Sun's azimuth and geometric altitude, in degrees,
// for a geographic point and an absolute time.
//
// The geom.Point ordinates are longitude (X), latitude (Y), and optional
// ellipsoidal height in metres (Z). Empty points, non-finite coordinates, or
// coordinates outside the geographic ranges return (NaN, NaN). Longitude is
// restricted to [-180, 180] and latitude to [-90, 90].
//
// No error can be returned because the requested API has no error result.
func SunPosition(point geom.Point, instant time.Time) (azimuth, altitude float64) {
	invalid := math.NaN()
	if point.Empty() {
		return invalid, invalid
	}

	if point.SRID() == _SRIDLV95 {
		npoint, err := toWGS84(&point)
		if err != nil {
			return invalid, invalid
		}
		point = *npoint
	}

	longitude := point.X()
	latitude := point.Y()
	height := point.Z()
	if !finite(longitude) || !finite(latitude) || !finite(height) ||
		longitude < -180 || longitude > 180 || latitude < -90 || latitude > 90 {
		return invalid, invalid
	}

	return sunPosition(latitude, longitude, height, instant)
}

const (
	_SRIDLV95  = 2056
	_SRIDWGS84 = 4326
)

// ToWGS84 converts a point from Swiss LV95 (EPSG:2056)
// to WGS84 longitude/latitude (EPSG:4326).
func toWGS84(point *geom.Point) (*geom.Point, error) {
	if point == nil {
		return nil, fmt.Errorf("point cannot be nil")
	}

	if point.SRID() != _SRIDLV95 {
		return nil, fmt.Errorf("expected SRID %d, got %d",
			_SRIDLV95, point.SRID())
	}

	// LV95 coordinates:
	// X = Easting
	// Y = Northing
	e := point.X()
	n := point.Y()

	// Convert to the auxiliary Swiss coordinate system.
	y := (e - 2600000.0) / 1000000.0
	x := (n - 1200000.0) / 1000000.0

	// Longitude in units of 10,000 arc-seconds.
	lon := 2.6779094 +
		4.728982*y +
		0.791484*y*x +
		0.1306*y*x*x -
		0.0436*y*y*y

	// Latitude in units of 10,000 arc-seconds.
	lat := 16.9023892 +
		3.238272*x -
		0.270978*y*y -
		0.002528*x*x -
		0.0447*y*y*x -
		0.0140*x*x*x

	// Convert to decimal degrees.
	lon *= 100.0 / 36.0
	lat *= 100.0 / 36.0

	// Return a new point (longitude, latitude).
	return geom.NewPoint(geom.XY).
		SetSRID(_SRIDWGS84).
		MustSetCoords(geom.Coord{lon, lat}), nil
}

func finite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}
