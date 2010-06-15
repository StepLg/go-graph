package graph

type VertexId uint

type Vertexes []VertexId

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

type VertexesIterable interface {
	VertexesIter() <-chan VertexId
}

type VertexesChecker interface {
	// Check node existance in graph
	CheckNode(node VertexId) bool
}

type GraphVertexesWriter interface {
	// Adding single node to graph
	AddNode(node VertexId)
}

type GraphVertexesReader interface {
	VertexesChecker
	// Getting nodes count in graph
	Order() int
}

type GraphVertexesRemover interface {
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
	GetSources() VertexesIterable
	
	// Getting all graph sinks.
	GetSinks() VertexesIterable
	
	// Getting node accessors
	GetAccessors(node VertexId) VertexesIterable
	
	// Getting node predecessors
	GetPredecessors(node VertexId) VertexesIterable
	
	// Checking arrow existance between node1 and node2
	//
	// node1 and node2 must exist in graph or error will be returned
	CheckArc(node1, node2 VertexId) bool
}

type DirectedGraphReader interface {
	GraphVertexesReader
	DirectedGraphArcsReader
	VertexesIterable
}

// Interface representing directed graph
type DirectedGraph interface {
	GraphVertexesWriter
	GraphVertexesReader
	GraphVertexesRemover

	ConnectionsIterable
	VertexesIterable
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
	GetNeighbours(node VertexId) VertexesIterable
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
	GraphVertexesReader
	UndirectedGraphEdgesReader
	ConnectionsIterable
	VertexesIterable
}

// Interface representing undirected graph
type UndirectedGraph interface {
	GraphVertexesWriter
	GraphVertexesReader
	GraphVertexesRemover

	ConnectionsIterable
	VertexesIterable

	UndirectedGraphEdgesWriter
	UndirectedGraphEdgesRemover
	UndirectedGraphEdgesReader
}

type MixedGraphSpecificReader interface {
	CheckEdgeType(tail, head VertexId) MixedConnectionType
	ConnectionsCnt() int
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

	VertexesIterable
	GraphVertexesReader
}

type MixedGraphWriter interface {
	GraphVertexesWriter
	UndirectedGraphEdgesWriter
	DirectedGraphArcsWriter
}

type MixedGraph interface {
	GraphVertexesWriter
	GraphVertexesReader
	GraphVertexesRemover
	
	ConnectionsIterable
	VertexesIterable

	UndirectedGraphEdgesWriter
	UndirectedGraphEdgesRemover
	UndirectedGraphEdgesReader

	DirectedGraphArcsWriter
	DirectedGraphArcsRemover
	DirectedGraphArcsReader	

	MixedGraphSpecificReader
}
