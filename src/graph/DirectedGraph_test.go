package graph

import (
	"testing"
	"runtime"
	"strings"
	"github.com/StepLg/go-erx/src/erx"
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
)

func init() {
	// adding to erx directory prefix to cut from file names
	_, file, _, _ := runtime.Caller(0)
	dirName := file[0:strings.LastIndex(file, "/")]
	prevDirName := dirName[0:strings.LastIndex(dirName, "/")+1]
	erx.AddPathCut(prevDirName)
}

func DirectedGraphSpec(c gospec.Context, graphCreator func() DirectedGraph) {
	gr := graphCreator()
	
	c.Specify("Empty directed graph", func() {
		c.Specify("contain no nodes", func() {
			c.Expect(gr.Order(), Equals, 0)
		})
		c.Specify("contain no edges", func() {
			c.Expect(gr.ArcsCnt(), Equals, 0)
		})
	})
	
	c.Specify("New node in empty graph", func() {
		VertexId := VertexId(1)
		gr.AddNode(VertexId)
				
		c.Specify("changing nodes count", func() {
			c.Expect(gr.Order(), Equals, 1)
		})
		
		c.Specify("doesn't change arrows count", func() {
			c.Expect(gr.ArcsCnt(), Equals, 0)
		})
		
		c.Specify("no accessors", func() {
			accessors := CollectVertexes(gr.GetAccessors(VertexId))
			c.Expect(len(accessors), Equals, 0)
		})
		
		c.Specify("no predecessors", func() {
			c.Expect(len(CollectVertexes(gr.GetPredecessors(VertexId))), Equals, 0)
		})
		
		c.Specify("node becomes a source", func() {
			sources := CollectVertexes(gr.GetSources())
			c.Expect(sources, ContainsExactly, Values(VertexId))
		})

		c.Specify("node becomes a sink", func() {
			sinks := CollectVertexes(gr.GetSinks())
			c.Expect(sinks, ContainsExactly, Values(VertexId))
		})

	})
	
	c.Specify("New arrow in empty graph", func() {
		vertexId := VertexId(1)
		anotherVertexId := VertexId(2)
		gr.AddArc(vertexId, anotherVertexId)
		c.Expect(gr.CheckArc(vertexId, anotherVertexId), Equals, true)

		c.Specify("changing nodes count", func() {
			c.Expect(gr.Order(), Equals, 2)
		})
		
		c.Specify("changing arrows count", func() {
			c.Expect(gr.ArcsCnt(), Equals, 1)
		})
		
		c.Specify("correct accessors in arrow start", func() {
			c.Expect(CollectVertexes(gr.GetAccessors(vertexId)), ContainsExactly, Values(anotherVertexId))
		})

		c.Specify("correct predecessors in arrow start", func() {
			c.Expect(len(CollectVertexes(gr.GetPredecessors(vertexId))), Equals, 0)
		})

		c.Specify("correct accessors in arrow end", func() {
			c.Expect(len(CollectVertexes(gr.GetAccessors(anotherVertexId))), Equals, 0)
		})

		c.Specify("correct predecessors in arrow end", func() {
			c.Expect(CollectVertexes(gr.GetPredecessors(anotherVertexId)), ContainsExactly, Values(vertexId))
		})
		
		c.Specify("arrow start becomes a source", func() {
			c.Expect(CollectVertexes(gr.GetSources()), ContainsExactly, Values(vertexId))
		})

		c.Specify("arrow end becomes a sink", func() {
			c.Expect(CollectVertexes(gr.GetSinks()), ContainsExactly, Values(anotherVertexId))
		})
	})
	
	c.Specify("A bit more complex example", func() {
		gr.AddArc(1, 2)
		gr.AddArc(2, 3)
		gr.AddArc(3, 1)
		gr.AddArc(1, 5)
		gr.AddArc(4, 5)
		gr.AddArc(6, 2)
		gr.AddArc(1, 7)
		
		c.Specify("checking nodes count", func() {
			c.Expect(gr.Order(), Equals, 7)
		})
		
		c.Specify("checking arrows count", func() {
			c.Expect(gr.ArcsCnt(), Equals, 7)
		})
		
		c.Specify("checking sources", func() {
			sources := CollectVertexes(gr.GetSources())
			c.Expect(sources, ContainsExactly, Values(VertexId(4), VertexId(6)))
			
			c.Specify("every source hasn't any predecessors", func() {
				for _, vertexId := range sources {
					c.Expect(len(CollectVertexes(gr.GetPredecessors(vertexId))), Equals, 0)
				}
			})
		})
		
		c.Specify("checking sinks", func() {
			sinks := CollectVertexes(gr.GetSinks())
			c.Expect(sinks, ContainsExactly, Values(VertexId(5), VertexId(7)))

			c.Specify("every sink hasn't any accessors", func() {
				for _, vertexId := range sinks {
					c.Expect(len(CollectVertexes(gr.GetAccessors(vertexId))), Equals, 0)
				}
			})
		})
		
		c.Specify("checking accessors in intermediate node", func() {
			accessors := CollectVertexes(gr.GetAccessors(1))
			c.Expect(accessors, ContainsExactly, Values(VertexId(2), VertexId(5), VertexId(7)))
			
			c.Specify("every accessor has this node in predecessors", func() {
				for _, vertexId := range accessors {
					c.Expect(CollectVertexes(gr.GetPredecessors(vertexId)), Contains, VertexId(1))
				}
			})
		})

		c.Specify("checking predecessors in intermediate node", func() {
			predecessors := CollectVertexes(gr.GetPredecessors(5))
			c.Expect(predecessors, ContainsExactly, Values(VertexId(1), VertexId(4)))
			
			c.Specify("every predecessor has this node in accessors", func() {
				for _, vertexId := range predecessors {
					c.Expect(CollectVertexes(gr.GetAccessors(vertexId)), Contains, VertexId(5))
				}
			})
		})

	})
}

func TestDirectedGraphSpec(t *testing.T) {
	r := gospec.NewRunner()

	// paramenerized test creator
	cr := func(graphCreator func() DirectedGraph) func (c gospec.Context) {
		return func(c gospec.Context){
			DirectedGraphSpec(c, graphCreator)
		}
	}
	
	r.AddNamedSpec("DirectedGraph(DirectedMap)", cr(func() DirectedGraph {
		return DirectedGraph(NewDirectedMap())
	}))
	r.AddNamedSpec("DirectedGraph(MixedMatrix)", cr(func() DirectedGraph {
		return DirectedGraph(NewMixedMatrix(10))
	}))
	r.AddNamedSpec("DirectedGraph(MixedMap)", cr(func() DirectedGraph {
		return DirectedGraph(NewMixedMap())
	}))
	gospec.MainGoTest(r, t)
}
