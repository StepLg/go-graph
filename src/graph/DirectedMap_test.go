package graph

import (
	"testing"
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
	"runtime"
	"strings"
	"github.com/StepLg/go-erx/src/erx"
)

func init() {
	// adding to erx directory prefix to cut from file names
	_, file, _, _ := runtime.Caller(0)
	dirName := file[0:strings.LastIndex(file, "/")]
	prevDirName := dirName[0:strings.LastIndex(dirName, "/")+1]
	erx.AddPathCut(prevDirName)
}

func DirectedMapSpec(c gospec.Context) {
	gr := NewDirectedMap()
	
	c.Specify("Empty directed graph", func() {
		c.Specify("contain no nodes", func() {
			c.Expect(gr.NodesCnt(), Equals, 0)
		})
		c.Specify("contain no edges", func() {
			c.Expect(gr.ArcsCnt(), Equals, 0)
		})
	})
	
	c.Specify("New node in empty graph", func() {
		nodeId := NodeId(1)
		gr.AddNode(nodeId)
				
		c.Specify("changing nodes count", func() {
			c.Expect(gr.NodesCnt(), Equals, 1)
		})
		
		c.Specify("doesn't change arrows count", func() {
			c.Expect(gr.ArcsCnt(), Equals, 0)
		})
		
		c.Specify("no accessors", func() {
			accessors := gr.GetAccessors(nodeId)
			c.Expect(len(accessors), Equals, 0)
		})
		
		c.Specify("no predecessors", func() {
			c.Expect(len(gr.GetPredecessors(nodeId)), Equals, 0)
		})
		
		c.Specify("node becomes a source", func() {
			sources := gr.GetSources()
			c.Expect(sources, ContainsExactly, Values(nodeId))
		})

		c.Specify("node becomes a sink", func() {
			sinks := gr.GetSinks()
			c.Expect(sinks, ContainsExactly, Values(nodeId))
		})

	})
	
	c.Specify("New arrow in empty graph", func() {
		nodeId := NodeId(1)
		anotherNodeId := NodeId(2)
		gr.AddArc(nodeId, anotherNodeId)
		c.Expect(gr.CheckArc(nodeId, anotherNodeId), Equals, true)

		c.Specify("changing nodes count", func() {
			c.Expect(gr.NodesCnt(), Equals, 2)
		})
		
		c.Specify("changing arrows count", func() {
			c.Expect(gr.ArcsCnt(), Equals, 1)
		})
		
		c.Specify("correct accessors in arrow start", func() {
			c.Expect(gr.GetAccessors(nodeId), ContainsExactly, Values(anotherNodeId))
		})

		c.Specify("correct predecessors in arrow start", func() {
			c.Expect(len(gr.GetPredecessors(nodeId)), Equals, 0)
		})

		c.Specify("correct accessors in arrow end", func() {
			c.Expect(len(gr.GetAccessors(anotherNodeId)), Equals, 0)
		})

		c.Specify("correct predecessors in arrow end", func() {
			c.Expect(gr.GetPredecessors(anotherNodeId), ContainsExactly, Values(nodeId))
		})
		
		c.Specify("arrow start becomes a source", func() {
			c.Expect(gr.GetSources(), ContainsExactly, Values(nodeId))
		})

		c.Specify("arrow end becomes a sink", func() {
			c.Expect(gr.GetSinks(), ContainsExactly, Values(anotherNodeId))
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
			c.Expect(gr.NodesCnt(), Equals, 7)
		})
		
		c.Specify("checking arrows count", func() {
			c.Expect(gr.ArcsCnt(), Equals, 7)
		})
		
		c.Specify("checking sources", func() {
			sources := gr.GetSources()
			c.Expect(sources, ContainsExactly, Values(NodeId(4), NodeId(6)))
			
			c.Specify("every source hasn't any predecessors", func() {
				for _, nodeId := range sources {
					c.Expect(len(gr.GetPredecessors(nodeId)), Equals, 0)
				}
			})
		})
		
		c.Specify("checking sinks", func() {
			sinks := gr.GetSinks()
			c.Expect(sinks, ContainsExactly, Values(NodeId(5), NodeId(7)))

			c.Specify("every sink hasn't any accessors", func() {
				for _, nodeId := range sinks {
					c.Expect(len(gr.GetAccessors(nodeId)), Equals, 0)
				}
			})
		})
		
		c.Specify("checking accessors in intermediate node", func() {
			accessors := gr.GetAccessors(1)
			c.Expect(accessors, ContainsExactly, Values(NodeId(2), NodeId(5), NodeId(7)))
			
			c.Specify("every accessor has this node in predecessors", func() {
				for _, nodeId := range accessors {
					c.Expect(gr.GetPredecessors(nodeId), Contains, NodeId(1))
				}
			})
		})

		c.Specify("checking predecessors in intermediate node", func() {
			predecessors := gr.GetPredecessors(5)
			c.Expect(predecessors, ContainsExactly, Values(NodeId(1), NodeId(4)))
			
			c.Specify("every predecessor has this node in accessors", func() {
				for _, nodeId := range predecessors {
					c.Expect(gr.GetAccessors(nodeId), Contains, NodeId(5))
				}
			})
		})

	})
}

func TestDirectedMapSpec(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(DirectedMapSpec)
	gospec.MainGoTest(r, t)
}
