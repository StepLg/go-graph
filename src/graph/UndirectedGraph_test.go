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
			c.Expect(gr.NodesCnt(), Equals, 0)
		})
		c.Specify("contain no edges", func() {
			c.Expect(gr.EdgesCnt(), Equals, 0)
		})
	})

	
	c.Specify("New node in empty graph", func() {
		nodeId := NodeId(1)
		err := gr.AddNode(nodeId)
		c.Expect(err, IsNil)
				
		c.Specify("changing nodes count", func() {
			c.Expect(gr.NodesCnt(), Equals, 1)
		})
		
		c.Specify("doesn't change arrows count", func() {
			c.Expect(gr.EdgesCnt(), Equals, 0)
		})
		
		c.Specify("no neighbours", func() {
			accessors, err := gr.GetNeighbours(nodeId)
			c.Expect(err, IsNil)
			c.Expect(len(accessors), Equals, 0)
		})
	})

	c.Specify("New edge in empty graph", func() {
		n1 := NodeId(1)
		n2 := NodeId(2)
		c.Expect(gr.AddEdge(n1, n2), IsNil)
				
		c.Specify("changing nodes count", func() {
			c.Expect(gr.NodesCnt(), Equals, 2)
		})
		
		c.Specify("changing edges count", func() {
			c.Expect(gr.EdgesCnt(), Equals, 1)
		})
		
		c.Specify("neighbours", func() {
			neighbours, err := gr.GetNeighbours(n1)
			c.Expect(err, IsNil)
			c.Expect(neighbours, ContainsExactly, Values(n2))

			neighbours, err = gr.GetNeighbours(n2)
			c.Expect(err, IsNil)
			c.Expect(neighbours, ContainsExactly, Values(n1))

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
	gospec.MainGoTest(r, t)
}
