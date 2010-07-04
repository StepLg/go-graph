package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)

// Copy graph og to rg except args i->j, where exists non direct path i->...->j
//
// Graph rg contains all vertexes from original graph gr and arcs i->j, if there
// doesn't exist path in original graph from i to j, which contains at least
// 3 vertexes
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

// Split mixed graph to independed subraphs.
//
// Each result subgraph contain only those vertexes, which are connected, and
// which are disjoint with any vertex from any other subgraph.
//
// @todo: Add creator function to control type of new created graphs.
// @todo: Here could be done lot's of optimisation: recolor the smallest color
// and parallelisation.
func SplitGraphToIndependentSubgraphs_mixed(gr MixedGraphReader) []MixedGraph {
	// vertexes color map. Two vertexes have same color if and only if there exist
	// any path (which could contain both arcs and edges) from one to another
	nodesColor := make(map[VertexId]int)
	curColor := 0
	
	// coloring vertexes
	for curNode := range gr.GetSources().VertexesIter() {
		if _, ok := nodesColor[curNode]; ok {
			// node already visited
			continue
		}
		splitGraphToIndependentSubgraphs_helper(curNode, curColor, NewMgraphOutNeighboursExtractor(gr), nodesColor)
		curColor++
	}
	
	// get total vertexes number of each subgraph
	colors := make(map[int]int)
	for _, color := range nodesColor {
		if _, ok := colors[color]; !ok {
			colors[color] = 0
		}
		colors[color]++
	}

	// making subgraphs map with color as a key	
	result := make(map[int]MixedGraph, len(colors))
	for color, nodesCnt := range colors {
		result[color] = NewMixedMatrix(nodesCnt)
	}
	
	// copying nodes to subgraphs
	for node := range gr.VertexesIter() {
		result[nodesColor[node]].AddNode(node)
	}
	
	// copying arcs to subgraphs
	for arc := range gr.ArcsIter() {
		result[nodesColor[arc.Tail]].AddArc(arc.Tail, arc.Head)
	}
	
	// copying edges to subgraphs
	for edge := range gr.EdgesIter() {
		result[nodesColor[edge.Tail]].AddEdge(edge.Tail, edge.Head)
	}
	
	// making result slice
	resultSlice := make([]MixedGraph, len(result))
	i := 0
	for _, subgraph := range result {
		resultSlice[i] = subgraph
		i++
	}
	return resultSlice
}

// Split directed graph to independed subraphs.
//
// Each result subgraph contain only those vertexes, which are connected, and
// which are disjoint with any vertex from any other subgraph.
//
// @todo: Add creator function to control type of new created graphs
// @todo: Here could be done lot's of optimisation: recolor the smallest color
// and parallelisation.
func SplitGraphToIndependentSubgraphs_directed(gr DirectedGraphReader) []DirectedGraph {
	// vertexes color map. Two vertexes have same color if and only if there exist
	// any path (which could contain both arcs and edges) from one to another
	nodesColor := make(map[VertexId]int)
	curColor := 0
	
	// coloring vertexes
	for curNode := range gr.GetSources().VertexesIter() {
		splitGraphToIndependentSubgraphs_helper(curNode, curColor, NewDgraphOutNeighboursExtractor(gr), nodesColor)
		curColor++
	}
	
	// copying nodes to subgraphs
	result := make(map[int]DirectedGraph)
	for node := range gr.VertexesIter() {
		var subgr DirectedGraph
		var ok bool
		if subgr, ok = result[nodesColor[node]]; !ok {
			subgr = NewDirectedMap()
			result[nodesColor[node]] = subgr
		}
		subgr.AddNode(node)
	}
	
	// copying arcs to subgraphs
	for arc := range gr.ArcsIter() {
		result[nodesColor[arc.Tail]].AddArc(arc.Tail, arc.Head)
	}
	
	// making result slice
	resultSlice := make([]DirectedGraph, len(result))
	i := 0
	for _, subgraph := range result {
		resultSlice[i] = subgraph
		i++
	}
	return resultSlice
}

// Split undirected graph to independed subraphs.
//
// Each result subgraph contain only those vertexes, which are connected, and
// which are disjoint with any vertex from any other subgraph.
//
// @todo: Add creator function to control type of new created graphs
// @todo: Here could be done lot's of optimisation: recolor the smallest color
// and parallelisation.
func SplitGraphToIndependentSubgraphs_undirected(gr UndirectedGraphReader) []UndirectedGraph {
	// vertexes color map. Two vertexes have same color if and only if there exist
	// any path (which could contain both arcs and edges) from one to another
	nodesColor := make(map[VertexId]int)
	curColor := 0
	
	// coloring vertexes
	for curNode := range gr.VertexesIter() {
		if _, ok := nodesColor[curNode]; ok {
			// node already visited
			continue
		}
		splitGraphToIndependentSubgraphs_helper(curNode, curColor, NewUgraphOutNeighboursExtractor(gr), nodesColor)
		curColor++
	}
	
	// copying nodes to subgraphs
	result := make(map[int]UndirectedGraph)
	for node := range gr.VertexesIter() {
		var subgr UndirectedGraph
		var ok bool
		if subgr, ok = result[nodesColor[node]]; !ok {
			subgr = NewUndirectedMap()
			result[nodesColor[node]] = subgr
		}
		subgr.AddNode(node)
	}
	
	// copying edges to subgraphs
	for edge := range gr.EdgesIter() {
		result[nodesColor[edge.Tail]].AddEdge(edge.Tail, edge.Head)
	}
	
	// making result slice
	resultSlice := make([]UndirectedGraph, len(result))
	i := 0
	for _, subgraph := range result {
		resultSlice[i] = subgraph
		i++
	}
	return resultSlice
}

// Helper function for splitting graph to independent subgraphs.
//
// Function assign color to each vertex, accessible from given one. If some
// vertex already has different color, then all vertexes with this calor are
// changing color to new.
func splitGraphToIndependentSubgraphs_helper(node VertexId, color int, gr OutNeighboursExtractor, nodesColor map[VertexId]int) {
	nodesColor[node] = color
	for next := range gr.GetOutNeighbours(node).VertexesIter() {
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
			splitGraphToIndependentSubgraphs_helper(next, color, gr, nodesColor)
		}
	}
	return
}
