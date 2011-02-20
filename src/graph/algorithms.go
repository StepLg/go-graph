package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)

// Topological sort of directed graph
//
// Return nodes in topological order. If graph has cycles, then hasCycles==true 
// and nodes==nil in function result.
func TopologicalSort(gr DirectedGraphReader) (nodes []VertexId, hasCycles bool) {
	hasCycles = false
	nodes = make([]VertexId, gr.Order())
	pos := len(nodes)
	// map of node status. If node doesn't present in map - white color,
	// node in map with false value - grey color, and with true value - black color
	status := make(map[VertexId]bool)
	for source := range gr.GetSources().VertexesIter() {
		pos, hasCycles = topologicalSortHelper(gr, source, nodes[0:pos], status)
		if hasCycles {
			nodes = nil
			return
		}
	}
	if pos!=0 {
		// cycle without path from any source to this cycle
		nodes = nil
		hasCycles = true
	}
	return
}

func topologicalSortHelper(gr DirectedGraphReader, curNode VertexId, nodes []VertexId, status map[VertexId]bool) (pos int, hasCycles bool) {
	if isBlack, ok := status[curNode]; ok {
		err := erx.NewError("Internal error in topological sort: node already in status map")
		err.AddV("node id", curNode)
		err.AddV("status in map", isBlack)
		panic(err)
	}
	hasCycles = false
	status[curNode] = false
	pos = len(nodes)
	for accessor := range gr.GetAccessors(curNode).VertexesIter() {
		if isBlack, ok := status[accessor]; ok {
			if !isBlack {
				// cycle detected!
				hasCycles = true
				return
			} else {
				// we have already visited this node
				continue
			}
		}
		pos, hasCycles = topologicalSortHelper(gr, accessor, nodes[0:pos], status)
		if hasCycles {
			return
		}
	}
	status[curNode] = true
	pos--
	nodes[pos] = curNode
	return
}

