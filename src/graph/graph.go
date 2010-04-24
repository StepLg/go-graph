package graph

import (
	"erx"
	"runtime"
	"strings"
)

type NodeId uint

type Nodes []NodeId

// Interface representing directed graph
type DirectedGraph interface {
	// Adding single node to graph
	AddNode(node NodeId) erx.Error

	// Adding arrow to graph.
	AddArrow(from, to NodeId) erx.Error
	
	// Removing arrow between 'from' and 'to' nodes
	RemoveArrow(from, to NodeId) erx.Error
	
	// Getting all graph sources.
	GetSources() (Nodes, erx.Error)
	
	// Getting all graph sinks.
	GetSinks() (Nodes, erx.Error)
	
	// Getting node accessors
	GetAccessors(node NodeId) (Nodes, erx.Error)
	
	// Getting node predecessors
	GetPredecessors(node NodeId) (Nodes, erx.Error)
	
	// Checking arrow existance between node1 and node2
	//
	// node1 and node2 must exist in graph or error will be returned
	CheckArrow(node1, node2 NodeId) (bool, erx.Error)

	ArrowsIter() <-chan Arrow
}

type Arrow struct {
	From NodeId
	To   NodeId
}

type DirectedArrowsIterable interface {
	ArrowsIter() <-chan Arrow
}

// Interface representing undirected graph
type UndirectedGraph interface {
	// Adding single node to graph
	AddNode(node NodeId) erx.Error

	// Nodes count in graph
	NodesCnt() int

	// Arrows count in graph
	ArrowsCnt() int

	// Adding new edge to graph
	AddEdge(node1, node2 NodeId) (erx.Error)
	
	// Removing edge, connecting node1 and node2
	RemoveEdge(node1, node2 NodeId) erx.Error
	
	// Checking edge existance between node1 and node2
	//
	// node1 and node2 must exist in graph or error will be returned
	CheckEdge(node1, node2 NodeId) (bool, erx.Error)
	
	// Getting all nodes, connected to given one
	GetNeighbours(node NodeId) (Nodes, erx.Error)
	
	EdgesIter() <-chan Arrow
}

func init() {
	// adding to erx directory prefix to cut from file names
	_, file, _, _ := runtime.Caller(0)
	dirName := file[0:strings.LastIndex(file, "/")]
	prevDirName := dirName[0:strings.LastIndex(dirName, "/")+1]
	erx.AddPathCut(prevDirName)
}
