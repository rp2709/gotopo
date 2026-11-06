package viewsheds_generation

import (
	"context"
	"log/slog"
	"math"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rp2709/gotopo/gotopo"
	"github.com/rp2709/gotopo/vrt_reader"
	"github.com/twpayne/go-geom"
	pgxgeom "github.com/twpayne/pgx-geom"
)

func Run(config GenerateViewshedsConfig) error {
	slog.Info("Connecting to database")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, config.Db.ConnectionString())
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	err = pgxgeom.Register(ctx, conn)
	if err != nil {
		return err
	}

	queries := gotopo.New(conn)

	if config.CreateWorkAreas {
		return createWorkAreas(config, queries, ctx)
	}
	return computeViewsheds(config, queries, ctx)
}

func computeViewsheds(config GenerateViewshedsConfig, queries *gotopo.Queries, ctx context.Context) error {
	slog.Info("Opening VRT raster dataset")
	ds, err := vrt_reader.OpenDataset(config.VrtFilePath, 0)
	if err != nil {
		return err
	}

	directions, err := computeDirectionsVectors(0.0, 0.0, 16)
	if err != nil {
		return err
	}

	for workArea, err := queries.GetNextWorkArea(ctx); err == nil; workArea, err = queries.GetNextWorkArea(ctx) {
		points, err := queries.GetSamplePointsWithin(ctx, workArea.Shape)
		if err != nil {
			return err
		}
		bounds := workArea.Shape.Bounds()

		ds.LoadArea(
			bounds.Min(0)-config.MaximumDistance,
			bounds.Max(0)+config.MaximumDistance,
			bounds.Min(1)-config.MaximumDistance,
			bounds.Min(1)+config.MaximumDistance,
		)

		slog.Info("Calculating viewsheds", "samples", len(points))
		for _, point := range points {
			_ = viewshed(ds, *point.Location, directions)
		}

		slog.Info("Marking area as done", "Area", workArea.Num)
		err = queries.UpdateWorkArea(ctx, gotopo.UpdateWorkAreaParams{
			Done: pgtype.Bool{Bool: true, Valid: true},
			Num:  workArea.Num,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func createWorkAreas(config GenerateViewshedsConfig, queries *gotopo.Queries, ctx context.Context) error {
	samplePointsExtent, err := queries.SamplePointsExtent(ctx)
	if err != nil {
		return err
	}

	samplePointsBounds := samplePointsExtent.Bounds()
	dimX := samplePointsBounds.Max(0) - samplePointsBounds.Min(0)
	dimY := samplePointsBounds.Max(1) - samplePointsBounds.Min(1)

	countX := math.Round(dimX/config.MaximumWorkAreaSize + 0.5) // + 0.5 to round up
	countY := math.Round(dimY/config.MaximumWorkAreaSize + 0.5)

	workAreaDimX := dimX / countX
	workAreaDimY := dimY / countY

	println("countX:", countX)
	println("countY:", countY)
	println("workAreaDimX:", workAreaDimX)
	println("workAreaDimY:", workAreaDimY)
	println("total working areas:", countX*countY)

	areas := make([]gotopo.InsertWorkAreasParams, 0, int(countX*countY))
	var totalSamples uint64 = 0
	for i := range int(countX) {
		minX, maxX := samplePointsBounds.Min(0)+float64(i)*workAreaDimX, samplePointsBounds.Min(0)+float64(i+1)*workAreaDimX
		for j := range int(countY) {
			minY, maxY := samplePointsBounds.Min(1)+float64(j)*workAreaDimY, samplePointsBounds.Min(1)+float64(j+1)*workAreaDimY
			workArea := geom.NewPolygon(geom.XY).SetSRID(samplePointsExtent.SRID()).MustSetCoords([][]geom.Coord{
				{{minX, maxY}, {maxX, maxY}, {maxX, minY}, {minX, minY}, {minX, maxY}},
			})
			points, err := queries.GetSamplePointsWithin(ctx, &workArea)
			if err != nil {
				return err
			}
			if len(points) == 0 {
				continue
			}
			totalSamples += uint64(len(points))
			areas = append(areas, gotopo.InsertWorkAreasParams{
				Num:   int32(len(areas)),
				Shape: workArea,
			})
		}
	}
	println("Areas with sample points: ", len(areas))
	println("Average samples per area: ", totalSamples/uint64(len(areas)))

	if err = queries.ClearWorkAreas(ctx); err != nil {
		return err
	}

	if _, err = queries.InsertWorkAreas(ctx, areas); err != nil {
		return err
	}

	return nil
}

func viewshed(ds *vrt_reader.RasterDataset, point geom.Point, directions []vector3) []bool {
	v := make([]bool, len(directions))
	for i, d := range directions {
		v[i] = obstructed(ds, point, d, 40000, 4500, NextDistance)
	}
	return v
}

type GenerateViewshedsConfig struct {
	Db *gotopo.DatabaseConfig

	MaximumWorkAreaSize float64
	MaximumDistance     float64
	VrtFilePath         string
	CreateWorkAreas     bool
}

func NextDistance(x float64) float64 {
	return 1.08*x + 1.5
}
