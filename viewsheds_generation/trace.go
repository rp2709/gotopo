package viewsheds_generation

import (
	"github.com/rp2709/gotopo/vrt_reader"
	"github.com/twpayne/go-geom"
)

type StepDistance func(step int) float64

func obstructed(ds *vrt_reader.RasterDataset, source geom.Point, direction vector3, maxDistance float64, maxAltitude float64, distanceAt StepDistance) bool {
	sourceAltitude := float64(ds.ValueAt(source.X(), source.Y()))
	step := 1
	distance := distanceAt(step)
	for distance < maxDistance {
		stepAltitude := sourceAltitude + direction.z*distance
		if stepAltitude > maxAltitude {
			break
		}
		localAltitude := float64(ds.ValueAt(source.X()+distance*direction.x, source.Y()+distance*direction.y))
		if stepAltitude <= localAltitude {
			return true // found some piece of terrain higher than our sight
		}
		step++
		distance = distanceAt(step)
	}
	return false
}

func LinearStepDistance(a float64, b float64) StepDistance {
	return func(step int) float64 {
		return a*float64(step) + b
	}
}

func QuadraticStepDistance(a float64, b float64, c float64) StepDistance {
	return func(step int) float64 {
		return a*float64(step*step) + b*float64(step) + c
	}
}
