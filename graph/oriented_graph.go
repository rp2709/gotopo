package graph

import (
	"container/heap"
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
	if _, ok := g[from]; !ok {
		return nil, fmt.Errorf("from is not a vertex in the graph")
	}
	if _, ok := g[to]; !ok {
		return nil, fmt.Errorf("to is not a vertex in the graph")
	}

	dist := make(map[int32]float64, len(g))
	ts := make(map[int32]time.Time, len(g))
	pred := make(map[int32]*gotopo.EdgeAgg, len(g))
	for id := range g {
		dist[id] = math.MaxFloat64
		pred[id] = nil
	}
	dist[from] = 0
	ts[from] = t

	remaining := make(map[int32]Vertex, len(g))
	maps.Copy(remaining, g)

	queue := newPriorityQueue(dist)
	for vid := range g {
		heap.Push(queue, vid)
	}

	for len(remaining) > 0 {
		u := popVertexWithSmallestDistance(remaining, queue)
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
				// The queue shares dist as its priority map, so a relaxation
				// is a decrease-key: restore the heap invariant for v.
				if i, ok := queue.index[v]; ok {
					heap.Fix(queue, i)
				}
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

func popVertexWithSmallestDistance(remaining Graph, queue *priorityQueue) Vertex {
	id := heap.Pop(queue).(int32)
	v := remaining[id]
	delete(remaining, id)
	return v
}

// priorityQueue is an indexed binary min-heap of vertex ids. priorities is
// shared with the algorithm's dist map, so decreasing dist[v] changes v's
// key; index lets heap.Fix restore the heap order after such a change.
type priorityQueue struct {
	data       []int32
	priorities map[int32]float64
	index      map[int32]int // vertex id -> position in data
}

func newPriorityQueue(priorities map[int32]float64) *priorityQueue {
	return &priorityQueue{
		data:       make([]int32, 0, len(priorities)),
		priorities: priorities,
		index:      make(map[int32]int, len(priorities)),
	}
}

func (pq *priorityQueue) Len() int { return len(pq.data) }

func (pq *priorityQueue) Less(i, j int) bool {
	return pq.priorities[pq.data[i]] < pq.priorities[pq.data[j]]
}

func (pq *priorityQueue) Swap(i, j int) {
	pq.data[i], pq.data[j] = pq.data[j], pq.data[i]
	pq.index[pq.data[i]] = i
	pq.index[pq.data[j]] = j
}

func (pq *priorityQueue) Push(x any) {
	vid := x.(int32)
	pq.index[vid] = len(pq.data)
	pq.data = append(pq.data, vid)
}

func (pq *priorityQueue) Pop() any {
	vid := pq.data[len(pq.data)-1]
	delete(pq.index, vid)
	pq.data = pq.data[:len(pq.data)-1]
	return vid
}
