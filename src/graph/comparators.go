package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)

// Check two mixed graph equality
func MixedGraphsEquals(gr1, gr2 MixedGraphReader) bool {
	// checking nodes equality
	if !GraphIncludeNodes(gr1, gr2) || !GraphIncludeNodes(gr2, gr1) {
		return false
	}
	
	// checking connections equality
	if !MixedGraphIncludeConnections(gr1, gr2) || !MixedGraphIncludeConnections(gr2, gr1) {
		return false
	}
	
	return true
}

// Check if graph gr include all connections
//
// Warning!!! Due to channels issue 296: http://code.google.com/p/go/issues/detail?id=296
// goroutine will block if not all nodes exists in graph
func MixedGraphIncludeConnections(gr MixedGraphReader, connections TypedConnectionsIterable) bool {
	for conn := range connections.TypedConnectionsIter() {
		switch conn.Type {
			case CT_UNDIRECTED:
				if !gr.CheckEdge(conn.Tail, conn.Head) {
					return false
				}
			case CT_DIRECTED:
				if !gr.CheckArc(conn.Tail, conn.Head) {
					return false
				}
			default:
				err := erx.NewError("Internal error: unknown connection type")
				panic(err)
		}
	}
	return true
}

// Check if graph gr include all nodes from nodesToCheck
//
// Warning!!! Due to channels issue 296: http://code.google.com/p/go/issues/detail?id=296
// goroutine will block if not all nodes exists in graph
func GraphIncludeNodes(gr NodesChecker, nodesToCheck NodesIterable) bool {
	for node := range nodesToCheck.NodesIter() {
		if !gr.CheckNode(node) {
			return false
		}
	}
	return true
} 
