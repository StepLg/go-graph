package graph

import (
	"sort"

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
	size int
}

// Create new simple nodes priority queue
//
// size is maximum number of nodes, which queue can store simultaneously
func newPriorityQueueSimple(size int) *nodesPriorityQueueSimple {
	if size<=0 {
		err := erx.NewError("Can't create priority queue with non-positive size.")
		err.AddV("size", size)
		panic(err)
	}
	
	q := &nodesPriorityQueueSimple {
		data: make(nodesPriority, size),
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
	for id, data := range q.data[0:q.size] {
		if data.Node == node {
			found = true
			if priority > q.data[id].Priority { 
				q.data[id].Priority = priority
			}
			break
		}
	}
	
	if !found {
		if q.size==len(q.data) {
			err := erx.NewError("Not enough space to add new node.")
			err.AddV("size", len(q.data))
			panic(err)
		}
		q.data[q.size].Node = node
		q.data[q.size].Priority = priority
		q.size++
	}
	
	sort.Sort(q.data[0:q.size])
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
