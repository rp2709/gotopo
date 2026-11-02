package viewsheds_generation

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"math/bits"
)

type vector3 struct {
	x, y, z float64
}

// computeDirectionsVectors distributes direction vectors over a hemisphere.
//
// The number of elevation steps is floor(log2(2*sizeBytes + 1)).
// Each elevation step doubles the number of azimuth samples.
func computeDirectionsVectors(
	deadzone float64,
	zenithDeadzone float64,
	sizeBytes int64,
) ([]vector3, error) {
	if sizeBytes <= 0 {
		return nil, fmt.Errorf("sizeBytes must be greater than zero")
	}

	if math.IsNaN(deadzone) || math.IsInf(deadzone, 0) ||
		math.IsNaN(zenithDeadzone) || math.IsInf(zenithDeadzone, 0) {
		return nil, fmt.Errorf("deadzone angles must be finite")
	}

	if deadzone < 0 || zenithDeadzone < 0 ||
		deadzone+zenithDeadzone >= math.Pi/2 {
		return nil, fmt.Errorf(
			"deadzone angles must be non-negative and sum to less than pi/2")
	}

	elevationStepCount := bits.Len64(2*uint64(sizeBytes)+1) - 1

	if elevationStepCount > bits.UintSize-3 {
		return nil, fmt.Errorf(
			"--size-bytes is too large to index the resulting directions on this platform")
	}

	elevationStep := (math.Pi/2 - deadzone - zenithDeadzone) /
		float64(elevationStepCount)

	elevationOffset := deadzone + elevationStep/2

	var directions []vector3

	for zsIndex := range elevationStepCount {
		elevationAngle := float64(zsIndex)*elevationStep +
			elevationOffset

		azimuthStepCount := 4 << zsIndex
		azimuthStep := 2 * math.Pi / float64(azimuthStepCount)
		azimuthOffset := azimuthStep / 2

		for asIndex := range azimuthStepCount {
			azimuthAngle := float64(asIndex)*azimuthStep +
				azimuthOffset

			directions = append(
				directions,
				sphericalToVector3(azimuthAngle, elevationAngle),
			)
		}
	}

	return directions, nil
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
