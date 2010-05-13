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
				pickNode, pickPriority := q.Pick()
				c.Expect(pickNode, Equals, node)
				c.Expect(pickPriority, Equals, priority)
				c.Specify("not empty", func() {
					c.Expect(q.Empty(), IsFalse)
					c.Expect(q.Size(), Equals, 1)				
				})
			})
			
			c.Specify("empty after next", func() {
				pickNode, pickPriority := q.Next()
				c.Expect(pickNode, Equals, node)
				c.Expect(pickPriority, Equals, priority)
				c.Expect(q.Empty(), IsTrue)
				c.Expect(q.Size(), Equals, 0)
			})
		})
	})
	
	c.Specify("Several items with priorities", func() {
		n1 := NodeId(1)
		p1 := 1.0
		n2 := NodeId(2)
		p2 := 2.0
		n3 := NodeId(3)
		p3 := 0.5
		n4 := NodeId(4)
		p4 := 1.5
		
		q.Add(n1, p1)
		q.Add(n2, p2)
		q.Add(n3, p3)
		q.Add(n4, p4)
		
		c.Expect(q.Size(), Equals, 4)
		node, prior := q.Next(); 
		c.Expect(node, Equals, n2)
		c.Expect(prior, Equals, p2)
		node, prior = q.Next(); 
		c.Expect(node, Equals, n4)
		c.Expect(prior, Equals, p4)
		node, prior = q.Next(); 
		c.Expect(node, Equals, n1)
		c.Expect(prior, Equals, p1)
		node, prior = q.Next(); 
		c.Expect(node, Equals, n3)
		c.Expect(prior, Equals, p3)
	})
	
	c.Specify("Manipulating items priority", func() {
		n1 := NodeId(1)
		p1 := 1.0
		n2 := NodeId(2)
		p2 := 2.0
		n3 := NodeId(3)
		p3 := 0.5
		n4 := NodeId(4)
		p4 := 1.5
		
		q.Add(n1, p1)
		q.Add(n2, p2)
		q.Add(n3, p3)
		q.Add(n4, p4)

		c.Specify("Do not decrease priority", func() {
			q.Add(n4, p4 - 1.0)
			
			c.Expect(q.Size(), Equals, 4)
			node, prior := q.Next(); 
			c.Expect(node, Equals, n2)
			c.Expect(prior, Equals, p2)
			node, prior = q.Next(); 
			c.Expect(node, Equals, n4)
			c.Expect(prior, Equals, p4)
			node, prior = q.Next(); 
			c.Expect(node, Equals, n1)
			c.Expect(prior, Equals, p1)
			node, prior = q.Next(); 
			c.Expect(node, Equals, n3)
			c.Expect(prior, Equals, p3)
		})

		c.Specify("Change middle to top", func() {
			p4 = 3.0
			q.Add(n4, p4)
			
			c.Expect(q.Size(), Equals, 4)
			node, prior := q.Next(); 
			c.Expect(node, Equals, n4)
			c.Expect(prior, Equals, p4)
			node, prior = q.Next(); 
			c.Expect(node, Equals, n2)
			c.Expect(prior, Equals, p2)
			node, prior = q.Next(); 
			c.Expect(node, Equals, n1)
			c.Expect(prior, Equals, p1)
			node, prior = q.Next(); 
			c.Expect(node, Equals, n3)
			c.Expect(prior, Equals, p3)
		})
	})
	
	c.Specify("Push more items than initial size", func() {
		n1 := NodeId(1)
		p1 := 1.0
		n2 := NodeId(2)
		p2 := 2.0
		n3 := NodeId(3)
		p3 := 0.5
		n4 := NodeId(4)
		p4 := 1.5
		n5 := NodeId(6)
		p5 := 1.6
		n6 := NodeId(7)
		p6 := 1.7
		
		
		q.Add(n1, p1)
		q.Add(n2, p2)
		q.Add(n3, p3)
		q.Add(n4, p4)
		q.Add(n5, p5)
		q.Add(n6, p6)
		
		c.Expect(q.Size(), Equals, 6)
		node, prior := q.Next(); 
		c.Expect(node, Equals, n2)
		c.Expect(prior, Equals, p2)
		node, prior = q.Next(); 
		c.Expect(node, Equals, n6)
		c.Expect(prior, Equals, p6)
		node, prior = q.Next(); 
		c.Expect(node, Equals, n5)
		c.Expect(prior, Equals, p5)
		node, prior = q.Next(); 
		c.Expect(node, Equals, n4)
		c.Expect(prior, Equals, p4)
		node, prior = q.Next(); 
		c.Expect(node, Equals, n1)
		c.Expect(prior, Equals, p1)
		node, prior = q.Next(); 
		c.Expect(node, Equals, n3)
		c.Expect(prior, Equals, p3)
	})
}

func TestStuff(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(NodesPriorityQueueSpec)
	gospec.MainGoTest(r, t)
}
