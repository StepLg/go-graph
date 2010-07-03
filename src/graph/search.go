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

// Generic check path algorithm for all graph types
// 
// Checking path between from and to nodes, using getNeighbours function
// to figure out connected nodes on each step of algorithm.
// 
// stopFunc is used to cut bad paths using user-defined criteria
// 
// weightFunction calculates total path weight
// 
// As a result CheckPathDijkstra returns total weight of path, if it exists.
func CheckPathDijkstra(neighboursExtractor OutNeighboursExtractor, from, to VertexId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) (float64, bool) {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Check path graph with Dijkstra algorithm", e)
			err.AddV("from", from)
			err.AddV("to", to)
			panic(err)
		}
	}()
	
	if from==to {
		return 0.0, true
	}
	
	q := newPriorityQueueSimple(10)
	q.Add(from, 0.0)
	
	for !q.Empty() {
		curNode, curWeight := q.Next()
		curWeight = -curWeight // because we inverse weight in priority queue
	
		for nextNode := range neighboursExtractor.GetOutNeighbours(curNode).VertexesIter() {
			arcWeight := weightFunction(curNode, nextNode)
			if arcWeight < 0 {
				err := erx.NewError("Negative weight detected")
				err.AddV("head", curNode)
				err.AddV("tail", nextNode)
				err.AddV("weight", arcWeight)
				panic(err)
			}
			nextWeight := curWeight + arcWeight
			if nextNode==to {
				return nextWeight, true
			}
			if stopFunc==nil || !stopFunc(nextNode, nextWeight) {
				q.Add(nextNode, -nextWeight)
			}
		}
	}
	
	return -1.0, false
}

type CheckDirectedPath func(gr DirectedGraphArcsReader, from, to VertexId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool

func CheckDirectedPathDijkstra(gr DirectedGraphArcsReader, from, to VertexId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool {
	_, pathExists := CheckPathDijkstra(NewDgraphOutNeighboursExtractor(gr), from, to, stopFunc, weightFunction)
	return pathExists
}

type CheckUndirectedPath func(gr UndirectedGraphEdgesReader, from, to VertexId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool

func CheckUndirectedPathDijkstra(gr UndirectedGraphEdgesReader, from, to VertexId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool {
	_, pathExists := CheckPathDijkstra(NewUgraphOutNeighboursExtractor(gr), from, to, stopFunc, weightFunction)
	return pathExists
}

type CheckMixedPath func(gr MixedGraphConnectionsReader, from, to VertexId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool

func CheckMixedPathDijkstra(gr MixedGraphConnectionsReader, from, to VertexId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool {
	_, pathExists := CheckPathDijkstra(NewMgraphOutNeighboursExtractor(gr), from, to, stopFunc, weightFunction)
	return pathExists
}

// Get all paths from one node to another
//
// This algorithms doesn't take any loops into paths.
func GetAllPaths(neighboursExtractor OutNeighboursExtractor, from, to VertexId) <-chan []VertexId {
	curPath := make([]VertexId, 10)
	nodesStatus := make(map[VertexId]bool)
	ch := make(chan []VertexId)
	go getAllPaths_helper(neighboursExtractor, from, to, curPath, 0, nodesStatus, ch, true)
	return ch
}

func getAllPaths_helper(neighboursExtractor OutNeighboursExtractor, from, to VertexId, curPath []VertexId, pathPos int, nodesStatus map[VertexId]bool, ch chan []VertexId, closeChannel bool) {
	if _, ok := nodesStatus[from]; ok {
		return
	}
	if pathPos==len(curPath) {
		// reallocate curPath slice to add new elements
		tmp := make([]VertexId, 2*pathPos)
		copy(tmp, curPath)
		curPath = tmp
	}
	
	curPath[pathPos] = from

	if from==to { 
		if pathPos>0 {
			pathCopy := make([]VertexId, pathPos+1)
			copy(pathCopy, curPath[0:pathPos+1])
			ch <- pathCopy
		}
		return
	}
	nodesStatus[from] = true
	
	for nextNode := range neighboursExtractor.GetOutNeighbours(from).VertexesIter() {
		getAllPaths_helper(neighboursExtractor, nextNode, to, curPath, pathPos+1, nodesStatus, ch, false)
	}
	
	nodesStatus[from] = false, false
	
	if closeChannel {
		close(ch)
	}
	return
}

func GetAllDirectedPaths(gr DirectedGraphArcsReader, from, to VertexId) <-chan []VertexId {
	return GetAllPaths(NewDgraphOutNeighboursExtractor(gr), from, to)
}

func GetAllUndirectedPaths(gr UndirectedGraphEdgesReader, from, to VertexId) <-chan []VertexId {
	return GetAllPaths(NewUgraphOutNeighboursExtractor(gr), from, to)
}

func GetAllMixedPaths(gr MixedGraphConnectionsReader, from, to VertexId) <-chan []VertexId {
	return GetAllPaths(NewMgraphOutNeighboursExtractor(gr), from, to)
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


// Compute multi-source shortest paths with Bellman-Ford algorithm
//
// Returs map, contains all nodes from graph. If there is no path from source to node in map
// then value for this node is math.MaxFloat64
//
// Returns nil if there are negative cycles. 
func BellmanFordMultiSource(gr DirectedGraphReader, sources Vertexes, weightFunc ConnectionWeightFunc) PathMarks {
	marks := make(PathMarks)
	for vertex := range gr.VertexesIter() {
		marks[vertex] = &VertexPathMark{Weight: math.MaxFloat64, PrevVertex: 0}
	}
	
	for _, vertex := range sources {
		marks[vertex].Weight = 0.0
	}
	
	nodesCnt := gr.Order()
	for i:=0; i<nodesCnt; i++ {
		for conn := range gr.ArcsIter() {
			possibleWeight := marks[conn.Tail].Weight + weightFunc(conn.Tail, conn.Head)
			if marks[conn.Head].Weight > possibleWeight {
				marks[conn.Head].PrevVertex = conn.Tail
				marks[conn.Head].Weight = possibleWeight
			}
		}
	}
	
	for conn := range gr.ArcsIter() {
		if marks[conn.Head].Weight > marks[conn.Tail].Weight + weightFunc(conn.Tail, conn.Head) {
			return nil
		}
	}
	
	return marks
}

func BellmanFordSingleSource(gr DirectedGraphReader, source VertexId, weightFunc ConnectionWeightFunc) PathMarks {
	return BellmanFordMultiSource(gr, Vertexes{source}, weightFunc)
}

// Compute multi-source shortest paths with Bellman-Ford algorithm
//
// Returs map, contains all accessiable vertexes from sources vertexes with
// minimal path weight.
//
// Returns nil if there are negative cycles. 
func BellmanFordLightMultiSource(gr OutNeighboursExtractor, sources Vertexes, weightFunc ConnectionWeightFunc) PathMarks {
	marks := make(PathMarks)
	for _, vertex := range sources {
		marks[vertex] = &VertexPathMark{Weight: 0.0, PrevVertex: 0}
	}
	
	for i:=0; i<len(marks); i++ {
		for vertex, vertexInfo := range marks {
			for nextVertex := range gr.GetOutNeighbours(vertex).VertexesIter() {
				possibleWeight := vertexInfo.Weight + weightFunc(vertex, nextVertex)
				if nextVertexInfo, ok := marks[nextVertex]; !ok || nextVertexInfo.Weight > possibleWeight {
					marks[vertex] = &VertexPathMark{Weight: possibleWeight, PrevVertex: vertex}
				}
			}
		}
	}
	
	// checking for negative cycles
	for vertex, vertexInfo := range marks {
		for nextVertex := range gr.GetOutNeighbours(vertex).VertexesIter() {
			if marks[nextVertex].Weight > vertexInfo.Weight + weightFunc(vertex, nextVertex) {
				return nil
			}
		}
	}
	
	return marks
}

func BellmanFordLightSingleSource(gr OutNeighboursExtractor, source VertexId, weightFunc ConnectionWeightFunc) PathMarks {
	return BellmanFordLightMultiSource(gr, Vertexes{source}, weightFunc)
}
