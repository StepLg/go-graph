package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)

// Arcs filter in DirectedGraphReader
//
// This is arcs filter for DirectedGraphReader. Initialize it with arcs, which need to be filtered
// and they never appeared in GetAccessors, GetPredecessors, CheckArc and Iter functions.
//
type DirectedGraphArcsFilter struct {
	DirectedGraphArcsReader
	arcs []Connection
}

// Create arcs filter with array of filtering arcs
func NewDirectedGraphArcsFilter(g DirectedGraphArcsReader, arcs []Connection) *DirectedGraphArcsFilter {
	filter := &DirectedGraphArcsFilter{
		DirectedGraphArcsReader: g,
		arcs: arcs,
	}
	return filter
}

// Create arcs filter with single arc
func NewDirectedGraphArcFilter(g DirectedGraphArcsReader, tail, head NodeId) *DirectedGraphArcsFilter {
	filter := &DirectedGraphArcsFilter{
		DirectedGraphArcsReader: g,
		arcs: make([]Connection, 1),
	}
	filter.arcs[0].Tail = tail
	filter.arcs[0].Head = head
	return filter	
}

// Getting node accessors
func (filter *DirectedGraphArcsFilter) GetAccessors(node NodeId) NodesIterable {
	iterator := func() <-chan NodeId {
		ch := make(chan NodeId)
		go func() {
			AccessorsLoop: 
			for accessor := range filter.DirectedGraphArcsReader.GetAccessors(node).NodesIter() {
				if !filter.IsArcFiltering(node, accessor) {
					ch <- accessor
				}
			}
			close(ch)
		}()
		return ch
	}
	
	return NodesIterable(&nodesIterableLambdaHelper{iterFunc:iterator})
}

// Getting node predecessors
func (filter *DirectedGraphArcsFilter) GetPredecessors(node NodeId) NodesIterable {
	iterator := func() <-chan NodeId {
		ch := make(chan NodeId)
		go func() {
			for predecessor := range filter.DirectedGraphArcsReader.GetPredecessors(node).NodesIter() {
				if !filter.IsArcFiltering(predecessor, node) {
					ch <- predecessor
				}
			}
			close(ch)
		}()
		return ch
	}
	
	return NodesIterable(&nodesIterableLambdaHelper{iterFunc:iterator})
}

// Checking arrow existance between node1 and node2
//
// node1 and node2 must exist in graph or error will be returned
func (filter *DirectedGraphArcsFilter) CheckArc(node1, node2 NodeId) bool {
	res := filter.DirectedGraphArcsReader.CheckArc(node1, node2)
	if res {
		for _, filteringConnection := range filter.arcs {
			if filteringConnection.Tail==node1 && filteringConnection.Head==node2 {
				res = false
				break
			}
		}
	}
	return res
}

func (filter *DirectedGraphArcsFilter) ArcsIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for conn := range filter.DirectedGraphArcsReader.ArcsIter() {
			if !filter.IsArcFiltering(conn.Tail, conn.Head) {
				ch <- conn
			}
		}
		close(ch)
	}()
	return ch
}

func (filter *DirectedGraphArcsFilter) IsArcFiltering(tail, head NodeId) bool {
	for _, filteringConnection := range filter.arcs {
		if filteringConnection.Head==head && filteringConnection.Tail==tail {
			return true
		}
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////

// Arcs filter in DirectedGraphReader
//
// This is arcs filter for DirectedGraphReader. Initialize it with arcs, which need to be filtered
// and they never appeared in GetAccessors, GetPredecessors, CheckArc and Iter functions.
//
type UndirectedGraphEdgesFilter struct {
	UndirectedGraphEdgesReader
	edges []Connection
}

// Create arcs filter with array of filtering arcs
func NewUndirectedGraphEdgesFilter(g UndirectedGraphEdgesReader, edges []Connection) *UndirectedGraphEdgesFilter {
	filter := &UndirectedGraphEdgesFilter{
		UndirectedGraphEdgesReader: g,
		edges: edges,
	}

	// all tails must be not greater than heads
	for i:=0; i<len(filter.edges); i++ {
		if filter.edges[i].Tail>filter.edges[i].Head {
			filter.edges[i].Tail, filter.edges[i].Head = filter.edges[i].Head, filter.edges[i].Tail
		}
	}
	return filter
}

// Create arcs filter with single arc
func NewUndirectedGraphEdgeFilter(g UndirectedGraphEdgesReader, tail, head NodeId) *UndirectedGraphEdgesFilter {
	filter := &UndirectedGraphEdgesFilter{
		UndirectedGraphEdgesReader: g,
		edges: make([]Connection, 1),
	}
	
	// tail must be not greater than head
	if tail<head {
		filter.edges[0].Tail = tail
		filter.edges[0].Head = head
	} else {
		filter.edges[0].Tail = head
		filter.edges[0].Head = tail
	}
	return filter	
}

// Getting node neighbours
func (filter *UndirectedGraphEdgesFilter) GetNeighbours(node NodeId) NodesIterable {
	iterator := func() <-chan NodeId {
		ch := make(chan NodeId)
		go func() {
			for neighbour := range filter.UndirectedGraphEdgesReader.GetNeighbours(node).NodesIter() {
				if !filter.IsEdgeFiltering(node, neighbour) {
					ch <- neighbour
				}
			}
			close(ch)
		}()
		return ch
	}
	
	return NodesIterable(&nodesIterableLambdaHelper{iterFunc:iterator})
}

// Checking edge existance between node1 and node2
//
// node1 and node2 must exist in graph or error will be returned
func (filter *UndirectedGraphEdgesFilter) CheckEdge(node1, node2 NodeId) bool {
	res := filter.UndirectedGraphEdgesReader.CheckEdge(node1, node2)
	if res {
		res = !filter.IsEdgeFiltering(node1, node2)
	}
	return res
}

func (filter *UndirectedGraphEdgesFilter) EdgesIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for conn := range filter.UndirectedGraphEdgesReader.EdgesIter() {
			if !filter.IsEdgeFiltering(conn.Tail, conn.Head) {
				ch <- conn
			}
		}
		close(ch)
	}()
	return ch
}

func (filter *UndirectedGraphEdgesFilter) IsEdgeFiltering(tail, head NodeId) bool {
	if head<tail {
		tail, head = head, tail
	}
	for _, filteringConnection := range filter.edges {
		if filteringConnection.Head==head && filteringConnection.Tail==tail {
			return true
		}
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////

// Arcs filter in MixedGraphReader
//
// This is arcs filter for MixedGraphReader.
type MixedGraphConnectionsFilter struct {
	gr MixedGraphConnectionsReader
	*DirectedGraphArcsFilter
	*UndirectedGraphEdgesFilter
}

func NewMixedGraphArcsFilter(g MixedGraphConnectionsReader, arcs []Connection, edges []Connection) *MixedGraphConnectionsFilter {
	filter := &MixedGraphConnectionsFilter{
		gr: g,
		DirectedGraphArcsFilter: NewDirectedGraphArcsFilter(g, arcs),
		UndirectedGraphEdgesFilter: NewUndirectedGraphEdgesFilter(g, edges),
	}
	return filter
}

func (filter *MixedGraphConnectionsFilter) ConnectionsIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for conn := range filter.TypedConnectionsIter() {
			ch <- conn.Connection
		}
		close(ch)
	}()
	return ch
}

func (filter *MixedGraphConnectionsFilter) TypedConnectionsIter() <-chan TypedConnection {
	ch := make(chan TypedConnection)
	go func() {
		for conn := range filter.gr.TypedConnectionsIter() {
			needToFilter := false
			switch conn.Type {
				case CT_UNDIRECTED:
					needToFilter = filter.UndirectedGraphEdgesFilter.IsEdgeFiltering(conn.Tail, conn.Head)
				case CT_DIRECTED:
					needToFilter = filter.DirectedGraphArcsFilter.IsArcFiltering(conn.Tail, conn.Head)
				default: 
					err := erx.NewError("Internal error: got unknown mixed connection type")
					panic(err)
			}
			if !needToFilter {
				ch <- conn
			}
		}
		close(ch)
	}()
	return ch
}

func (filter *MixedGraphConnectionsFilter) CheckEdgeType(tail NodeId, head NodeId) MixedConnectionType {
	res := filter.gr.CheckEdgeType(tail, head)
	if res!=CT_NONE {
		switch res {
			case CT_UNDIRECTED:
				if filter.UndirectedGraphEdgesFilter.IsEdgeFiltering(tail, head) {
					res = CT_NONE
				}
			case CT_DIRECTED:
				if filter.DirectedGraphArcsFilter.IsArcFiltering(tail, head) {
					res = CT_NONE
				}
			default: 
				err := erx.NewError("Internal error: got unknown mixed connection type")
				panic(err)
		}
	}
	return res
}
