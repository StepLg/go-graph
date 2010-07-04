package graph

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

func generateDirectedGraph1() DirectedGraph {
	gr := NewDirectedMap()
	gr.AddArc(1, 2)
	gr.AddArc(2, 3)
	gr.AddArc(3, 4)
	gr.AddArc(2, 4)
	gr.AddArc(4, 5)
	gr.AddArc(1, 6)
	gr.AddArc(2, 6)
	
	return gr
}

func genDgr2IndependentSubGr() (DirectedGraphReader, DirectedGraphReader, DirectedGraphReader) {
	gr1 := NewDirectedMap()
	ReadArcsLine(gr1, "1-2-3-4")
	ReadArcsLine(gr1, "2-6")
	ReadArcsLine(gr1, "5-4-2")
	
	gr2 := NewDirectedMap()
	ReadArcsLine(gr2, "10-11-12")
	ReadArcsLine(gr2, "14-13-11")
	ReadArcsLine(gr2, "15-12-16-17")
	
	gr_merged := NewDirectedMap()
	CopyDirectedGraph(gr1, gr_merged)
	CopyDirectedGraph(gr2, gr_merged)
	
	return gr1, gr2, gr_merged
}

func genUgr2IndependentSubGr() (UndirectedGraphReader, UndirectedGraphReader, UndirectedGraphReader) {
	gr1 := NewUndirectedMap()
	ReadEdgesLine(gr1, "1-2-3-4")
	ReadEdgesLine(gr1, "2-6")
	ReadEdgesLine(gr1, "5-4-2")
	
	gr2 := NewUndirectedMap()
	ReadEdgesLine(gr2, "10-11-12")
	ReadEdgesLine(gr2, "14-13-11")
	ReadEdgesLine(gr2, "15-12-16-17")
	
	gr_merged := NewUndirectedMap()
	CopyUndirectedGraph(gr1, gr_merged)
	CopyUndirectedGraph(gr2, gr_merged)
	
	return gr1, gr2, gr_merged
}
