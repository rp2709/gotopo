package quadtree

import "slices"

type Point struct {
	X, Y float64
}

type PlanarIndexable interface {
	Coordinates() Point
}

type QuadTree[T PlanarIndexable] struct {
	capacity int
	size     int
	maxDepth int
	root     *node[T]
}

func (q *QuadTree[T]) Size() int     { return q.size }
func (q *QuadTree[T]) MaxDepth() int { return q.maxDepth }
func (q *QuadTree[T]) Capacity() int { return q.capacity }

func (q *QuadTree[T]) Insert(object T) {
	q.root.insert(object)
	q.size++
}

func (q *QuadTree[T]) GetWithin(distance float64, from Point) []T {
	return q.root.getWithin(distance, from)
}

// NewQuadTree creates a new quadtree with capacity at each leaf, set to contain points within range
func NewQuadTree[T PlanarIndexable](capacity int, minX float64, maxX float64, minY float64, maxY float64) *QuadTree[T] {
	rootBoundingBox := boundingBox{
		Center: Point{(minX + maxX) / 2, (minY + maxY) / 2},
		DimX:   (maxX - minX) / 2,
		DimY:   (maxY - minY) / 2,
	}
	qt := &QuadTree[T]{
		capacity: capacity,
		size:     0,
		maxDepth: 0,
	}

	qt.root = newNode(rootBoundingBox, qt, 0)
	return qt
}

// center with half dimensions
type boundingBox struct {
	Center Point
	DimX   float64
	DimY   float64
}

type node[T PlanarIndexable] struct {
	boundingBox

	isLeaf bool
	depth  int

	nw *node[T]
	ne *node[T]
	sw *node[T]
	se *node[T]

	objects []T

	quadtree *QuadTree[T] // quadtree to which this node belongs to
}

func newNode[T PlanarIndexable](box boundingBox, quadtree *QuadTree[T], depth int) *node[T] {
	if depth > quadtree.maxDepth {
		quadtree.maxDepth = depth
	}

	return &node[T]{
		boundingBox: box,
		isLeaf:      true,
		objects:     make([]T, 0, quadtree.capacity),
		quadtree:    quadtree,
		depth:       depth,
	}
}

func (n *node[T]) insert(object T) {
	if n.isLeaf {
		// easy if still enough place within this leaf
		if len(n.objects) < n.quadtree.capacity {
			n.objects = append(n.objects, object)
			return
		}
		// capacity reached, split node in 4
		n.split()
	}
	// at this point node is not a leaf
	n.quadrant(object.Coordinates()).insert(object)
}

// returns the child node that should contain object
func (n *node[T]) quadrant(p Point) *node[T] {
	if n.isLeaf {
		panic("leaf node has no children")
	}
	if p.X < n.Center.X { //west
		if p.Y < n.Center.Y { //south
			return n.sw
		}
		//north
		return n.nw
	}
	//east
	if p.Y < n.Center.Y { //south
		return n.se
	}
	//north
	return n.ne
}

// splits this leaf into 4 quadrants
func (n *node[T]) split() {
	if !n.isLeaf {
		panic("cannot split a non-leaf node")
	}

	n.nw = newNode[T](n.boundingBox.nw(), n.quadtree, n.depth+1)
	n.ne = newNode[T](n.boundingBox.ne(), n.quadtree, n.depth+1)
	n.sw = newNode[T](n.boundingBox.sw(), n.quadtree, n.depth+1)
	n.se = newNode[T](n.boundingBox.se(), n.quadtree, n.depth+1)

	n.isLeaf = false
	for _, object := range n.objects {
		n.insert(object)
	}

	n.objects = nil
}

func (b boundingBox) nw() boundingBox {
	dimx := b.DimX / 2
	dimy := b.DimY / 2
	return boundingBox{
		DimX:   dimx,
		DimY:   dimy,
		Center: Point{b.Center.X - dimx, b.Center.Y + dimy},
	}
}

func (b boundingBox) ne() boundingBox {
	dimx := b.DimX / 2
	dimy := b.DimY / 2
	return boundingBox{
		DimX:   dimx,
		DimY:   dimy,
		Center: Point{b.Center.X + dimx, b.Center.Y + dimy},
	}
}

func (b boundingBox) sw() boundingBox {
	dimx := b.DimX / 2
	dimy := b.DimY / 2
	return boundingBox{
		DimX:   dimx,
		DimY:   dimy,
		Center: Point{b.Center.X - dimx, b.Center.Y - dimy},
	}
}

func (b boundingBox) se() boundingBox {
	dimx := b.DimX / 2
	dimy := b.DimY / 2
	return boundingBox{
		DimX:   dimx,
		DimY:   dimy,
		Center: Point{b.Center.X + dimx, b.Center.Y - dimy},
	}
}

func (n *node[T]) getWithin(distance float64, from Point) []T {
	if n.isLeaf {
		within := make([]T, 0, cap(n.objects))
		for _, object := range n.objects {
			if object.Coordinates().isWithin(distance, from) {
				within = append(within, object)
			}
		}
		return within
	}
	// take the four corners
	corners := []Point{
		{from.X + distance, from.Y + distance},
		{from.X - distance, from.Y + distance},
		{from.X + distance, from.Y - distance},
		{from.X - distance, from.Y - distance},
	}
	quadrants := make([]*node[T], 0, 4)
	for _, corner := range corners {
		quadrant := n.quadrant(corner)
		if !slices.Contains(quadrants, quadrant) {
			quadrants = append(quadrants, quadrant)
		}
	}
	within := make([]T, 0)
	for _, quadrant := range quadrants {
		within = append(within, quadrant.getWithin(distance, from)...)
	}
	return within
}

func (p1 Point) isWithin(distance float64, p2 Point) bool {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	dxs := dx * dx
	dys := dy * dy
	ds := dxs + dys
	wds := distance * distance
	return ds <= wds
}
