package graph

import (
	"erx"
	. "container/vector"
)

type DirectedMap struct {
	directArrows map[NodeId]map[NodeId]bool
	reversedArrows map[NodeId]map[NodeId]bool
	arrowsCnt int
}

func NewDirectedMap() *DirectedMap {
	g := new(DirectedMap)
	g.directArrows = make(map[NodeId]map[NodeId]bool)
	g.reversedArrows = make(map[NodeId]map[NodeId]bool)
	g.arrowsCnt = 0
	return g
}

func (g *DirectedMap) NodesCnt() int {
	return len(g.directArrows)
}

func (g *DirectedMap) ArrowsCnt() int {
	return g.arrowsCnt
}

// Adding single node to graph
func (g *DirectedMap) AddNode(node NodeId) erx.Error {
	var err erx.Error
	if _, ok := g.directArrows[node]; ok {
		err = erx.NewError("Node already exists.")
		goto Error
	}
	
	g.directArrows[node] = make(map[NodeId]bool)
	g.reversedArrows[node] = make(map[NodeId]bool)
	
	return nil
	Error:
	err = erx.NewSequent("", err)
	err.AddV("node id", node)
	return err
}

func (g *DirectedMap) touchNode(node NodeId) {
	if _, ok := g.directArrows[node]; !ok {
		g.directArrows[node] = make(map[NodeId]bool)
		g.reversedArrows[node] = make(map[NodeId]bool)
	}
}

// Adding arrow to graph.
func (g *DirectedMap) AddArrow(from, to NodeId) (err erx.Error) {
	g.touchNode(from)
	g.touchNode(to)
	
	if direction, ok := g.directArrows[from][to]; ok && direction {
		err = erx.NewError("Duplicate arrow.")
		goto Error
	}
	
	g.directArrows[from][to] = true
	g.reversedArrows[to][from] = true
	g.arrowsCnt++	
	return
	
	Error:
	err = erx.NewSequent("Can't add arrow", err)
	err.AddV("from", from)
	err.AddV("to", to)
	return
}

// Removing arrow  'from' and 'to' nodes
func (g *DirectedMap) RemoveArrow(from, to NodeId) (err erx.Error) {
	connectedNodes, ok := g.directArrows[from]
	if !ok {
		err = erx.NewError("From node doesn't exist.")
		goto Error
	}
	
	if _, ok = connectedNodes[to]; ok {
		err = erx.NewError("To node doesn't exist.")
		goto Error
	}
	
	g.directArrows[from][to] = false, false
	g.reversedArrows[to][from] = false, false
	g.arrowsCnt--
	return nil
	
	Error:
	err = erx.NewSequent("Can't remove arrow.", err)
	err.AddV("from", from)
	err.AddV("to", to)
	return err
}

// Getting all graph sources.
func (g *DirectedMap) GetSources() (result Nodes, err erx.Error) {
	resultVector := new(Vector)
	for nodeId, predecessors := range g.reversedArrows {
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
	for nodeId, accessors := range g.directArrows {
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
	if accessorsMap, ok := g.directArrows[node]; ok {
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
	if predecessorsMap, ok := g.reversedArrows[node]; ok {
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

func (g *DirectedMap) CheckArrow(from, to NodeId) (isExist bool, err erx.Error) {
	connectedNodes, ok := g.directArrows[from]
	if !ok {
		err = erx.NewError("From node doesn't exist.")
		goto Error
	}
	
	if _, ok = g.reversedArrows[to]; !ok {
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

func (g *DirectedMap) ArrowsIter() <-chan Arrow {
	ch := make(chan Arrow)
	go func() {
		for from, connectedNodes := range g.directArrows {
			for to, _ := range connectedNodes {
				ch <- Arrow{from, to}
			}
		}
		close(ch)
	}()
	return ch
}
