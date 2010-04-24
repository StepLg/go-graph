package graph

import (
	"erx"
	. "exp/iterable"
)

type arrowsIterable struct {
	arrows DirectedArrowsIterable
}

func (ai arrowsIterable) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for arr := range ai.arrows.ArrowsIter() {
			ch <- arr
		}
	}()
	return ch
}

func ArrowsToGenericIter(arrIter DirectedArrowsIterable) Iterable {
	return arrowsIterable{arrIter}
}

func CopyDirectedGraph(arrIter DirectedArrowsIterable, gr DirectedGraph) erx.Error {
	wheel := erx.NewError("Can't copy directed graph")
	for arrow := range arrIter.ArrowsIter() {
		err := gr.AddArrow(arrow.From, arrow.To)
		if err!=nil {
			wheel.AddE(err)
		}
	}
	
	if wheel.Errors().Len()>0 {
		return wheel
	}
	
	return nil
}
