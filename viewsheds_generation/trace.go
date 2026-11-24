package viewsheds_generation

import (
	"math"

	"github.com/rp2709/gotopo/vrt_reader"
	"github.com/twpayne/go-geom"
)

func obstructed(dsm *vrt_reader.RasterDataset, dem *vrt_reader.RasterDataset, source geom.Point, direction vector3, maxDistance float64, maxAltitude float64, nextDistance func(float64) float64) bool {
	// take sample point altitude from terrain, not surface, since in a dense forest, the surface model would give the altitude of the canopy, not the ground
	sourceAltitude := float64(dem.ValueAt(source.X(), source.Y()))

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
		localAltitude := float64(dsm.ValueAt(source.X()+distanceAlongDirection*direction.x, source.Y()+distanceAlongDirection*direction.y))
		if stepAltitude <= localAltitude {
			return true // found some piece of terrain higher than our sight
		}
	}
	return false
}

func viewshed(dsm *vrt_reader.RasterDataset, dem *vrt_reader.RasterDataset, point geom.Point, maxDistance float64, maxAltitude float64, method PackingMethod) sliceRepresentation {
	v := make(sliceRepresentation, len(method.Directions()))
	for i, d := range method.Directions() {
		v[i] = obstructed(dsm, dem, point, d, maxDistance, maxAltitude, nextDistance)
	}
	return v
}

func nextDistance(x float64) float64 {
	return 1.08*x + 1.5
}
