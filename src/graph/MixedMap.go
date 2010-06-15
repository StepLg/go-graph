package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)

type MixedMap struct {
	connections map[VertexId]map[VertexId]MixedConnectionType
	arcsCnt int
	edgesCnt int
}

func NewMixedMap() *MixedMap {
	g := &MixedMap {
		connections: make(map[VertexId]map[VertexId]MixedConnectionType),
		arcsCnt: 0,
		edgesCnt: 0,
	}
	return g
}

///////////////////////////////////////////////////////////////////////////////
// ConnectionsIterable

func (g *MixedMap) ConnectionsIter() <-chan Connection {
	ch := make(chan Connection)
	panic(erx.NewError("Function doesn't implemented yet"))
	return ch
}

///////////////////////////////////////////////////////////////////////////////
// VertexesIterable

func (g *MixedMap) VertexesIter() <-chan VertexId {
	ch := make(chan VertexId)
	go func() {
		for from, _ := range g.connections {
			ch <- from
		}
		close(ch)
	}()
	return ch
}

///////////////////////////////////////////////////////////////////////////////
// VertexesChecker

func (g *MixedMap) CheckNode(node VertexId) (exists bool) {
	_, exists = g.connections[node]
	return
}

///////////////////////////////////////////////////////////////////////////////
// GraphVertexesWriter

// Adding single node to graph
func (g *MixedMap) AddNode(node VertexId) {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Add node to graph.", e)
			err.AddV("node id", node)
			panic(err)
		}
	}()

	if _, ok := g.connections[node]; ok {
		panic(erx.NewError("Node already exists."))
	}
	
	g.connections[node] = make(map[VertexId]MixedConnectionType)

	return
}

///////////////////////////////////////////////////////////////////////////////
// GraphVertexesRemover

func (g *MixedMap) RemoveNode(node VertexId) {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Remove node from graph.", e)
			err.AddV("node id", node)
			panic(err)
		}
	}()

	_, ok := g.connections[node]
	if !ok {
		panic(erx.NewError("Node doesn't exist."))
	}
	
	g.connections[node] = nil, false
	for _, connectedVertexes := range g.connections {
		connectedVertexes[node] = CT_NONE, false
	}
	return
}

///////////////////////////////////////////////////////////////////////////////
// DirectedGraphArcsWriter

func (g *MixedMap) touchNode(node VertexId) {
	if _, ok := g.connections[node]; !ok {
		g.connections[node] = make(map[VertexId]MixedConnectionType)
	}
}

// Adding arrow to graph.
func (g *MixedMap) AddArc(from, to VertexId) {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Add arc to graph.", e)
			err.AddV("tail", from)
			err.AddV("head", to)
			panic(err)
		}
	}()

	g.touchNode(from)
	g.touchNode(to)
	
	if direction, ok := g.connections[from][to]; ok {
		err := erx.NewError("Duplicate connection.")
		err.AddV("connection type", direction)
		panic(err)
	}
	
	g.connections[from][to] = CT_DIRECTED
	g.connections[to][from] = CT_DIRECTED_REVERSED
	g.arcsCnt++
	return	
}

///////////////////////////////////////////////////////////////////////////////
// DirectedGraphArcsRemover

// Removing arrow  'from' and 'to' nodes
func (g *MixedMap) RemoveArc(from, to VertexId) {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Remove arc from graph.", e)
			err.AddV("tail", from)
			err.AddV("head", to)
			panic(err)
		}
	}()

	if _, ok := g.connections[from]; ok {
		panic(erx.NewError("Tail node doesn't exist."))
	}
	
	if _, ok := g.connections[to]; ok {
		panic(erx.NewError("Head node doesn't exist."))
	}
	
	if dir, ok := g.connections[from][to]; !ok || dir!=CT_DIRECTED {
		panic(erx.NewError("Arc doesn't exist."))
	}
	
	g.connections[from][to] = CT_NONE, false
	g.connections[to][from] = CT_NONE, false
	g.arcsCnt--
	
	return
}

///////////////////////////////////////////////////////////////////////////////
// DirectedGraphReader

func (g *MixedMap) Order() int {
	return len(g.connections)
}

func (g *MixedMap) ArcsCnt() int {
	return g.arcsCnt
}

// Getting all graph sources.
func (g *MixedMap) GetSources() VertexesIterable {
	iterator := func() <-chan VertexId {
		ch := make(chan VertexId)
		
		go func() {
			for VertexId, connections := range g.connections {
				isSource := true
				for _, connType := range connections {
					if connType==CT_DIRECTED_REVERSED {
						isSource = false
						break
					}
				}
				if isSource {
					ch <- VertexId
				}
			}

			close(ch)
		}()
		
		return ch
	}
	
	return VertexesIterable(&nodesIterableLambdaHelper{iterFunc:iterator})
}

// Getting all graph sinks.
func (g *MixedMap) GetSinks() VertexesIterable {
	iterator := func() <-chan VertexId {
		ch := make(chan VertexId)
		
		go func() {
			for VertexId, connections := range g.connections {
				isSink := true
				for _, connType := range connections {
					if connType==CT_DIRECTED {
						isSink = false
						break
					}
				}
				if isSink {
					ch <- VertexId
				}
			}

			close(ch)
		}()
		
		return ch
	}
	
	return VertexesIterable(&nodesIterableLambdaHelper{iterFunc:iterator})
}

// Getting node accessors
func (g *MixedMap) GetAccessors(node VertexId) VertexesIterable {
	iterator := func() <-chan VertexId {
		ch := make(chan VertexId)
		
		go func() {
			defer func() {
				if e := recover(); e!=nil {
					err := erx.NewSequent("Getting node accessors.", e)
					err.AddV("node id", node)
					panic(err)
				}
			}()
		
			accessorsMap, ok := g.connections[node]
			if !ok {
				panic(erx.NewError("Node doesn't exists."))
			}
			
			for VertexId, connType := range accessorsMap {
				if connType==CT_DIRECTED {
					ch <- VertexId
				}
			}
			
			close(ch)
		}()
		
		return ch
	}
	
	return VertexesIterable(&nodesIterableLambdaHelper{iterFunc:iterator})
}

// Getting node predecessors
func (g *MixedMap) GetPredecessors(node VertexId) VertexesIterable {
	iterator := func() <-chan VertexId {
		ch := make(chan VertexId)
		
		go func() {
			defer func() {
				if e := recover(); e!=nil {
					err := erx.NewSequent("Getting node predecessors.", e)
					err.AddV("node id", node)
					panic(err)
				}
			}()
		
			accessorsMap, ok := g.connections[node]
			if !ok {
				panic(erx.NewError("Node doesn't exists."))
			}
			
			for VertexId, connType := range accessorsMap {
				if connType==CT_DIRECTED_REVERSED {
					ch <- VertexId
				}
			}
	
			close(ch)
		}()
		
		return ch
	}
	
	return VertexesIterable(&nodesIterableLambdaHelper{iterFunc:iterator}) 
}

func (g *MixedMap) CheckArc(from, to VertexId) (isExist bool) {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Checking arc existance in graph.", e)
			err.AddV("tail", from)
			err.AddV("head", to)
			panic(err)
		}
	}()
	
	connectedVertexes, ok := g.connections[from]
	if !ok {
		panic(erx.NewError("Tail node doesn't exist."))
	}
	
	if _, ok = g.connections[to]; !ok {
		panic(erx.NewError("Head node doesn't exist."))
	}
	
	connType, ok := connectedVertexes[to]

	return ok && connType==CT_DIRECTED
}

func (g *MixedMap) ArcsIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for from, connectedVertexes := range g.connections {
			for to, connType := range connectedVertexes {
				if connType!=CT_UNDIRECTED {
					ch <- Connection{from, to}
				}
			}
		}
		close(ch)
	}()
	return ch
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphEdgesWriter

// Adding edge to graph.
func (g *MixedMap) AddEdge(from, to VertexId) {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Add edge to graph.", e)
			err.AddV("node 1", from)
			err.AddV("node 2", to)
			panic(err)
		}
	}()

	g.touchNode(from)
	g.touchNode(to)
	
	if direction, ok := g.connections[from][to]; ok {
		err := erx.NewError("Duplicate connection.")
		err.AddV("connection type", direction)
		panic(err)
	}
	
	g.connections[from][to] = CT_UNDIRECTED
	g.connections[to][from] = CT_UNDIRECTED
	g.edgesCnt++

	return
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphEdgesRemover

// Removing arrow  'from' and 'to' nodes
func (g *MixedMap) RemoveEdge(from, to VertexId) {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Removing edge from graph.", e)
			err.AddV("node 1", from)
			err.AddV("node 2", to)
			panic(err)
		}
	}()
	
	if _, ok := g.connections[from]; !ok {
		panic(erx.NewError("First node doesn't exists"))
	}
	
	if _, ok := g.connections[to]; !ok {
		panic(erx.NewError("Second node doesn't exists"))
	}
	
	g.connections[from][to] = CT_NONE, false
	g.connections[to][from] = CT_NONE, false
	g.edgesCnt--

	return
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphReader

func (g *MixedMap) EdgesCnt() int {
	return g.edgesCnt
}

// Getting node predecessors
func (g *MixedMap) GetNeighbours(node VertexId) VertexesIterable {
	iterator := func() <-chan VertexId {
		ch := make(chan VertexId)
		
		go func() {
			defer func() {
				if e:=recover(); e!=nil {
					err := erx.NewSequent("Get node neighbours.", e)
					err.AddV("node id", node)
					panic(err)
				}
			}()
			
			if connectedMap, ok := g.connections[node]; ok {
				for VertexId, connType := range connectedMap {
					if connType==CT_UNDIRECTED {
						ch <- VertexId
					}
				}
			} else {
				panic(erx.NewError("Node doesn't exists."))
			}

			close(ch)
		}()
		
		return ch
	}
	
	return VertexesIterable(&nodesIterableLambdaHelper{iterFunc:iterator}) 
}

func (g *MixedMap) CheckEdge(from, to VertexId) bool {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Check edge existance in graph.", e)
			err.AddV("node 1", from)
			err.AddV("node 2", to)
			panic(err)
		}
	}()

	connectedVertexes, ok := g.connections[from]
	if !ok {
		panic(erx.NewError("Fist node doesn't exist."))
	}
	
	if _, ok = g.connections[to]; !ok {
		panic(erx.NewError("Second node doesn't exist."))
	}
	
	direction, ok := connectedVertexes[to]
	
	return ok && direction==CT_UNDIRECTED
}

func (g *MixedMap) EdgesIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for from, connectedVertexes := range g.connections {
			for to, connType := range connectedVertexes {
				if from<to && connType==CT_UNDIRECTED {
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
// MixedGraphSpecificReader

func (g *MixedMap) CheckEdgeType(tail VertexId, head VertexId) MixedConnectionType {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Check edge type in mixed graph.", e)
			err.AddV("tail", tail)
			err.AddV("head", head)
			panic(err)
		}
	}()
	
	connectedVertexes, ok := g.connections[tail]
	if !ok {
		panic(erx.NewError("Fist node doesn't exist."))
	}
	
	if _, ok = g.connections[head]; !ok {
		panic(erx.NewError("Second node doesn't exist."))
	}
	
	direction, ok := connectedVertexes[head]
	if !ok {
		direction = CT_NONE
	}
	
	return direction
}

func (g *MixedMap) ConnectionsCnt() int {
	return g.arcsCnt + g.edgesCnt
}

func (g *MixedMap) TypedConnectionsIter() <-chan TypedConnection {
	ch := make(chan TypedConnection)
	go func() {
		for from, connectedVertexes := range g.connections {
			for to, connType := range connectedVertexes {
				switch connType {
					case CT_NONE:
					case CT_UNDIRECTED:
						if from<to {
							ch <- TypedConnection{Connection:Connection{Tail: from, Head:to}, Type:CT_UNDIRECTED}
						} 
					case CT_DIRECTED:
						ch <- TypedConnection{Connection:Connection{Tail: from, Head:to}, Type:CT_DIRECTED}
					case CT_DIRECTED_REVERSED:
					default:
						err := erx.NewError("Internal error: wrong connection type in mixed graph matrix")
						err.AddV("connection type", connType)
						err.AddV("tail node", from)
						err.AddV("head node", to)
						panic(err)
				}
			}
		}
		close(ch)
	}()
	return ch
}
