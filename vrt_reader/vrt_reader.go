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
	band      [][]string
	tileSizeX int
	tileSizeY int
	nodata    float64
	cache     imageCache
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

	ds.cache = imageCache{
		maxSize: 50,
		time:    0,
		cache:   make(map[string]*cacheEntry),
	}

	return &ds, nil
}

func (ds *RasterDataset) ValueAt(x, y float64) float32 {
	X, Y := ds.gt.applyTransform(x, y)

	// Check raster bounds before calculating tile indices.
	if math.IsNaN(X) || math.IsNaN(Y) ||
		math.IsInf(X, 0) || math.IsInf(Y, 0) ||
		X < 0 || Y < 0 ||
		X >= float64(ds.dataset.RasterXSize) ||
		Y >= float64(ds.dataset.RasterYSize) {
		return float32(ds.nodata)
	}

	pixelX := int(math.Floor(X))
	pixelY := int(math.Floor(Y))

	tileX := pixelX / ds.tileSizeX
	tileY := pixelY / ds.tileSizeY

	if tileX >= len(ds.band) || tileY >= len(ds.band[tileX]) {
		return float32(ds.nodata)
	}

	file := ds.band[tileX][tileY]
	if file == "" {
		return float32(ds.nodata)
	}

	value, err := ds.readValue(file, x, y)
	if err != nil {
		slog.Warn("Error reading value",
			"x", x,
			"y", y,
			"err", err,
		)
		return float32(ds.nodata)
	}

	return value
}

type cacheEntry struct {
	image        *geotiff.GeoTIFF
	lastAccessed uint64
}
type imageCache struct {
	maxSize int
	time    uint64
	cache   map[string]*cacheEntry
}

func (cache *imageCache) get(path string) (*geotiff.GeoTIFF, error) {
	cache.time++
	if img, ok := cache.cache[path]; ok {
		img.lastAccessed = cache.time
		return img.image, nil
	}

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := geotiff.Read(file)
	if err != nil {
		return nil, err
	}

	cache.cache[path] = &cacheEntry{img, cache.time}

	if len(cache.cache) > cache.maxSize {
		var minTime uint64 = math.MaxUint64
		minP := ""
		for p, entry := range cache.cache {
			if entry.lastAccessed < minTime {
				minP = p
				minTime = entry.lastAccessed
			}
		}
		delete(cache.cache, minP)
	}

	return img, nil
}

func (ds *RasterDataset) readValue(path string, x float64, y float64) (float32, error) {
	img, err := ds.cache.get(path)
	if err != nil {
		return 0, err
	}

	value, err := img.AtCoord(x, y, false)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (ds *RasterDataset) bandMatrix(band int, vrtFilePath string) [][]string {
	dimX := (ds.dataset.RasterXSize + ds.tileSizeX - 1) / ds.tileSizeX
	dimY := (ds.dataset.RasterYSize + ds.tileSizeY - 1) / ds.tileSizeY

	matrix := make([][]string, dimX)
	for i := range matrix {
		matrix[i] = make([]string, dimY)
	}

	vrtDirectory := filepath.Dir(vrtFilePath)

	for _, tile := range ds.dataset.RasterBands[band].Sources {
		x := int(tile.DstRect.XOff) / ds.tileSizeX
		y := int(tile.DstRect.YOff) / ds.tileSizeY

		path := tile.SourceFilename.Value

		if tile.SourceFilename.RelativeToVRT {
			path = filepath.Join(vrtDirectory, path)
		}

		matrix[x][y] = path
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
