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

func (filter *DirectedGraphArcsFilter) ConnectionsIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for conn := range filter.DirectedGraphArcsReader.ArcsIter() {
			for _, filteringConnection := range filter.arcs {
				if filteringConnection.Head==conn.Head && filteringConnection.Tail==conn.Tail {
					continue
				}
			}
		}
		close(ch)
	}()
	return ch
}

///////////////////////////////////////////////////////////////////////////////
/*
// Arcs filter in MixedGraphReader
//
// This is arcs filter for MixedGraphReader. Initialize it with arcs, which need to be filtered
// and they never appeared in GetAccessors, GetPredecessors, CheckArc and Iter functions.
//
// Be careful! Filter doesn't affect GetSources and GetSinks functions. Also it doesn't recalculate
// dangling vertexes.
type MixedGraphConnectionsFilter struct {
	*DirectedGraphArcsFilter	
}

func NewMixedGraphArcsFilter(g DirectedGraphReader, arcs []Connection, edges []Connection) *DirectedGraphArcsFilter {
	filter := &DirectedGraphArcsFilter{
		DirectedGraphArcsFilter: NewDirectedGraphArcsFilter(g, arcs),
	}
	return filter
}
*/