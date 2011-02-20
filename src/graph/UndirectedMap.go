package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)

type UndirectedMap struct {
	edges map[VertexId]map[VertexId]bool
	edgesCnt int
}

func NewUndirectedMap() *UndirectedMap {
	g := new(UndirectedMap)
	g.edges = make(map[VertexId]map[VertexId]bool)
	g.edgesCnt = 0
	return g
}

///////////////////////////////////////////////////////////////////////////////
// ConnectionsIterable

func (g *UndirectedMap) ConnectionsIter() <-chan Connection {
	return g.EdgesIter()
}

///////////////////////////////////////////////////////////////////////////////
// VertexesIterable

func (g *UndirectedMap) VertexesIter() <-chan VertexId {
	ch := make(chan VertexId)
	go func() {
		for from, _ := range g.edges {
			ch <- from
		}
		close(ch)
	}()
	return ch
}

///////////////////////////////////////////////////////////////////////////////
// VertexesChecker

func (g *UndirectedMap) CheckNode(node VertexId) (exists bool) {
	_, exists = g.edges[node]
	return
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphVertexesWriter

// Adding single node to graph
func (g *UndirectedMap) AddNode(node VertexId) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Add node to graph.", err, 1)
		res.AddV("node id", node)
		return
	}

	if _, ok := g.edges[node]; ok {
		panic(makeError(erx.NewError("Node already exists.")))
	}
	
	g.edges[node] = make(map[VertexId]bool)

	return	
}

///////////////////////////////////////////////////////////////////////////////
// GraphVertexesRemover

func (g *UndirectedMap) RemoveNode(node VertexId) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Remove node from graph.", err, 1)
		res.AddV("node id", node)
		return
	}

	if _, ok := g.edges[node]; !ok {
		panic(makeError(erx.NewError("Node doesn't exist.")))
	}
	
	g.edges[node] = nil, false
	for _, connectedVertexes := range g.edges {
		connectedVertexes[node] = false, false
	}
	
	return
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphEdgesWriter

func (g *UndirectedMap) touchNode(node VertexId) {
	if _, ok := g.edges[node]; !ok {
		g.edges[node] = make(map[VertexId]bool)
	}
}

// Adding arrow to graph.
func (g *UndirectedMap) AddEdge(from, to VertexId) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Add edge to graph.", err, 1)
		res.AddV("node 1", from)
		res.AddV("node 2", to)
		return
	}

	g.touchNode(from)
	g.touchNode(to)
	
	if direction, ok := g.edges[from][to]; ok && direction {
		panic(makeError(erx.NewError("Duplicate arrow.")))
	}
	
	g.edges[from][to] = true
	g.edges[to][from] = true
	g.edgesCnt++	

	return
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphEdgesRemover

// Removing arrow  'from' and 'to' nodes
func (g *UndirectedMap) RemoveEdge(from, to VertexId) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Remove edge from graph.", err, 1)
		res.AddV("node 1", from)
		res.AddV("node 2", to)
		return
	}
	connectedVertexes, ok := g.edges[from]
	if !ok {
		panic(makeError(erx.NewError("First node doesn't exists")))
	}
	
	if _, ok = connectedVertexes[to]; ok {
		panic(makeError(erx.NewError("Second node doesn't exists")))
	}
	
	g.edges[from][to] = false, false
	g.edges[to][from] = false, false
	g.edgesCnt--

	return
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphReader

func (g *UndirectedMap) Order() int {
	return len(g.edges)
}

func (g *UndirectedMap) EdgesCnt() int {
	return g.edgesCnt
}

func (g *UndirectedMap) CheckEdge(from, to VertexId) (isExist bool) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Check edge existance in graph.", err, 1)
		res.AddV("node 1", from)
		res.AddV("node 2", to)
		return
	}

	connectedVertexes, ok := g.edges[from]
	if !ok {
		panic(makeError(erx.NewError("Fist node doesn't exist.")))
	}
	
	if _, ok = g.edges[to]; !ok {
		panic(makeError(erx.NewError("Second node doesn't exist.")))
	}
	
	_, isExist = connectedVertexes[to]
	
	return
}

func (g *UndirectedMap) EdgesIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for from, connectedVertexes := range g.edges {
			for to, _ := range connectedVertexes {
				if from<to {
					// each edge has a duplicate, so we need to 
					// push only one edge to channel
					ch <- Connection{from, to}
				}
			}
		}
		close(ch)
	}()
	return ch
}
