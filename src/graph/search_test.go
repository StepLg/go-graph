package graph

import (
	"testing"
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
)

func CheckDirectedPathSpec(c gospec.Context, checkPathFunction CheckDirectedPath) {
	gr := generateDirectedGraph1()
	
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
		c.Expect(checkPathFunction(gr, 1, 5, func(node VertexId, weight float64) bool {
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
		c.Expect(checkPathFunction(gr, 1, 5, func(node VertexId, weight float64) bool {
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

func BellmanFordSingleSourceSpec(c gospec.Context) {
	gr := generateDirectedGraph1()
	
	marks := BellmanFordSingleSource(gr, VertexId(2), SimpleWeightFunc)
	c.Expect(len(marks), Equals, gr.Order())

	c.Expect(PathFromMarks(marks, VertexId(6)), ContainsExactly, Values(VertexId(2), VertexId(6)))
	c.Expect(PathFromMarks(marks, VertexId(5)), ContainsExactly, Values(VertexId(2), VertexId(4), VertexId(5)))
	c.Expect(PathFromMarks(marks, VertexId(1)), ContainsExactly, Values())
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
	r.AddSpec(BellmanFordSingleSourceSpec)


	gospec.MainGoTest(r, t)
}
