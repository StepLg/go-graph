package graph

import (
	"erx"
	"runtime"
	"strings"
)

type NodeId uint

type ConnectionId uint

type Nodes []NodeId

// Interface representing directed graph
type DirectedGraph interface {
	// Adding single node to graph
	AddNode(node NodeId) erx.Error

	// Adding arrow to graph.
	AddArrow(from, to NodeId) (ConnectionId, erx.Error)
	
	// Removing arrow between 'from' and 'to' nodes
	RemoveArrowBetween(from, to NodeId) erx.Error
	
	// Removing arrow with specific Id
	RemoveArrow(id ConnectionId) erx.Error
	
	// Getting all graph sources.
	GetSources() (Nodes, erx.Error)
	
	// Getting all graph sinks.
	GetSinks() (Nodes, erx.Error)
	
	// Getting node accessors
	GetAccessors(node NodeId) (Nodes, erx.Error)
	
	// Getting node predecessors
	GetPredecessors(node NodeId) (Nodes, erx.Error)
}

// Interface representing undirected graph
type UndirectedGraph interface {
	// Adding new edge to graph
	AddEdge(node1, node2 NodeId) (ConnectionId, erx.Error)
	
	// Removing edge, connecting node1 and node2
	RemoveEdgeBetween(node1, node2 NodeId) erx.Error
	
	// Checking edge existance between node1 and node2
	//
	// node1 and node2 must exist in graph or error will be returned
	CheckEdgeBetween(node1, node2 NodeId) (bool, erx.Error)
	
	// Removing edge with specific id
	RemoveEdge(id ConnectionId) erx.Error
	
	// Getting all nodes, connected to given one
	GetConnected(node NodeId) (Nodes, erx.Error)
}

/*
// Nodes iterator over the graph
type NodesIterator interface {
	Iter() <- NodeId
}
*/

func init() {
	// adding to erx directory prefix to cut from file names
	_, file, _, _ := runtime.Caller(0)
	dirName := file[0:strings.LastIndex(file, "/")]
	prevDirName := dirName[0:strings.LastIndex(dirName, "/")+1]
	erx.AddPathCut(prevDirName)
}
