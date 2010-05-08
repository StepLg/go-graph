package graph

import (
	"testing"
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
)

func NodesPriorityQueueSpec(c gospec.Context) {
	q := newPriorityQueueSimple(5)
	
	c.Specify("Empty queue", func() {
		c.Specify("is empty", func() {
			c.Expect(q.Empty(), IsTrue)
			c.Expect(q.Size(), Equals, 0)
		})
		
		c.Specify("after add", func() {
			node := NodeId(1)
			priority := 0.5
			q.Add(node, priority)
			
			c.Specify("no longer empty", func() {
				c.Expect(q.Empty(), IsFalse)
				c.Expect(q.Size(), Equals, 1)				
			})
			
			c.Specify("can make pick", func() {
				c.Expect(q.Pick(), Equals, node)
				c.Specify("not empty", func() {
					c.Expect(q.Empty(), IsFalse)
					c.Expect(q.Size(), Equals, 1)				
				})
			})
			
			c.Specify("empty after next", func() {
				c.Expect(q.Next(), Equals, node)
				c.Expect(q.Empty(), IsTrue)
				c.Expect(q.Size(), Equals, 0)
			})
		})
	})
	
	c.Specify("Several items with priorities", func() {
		n1 := NodeId(1)
		n2 := NodeId(2)
		n3 := NodeId(3)
		n4 := NodeId(4)
		
		q.Add(n1, 1.0)
		q.Add(n2, 2.0)
		q.Add(n3, 0.5)
		q.Add(n4, 1.5)
		
		c.Expect(q.Size(), Equals, 4)
		c.Expect(q.Next(), Equals, n2)
		c.Expect(q.Next(), Equals, n4)
		c.Expect(q.Next(), Equals, n1)
		c.Expect(q.Next(), Equals, n3)
	})
}

func TestStuff(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(NodesPriorityQueueSpec)
	gospec.MainGoTest(r, t)
}
