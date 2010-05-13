package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)

type MixedMatrix struct {
	nodes []MixedConnectionType
	size int
	nodeIds map[NodeId]int // internal node ids, used in nodes array
	edgesCnt int
	arcsCnt int
}

func NewMixedMatrix(size int) *MixedMatrix {
	if size<=0 {
		panic(erx.NewError("Trying to create mixed matrix graph with zero size"))
	}
	g := new(MixedMatrix)
	g.nodes = make([]MixedConnectionType, size*(size-1)/2)
	g.size = size
	g.nodeIds = make(map[NodeId]int)
	return g
}

///////////////////////////////////////////////////////////////////////////////
// GraphNodesWriter

// Adding single node to graph
func (gr *MixedMatrix) AddNode(node NodeId) {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Add node to graph.", e)
			err.AddV("node id", node)
			panic(err)
		}
	}()
	
	if _, ok := gr.nodeIds[node]; ok {
		panic(erx.NewError("Node already exists."))
	}
	
	if len(gr.nodeIds) == gr.size {
		panic(erx.NewError("Not enough space to add new node"))
	}
	
	gr.nodeIds[node] = len(gr.nodeIds)
}

///////////////////////////////////////////////////////////////////////////////
// GraphNodesReader

// Getting nodes count in graph
func (gr *MixedMatrix) NodesCnt() int {
	return len(gr.nodeIds)
}

///////////////////////////////////////////////////////////////////////////////
// GraphNodesRemover

// Removing node from graph
func (gr *MixedMatrix) RemoveNode(node NodeId) {
	panic("Function RemoveNode doesn't implement in MixedMatrix graph yet.")
}
	
///////////////////////////////////////////////////////////////////////////////
// ConnectionsIterable
func (gr *MixedMatrix) ConnectionsIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for from, _ := range gr.nodeIds {
			for to, _ := range gr.nodeIds {
				if from>=to {
					continue
				}
				
				conn := gr.getConnectionId(from, to, false)
				if gr.nodes[conn]!=CT_NONE {
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
func (gr *MixedMatrix) NodesIter() <-chan NodeId {
	ch := make(chan NodeId)
	go func() {
		for nodeId, _ := range gr.nodeIds {
			ch <- nodeId
		}
		close(ch)
	}()
	return ch
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphEdgesWriter

// Adding new edge to graph
func (gr *MixedMatrix) AddEdge(node1, node2 NodeId) {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Add edge to mixed graph.", e)
			err.AddV("node 1", node1)
			err.AddV("node 2", node2)
			panic(err)
		}
	}()

	conn := gr.getConnectionId(node1, node2, true)
	if gr.nodes[conn]!=CT_NONE {
		err := erx.NewError("Duplicate edge.")
		err.AddV("connection id", conn)
		err.AddV("type", gr.nodes[conn]) 
		panic(err)
	}
	
	gr.nodes[conn] = CT_UNDIRECTED
	gr.edgesCnt++
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphEdgesRemover

// Removing edge, connecting node1 and node2
func (gr *MixedMatrix) RemoveEdge(node1, node2 NodeId) {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Remove edge from mixed graph.", e)
			err.AddV("node 1", node1)
			err.AddV("node 2", node2)
			panic(err)
		}
	}()

	conn := gr.getConnectionId(node1, node2, true)
	if gr.nodes[conn]!=CT_UNDIRECTED {
		err := erx.NewError("Edge doesn't exists.")
		err.AddV("connection id", conn)
		err.AddV("type", gr.nodes[conn])
		panic(err)
	}
	
	gr.nodes[conn] = CT_NONE
	gr.edgesCnt--
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphReader

// Edges count in graph
func (gr *MixedMatrix) EdgesCnt() int {
	return gr.edgesCnt
}

// Checking edge existance between node1 and node2
//
// node1 and node2 must exist in graph or error will be returned
func (gr *MixedMatrix) CheckEdge(node1, node2 NodeId) bool {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Check edge in mixed graph.", e)
			err.AddV("node 1", node1)
			err.AddV("node 2", node2)
			panic(err)
		}
	}()

	if node1==node2 {
		return false
	}
	
	return gr.nodes[gr.getConnectionId(node1, node2, false)]==CT_UNDIRECTED
}

// Getting all nodes, connected to given one
func (gr *MixedMatrix) GetNeighbours(node NodeId) Nodes {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Get node neighbours in mixed graph.", e)
			err.AddV("node", node)
			panic(err)
		}
	}()
	
	if _, ok := gr.nodeIds[node]; !ok {
		panic(erx.NewError("Unknown node."))
	}
	
	result := make([]NodeId, gr.size)
	ind := 0
	{
		var connId int
		for aNode, _ := range gr.nodeIds {
			if aNode==node {
				continue
			}
			connId= gr.getConnectionId(node, aNode, false)
			
			if gr.nodes[connId]==CT_UNDIRECTED {
				result[ind] = aNode
				ind++
			}
		}
	}
	
	return result[0:ind]
}

///////////////////////////////////////////////////////////////////////////////
// DirectedGraphArcsWriter

// Adding directed arc to graph
func (gr *MixedMatrix) AddArc(tail, head NodeId) {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Add arc to mixed graph.", e)
			err.AddV("tail", tail)
			err.AddV("head", head)
			panic(err)
		}
	}()

	conn := gr.getConnectionId(tail, head, true)
	if gr.nodes[conn]!=CT_NONE {
		err := erx.NewError("Duplicate edge.")
		err.AddV("connection id", conn)
		err.AddV("type", gr.nodes[conn]) 
		panic(err)
	}
	
	if tail<head {
		gr.nodes[conn] = CT_DIRECTED
	} else {
		gr.nodes[conn] = CT_DIRECTED_REVERSED
	}
	
	gr.arcsCnt++
}

///////////////////////////////////////////////////////////////////////////////
// DirectedGraphArcsRemover

// Removding directed arc
func (gr *MixedMatrix) RemoveArc(tail, head NodeId) {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Remove arc from mixed graph.", e)
			err.AddV("tail", tail)
			err.AddV("head", head)
			panic(err)
		}
	}()

	conn := gr.getConnectionId(tail, head, true)
	expectedType := CT_NONE
	if tail<head {
		expectedType = CT_DIRECTED
	} else {
		expectedType = CT_DIRECTED_REVERSED
	}
	
	if gr.nodes[conn]!=expectedType {
		err := erx.NewError("Arc doesn't exists.")
		err.AddV("connection id", conn)
		err.AddV("type", gr.nodes[conn])
		panic(err)
	}
	
	gr.nodes[conn] = CT_NONE
	gr.arcsCnt--
}

///////////////////////////////////////////////////////////////////////////////
// DirectedGraphReader

// Getting arcs count in graph
func (gr *MixedMatrix) ArcsCnt() int {
	return gr.arcsCnt
}

// Getting all graph sources.
func (gr *MixedMatrix) GetSources() Nodes {
	defer func() {
		if e := recover(); e!=nil {
			panic(erx.NewSequent("Get sources in mixed graph.", e))
		}
	}()
	
	result := make([]NodeId, gr.size)
	ind := 0
	{
		for tailNode, _ := range gr.nodeIds {
			hasPredecessors := false
			for headNode, _ := range gr.nodeIds {
				if tailNode==headNode {
					continue
				}

				checkingType := CT_NONE
				if tailNode < headNode {
					checkingType = CT_DIRECTED_REVERSED
				} else {
					checkingType = CT_DIRECTED
				}
			
				connId := gr.getConnectionId(tailNode, headNode, false)

				if gr.nodes[connId]==checkingType {
					hasPredecessors = true
					break
				}
			}
			if !hasPredecessors {
				result[ind] = tailNode
				ind++
			}
		}
	}
	
	return result[0:ind]
}

// Getting all graph sinks.
func (gr *MixedMatrix) GetSinks() Nodes {
	defer func() {
		if e := recover(); e!=nil {
			panic(erx.NewSequent("Get sinks in mixed graph.", e))
		}
	}()
	
	result := make([]NodeId, gr.size)
	ind := 0
	{
		for tailNode, _ := range gr.nodeIds {
			hasAccessors := false
			for headNode, _ := range gr.nodeIds {
				if tailNode==headNode {
					continue
				}

				checkingType := CT_NONE
				if tailNode < headNode {
					checkingType = CT_DIRECTED
				} else {
					checkingType = CT_DIRECTED_REVERSED
				}
			
				connId := gr.getConnectionId(tailNode, headNode, false)

				if gr.nodes[connId]==checkingType {
					hasAccessors = true
					break
				}
			}
			if !hasAccessors {
				result[ind] = tailNode
				ind++
			}
		}
	}
	
	return result[0:ind]
}

// Getting node accessors
func (gr *MixedMatrix) GetAccessors(node NodeId) Nodes {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Get node accessors in mixed graph.", e)
			err.AddV("node", node)
			panic(err)
		}
	}()
	
	result := make([]NodeId, gr.size)
	ind := 0
	{
		for headNode, _ := range gr.nodeIds {
			if node==headNode {
				// skipping loops
				continue
			}

			checkingType := CT_NONE
			if node < headNode {
				checkingType = CT_DIRECTED
			} else {
				checkingType = CT_DIRECTED_REVERSED
			}

			connId := gr.getConnectionId(node, headNode, false)
			
			if gr.nodes[connId]==checkingType {
				result[ind] = headNode
				ind++
			}
		}
	}
	
	return result[0:ind]
}

// Getting node predecessors
func (gr *MixedMatrix) GetPredecessors(node NodeId) Nodes {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Get node accessors in mixed graph.", e)
			err.AddV("node", node)
			panic(err)
		}
	}()
	
	result := make([]NodeId, gr.size)
	ind := 0
	{
		for tailNode, _ := range gr.nodeIds {
			if node==tailNode {
				// skipping loops
				continue
			}
		
			checkingType := CT_NONE
			if node < tailNode {
				checkingType = CT_DIRECTED_REVERSED
			} else {
				checkingType = CT_DIRECTED
			}


			connId := gr.getConnectionId(node, tailNode, false)
			
			if gr.nodes[connId]==checkingType {
				result[ind] = tailNode
				ind++
			}
		}
	}
	
	return result[0:ind]
}

// Checking arrow existance between node1 and node2
//
// node1 and node2 must exist in graph or error will be returned
func (gr *MixedMatrix) CheckArc(tail, head NodeId) bool {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Check arc in mixed graph.", e)
			err.AddV("tail", tail)
			err.AddV("head", head)
			panic(err)
		}
	}()
	
	checkingType := CT_NONE
	if tail < head {
		checkingType = CT_DIRECTED
	} else {
		checkingType = CT_DIRECTED_REVERSED
	}
	
	return gr.nodes[gr.getConnectionId(tail, head, false)]==checkingType
}

///////////////////////////////////////////////////////////////////////////////
// MixedGraphSpecificReader

// Iterate over only undirected edges
func (gr *MixedMatrix) EdgesIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for from, _ := range gr.nodeIds {
			for to, _ := range gr.nodeIds {
				if from>=to {
					continue
				}
				
				if gr.nodes[gr.getConnectionId(from, to, false)]==CT_UNDIRECTED {
					ch <- Connection{from, to}
				}
			}
		}
		close(ch)
	}()
	return ch
}
	
// Iterate over only directed arcs
func (gr *MixedMatrix) ArcsIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for from, _ := range gr.nodeIds {
			for to, _ := range gr.nodeIds {
				if from>=to {
					continue
				}
				
				conn := gr.getConnectionId(from, to, false)
				if gr.nodes[conn]==CT_DIRECTED {
					ch <- Connection{from, to}
				}
				if gr.nodes[conn]==CT_DIRECTED_REVERSED {
					ch <- Connection{to, from}
				}
			}
		}
		close(ch)
	}()
	return ch
}

func (gr *MixedMatrix) CheckEdgeType(tail NodeId, head NodeId) MixedConnectionType {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Check edge type in mixed graph.", e)
			err.AddV("tail", tail)
			err.AddV("head", head)
			panic(err)
		}
	}()
	
	conn := gr.getConnectionId(tail, head, false)
	return gr.nodes[conn]
}

func (gr *MixedMatrix) TypedConnectionsIter() <-chan TypedConnection {
	ch := make(chan TypedConnection)
	go func() {
		for from, _ := range gr.nodeIds {
			for to, _ := range gr.nodeIds {
				if from>=to {
					continue
				}
				
				conn := gr.getConnectionId(from, to, false)
				switch gr.nodes[conn] {
					case CT_NONE:
					case CT_UNDIRECTED:
						ch <- TypedConnection{Connection:Connection{Tail: from, Head:to}, Type:CT_UNDIRECTED} 
					case CT_DIRECTED:
						ch <- TypedConnection{Connection:Connection{Tail: from, Head:to}, Type:CT_DIRECTED}
					case CT_DIRECTED_REVERSED:
						ch <- TypedConnection{Connection:Connection{Tail: to, Head:from}, Type:CT_DIRECTED}
					default:
						err := erx.NewError("Internal error: wrong connection type in mixed graph matrix")
						err.AddV("connection type", gr.nodes[conn])
						err.AddV("connection id", conn)
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

func (gr *MixedMatrix) getConnectionId(node1, node2 NodeId, create bool) int {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Calculating connection id.", e)
			err.AddV("node 1", node1)
			err.AddV("node 2", node2)
			panic(err)
		}
	}()
	
	var id1, id2 int
	node1Exist := false
	node2Exist := false
	id1, node1Exist = gr.nodeIds[node1]
	id2, node2Exist = gr.nodeIds[node2]
	
	// checking for errors
	{
		if node1==node2 {
			panic(erx.NewError("Equal nodes."))
		}
		if !create {
			if !node1Exist {
				panic(erx.NewError("First node doesn't exist in graph"))
			}
			if !node2Exist {
				panic(erx.NewError("Second node doesn't exist in graph"))
			}
		} else if !node1Exist || !node2Exist {
			if node1Exist && node2Exist {
				if gr.size - len(gr.nodeIds) < 2 {
					panic(erx.NewError("Not enough space to create two new nodes."))
				}
			} else {
				if gr.size - len(gr.nodeIds) < 1 {
					panic(erx.NewError("Not enough space to create new node."))
				}
			}
		}
	}
	
	if !node1Exist {
		id1 = int(len(gr.nodeIds))
		gr.nodeIds[node1] = id1
	}

	if !node2Exist {
		id2 = int(len(gr.nodeIds))
		gr.nodeIds[node2] = id2
	}
	
	// switching id1, id2 in order to id1 < id2
	if id1>id2 {
		id1, id2 = id2, id1
	}
	
	// id from upper triangle matrix, stored in vector
	connId := id1*(gr.size-1) + id2 - 1 - id1*(id1+1)/2
	return connId 
}
