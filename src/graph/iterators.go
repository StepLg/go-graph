package graph

import (
	. "exp/iterable"

	"github.com/StepLg/go-erx/src/erx"
)

type connectionsIterableHelper struct {
	connIter ConnectionsIterable
}

func (helper *connectionsIterableHelper) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for arr := range helper.connIter.ConnectionsIter() {
			ch <- arr
		}
		close(ch)
	}()
	return ch
}

type connectionsGenericIterableHelper struct {
	iter Iterable
}

func (helper *connectionsGenericIterableHelper) ConnectionsIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for arr := range helper.iter.Iter() {
			ch <- arr.(Connection)
		}
		close(ch)
	}()
	return ch
}

// Transform connections iterable to generic iterable object.
func ConnectionsToGenericIter(connIter ConnectionsIterable) Iterable {
	return Iterable(&connectionsIterableHelper{connIter:connIter})
}

// Transform generic iterable to connections iterable object.
func GenericToConnectionsIter(iter Iterable) ConnectionsIterable {
	return ConnectionsIterable(&connectionsGenericIterableHelper{iter:iter})
}

type nodesIterableHelper struct {
	nodesIter VertexesIterable
}

func (helper *nodesIterableHelper) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for node := range helper.nodesIter.VertexesIter() {
			ch <- node
		}
		close(ch)
	}()
	return ch
}

type nodesGenericIterableHelper struct {
	iter Iterable
}

func (helper *nodesGenericIterableHelper) VertexesIter() <-chan VertexId {
	ch := make(chan VertexId)
	go func() {
		for node := range helper.iter.Iter() {
			ch <- node.(VertexId)
		}
		close(ch)
	}()
	return ch
}

// Transform vertexes iterable to generic iterable object.
func VertexesToGenericIter(nodesIter VertexesIterable) Iterable {
	return Iterable(&nodesIterableHelper{nodesIter:nodesIter})
}

// Transform generic iterator to vertexes iterable
func GenericToVertexesIter(iter Iterable) VertexesIterable {
	return VertexesIterable(&nodesGenericIterableHelper{iter:iter})
}

// Collect all vertexes from iterator to slice.
func CollectVertexes(iter VertexesIterable) []VertexId {
	res := make([]VertexId, 10)
	i := 0
	for node := range iter.VertexesIter() {
		if i==len(res) {
			tmp := make([]VertexId, 2*i)
			copy(tmp, res)
			res = tmp
		}
		res[i] = node
		i++
	}
	
	return res[0:i]
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

// Copy all arcs from iterator to directed graph
//
// todo: merge with CopyUndirectedGraph
func CopyDirectedGraph(connIter ConnectionsIterable, gr DirectedGraphArcsWriter) {
	// wheel := erx.NewError("Can't copy directed graph")
	for arrow := range connIter.ConnectionsIter() {
		gr.AddArc(arrow.Tail, arrow.Head)
	}
	return
}

// Copy all arcs from iterator to directed graph
//
// todo: add VertexesIterable interface and copy all nodes before copying arcs
func CopyUndirectedGraph(connIter ConnectionsIterable, gr UndirectedGraphEdgesWriter) {
	// wheel := erx.NewError("Can't copy directed graph")
	for arrow := range connIter.ConnectionsIter() {
		gr.AddEdge(arrow.Tail, arrow.Head)
	}
	return
}

// Copy all connections from iterator to mixed graph
//
// todo: merge with CopyDirectedGraph
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

// helper struct for ArcsToConnIterable
type arcsToConnIterable_helper struct {
	gr DirectedGraphArcsReader
}

func (helper *arcsToConnIterable_helper) ConnectionsIter() <-chan Connection {
	return helper.gr.ArcsIter()
}

// Convert arcs iterator to connections iterator.
func ArcsToConnIterable(gr DirectedGraphArcsReader) ConnectionsIterable {
	return &arcsToConnIterable_helper{gr}
}

// helper struct for EdgesToConnIterable
type edgesToConnIterable_helper struct {
	gr UndirectedGraphEdgesReader
}

func (helper *edgesToConnIterable_helper) ConnectionsIter() <-chan Connection {
	return helper.gr.EdgesIter()
}

// Convert edges iterator to connections iterator.
func EdgesToConnIterable(gr UndirectedGraphEdgesReader) ConnectionsIterable {
	return &edgesToConnIterable_helper{gr}
}

// helper struct for ArcsToTypedConnIterable
type arcsToTypedConnIterable_helper struct {
	gr DirectedGraphArcsReader
}

func (helper *arcsToTypedConnIterable_helper) TypedConnectionsIter() <-chan TypedConnection {
	ch := make(chan TypedConnection)
	go func() {
		for conn := range helper.gr.ArcsIter() {
			ch <- TypedConnection{Connection: conn, Type: CT_DIRECTED}
		}
	}()
	return ch
}

// Convert arcs iterator to typed connections iterator.
func ArcsToTypedConnIterable(gr DirectedGraphArcsReader) TypedConnectionsIterable {
	return &arcsToTypedConnIterable_helper{gr}
}

// helper struct for EdgesToTypedConnIterable
type edgesToTypedConnIterable_helper struct {
	gr UndirectedGraphEdgesReader
}

func (helper *edgesToTypedConnIterable_helper) TypedConnectionsIter() <-chan TypedConnection {
	ch := make(chan TypedConnection)
	go func() {
		for conn := range helper.gr.EdgesIter() {
			ch <- TypedConnection{Connection: conn, Type: CT_UNDIRECTED}
		}
	}()
	return ch
}

// Convert edges iterator to typed connections iterator.
func EdgesToTypedConnIterable(gr UndirectedGraphEdgesReader) TypedConnectionsIterable {
	return &edgesToTypedConnIterable_helper{gr}
}


// Vertexes iterable object for function, which returns VertexId channel.
//
// Internal use only at this moment. Don't know what for it could be used
// outside.
type nodesIterableLambdaHelper struct {
	iterFunc func() <-chan VertexId
}

func (helper *nodesIterableLambdaHelper) VertexesIter() <-chan VertexId {
	return helper.iterFunc()
}
