package viewsheds_generation

import (
	"math"

	"github.com/jackc/pgx/v5/pgtype"
)

// EvenAnglePacking maps a grid of directions (nbStepsZenith rows x
// nbStepsAzimuth columns) onto the bits of a Postgres bit string.
//
// Bit index = zenithIndex*nbStepsAzimuth + azimuthIndex, bit 0 being the
// MOST significant bit of byte 0 (Postgres / pgtype.Bits order, same as
// SQL get_bit(bits, n)).
//
// Assumes smallestPairFactors returns integer types convertible to int.
type EvenAnglePacking struct {
	sizeBits                int
	minimumAngleFromZenith  float64
	minimumAngleFromHorizon float64

	nbStepsAzimuth     int
	nbStepsZenith      int
	zenithAngleToCover float64
	zenithStep         float64
	azimuthStep        float64
	zenithOffset       float64
	azimuthOffset      float64

	directions []vector3
}

// NewEvenAnglePacking computes every derived value once, so the value is
// immutable afterwards (value receivers are then fine and it is safe to
// share between goroutines).
func NewEvenAnglePacking(sizeBits int, minAngleFromZenith, minAngleFromHorizon float64) EvenAnglePacking {
	evp := EvenAnglePacking{
		sizeBits:                sizeBits,
		minimumAngleFromZenith:  minAngleFromZenith,
		minimumAngleFromHorizon: minAngleFromHorizon,
	}

	evp.nbStepsAzimuth, evp.nbStepsZenith = smallestPairFactors(sizeBits)

	evp.zenithAngleToCover = math.Pi/2 - minAngleFromZenith - minAngleFromHorizon
	evp.zenithStep = evp.zenithAngleToCover / float64(evp.nbStepsZenith)
	evp.azimuthStep = 2 * math.Pi / float64(evp.nbStepsAzimuth)
	evp.zenithOffset = minAngleFromHorizon + evp.zenithStep/2
	evp.azimuthOffset = evp.azimuthStep / 2

	evp.directions = make([]vector3, 0, sizeBits)
	for zi := range evp.nbStepsZenith {
		for ai := range evp.nbStepsAzimuth {
			evp.directions = append(evp.directions,
				sphericalToVector3(evp.azimuthAt(ai), math.Pi/2-evp.zenithAt(zi)))
		}
	}
	return evp
}

// Directions returns the direction of each bit, in bit order.
func (evp EvenAnglePacking) Directions() []vector3 {
	return evp.directions
}

// Mask returns a bit string with only the bit of the cell nearest to
// (azimuth, zenith) set. Azimuth wraps around; zenith is clamped to the
// covered range.
func (evp EvenAnglePacking) Mask(azimuth float64, zenith float64) pgtype.Bits {
	fa, fz := evp.gridCoords(azimuth, zenith)
	ai := int(math.Round(fa))
	zi := max(0, min(int(math.Round(fz)), evp.nbStepsZenith-1))

	bytes := make([]byte, (evp.sizeBits+7)/8)
	setBit(bytes, evp.bitIndex(ai, zi))

	return pgtype.Bits{
		Bytes: bytes,
		Len:   int32(evp.sizeBits),
		Valid: true,
	}
}

// Interpolate bilinearly interpolates the 0/1 bits of the four cells
// surrounding (azimuth, zenith). The result is in [0, 1]: exactly 0 or 1 at
// a cell centre, a fraction in between. Azimuth wraps; zenith is clamped
// (no extrapolation past the outermost cell centres).
// Returns NaN if viewshed is invalid or has the wrong length.
func (evp EvenAnglePacking) Interpolate(viewshed pgtype.Bits, azimuth float64, zenith float64) float64 {
	if !viewshed.Valid || int(viewshed.Len) != evp.sizeBits || len(viewshed.Bytes)*8 < evp.sizeBits {
		return math.NaN()
	}

	fa, fz := evp.gridCoords(azimuth, zenith)
	fz = max(0, min(fz, float64(evp.nbStepsZenith-1)))

	a0, z0 := math.Floor(fa), math.Floor(fz)
	ta, tz := fa-a0, fz-z0
	ai, zi0 := int(a0), int(z0)
	zi1 := min(zi0+1, evp.nbStepsZenith-1)

	bit := func(ai, zi int) float64 {
		if getBit(viewshed.Bytes, evp.bitIndex(ai, zi)) {
			return 1
		}
		return 0
	}

	low := lerp(bit(ai, zi0), bit(ai+1, zi0), ta)
	high := lerp(bit(ai, zi1), bit(ai+1, zi1), ta)
	return lerp(low, high, tz)
}

// azimuthAt / zenithAt: angle of the center of a cell (used by Directions).
func (evp EvenAnglePacking) azimuthAt(ai int) float64 {
	return evp.azimuthOffset + float64(ai)*evp.azimuthStep
}

func (evp EvenAnglePacking) zenithAt(zi int) float64 {
	return evp.zenithOffset + float64(zi)*evp.zenithStep
}

// gridCoords is the inverse of azimuthAt/zenithAt: the fractional grid
// position of an angle, where integer values are cell centres.
func (evp EvenAnglePacking) gridCoords(azimuth, zenith float64) (fa, fz float64) {
	return (azimuth - evp.azimuthOffset) / evp.azimuthStep,
		(zenith - evp.zenithOffset) / evp.zenithStep
}

// bitIndex converts grid indices to a bit index. Azimuth wraps around
// (any int, including negative); zenith must already be in range.
func (evp EvenAnglePacking) bitIndex(ai, zi int) int {
	n := evp.nbStepsAzimuth
	return zi*n + ((ai%n)+n)%n
}

// setBit / getBit use Postgres bit order: bit 0 is the MSB of byte 0.
func setBit(b []byte, i int) { b[i>>3] |= byte(0x80) >> (i & 7) }

func getBit(b []byte, i int) bool { return b[i>>3]&(byte(0x80)>>(i&7)) != 0 }

func lerp(a, b, t float64) float64 { return a + (b-a)*t }

func smallestPairFactors(n int) (g, p int) {
	smallestSum := int(math.MaxInt)
	sqrt := int(math.Sqrt(float64(n)))

	for i := 1; i <= sqrt; i++ {
		if n%i == 0 {
			dividend := n / i
			sum := i + dividend
			if sum < smallestSum {
				smallestSum = sum
				p = i
				g = dividend
			}
		}
	}
	return g, p
}
