package sample_edges

import (
	"context"
	"flag"
	"log/slog"
	"math"

	"github.com/rp2709/gotopo/gotopo"
	"github.com/rp2709/gotopo/schema"
	"github.com/twpayne/go-geom"
)

func Run(dbConfig gotopo.DatabaseConfig, args []string) error {

	// sample edges flags
	var sampleEdgesConfig SampleEdgesConfig
	sampleEdgesFlags := flag.NewFlagSet("sample-edges", flag.ExitOnError)
	sampleEdgesFlags.Float64Var(&sampleEdgesConfig.SampleInterval, "sample-interval", 20.0, "distance between every sample point in meters")

	err := sampleEdgesFlags.Parse(args)
	if err != nil {
		return err
	}

	ctx := context.Background()
	queries, conn, err := dbConfig.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	slog.Info("Loading edges")
	// read all edges from DB, takes time
	var edges []gotopo.GetAllEdgesRow
	edges, err = queries.GetAllEdges(ctx)
	if err != nil {
		return err
	}

	var samplePoints []gotopo.InsertSamplePointParams

	slog.Info("Placing sample points")
	// run along edges and place sample points every sampleInterval
	for _, edge := range edges {
		distanceSinceLastSample := sampleEdgesConfig.SampleInterval / 2 // With a vertex between edges and an interval of 4: s--x--s and not s----x----s gives a distance of 4 units between s
		coords := edge.Shape.Coords()
		for i := range len(coords) - 1 {
			dx := coords[i+1].X() - coords[i].X()
			dy := coords[i+1].Y() - coords[i].Y()
			segmentLength := math.Sqrt(dx*dx + dy*dy)
			distanceSinceLastSample += segmentLength

			for distanceSinceLastSample >= sampleEdgesConfig.SampleInterval {
				distanceSinceLastSample -= sampleEdgesConfig.SampleInterval
				distanceFromPointToSample := segmentLength - distanceSinceLastSample
				spX := coords[i].X() + dx/segmentLength*distanceFromPointToSample
				spY := coords[i].Y() + dy/segmentLength*distanceFromPointToSample

				samplePoints = append(samplePoints, gotopo.InsertSamplePointParams{
					Eid:      edge.Eid,
					Location: geom.NewPoint(geom.XY).SetSRID(edge.Shape.SRID()).MustSetCoords(geom.Coord{spX, spY}),
				})
			}
		}
	}

	slog.Info("Writing schema to db")
	_, err = conn.Exec(ctx, schema.SamplePointsSchema())
	if err != nil {
		return err
	}

	slog.Info("Writing sample points to db", "number", len(samplePoints))
	_, err = queries.InsertSamplePoint(ctx, samplePoints)
	if err != nil {
		return err
	}

	slog.Info("Done")
	return nil
}

type SampleEdgesConfig struct {
	SampleInterval float64
}
