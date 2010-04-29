package graph

import (
	"erx"
	"testing"
	"github.com/orfjackal/gospec/src/gospec"
)

func CheckErx(t *testing.T, err erx.Error) {
	if err!=nil {
		formatter := erx.NewStringFormatter("  ")
		t.Error("Erx error\n" + formatter.Format(err))
	}
}

func CheckEq(t *testing.T, l, r interface{}, msg string) {
	if (l!=r) {
		t.Error(msg)
	}
}

func NodesBinarySearch(nodes Nodes, n NodeId) bool {
	if len(nodes)==0 {
		return false
	}
	
	switch len(nodes) {
		case 0 : return false
		case 1 : return nodes[0]==n
		default:
			m := len(nodes)>>1
			if nodes[m]==n {
				return true
			} else if nodes[m]<n {
				return NodesBinarySearch(nodes[0:m], n)
			} else {
				return NodesBinarySearch(nodes[m+1:], n)
			}
	}
	return false
}

func Test_MatrixUndirected_Edges(t *testing.T) {
	gr := NewUndirectedGraphMatrix(5)
	var err erx.Error
	CheckErx(t, gr.AddEdge(1, 2))
	CheckErx(t, gr.AddEdge(2, 3))
	CheckErx(t, gr.AddEdge(2, 4))
	CheckErx(t, gr.AddEdge(1, 4))
	
	var conn bool
	conn, err = gr.CheckEdge(1, 2)
	CheckErx(t, err)
	CheckEq(t, conn, true, "Missing edge 1--2")
	conn, err = gr.CheckEdge(3, 2)
	CheckErx(t, err)
	CheckEq(t, conn, true, "Missing edge 3--2")
	conn, err = gr.CheckEdge(4, 2)
	CheckErx(t, err)
	CheckEq(t, conn, true, "Missing edge 4--2")
	conn, err = gr.CheckEdge(1, 4)
	CheckErx(t, err)
	CheckEq(t, conn, true, "Missing edge 1--4")

	conn, err = gr.CheckEdge(1, 3)
	CheckErx(t, err)
	CheckEq(t, conn, false, "Phantom edge 1--3")
	conn, err = gr.CheckEdge(4, 3)
	CheckErx(t, err)
	CheckEq(t, conn, false, "Phantom edge 4--3")
}

func MatrixUndirectedInternalSpec(c gospec.Context) {
	gr := NewUndirectedGraphMatrix(5)
	
	c.Specify("Correct edge ids generation to first node", func() {
		for i := 2; i < gr.GetCapacity(); i++ {
			connId, err := gr.getConnectionId(1, NodeId(i), true)
			c.Expect(err, gospec.IsNil)
			c.Expect(int(connId), gospec.Equals, i-2)
			c.Expect(gr.GetSize(), gospec.Equals, i)
		}
	})
}

func TestMatrixUndirectedSpec(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(MatrixUndirectedInternalSpec)
	gospec.MainGoTest(r, t)
}
