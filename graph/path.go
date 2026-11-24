package graph

import (
	"context"
	"time"

	"github.com/rp2709/gotopo/gotopo"
	"github.com/twpayne/go-geom"
)

func Run(db gotopo.DatabaseConfig, _ []string) error {
	ctx := context.Background()
	queries, conn, err := db.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	startPoint := geom.NewPoint(geom.XY).SetSRID(2056).MustSetCoords(geom.Coord{2540426.90, 1187322.16})
	endPoint := geom.NewPoint(geom.XY).SetSRID(2056).MustSetCoords(geom.Coord{2540374.29, 1181158.43})

	startVertex, err := queries.GetClosestVertexWithin(ctx, gotopo.GetClosestVertexWithinParams{
		Point:       startPoint,
		MaxDistance: 50.0,
	})
	if err != nil {
		return err
	}

	endVertex, err := queries.GetClosestVertexWithin(ctx, gotopo.GetClosestVertexWithinParams{
		Point:       endPoint,
		MaxDistance: 50.0,
	})
	if err != nil {
		return err
	}

	margin := 100.0
	pathArea := geom.NewBounds(geom.XY).SetCoords(
		geom.Coord{min(startPoint.X(), endPoint.X()) - margin, min(startPoint.Y(), endPoint.Y()) - margin},
		geom.Coord{max(startPoint.X(), endPoint.X()) + margin, max(startPoint.Y(), endPoint.Y()) + margin},
	).Polygon().SetSRID(2056)

	graph, err := queries.GetGraphWithin(ctx, pathArea)
	if err != nil {
		return err
	}

	path, err := ShortestPath(graph, startVertex, endVertex, time.Now(), DefaultMetrics, DefaultEdgeDuration)
	if err != nil {
		return err
	}

	println(path)
	return nil
}
