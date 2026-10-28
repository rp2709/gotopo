package main

import (
	"context"
	"gotopo/gotopo"
	"log"
	"math"

	"github.com/twpayne/go-geom"
)

// applies a series of transformations to the edges of a graph to get directional location of edges

// DoubleEdges : double each edge (with exceptions) and associate geometry (use gotopo/InsertEdgesParams)
func DoubleEdges(simpleEdges *[]Edge, ctx context.Context, queries *gotopo.Queries) *[]gotopo.InsertEdgeParams {
	// get all roads from db
	roads, err := queries.GetRoads(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// construct index of roads by gid
	roadsIndex := make(map[int32]*gotopo.GetRoadsRow)
	for i := range roads {
		roadsIndex[roads[i].Gid.Int32] = &roads[i]
	}

	// construct slice of edges
	edges := make([]gotopo.InsertEdgeParams, 0, len(*simpleEdges)*2) // we will potentially double the number of simple edges
	for _, edge := range *simpleEdges {
		correspondingRoad, ok := roadsIndex[edge.gid]
		if !ok {
			continue
		}
		// add in original direction
		edges = append(edges, gotopo.InsertEdgeParams{
			Gid:      edge.gid,
			StartVid: edge.start,
			EndVid:   edge.end,
			Geom:     *correspondingRoad.Geom.Clone(), // deep copy: the geometry is modified later
		})

		if shouldSplit(correspondingRoad) {
			// add in reverse direction
			// deep copy, otherwise Reverse() would also reverse the original (shared flat coords)
			reversedGeometry := *correspondingRoad.Geom.Clone()
			reversedGeometry.Reverse()
			edges = append(edges, gotopo.InsertEdgeParams{
				Gid:      edge.gid,
				StartVid: edge.end,
				EndVid:   edge.start,
				Geom:     reversedGeometry,
			})
		}
	}

	offsetGeometry(&edges, &roadsIndex)

	return &edges
}

func shouldSplit(road *gotopo.GetRoadsRow) bool {
	return true
}

// offset_geometry : offset each edge by a small amount to get directional location of edges (use gotopo/InsertEdgesParams)
func offsetGeometry(edges *[]gotopo.InsertEdgeParams, roadIndex *map[int32]*gotopo.GetRoadsRow) {
	for i := range *edges {
		edge := &(*edges)[i] // pointer to be able to modify the element
		displacement := roadWidth((*roadIndex)[edge.Gid]) / 2

		// Coords() returns a deep copy of all linestrings (each point has the full stride, e.g. X, Y, Z)
		lines := edge.Geom.Coords()
		newLines := make([][]geom.Coord, len(lines))
		for l, line := range lines {
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
			newLines[l] = newLine
		}

		// set coordinates on the whole MultiLineString (LineString(i) only returns a temporary copy)
		if _, err := edge.Geom.SetCoords(newLines); err != nil {
			log.Fatal(err)
		}
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
