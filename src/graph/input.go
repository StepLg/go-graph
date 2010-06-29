package graph

import (
	"strconv"
	"strings"
	
	"github.com/StepLg/go-erx/src/erx"
)

func ReadEdgesLine(gr UndirectedGraphEdgesWriter, line string) {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Parsing graph edges from line.", e)
			err.AddV("line", line)
			panic(err)
		}
	}()
	var prevVertexId VertexId
	hasPrev := false
	for _, nodeAsStr := range strings.Split(line, "-", 0) {
		nodeAsStr = strings.Trim(nodeAsStr, " \t\n")
		nodeAsInt, err := strconv.Atoi(nodeAsStr)
		if err!=nil {
			errErx := erx.NewSequent("Can't parse node id.", err)
			errErx.AddV("chunk", nodeAsStr)
			panic(errErx)
		}
		
		VertexId := VertexId(nodeAsInt)
		if hasPrev {
			gr.AddEdge(prevVertexId, VertexId)
		} else {
			hasPrev = true
		}
		prevVertexId = VertexId
	}
	return
}

func ReadArcsLine(gr DirectedGraphArcsWriter, line string) {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Parsing graph arcs from line.", e)
			err.AddV("line", line)
			panic(err)
		}
	}()
	var prevVertexId VertexId
	hasPrev := false
	for _, nodeAsStr := range strings.Split(line, "-", 0) {
		nodeAsStr = strings.Trim(nodeAsStr, " \t\n")
		nodeAsInt, err := strconv.Atoi(nodeAsStr)
		if err!=nil {
			errErx := erx.NewSequent("Can't parse node id.", err)
			errErx.AddV("chunk", nodeAsStr)
			panic(errErx)
		}
		
		VertexId := VertexId(nodeAsInt)
		if hasPrev {
			gr.AddArc(prevVertexId, VertexId)
		} else {
			hasPrev = true
		}
		prevVertexId = VertexId
	}
	return
}
