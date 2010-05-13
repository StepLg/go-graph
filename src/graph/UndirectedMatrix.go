package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)

type UndirectedMatrix struct {
	nodes []bool
	size int
	nodeIds map[NodeId]int // internal node ids, used in nodes array
	edgesCnt int
}

// Creating new undirected graph with matrix storage.
//
// size means maximum number of nodes, used in graph. Trying to add
// more nodes, than this size will cause an error. 
func NewUndirectedMatrix(size int) *UndirectedMatrix {
	if size<=0 {
		return nil
	}
	g := new(UndirectedMatrix)
	g.nodes = make([]bool, size*(size-1)/2)
	g.size = size
	g.nodeIds = make(map[NodeId]int)
	g.edgesCnt = 0
	return g
}

// Maximum graph capacity
//
// Maximum nodes count graph can handle
func (g *UndirectedMatrix) GetCapacity() int {
	return int(g.size)
}

///////////////////////////////////////////////////////////////////////////////
// ConnectionsIterable

func (g *UndirectedMatrix) ConnectionsIter() <-chan Connection {
	return g.EdgesIter()
}

///////////////////////////////////////////////////////////////////////////////
// NodesIterable

func (g *UndirectedMatrix) NodesIter() <-chan NodeId {
	ch := make(chan NodeId)
	go func() {
		for nodeId, _ := range g.nodeIds {
			ch <- nodeId
		}
		close(ch)
	}()
	return ch
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphNodesWriter

// Adding single node to graph
func (g *UndirectedMatrix) AddNode(node NodeId) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Add node to graph.", err, 1)
		res.AddV("node id", node)
		return
	}

	if _, ok := g.nodeIds[node]; ok {
		panic(makeError(erx.NewError("Node already exists.")))
	}
	
	g.nodeIds[node] = len(g.nodeIds)

	return	
}

///////////////////////////////////////////////////////////////////////////////
// GraphNodesRemover

func (g *UndirectedMatrix) RemoveNode(node NodeId) {
	panic(erx.NewError("Function doesn't implemented yet."))
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphEdgesWriter

// Adding new edge to graph
func (g *UndirectedMatrix) AddEdge(node1, node2 NodeId) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Add edge to graph.", err, 1)
		res.AddV("node 1", node1)
		res.AddV("node 2", node2)
		return
	}
	
	defer func() {
		// warning! such code generates wrong file/line info about error!
		// see http://groups.google.com/group/golang-nuts/browse_thread/thread/66bd57dcdac63aa
		// for details
		if err := recover(); err!=nil {
			panic(makeError(err))
		}
	}()

	var conn int
	conn = g.getConnectionId(node1, node2, true)
	
	if g.nodes[conn] {
		err := erx.NewError("Duplicate edge.")
		err.AddV("connection id", conn)
		panic(makeError(err))
	}
	g.nodes[conn] = true
	g.edgesCnt++
	
	return
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphEdgesRemover

// Removing edge, connecting node1 and node2
func (g *UndirectedMatrix) RemoveEdge(node1, node2 NodeId) {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Remove edge from graph.", err, 1)
		res.AddV("node 1", node1)
		res.AddV("node 2", node2)
		return
	}
	
	defer func() {
		// warning! such code generates wrong file/line info about error!
		// see http://groups.google.com/group/golang-nuts/browse_thread/thread/66bd57dcdac63aa
		// for details
		if err := recover(); err!=nil {
			panic(makeError(err))
		}
	}()

	
	conn := g.getConnectionId(node1, node2, false)
	
	if (!g.nodes[conn]) {
		panic(erx.NewError("Edge doesn't exist."))
	}
	
	g.nodes[conn] = false
	g.edgesCnt--
	
	return
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphReader

// Current nodes count in graph
func (g *UndirectedMatrix) NodesCnt() int {
	return int(len(g.nodeIds))
}

// Current nodes count in graph
func (g *UndirectedMatrix) EdgesCnt() int {
	return g.edgesCnt
}


// Getting all nodes, connected to given one
func (g *UndirectedMatrix) GetNeighbours(node NodeId) Nodes {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Get node neighbours.", err, 1)
		res.AddV("node", node)
		return
	}
	
	defer func() {
		// warning! such code generates wrong file/line info about error!
		// see http://groups.google.com/group/golang-nuts/browse_thread/thread/66bd57dcdac63aa
		// for details
		if err := recover(); err!=nil {
			panic(makeError(err))
		}
	}()
	
	if _, ok := g.nodeIds[node]; !ok {
		panic(makeError(erx.NewError("Unknown node.")))
	}
	
	result := make([]NodeId, g.size)
	ind := 0
	{
		var connId int
		for aNode, _ := range g.nodeIds {
			if aNode==node {
				continue
			}
			connId= g.getConnectionId(node, aNode, false)
			
			if g.nodes[connId] {
				result[ind] = aNode
				ind++
			}
		}
	}
	
	return result[0:ind]
}

func (g *UndirectedMatrix) EdgesIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for from, _ := range g.nodeIds {
			for to, _ := range g.nodeIds {
				if from<to && g.CheckEdge(from, to) {
					ch <- Connection{from, to}
				}
			}
		}
		close(ch)
	}()
	return ch
}

func (g *UndirectedMatrix) CheckEdge(node1, node2 NodeId) bool {
	defer func() {
		// warning! such code generates wrong file/line info about error!
		// see http://groups.google.com/group/golang-nuts/browse_thread/thread/66bd57dcdac63aa
		// for details
		if err := recover(); err!=nil {
			errErx := erx.NewSequent("Checking edge", err)
			errErx.AddV("node 1", node1)
			errErx.AddV("node 2", node2)
			panic(errErx)
		}
	}()

	return g.nodes[g.getConnectionId(node1, node2, false)]
}

func (g *UndirectedMatrix) getConnectionId(node1, node2 NodeId, create bool) int {
	makeError := func(err interface{}) (res erx.Error) {
		res = erx.NewSequentLevel("Calculating connection id.", err, 1)
		res.AddV("node 1", node1)
		res.AddV("node 2", node2)
		return
	}
	
	defer func() {
		// warning! such code generates wrong file/line info about error!
		// see http://groups.google.com/group/golang-nuts/browse_thread/thread/66bd57dcdac63aa
		// for details
		if err := recover(); err!=nil {
			panic(makeError(err))
		}
	}()
	
	var id1, id2 int
	node1Exist := false
	node2Exist := false
	id1, node1Exist = g.nodeIds[node1]
	id2, node2Exist = g.nodeIds[node2]
	
	// checking for errors
	{
		if node1==node2 {
			panic(makeError(erx.NewError("Equal nodes.")))
		}
		if !create {
			if !node1Exist {
				panic(makeError(erx.NewError("First node doesn't exist in graph")))
			}
			if !node2Exist {
				panic(makeError(erx.NewError("Second node doesn't exist in graph")))
			}
		} else if !node1Exist || !node2Exist {
			if node1Exist && node2Exist {
				if g.size - len(g.nodeIds) < 2 {
					panic(makeError(erx.NewError("Not enough space to create two new nodes.")))
				}
			} else {
				if g.size - len(g.nodeIds) < 1 {
					panic(makeError(erx.NewError("Not enough space to create new node.")))
				}
			}
		}
	}
	
	if !node1Exist {
		id1 = int(len(g.nodeIds))
		g.nodeIds[node1] = id1
	}

	if !node2Exist {
		id2 = int(len(g.nodeIds))
		g.nodeIds[node2] = id2
	}
	
	// switching id1, id2 in order to id1 < id2
	if id1>id2 {
		id1, id2 = id2, id1
	}
	
	// id from upper triangle matrix, stored in vector
	connId := id1*(g.size-1) + id2 - 1 - id1*(id1+1)/2
	return connId 
}
