package main

import (
	"context"
	"gotopo/gotopo"
	"gotopo/quadtree"
	"log"

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

	graph := ConstructGraph(queries, ctx, quadtree.Point{X: 2540215.68, Y: 1181115.23}, 0.5)

	locatedEdges := DoubleEdges(&graph.Edges, ctx, queries)

	SaveGraph(ctx, conn, nil, locatedEdges, graph.Srid)
}
