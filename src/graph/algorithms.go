package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)

// Copy graph og to rg except args i->j, where exists non direct path i->...->j
func ReduceDirectPaths(og DirectedGraphReader, rg DirectedGraphArcsWriter, stopFunc func(from, to NodeId, weight float) bool) {
	var checkStopFunc StopFunc
	for conn := range og.ArcsIter() {
		filteredGraph := NewDirectedGraphArcFilter(og, conn.Tail, conn.Head)
		if stopFunc!=nil {
			checkStopFunc = func(node NodeId, weight float) bool {
				return stopFunc(conn.Tail, node, weight)
			}
		} else {
			checkStopFunc = nil
		}
		if !CheckDirectedPathDijkstra(filteredGraph, conn.Tail, conn.Head, checkStopFunc, SimpleWeightFunc) {
			rg.AddArc(conn.Tail, conn.Head)
		}
	}
}

func TopologicalSort(gr DirectedGraphReader) (nodes []NodeId, hasCycles bool) {
	hasCycles = false
	nodes = make([]NodeId, gr.NodesCnt())
	pos := len(nodes)
	// map of node status. If node doesn't present in map - white color,
	// node in map with false value - grey color, and with true value - black color
	status := make(map[NodeId]bool)
	for _, source := range gr.GetSources() {
		pos, hasCycles = topologicalSortHelper(gr, source, nodes[0:pos], status)
		if hasCycles {
			nodes = nil
			return
		}
	}
	return
}

func topologicalSortHelper(gr DirectedGraphReader, curNode NodeId, nodes []NodeId, status map[NodeId]bool) (pos int, hasCycles bool) {
	if isBlack, ok := status[curNode]; ok {
		err := erx.NewError("Internal error in topological sort: node already in status map")
		err.AddV("node id", curNode)
		err.AddV("status in map", isBlack)
		panic(err)
	}
	hasCycles = false
	status[curNode] = false
	pos = len(nodes)
	for _, accessor := range gr.GetAccessors(curNode) {
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
