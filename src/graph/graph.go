package graph

type VertexId uint

type Nodes []VertexId

type Connection struct {
	Tail VertexId
	Head VertexId
}

type TypedConnection struct {
	Connection
	Type MixedConnectionType
}

type ConnectionsIterable interface {
	ConnectionsIter() <-chan Connection
}

type EdgesIterable interface {
	EdgesIter() <-chan Connection
}

type ArcsIterable interface {
	ArcsIter() <-chan Connection
}

type TypedConnectionsIterable interface {
	TypedConnectionsIter() <-chan TypedConnection
}

type NodesIterable interface {
	NodesIter() <-chan VertexId
}

type NodesChecker interface {
	// Check node existance in graph
	CheckNode(node VertexId) bool
}

type GraphNodesWriter interface {
	// Adding single node to graph
	AddNode(node VertexId)
}

type GraphNodesReader interface {
	NodesChecker
	// Getting nodes count in graph
	NodesCnt() int
}

type GraphNodesRemover interface {
	// Removing node from graph
	RemoveNode(node VertexId)
}

type DirectedGraphArcsWriter interface {
	// Adding directed arc to graph
	AddArc(from, to VertexId)
}

type DirectedGraphArcsRemover interface {
	// Removding directed arc
	RemoveArc(from, to VertexId)
}

type DirectedGraphArcsReader interface {
	ArcsIterable
	
	// Getting arcs count in graph
	ArcsCnt() int

	// Getting all graph sources.
	GetSources() NodesIterable
	
	// Getting all graph sinks.
	GetSinks() NodesIterable
	
	// Getting node accessors
	GetAccessors(node VertexId) NodesIterable
	
	// Getting node predecessors
	GetPredecessors(node VertexId) NodesIterable
	
	// Checking arrow existance between node1 and node2
	//
	// node1 and node2 must exist in graph or error will be returned
	CheckArc(node1, node2 VertexId) bool
}

type DirectedGraphReader interface {
	GraphNodesReader
	DirectedGraphArcsReader
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
	EdgesIterable

	// Arrows count in graph
	EdgesCnt() int

	// Checking edge existance between node1 and node2
	//
	// node1 and node2 must exist in graph or error will be returned
	CheckEdge(node1, node2 VertexId) bool

	// Getting all nodes, connected to given one
	GetNeighbours(node VertexId) NodesIterable
}

type UndirectedGraphEdgesWriter interface {
	// Adding new edge to graph
	AddEdge(node1, node2 VertexId)	
}

type UndirectedGraphEdgesRemover interface {
	// Removing edge, connecting node1 and node2
	RemoveEdge(node1, node2 VertexId)
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

type MixedGraphSpecificReader interface {
	CheckEdgeType(tail, head VertexId) MixedConnectionType
	TypedConnectionsIterable
}

type MixedGraphConnectionsReader interface {
	ConnectionsIterable
	UndirectedGraphEdgesReader
	DirectedGraphArcsReader
	MixedGraphSpecificReader
}

type MixedGraphReader interface {
	MixedGraphConnectionsReader

	NodesIterable
	GraphNodesReader
}

type MixedGraphWriter interface {
	GraphNodesWriter
	UndirectedGraphEdgesWriter
	DirectedGraphArcsWriter
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

	MixedGraphSpecificReader
}
