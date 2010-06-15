package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)

// Check two mixed graph equality
func MixedGraphsEquals(gr1, gr2 MixedGraphReader) bool {
	// checking nodes equality
	if !GraphIncludeVertexes(gr1, gr2) || !GraphIncludeVertexes(gr2, gr1) {
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
func GraphIncludeVertexes(gr VertexesChecker, nodesToCheck VertexesIterable) bool {
	for node := range nodesToCheck.VertexesIter() {
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
	if !GraphIncludeVertexes(gr1, gr2) {
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
	if !GraphIncludeVertexes(gr1, gr2) {
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
	if !GraphIncludeVertexes(gr1, gr2) || !GraphIncludeVertexes(gr2, gr1) {
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
	if !GraphIncludeVertexes(gr1, gr2) || !GraphIncludeVertexes(gr2, gr1) {
		return false
	}
	
	if !GraphIncludeEdges(gr1, gr2) || !GraphIncludeEdges(gr2, gr1) {
		return false
	}
	
	return true
}

// Interface for ContainPath function.
type NodeAndConnectionChecker interface {
	// Check if node exist in graph.
	CheckNode(VertexId) bool
	// Check if there is connection between from and to nodes in graph.
	CheckConnection(from, to VertexId) bool
}

// Generic function to check if graph contain specific path.
// 
// First argument gr is an interface with two functions to check node existance and 
// connection existance between two nodes in graph.
// 
// unexistNodePanic flag is used to point wether or not to panic if we figure out that
// one of the nodes in path doesn't exist in graph. If unexistNodePanic is false, then
// result of the function will be false.
func ContainPath(gr NodeAndConnectionChecker, path []VertexId, unexistNodePanic bool) bool {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Checking if graph contain path", e)
			err.AddV("path", path)
			err.AddV("panic if node doesn't exist in graph", unexistNodePanic)
			panic(err)
		}
	}()
	if len(path)==0 {
		// emty path always exists
		return true
	}
	
	prev := path[0]
	if !gr.CheckNode(prev) {
		if unexistNodePanic {
			err := erx.NewError("Node doesn't exist in graph.")
			err.AddV("node", prev)
			panic(err)
		}
		return false
	}
	
	if len(path)==1 {
		return true
	}
	
	for i:=1; i<len(path); i++ {
		cur := path[i]
		if !gr.CheckNode(cur) {
			if unexistNodePanic {
				err := erx.NewError("Node doesn't exist in graph.")
				err.AddV("node", cur)
				panic(err)
			}
			return false
		}
		if !gr.CheckConnection(prev, cur) {
			return false
		}
		prev = cur
	}
	return true
}

// Helper NodeAndConnectionChecker realisation for UndirectedGraph 
type undirectedNodeAndConnectionChecker struct {
	gr UndirectedGraphReader
}

func (checker *undirectedNodeAndConnectionChecker) CheckNode(node VertexId) bool {
	return checker.gr.CheckNode(node)
}

func (checker *undirectedNodeAndConnectionChecker) CheckConnection(from, to VertexId) bool {
	return checker.gr.CheckEdge(from, to)
}

// Helper NodeAndConnectionChecker realisation for DirectedGraph 
type directedNodeAndConnectionChecker struct {
	gr DirectedGraphReader
}

func (checker *directedNodeAndConnectionChecker) CheckNode(node VertexId) bool {
	return checker.gr.CheckNode(node)
}

func (checker *directedNodeAndConnectionChecker) CheckConnection(from, to VertexId) bool {
	return checker.gr.CheckArc(from, to)
}

// Helper NodeAndConnectionChecker realisation for MixedGraph 
type mixedNodeAndConnectionChecker struct {
	gr MixedGraphReader
}

func (checker *mixedNodeAndConnectionChecker) CheckNode(node VertexId) bool {
	return checker.gr.CheckNode(node)
}

func (checker *mixedNodeAndConnectionChecker) CheckConnection(from, to VertexId) bool {
	connType := checker.gr.CheckEdgeType(from, to)
	return connType==CT_DIRECTED || connType==CT_UNDIRECTED
}

// Check if undirected graph contain specific path.
// 
// unexistNodePanic flag is used to point wether or not to panic if we figure out that
// one of the nodes in path doesn't exist in graph. If unexistNodePanic is false, then
// result of the function will be false.
func ContainUndirectedPath(gr UndirectedGraphReader, path []VertexId, unexistNodePanic bool) bool {
	return ContainPath(NodeAndConnectionChecker(&undirectedNodeAndConnectionChecker{gr:gr}), path, unexistNodePanic)
}

// Check if directed graph contain specific path.
// 
// unexistNodePanic flag is used to point wether or not to panic if we figure out that
// one of the nodes in path doesn't exist in graph. If unexistNodePanic is false, then
// result of the function will be false.
func ContainDirectedPath(gr DirectedGraphReader, path []VertexId, unexistNodePanic bool) bool {
	return ContainPath(NodeAndConnectionChecker(&directedNodeAndConnectionChecker{gr:gr}), path, unexistNodePanic)
}

// Check if mixed graph contain specific path.
// 
// Connection between two nodes in path "exist" in graph if and only if there is undirected connection
// between these two nodes or there is a directed connection from previous to next nodes
// in path.
// 
// unexistNodePanic flag is used to point wether or not to panic if we figure out that
// one of the nodes in path doesn't exist in graph. If unexistNodePanic is false, then
// result of the function will be false.
func ContainMixedPath(gr MixedGraphReader, path []VertexId, unexistNodePanic bool) bool {
	return ContainPath(NodeAndConnectionChecker(&mixedNodeAndConnectionChecker{gr:gr}), path, unexistNodePanic)
}
