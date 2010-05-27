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
// goroutine will block if not all connections exists in graph
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

// Check if graph gr include all edges from edgesToCheck
//
// Warning!!! Due to channels issue 296: http://code.google.com/p/go/issues/detail?id=296
// goroutine will block if function result is false
func GraphIncludeEdges(gr UndirectedGraphReader, edgesToCheck EdgesIterable) bool {
	for conn := range edgesToCheck.EdgesIter() {
		if !gr.CheckEdge(conn.Tail, conn.Head) {
			return false
		}
	}
	return true
}

// Check if graph gr include all arcs from edgesToCheck
//
// Warning!!! Due to channels issue 296: http://code.google.com/p/go/issues/detail?id=296
// goroutine will block if function result is false
func GraphIncludeArcs(gr DirectedGraphReader, arcsToCheck ArcsIterable) bool {
	for conn := range arcsToCheck.ArcsIter() {
		if !gr.CheckArc(conn.Tail, conn.Head) {
			return false
		}
	}
	return true
}

// Check if graph gr include all arcs and all nodes from gr2
//
// Warning!!! Due to channels issue 296: http://code.google.com/p/go/issues/detail?id=296
// goroutine will block if function result is false
func DirectedGraphInclude(gr1, gr2 DirectedGraphReader) bool {
	if !GraphIncludeNodes(gr1, gr2) {
		return false
	}
	
	if !GraphIncludeArcs(gr1, gr2) {
		return false
	}
	
	return true
}

// Check if graph gr include all edges and all nodes from gr2
//
// Warning!!! Due to channels issue 296: http://code.google.com/p/go/issues/detail?id=296
// goroutine will block if function result is false
func UndirectedGraphInclude(gr1, gr2 DirectedGraphReader) bool {
	if !GraphIncludeNodes(gr1, gr2) {
		return false
	}
	
	if !GraphIncludeArcs(gr1, gr2) {
		return false
	}
	
	return true
}

// Check if two directed grahps are equal
//
// Warning!!! Due to channels issue 296: http://code.google.com/p/go/issues/detail?id=296
// goroutine will block if function result is false
func DirectedGraphsEquals(gr1, gr2 DirectedGraphReader) bool {
	if !GraphIncludeNodes(gr1, gr2) || !GraphIncludeNodes(gr2, gr1) {
		return false
	}
	
	if !GraphIncludeArcs(gr1, gr2) || !GraphIncludeArcs(gr2, gr1) {
		return false
	}
	
	return true
}

// Check if two undirected grahps are equal
//
// Warning!!! Due to channels issue 296: http://code.google.com/p/go/issues/detail?id=296
// goroutine will block if function result is false
func UndirectedGraphsEquals(gr1, gr2 UndirectedGraphReader) bool {
	if !GraphIncludeNodes(gr1, gr2) || !GraphIncludeNodes(gr2, gr1) {
		return false
	}
	
	if !GraphIncludeEdges(gr1, gr2) || !GraphIncludeEdges(gr2, gr1) {
		return false
	}
	
	return true
}
