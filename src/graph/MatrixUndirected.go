package graph

import "erx"

type UndirectedGraphMatrix struct {
	nodes []bool
	size int
	nodeIds map[NodeId]int // internal node ids, used in nodes array
}

// Creating new undirected graph with matrix storage.
//
// size means maximum number of nodes, used in graph. Trying to add
// more nodes, than this size will cause an error. 
func NewUndirectedGraphMatrix(size int) *UndirectedGraphMatrix {
	if size<=0 {
		return nil
	}
	g := new(UndirectedGraphMatrix)
	g.nodes = make([]bool, size*(size-1)/2)
	g.size = size
	g.nodeIds = make(map[NodeId]int)
	return g
}

// Maximum graph capacity
//
// Maximum nodes count graph can handle
func (g *UndirectedGraphMatrix) GetCapacity() int {
	return int(g.size)
}

// Current nodes count in graph
func (g *UndirectedGraphMatrix) GetSize() int {
	return int(len(g.nodeIds))
}

// Adding new edge to graph
func (g *UndirectedGraphMatrix) AddEdge(node1, node2 NodeId) (err erx.Error) {
	var conn int
	conn, err = g.getConnectionId(node1, node2, true)
	if nil!=err {
		goto Error
	}
	
	if g.nodes[conn] {
		err = erx.NewError("Duplicate edge.")
		err.AddV("edge id", conn)
		goto Error
	}
	g.nodes[conn] = true
	
	return
	
	Error:
	err = erx.NewSequent("Can't add new edge to graph.", err)
	err.AddV("from node", node1)
	err.AddV("to node", node2)
	return
}

// Removing edge, connecting node1 and node2
func (g *UndirectedGraphMatrix) RemoveEdgeBetween(node1, node2 NodeId) erx.Error {
	var conn int
	var err erx.Error
	conn, err = g.getConnectionId(node1, node2, false)
	if nil!=err {
		goto Error
	}
	if (!g.nodes[conn]) {
		err = erx.NewError("Edge doesn't exist.")
		goto Error
	}
	g.nodes[conn] = false
	return nil
	
	Error:
	err = erx.NewSequent("Can't remove edge from graph.", err)
	err.AddV("from node", node1)
	err.AddV("to node", node2)
	return err
}

// Removing edge with specific id
func (g *UndirectedGraphMatrix) RemoveEdge(id int) erx.Error {
	var err erx.Error
	if int(id) >= len(g.nodes) {
		err = erx.NewError("Edge id out of bounds.")
		err.AddV("total edges count", len(g.nodes))
		goto Error
	}
	if !g.nodes[id] {
		err = erx.NewError("Edge doesn't exist.")
		goto Error
	}
	g.nodes[id] = false
	
	Error:
	err = erx.NewSequent("Can't remove edge from graph.", err)
	err.AddV("edge id", id)
	return err
}

// Getting all nodes, connected to given one
func (g *UndirectedGraphMatrix) GetConnected(node NodeId) (Nodes, erx.Error) {
	var err erx.Error
	if _, ok := g.nodeIds[node]; !ok {
		err = erx.NewError("Unknown node.")
		goto Error
	}
	
	result := make([]NodeId, g.size)
	ind := 0
	{
		var connId int
		for aNode, _ := range g.nodeIds {
			if aNode==node {
				continue
			}
			connId, err = g.getConnectionId(node, aNode, false)
			if nil!=err {
				goto Error
			}
			
			if g.nodes[connId] {
				result[ind] = aNode
				ind++
			}
		}
	}
	
	return result[0:ind], nil
	
	Error:
	err = erx.NewSequent("Can't find connected nodes.", err)
	err.AddV("node", node)
	return nil, err
}

func (g *UndirectedGraphMatrix) CheckEdgeBetween(node1, node2 NodeId) (bool, erx.Error) {
	var conn int
	var err erx.Error
	conn, err = g.getConnectionId(node1, node2, false)
	if nil!=err {
		goto Error
	}
	return g.nodes[conn], nil
		
	Error:
	err = erx.NewSequent("Can't check edge in graph.", err)
	err.AddV("from node", node1)
	err.AddV("to node", node2)
	return false, err
	
}

func (g *UndirectedGraphMatrix) getConnectionId(node1, node2 NodeId, create bool) (int, erx.Error) {
	var id1, id2 int
	node1Exist := false
	node2Exist := false
	id1, node1Exist = g.nodeIds[node1]
	id2, node2Exist = g.nodeIds[node2]
	
	// checking for errors
	{
		err := erx.NewError("Can't get edge ID.")
		err.AddV("from node", node1)
		err.AddV("to node", node2)
		err.AddV("create flag", create)
		if node1==node2 {
			err.AddE(erx.NewError("Equal nodes."))
		} else if !create {
			if !node1Exist {
				err.AddE(erx.NewError("Node " + string(node1) + "doesn't exist in graph"))
			}
			if !node2Exist {
				err.AddE(erx.NewError("Node " + string(node2) + "doesn't exist in graph"))
			}
		} else if !node1Exist || !node2Exist {
			if node1Exist && node2Exist {
				if g.size - len(g.nodeIds) < 2 {
					err.AddE(erx.NewError("Not enough space to create two new nodes."))
				}
			} else {
				if g.size - len(g.nodeIds) < 1 {
					err.AddE(erx.NewError("Not enough space to create new node."))
				}
			}
		}
		if err.Errors().Front()!=nil {
			return 0, err
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
	connId := id1*g.size + id2 - 1 - id1*(id1+1)/2
	return connId , nil	
}
