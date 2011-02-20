package graph

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
func NewDirectedGraphArcFilter(g DirectedGraphArcsReader, tail, head VertexId) *DirectedGraphArcsFilter {
	filter := &DirectedGraphArcsFilter{
		DirectedGraphArcsReader: g,
		arcs: make([]Connection, 1),
	}
	filter.arcs[0].Tail = tail
	filter.arcs[0].Head = head
	return filter	
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
func NewUndirectedGraphEdgeFilter(g UndirectedGraphEdgesReader, tail, head VertexId) *UndirectedGraphEdgesFilter {
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

