package graph

import (
	"testing"
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
)

func CheckPathSpec(c gospec.Context, checkPathFunction CheckPath) {
	gr := NewDirectedMap()
	gr.AddArc(1, 2)
	gr.AddArc(2, 3)
	gr.AddArc(3, 4)
	gr.AddArc(2, 4)
	gr.AddArc(4, 5)
	gr.AddArc(1, 6)
	gr.AddArc(2, 6)
	
	c.Specify("Check path to self", func() {
		c.Expect(checkPathFunction(gr, 1, 1, -1, SimpleWeightFunc), IsTrue) 
		c.Expect(checkPathFunction(gr, 6, 6, -1, SimpleWeightFunc), IsTrue)
	})
	
	c.Specify("Check neighbours path", func() {
		c.Expect(checkPathFunction(gr, 1, 2, -1, SimpleWeightFunc), IsTrue)
		c.Expect(checkPathFunction(gr, 2, 4, -1, SimpleWeightFunc), IsTrue)
		c.Expect(checkPathFunction(gr, 1, 6, -1, SimpleWeightFunc), IsTrue)
	})
	
	c.Specify("Check reversed neighbours", func() {
		c.Expect(checkPathFunction(gr, 6, 1, -1, SimpleWeightFunc), IsFalse)
		c.Expect(checkPathFunction(gr, 4, 3, -1, SimpleWeightFunc), IsFalse)
		c.Expect(checkPathFunction(gr, 5, 4, -1, SimpleWeightFunc), IsFalse)
	})
	
	c.Specify("Check long path", func() {
		c.Expect(checkPathFunction(gr, 1, 6, -1, SimpleWeightFunc), IsTrue)
		c.Expect(checkPathFunction(gr, 1, 5, -1, SimpleWeightFunc), IsTrue)
	})
	
	c.Specify("Check weight limit", func() {
		c.Expect(checkPathFunction(gr, 1, 5, 2, SimpleWeightFunc), IsFalse)
	})
}

func TestSearch(t *testing.T) {
	r := gospec.NewRunner()

	// paramenerized test creator
	cr := func(checkPathFunction CheckPath) func (c gospec.Context) {
		return func(c gospec.Context){
			CheckPathSpec(c, checkPathFunction)
		}
	}
	r.AddNamedSpec("CheckPath(Dijkstra)", cr(CheckPathDijkstra))

	gospec.MainGoTest(r, t)
}
