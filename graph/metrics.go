package graph

import (
	"time"

	"github.com/rp2709/gotopo/gotopo"
	"github.com/rp2709/gotopo/schema"
)

var roadTypeScores = map[schema.ObjectType]float64{
	schema.OBJECT_TYPE_SORTIE_AUTOROUTE:            1.0,
	schema.OBJECT_TYPE_ENTREE_AUTOROUTE:            1.0,
	schema.OBJECT_TYPE_AUTOROUTE:                   1.0,
	schema.OBJECT_TYPE_AIRE_DE_REPOS:               1.0,
	schema.OBJECT_TYPE_LIAISON:                     1.0,
	schema.OBJECT_TYPE_ACCES:                       1.0,
	schema.OBJECT_TYPE_ACCES_DE_SERVICE:            1.0,
	schema.OBJECT_TYPE_ROUTE_10M:                   0.6,
	schema.OBJECT_TYPE_ROUTE_8M:                    0.5,
	schema.OBJECT_TYPE_ROUTE_6M:                    0.4,
	schema.OBJECT_TYPE_ROUTE_4M:                    0.3,
	schema.OBJECT_TYPE_ROUTE_3M:                    0.05,
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

func distanceMetric(e gotopo.EdgeAgg, _ time.Time) float64 {
	return 1.0
}

func roadTypeMetric(e gotopo.EdgeAgg, _ time.Time) float64 {
	return roadTypeScores[e.RoadType]
}

func surfaceTypeMetric(e gotopo.EdgeAgg, _ time.Time) float64 {
	return surfaceTypeScores[e.SurfaceType]
}

var DefaultMetrics = Metrics{
	Metric{
		name:     "distance",
		evaluate: distanceMetric,
		weight:   0.5,
	},
	Metric{
		name:     "road type",
		evaluate: roadTypeMetric,
		weight:   0.25,
	},
	Metric{
		name:     "surface type",
		evaluate: surfaceTypeMetric,
		weight:   0.25,
	},
}

func DefaultEdgeDuration(e gotopo.EdgeAgg) time.Duration {
	return time.Duration(float64(time.Hour.Nanoseconds()) * (e.Length / 15000.0))
}
