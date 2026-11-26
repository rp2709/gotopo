package vrt_reader

import (
	"log/slog"
	"os"
	"runtime"
	"sync"

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

func loadWorker(toLoad <-chan imageToLoad, toAddToCache chan<- loadedImage) {
	for e := range toLoad {
		file, err := os.OpenFile(e.file, os.O_RDONLY, os.ModePerm)
		if err != nil {
			toAddToCache <- loadedImage{imageToLoad: e, err: err}
			continue
		}

		e.img, err = geotiff.Read(file)
		_ = file.Close()
		toAddToCache <- loadedImage{
			imageToLoad: e,
			err:         err,
		}
	}
}

func (ds *RasterDataset) LoadArea(minX, maxX, minY, maxY float64) {
	tlX, tlY := ds.tile(minX, maxY)
	brX, brY := ds.tile(maxX, minY)

	// cleanup outside of loaded area
	for i := range len(ds.band) {
		for j := range len(ds.band[i]) {
			if !(i >= tlX && i <= brX && j >= tlY && j <= brY) {
				ds.band[i][j].img = nil
			}
		}
	}

	if len(ds.band) == 0 {
		return
	}

	loadMinX, loadMaxX := max(tlX, 0), min(brX, len(ds.band)-1)
	if loadMinX > loadMaxX {
		return
	}
	loadMinY, loadMaxY := max(tlY, 0), min(brY, len(ds.band[0])-1)
	if loadMinY > loadMaxY {
		return
	}

	workerCount := min(runtime.NumCPU(), (loadMaxX-loadMinX+1)*(loadMaxY-loadMinY+1))
	toLoad := make(chan imageToLoad)
	toAddToCache := make(chan loadedImage)

	var workers sync.WaitGroup
	workers.Add(workerCount)
	for range workerCount {
		go func() {
			defer workers.Done()
			loadWorker(toLoad, toAddToCache)
		}()
	}

	go func() {
		defer close(toLoad)
		for i := loadMinX; i <= loadMaxX; i++ {
			for j := loadMinY; j <= loadMaxY && j < len(ds.band[i]); j++ {
				tile := ds.band[i][j]
				if tile.file != "" && tile.img == nil {
					toLoad <- imageToLoad{tile, i, j}
				}
			}
		}
	}()

	go func() {
		workers.Wait()
		close(toAddToCache)
	}()

	for tile := range toAddToCache {
		if tile.err != nil {
			slog.Warn("Error opening a tile file", "File", tile.file, "Error", tile.err)
			continue
		}
		ds.band[tile.x][tile.y] = tile.datasetTile
	}
}
