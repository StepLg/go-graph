package graph

import (
	"testing"
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
)

func ComparatorsSpec(c gospec.Context) {
	gr := NewMixedMatrix(10)
	gr.AddArc(1, 2)
	gr.AddArc(2, 3)
	gr.AddArc(3, 4)
	gr.AddArc(2, 4)
	gr.AddArc(4, 5)
	gr.AddArc(1, 6)
	gr.AddArc(2, 6)
	gr.AddEdge(3, 6)
	gr.AddEdge(2, 7)
	gr.AddArc(7, 4)

	c.Specify("Graph copy", func() {
		grcopy := NewMixedMatrix(gr.NodesCnt())
		CopyMixedGraph(gr, grcopy)

		c.Specify("includes must be true in both ways", func() {
			c.Expect(GraphIncludeNodes(gr, grcopy), IsTrue)
			c.Expect(GraphIncludeNodes(grcopy, gr), IsTrue)

			c.Expect(MixedGraphIncludeConnections(gr, grcopy), IsTrue)
			c.Expect(MixedGraphIncludeConnections(grcopy, gr), IsTrue)
		})
		
		c.Specify("must be equal to original", func() {
			c.Expect(MixedGraphsEquals(gr, grcopy), IsTrue)
			c.Expect(MixedGraphsEquals(grcopy, gr), IsTrue)
		})
	})
	
	c.Specify("Graph copy with additional connection", func() {
		grcopy := NewMixedMatrix(gr.NodesCnt())
		CopyMixedGraph(gr, grcopy)
		grcopy.AddEdge(4, 6)
		
		c.Specify("includes original as a subgraph", func() {
			c.Expect(GraphIncludeNodes(gr, grcopy), IsTrue)
			c.Expect(GraphIncludeNodes(grcopy, gr), IsTrue)

			c.Expect(MixedGraphIncludeConnections(gr, grcopy), IsFalse)
			c.Expect(MixedGraphIncludeConnections(grcopy, gr), IsTrue)
		})
		
		c.Specify("must not be equal to original", func() {
			c.Expect(MixedGraphsEquals(gr, grcopy), IsFalse)
			c.Expect(MixedGraphsEquals(grcopy, gr), IsFalse)
		})
	})
}

func TestComparators(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(ComparatorsSpec)
	gospec.MainGoTest(r, t)
}
