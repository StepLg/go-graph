package graph

import (
	"testing"
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
)

func ReduceDirectPathsSpec(c gospec.Context) {
	c.Specify("Reduced triangle", func() {
		gr := NewDirectedMap()
		gr.AddArc(1, 2)
		gr.AddArc(2, 3)
		gr.AddArc(1, 3)
		
		rgr := NewDirectedMap()
		ReduceDirectPaths(gr, rgr, nil)
		
		expectedGraph := NewDirectedMap()
		expectedGraph.AddArc(1, 2)
		expectedGraph.AddArc(2, 3)
		c.Expect(rgr, DirectedGraphEquals, expectedGraph)
	})
	
	c.Specify("A bit more complex example", func() {
		gr := NewDirectedMap()
		gr.AddArc(1, 2)
		gr.AddArc(2, 3)
		gr.AddArc(3, 4)
		gr.AddArc(2, 4)
		gr.AddArc(4, 5)
		gr.AddArc(1, 6)
		gr.AddArc(2, 6)
		
		rgr := NewDirectedMap()
		ReduceDirectPaths(gr, rgr, nil)
		
		expectedGraph := NewDirectedMap()
		expectedGraph.AddArc(1, 2)
		expectedGraph.AddArc(2, 3)
		expectedGraph.AddArc(3, 4)
		expectedGraph.AddArc(4, 5)
		expectedGraph.AddArc(2, 6)
		c.Expect(rgr, DirectedGraphEquals, expectedGraph)
	})
}

func TopologicalSortSpec(c gospec.Context) {
	gr := NewDirectedMap()
	c.Specify("Single node graph", func() {
		gr.AddNode(VertexId(1))
		nodes, hasCycle := TopologicalSort(gr)
		c.Expect(hasCycle, IsFalse)
		c.Expect(nodes, ContainsExactly, Values(VertexId(1)))
	})
	
	c.Specify("Simple two nodes graph", func() {
		gr.AddArc(1, 2)
		nodes, hasCycle := TopologicalSort(gr)
		c.Expect(hasCycle, IsFalse)
		c.Expect(nodes, ContainsExactly, Values(VertexId(1), VertexId(2)))
	})
	
	c.Specify("Pseudo loops", func() {
		gr.AddArc(1, 2)
		gr.AddArc(2, 3)
		gr.AddArc(1, 4)
		gr.AddArc(4, 3)
		
		_, hasCycle := TopologicalSort(gr)
		c.Expect(hasCycle, IsFalse)
	})
}

func SplitGraphToIndependentSubgraphs_mixedSpec(c gospec.Context) {
	c.Specify("Single node graph", func() {
		subgr1 := NewMixedMatrix(3)
		subgr1.AddArc(1, 2)
		subgr1.AddEdge(1, 3)
		subgr2 := NewMixedMatrix(3)
		subgr2.AddArc(4, 5)
		subgr2.AddArc(5, 6)
		subgr2.AddArc(4, 6)
		
		gr := NewMixedMatrix(6)
		CopyMixedGraph(subgr1, gr)
		CopyMixedGraph(subgr2, gr)
		
		subgraphs := SplitGraphToIndependentSubgraphs_mixed(gr)
		c.Expect(len(subgraphs), Equals, 2)
		// @todo: Add ContainsGraph comparator to check if slice contains a graph
		c.Expect(MixedGraphsEquals(subgr1, subgraphs[1]), IsTrue)
		c.Expect(MixedGraphsEquals(subgr2, subgraphs[0]), IsTrue)
	})
}

func SplitGraphToIndependentSubgraphs_directedSpec(c gospec.Context) {
	c.Specify("Directed graph with 2 independent parts", func() {
		gr1, gr2, gr_merged := genDgr2IndependentSubGr()
		subgraphs := SplitGraphToIndependentSubgraphs_directed(gr_merged)
		c.Expect(len(subgraphs), Equals, 2)
		if subgraphs[0].CheckNode(VertexId(1)) {
			c.Expect(DirectedGraphsEquals(subgraphs[0], gr1), IsTrue)
			c.Expect(DirectedGraphsEquals(subgraphs[1], gr2), IsTrue)
		} else {
			c.Expect(DirectedGraphsEquals(subgraphs[0], gr2), IsTrue)
			c.Expect(DirectedGraphsEquals(subgraphs[1], gr1), IsTrue)
		}
	})
}

func SplitGraphToIndependentSubgraphs_undirectedSpec(c gospec.Context) {
	c.Specify("Undirected graph with 2 independent parts", func() {
		gr1, gr2, gr_merged := genUgr2IndependentSubGr()
		subgraphs := SplitGraphToIndependentSubgraphs_undirected(gr_merged)
		c.Expect(len(subgraphs), Equals, 2)
		if subgraphs[0].CheckNode(VertexId(1)) {
			c.Expect(UndirectedGraphsEquals(subgraphs[0], gr1), IsTrue)
			c.Expect(UndirectedGraphsEquals(subgraphs[1], gr2), IsTrue)
		} else {
			c.Expect(UndirectedGraphsEquals(subgraphs[0], gr2), IsTrue)
			c.Expect(UndirectedGraphsEquals(subgraphs[1], gr1), IsTrue)
		}
	})
}


func TestAlgorithms(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(ReduceDirectPathsSpec)
	r.AddSpec(TopologicalSortSpec)
	r.AddSpec(SplitGraphToIndependentSubgraphs_mixedSpec)
	r.AddSpec(SplitGraphToIndependentSubgraphs_directedSpec)
	r.AddSpec(SplitGraphToIndependentSubgraphs_undirectedSpec)
	gospec.MainGoTest(r, t)
}
