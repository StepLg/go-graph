package graph

import (
	"testing"
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
)

func VertexesPriorityQueueSpec(c gospec.Context) {
	q := newPriorityQueueSimple(5)
	
	c.Specify("Empty queue", func() {
		c.Specify("is empty", func() {
			c.Expect(q.Empty(), IsTrue)
			c.Expect(q.Size(), Equals, 0)
		})
		
		c.Specify("after add", func() {
			node := VertexId(1)
			priority := float64(0.5)
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
		n1 := VertexId(1)
		p1 := float64(1.0)
		n2 := VertexId(2)
		p2 := float64(2.0)
		n3 := VertexId(3)
		p3 := float64(0.5)
		n4 := VertexId(4)
		p4 := float64(1.5)
		
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
		n1 := VertexId(1)
		p1 := float64(1.0)
		n2 := VertexId(2)
		p2 := float64(2.0)
		n3 := VertexId(3)
		p3 := float64(0.5)
		n4 := VertexId(4)
		p4 := float64(1.5)
		
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
			p4 = float64(3.0)
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
		n1 := VertexId(1)
		p1 := float64(1.0)
		n2 := VertexId(2)
		p2 := float64(2.0)
		n3 := VertexId(3)
		p3 := float64(0.5)
		n4 := VertexId(4)
		p4 := float64(1.5)
		n5 := VertexId(6)
		p5 := float64(1.6)
		n6 := VertexId(7)
		p6 := float64(1.7)
		
		
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

func MatrixIndexerSpec(c gospec.Context) {
	size := 100
	usedIds := make(map[int]bool)
	nodesIds := make(map[VertexId]int)
	for i:=0; i<size; i++ {
		for j:=0; j<i; j++ {
			connId := matrixConnectionsIndexer(VertexId(i), VertexId(j), nodesIds, size, true)
			_, ok := usedIds[connId]
			c.Expect(ok, IsFalse)
			usedIds[connId] = true
		}
	}
}

func TestStuff(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(VertexesPriorityQueueSpec)
	r.AddSpec(MatrixIndexerSpec)
	gospec.MainGoTest(r, t)
}
