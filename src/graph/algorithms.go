package graph

func ReduceDirectPaths(og DirectedGraphReader, rg DirectedGraphArcsWriter) {
	for conn := range og.ConnectionsIter() {
		filteredGraph := NewArcFilter(og, conn.Tail, conn.Head)
		if !CheckPathDijkstra(filteredGraph, conn.Tail, conn.Head, -1, SimpleWeightFunc) {
			rg.AddArc(conn.Tail, conn.Head)
		}
	}
}
