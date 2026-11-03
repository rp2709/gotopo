package viewsheds_generation

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/rp2709/gotopo/vrt_reader"
	"github.com/twpayne/go-geom"
	pgxgeom "github.com/twpayne/pgx-geom"
)

const connectionStringFmt = "user=%s dbname=%s sslmode=prefer host=%s password=%s port=%d"

func Run(config GenerateViewshedsConfig) error {
	previousLogLevel := slog.SetLogLoggerLevel(config.LogLevel)
	defer slog.SetLogLoggerLevel(previousLogLevel)

	slog.Info("Connecting to database")
	ctx := context.Background()
	conn, err := pgx.Connect(
		ctx, fmt.Sprintf(
			connectionStringFmt,
			config.Dbuser,
			config.Dbname,
			config.Dbhost,
			config.Dbpwd,
			config.Dbport,
		),
	)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	err = pgxgeom.Register(ctx, conn)
	if err != nil {
		return err
	}

	//queries := gotopo.New(conn)

	slog.Info("Opening VRT raster dataset")
	ds, err := vrt_reader.OpenDataset(config.VrtFilePath, 0)
	if err != nil {
		return err
	}

	directions, err := computeDirectionsVectors(0.0, 0.0, 16)
	if err != nil {
		return err
	}
	/*
		const limit = 5000
		slog.Info("Loading sample points from db", "limit", limit)
		samplePoints, err := queries.GetSomeSamplePoints(ctx, limit)
		if err != nil {
			return err
		}

		slog.Info("Start computing hemispherical viewsheds")
		ti := time.Now()
		for _, point := range samplePoints {
			_ = viewshed(ds, *point.Location, directions)
		}
		tf := time.Now()
		duration := tf.Sub(ti)
		slog.Info("Finished computing hemispherical viewsheds", "duration", duration.String())
	*/
	err = SaveViewshedPNG(
		"viewshed.png",
		viewshed(ds, *geom.NewPoint(geom.XY).SetSRID(2056).MustSetCoords(geom.Coord{2531882.927, 1184492.072}), directions),
		500,
		0.0,
		0.0,
	)
	if err != nil {
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
	LogLevel slog.Level
	Dbuser   string
	Dbname   string
	Dbpwd    string
	Dbhost   string
	Dbport   int

	VrtFilePath string
}

func NextDistance(x float64) float64 {
	return 1.08*x + 1.5
}
