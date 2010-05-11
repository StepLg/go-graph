package graph

func ReduceDirectPaths(og DirectedGraphReader, rg DirectedGraphArcsWriter, stopFunc func(from, to NodeId, weight float) bool) {
	var checkStopFunc StopFunc
	for conn := range og.ConnectionsIter() {
		filteredGraph := NewArcFilter(og, conn.Tail, conn.Head)
		if stopFunc!=nil {
			checkStopFunc = func(node NodeId, weight float) bool {
				return stopFunc(conn.Tail, node, weight)
			}
		} else {
			checkStopFunc = nil
		}
		if !CheckDirectedPathDijkstra(filteredGraph, conn.Tail, conn.Head, checkStopFunc, SimpleWeightFunc) {
			rg.AddArc(conn.Tail, conn.Head)
		}
	}
}
