package graph

import (
	"erx"
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
func (g *DirectedMap) AddNode(node NodeId) erx.Error {
	var err erx.Error
	if _, ok := g.directArcs[node]; ok {
		err = erx.NewError("Node already exists.")
		goto Error
	}
	
	g.directArcs[node] = make(map[NodeId]bool)
	g.reversedArcs[node] = make(map[NodeId]bool)
	
	return nil
	Error:
	err = erx.NewSequent("", err)
	err.AddV("node id", node)
	return err
}

///////////////////////////////////////////////////////////////////////////////
// GraphNodesRemover

func (g *DirectedMap) RemoveNode(node NodeId) erx.Error {
	var err erx.Error
	_, okDirect := g.directArcs[node]
	_, okReversed := g.reversedArcs[node]
	if !okDirect && !okReversed {
		err = erx.NewError("Node doesn't exist.")
		goto Error
	}
	
	g.directArcs[node] = nil, false
	g.reversedArcs[node] = nil, false
	for _, connectedNodes := range g.directArcs {
		connectedNodes[node] = false, false
	}
	for _, connectedNodes := range g.reversedArcs {
		connectedNodes[node] = false, false
	}
	
	return nil
	Error:
	err = erx.NewSequent("Can't remove node from undirected graph.", err)
	err.AddV("node id", node)
	return err
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
func (g *DirectedMap) AddArc(from, to NodeId) (err erx.Error) {
	g.touchNode(from)
	g.touchNode(to)
	
	if direction, ok := g.directArcs[from][to]; ok && direction {
		err = erx.NewError("Duplicate arrow.")
		goto Error
	}
	
	g.directArcs[from][to] = true
	g.reversedArcs[to][from] = true
	g.arcsCnt++	
	return
	
	Error:
	err = erx.NewSequent("Can't add arrow", err)
	err.AddV("from", from)
	err.AddV("to", to)
	return
}

///////////////////////////////////////////////////////////////////////////////
// DirectedGraphArcsRemover

// Removing arrow  'from' and 'to' nodes
func (g *DirectedMap) RemoveArc(from, to NodeId) (err erx.Error) {
	connectedNodes, ok := g.directArcs[from]
	if !ok {
		err = erx.NewError("From node doesn't exist.")
		goto Error
	}
	
	if _, ok = connectedNodes[to]; ok {
		err = erx.NewError("To node doesn't exist.")
		goto Error
	}
	
	g.directArcs[from][to] = false, false
	g.reversedArcs[to][from] = false, false
	g.arcsCnt--
	return nil
	
	Error:
	err = erx.NewSequent("Can't remove arrow.", err)
	err.AddV("from", from)
	err.AddV("to", to)
	return err
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
func (g *DirectedMap) GetSources() (result Nodes, err erx.Error) {
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
func (g *DirectedMap) GetSinks() (result Nodes, err erx.Error) {
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
func (g *DirectedMap) GetAccessors(node NodeId) (accessors Nodes, err erx.Error) {
	if accessorsMap, ok := g.directArcs[node]; ok {
		accessors = make(Nodes, len(accessorsMap))
		id := 0
		for nodeId, _ := range accessorsMap {
			accessors[id] = nodeId
			id++
		}
	} else {
		err = erx.NewError("Node doesn't exists.")
	}
	
	if err!=nil {
		err = erx.NewSequent("Can't get node accessors.", err)
		err.AddV("node", node)
	}
	
	return
}

// Getting node predecessors
func (g *DirectedMap) GetPredecessors(node NodeId) (predecessors Nodes, err erx.Error) {
	if predecessorsMap, ok := g.reversedArcs[node]; ok {
		predecessors = make(Nodes, len(predecessorsMap))
		id := 0
		for nodeId, _ := range predecessorsMap {
			predecessors[id] = nodeId
			id++
		}
	} else {
		err = erx.NewError("Node doesn't exists.")
	}
	
	if err!=nil {
		err = erx.NewSequent("Can't get node predecessors.", err)
		err.AddV("node", node)
	}
	
	return
}

func (g *DirectedMap) CheckArc(from, to NodeId) (isExist bool, err erx.Error) {
	connectedNodes, ok := g.directArcs[from]
	if !ok {
		err = erx.NewError("From node doesn't exist.")
		goto Error
	}
	
	if _, ok = g.reversedArcs[to]; !ok {
		err = erx.NewError("To node doesn't exist.")
		goto Error
	}
	
	_, isExist = connectedNodes[to]
	return
	
	Error:
	err = erx.NewSequent("Can't check arrow existance.", err)
	err.AddV("from", from)
	err.AddV("to", to)
	return
}
