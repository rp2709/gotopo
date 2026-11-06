package vrt_reader

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rp2709/geotiff/geotiff"
)

type RasterDataset struct {
	srid      uint32 // in what crs the coordinates passed to ValueAt are. Initialized with value read from vrt file
	dataset   VRTDataset
	gt        geoTransform
	band      [][]datasetTile
	tileSizeX int
	tileSizeY int
	nodata    float64
}

type datasetTile struct {
	file string
	img  *geotiff.GeoTIFF
}

func OpenDataset(vrtFilePath string, band int) (*RasterDataset, error) {
	ds := RasterDataset{}
	data, err := os.ReadFile(vrtFilePath)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &ds.dataset)
	if err != nil {
		return nil, err
	}

	if !ds.dataset.RasterBands[band].allSourcesSameDimensions() {
		return nil, fmt.Errorf("Band's sources are not all the same dimensions")
	}

	ds.nodata, err = strconv.ParseFloat(ds.dataset.RasterBands[band].NoDataValue, 64)
	if err != nil {
		return nil, err
	}

	ds.gt, err = parseGeoTransform(ds.dataset.GeoTransform)
	if err != nil {
		return nil, err
	}

	ds.tileSizeX = int(ds.dataset.RasterBands[band].Sources[0].SrcRect.XSize)
	ds.tileSizeY = int(ds.dataset.RasterBands[band].Sources[0].SrcRect.YSize)

	ds.band = ds.bandMatrix(band, vrtFilePath)

	return &ds, nil
}

func (ds *RasterDataset) ValueAt(x, y float64) float32 {

	tileX, tileY := ds.tile(x, y)
	if tileX < 0 || tileY < 0 || tileX >= len(ds.band) || tileY >= len(ds.band[tileX]) {
		slog.Debug("Tile outside of dataset", "tileX", tileX, "tileY", tileY)
		return float32(ds.nodata)
	}

	tile := ds.band[tileX][tileY]
	if tile.img == nil {
		if tile.file != "" {
			slog.Warn("Tile image is not loaded", "tileX", tileX, "tileY", tileY, "path", tile.file)
		}
		return float32(ds.nodata)
	}

	value, err := tile.img.AtCoord(x, y, false)
	if err != nil {
		// some tiles may be smaller at the edges, only print on debug
		if tileX == 0 || tileY == 0 || tileX == ds.tileSizeX || tileY == ds.tileSizeY {
			slog.Debug("Error reading value",
				"x", x,
				"y", y,
				"err", err,
			)
		} else {
			slog.Warn("Error reading value",
				"x", x,
				"y", y,
				"err", err,
			)
		}
		return float32(ds.nodata)
	}

	return value
}

func (ds *RasterDataset) bandMatrix(band int, vrtFilePath string) [][]datasetTile {
	dimX := (ds.dataset.RasterXSize + ds.tileSizeX - 1) / ds.tileSizeX
	dimY := (ds.dataset.RasterYSize + ds.tileSizeY - 1) / ds.tileSizeY

	matrix := make([][]datasetTile, dimX)
	for i := range matrix {
		matrix[i] = make([]datasetTile, dimY)
	}

	vrtDirectory := filepath.Dir(vrtFilePath)

	for _, tile := range ds.dataset.RasterBands[band].Sources {
		x := int(tile.DstRect.XOff) / ds.tileSizeX
		y := int(tile.DstRect.YOff) / ds.tileSizeY

		path := tile.SourceFilename.Value

		if tile.SourceFilename.RelativeToVRT {
			path = filepath.Join(vrtDirectory, path)
		}

		matrix[x][y] = datasetTile{file: path}
	}
	return matrix
}

func (band VRTRasterBand) allSourcesSameDimensions() bool {
	var srcRect *Rect
	for _, src := range band.Sources {
		if srcRect == nil {
			srcRect = &src.SrcRect
		} else if *srcRect != src.SrcRect {
			return false
		}
	}
	return true
}

func (ds *RasterDataset) tile(x, y float64) (int, int) {
	X, Y := ds.gt.applyTransform(x, y)

	pixelX := int(math.Floor(X))
	pixelY := int(math.Floor(Y))

	tileX := pixelX / ds.tileSizeX
	tileY := pixelY / ds.tileSizeY
	return tileX, tileY
}

func (gt geoTransform) applyTransform(x float64, y float64) (float64, float64) {
	det := gt[1]*gt[5] - gt[2]*gt[4]
	if det == 0 {
		panic("Invalid transform")
	}

	dx := x - gt[0]
	dy := y - gt[3]

	pixel := (gt[5]*dx - gt[2]*dy) / det
	line := (-gt[4]*dx + gt[1]*dy) / det

	return pixel, line
}

func parseGeoTransform(s string) (geoTransform, error) {
	var result geoTransform

	parts := strings.Split(s, ",")

	if len(parts) != 6 {
		return result, fmt.Errorf(
			"invalid GeoTransform: expected 6 values, got %d",
			len(parts),
		)
	}

	for i, part := range parts {
		value, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return result, err
		}
		result[i] = value
	}

	return result, nil
}

type geoTransform [6]float64
