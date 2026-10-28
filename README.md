# Gotopo
The swissTLM3D dataset published by Swisstopo represents roads with simple polylines and doesn't represent connectivity between segments.
*Gotopo* constructs a graph representation of the read network. The graph is oriented and each edge stores the id of the original swissTLM3D road, start and end vertices and geometry.
The geometry of edges is adjusted so it represents the actual path of a cyclist on the road.
