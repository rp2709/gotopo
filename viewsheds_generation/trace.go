package viewsheds_generation

import (
	"math"

	"github.com/rp2709/gotopo/vrt_reader"
	"github.com/twpayne/go-geom"
)

func obstructed(ds *vrt_reader.RasterDataset, source geom.Point, direction vector3, maxDistance float64, maxAltitude float64, nextDistance func(float64) float64) bool {
	sourceAltitude := float64(ds.ValueAt(source.X(), source.Y()))
	horizontalDirectionLength := math.Sqrt(direction.x*direction.x + direction.y*direction.y)
	if horizontalDirectionLength == 0 {
		return false
	}

	for distance := 0.5; distance < maxDistance; distance = nextDistance(distance) {
		distanceAlongDirection := distance / horizontalDirectionLength
		stepAltitude := sourceAltitude + direction.z*distanceAlongDirection
		if stepAltitude > maxAltitude {
			break
		}
		localAltitude := float64(ds.ValueAt(source.X()+distanceAlongDirection*direction.x, source.Y()+distanceAlongDirection*direction.y))
		if stepAltitude <= localAltitude {
			return true // found some piece of terrain higher than our sight
		}
	}
	return false
}

func viewshed(ds *vrt_reader.RasterDataset, point geom.Point, directions []vector3) []bool {
	v := make([]bool, len(directions))
	for i, d := range directions {
		v[i] = obstructed(ds, point, d, 40000, 4500, nextDistance)
	}
	return v
}

func nextDistance(x float64) float64 {
	return 1.08*x + 1.5
}
