package viewsheds_generation

import "github.com/rp2709/gotopo/vrt_reader"

func Run(config GenerateViewshedsConfig) error {
	ds, err := vrt_reader.OpenDataset(config.VrtFilePath, 0)
	if err != nil {
		return err
	}

	print(ds.ValueAt(2540218, 1181114))

	return nil
}

type GenerateViewshedsConfig struct {
	VrtFilePath string
}
