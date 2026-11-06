package vrt_reader

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/rp2709/geotiff/geotiff"
)

type imageToLoad struct {
	datasetTile
	x, y int
}
type loadedImage struct {
	imageToLoad
	err error
}

func loadWorker(toLoad chan imageToLoad, toAddToCache chan loadedImage) {
	for e := range toLoad {
		file, err := os.OpenFile(e.file, os.O_RDONLY, os.ModePerm)
		if err != nil {
			toAddToCache <- loadedImage{err: err}
		}

		e.img, err = geotiff.Read(file)
		_ = file.Close()
		if err != nil {
			toAddToCache <- loadedImage{err: err}
		}

		toAddToCache <- loadedImage{
			imageToLoad: e,
			err:         nil,
		}
	}
}

func (ds *RasterDataset) LoadArea(minX, maxX, minY, maxY float64) {
	tlX, tlY := ds.tile(minX, maxY)
	brX, brY := ds.tile(maxX, minY)

	toLoad := make(chan imageToLoad, (brX-tlX)*(brY-tlY))
	toAddToCache := make(chan loadedImage, (brX-tlX)*(brY-tlY))

	// cleanup outside of loaded area
	for i := range len(ds.band) {
		for j := range len(ds.band[i]) {
			if !(i >= tlX && i <= brX && j >= brY && j <= tlY) {
				ds.band[i][j].img = nil
			}
		}
	}

	// load new tiles
	fileCount := 0
	for i := tlX; i <= brX; i++ {
		for j := tlY; j <= brY; j++ {
			if i < 0 || j < 0 || i >= len(ds.band) || j >= len(ds.band[i]) {
				continue
			}
			tile := ds.band[i][j]
			if tile.file != "" && tile.img == nil {
				toLoad <- imageToLoad{tile, i, j}
				fileCount++
			}
		}
	}
	close(toLoad)

	for range runtime.NumCPU() {
		go loadWorker(toLoad, toAddToCache)
	}

	for range fileCount {
		tile := <-toAddToCache
		if tile.err != nil {
			slog.Warn("Error opening a tile file", "File", tile.file, "Error", tile.err)
			continue
		}
		ds.band[tile.x][tile.y] = tile.datasetTile
	}
}
