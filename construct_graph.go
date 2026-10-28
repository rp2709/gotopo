package main

import (
	"context"
	"gotopo/gotopo"
	"gotopo/quadtree"
	"gotopo/queue"
	"log"
	"math"
)

const none int32 = -1

/** ConstructGraph reads roads from database and constructs graph
 * conceptually, it starts at a given location and treats points in cluster
 * a cluster is a set of points that are within epsilon distance of each other
 */
func ConstructGraph(queries *gotopo.Queries, ctx context.Context, startingLocation quadtree.Point, epsilon float64) *Graph {
	roadPoints, srid := getRoadPointsAndSrid(queries, ctx)

	graph := Graph{
		Vertices: make([]Vertex, 0),
		Edges:    make([]Edge, 0),
		Srid:     srid,
	}

	var nextVertexId int32 = 0

	roadPointsIndex, exploredIndex := getRoadIndex(roadPoints)

	// queue fifo of RoadPoints to explore next
	opened := queue.NewRingBuffer[*RoadPoint](50000)

	incompleteEdges := make(map[int32]Edge)

	startRoadPoint := roadPointsIndex.GetNearest(startingLocation)
	opened.Push(startRoadPoint)
	exploredIndex.Insert(startRoadPoint)

	for opened.Len() > 0 { // while there are opened ends
		// get cluster that may be considered as a single vertex
		cluster := roadPointsIndex.GetWithin(epsilon, opened.Pop().Point)
		vertex := Vertex{
			Location: cluster[0].Point,
			Vid:      nextVertexId,
		}
		nextVertexId++
		graph.Vertices = append(graph.Vertices, vertex)

		// create new edge or attach to incomplete
		for _, point := range cluster {
			if incompleteEdge, ok := incompleteEdges[point.Gid]; ok {
				// attach to existing incomplete edge
				edge := incompleteEdge
				if point.Start {
					if incompleteEdge.start != none {
						panic("incomplete edge already has a start")
					}
					edge.start = vertex.Vid
				} else {
					if incompleteEdge.end != none {
						panic("incomplete edge already has an end")
					}
					edge.end = vertex.Vid
				}
				delete(incompleteEdges, point.Gid)
				graph.Edges = append(graph.Edges, edge)
			} else {
				// create new edge
				if point.Start {
					incompleteEdges[point.Gid] = Edge{vertex.Vid, none, point.Gid}
				} else {
					incompleteEdges[point.Gid] = Edge{none, vertex.Vid, point.Gid}
				}
				// add other end to the queue if not equivalent to other opened point
				if len(exploredIndex.GetWithin(epsilon, point.Other.Point)) == 0 {
					opened.Push(point.Other)
					exploredIndex.Insert(point.Other)
				}
			}
		}
	}
	return &graph
}

type Graph struct {
	Srid     int
	Vertices []Vertex
	Edges    []Edge
}

type Vertex struct {
	Vid      int32
	Location quadtree.Point
}

func (v *Vertex) Coordinates() quadtree.Point { return v.Location }

type Edge struct {
	start int32
	end   int32
	gid   int32
}

func getRoadIndex(roadPoints *[]RoadPoint) (*quadtree.QuadTree[*RoadPoint], *quadtree.QuadTree[*RoadPoint]) {
	minX := math.MaxFloat64
	minY := math.MaxFloat64
	maxX := -math.MaxFloat64
	maxY := -math.MaxFloat64

	for _, roadPoint := range *roadPoints {
		minX = min(minX, roadPoint.Point.X)
		minY = min(minY, roadPoint.Point.Y)
		maxX = max(maxX, roadPoint.Point.X)
		maxY = max(maxY, roadPoint.Point.Y)
	}

	roadIndex := quadtree.NewQuadTree[*RoadPoint](20, minX, maxX, minY, maxY)

	for _, roadPoint := range *roadPoints {
		roadIndex.Insert(&roadPoint)
	}

	exploredIndex := quadtree.NewQuadTree[*RoadPoint](20, minX, maxX, minY, maxY)

	return roadIndex, exploredIndex
}

type RoadPoint struct {
	Start bool
	Other *RoadPoint
	Point quadtree.Point
	Gid   int32 //id of the original swisstlm3d strasse
}

func (r RoadPoint) Coordinates() quadtree.Point {
	return r.Point
}

func getRoadPointsAndSrid(queries *gotopo.Queries, ctx context.Context) (*[]RoadPoint, int) {
	roads, err := queries.GetRoads(ctx)
	if err != nil {
		log.Fatal(err)
	}
	roadPoints := make([]RoadPoint, 0, 2*len(roads)) // important to pre-allocate so pointers are stable (for indexing)
	for _, road := range roads {
		coords := road.Geom.LineString(0).Coords()
		start := coords[0]
		end := coords[len(coords)-1]

		roadPoints = append(roadPoints,
			RoadPoint{
				Point: quadtree.Point{X: start.X(), Y: start.Y()},
				Gid:   road.Gid.Int32,
				Start: true,
			},
			RoadPoint{
				Point: quadtree.Point{X: end.X(), Y: end.Y()},
				Gid:   road.Gid.Int32,
				Start: false,
			})

		// link the two ends of the segment
		roadPoints[len(roadPoints)-1].Other = &roadPoints[len(roadPoints)-2]
		roadPoints[len(roadPoints)-2].Other = &roadPoints[len(roadPoints)-1]
	}
	return &roadPoints, roads[0].Geom.SRID()
}
