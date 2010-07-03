package graph

import (
	. "exp/iterable"
)	

// Extract all vertexes, which are accessible from given node.
type OutNeighboursExtractor interface {
	GetOutNeighbours(node VertexId) VertexesIterable
}

type dgraphOutNeighboursExtractor struct {
	dgraph DirectedGraphArcsReader
}

func (e *dgraphOutNeighboursExtractor) GetOutNeighbours(node VertexId) VertexesIterable {
	return e.dgraph.GetAccessors(node)
}

// Extract all vertexes, accessible from given node in directed graph.
//
// In fact, interface function is a synonym to DirectedGraphArcsReader.GetAccessors() function
func NewDgraphOutNeighboursExtractor(gr DirectedGraphArcsReader) OutNeighboursExtractor {
	return OutNeighboursExtractor(&dgraphOutNeighboursExtractor{dgraph:gr})
}

type ugraphOutNeighboursExtractor struct {
	ugraph UndirectedGraphEdgesReader
}

func (e *ugraphOutNeighboursExtractor) GetOutNeighbours(node VertexId) VertexesIterable {
	return e.ugraph.GetNeighbours(node)
}

// Extract all vertexes, accessible from given node in undirected graph.
//
// In fact, interface function is a synonym to UndirectedGraphEdgesReader.GetNeighbours() function
func NewUgraphOutNeighboursExtractor(gr UndirectedGraphEdgesReader) OutNeighboursExtractor {
	return OutNeighboursExtractor(&ugraphOutNeighboursExtractor{ugraph:gr})
}

type mgraphOutNeighboursExtractor struct {
	mgraph MixedGraphConnectionsReader
}

func (e *mgraphOutNeighboursExtractor) GetOutNeighbours(node VertexId) VertexesIterable {
	return GenericToVertexesIter(Chain(&[...]Iterable{
		VertexesToGenericIter(e.mgraph.GetAccessors(node)), 
		VertexesToGenericIter(e.mgraph.GetNeighbours(node)),
	}))
}

// Extract all vertexes, accessible from given node in mixed graph.
//
// In fact, interface function is chain of MixedGraphConnectionsReader.GetAccessors()
// and MixedGraphConnectionsReader.GetNeighbours() functions.
func NewMgraphOutNeighboursExtractor(gr MixedGraphConnectionsReader) OutNeighboursExtractor {
	return OutNeighboursExtractor(&mgraphOutNeighboursExtractor{mgraph:gr})
}

////////////////////////////////////////////////////////////////////////////////


// Extract all vertexes, from which accessibe given node.
type InNeighboursExtractor interface {
	GetInNeighbours(node VertexId) VertexesIterable
}

type dgraphInNeighboursExtractor struct {
	dgraph DirectedGraphArcsReader
}

// Extract all vertexes, from which accessible given node in directed graph.
//
// In fact, interface function is a synonym to DirectedGraphArcsReader.GetPredecessors() function
func (e *dgraphInNeighboursExtractor) GetInNeighbours(node VertexId) VertexesIterable {
	return e.dgraph.GetPredecessors(node)
}

func NewDgraphInNeighboursExtractor(gr DirectedGraphArcsReader) InNeighboursExtractor {
	return InNeighboursExtractor(&dgraphInNeighboursExtractor{dgraph:gr})
}

type ugraphInNeighboursExtractor struct {
	ugraph UndirectedGraphEdgesReader
}

func (e *ugraphInNeighboursExtractor) GetInNeighbours(node VertexId) VertexesIterable {
	return e.ugraph.GetNeighbours(node)
}

// Extract all vertexes, from which accessible given node in undirected graph.
//
// In fact, interface function is a synonym to UndirectedGraphEdgesReader.GetNeighbours() function
func NewUgraphInNeighboursExtractor(gr UndirectedGraphEdgesReader) InNeighboursExtractor {
	return InNeighboursExtractor(&ugraphInNeighboursExtractor{ugraph:gr})
}

type mgraphInNeighboursExtractor struct {
	mgraph MixedGraphConnectionsReader
}

func (e *mgraphInNeighboursExtractor) GetInNeighbours(node VertexId) VertexesIterable {
	return GenericToVertexesIter(Chain(&[...]Iterable{
		VertexesToGenericIter(e.mgraph.GetPredecessors(node)), 
		VertexesToGenericIter(e.mgraph.GetNeighbours(node)),
	}))
}

// Extract all vertexes, accessible from given node in mixed graph.
//
// In fact, interface function is chain of MixedGraphConnectionsReader.GetPredecessors()
// and MixedGraphConnectionsReader.GetNeighbours() functions.
func NewMgraphInNeighboursExtractor(gr MixedGraphConnectionsReader) InNeighboursExtractor {
	return InNeighboursExtractor(&mgraphInNeighboursExtractor{mgraph:gr})
}
