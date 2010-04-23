package graph

import (
	"testing"
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
)


func DirectedMapSpec(c gospec.Context) {
	gr := NewDirectedMap()
	
	c.Specify("Empty directed graph", func() {
		c.Specify("contain no nodes", func() {
			c.Expect(gr.NodesCnt(), Equals, 0)
		})
		c.Specify("contain no edges", func() {
			c.Expect(gr.ArrowsCnt(), Equals, 0)
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
			c.Expect(gr.ArrowsCnt(), Equals, 0)
		})
		
		c.Specify("no accessors", func() {
			accessors, err := gr.GetAccessors(nodeId)
			c.Expect(err, IsNil)
			c.Expect(len(accessors), Equals, 0)
		})
		
		c.Specify("no predecessors", func() {
			predecessors, err := gr.GetPredecessors(nodeId)
			c.Expect(err, IsNil)
			c.Expect(len(predecessors), Equals, 0)
		})
		
		c.Specify("node becomes a source", func() {
			sources, err := gr.GetSources()
			c.Expect(err, IsNil)
			c.Expect(sources, ContainsExactly, Values(nodeId))
		})

		c.Specify("node becomes a sink", func() {
			sinks, err := gr.GetSinks()
			c.Expect(err, IsNil)
			c.Expect(sinks, ContainsExactly, Values(nodeId))
		})

	})
	
	c.Specify("New arrow in empty graph", func() {
		nodeId := NodeId(1)
		anotherNodeId := NodeId(2)
		err := gr.AddArrow(nodeId, anotherNodeId)
		c.Expect(err, IsNil)
		
		c.Specify("changing nodes count", func() {
			c.Expect(gr.NodesCnt(), Equals, 2)
		})
		
		c.Specify("changing arrows count", func() {
			c.Expect(gr.ArrowsCnt(), Equals, 1)
		})
		
		c.Specify("correct accessors in arrow start", func() {
			accessors, err := gr.GetAccessors(nodeId)
			c.Expect(err, IsNil)
			c.Expect(accessors, ContainsExactly, Values(anotherNodeId))
		})

		c.Specify("correct predecessors in arrow start", func() {
			predecessors, err := gr.GetPredecessors(nodeId)
			c.Expect(err, IsNil)
			c.Expect(len(predecessors), Equals, 0)
		})

		c.Specify("correct accessors in arrow end", func() {
			accessors, err := gr.GetAccessors(anotherNodeId)
			c.Expect(err, IsNil)
			c.Expect(len(accessors), Equals, 0)
		})

		c.Specify("correct predecessors in arrow end", func() {
			predecessors, err := gr.GetPredecessors(anotherNodeId)
			c.Expect(err, IsNil)
			c.Expect(predecessors, ContainsExactly, Values(nodeId))
		})
		
		c.Specify("arrow start becomes a source", func() {
			sources, err := gr.GetSources()
			c.Expect(err, IsNil)
			c.Expect(sources, ContainsExactly, Values(nodeId))
		})

		c.Specify("arrow end becomes a sink", func() {
			sinks, err := gr.GetSinks()
			c.Expect(err, IsNil)
			c.Expect(sinks, ContainsExactly, Values(anotherNodeId))
		})
	})
	
	c.Specify("A bit more complex example", func() {
		c.Expect(gr.AddArrow(1, 2), IsNil)
		c.Expect(gr.AddArrow(2, 3), IsNil)
		c.Expect(gr.AddArrow(3, 1), IsNil)
		c.Expect(gr.AddArrow(1, 5), IsNil)
		c.Expect(gr.AddArrow(4, 5), IsNil)
		c.Expect(gr.AddArrow(6, 2), IsNil)
		c.Expect(gr.AddArrow(1, 7), IsNil)
		
		c.Specify("checking nodes count", func() {
			c.Expect(gr.NodesCnt(), Equals, 7)
		})
		
		c.Specify("checking arrows count", func() {
			c.Expect(gr.ArrowsCnt(), Equals, 7)
		})
		
		c.Specify("checking sources", func() {
			sources, err := gr.GetSources()
			c.Expect(err, IsNil)
			c.Expect(sources, ContainsExactly, Values(NodeId(4), NodeId(6)))
			
			c.Specify("every source hasn't any predecessors", func() {
				for _, nodeId := range sources {
					predecessors, err := gr.GetPredecessors(nodeId)
					c.Expect(err, IsNil)
					c.Expect(len(predecessors), Equals, 0)
				}
			})
		})
		
		c.Specify("checking sinks", func() {
			sinks, err := gr.GetSinks()
			c.Expect(err, IsNil)
			c.Expect(sinks, ContainsExactly, Values(NodeId(5), NodeId(7)))

			c.Specify("every sink hasn't any accessors", func() {
				for _, nodeId := range sinks {
					accessors, err := gr.GetAccessors(nodeId)
					c.Expect(err, IsNil)
					c.Expect(len(accessors), Equals, 0)
				}
			})
		})
		
		c.Specify("checking accessors in intermediate node", func() {
			accessors, err := gr.GetAccessors(1)
			c.Expect(err, IsNil)
			c.Expect(accessors, ContainsExactly, Values(NodeId(2), NodeId(5), NodeId(7)))
			
			c.Specify("every accessor has this node in predecessors", func() {
				for _, nodeId := range accessors {
					predecessors, err := gr.GetPredecessors(nodeId)
					c.Expect(err, IsNil)
					c.Expect(predecessors, Contains, NodeId(1))
				}
			})
		})

		c.Specify("checking predecessors in intermediate node", func() {
			predecessors, err := gr.GetPredecessors(5)
			c.Expect(err, IsNil)
			c.Expect(predecessors, ContainsExactly, Values(NodeId(1), NodeId(4)))
			
			c.Specify("every predecessor has this node in accessors", func() {
				for _, nodeId := range predecessors {
					accessors, err := gr.GetAccessors(nodeId)
					c.Expect(err, IsNil)
					c.Expect(accessors, Contains, NodeId(5))
				}
			})
		})

	})
}

func TestDirectedMapSpec(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec("DirectedMapSpec", DirectedMapSpec)
	gospec.MainGoTest(r, t)
}
