package viewsheds_generation

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
)

type vector3 struct {
	x, y, z float64
}

// sphericalToVector3 converts spherical angles into a Cartesian vector.
// Both angles are measured in radians.
func sphericalToVector3(azimuth, zenithAngle float64) vector3 {
	return vector3{
		x: math.Sin(zenithAngle) * math.Cos(azimuth),
		y: math.Sin(zenithAngle) * math.Sin(azimuth),
		z: math.Cos(zenithAngle),
	}
}

// writeDirections prints list of coma separated triples.
func writeDirections(w io.Writer, directions []vector3) error {
	writer := bufio.NewWriter(w)
	print('[')
	for _, v := range directions {
		if _, err := fmt.Fprintf(
			writer, "(%.17g,%.17g,%.17g),", v.x, v.y, v.z,
		); err != nil {
			return err
		}
	}
	print(']')

	return writer.Flush()
}

func VisualizeViewshed(
	filename string,
	bits pgtype.Bits,
	size int,
	method PackingMethod,
) error {
	// Validate input.
	if !bits.Valid {
		return fmt.Errorf("viewshed bit string is NULL")
	}

	if bits.Len <= 0 {
		return fmt.Errorf("viewshed bit string must contain at least one bit")
	}

	if size == 0 {
		return fmt.Errorf("image size must be greater than zero")
	}

	if len(method.Directions()) != int(bits.Len) {
		return fmt.Errorf(
			"direction count mismatch: expected %d, got %d",
			bits.Len,
			len(method.Directions()),
		)
	}

	// Create a square grayscale image.
	img := image.NewGray(image.Rect(0, 0, int(size), int(size)))

	center := float64(size-1) / 2
	radius := float64(size) / 2

	for y := 0; y < int(size); y++ {
		for x := 0; x < int(size); x++ {

			// Map the pixel onto the projected hemisphere.
			// Image Y increases downwards, hence the inverted Y coordinate.
			px := (float64(x) - center) / radius
			py := (center - float64(y)) / radius

			// Distance from the center of the projected hemisphere.
			radiusSquared := px*px + py*py

			// Pixels outside the hemisphere remain black.
			if radiusSquared > 1 {
				continue
			}

			// Orthographic projection onto the upper unit hemisphere.
			pz := math.Sqrt(1 - radiusSquared)

			// Minimum angle from horizon: the outer circular band is always obstructed.
			angleFromHorizon := math.Asin(pz)
			if angleFromHorizon < minimumAngleFromHorizon {
				img.SetGray(x, y, color.Gray{Y: 75})
				continue
			}

			// Minimum angle from zenith: the central disk is always visible.
			zenithAngle := math.Acos(pz)
			if zenithAngle < minimumAngleFromZenith {
				img.SetGray(x, y, color.Gray{Y: 255})
				continue
			}

			pixelDirection := vector3{
				x: px,
				y: py,
				z: pz,
			}

			// Find the sampled direction with the smallest angular distance.
			closestIndex := 0
			maxDotProduct := math.Inf(-1)

			for i, direction := range method.Directions() {
				dotProduct :=
					pixelDirection.x*direction.x +
						pixelDirection.y*direction.y +
						pixelDirection.z*direction.z

				if dotProduct > maxDotProduct {
					maxDotProduct = dotProduct
					closestIndex = i
				}
			}

			// PostgreSQL bit strings are stored MSB-first.
			// Bit 0 is the most significant bit of the first byte.
			bitIndex := int32(closestIndex)
			bitMask := byte(0x80 >> uint(bitIndex%8))

			isVisible := bits.Bytes[bitIndex/8]&bitMask != 0

			gray := uint8(240) // White for 0.
			if isVisible {
				gray = 128 // Gray for 1.
			}

			img.SetGray(x, y, color.Gray{Y: gray})
		}
	}

	// Encode and save as PNG.
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create image file: %w", err)
	}

	if err := png.Encode(file, img); err != nil {
		_ = file.Close()
		return fmt.Errorf("encode PNG: %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("close image file: %w", err)
	}

	return nil
}
