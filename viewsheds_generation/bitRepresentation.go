package viewsheds_generation

import (
	"github.com/jackc/pgx/v5/pgtype"
)

/*
Abstraction of the method to pack the viewshed in a bitset
It defines all operations related to a specific representation
*/
type PackingMethod interface {
	/*
		Returns the direction vectors representing the direction in which we will compute obstruction
	*/
	Directions() []vector3
	/*
		Constructs a mask with only bit corresponding to the spherical coordinates
		This is intended to be used once and then passed to postgresql, which handles Bits efficiently.
		It is faster to use Interpolate for direct use since there is no allocation of pgtype.Bits
	*/
	Mask(azimuth float64, zenith float64) pgtype.Bits
	/*
		Bilinear interpolation of the value of the viewshed for some spherical coordinates
		Returns a number between 0 and 1 inclusive representing how obstructed a point on the hemisphere is, 0 being not at all and 1 completely obstructed
	*/
	Interpolate(viewshed pgtype.Bits, azimuth float64, zenith float64) float64
}

type sliceRepresentation []bool

func (sr sliceRepresentation) toBits() pgtype.Bits {
	bits := pgtype.Bits{}

	bits.Len = int32(len(sr))
	bits.Valid = true
	bits.Bytes = make([]byte, (int(bits.Len)+7)/8)

	for i, b := range sr {
		if b {
			bits.Bytes[i/8] |= 1 << uint(7-i%8)
		}
	}

	return bits
}
