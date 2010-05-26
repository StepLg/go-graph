package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)

type UndirectedMap struct {
	edges map[NodeId]map[NodeId]bool
	edgesCnt int
}

func NewUndirectedMap() *UndirectedMap {
	g := new(UndirectedMap)
	g.edges = make(map[NodeId]map[NodeId]bool)
	g.edgesCnt = 0
	return g
}

///////////////////////////////////////////////////////////////////////////////
// ConnectionsIterable

func (g *UndirectedMap) ConnectionsIter() <-chan Connection {
	return g.EdgesIter()
}

///////////////////////////////////////////////////////////////////////////////
// NodesIterable

func (g *UndirectedMap) NodesIter() <-chan NodeId {
	ch := make(chan NodeId)
	go func() {
		for from, _ := range g.edges {
			ch <- from
		}
		close(ch)
	}()
	return ch
}

///////////////////////////////////////////////////////////////////////////////
// NodesChecker

func (g *UndirectedMap) CheckNode(node NodeId) (exists bool) {
	_, exists = g.edges[node]
	return
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphNodesWriter

// Adding single node to graph
func (g *UndirectedMap) AddNode(node NodeId) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Add node to graph.", err, 1)
		res.AddV("node id", node)
		return
	}

	if _, ok := g.edges[node]; ok {
		panic(makeError(erx.NewError("Node already exists.")))
	}
	
	g.edges[node] = make(map[NodeId]bool)

	return	
}

///////////////////////////////////////////////////////////////////////////////
// GraphNodesRemover

func (g *UndirectedMap) RemoveNode(node NodeId) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Remove node from graph.", err, 1)
		res.AddV("node id", node)
		return
	}

	if _, ok := g.edges[node]; !ok {
		panic(makeError(erx.NewError("Node doesn't exist.")))
	}
	
	g.edges[node] = nil, false
	for _, connectedNodes := range g.edges {
		connectedNodes[node] = false, false
	}
	
	return
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphEdgesWriter

func (g *UndirectedMap) touchNode(node NodeId) {
	if _, ok := g.edges[node]; !ok {
		g.edges[node] = make(map[NodeId]bool)
	}
}

// Adding arrow to graph.
func (g *UndirectedMap) AddEdge(from, to NodeId) {
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
func (g *UndirectedMap) RemoveEdge(from, to NodeId) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Remove edge from graph.", err, 1)
		res.AddV("node 1", from)
		res.AddV("node 2", to)
		return
	}
	connectedNodes, ok := g.edges[from]
	if !ok {
		panic(makeError(erx.NewError("First node doesn't exists")))
	}
	
	if _, ok = connectedNodes[to]; ok {
		panic(makeError(erx.NewError("Second node doesn't exists")))
	}
	
	g.edges[from][to] = false, false
	g.edges[to][from] = false, false
	g.edgesCnt--

	return
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphReader

func (g *UndirectedMap) NodesCnt() int {
	return len(g.edges)
}

func (g *UndirectedMap) EdgesCnt() int {
	return g.edgesCnt
}

// Getting node predecessors
func (g *UndirectedMap) GetNeighbours(node NodeId) NodesIterable {
	iterator := func() <-chan NodeId {
		ch := make(chan NodeId)
		go func() {
			if connectedMap, ok := g.edges[node]; ok {
				for nodeId, _ := range connectedMap {
					ch <- nodeId
				}
			} else {
				panic(erx.NewError("Node doesn't exists."))
			}
			close(ch)
		}()
		return ch
	}
	
	return NodesIterable(&nodesIterableLambdaHelper{iterFunc:iterator})
}

func (g *UndirectedMap) CheckEdge(from, to NodeId) (isExist bool) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Check edge existance in graph.", err, 1)
		res.AddV("node 1", from)
		res.AddV("node 2", to)
		return
	}

	connectedNodes, ok := g.edges[from]
	if !ok {
		panic(makeError(erx.NewError("Fist node doesn't exist.")))
	}
	
	if _, ok = g.edges[to]; !ok {
		panic(makeError(erx.NewError("Second node doesn't exist.")))
	}
	
	_, isExist = connectedNodes[to]
	
	return
}

func (g *UndirectedMap) EdgesIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for from, connectedNodes := range g.edges {
			for to, _ := range connectedNodes {
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
