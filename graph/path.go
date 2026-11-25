package graph

import (
	"context"
	"log/slog"
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
	slog.Info("found starting vertex", "vertex_id", startVertex)

	endVertex, err := queries.GetClosestVertexWithin(ctx, gotopo.GetClosestVertexWithinParams{
		Point:       endPoint,
		MaxDistance: 50.0,
	})
	if err != nil {
		return err
	}
	slog.Info("found ending vertex", "vertex_id", endVertex)

	margin := 5000.0

	minX := min(startPoint.X(), endPoint.X()) - margin
	minY := min(startPoint.Y(), endPoint.Y()) - margin
	maxX := max(startPoint.X(), endPoint.X()) + margin
	maxY := max(startPoint.Y(), endPoint.Y()) + margin

	pathArea := geom.NewPolygon(geom.XY).SetSRID(2056).MustSetCoords([][]geom.Coord{{
		{minX, minY},
		{maxX, minY},
		{maxX, maxY},
		{minX, maxY},
		{minX, minY}, // ring must be closed
	}})

	slog.Debug("Start loading graph")
	rows, err := queries.GetGraphWithin(ctx, pathArea)
	if err != nil {
		return err
	}
	slog.Info("loaded graph", "rows", len(rows))

	g := FromRows(rows)
	slog.Info("starting to compute shortest path")
	path, err := ShortestPath(g, startVertex, endVertex, time.Now(), DefaultMetrics, DefaultEdgeDuration)
	if err != nil {
		return err
	}

	slog.Info("shortest path found", "length", len(path))
	pathMultiLineString := geom.NewMultiLineString(geom.XYZ).SetSRID(2056)
	for _, e := range path {
		err = pathMultiLineString.Push(e.Shape)
		if err != nil {
			return err
		}
	}

	err = queries.AddPath(ctx, pathMultiLineString)
	if err != nil {
		return err
	}
	return nil
}
