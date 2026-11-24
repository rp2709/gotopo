package graph

import (
	"fmt"
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
	return fmt.Sprintf("{name:%s, weight:%f%}", m.name, m.weight*100)
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

func ShortestPath(g []gotopo.GetGraphWithinRow, from int32, to int32, t time.Time, metrics Metrics, measureDuration func(agg gotopo.EdgeAgg) time.Duration) (path Path, err error) {
	dist := make(map[int32]float64, len(g))
	ts := make(map[int32]time.Time, len(g))
	pred := make(map[int32]*gotopo.EdgeAgg, len(g))
	for _, u := range g {
		dist[u.VertexID] = math.MaxFloat64
		pred[u.VertexID] = nil
	}
	dist[from] = 0
	ts[from] = t

	L := make(map[int32]gotopo.GetGraphWithinRow, len(g))
	for _, u := range g {
		L[u.VertexID] = u
	}

	if _, ok := L[from]; !ok {
		return nil, fmt.Errorf("from is not a vertex in the graph")
	}

	if _, ok := L[to]; !ok {
		return nil, fmt.Errorf("to is not a vertex in the graph")
	}

	for len(L) > 0 {
		u := popVertexWithSmallestDistance(L, dist)
		if dist[u.VertexID] == math.MaxFloat64 {
			return nil, fmt.Errorf("vertices are not connected")
		}

		if u.VertexID == to {
			break // we found the destination, but must still reconstruct path
		}

		for _, e := range u.Edges {
			v := e.NeighbourID
			eCost := metrics.Evaluate(e, ts[u.VertexID])
			if dist[v] > dist[u.VertexID]+eCost {
				dist[v] = dist[u.VertexID] + eCost
				ts[v] = ts[u.VertexID].Add(measureDuration(e))
				pred[v] = &e
			}
		}
	}

	//reconstruct path
	path = make([]gotopo.EdgeAgg, 0)
	current := to
	for e := pred[current]; e != nil; e = pred[current] {
		path = append(path, *e)
	}
	slices.Reverse(path)
	return path, nil
}

func popVertexWithSmallestDistance(queue map[int32]gotopo.GetGraphWithinRow, distance map[int32]float64) gotopo.GetGraphWithinRow {
	// find vertex with smallest distance from start
	minDist := math.MaxFloat64
	var minVertexId int32 = -1
	for k, _ := range queue {
		distV := distance[k]
		if distV < minDist {
			minDist = distV
			minVertexId = k
		}
	}
	u := queue[minVertexId]
	delete(queue, minVertexId)
	return u
}
