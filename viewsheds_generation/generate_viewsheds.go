package viewsheds_generation

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rp2709/gotopo/gotopo"
	"github.com/rp2709/gotopo/vrt_reader"
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

	queries := gotopo.New(conn)

	slog.Info("Opening VRT raster dataset")
	ds, err := vrt_reader.OpenDataset(config.VrtFilePath, 0)
	if err != nil {
		return err
	}

	directions, err := computeDirectionsVectors(0.0, 0.0, 8)
	if err != nil {
		return err
	}

	const limit = 2000
	slog.Info("Loading sample points from db", "limit", limit)
	samplePoints, err := queries.GetSomeSamplePoints(ctx, limit)
	if err != nil {
		return err
	}

	slog.Info("Start computing hemispherical viewsheds")
	ti := time.Now()
	for _, point := range samplePoints {
		viewshed := make([]bool, len(directions))
		for i, d := range directions {
			viewshed[i] = obstructed(ds, *point.Location, d, 5000, 4500, 2)
		}

	}
	tf := time.Now()
	duration := tf.Sub(ti)
	println(duration.String())

	return nil
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
