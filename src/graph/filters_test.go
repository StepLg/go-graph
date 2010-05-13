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
		ftail := NodeId(2)
		fhead := NodeId(3)
		f := NewDirectedGraphArcFilter(gr, ftail, fhead)
		
		c.Specify("shouldn't be checked", func() {
			c.Expect(f.CheckArc(ftail, fhead), IsFalse)
		})
		
		c.Specify("shouldn't appear in accessors", func() {
			c.Expect(f.GetAccessors(NodeId(ftail)), Not(Contains), fhead)
		})
		c.Specify("shouldn't appear in predecessors", func() {
			c.Expect(f.GetPredecessors(NodeId(fhead)), Not(Contains), ftail)
		})
		c.Specify("shouldn't appear in iterator", func() {
			for conn := range f.ConnectionsIter() {
				c.Expect(conn.Tail==ftail && conn.Head==fhead, IsFalse)
			}
		})
	})
}

func TestGraphFilters(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(DirectedGraphArcsFilterSpec)
	gospec.MainGoTest(r, t)
}
