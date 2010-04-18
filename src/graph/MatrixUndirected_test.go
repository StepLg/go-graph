package graph

import (
	"erx"
	"testing"
)

func CheckErx(t *testing.T, err erx.Error) {
	if err!=nil {
		formatter := erx.NewStringFormatter("  ")
		t.Error("Erx error\n" + formatter.Format(err))
	}
}

func CheckEqBool(t *testing.T, l, r bool, msg string) {
	if (l!=r) {
		t.Error(msg)
	}
}

func Test_Edges(t *testing.T) {
	gr := NewUndirectedGraphMatrix(5)
	var err erx.Error
	_, err = gr.AddEdge(1, 2)
	CheckErx(t, err)
	_, err = gr.AddEdge(2, 3)
	CheckErx(t, err)
	_, err = gr.AddEdge(2, 4)
	CheckErx(t, err)
	_, err = gr.AddEdge(1, 4)
	CheckErx(t, err)
	
	var conn bool
	conn, err = gr.CheckEdgeBetween(1, 2)
	CheckErx(t, err)
	CheckEqBool(t, conn, true, "Missing edge 1--2")
	conn, err = gr.CheckEdgeBetween(3, 2)
	CheckErx(t, err)
	CheckEqBool(t, conn, true, "Missing edge 3--2")
	conn, err = gr.CheckEdgeBetween(4, 2)
	CheckErx(t, err)
	CheckEqBool(t, conn, true, "Missing edge 4--2")
	conn, err = gr.CheckEdgeBetween(1, 4)
	CheckErx(t, err)
	CheckEqBool(t, conn, true, "Missing edge 1--4")

	conn, err = gr.CheckEdgeBetween(1, 3)
	CheckErx(t, err)
	CheckEqBool(t, conn, false, "Phantom edge 1--3")
	conn, err = gr.CheckEdgeBetween(4, 3)
	CheckErx(t, err)
	CheckEqBool(t, conn, false, "Phantom edge 4--3")
}
