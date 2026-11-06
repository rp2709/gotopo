package viewsheds_generation

import (
	"context"
	"log/slog"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rp2709/gotopo/gotopo"
	"github.com/rp2709/gotopo/vrt_reader"
	"github.com/twpayne/go-geom"
	pgxgeom "github.com/twpayne/pgx-geom"
)

func Run(config GenerateViewshedsConfig) error {
	previousLogLevel := slog.SetLogLoggerLevel(config.LogLevel)
	defer slog.SetLogLoggerLevel(previousLogLevel)

	slog.Info("Connecting to database")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, config.ConnectionString())
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	err = pgxgeom.Register(ctx, conn)
	if err != nil {
		return err
	}

	queries := gotopo.New(conn)

	slog.Info("Opening VRT raster dataset")
	ds, err := vrt_reader.OpenDataset(config.VrtFilePath, 0)
	if err != nil {
		return err
	}

	directions, err := computeDirectionsVectors(0.0, 0.0, 16)
	if err != nil {
		return err
	}

	maximumWorkAreaSize := 10000.0
	maxDistance := 40000.0

	samplePointsExtent, err := queries.SamplePointsExtent(ctx)
	if err != nil {
		return err
	}

	samplePointsBounds := samplePointsExtent.Bounds()
	dimX := samplePointsBounds.Max(0) - samplePointsBounds.Min(0)
	dimY := samplePointsBounds.Max(1) - samplePointsBounds.Min(1)

	countX := math.Round(dimX/maximumWorkAreaSize + 0.5)
	countY := math.Round(dimY/maximumWorkAreaSize + 0.5)

	workAreaDimX := dimX / countX
	workAreaDimY := dimY / countY

	println("countX:", countX)
	println("countY:", countY)
	println("workAreaDimX:", workAreaDimX)
	println("workAreaDimY:", workAreaDimY)
	println("total working areas:", countX*countY)

	var areasWithSamples uint64 = 0
	var totalSamples uint64 = 0
	ti := time.Now()
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
			areasWithSamples++
			totalSamples += uint64(len(points))
			slog.Info("Loading area")
			ds.LoadArea(
				minX-maxDistance,
				maxX+maxDistance,
				minY-maxDistance,
				maxY+maxDistance,
			)
			slog.Info("Calculating viewsheds", "samples", len(points))
			for _, point := range points {
				_ = viewshed(ds, *point.Location, directions)
			}
		}
		break
	}
	tf := time.Now()
	println("Areas with sample points: ", areasWithSamples)
	println("Average samples per area: ", totalSamples/areasWithSamples)
	println("Average time per sample point: ",
		(time.Duration(tf.Sub(ti).Nanoseconds()/int64(totalSamples)) * time.Nanosecond).String())

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
	LogLevel slog.Level
	gotopo.DatabaseConfig

	VrtFilePath string
}

func NextDistance(x float64) float64 {
	return 1.08*x + 1.5
}
