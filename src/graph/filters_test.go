package graph

import (
	"testing"
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
)

func DirectedGraphArcsFilterSpec(c gospec.Context) {
	gr := NewDirectedMap()
	gr.AddArc(1, 2)
	gr.AddArc(2, 3)
	gr.AddArc(3, 4)
	gr.AddArc(2, 4)
	gr.AddArc(4, 5)
	gr.AddArc(1, 6)
	gr.AddArc(2, 6)

	c.Specify("Single filtered arc", func() {
		ftail := VertexId(2)
		fhead := VertexId(3)
		f := NewDirectedGraphArcFilter(gr, ftail, fhead)
		
		c.Specify("shouldn't be checked", func() {
			c.Expect(f.CheckArc(ftail, fhead), IsFalse)
		})
		
		c.Specify("shouldn't appear in accessors", func() {
			c.Expect(CollectNodes(f.GetAccessors(VertexId(ftail))), Not(Contains), fhead)
		})
		c.Specify("shouldn't appear in predecessors", func() {
			c.Expect(CollectNodes(f.GetPredecessors(VertexId(fhead))), Not(Contains), ftail)
		})
		c.Specify("shouldn't appear in iterator", func() {
			for conn := range f.ArcsIter() {
				c.Expect(conn.Tail==ftail && conn.Head==fhead, IsFalse)
			}
		})
	})
}

func UndirectedGraphEdgesFilterSpec(c gospec.Context) {
	gr := NewUndirectedMap()
	gr.AddEdge(1, 2)
	gr.AddEdge(2, 3)
	gr.AddEdge(3, 4)
	gr.AddEdge(2, 4)
	gr.AddEdge(4, 5)
	gr.AddEdge(1, 6)
	gr.AddEdge(2, 6)

	c.Specify("Single filtered arc", func() {
		ftail := VertexId(3)
		fhead := VertexId(2)
		f := NewUndirectedGraphEdgeFilter(gr, ftail, fhead)
		
		c.Specify("should be filtered", func() {
			c.Expect(f.IsEdgeFiltering(ftail, fhead), IsTrue)
			c.Expect(f.IsEdgeFiltering(fhead, ftail), IsTrue)
		})
		
		c.Specify("shouldn't be checked", func() {
			c.Expect(f.CheckEdge(ftail, fhead), IsFalse)
			c.Expect(f.CheckEdge(fhead, ftail), IsFalse)
		})
		
		c.Specify("shouldn't appear in neighbours", func() {
			c.Expect(CollectNodes(f.GetNeighbours(VertexId(ftail))), Not(Contains), fhead)
			c.Expect(CollectNodes(f.GetNeighbours(VertexId(fhead))), Not(Contains), ftail)
		})
		c.Specify("shouldn't appear in iterator", func() {
			for conn := range f.EdgesIter() {
				// iter always retur min node id as tail and max node id as head
				c.Expect(conn.Tail==fhead && conn.Head==ftail, IsFalse)
			}
		})
	})
}

func MixedGraphConnectionsFilterSpec(c gospec.Context) {
	gr := NewMixedMatrix(10)
	gr.AddArc(1, 2)
	gr.AddArc(2, 3)
	gr.AddArc(3, 4)
	gr.AddArc(2, 4)
	gr.AddArc(4, 5)
	gr.AddArc(1, 6)
	gr.AddArc(2, 6)

	c.Specify("Single filtered arc", func() {
		ftail := VertexId(2)
		fhead := VertexId(3)
		f := NewDirectedGraphArcFilter(gr, ftail, fhead)
		
		c.Specify("shouldn't be checked", func() {
			c.Expect(f.CheckArc(ftail, fhead), IsFalse)
		})
		
		c.Specify("shouldn't appear in accessors", func() {
			c.Expect(CollectNodes(f.GetAccessors(VertexId(ftail))), Not(Contains), fhead)
		})
		c.Specify("shouldn't appear in predecessors", func() {
			c.Expect(CollectNodes(f.GetPredecessors(VertexId(fhead))), Not(Contains), ftail)
		})
		c.Specify("shouldn't appear in iterator", func() {
			for conn := range f.ArcsIter() {
				c.Expect(conn.Tail==ftail && conn.Head==fhead, IsFalse)
			}
		})
	})
}
func TestGraphFilters(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(DirectedGraphArcsFilterSpec)
	r.AddSpec(UndirectedGraphEdgesFilterSpec)
	r.AddSpec(MixedGraphConnectionsFilterSpec)
	gospec.MainGoTest(r, t)
}
