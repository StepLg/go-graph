package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)

// Copy graph og to rg except args i->j, where exists non direct path i->...->j
func ReduceDirectPaths(og DirectedGraphReader, rg DirectedGraphArcsWriter, stopFunc func(from, to VertexId, weight float64) bool) {
	var checkStopFunc StopFunc
	for conn := range og.ArcsIter() {
		filteredGraph := NewDirectedGraphArcFilter(og, conn.Tail, conn.Head)
		if stopFunc!=nil {
			checkStopFunc = func(node VertexId, weight float64) bool {
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

func splitMixedGraph_helper(node VertexId, color int, gr MixedGraphReader, nodesColor map[VertexId]int) {
	nodesColor[node] = color
	// todo: neighbours and accesors as iterators
	for next := range gr.GetNeighbours(node).VertexesIter() {
		if nextColor, ok := nodesColor[next]; ok {
			if nextColor != color {
				// change all 'nextColor' nodes to 'color' nodes
				for k, v := range nodesColor {
					if v==nextColor {
						nodesColor[k] = color
					}
				}
			}
		} else {
			splitMixedGraph_helper(next, color, gr, nodesColor)
		}
	}
	for next := range gr.GetAccessors(node).VertexesIter() {
		if nextColor, ok := nodesColor[next]; ok {
			if nextColor != color {
				// change all 'nextColor' nodes to 'color' nodes
				for k, v := range nodesColor {
					if v==nextColor {
						nodesColor[k] = color
					}
				}
			}
		} else {
			splitMixedGraph_helper(next, color, gr, nodesColor)
		}
	}
	return
}

// Split mixed graph to independed subraphs
//
// @todo: Add creator function to control type of new created graphs
func SplitMixedGraph(gr MixedGraphReader) []MixedGraph {
	nodesColor := make(map[VertexId]int)
	curColor := 0
	
	for curNode := range gr.GetSources().VertexesIter() {
		if _, ok := nodesColor[curNode]; ok {
			// node already visited
			continue
		}
		splitMixedGraph_helper(curNode, curColor, gr, nodesColor)
		curColor++
	}
	
	// get total nodes count of each subgraph
	colors := make(map[int]int)
	for _, color := range nodesColor {
		if _, ok := colors[color]; !ok {
			colors[color] = 0
		}
		colors[color]++
	}
	
	result := make(map[int]MixedGraph, len(colors))
	for color, nodesCnt := range colors {
		result[color] = NewMixedMatrix(nodesCnt)
	}
	
	for node := range gr.VertexesIter() {
		result[nodesColor[node]].AddNode(node)
	}
	
	for arc := range gr.ArcsIter() {
		result[nodesColor[arc.Tail]].AddArc(arc.Tail, arc.Head)
	}
	
	for edge := range gr.EdgesIter() {
		result[nodesColor[edge.Tail]].AddEdge(edge.Tail, edge.Head)
	}
	
	resultSlice := make([]MixedGraph, len(result))
	i := 0
	for _, subgraph := range result {
		resultSlice[i] = subgraph
		i++
	}
	return resultSlice
}
