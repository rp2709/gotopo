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

// computeDirectionsVectors distributes direction vectors over a hemisphere.
func computeDirectionsVectors(
	minimumAngleFromHorizon float64,
	minimumAngleFromZenith float64,
	sizeBits uint32,
) ([]vector3, error) {
	directions := make([]vector3, 0, sizeBits)

	nbStepsAzimuth, nbStepsZenith := smallestPairFactors(sizeBits)
	zenithAngleToCover := math.Pi/2 - minimumAngleFromZenith - minimumAngleFromHorizon
	zenithStep := zenithAngleToCover / float64(nbStepsZenith)
	azimuthStep := math.Pi * 2 / float64(nbStepsAzimuth)
	zenithOffset := minimumAngleFromHorizon + zenithStep/2
	azimuthOffset := azimuthStep / 2

	for zi := range nbStepsZenith {
		zenithAngle := zenithOffset + float64(zi)*zenithStep
		for ai := range nbStepsAzimuth {
			azimuthAngle := azimuthOffset + float64(ai)*azimuthStep
			directions = append(directions, sphericalToVector3(azimuthAngle, math.Pi/2-zenithAngle))
		}
	}

	return directions, nil
}

func smallestPairFactors(n uint32) (g, p uint32) {
	smallestSum := uint32(math.MaxUint32)
	sqrt := uint32(math.Sqrt(float64(n)))

	for i := uint32(1); i <= sqrt; i++ {
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
	size uint32,
	minimumAngleFromHorizon float64,
	minimumAngleFromZenith float64,
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

	if uint64(size) > uint64(^uint(0)>>1) {
		return fmt.Errorf("image size exceeds maximum integer size")
	}

	expectedBytes := (int64(bits.Len) + 7) / 8
	if int64(len(bits.Bytes)) != expectedBytes {
		return fmt.Errorf(
			"invalid bit string: expected %d bytes, got %d",
			expectedBytes,
			len(bits.Bytes),
		)
	}

	// Generate exactly one direction vector per bit.
	directions, err := computeDirectionsVectors(
		minimumAngleFromHorizon,
		minimumAngleFromZenith,
		uint32(bits.Len),
	)
	if err != nil {
		return fmt.Errorf("compute direction vectors: %w", err)
	}

	if len(directions) != int(bits.Len) {
		return fmt.Errorf(
			"direction count mismatch: expected %d, got %d",
			bits.Len,
			len(directions),
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

			pixelDirection := vector3{
				x: px,
				y: py,
				z: pz,
			}

			// Find the sampled direction with the smallest angular distance.
			closestIndex := 0
			maxDotProduct := math.Inf(-1)

			for i, direction := range directions {
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

			gray := uint8(255) // White for 0.
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
