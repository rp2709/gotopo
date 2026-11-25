package viewsheds_generation

import (
	"context"
	"flag"
	"log/slog"
	"math"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rp2709/gotopo/gotopo"
	"github.com/rp2709/gotopo/schema"
	"github.com/rp2709/gotopo/vrt_reader"
	"github.com/twpayne/go-geom"
)

func doVisualizeViewsheds(_ gotopo.DatabaseConfig, viewshedConfig viewshedsConfigT, args []string) error {
	flags := flag.NewFlagSet("visualize", flag.ExitOnError)
	x := flags.Float64("X", 2531882.927, "x coordinate of test point")
	y := flags.Float64("Y", 1184492.072, "y coordinate of test point")
	size := flags.Int("size", 800, "size of the image in pixels")
	err := flags.Parse(args)
	if err != nil {
		return err
	}

	dsm, err := vrt_reader.OpenDataset(viewshedConfig.DSMFilePath, 0)
	if err != nil {
		return err
	}

	dem, err := vrt_reader.OpenDataset(viewshedConfig.DEMFilePath, 0)
	if err != nil {
		return err
	}

	point := *geom.NewPoint(geom.XY).SetSRID(2056).MustSetCoords(geom.Coord{*x, *y})

	dsm.LoadArea(
		point.X()-viewshedConfig.MaxDistance,
		point.X()+viewshedConfig.MaxDistance,
		point.Y()-viewshedConfig.MaxDistance,
		point.Y()+viewshedConfig.MaxDistance,
	)

	dem.LoadArea(
		point.X(),
		point.X(),
		point.Y(),
		point.Y(),
	)

	vsb := viewshed(dsm, dem, point, viewshedConfig.MaxDistance, viewshedConfig.MaxAltitude, packingMethod)
	return VisualizeViewshed(
		"viewshed.png",
		vsb.toBits(),
		*size,
		packingMethod,
	)
}

func doPrintDirections(_ gotopo.DatabaseConfig, _ viewshedsConfigT, _ []string) error {
	return writeDirections(os.Stdout, packingMethod.Directions())
}

func doGenerateViewsheds(db gotopo.DatabaseConfig, viewshedConfig viewshedsConfigT, _ []string) error {
	ctx := context.Background()
	queries, conn, err := db.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	dsm, err := vrt_reader.OpenDataset(viewshedConfig.DSMFilePath, 0)
	if err != nil {
		return err
	}

	dem, err := vrt_reader.OpenDataset(viewshedConfig.DEMFilePath, 0)
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

		dsm.LoadArea(
			bounds.Min(0)-viewshedConfig.MaxDistance,
			bounds.Max(0)+viewshedConfig.MaxDistance,
			bounds.Min(1)-viewshedConfig.MaxDistance,
			bounds.Max(1)+viewshedConfig.MaxDistance,
		)

		dem.LoadArea(
			bounds.Min(0),
			bounds.Max(0),
			bounds.Min(1),
			bounds.Max(1),
		)

		slog.Info("Calculating viewsheds", "samples", len(points))
		for _, point := range points {
			point.Viewshed = viewshed(dsm, dem, *point.Location, viewshedConfig.MaxDistance, viewshedConfig.MaxAltitude, packingMethod).toBits()
		}

		slog.Info("Saving viewsheds", "samples", len(points))
		params := make([]gotopo.InsertSamplePointWithViewshedParams, len(points))
		for i, point := range points {
			params[i] = gotopo.InsertSamplePointWithViewshedParams{
				Eid:      point.Eid,
				Location: point.Location,
				Viewshed: point.Viewshed,
			}
		}

		_, err = queries.InsertSamplePointWithViewshed(ctx, params)
		if err != nil {
			return err
		}

		err = queries.DeleteSamplePointsWithoutViewshedsWithin(ctx, workArea.Shape)
		if err != nil {
			return err
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

	_, err = conn.Exec(ctx, schema.WorkAreasSchema())
	if err != nil {
		return err
	}

	if _, err = queries.InsertWorkAreas(ctx, areas); err != nil {
		return err
	}

	return nil
}
