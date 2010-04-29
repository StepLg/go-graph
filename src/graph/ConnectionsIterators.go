package graph

import (
	. "exp/iterable"
)

type connectionsIterable struct {
	arrows ConnectionsIterable
}

func (ai connectionsIterable) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for arr := range ai.arrows.ConnectionsIter() {
			ch <- arr
		}
	}()
	return ch
}

func ArrowsToGenericIter(connIter ConnectionsIterable) Iterable {
	return connectionsIterable{connIter}
}

func CopyDirectedGraph(connIter ConnectionsIterable, gr DirectedGraphArcsWriter) {
	// wheel := erx.NewError("Can't copy directed graph")
	for arrow := range connIter.ConnectionsIter() {
		gr.AddArc(arrow.Tail, arrow.Head)
	}
	return
}

func BuildDirectedGraph(gr DirectedGraph, connIterable ConnectionsIterable , isCorrectOrder func(Connection) bool) {
	for arr := range connIterable.ConnectionsIter() {
		if isCorrectOrder(arr) {
			gr.AddArc(arr.Tail, arr.Head)
		} else {
			gr.AddArc(arr.Head, arr.Tail)
		}
	}
}
