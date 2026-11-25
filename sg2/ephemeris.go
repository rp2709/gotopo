package sg2

import (
	"math"
	"time"
)

const (
	secondsPerDay = 86400.0
	julianUnix    = 2440587.5
	julian1980    = 2444239.5 // 1980-01-01 00:00:00 UT

	// SG2 constants from Blanc and Wald (2012), and the corrected reference
	// implementation in solposx.
	stellarAberration = -9.933735e-5 // radians
	solarParallax     = 4.263521e-5  // radians

	// Reference ellipsoid used by SG2 for the topocentric parallax correction.
	earthSemiMajor  = 6378140.0
	earthFlattening = 1.0 / 298.257282697
)

type sinusoid struct {
	period    float64 // days
	amplitude float64 // radians
	phase     float64 // radians
}

// The ten terms of SG2's fitted Earth heliocentric longitude, Table 6.
var earthLongitudeTerms = [...]sinusoid{
	{period: 365.261278, amplitude: 3.401508e-2, phase: 1.600780},
	{period: 182.632412, amplitude: 3.486440e-4, phase: 1.662976},
	{period: 29.530634, amplitude: 3.136227e-5, phase: -1.195905},
	{period: 399.529850, amplitude: 3.578979e-5, phase: -1.042052},
	{period: 291.956812, amplitude: 2.676185e-5, phase: 2.012613},
	{period: 583.598201, amplitude: 2.333925e-5, phase: -2.867714},
	{period: 4652.629372, amplitude: 1.221214e-5, phase: 1.225038},
	{period: 1450.236684, amplitude: 1.217941e-5, phase: -0.828601},
	{period: 199.459709, amplitude: 1.343914e-5, phase: -3.108253},
	{period: 365.355291, amplitude: 8.499475e-4, phase: -2.353709},
}

func sunPosition(latitudeDeg, longitudeDeg, heightM float64, instant time.Time) (azimuth, altitude float64) {
	utc := instant.UTC()
	jdUT := julianDay(utc)
	jUT := jdUT - julian1980

	// The corrected SG2 implementation selects the Delta-T polynomial using
	// the mid-month decimal year, but evaluates it using the integer year.
	deltaT := deltaTSG2(utc.Year(), utc.Month())
	jTT := jUT + deltaT/secondsPerDay

	// Earth heliocentric longitude (SG2 Eq. 26).
	earthLongitude := jTT/58.130101 + 1.742145
	for _, term := range earthLongitudeTerms {
		earthLongitude += term.amplitude * math.Cos(2*math.Pi*jTT/term.period-term.phase)
	}
	earthLongitude = wrapRadians(earthLongitude)

	// Nutation in longitude and true obliquity (SG2 Eqs. 29–30).
	nutation := 8.329092e-5 * math.Cos(2*math.Pi*jTT/6791.164405+2.052757)
	obliquity := -6.216374e-9*jTT + 4.091383e-1 +
		4.456183e-5*math.Cos(2*math.Pi*jTT/6791.164405-2.660352)

	// Apparent geocentric solar longitude, declination, and right ascension.
	sunLongitude := wrapRadians(earthLongitude + math.Pi + nutation + stellarAberration)
	geocentricDeclination := math.Asin(clamp(math.Sin(sunLongitude)*math.Sin(obliquity), -1, 1))
	geocentricRA := math.Atan2(
		math.Sin(sunLongitude)*math.Cos(obliquity),
		math.Cos(sunLongitude),
	)

	// Mean sidereal time. The 6.300388099 coefficient is the corrected SG2
	// value; the journal paper prints a typographical error in this equation.
	meanSidereal := wrapRadians(6.300388099*jUT + 1.742079)

	latitude := degreesToRadians(latitudeDeg)
	longitude := degreesToRadians(longitudeDeg)

	// Observer coordinates and solar parallax (SG2 Eqs. 32–35).
	u := math.Atan((1 - earthFlattening) * math.Tan(latitude))
	x := math.Cos(u) + (heightM/earthSemiMajor)*math.Cos(latitude)
	y := (1-earthFlattening)*math.Sin(u) + (heightM/earthSemiMajor)*math.Sin(latitude)

	// Geocentric hour angle and topocentric parallax corrections.
	geocentricHourAngle := meanSidereal + nutation*math.Cos(obliquity) - geocentricRA + longitude
	deltaRA := -x * math.Sin(geocentricHourAngle) / math.Cos(geocentricDeclination) * solarParallax
	topocentricDeclination := geocentricDeclination +
		(x*math.Cos(geocentricHourAngle)*math.Sin(geocentricDeclination)-
			y*math.Cos(geocentricDeclination))*solarParallax
	topocentricHourAngle := geocentricHourAngle - deltaRA

	// Topocentric azimuth, eastward from true north (SG2 Eq. 20).
	azimuthRad := math.Atan2(
		math.Sin(topocentricHourAngle),
		math.Cos(topocentricHourAngle)*math.Sin(latitude)-
			math.Tan(topocentricDeclination)*math.Cos(latitude),
	) + math.Pi
	azimuthRad = wrapRadians(azimuthRad)

	// Geometric topocentric altitude (SG2 Eq. 21). Atmospheric refraction is
	// deliberately omitted: the core SG2 sun_position output is gamma_S0.
	altitudeRad := math.Asin(clamp(
		math.Sin(latitude)*math.Sin(topocentricDeclination)+
			math.Cos(latitude)*math.Cos(topocentricDeclination)*math.Cos(topocentricHourAngle),
		-1, 1,
	))

	return radiansToDegrees(azimuthRad), radiansToDegrees(altitudeRad)
}

func julianDay(t time.Time) float64 {
	return julianUnix + (float64(t.Unix())+float64(t.Nanosecond())*1e-9)/secondsPerDay
}

func deltaTSG2(year int, month time.Month) float64 {
	decimalYear := float64(year) + (float64(month)-0.5)/12.0
	switch {
	case decimalYear <= 1986.0:
		x := float64(year) - 1975.0
		return 45.45 + 1.067*x - x*x/260.0 - x*x*x/718.0
	case decimalYear <= 2005.0:
		x := float64(year) - 2000.0
		return 63.86 + 0.3345*x - 0.060374*x*x + 0.0017275*x*x*x +
			6.518e-4*x*x*x*x + 2.374e-5*x*x*x*x*x
	default:
		x := float64(year) - 2000.0
		return 63.48 + 0.2040*x + 0.005576*x*x
	}
}

func degreesToRadians(degrees float64) float64 { return degrees * math.Pi / 180.0 }
func radiansToDegrees(radians float64) float64 { return radians * 180.0 / math.Pi }

func wrapRadians(angle float64) float64 {
	wrapped := math.Mod(angle, 2*math.Pi)
	if wrapped < 0 {
		wrapped += 2 * math.Pi
	}
	return wrapped
}

func clamp(v, minValue, maxValue float64) float64 {
	if v < minValue {
		return minValue
	}
	if v > maxValue {
		return maxValue
	}
	return v
}
