package main

import (
	"gotopo/quadtree"
	"gotopo/queue"
	"math"
)

const none int32 = -1

func ConstructGraph(roadPoints *[]RoadPoint, startingLocation quadtree.Point, epsilon float64) *Graph {
	graph := Graph{
		Vertices: make([]Vertex, 0),
		Edges:    make([]Edge, 0),
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

	roadindex := quadtree.NewQuadTree[*RoadPoint](20, minX, maxX, minY, maxY)

	for _, roadPoint := range *roadPoints {
		roadindex.Insert(&roadPoint)
	}

	exploredindex := quadtree.NewQuadTree[*RoadPoint](20, minX, maxX, minY, maxY)

	return roadindex, exploredindex
}
