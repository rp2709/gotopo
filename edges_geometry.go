package main

import (
	"gotopo/gotopo"
	"log/slog"
	"math"

	"github.com/twpayne/go-geom"
)

// applies a series of transformations to the edges of a graph to get directional location of edges

// DoubleEdges : double each edge (with exceptions) and associate its offset geometry
func DoubleEdges(simpleEdges []gotopo.Edge, roads []gotopo.GetRoadsRow) []gotopo.Edge {
	// construct index of roads by gid (pointers into the shared slice, no copies)
	roadsIndex := make(map[int32]*gotopo.GetRoadsRow, len(roads))
	for i := range roads {
		roadsIndex[roads[i].Gid.Int32] = &roads[i]
	}

	// construct slice of edges
	edges := make([]gotopo.Edge, 0, len(simpleEdges)*2) // we will potentially double the number of simple edges
	for i := range simpleEdges {
		edge := &simpleEdges[i]
		correspondingRoad, ok := roadsIndex[edge.Gid]
		if !ok {
			slog.Warn("Skipping edge that does not correspond to a swisstlm3d strasse")
			continue
		}
		// shares the road's flat coordinates: only read from it
		line := correspondingRoad.Geom.LineString(0)

		// add in original direction
		edges = append(edges, gotopo.Edge{
			Gid:      edge.Gid,
			StartVid: edge.StartVid,
			EndVid:   edge.EndVid,
			Geom:     *line.Clone(),
		})

		if shouldSplit(correspondingRoad) {
			// add in reverse direction
			reversedLine := line.Clone()
			reversedLine.Reverse()

			edges = append(edges, gotopo.Edge{
				Gid:      edge.Gid,
				StartVid: edge.EndVid,
				EndVid:   edge.StartVid,
				Geom:     *reversedLine,
			})
		}
	}

	offsetGeometry(&edges, &roadsIndex)

	return edges
}

func shouldSplit(road *gotopo.GetRoadsRow) bool {
	return true
}

// offset_geometry : offset each edge by a small amount to get directional location of edges (use gotopo/InsertEdgesParams)
func offsetGeometry(edges *[]gotopo.Edge, roadIndex *map[int32]*gotopo.GetRoadsRow) {
	for i := range *edges {
		edge := &(*edges)[i] // pointer to be able to modify the element
		displacement := roadWidth((*roadIndex)[edge.Gid]) / 2

		line := edge.Geom.Coords()
		newLine := make([]geom.Coord, len(line))
		for c, current := range line {
			// copy the coordinate to keep every dimension (Z, M) so the stride matches the layout
			newCoord := make(geom.Coord, len(current))
			copy(newCoord, current)

			if len(line) > 1 {
				var displacementX, displacementY float64
				if c < len(line)-1 {
					displacementX, displacementY = displacementVector(displacement, current, line[c+1])
				} else {
					displacementX, displacementY = displacementVector(displacement, line[c-1], current)
				}
				newCoord[0] += displacementX
				newCoord[1] += displacementY
			}
			newLine[c] = newCoord
		}

		edge.Geom.MustSetCoords(newLine)
	}
}

func displacementVector(displacement float64, start geom.Coord, end geom.Coord) (float64, float64) {
	dx := end.X() - start.X()
	dy := end.Y() - start.Y()
	norm := math.Sqrt(dx*dx + dy*dy)
	if norm == 0 {
		return 0, 0
	}
	perpendicularX := dy / norm
	perpendicularY := -dx / norm
	displacementX := perpendicularX * displacement
	displacementY := perpendicularY * displacement
	return displacementX, displacementY
}

func roadWidth(road *gotopo.GetRoadsRow) float64 {
	return 6.0
}
