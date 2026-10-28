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
			continue
		}
		if correspondingRoad.Geom.NumLineStrings() == 0 {
			slog.Warn("Skipping edge of road without geometry", "gid", edge.Gid)
			continue
		}
		// shares the road's flat coordinates: only read from it
		line := correspondingRoad.Geom.LineString(0)
		if line.Stride() < 2 || len(line.FlatCoords()) < line.Stride() {
			slog.Warn("Skipping edge of road without geometry", "gid", edge.Gid)
			continue
		}
		srid := correspondingRoad.Geom.SRID()
		displacement := roadWidth(correspondingRoad) / 2

		// add in original direction
		edges = append(edges, gotopo.Edge{
			Gid:      edge.Gid,
			StartVid: edge.StartVid,
			EndVid:   edge.EndVid,
			Geom:     offsetLine(line, srid, displacement, false),
		})

		if shouldSplit(correspondingRoad) {
			// add in reverse direction
			edges = append(edges, gotopo.Edge{
				Gid:      edge.Gid,
				StartVid: edge.EndVid,
				EndVid:   edge.StartVid,
				Geom:     offsetLine(line, srid, displacement, true),
			})
		}
	}

	return edges
}

func shouldSplit(road *gotopo.GetRoadsRow) bool {
	return true
}

// offsetLine : build a new line (one fresh flat buffer) from source, optionally reversed, with each point
// offset perpendicularly to its segment. Displacements are always computed from the unmodified source.
func offsetLine(source *geom.LineString, srid int, displacement float64, reverse bool) geom.LineString {
	src := source.FlatCoords()
	stride := source.Stride()
	n := len(src) / stride

	// offset in src of the k-th point of the output line
	sourceOffset := func(k int) int {
		if reverse {
			return (n - 1 - k) * stride
		}
		return k * stride
	}

	dst := make([]float64, n*stride)
	for k := 0; k < n; k++ {
		current := sourceOffset(k)
		out := k * stride
		// keep every dimension (Z, M) so the stride matches the layout
		copy(dst[out:out+stride], src[current:current+stride])

		if n > 1 {
			var from, to int
			if k < n-1 {
				from, to = current, sourceOffset(k+1)
			} else {
				from, to = sourceOffset(k-1), current
			}
			displacementX, displacementY := displacementVector(displacement, src[from:from+stride], src[to:to+stride])
			dst[out] += displacementX
			dst[out+1] += displacementY
		}
	}

	return *geom.NewLineStringFlat(source.Layout(), dst).SetSRID(srid)
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
