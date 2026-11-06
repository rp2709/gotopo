package vrt_reader

import (
	"log/slog"
	"math"
	"os"
	"runtime"

	"github.com/rp2709/geotiff/geotiff"
)

type loadedImage struct {
	img      *geotiff.GeoTIFF
	filename string
	err      error
}

func loadWorker(toLoad chan string, toAddToCache chan loadedImage) {
	for filename := range toLoad {
		file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
		if err != nil {
			toAddToCache <- loadedImage{nil, filename, err}
		}

		img, err := geotiff.Read(file)
		_ = file.Close()
		if err != nil {
			toAddToCache <- loadedImage{nil, filename, err}
		}

		toAddToCache <- loadedImage{
			img:      img,
			filename: filename,
			err:      nil,
		}
	}
}

func (ds *RasterDataset) LoadArea(minX, maxX, minY, maxY float64) {
	tlX, tlY := ds.tile(minX, maxY)
	brX, brY := ds.tile(maxX, minY)

	toLoad := make(chan string, (brX-tlX)*(brY-tlY))
	toAddToCache := make(chan loadedImage, (brX-tlX)*(brY-tlY))

	for len(ds.cache.cache) > ds.cache.maxSize {
		var minTime uint64 = math.MaxUint64
		minP := ""
		for p, entry := range ds.cache.cache {
			if entry.lastAccessed < minTime {
				minP = p
				minTime = entry.lastAccessed
			}
		}
		delete(ds.cache.cache, minP)
	}

	fileCount := 0
	for i := tlX; i <= brX; i++ {
		for j := tlY; j <= brY; j++ {
			if i < 0 || j < 0 || i >= len(ds.band) || j >= len(ds.band[i]) {
				continue
			}
			file := ds.band[i][j]
			if file != "" && ds.cache.cache[file] == nil {
				toLoad <- file
				fileCount++
			}
		}
	}
	close(toLoad)

	for range runtime.NumCPU() {
		go loadWorker(toLoad, toAddToCache)
	}

	for range fileCount {
		file := <-toAddToCache
		if file.err != nil {
			slog.Warn("Error opening a tile file", "File", file.filename, "Error", file.err)
			continue
		}
		ds.cache.cache[file.filename] = &cacheEntry{image: file.img}
	}
}
