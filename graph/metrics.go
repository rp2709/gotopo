package graph

import (
	"math"
	"time"

	"github.com/rp2709/gotopo/gotopo"
	"github.com/rp2709/gotopo/schema"
	"github.com/rp2709/gotopo/sg2"
	"github.com/rp2709/gotopo/viewsheds_generation"
	"github.com/twpayne/go-geom"
)

var roadTypeScores = map[schema.ObjectType]float64{
	schema.OBJECT_TYPE_SORTIE_AUTOROUTE:            1.0,
	schema.OBJECT_TYPE_ENTREE_AUTOROUTE:            1.0,
	schema.OBJECT_TYPE_AUTOROUTE:                   1.0,
	schema.OBJECT_TYPE_AIRE_DE_REPOS:               1.0,
	schema.OBJECT_TYPE_LIAISON:                     1.0,
	schema.OBJECT_TYPE_ACCES:                       1.0,
	schema.OBJECT_TYPE_ACCES_DE_SERVICE:            1.0,
	schema.OBJECT_TYPE_ROUTE_10M:                   0.8,
	schema.OBJECT_TYPE_ROUTE_8M:                    0.7,
	schema.OBJECT_TYPE_ROUTE_6M:                    0.6,
	schema.OBJECT_TYPE_ROUTE_4M:                    0.2,
	schema.OBJECT_TYPE_ROUTE_3M:                    0.0,
	schema.OBJECT_TYPE_PLACE:                       0.0,
	schema.OBJECT_TYPE_TRAIN_AUTOBUS:               1.0,
	schema.OBJECT_TYPE_BAC:                         0.5,
	schema.OBJECT_TYPE_CHEMIN_2M:                   0.05,
	schema.OBJECT_TYPE_CHEMIN_1M:                   0.1,
	schema.OBJECT_TYPE_FRAGMENT_CHEMIN_1M:          0.0,
	schema.OBJECT_TYPE_FRAGMENT_CHEMIN_2M:          0.0,
	schema.OBJECT_TYPE_MARQUE_TRACEE:               0.5,
	schema.OBJECT_TYPE_SEMI_AUTOROUTE:              1.0,
	schema.OBJECT_TYPE_VIA_FERRATA:                 1.0,
	schema.OBJECT_TYPE_MOBILITE_DOUCE_NON_ATTRIBUE: 0.0,
}

var surfaceTypeScores = map[schema.SurfaceType]float64{
	schema.SURFACE_DURE:          0.0,
	schema.SURFACE_NATURE:        0.5,
	schema.SURFACE_INCONNU:       0.4,
	schema.SURFACE_PAS_DE_VALEUR: 0.4,
}

func sunExposureMetric(e gotopo.EdgeAgg, t time.Time) float64 {
	if len(e.Viewsheds) == 0 {
		return 0
	}

	edgeCenter := geom.NewPoint(geom.XY).SetSRID(e.Shape.SRID()).MustSetCoords(e.Shape.Coord(e.Shape.NumCoords() / 2))
	azimuth, altitude := sg2.SunPosition(*edgeCenter, t)
	packingMethod := viewsheds_generation.DefaultEvenPackingMethod

	score := 0.0
	validNb := 0
	for _, vs := range e.Viewsheds {
		if !vs.Valid {
			continue
		}
		validNb++
		score += 1.0 - packingMethod.Interpolate(vs, azimuth, altitude)
	}

	if validNb == 0 {
		return 0
	}

	return score / float64(validNb)
}

func distanceMetric(e gotopo.EdgeAgg, _ time.Time) float64 {
	return 1.0
}

func roadTypeMetric(e gotopo.EdgeAgg, _ time.Time) float64 {
	return roadTypeScores[e.RoadType]
}

func surfaceTypeMetric(e gotopo.EdgeAgg, _ time.Time) float64 {
	return surfaceTypeScores[e.SurfaceType]
}

func slopeMetric(e gotopo.EdgeAgg, _ time.Time) float64 {
	slopeScore := 0.0
	for i := 0; i < e.Shape.NumCoords()-1; i++ {
		currentAlti := e.Shape.FlatCoords()[i*e.Shape.Stride()+2]
		nextAlti := e.Shape.FlatCoords()[(i+1)*e.Shape.Stride()+2]
		dx, dy := e.Shape.Coords()[i+1].X()-e.Shape.Coords()[i].X(), e.Shape.Coords()[i+1].Y()-e.Shape.Coords()[i].Y()
		slope := (nextAlti - currentAlti) / math.Sqrt(dx*dx+dy*dy)
		if slope > 0 {
			slopeScore += slope
		}
	}
	// I would like to map this slope percentage to values where 1.0 is the steepest slope there is
	// and 0.0 flat or descending. I am assuming that there is no road with a slope greater than 30% and therefore multiply the slope by 3
	return 3 * slopeScore
}

var DefaultMetrics = Metrics{
	Metric{
		name:     "distance",
		evaluate: distanceMetric,
		weight:   0.20,
	},
	Metric{
		name:     "road type",
		evaluate: roadTypeMetric,
		weight:   0.20,
	},
	Metric{
		name:     "surface type",
		evaluate: surfaceTypeMetric,
		weight:   0.20,
	},
	Metric{
		name:     "slope",
		evaluate: slopeMetric,
		weight:   0.20,
	},
	Metric{
		name:     "sun exposure",
		evaluate: sunExposureMetric,
		weight:   0.20,
	},
}

func DefaultEdgeDuration(e gotopo.EdgeAgg) time.Duration {
	return time.Duration(float64(time.Hour.Nanoseconds()) * (e.Length / 15000.0))
}
