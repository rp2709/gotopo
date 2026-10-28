package graph_construction

import (
	"gotopo/gotopo"
	"gotopo/quadtree"
	"gotopo/queue"
	"log/slog"
	"math"

	"github.com/twpayne/go-geom"
)

const none int32 = -1

/** ConstructGraph constructs graph from roads
 * conceptually, it starts at a given location and treats points in cluster
 * a cluster is a set of points that are within epsilon distance of each other
 */
func ConstructGraph(roads []gotopo.GetRoadsRow, startingLocation quadtree.Point, epsilon float64, explorationLimit int64) Graph {
	roadPoints := getRoadPoints(roads)

	graph := Graph{
		Vertices: make([]gotopo.Vertex, 0),
		Edges:    make([]gotopo.Edge, 0),
	}

	if len(roadPoints) == 0 {
		return graph
	}

	var nextVertexId int32 = 0

	roadPointsIndex, exploredIndex := getRoadIndex(roadPoints)

	// queue fifo of RoadPoints to explore next
	opened := queue.NewRingBuffer[*RoadPoint](50000)

	incompleteEdges := make(map[int32]gotopo.Edge)

	startRoadPoint := roadPointsIndex.GetNearest(startingLocation)
	opened.Push(startRoadPoint)
	exploredIndex.Insert(startRoadPoint)

	var iterations int64 = 0
	for opened.Len() > 0 && iterations < explorationLimit { // while there are opened ends and we haven't reached the limit
		iterations++
		// get cluster that may be considered as a single vertex
		cluster := roadPointsIndex.GetWithin(epsilon, opened.Pop().Point)
		location := cluster[0].Point
		vertex := gotopo.Vertex{
			Vid:      nextVertexId,
			Location: *geom.NewPoint(geom.XY).SetSRID(cluster[0].Road.Geom.SRID()).MustSetCoords(geom.Coord{location.X, location.Y}),
		}
		nextVertexId++
		graph.Vertices = append(graph.Vertices, vertex)

		// create new edge or attach to incomplete
		for _, point := range cluster {
			gid := point.Road.Gid.Int32
			if edge, ok := incompleteEdges[gid]; ok {
				// attach to existing incomplete edge
				if point.Start {
					if edge.StartVid != none {
						panic("incomplete edge already has a start")
					}
					edge.StartVid = vertex.Vid
				} else {
					if edge.EndVid != none {
						panic("incomplete edge already has an end")
					}
					edge.EndVid = vertex.Vid
				}
				delete(incompleteEdges, gid)
				if edge.StartVid == none || edge.EndVid == none {
					panic("Incomplete edge marked as complete")
				}
				if edge.EndVid == edge.StartVid {
					slog.Warn("Skipping edge that forms a closed loop", "vertex_id", edge.StartVid)
					continue
				}
				graph.Edges = append(graph.Edges, edge)
			} else {
				// create new edge
				if point.Start {
					incompleteEdges[gid] = gotopo.Edge{Gid: gid, StartVid: vertex.Vid, EndVid: none}
				} else {
					incompleteEdges[gid] = gotopo.Edge{Gid: gid, StartVid: none, EndVid: vertex.Vid}
				}
				// add other end to the queue if not equivalent to other opened point
				if len(exploredIndex.GetWithin(epsilon, point.Other.Point)) == 0 {
					opened.Push(point.Other)
					exploredIndex.Insert(point.Other)
				}
			}
		}
	}
	return graph
}

type Graph struct {
	Vertices []gotopo.Vertex
	Edges    []gotopo.Edge // Geom empty until DoubleEdges
}

func getRoadIndex(roadPoints []RoadPoint) (*quadtree.QuadTree[*RoadPoint], *quadtree.QuadTree[*RoadPoint]) {
	minX := math.MaxFloat64
	minY := math.MaxFloat64
	maxX := -math.MaxFloat64
	maxY := -math.MaxFloat64

	for i := range roadPoints {
		p := &roadPoints[i].Point
		minX = min(minX, p.X)
		minY = min(minY, p.Y)
		maxX = max(maxX, p.X)
		maxY = max(maxY, p.Y)
	}

	roadIndex := quadtree.NewQuadTree[*RoadPoint](20, minX, maxX, minY, maxY)

	for i := range roadPoints {
		roadIndex.Insert(&roadPoints[i])
	}

	exploredIndex := quadtree.NewQuadTree[*RoadPoint](20, minX, maxX, minY, maxY)

	return roadIndex, exploredIndex
}

type RoadPoint struct {
	Start bool
	Other *RoadPoint
	Point quadtree.Point
	Road  *gotopo.GetRoadsRow // original swisstlm3d strasse
}

func (r *RoadPoint) Coordinates() quadtree.Point {
	return r.Point
}

func getRoadPoints(roads []gotopo.GetRoadsRow) []RoadPoint {
	roadPoints := make([]RoadPoint, 0, 2*len(roads)) // important to pre-allocate so pointers are stable
	for i := range roads {
		road := &roads[i]
		if road.Geom.NumLineStrings() == 0 {
			slog.Warn("Skipping road without geometry", "gid", road.Gid.Int32)
			continue
		}
		line := road.Geom.LineString(0)
		coords := line.Coords()
		roadPoints = append(roadPoints,
			RoadPoint{
				Point: quadtree.Point{X: coords[0].X(), Y: coords[0].Y()},
				Road:  road,
				Start: true,
			},
			RoadPoint{
				Point: quadtree.Point{X: coords[len(coords)-1].X(), Y: coords[len(coords)-1].Y()},
				Road:  road,
				Start: false,
			})

		// link the two ends of the segment
		roadPoints[len(roadPoints)-1].Other = &roadPoints[len(roadPoints)-2]
		roadPoints[len(roadPoints)-2].Other = &roadPoints[len(roadPoints)-1]
	}
	return roadPoints
}
