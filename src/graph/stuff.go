package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)

type MixedConnectionType uint8

const (
	CT_NONE MixedConnectionType = iota
	CT_UNDIRECTED
	CT_DIRECTED
	CT_DIRECTED_REVERSED
)

func (t MixedConnectionType) String() string {
	switch t {
		case CT_NONE : return "none"
		case CT_UNDIRECTED : return "undirected"
		case CT_DIRECTED : return "directed"
		case CT_DIRECTED_REVERSED : return "reversed"
	}
	
	return "unknown"
}

func NewUndirectedConnection(n1, n2 NodeId) TypedConnection {
	if n1>n2 {
		n1, n2 = n2, n1
	}
	return TypedConnection {
		Connection: Connection {
			Tail: n1,
			Head: n2,
		},
		Type: CT_UNDIRECTED,
	}
}

func NewDirectedConnection(tail, head NodeId) TypedConnection {
	return TypedConnection {
		Connection: Connection {
			Tail: tail,
			Head: head,
		},
		Type: CT_DIRECTED,
	}
}

// internal struct to store node with it's priority for priority queue
type priority_data_t struct {
	Node NodeId
	Priority float
}

type nodesPriority []priority_data_t

func (d nodesPriority) Less(i, j int) bool {
	return d[i].Priority < d[j].Priority 
}

func (d nodesPriority) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d nodesPriority) Len() int {
	return len(d)
}

// Nodes priority queue
type nodesPriorityQueue interface {
	// Add new item to queue
	Add(node NodeId, priority float)
	// Get item with max priority and remove it from the queue
	//
	// Panic if queue is empty
	Next() (NodeId, float)
	// Get item with max priority without removing it from the queue
	//
	// Panic if queue is empty
	Pick() (NodeId, float)
	// Total queue size
	Size() int
	// Check if queue is empty
	Empty() bool
}

// Very simple nodes priority queue
//
// Warning! It's very UNEFFICIENT!!!
type nodesPriorityQueueSimple struct {
	data nodesPriority
	nodesIndex map[NodeId]int
	size int
}

// Create new simple nodes priority queue
//
// size is maximum number of nodes, which queue can store simultaneously
func newPriorityQueueSimple(initialSize int) *nodesPriorityQueueSimple {
	if initialSize<=0 {
		err := erx.NewError("Can't create priority queue with non-positive size.")
		err.AddV("size", initialSize)
		panic(err)
	}
	
	q := &nodesPriorityQueueSimple {
		data: make(nodesPriority, initialSize),
		nodesIndex: make(map[NodeId]int),
		size: 0,
	}
	return q
}

// Add new item to queue
func (q *nodesPriorityQueueSimple) Add(node NodeId, priority float) {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("", e)
			err.AddV("node", node)
			err.AddV("priority", priority)
			panic(err)
		}
	}()
	
	found := false
	if id, ok := q.nodesIndex[node]; ok {
		if priority > q.data[id].Priority { 
			q.data[id].Priority = priority
			// changing position
			newId := id+1
			for q.data[newId].Priority<priority && newId<q.size {
				newId++
			}
			
			if newId > id+1 {
				// need to move
				copy(q.data[id:newId-1], q.data[id+1:newId])
				q.data[newId-1].Node = node
				q.data[newId-1].Priority = priority
			}
		}
		found = true
	}

	if !found {
		if q.size==len(q.data) {
			// resize
			// 2 is just a magic number
			newData := make(nodesPriority, 2*len(q.data))
			copy(newData, q.data)
			q.data = newData
		}
		id := 0
		for q.data[id].Priority<priority && id<q.size {
			id++
		}
		if id<q.size {
			copy(q.data[id+1:q.size+1], q.data[id:q.size])
		}
		q.data[id].Node = node
		q.data[id].Priority = priority
		q.nodesIndex[node] = id
		q.size++
	}
}

// Get item with max priority and remove it from the queue
//
// Panic if queue is empty
func (q *nodesPriorityQueueSimple) Next() (NodeId, float) {
	if q.Empty() {
		panic("Can't pick from empty queue.")
	}
	node := q.data[q.size-1].Node
	prior := q.data[q.size-1].Priority
	q.size--
	
	return node, prior
}

// Get item with max priority without removing it from the queue
//
// Panic if queue is empty
func (q *nodesPriorityQueueSimple) Pick() (NodeId, float) {
	if q.Empty() {
		panic("Can't pick from empty queue.")
	}
	node := q.data[q.size-1].Node
	prior := q.data[q.size-1].Priority
	return node, prior
}

// Total queue size
func (q *nodesPriorityQueueSimple) Size() int {
	return q.size
}

// Check if queue is empty
func (q *nodesPriorityQueueSimple) Empty() bool {
	return q.Size()==0
}

func matrixConnectionsIndexer(node1, node2 NodeId, nodeIds map[NodeId]int, size int, create bool) int {
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
	id1, node1Exist = nodeIds[node1]
	id2, node2Exist = nodeIds[node2]
	
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
				if size - len(nodeIds) < 2 {
					panic(erx.NewError("Not enough space to create two new nodes."))
				}
			} else {
				if size - len(nodeIds) < 1 {
					panic(erx.NewError("Not enough space to create new node."))
				}
			}
		}
	}
	
	if !node1Exist {
		id1 = int(len(nodeIds))
		nodeIds[node1] = id1
	}

	if !node2Exist {
		id2 = int(len(nodeIds))
		nodeIds[node2] = id2
	}
	
	// switching id1, id2 in order to id1 < id2
	if id1>id2 {
		id1, id2 = id2, id1
	}
	
	// id from upper triangle matrix, stored in vector
	connId := id1*(size-1) + id2 - 1 - id1*(id1+1)/2
	return connId 
}

type arcsToConnIterable_helper struct {
	gr DirectedGraphArcsReader
}

func (helper *arcsToConnIterable_helper) ConnectionsIter() <-chan Connection {
	return helper.gr.ArcsIter()
}

func ArcsToConnIterable(gr DirectedGraphArcsReader) ConnectionsIterable {
	return &arcsToConnIterable_helper{gr}
}

type edgesToConnIterable_helper struct {
	gr UndirectedGraphEdgesReader
}

func (helper *edgesToConnIterable_helper) ConnectionsIter() <-chan Connection {
	return helper.gr.EdgesIter()
}

func EdgesToConnIterable(gr UndirectedGraphEdgesReader) ConnectionsIterable {
	return &edgesToConnIterable_helper{gr}
}

type arcsToTypedConnIterable_helper struct {
	gr DirectedGraphArcsReader
}

func (helper *arcsToTypedConnIterable_helper) TypedConnectionsIter() <-chan TypedConnection {
	ch := make(chan TypedConnection)
	go func() {
		for conn := range helper.gr.ArcsIter() {
			ch <- TypedConnection{Connection: conn, Type: CT_DIRECTED}
		}
	}()
	return ch
}

func ArcsToTypedConnIterable(gr DirectedGraphArcsReader) TypedConnectionsIterable {
	return &arcsToTypedConnIterable_helper{gr}
}

type edgesToTypedConnIterable_helper struct {
	gr UndirectedGraphEdgesReader
}

func (helper *edgesToTypedConnIterable_helper) TypedConnectionsIter() <-chan TypedConnection {
	ch := make(chan TypedConnection)
	go func() {
		for conn := range helper.gr.EdgesIter() {
			ch <- TypedConnection{Connection: conn, Type: CT_UNDIRECTED}
		}
	}()
	return ch
}

func EdgesToTypedConnIterable(gr UndirectedGraphEdgesReader) TypedConnectionsIterable {
	return &edgesToTypedConnIterable_helper{gr}
}

// Helper struct to create nodes iterators with lambda functions
type nodesIterableLambdaHelper struct {
	iterFunc func() <-chan NodeId
}

func (helper *nodesIterableLambdaHelper) NodesIter() <-chan NodeId {
	return helper.iterFunc()
}
