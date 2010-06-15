package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)

type DirectedMap struct {
	directArcs map[VertexId]map[VertexId]bool
	reversedArcs map[VertexId]map[VertexId]bool
	arcsCnt int
}

func NewDirectedMap() *DirectedMap {
	g := new(DirectedMap)
	g.directArcs = make(map[VertexId]map[VertexId]bool)
	g.reversedArcs = make(map[VertexId]map[VertexId]bool)
	g.arcsCnt = 0
	return g
}

///////////////////////////////////////////////////////////////////////////////
// ConnectionsIterable

func (g *DirectedMap) ConnectionsIter() <-chan Connection {
	return g.ArcsIter()
}

///////////////////////////////////////////////////////////////////////////////
// NodesIterable

func (g *DirectedMap) NodesIter() <-chan VertexId {
	ch := make(chan VertexId)
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
// NodesChecker

func (g *DirectedMap) CheckNode(node VertexId) (exists bool) {
	_, exists = g.directArcs[node]
	return
}

///////////////////////////////////////////////////////////////////////////////
// GraphNodesWriter

// Adding single node to graph
func (g *DirectedMap) AddNode(node VertexId) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Add node to graph.", err, 1)
		res.AddV("node id", node)
		return
	}

	if _, ok := g.directArcs[node]; ok {
		panic(makeError(erx.NewError("Node already exists.")))
	}
	
	g.directArcs[node] = make(map[VertexId]bool)
	g.reversedArcs[node] = make(map[VertexId]bool)

	return	
}

///////////////////////////////////////////////////////////////////////////////
// GraphNodesRemover

func (g *DirectedMap) RemoveNode(node VertexId) {
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

func (g *DirectedMap) touchNode(node VertexId) {
	if _, ok := g.directArcs[node]; !ok {
		g.directArcs[node] = make(map[VertexId]bool)
		g.reversedArcs[node] = make(map[VertexId]bool)
	}
}

// Adding arrow to graph.
func (g *DirectedMap) AddArc(from, to VertexId) {
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
func (g *DirectedMap) RemoveArc(from, to VertexId) {
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
func (g *DirectedMap) GetSources() NodesIterable {
	iterator := func() <-chan VertexId {
		ch := make(chan VertexId)
		go func() {
			for VertexId, predecessors := range g.reversedArcs {
				if len(predecessors)==0 {
					ch <- VertexId
				}
			}
			close(ch)
		}()
		return ch
	}
	
	return NodesIterable(&nodesIterableLambdaHelper{iterFunc:iterator})
}

// Getting all graph sinks.
func (g *DirectedMap) GetSinks() NodesIterable {
	iterator := func() <-chan VertexId {
		ch := make(chan VertexId)
		go func() {
			for VertexId, accessors := range g.directArcs {
				if len(accessors)==0 {
					ch <- VertexId
				}
			}
			close(ch)
		}()
		return ch
	}
	
	return NodesIterable(&nodesIterableLambdaHelper{iterFunc:iterator})
}

// Getting node accessors
func (g *DirectedMap) GetAccessors(node VertexId) NodesIterable {
	iterator := func() <-chan VertexId {
		ch := make(chan VertexId)
		
		go func() {
			defer func() {
				if e := recover(); e!=nil {
					err := erx.NewSequent("Get node accessors in mixed graph.", e)
					err.AddV("node", node)
					panic(err)
				}
			}()
			accessorsMap, ok := g.directArcs[node]
			if !ok {
				panic(erx.NewError("Node doesn't exists."))
			}
			
			for VertexId, _ := range accessorsMap {
				ch <- VertexId
			}
			close(ch)
		}()
			
		return ch
	}
	
	return NodesIterable(&nodesIterableLambdaHelper{iterFunc:iterator})
}

// Getting node predecessors
func (g *DirectedMap) GetPredecessors(node VertexId) NodesIterable {
	iterator := func() <-chan VertexId {
		ch := make(chan VertexId)
		
		go func() {
			defer func() {
				if e := recover(); e!=nil {
					err := erx.NewSequent("Get node accessors in mixed graph.", e)
					err.AddV("node", node)
					panic(err)
				}
			}()
			accessorsMap, ok := g.reversedArcs[node]
			if !ok {
				panic(erx.NewError("Node doesn't exists."))
			}
			
			for VertexId, _ := range accessorsMap {
				ch <- VertexId
			}
			close(ch)
		}()
			
		return ch
	}
	
	return NodesIterable(&nodesIterableLambdaHelper{iterFunc:iterator})
}

func (g *DirectedMap) CheckArc(from, to VertexId) (isExist bool) {
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

func (g *DirectedMap) ArcsIter() <-chan Connection {
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
