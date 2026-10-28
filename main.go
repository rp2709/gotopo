package main

import (
	"context"
	"gotopo/gotopo"
	"gotopo/quadtree"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/twpayne/go-geom"
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
	log.Default().Println("Connected to database")
	err = pgxgeom.Register(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}

	queries := gotopo.New(conn)

	// road points live in this buffer
	log.Default().Println("Start loading roads")
	roadPoints := GetRoadPoints(queries, ctx)
	log.Default().Println("Loaded all roads")

	log.Default().Println("Start constructing graph")
	graph := ConstructGraph(roadPoints, quadtree.Point{2540215.68, 1181115.23}, 0.5)
	log.Default().Println("Graph created")
	clear(*roadPoints)

	// save graph to database
	println("n :", len(graph.Vertices))
	println("m :", len(graph.Edges))

	println("start writing vertices")
	fivePercentOfVertices := len(graph.Vertices) / 20
	for i, vertex := range graph.Vertices {
		v := gotopo.InsertVertexParams{
			Vid: vertex.Vid,
		}
		v.Location.SetSRID(2056).SetCoords(geom.Coord{vertex.Location.X, vertex.Location.Y})
		err := queries.InsertVertex(ctx, v)
		if err != nil {
			log.Fatal(err)
		}
		if i%fivePercentOfVertices == 0 {
			println(i*100/len(graph.Vertices), "%")
		}
	}

	println("done")
}

// ------------------ structs and helpers -----------------------------------

type RoadPoint struct {
	Start bool
	Other *RoadPoint
	Point quadtree.Point
	Gid   int32 //id of the original swisstlm3d strasse
}

func (r RoadPoint) Coordinates() quadtree.Point {
	return r.Point
}

func GetRoadPoints(queries *gotopo.Queries, ctx context.Context) *[]RoadPoint {
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
	return &roadPoints
}
