package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)


// Interface for CheckPath functions
//
// Require GraphNodesReader because of very unefficient nodesPriorityQueue realization
type ICheckPathGraph interface {
	DirectedGraphReader
	GraphNodesReader
}

type ConnectionWeightFunc func(head, tail NodeId) float

func SimpleWeightFunc(head, tail NodeId) float {
	return 1.0
}

type CheckPath func(gr ICheckPathGraph, from, to NodeId, maxWeight float, weightFunction ConnectionWeightFunc) bool

func CheckPathDijkstra(gr ICheckPathGraph, from, to NodeId, maxWeight float, weightFunction ConnectionWeightFunc) bool {
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
			if nextWeight < maxWeight || maxWeight < 0 {
				q.Add(nextNode, nextWeight)
			}
		}
	}
	
	return false
}
