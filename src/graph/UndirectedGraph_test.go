package graph

import (
	"testing"
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
)

func UndirectedGraphSpec(c gospec.Context, graphCreator func() UndirectedGraph) {
	gr := graphCreator()

	c.Specify("Empty undirected graph", func() {
		c.Specify("contain no nodes", func() {
			c.Expect(gr.VertexesCnt(), Equals, 0)
		})
		c.Specify("contain no edges", func() {
			c.Expect(gr.EdgesCnt(), Equals, 0)
		})
	})

	
	c.Specify("New node in empty graph", func() {
		VertexId := VertexId(1)
		gr.AddNode(VertexId)
				
		c.Specify("changing nodes count", func() {
			c.Expect(gr.VertexesCnt(), Equals, 1)
		})
		
		c.Specify("doesn't change edges count", func() {
			c.Expect(gr.EdgesCnt(), Equals, 0)
		})
		
		c.Specify("no neighbours", func() {
			c.Expect(len(CollectVertexes(gr.GetNeighbours(VertexId))), Equals, 0)
		})
	})

	c.Specify("New edge in empty graph", func() {
		n1 := VertexId(1)
		n2 := VertexId(2)
		gr.AddEdge(n1, n2)
				
		c.Specify("changing nodes count", func() {
			c.Expect(gr.VertexesCnt(), Equals, 2)
		})
		
		c.Specify("changing edges count", func() {
			c.Expect(gr.EdgesCnt(), Equals, 1)
		})
		
		c.Specify("neighbours", func() {
			c.Expect(CollectVertexes(gr.GetNeighbours(n1)), ContainsExactly, Values(n2))
			c.Expect(CollectVertexes(gr.GetNeighbours(n2)), ContainsExactly, Values(n1))
		})
	})
}

func TestUndirectedGraphSpec(t *testing.T) {
	r := gospec.NewRunner()
	
	// paramenerized test creator
	cr := func(graphCreator func() UndirectedGraph) func (c gospec.Context) {
		return func(c gospec.Context){
			UndirectedGraphSpec(c, graphCreator)
		}
	}
	
	r.AddNamedSpec("UndirectedGraph(Map)", cr(func() UndirectedGraph {
		return UndirectedGraph(NewUndirectedMap())
	}))
	r.AddNamedSpec("UndirectedGraph(Matrix)", cr(func() UndirectedGraph {
		return UndirectedGraph(NewUndirectedMatrix(10))
	}))
	r.AddNamedSpec("UndirectedGraph(MixedMatrix)", cr(func() UndirectedGraph {
		return UndirectedGraph(NewMixedMatrix(10))
	}))
	r.AddNamedSpec("UndirectedGraph(MixedMap)", cr(func() UndirectedGraph {
		return UndirectedGraph(NewMixedMap())
	}))
	gospec.MainGoTest(r, t)
}
