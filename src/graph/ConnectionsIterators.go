package graph

import (
	. "exp/iterable"

	"github.com/StepLg/go-erx/src/erx"
)

type connectionsIterable struct {
	arrows ConnectionsIterable
}

func (ai connectionsIterable) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for arr := range ai.arrows.ConnectionsIter() {
			ch <- arr
		}
	}()
	return ch
}

func ArrowsToGenericIter(connIter ConnectionsIterable) Iterable {
	return connectionsIterable{connIter}
}

// Copy all arcs from iterator to directed graph
//
// todo: add NodesIterable interface and copy all nodes before copying arcs
func CopyDirectedGraph(connIter ConnectionsIterable, gr DirectedGraphArcsWriter) {
	// wheel := erx.NewError("Can't copy directed graph")
	for arrow := range connIter.ConnectionsIter() {
		gr.AddArc(arrow.Tail, arrow.Head)
	}
	return
}

// Build directed graph from connecection iterator with order function
//
// For all connections from iterator check isCorrectOrder function 
// and add to directed graph connection in correct order
func BuildDirectedGraph(gr DirectedGraph, connIterable ConnectionsIterable , isCorrectOrder func(Connection) bool) {
	for arr := range connIterable.ConnectionsIter() {
		if isCorrectOrder(arr) {
			gr.AddArc(arr.Tail, arr.Head)
		} else {
			gr.AddArc(arr.Head, arr.Tail)
		}
	}
}

// Copy all connections from iterator to mixed graph
//
// todo: add NodesIterable interface and copy all nodes before copying connections
func CopyMixedGraph(from TypedConnectionsIterable, to MixedGraphWriter) {
	for conn := range from.TypedConnectionsIter() {
		switch conn.Type {
			case CT_UNDIRECTED:
				to.AddEdge(conn.Tail, conn.Head)
			case CT_DIRECTED:
				to.AddArc(conn.Tail, conn.Head)
			default:
				err := erx.NewError("Internal error: unknown connection type")
				panic(err)
		}
	}
}
