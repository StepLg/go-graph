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

// Transform connections iterable to generic iterable object.
func ConnectionsToGenericIter(connIter ConnectionsIterable) Iterable {
	return Iterable(&connectionsIterableHelper{connIter:connIter})
}

type nodesIterableHelper struct {
	nodesIter NodesIterable
}

func (helper *nodesIterableHelper) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for node := range helper.nodesIter.NodesIter() {
			ch <- node
		}
		close(ch)
	}()
	return ch
}

type nodesGenericIterableHelper struct {
	iter Iterable
}

func (helper *nodesGenericIterableHelper) NodesIter() <-chan NodeId {
	ch := make(chan NodeId)
	go func() {
		for node := range helper.iter.Iter() {
			ch <- node.(NodeId)
		}
		close(ch)
	}()
	return ch
}

func CollectNodes(iter NodesIterable) []NodeId {
	res := make([]NodeId, 10)
	i := 0
	for node := range iter.NodesIter() {
		if i==len(res) {
			tmp := make([]NodeId, 2*i)
			copy(tmp, res)
			res = tmp
		}
		res[i] = node
		i++
	}
	
	return res[0:i]
}

// Transform nodes iterable to generic iterable object.
func NodesToGenericIter(nodesIter NodesIterable) Iterable {
	return Iterable(&nodesIterableHelper{nodesIter:nodesIter})
}

func GenericToNodesIter(iter Iterable) NodesIterable {
	return NodesIterable(&nodesGenericIterableHelper{iter:iter})
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
