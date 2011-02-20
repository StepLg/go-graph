package graph

import (
	"math"

	"github.com/StepLg/go-erx/src/erx"
)

// Path mark, set by some of search algorithms.
type VertexPathMark struct {
	Weight float64 // Weight from one of source nodes to current node.
	PrevVertex VertexId // Previous node in path.
}

// Path weight from one of sources to node in map.
//
// Used as a result by some of minimal path search algorithms.
// To get real path from this marks map use PathFromMarks function. 
type PathMarks map[VertexId]*VertexPathMark

type ConnectionWeightFunc func(head, tail VertexId) float64

type StopFunc func(node VertexId, sumWeight float64) bool

func SimpleWeightFunc(head, tail VertexId) float64 {
	return float64(1.0)
}

// Retrieving path from path marks.
func PathFromMarks(marks PathMarks, destination VertexId) Vertexes {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Retrieving path from path marks.", e)
			err.AddV("marks", marks)
			err.AddV("destination", destination)
			panic(err)
		}
	}()
	destInfo, ok := marks[destination]
	if !ok || destInfo.Weight==math.MaxFloat64 {
		// no path from any source to destination
		return nil
	}
	
	curVertexInfo := destInfo
	path := make(Vertexes, 10)
	curPathPos := 0
	path[curPathPos] = destination
	curPathPos++
	for curVertexInfo.Weight > 0.0 {
		if len(path)==curPathPos {
			// reallocate memory for path
			tmp := make(Vertexes, 2*curPathPos)
			copy(tmp, path)
			path = tmp
		}
		path[curPathPos] = curVertexInfo.PrevVertex
		curPathPos++
		var ok bool
		curVertexInfo, ok = marks[curVertexInfo.PrevVertex]
		if !ok {
			err := erx.NewError("Can't find path mark info for vertex in path.")
			err.AddV("vertex", curVertexInfo.PrevVertex)
			err.AddV("cur path", path)
			panic(err)
		}
	}
	
	path = path[0:curPathPos]
	
	// reversing path
	pathLen := len(path)
	for i:=0; i<pathLen/2; i++ {
		path[i], path[pathLen-i-1] = path[pathLen-i-1], path[i]
	}
	
	return path
}

