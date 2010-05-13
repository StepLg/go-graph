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
func (filter *DirectedGraphArcsFilter) GetAccessors(node NodeId) Nodes {
	accessors := filter.DirectedGraphArcsReader.GetAccessors(node)
	newAccessorsLen := len(accessors)
	for _, filteringConnection := range filter.arcs {
		if node == filteringConnection.Tail {
			// need to remove filtering arc
			k := 0
			for k=0; k<newAccessorsLen; k++ {
				if accessors[k]==filteringConnection.Head {
					break
				}
			}
			if k<newAccessorsLen {
				copy(accessors[k:newAccessorsLen-1], accessors[k+1:newAccessorsLen])
				newAccessorsLen--
			}
		}
	}
	return accessors[0:newAccessorsLen]
}

// Getting node predecessors
func (filter *DirectedGraphArcsFilter) GetPredecessors(node NodeId) Nodes {
	predecessors := filter.DirectedGraphArcsReader.GetAccessors(node)
	newPredecessorsLen := len(predecessors)
	for _, filteringConnection := range filter.arcs {
		if node == filteringConnection.Head {
			// need to remove filtering arc
			k := 0
			for k=0; k<newPredecessorsLen; k++ {
				if predecessors[k]==filteringConnection.Tail {
					break
				}
			}
			if k<newPredecessorsLen {
				copy(predecessors[k:newPredecessorsLen-1], predecessors[k+1:newPredecessorsLen])
				newPredecessorsLen--
			}
		}
	}
	return predecessors[0:newPredecessorsLen]
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
			if filter.IsArcFiltering(conn.Tail, conn.Head) {
				continue
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
	for _, conn := range edges {
		if conn.Tail>conn.Head {
			conn.Tail, conn.Head = conn.Head, conn.Tail
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
func (filter *UndirectedGraphEdgesFilter) GetNeighbours(node NodeId) Nodes {
	neighbours := filter.UndirectedGraphEdgesReader.GetNeighbours(node)
	newNeighboursLen := len(neighbours)
	for _, filteringConnection := range filter.edges {
		if node == filteringConnection.Tail {
			// need to remove filtering edge
			k := 0
			for k=0; k<newNeighboursLen; k++ {
				if neighbours[k]==filteringConnection.Head {
					break
				}
			}
			if k<newNeighboursLen {
				copy(neighbours[k:newNeighboursLen-1], neighbours[k+1:newNeighboursLen])
				newNeighboursLen--
			}
		}
	}
	return neighbours[0:newNeighboursLen]
}

// Checking edge existance between node1 and node2
//
// node1 and node2 must exist in graph or error will be returned
func (filter *UndirectedGraphEdgesFilter) CheckEdge(node1, node2 NodeId) bool {
	res := filter.UndirectedGraphEdgesReader.CheckEdge(node1, node2)
	if res {
		for _, filteringConnection := range filter.edges {
			if filteringConnection.Tail==node1 && filteringConnection.Head==node2 {
				res = false
				break
			}
		}
	}
	return res
}

func (filter *UndirectedGraphEdgesFilter) EdgesIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for conn := range filter.UndirectedGraphEdgesReader.EdgesIter() {
			if filter.IsEdgeFiltering(conn.Tail, conn.Head) {
				continue
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
