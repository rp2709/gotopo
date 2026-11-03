package viewsheds_generation

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"math/bits"
	"os"
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

// RenderViewshed renders a hemispherical viewshed using cell-based
// rasterization.
//
// directions must contain obstruction results in the same order as
// computeDirectionsVectors.
//
// true  = obstructed
// false = visible
//
// size is the width and height of the resulting square image.
//
// deadzone excludes the area around the zenith.
// zenithDeadzone excludes the area around the horizon.
func RenderViewshed(
	directions []bool,
	size int,
	deadzone float64,
	zenithDeadzone float64,
) (*image.RGBA, error) {
	if size < 2 {
		return nil, fmt.Errorf("size must be at least 2")
	}

	if math.IsNaN(deadzone) || math.IsInf(deadzone, 0) ||
		math.IsNaN(zenithDeadzone) || math.IsInf(zenithDeadzone, 0) {
		return nil, fmt.Errorf("deadzone angles must be finite")
	}

	if deadzone < 0 || zenithDeadzone < 0 ||
		deadzone+zenithDeadzone >= math.Pi/2 {
		return nil, fmt.Errorf(
			"deadzone angles must be non-negative and sum to less than pi/2",
		)
	}

	// Recover the number of elevation bands.
	//
	// Total directions = 4 * (2^n - 1)
	q := len(directions)/4 + 1

	if len(directions) == 0 ||
		len(directions)%4 != 0 ||
		q&(q-1) != 0 {
		return nil, fmt.Errorf(
			"invalid direction count: got %d",
			len(directions),
		)
	}

	elevationStepCount := bits.Len(uint(q)) - 1

	elevationStep := (math.Pi/2 - deadzone - zenithDeadzone) /
		float64(elevationStepCount)

	// Precompute the starting index of each elevation band.
	bandOffsets := make([]int, elevationStepCount)

	offset := 0

	for band := 0; band < elevationStepCount; band++ {
		bandOffsets[band] = offset
		offset += 4 << band
	}

	// Colors.
	visibleColor := color.RGBA{
		R: 65, G: 180, B: 110, A: 255,
	}
	obstructedColor := color.RGBA{
		R: 210, G: 75, B: 65, A: 255,
	}
	deadzoneColor := color.RGBA{
		R: 65, G: 65, B: 70, A: 255,
	}
	backgroundColor := color.RGBA{
		R: 30, G: 30, B: 35, A: 255,
	}

	img := image.NewRGBA(image.Rect(0, 0, size, size))

	cx := float64(size-1) / 2
	cy := float64(size-1) / 2
	radius := math.Min(cx, cy)

	maxZenithAngle := math.Pi/2 - zenithDeadzone

	// Rasterize the hemisphere.
	//
	// Each pixel is mapped back to spherical coordinates.
	// Its angular coordinates identify exactly one cell.
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x) - cx
			dy := cy - float64(y)

			// Radial distance from the zenith.
			r := math.Hypot(dx, dy)

			if r > radius {
				img.SetRGBA(x, y, backgroundColor)
				continue
			}

			// Convert radial distance to zenith angle.
			zenithAngle := (r / radius) * (math.Pi / 2)

			// Outside the sampled hemisphere.
			if zenithAngle < deadzone ||
				zenithAngle >= maxZenithAngle {
				img.SetRGBA(x, y, deadzoneColor)
				continue
			}

			// Determine the elevation band.
			band := int((zenithAngle - deadzone) / elevationStep)

			if band >= elevationStepCount {
				band = elevationStepCount - 1
			}

			// Recover azimuth.
			//
			// The Y coordinate is inverted because image
			// coordinates increase downwards.
			azimuth := math.Atan2(dy, dx)

			if azimuth < 0 {
				azimuth += 2 * math.Pi
			}

			// Each elevation band has its own azimuth resolution.
			azimuthCount := 4 << band
			azimuthStep := 2 * math.Pi / float64(azimuthCount)

			// Azimuth cell boundaries are multiples of azimuthStep.
			// Samples themselves are located at cell centers.
			azimuthIndex := int(azimuth / azimuthStep)

			// Protect against floating-point rounding at 2*pi.
			if azimuthIndex >= azimuthCount {
				azimuthIndex = 0
			}

			// Find the corresponding boolean.
			directionIndex := bandOffsets[band] + azimuthIndex

			if directions[directionIndex] {
				img.SetRGBA(x, y, obstructedColor)
			} else {
				img.SetRGBA(x, y, visibleColor)
			}
		}
	}

	return img, nil
}

// SaveViewshedPNG renders and saves a viewshed as PNG.
func SaveViewshedPNG(
	filename string,
	directions []bool,
	size int,
	deadzone float64,
	zenithDeadzone float64,
) error {
	img, err := RenderViewshed(
		directions,
		size,
		deadzone,
		zenithDeadzone,
	)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
