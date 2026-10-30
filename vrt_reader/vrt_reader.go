package vrt_reader

import "os"

type RasterDataset struct {
	srid uint32 // in what crs the coordinates passed to ValueAt are. Initialized with value read from vrt file

}

func OpenDataset(vrtFilePath string) (*RasterDataset, error) {
	// read vrt file itself and verify structure
	file, err := os.OpenFile(vrtFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// initialize array of tiff image paths

	return &RasterDataset{}, nil
}

func (ds *RasterDataset) ValueAt(x float64, y float64) float64 {
	return 0
}

func (ds *RasterDataset) GetSRID() uint32 {
	return ds.srid
}

func (ds *RasterDataset) SetSRID(x uint32) {
	// maybe recalculate entire index to avoid transforming between crs
}
