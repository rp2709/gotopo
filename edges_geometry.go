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
		correspondingRoad := roadsIndex[edge.gid]
		// add in original direction
		edges = append(edges, gotopo.InsertEdgeParams{
			Gid:      edge.gid,
			StartVid: edge.start,
			EndVid:   edge.end,
			Geom:     correspondingRoad.Geom,
		})

		if shouldSplit(correspondingRoad) {
			// add in reverse direction
			reversedGeometry := correspondingRoad.Geom
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
		newCoords := make([]geom.Coord, len(edge.Geom.LineString(0).Coords()))
		for coord := range edge.Geom.LineString(0).Coords() {
			displacement := roadWidth((*roadIndex)[edge.Gid]) / 2
			current := edge.Geom.LineString(0).Coords()[coord]
			var displacementX float64
			var displacementY float64

			if coord < len(newCoords)-1 {
				displacementX, displacementY = displacementVector(displacement, current, edge.Geom.LineString(0).Coords()[coord+1])
			} else {
				displacementX, displacementY = displacementVector(displacement, edge.Geom.LineString(0).Coords()[coord-1], current)
			}

			newCoords[coord] = geom.Coord{
				current.X() + displacementX,
				current.Y() + displacementY,
			}
		}
		edge.Geom.LineString(0).MustSetCoords(newCoords)
	}
}

func displacementVector(displacement float64, start geom.Coord, end geom.Coord) (float64, float64) {
	dx := end.X() - start.X()
	dy := end.Y() - start.Y()
	norm := math.Sqrt(dx*dx + dy*dy)
	perpendicularX := dy / norm
	perpendicularY := -dx / norm
	displacementX := perpendicularX * displacement
	displacementY := perpendicularY * displacement
	return displacementX, displacementY
}

func roadWidth(road *gotopo.GetRoadsRow) float64 {
	return 6.0
}
