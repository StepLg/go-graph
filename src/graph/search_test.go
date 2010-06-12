package graph

import (
	"testing"
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
)

func generateMixedGraph1() MixedGraph {
	gr := NewMixedMatrix(6)
	gr.AddArc(1, 2)
	gr.AddArc(2, 3)
	gr.AddArc(3, 4)
	gr.AddArc(2, 4)
	gr.AddArc(4, 5)
	gr.AddArc(1, 6)
	gr.AddArc(2, 6)
	gr.AddEdge(6, 4)

	return gr	
}

func CheckDirectedPathSpec(c gospec.Context, checkPathFunction CheckDirectedPath) {
	gr := NewDirectedMap()
	gr.AddArc(1, 2)
	gr.AddArc(2, 3)
	gr.AddArc(3, 4)
	gr.AddArc(2, 4)
	gr.AddArc(4, 5)
	gr.AddArc(1, 6)
	gr.AddArc(2, 6)
	
	c.Specify("Check path to self", func() {
		c.Expect(checkPathFunction(gr, 1, 1, nil, SimpleWeightFunc), IsTrue) 
		c.Expect(checkPathFunction(gr, 6, 6, nil, SimpleWeightFunc), IsTrue)
	})
	
	c.Specify("Check neighbours path", func() {
		c.Expect(checkPathFunction(gr, 1, 2, nil, SimpleWeightFunc), IsTrue)
		c.Expect(checkPathFunction(gr, 2, 4, nil, SimpleWeightFunc), IsTrue)
		c.Expect(checkPathFunction(gr, 1, 6, nil, SimpleWeightFunc), IsTrue)
	})
	
	c.Specify("Check reversed neighbours", func() {
		c.Expect(checkPathFunction(gr, 6, 1, nil, SimpleWeightFunc), IsFalse)
		c.Expect(checkPathFunction(gr, 4, 3, nil, SimpleWeightFunc), IsFalse)
		c.Expect(checkPathFunction(gr, 5, 4, nil, SimpleWeightFunc), IsFalse)
	})
	
	c.Specify("Check long path", func() {
		c.Expect(checkPathFunction(gr, 1, 6, nil, SimpleWeightFunc), IsTrue)
		c.Expect(checkPathFunction(gr, 1, 5, nil, SimpleWeightFunc), IsTrue)
	})
	
	c.Specify("Check weight limit", func() {
		c.Expect(checkPathFunction(gr, 1, 5, func(node NodeId, weight float) bool {
			return weight < 2.0
		}, SimpleWeightFunc), IsFalse)
	})
}

func CheckMixedPathSpec(c gospec.Context, checkPathFunction CheckMixedPath) {
	gr := generateMixedGraph1()
	
	c.Specify("Check path to self", func() {
		c.Expect(checkPathFunction(gr, 1, 1, nil, SimpleWeightFunc), IsTrue) 
		c.Expect(checkPathFunction(gr, 6, 6, nil, SimpleWeightFunc), IsTrue)
		c.Expect(checkPathFunction(gr, 4, 4, nil, SimpleWeightFunc), IsTrue)
	})
	
	c.Specify("Check directed neighbours path", func() {
		c.Expect(checkPathFunction(gr, 1, 2, nil, SimpleWeightFunc), IsTrue)
		c.Expect(checkPathFunction(gr, 2, 4, nil, SimpleWeightFunc), IsTrue)
		c.Expect(checkPathFunction(gr, 1, 6, nil, SimpleWeightFunc), IsTrue)
	})

	c.Specify("Check undirected neighbours path", func() {
		c.Expect(checkPathFunction(gr, 4, 6, nil, SimpleWeightFunc), IsTrue)
		c.Expect(checkPathFunction(gr, 6, 4, nil, SimpleWeightFunc), IsTrue)
	})
		
	c.Specify("Check reversed directed neighbours", func() {
		c.Expect(checkPathFunction(gr, 6, 1, nil, SimpleWeightFunc), IsFalse)
		c.Expect(checkPathFunction(gr, 4, 3, nil, SimpleWeightFunc), IsFalse)
		c.Expect(checkPathFunction(gr, 5, 4, nil, SimpleWeightFunc), IsFalse)
	})
	
	c.Specify("Check long path", func() {
		c.Expect(checkPathFunction(gr, 1, 6, nil, SimpleWeightFunc), IsTrue)
		c.Expect(checkPathFunction(gr, 1, 5, nil, SimpleWeightFunc), IsTrue)
		c.Expect(checkPathFunction(gr, 6, 5, nil, SimpleWeightFunc), IsTrue)
		c.Expect(checkPathFunction(gr, 3, 6, nil, SimpleWeightFunc), IsTrue)
		
		c.Expect(checkPathFunction(gr, 6, 3, nil, SimpleWeightFunc), IsFalse)
	})
	
	c.Specify("Check weight limit", func() {
		c.Expect(checkPathFunction(gr, 1, 5, func(node NodeId, weight float) bool {
			return weight < 2.0
		}, SimpleWeightFunc), IsFalse)
	})
}

func GetAllMixedPathsSpec(c gospec.Context) {
	gr := generateMixedGraph1()
	/*
	[1 2 4 6]          
	[1 2 6]            
	[1 2 3 4 6]        
	[1 6]
	*/

	pathsCnt := 0
	
	for path := range GetAllMixedPaths(gr, 1, 6) {
		pathsCnt++
		c.Expect(ContainMixedPath(gr, path, true), IsTrue)
	}
	
	c.Expect(pathsCnt, Equals, 4)
}

func TestSearch(t *testing.T) {
	r := gospec.NewRunner()

	{
		// paramenerized test creator
		cr := func(checkPathFunction CheckDirectedPath) func (c gospec.Context) {
			return func(c gospec.Context){
				CheckDirectedPathSpec(c, checkPathFunction)
			}
		}
		r.AddNamedSpec("CheckDirectedPath(Dijkstra)", cr(CheckDirectedPathDijkstra))
	}
	{
		// paramenerized test creator
		cr := func(checkPathFunction CheckMixedPath) func (c gospec.Context) {
			return func(c gospec.Context){
				CheckMixedPathSpec(c, checkPathFunction)
			}
		}
		r.AddNamedSpec("CheckMixedPath(Dijkstra)", cr(CheckMixedPathDijkstra))
	}
	
	r.AddSpec(GetAllMixedPathsSpec)


	gospec.MainGoTest(r, t)
}
