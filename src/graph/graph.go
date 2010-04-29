package graph

import (
	"erx"
	"runtime"
	"strings"
)

type NodeId uint

type Nodes []NodeId

type Connection struct {
	Tail NodeId
	Head NodeId
}

type ConnectionsIterable interface {
	ConnectionsIter() <-chan Connection
}

type NodesIterable interface {
	NodesIter() <-chan NodeId
}

type GraphNodesWriter interface {
	// Adding single node to graph
	AddNode(node NodeId) erx.Error
}

type GraphNodesReader interface {
	// Getting nodes count in graph
	NodesCnt() int
}

type GraphNodesRemover interface {
	// Removing node from graph
	RemoveNode(node NodeId) erx.Error
}

type DirectedGraphArcsWriter interface {
	// Adding directed arc to graph
	AddArc(from, to NodeId) erx.Error
}

type DirectedGraphArcsRemover interface {
	// Removding directed arc
	RemoveArc(from, to NodeId) erx.Error
}

type DirectedGraphReader interface {
	// Getting arcs count in graph
	ArcsCnt() int

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
	CheckArc(node1, node2 NodeId) (bool, erx.Error)
	
}

// Interface representing directed graph
type DirectedGraph interface {
	GraphNodesWriter
	GraphNodesReader
	GraphNodesRemover

	ConnectionsIterable
	NodesIterable
	DirectedGraphArcsWriter
	DirectedGraphArcsRemover
	DirectedGraphReader
}

type UndirectedGraphReader interface {
	// Arrows count in graph
	EdgesCnt() int

	// Checking edge existance between node1 and node2
	//
	// node1 and node2 must exist in graph or error will be returned
	CheckEdge(node1, node2 NodeId) (bool, erx.Error)

	// Getting all nodes, connected to given one
	GetNeighbours(node NodeId) (Nodes, erx.Error)
}

type UndirectedGraphEdgesWriter interface {
	// Adding new edge to graph
	AddEdge(node1, node2 NodeId) (erx.Error)	
}

type UndirectedGraphEdgesRemover interface {
	// Removing edge, connecting node1 and node2
	RemoveEdge(node1, node2 NodeId) erx.Error
}

// Interface representing undirected graph
type UndirectedGraph interface {
	GraphNodesWriter
	GraphNodesReader
	GraphNodesRemover

	ConnectionsIterable
	NodesIterable

	UndirectedGraphEdgesWriter
	UndirectedGraphEdgesRemover
	UndirectedGraphReader
}

type MixedGraph interface {
	GraphNodesWriter
	GraphNodesReader
	GraphNodesRemover
	
	ConnectionsIterable
	NodesIterable

	UndirectedGraphEdgesWriter
	UndirectedGraphEdgesRemover
	UndirectedGraphReader

	DirectedGraphArcsWriter
	DirectedGraphArcsRemover
	DirectedGraphReader
	
	// Iterate over only undirected edges
	EdgesIter() ConnectionsIterable
	
	// Iterate over only directed arcs
	ArcsIter() ConnectionsIterable
}

func init() {
	// adding to erx directory prefix to cut from file names
	_, file, _, _ := runtime.Caller(0)
	dirName := file[0:strings.LastIndex(file, "/")]
	prevDirName := dirName[0:strings.LastIndex(dirName, "/")+1]
	erx.AddPathCut(prevDirName)
}
