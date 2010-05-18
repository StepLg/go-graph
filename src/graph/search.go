package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)


type ConnectionWeightFunc func(head, tail NodeId) float

type StopFunc func(node NodeId, sumWeight float) bool

func SimpleWeightFunc(head, tail NodeId) float {
	return 1.0
}

type CheckDirectedPath func(gr DirectedGraphArcsReader, from, to NodeId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool

func CheckDirectedPathDijkstra(gr DirectedGraphArcsReader, from, to NodeId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Check path in directed graph with Dijkstra algorithm", e)
			err.AddV("from", from)
			err.AddV("to", to)
			panic(err)
		}
	}()
	
	if from==to {
		return true
	}
	
	q := newPriorityQueueSimple(10)
	q.Add(from, 0)
	
	for !q.Empty() {
		curNode, curWeight := q.Next()
		curWeight = -curWeight // because we inverse weight in priority queue
		accessors := gr.GetAccessors(curNode)
		for _, nextNode := range accessors {
			if nextNode==to {
				return true
			}
			arcWeight := weightFunction(curNode, nextNode)
			if arcWeight < 0 {
				err := erx.NewError("Negative weight detected")
				err.AddV("head", curNode)
				err.AddV("tail", nextNode)
				err.AddV("weight", arcWeight)
				panic(err)
			}
			nextWeight := curWeight + arcWeight
			if stopFunc==nil || !stopFunc(nextNode, nextWeight) {
				q.Add(nextNode, -nextWeight)
			}
		}
	}
	
	return false
}

type CheckMixedPath func(gr MixedGraphConnectionsReader, from, to NodeId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool

func CheckMixedPathDijkstra(gr MixedGraphConnectionsReader, from, to NodeId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Check path in mixed graph with Dijkstra algorithm", e)
			err.AddV("from", from)
			err.AddV("to", to)
			panic(err)
		}
	}()
	
	if from==to {
		return true
	}
	
	q := newPriorityQueueSimple(10)
	q.Add(from, 0)
	
	for !q.Empty() {
		curNode, curWeight := q.Next()
		curWeight = -curWeight // because we inverse weight in priority queue
		
		// todo: implement GetAccessors and GetNeighbours as channels instead of slices
		accessors := gr.GetAccessors(curNode)
		neighbours := gr.GetNeighbours(curNode)
		
		if len(accessors)+len(neighbours)==0 {
			continue
		}
		
		nextNodes := make([]NodeId, len(accessors) + len(neighbours))
		if len(accessors)!=0 {
			copy(nextNodes[0:len(accessors)], accessors)
		}
		if len(neighbours)!=0 {
			copy(nextNodes[len(accessors):], neighbours)
		}
		
		for _, nextNode := range nextNodes {
			if nextNode==to {
				return true
			}
			arcWeight := weightFunction(curNode, nextNode)
			if arcWeight < 0 {
				err := erx.NewError("Negative weight detected")
				err.AddV("head", curNode)
				err.AddV("tail", nextNode)
				err.AddV("weight", arcWeight)
				panic(err)
			}
			nextWeight := curWeight + arcWeight
			if stopFunc==nil || !stopFunc(nextNode, nextWeight) {
				q.Add(nextNode, -nextWeight)
			}
		}
	}
	
	return false
}

// Get all mixed paths from one node to another
//
// This algorithms doesn't take any loops into paths. So maximum path length is 
// nodes count in graph.
func GetAllMixedPaths(gr MixedGraphReader, from, to NodeId) <-chan []NodeId {
	curPath := make([]NodeId, gr.NodesCnt())
	nodesStatus := make(map[NodeId]bool)
	ch := make(chan []NodeId)
	go getAllMixedPaths_helper(gr, from, to, curPath, 0, nodesStatus, ch, true)
	return ch
}

func getAllMixedPaths_helper(gr MixedGraphReader, from, to NodeId, curPath []NodeId, pathPos int, nodesStatus map[NodeId]bool, ch chan []NodeId, closeChannel bool) {
	if _, ok := nodesStatus[from]; ok {
		return
	}
	curPath[pathPos] = from

	if from==to { 
		if pathPos>0 {
			pathCopy := make([]NodeId, pathPos+1)
			copy(pathCopy, curPath[0:pathPos+1])
			ch <- pathCopy
		}
		return
	}
	nodesStatus[from] = true
	
	for _, neighbour := range gr.GetNeighbours(from) {
		getAllMixedPaths_helper(gr, neighbour, to, curPath, pathPos+1, nodesStatus, ch, false)
	}
	for _, accessor := range gr.GetAccessors(from) {
		getAllMixedPaths_helper(gr, accessor, to, curPath, pathPos+1, nodesStatus, ch, false)
	}
	
	nodesStatus[from] = false, false
	
	if closeChannel {
		close(ch)
	}
	return
}
