package viewsheds_generation

import (
	"time"

	"github.com/rp2709/gotopo/vrt_reader"
	"github.com/twpayne/go-geom"
)

func Run(config GenerateViewshedsConfig) error {
	ds, err := vrt_reader.OpenDataset(config.VrtFilePath, 0)
	if err != nil {
		return err
	}

	directions, err := computeDirectionsVectors(0.0, 0.0, 8)
	if err != nil {
		return err
	}

	ti := time.Now()
	viewshed := make([]bool, len(directions))
	for i, d := range directions {
		viewshed[i] = obstructed(ds, *geom.NewPoint(geom.XY).SetSRID(2056).MustSetCoords(geom.Coord{2540218, 1181114}), d, 30000, 4500, LinearStepDistance(5, 0))
	}
	tf := time.Now()
	duration := tf.Sub(ti)
	println(duration.String())

	return nil
}

type GenerateViewshedsConfig struct {
	VrtFilePath string
}
