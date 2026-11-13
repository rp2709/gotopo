package viewsheds_generation

import (
	"context"
	"flag"
	"log/slog"
	"math"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rp2709/gotopo/gotopo"
	"github.com/rp2709/gotopo/vrt_reader"
	"github.com/twpayne/go-geom"
)

func doVisualizeViewsheds(_ gotopo.DatabaseConfig, viewshedConfig viewshedsConfigT, args []string) error {
	flags := flag.NewFlagSet("visualize", flag.ExitOnError)
	x := flags.Float64("X", 2531882.927, "x coordinate of test point")
	y := flags.Float64("Y", 1184492.072, "y coordinate of test point")
	err := flags.Parse(args)
	if err != nil {
		return err
	}

	ds, err := vrt_reader.OpenDataset(viewshedConfig.VrtFilePath, 0)
	if err != nil {
		return err
	}

	directions, err := computeDirectionsVectors(minimumAngleFromHorizon, minimumAngleFromZenith, viewshedSizeBits)
	if err != nil {
		return err
	}

	point := *geom.NewPoint(geom.XY).SetSRID(2056).MustSetCoords(geom.Coord{*x, *y})

	ds.LoadArea(
		point.X()-viewshedConfig.MaxDistance,
		point.X()+viewshedConfig.MaxDistance,
		point.Y()-viewshedConfig.MaxDistance,
		point.Y()+viewshedConfig.MaxDistance,
	)

	vsb := viewshed(ds, point, directions)
	return VisualizeViewshed(
		"viewshed.png",
		SliceOfBoolToBits(vsb),
		2000,
		minimumAngleFromHorizon,
		minimumAngleFromZenith,
	)
}

func doPrintDirections(_ gotopo.DatabaseConfig, _ viewshedsConfigT, _ []string) error {
	directions, err := computeDirectionsVectors(minimumAngleFromHorizon, minimumAngleFromZenith, viewshedSizeBits)
	if err != nil {
		return err
	}
	return writeDirections(os.Stdout, directions)
}

func doGenerateViewsheds(db gotopo.DatabaseConfig, viewshedConfig viewshedsConfigT, _ []string) error {
	ctx := context.Background()
	queries, conn, err := db.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	ds, err := vrt_reader.OpenDataset(viewshedConfig.VrtFilePath, 0)
	if err != nil {
		return err
	}

	directions, err := computeDirectionsVectors(minimumAngleFromHorizon, minimumAngleFromZenith, viewshedSizeBits)
	if err != nil {
		return err
	}

	var workArea gotopo.WorkArea
	for workArea, err = queries.GetNextWorkArea(ctx); err == nil; workArea, err = queries.GetNextWorkArea(ctx) {
		points, err := queries.GetSamplePointsWithin(ctx, workArea.Shape)
		if err != nil {
			return err
		}
		bounds := workArea.Shape.Bounds()

		ds.LoadArea(
			bounds.Min(0)-viewshedConfig.MaxDistance,
			bounds.Max(0)+viewshedConfig.MaxDistance,
			bounds.Min(1)-viewshedConfig.MaxDistance,
			bounds.Min(1)+viewshedConfig.MaxDistance,
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
	return err
}

func doCreateWorkAreas(db gotopo.DatabaseConfig, _ viewshedsConfigT, args []string) error {
	createWorkAreaFlags := flag.NewFlagSet("create-work-areas", flag.ExitOnError)
	maximumWorkAreaSize := createWorkAreaFlags.Float64("max-size", 10000, "maximum size of work areas")
	err := createWorkAreaFlags.Parse(args)
	if err != nil {
		return err
	}

	ctx := context.Background()
	queries, conn, err := db.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	samplePointsExtent, err := queries.SamplePointsExtent(ctx)
	if err != nil {
		return err
	}

	samplePointsBounds := samplePointsExtent.Bounds()
	dimX := samplePointsBounds.Max(0) - samplePointsBounds.Min(0)
	dimY := samplePointsBounds.Max(1) - samplePointsBounds.Min(1)

	countX := math.Round(dimX / *maximumWorkAreaSize + 0.5) // + 0.5 to round up
	countY := math.Round(dimY / *maximumWorkAreaSize + 0.5)

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
