package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)


type ConnectionWeightFunc func(head, tail NodeId) float

type StopFunc func(node NodeId, sumWeight float) bool

func SimpleWeightFunc(head, tail NodeId) float {
	return 1.0
}

type CheckPath func(gr DirectedGraphReader, from, to NodeId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool

func CheckPathDijkstra(gr DirectedGraphReader, from, to NodeId, stopFunc StopFunc, weightFunction ConnectionWeightFunc) bool {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Check path in graph with Dijkstra algorithm", e)
			err.AddV("from", from)
			err.AddV("to", to)
			panic(err)
		}
	}()
	
	if from==to {
		return true
	}
	
	q := newPriorityQueueSimple(gr.NodesCnt())
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
			if stopFunc==nil || stopFunc(nextNode, nextWeight) {
				q.Add(nextNode, -nextWeight)
			}
		}
	}
	
	return false
}
