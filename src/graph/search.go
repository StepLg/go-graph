package graph

import (
	. "exp/iterable"

	"github.com/StepLg/go-erx/src/erx"
)


type ConnectionWeightFunc func(head, tail NodeId) float

type StopFunc func(node NodeId, sumWeight float) bool

func SimpleWeightFunc(head, tail NodeId) float {
	return 1.0
}

type AllNeighboursExtractor interface {
	GetAllNeighbours(node NodeId) NodesIterable
}

type allDirectedNeighboursExtractor struct {
	dgraph DirectedGraphArcsReader
}

func (e *allDirectedNeighboursExtractor) GetAllNeighbours(node NodeId) NodesIterable {
	return e.dgraph.GetAccessors(node)
}

func NewDirectedNeighboursExtractor(gr DirectedGraphArcsReader) AllNeighboursExtractor {
	return AllNeighboursExtractor(&allDirectedNeighboursExtractor{dgraph:gr})
}

type allUndirectedNeighboursExtractor struct {
	ugraph UndirectedGraphEdgesReader
}

func (e *allUndirectedNeighboursExtractor) GetAllNeighbours(node NodeId) NodesIterable {
	return e.ugraph.GetNeighbours(node)
}

func NewUndirectedNeighboursExtractor(gr UndirectedGraphEdgesReader) AllNeighboursExtractor {
	return AllNeighboursExtractor(&allUndirectedNeighboursExtractor{ugraph:gr})
}

type allMixedNeighboursExtractor struct {
	mgraph MixedGraphConnectionsReader
}

func (e *allMixedNeighboursExtractor) GetAllNeighbours(node NodeId) NodesIterable {
	return GenericToNodesIter(Chain(&[...]Iterable{
		NodesToGenericIter(e.mgraph.GetAccessors(node)), 
		NodesToGenericIter(e.mgraph.GetNeighbours(node)),
	}))
}

func NewMixedNeighboursExtractor(gr MixedGraphConnectionsReader) AllNeighboursExtractor {
	return AllNeighboursExtractor(&allMixedNeighboursExtractor{mgraph:gr})
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
func CheckPathDijkstra(neighboursExtractor AllNeighboursExtractor, from, to NodeId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) (float, bool) {
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
	
		for nextNode := range neighboursExtractor.GetAllNeighbours(curNode).NodesIter() {
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

type CheckDirectedPath func(gr DirectedGraphArcsReader, from, to NodeId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool

func CheckDirectedPathDijkstra(gr DirectedGraphArcsReader, from, to NodeId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool {
	_, pathExists := CheckPathDijkstra(NewDirectedNeighboursExtractor(gr), from, to, stopFunc, weightFunction)
	return pathExists
}

type CheckUndirectedPath func(gr UndirectedGraphEdgesReader, from, to NodeId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool

func CheckUndirectedPathDijkstra(gr UndirectedGraphEdgesReader, from, to NodeId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool {
	_, pathExists := CheckPathDijkstra(NewUndirectedNeighboursExtractor(gr), from, to, stopFunc, weightFunction)
	return pathExists
}

type CheckMixedPath func(gr MixedGraphConnectionsReader, from, to NodeId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool

func CheckMixedPathDijkstra(gr MixedGraphConnectionsReader, from, to NodeId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool {
	_, pathExists := CheckPathDijkstra(NewMixedNeighboursExtractor(gr), from, to, stopFunc, weightFunction)
	return pathExists
}

// Get all paths from one node to another
//
// This algorithms doesn't take any loops into paths.
func GetAllPaths(neighboursExtractor AllNeighboursExtractor, from, to NodeId) <-chan []NodeId {
	curPath := make([]NodeId, 10)
	nodesStatus := make(map[NodeId]bool)
	ch := make(chan []NodeId)
	go getAllPaths_helper(neighboursExtractor, from, to, curPath, 0, nodesStatus, ch, true)
	return ch
}

func getAllPaths_helper(neighboursExtractor AllNeighboursExtractor, from, to NodeId, curPath []NodeId, pathPos int, nodesStatus map[NodeId]bool, ch chan []NodeId, closeChannel bool) {
	if _, ok := nodesStatus[from]; ok {
		return
	}
	if pathPos==len(curPath) {
		// reallocate curPath slice to add new elements
		tmp := make([]NodeId, 2*pathPos)
		copy(tmp, curPath)
		curPath = tmp
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
	
	for nextNode := range neighboursExtractor.GetAllNeighbours(from).NodesIter() {
		getAllPaths_helper(neighboursExtractor, nextNode, to, curPath, pathPos+1, nodesStatus, ch, false)
	}
	
	nodesStatus[from] = false, false
	
	if closeChannel {
		close(ch)
	}
	return
}
