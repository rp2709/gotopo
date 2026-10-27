package main

import (
	"context"
	"gotopo/gotopo"
	"gotopo/quadtree"
	"log"
	"math"

	"github.com/jackc/pgx/v5"
	pgxgeom "github.com/twpayne/pgx-geom"
)

const ConnectionString = "user=postgres dbname=swisstlm3d sslmode=prefer host=localhost password=postgres"

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	err = pgxgeom.Register(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}

	queries := gotopo.New(conn)

	roadindex := getRoadIndex(GetRoadPoints(queries, ctx))

	println("index size", roadindex.Size())
	println("index max depth", roadindex.MaxDepth())

	clicked_point := quadtree.Point{2540215.90, 1181115.26}
	adjacentRoads := roadindex.GetWithin(3, clicked_point)

	println("adjacent roads")
	for _, road := range adjacentRoads {
		println(road.Road.Gid)
	}
}

type Road struct {
	Gid   int32
	Start *RoadPoint
	End   *RoadPoint
}

type RoadPoint struct {
	Road  *Road
	Point quadtree.Point
}

func (r RoadPoint) Coordinates() quadtree.Point {
	return r.Point
}

func getRoadIndex(roadPoints *[]RoadPoint) *quadtree.QuadTree[RoadPoint] {
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

	roadindex := quadtree.NewQuadTree[RoadPoint](20, minX, maxX, minY, maxY)

	for _, roadPoint := range *roadPoints {
		roadindex.Insert(roadPoint)
	}

	return roadindex
}

func GetRoadPoints(queries *gotopo.Queries, ctx context.Context) *[]RoadPoint {
	roads, err := queries.GetRoads(ctx)
	if err != nil {
		log.Fatal(err)
	}
	roadPoints := make([]RoadPoint, 0, 2*len(roads))
	for _, road := range roads {
		rs := &Road{Gid: road.Gid.Int32}
		coords := road.Geom.LineString(0).Coords()
		start := coords[0]
		end := coords[len(coords)-1]

		roadPoints = append(roadPoints,
			RoadPoint{
				Road:  rs,
				Point: quadtree.Point{X: start.X(), Y: start.Y()},
			},
			RoadPoint{
				Road:  rs,
				Point: quadtree.Point{X: end.X(), Y: end.Y()},
			})
	}
	return &roadPoints
}
