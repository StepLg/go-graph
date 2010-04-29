package graph

import (
	"github.com/StepLg/go-erx/src/erx"
	. "container/vector"
)

type DirectedMap struct {
	directArcs map[NodeId]map[NodeId]bool
	reversedArcs map[NodeId]map[NodeId]bool
	arcsCnt int
}

func NewDirectedMap() *DirectedMap {
	g := new(DirectedMap)
	g.directArcs = make(map[NodeId]map[NodeId]bool)
	g.reversedArcs = make(map[NodeId]map[NodeId]bool)
	g.arcsCnt = 0
	return g
}

///////////////////////////////////////////////////////////////////////////////
// ConnectionsIterable

func (g *DirectedMap) ConnectionsIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for from, connectedNodes := range g.directArcs {
			for to, _ := range connectedNodes {
				ch <- Connection{from, to}
			}
		}
		close(ch)
	}()
	return ch
}

///////////////////////////////////////////////////////////////////////////////
// NodesIterable

func (g *DirectedMap) NodesIter() <-chan NodeId {
	ch := make(chan NodeId)
	go func() {
		for from, _ := range g.directArcs {
			ch <- from
		}
		
		for to, _ := range g.reversedArcs {
			// need to prevent duplicating node ids
			if _, ok := g.directArcs[to]; !ok {
				ch <- to
			}
		}
		close(ch)
	}()
	return ch
}

///////////////////////////////////////////////////////////////////////////////
// GraphNodesWriter

// Adding single node to graph
func (g *DirectedMap) AddNode(node NodeId) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Add node to graph.", err, 1)
		res.AddV("node id", node)
		return
	}

	if _, ok := g.directArcs[node]; ok {
		panic(makeError(erx.NewError("Node already exists.")))
	}
	
	g.directArcs[node] = make(map[NodeId]bool)
	g.reversedArcs[node] = make(map[NodeId]bool)

	return	
}

///////////////////////////////////////////////////////////////////////////////
// GraphNodesRemover

func (g *DirectedMap) RemoveNode(node NodeId) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Remove node from graph.", err, 1)
		res.AddV("node id", node)
		return
	}

	_, okDirect := g.directArcs[node]
	_, okReversed := g.reversedArcs[node]
	if !okDirect && !okReversed {
		panic(makeError(erx.NewError("Node doesn't exist.")))
	}
	
	g.directArcs[node] = nil, false
	g.reversedArcs[node] = nil, false
	for _, connectedNodes := range g.directArcs {
		connectedNodes[node] = false, false
	}
	for _, connectedNodes := range g.reversedArcs {
		connectedNodes[node] = false, false
	}
	return
}

///////////////////////////////////////////////////////////////////////////////
// DirectedGraphArcsWriter

func (g *DirectedMap) touchNode(node NodeId) {
	if _, ok := g.directArcs[node]; !ok {
		g.directArcs[node] = make(map[NodeId]bool)
		g.reversedArcs[node] = make(map[NodeId]bool)
	}
}

// Adding arrow to graph.
func (g *DirectedMap) AddArc(from, to NodeId) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Add arc to graph.", err, 1)
		res.AddV("tail", from)
		res.AddV("head", to)
		return
	}

	g.touchNode(from)
	g.touchNode(to)
	
	if direction, ok := g.directArcs[from][to]; ok && direction {
		panic(makeError(erx.NewError("Duplicate arrow.")))
	}
	
	g.directArcs[from][to] = true
	g.reversedArcs[to][from] = true
	g.arcsCnt++
	return	
}

///////////////////////////////////////////////////////////////////////////////
// DirectedGraphArcsRemover

// Removing arrow  'from' and 'to' nodes
func (g *DirectedMap) RemoveArc(from, to NodeId) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Remove arc from graph.", err, 1)
		res.AddV("tail", from)
		res.AddV("head", to)
		return
	}

	connectedNodes, ok := g.directArcs[from]
	if !ok {
		panic(makeError(erx.NewError("Tail node doesn't exist.")))
	}
	
	if _, ok = connectedNodes[to]; ok {
		panic(makeError(erx.NewError("Head node doesn't exist.")))
	}
	
	g.directArcs[from][to] = false, false
	g.reversedArcs[to][from] = false, false
	g.arcsCnt--
	
	return
}

///////////////////////////////////////////////////////////////////////////////
// DirectedGraphReader

func (g *DirectedMap) NodesCnt() int {
	return len(g.directArcs)
}

func (g *DirectedMap) ArcsCnt() int {
	return g.arcsCnt
}

// Getting all graph sources.
func (g *DirectedMap) GetSources() (result Nodes) {
	resultVector := new(Vector)
	for nodeId, predecessors := range g.reversedArcs {
		if len(predecessors)==0 {
			resultVector.Push(nodeId)
		}
	}

	result = make(Nodes, resultVector.Len())
	for i:=0; i<resultVector.Len(); i++ {
		
		if nId, ok := resultVector.At(i).(NodeId); ok {
			// must allways be true! lack of generics :(
			result[i] = nId
		}
	}	
	return
}

// Getting all graph sinks.
func (g *DirectedMap) GetSinks() (result Nodes) {
	resultVector := new(Vector)
	for nodeId, accessors := range g.directArcs {
		if len(accessors)==0 {
			resultVector.Push(nodeId)
		}
	}

	result = make(Nodes, resultVector.Len())
	for i:=0; i<resultVector.Len(); i++ {
		
		if nId, ok := resultVector.At(i).(NodeId); ok {
			// must allways be true! lack of generics :(
			result[i] = nId
		}
	}	
	return
}

// Getting node accessors
func (g *DirectedMap) GetAccessors(node NodeId) (accessors Nodes) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Getting node accessors.", err, 1)
		res.AddV("node id", node)
		return
	}

	accessorsMap, ok := g.directArcs[node]
	if !ok {
		panic(makeError(erx.NewError("Node doesn't exists.")))
	}
	
	accessors = make(Nodes, len(accessorsMap))
	id := 0
	for nodeId, _ := range accessorsMap {
		accessors[id] = nodeId
		id++
	}
	
	return
}

// Getting node predecessors
func (g *DirectedMap) GetPredecessors(node NodeId) (predecessors Nodes) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Getting node predecessors.", err, 1)
		res.AddV("node id", node)
		return
	}

	predecessorsMap, ok := g.reversedArcs[node]
	if !ok {
		panic(makeError(erx.NewError("Node doesn't exists.")))
	}

	predecessors = make(Nodes, len(predecessorsMap))
	id := 0
	for nodeId, _ := range predecessorsMap {
		predecessors[id] = nodeId
		id++
	}
	
	return
}

func (g *DirectedMap) CheckArc(from, to NodeId) (isExist bool) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Checking arc existance in graph.", err, 1)
		res.AddV("tail", from)
		res.AddV("head", to)
		return
	}
	
	connectedNodes, ok := g.directArcs[from]
	if !ok {
		panic(makeError(erx.NewError("From node doesn't exist.")))
	}
	
	if _, ok = g.reversedArcs[to]; !ok {
		panic(makeError(erx.NewError("To node doesn't exist.")))
	}
	
	_, isExist = connectedNodes[to]

	return
}
