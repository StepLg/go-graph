package graph

import (
	"fmt"
	. "exp/iterable"

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
		for nextNode := range gr.GetAccessors(curNode).NodesIter() {
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
		
		iter := Chain(&[...]Iterable{
			NodesToGenericIter(gr.GetAccessors(curNode)), 
			NodesToGenericIter(gr.GetNeighbours(curNode)),
		})

		for nextNode := range iter.Iter() {
			nextNode, ok := nextNode.(NodeId)
			if !ok {
				err := erx.NewError("Generics type assertation.")
				err.AddV("expected type", "graph.NodeId")
				err.AddV("got type", fmt.Sprintf("%T", nextNode))
				panic(err)
			}
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
	
	iter := Chain(&[...]Iterable{
		NodesToGenericIter(gr.GetAccessors(from)), 
		NodesToGenericIter(gr.GetNeighbours(from)),
	})
	
	for nextNode := range iter.Iter() {
		nextNode, ok := nextNode.(NodeId)
		if !ok {
			err := erx.NewError("Generics type assertation.")
			err.AddV("expected type", "graph.NodeId")
			err.AddV("got type", fmt.Sprintf("%T", nextNode))
			panic(err)
		}
		getAllMixedPaths_helper(gr, nextNode, to, curPath, pathPos+1, nodesStatus, ch, false)
	}
	
	nodesStatus[from] = false, false
	
	if closeChannel {
		close(ch)
	}
	return
}
