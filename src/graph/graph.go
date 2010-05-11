package graph

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
	AddNode(node NodeId)
}

type GraphNodesReader interface {
	// Getting nodes count in graph
	NodesCnt() int
}

type GraphNodesRemover interface {
	// Removing node from graph
	RemoveNode(node NodeId)
}

type DirectedGraphArcsWriter interface {
	// Adding directed arc to graph
	AddArc(from, to NodeId)
}

type DirectedGraphArcsRemover interface {
	// Removding directed arc
	RemoveArc(from, to NodeId)
}

type DirectedGraphArcsReader interface {
	// Getting arcs count in graph
	ArcsCnt() int

	// Getting all graph sources.
	GetSources() Nodes
	
	// Getting all graph sinks.
	GetSinks() Nodes
	
	// Getting node accessors
	GetAccessors(node NodeId) Nodes
	
	// Getting node predecessors
	GetPredecessors(node NodeId) Nodes
	
	// Checking arrow existance between node1 and node2
	//
	// node1 and node2 must exist in graph or error will be returned
	CheckArc(node1, node2 NodeId) bool
}

type DirectedGraphReader interface {
	GraphNodesReader
	DirectedGraphArcsReader
	ConnectionsIterable
	NodesIterable
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
	DirectedGraphArcsReader
}

type UndirectedGraphEdgesReader interface {
	// Arrows count in graph
	EdgesCnt() int

	// Checking edge existance between node1 and node2
	//
	// node1 and node2 must exist in graph or error will be returned
	CheckEdge(node1, node2 NodeId) bool

	// Getting all nodes, connected to given one
	GetNeighbours(node NodeId) Nodes
}

type UndirectedGraphEdgesWriter interface {
	// Adding new edge to graph
	AddEdge(node1, node2 NodeId)	
}

type UndirectedGraphEdgesRemover interface {
	// Removing edge, connecting node1 and node2
	RemoveEdge(node1, node2 NodeId)
}

type UndirectedGraphReader interface {
	GraphNodesReader
	UndirectedGraphEdgesReader
	ConnectionsIterable
	NodesIterable
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
	UndirectedGraphEdgesReader
}

type MixedGraph interface {
	GraphNodesWriter
	GraphNodesReader
	GraphNodesRemover
	
	ConnectionsIterable
	NodesIterable

	UndirectedGraphEdgesWriter
	UndirectedGraphEdgesRemover
	UndirectedGraphEdgesReader

	DirectedGraphArcsWriter
	DirectedGraphArcsRemover
	DirectedGraphArcsReader
	
	// Iterate over only undirected edges
	EdgesIter() <-chan Connection 
	
	// Iterate over only directed arcs
	ArcsIter() <-chan Connection
	
	CheckEdgeType(tail, head NodeId) MixedConnectionType
}
