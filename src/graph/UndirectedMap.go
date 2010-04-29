package graph

import (
	"erx"
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
// UndirectedGraphNodesWriter

// Adding single node to graph
func (g *UndirectedMap) AddNode(node NodeId) erx.Error {
	var err erx.Error
	if _, ok := g.edges[node]; ok {
		err = erx.NewError("Node already exists.")
		goto Error
	}
	
	g.edges[node] = make(map[NodeId]bool)
	
	return nil
	Error:
	err = erx.NewSequent("Can't add single node to undirected graph.", err)
	err.AddV("node id", node)
	return err
}

///////////////////////////////////////////////////////////////////////////////
// GraphNodesRemover

func (g *UndirectedMap) RemoveNode(node NodeId) erx.Error {
	var err erx.Error
	if _, ok := g.edges[node]; !ok {
		err = erx.NewError("Node doesn't exist.")
		goto Error
	}
	
	g.edges[node] = nil, false
	for _, connectedNodes := range g.edges {
		connectedNodes[node] = false, false
	}
	
	return nil
	Error:
	err = erx.NewSequent("Can't remove node from undirected graph.", err)
	err.AddV("node id", node)
	return err
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphEdgesWriter

func (g *UndirectedMap) touchNode(node NodeId) {
	if _, ok := g.edges[node]; !ok {
		g.edges[node] = make(map[NodeId]bool)
	}
}

// Adding arrow to graph.
func (g *UndirectedMap) AddEdge(from, to NodeId) (err erx.Error) {
	g.touchNode(from)
	g.touchNode(to)
	
	if direction, ok := g.edges[from][to]; ok && direction {
		err = erx.NewError("Duplicate edge.")
		goto Error
	}
	
	g.edges[from][to] = true
	g.edges[to][from] = true
	g.edgesCnt++	
	return
	
	Error:
	err = erx.NewSequent("Can't add edge.", err)
	err.AddV("from", from)
	err.AddV("to", to)
	return
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphEdgesRemover

// Removing arrow  'from' and 'to' nodes
func (g *UndirectedMap) RemoveEdge(from, to NodeId) (err erx.Error) {
	connectedNodes, ok := g.edges[from]
	if !ok {
		err = erx.NewError("From node doesn't exist.")
		goto Error
	}
	
	if _, ok = connectedNodes[to]; ok {
		err = erx.NewError("To node doesn't exist.")
		goto Error
	}
	
	g.edges[from][to] = false, false
	g.edges[to][from] = false, false
	g.edgesCnt--
	return nil
	
	Error:
	err = erx.NewSequent("Can't remove arrow.", err)
	err.AddV("from", from)
	err.AddV("to", to)
	return err
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
func (g *UndirectedMap) GetNeighbours(node NodeId) (connected Nodes, err erx.Error) {
	if connectedMap, ok := g.edges[node]; ok {
		connected = make(Nodes, len(connectedMap))
		id := 0
		for nodeId, _ := range connectedMap {
			connected[id] = nodeId
			id++
		}
	} else {
		err = erx.NewError("Node doesn't exists.")
	}
	
	if err!=nil {
		err = erx.NewSequent("Can't get node neighbours.", err)
		err.AddV("node", node)
	}
	
	return
}

func (g *UndirectedMap) CheckEdge(from, to NodeId) (isExist bool, err erx.Error) {
	connectedNodes, ok := g.edges[from]
	if !ok {
		err = erx.NewError("From node doesn't exist.")
		goto Error
	}
	
	if _, ok = g.edges[to]; !ok {
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
