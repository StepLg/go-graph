package graph

import (
	"testing"
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
)

func MixedGraphSpec(c gospec.Context, graphCreator func() MixedGraph) {
	gr := graphCreator()
	c.Specify("After adding new edge", func() {
		tail := VertexId(1)
		head := VertexId(2)
		gr.AddEdge(tail, head)
		c.Specify("contain exactly two nodes", func() {
			c.Expect(gr.Order(), Equals, 2)
		})
		c.Specify("contain single edge", func() {
			c.Expect(gr.EdgesCnt(), Equals, 1)
		})
		c.Specify("contain no arcs", func() {
			c.Expect(gr.ArcsCnt(), Equals, 0)
		})
		c.Specify("has one connection with type 'undirected'", func() {
			c.Expect(gr.CheckEdge(tail, head), IsTrue)
			c.Expect(gr.CheckEdge(head, tail), IsTrue)
			c.Expect(gr.CheckArc(tail, head), IsFalse)
			c.Expect(gr.CheckArc(head, tail), IsFalse)
			c.Expect(gr.CheckEdgeType(tail, head), Equals, CT_UNDIRECTED)
			c.Expect(gr.CheckEdgeType(head, tail), Equals, CT_UNDIRECTED)
		})
	})
	c.Specify("After adding new arc", func() {
		tail := VertexId(1)
		head := VertexId(2)
		gr.AddArc(tail, head)
		c.Specify("contain exactly two nodes", func() {
			c.Expect(gr.Order(), Equals, 2)
		})
		c.Specify("contain no edges", func() {
			c.Expect(gr.EdgesCnt(), Equals, 0)
		})
		c.Specify("contain single arc", func() {
			c.Expect(gr.ArcsCnt(), Equals, 1)
		})
		c.Specify("has one connection with type 'directed'", func() {
			c.Expect(gr.CheckEdge(tail, head), IsFalse)
			c.Expect(gr.CheckEdge(head, tail), IsFalse)
			c.Expect(gr.CheckArc(tail, head), IsTrue)
			c.Expect(gr.CheckArc(head, tail), IsFalse)
			c.Expect(gr.CheckEdgeType(tail, head), Equals, CT_DIRECTED)
			c.Expect(gr.CheckEdgeType(head, tail), Equals, CT_DIRECTED_REVERSED)
		})
	})
}

func TestMixedGraphSpec(t *testing.T) {
	r := gospec.NewRunner()
	
	// paramenerized test creator
	cr := func(graphCreator func() MixedGraph) func (c gospec.Context) {
		return func(c gospec.Context){
			MixedGraphSpec(c, graphCreator)
		}
	}
	
	r.AddNamedSpec("MixedGraph(MixedMap)", cr(func() MixedGraph {
		return MixedGraph(NewMixedMap())
	}))
	r.AddNamedSpec("MixedGraph(MixedMatrix)", cr(func() MixedGraph {
		return MixedGraph(NewMixedMatrix(10))
	}))
	
	gospec.MainGoTest(r, t)
}
