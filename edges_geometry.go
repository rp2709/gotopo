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
		displacement := roadWidth(correspondingRoad) / 2

		// add in original direction
		edges = append(edges, gotopo.Edge{
			Gid:      edge.Gid,
			StartVid: edge.StartVid,
			EndVid:   edge.EndVid,
			Geom:     offsetLine(line, displacement, false),
		})

		if shouldSplit(correspondingRoad) {
			// add in reverse direction
			edges = append(edges, gotopo.Edge{
				Gid:      edge.Gid,
				StartVid: edge.EndVid,
				EndVid:   edge.StartVid,
				Geom:     offsetLine(line, displacement, true),
			})
		}
	}

	return edges
}

func shouldSplit(road *gotopo.GetRoadsRow) bool {
	return true
}

func offsetLine(line *geom.LineString, displacement float64, reverse bool) geom.LineString {
	var displacementX, displacementY float64
	newCoords := make([]geom.Coord, len(line.Coords()))

	index := func(it int) int {
		if reverse {
			return len(line.Coords()) - it - 1
		}
		return it
	}

	for iteration := range line.Coords() {
		current := line.Coords()[index(iteration)]
		if iteration < len(line.Coords())-1 {
			displacementX, displacementY = displacementVector(displacement, current, line.Coords()[index(iteration+1)])
		} else {
			displacementX, displacementY = displacementVector(displacement, line.Coords()[index(iteration-1)], current)
		}

		newCoords[index(iteration)] = geom.Coord{
			current.X() + displacementX,
			current.Y() + displacementY}
	}

	return *geom.NewLineString(geom.XY).SetSRID(line.SRID()).MustSetCoords(newCoords)
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
