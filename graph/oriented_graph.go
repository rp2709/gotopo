package graph

import (
	"fmt"
	"maps"
	"math"
	"slices"
	"time"

	"github.com/rp2709/gotopo/gotopo"
)

type Path []gotopo.EdgeAgg

type Metric struct {
	name     string
	evaluate func(agg gotopo.EdgeAgg, t time.Time) float64
	weight   float64
}

func (m Metric) String() string {
	return fmt.Sprintf("{name:%s, weight:%f%%}", m.name, m.weight*100)
}

type Metrics []Metric

/*
Every metric is evaluated, balanced with its weight and then weighted by the edge length
*/
func (metrics Metrics) Evaluate(agg gotopo.EdgeAgg, t time.Time) float64 {
	result := 0.0
	for _, metric := range metrics {
		result += metric.evaluate(agg, t) * metric.weight
	}
	return result * agg.Length
}

func ShortestPath(g Graph, from int32, to int32, t time.Time, metrics Metrics, measureDuration func(agg gotopo.EdgeAgg) time.Duration) (path Path, err error) {
	dist := make(map[int32]float64, len(g))
	ts := make(map[int32]time.Time, len(g))
	pred := make(map[int32]*gotopo.EdgeAgg, len(g))
	for id := range g {
		dist[id] = math.MaxFloat64
		pred[id] = nil
	}
	dist[from] = 0
	ts[from] = t

	remaining := make(Graph, len(g))
	maps.Copy(remaining, g)

	if _, ok := remaining[from]; !ok {
		return nil, fmt.Errorf("from is not a vertex in the graph")
	}

	if _, ok := remaining[to]; !ok {
		return nil, fmt.Errorf("to is not a vertex in the graph")
	}

	for len(remaining) > 0 {
		u := popVertexWithSmallestDistance(remaining, dist)
		if dist[u.ID] == math.MaxFloat64 {
			return nil, fmt.Errorf("vertices are not connected")
		}

		if u.ID == to {
			break // we found the destination, but must still reconstruct path
		}

		for _, e := range u.Edges {
			v := e.NeighbourID
			eCost := metrics.Evaluate(e, ts[u.ID])
			if dist[v] > dist[u.ID]+eCost {
				dist[v] = dist[u.ID] + eCost
				ts[v] = ts[u.ID].Add(measureDuration(e))
				pred[v] = &e
			}
		}
	}

	// reconstruct path by following predecessors back to the start
	path = make(Path, 0)
	for current := to; ; {
		e := pred[current]
		if e == nil {
			break
		}
		path = append(path, *e)
		current = e.StartVid
	}
	slices.Reverse(path)
	return path, nil
}

func popVertexWithSmallestDistance(queue Graph, distance map[int32]float64) Vertex {
	// find vertex with smallest distance from start
	minDist := math.MaxFloat64
	var minVertexID int32 = -1
	for k := range queue {
		distV := distance[k]
		if distV < minDist || minVertexID == -1 {
			minDist = distV
			minVertexID = k
		}
	}
	u := queue[minVertexID]
	delete(queue, minVertexID)
	return u
}
