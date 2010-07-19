package graph

import (
	"fmt"
	"io"
	"strings"
	"strconv"
)

func (node VertexId) String() string {
	return strconv.Itoa(int(uint(node)))
}

func (conn Connection) String() string {
	return fmt.Sprintf("%v->%v", conn.Tail, conn.Head)
}

// Typed connection as a string.
//
// Format:
//  * 1--2 for undirected connection
//  * 1->2 for directed connection
//  * 1<-2 for reversed directed connection
//  * 1><2 if there is no connection between 1 and 2
//  * 1!!2 for unexpected error (if function can't recognize connection type) 
func (conn TypedConnection) String() string {
	switch conn.Type {
		case CT_UNDIRECTED:
			return fmt.Sprintf("%v--%v", conn.Tail, conn.Head)
		case CT_DIRECTED:
			return fmt.Sprintf("%v->%v", conn.Tail, conn.Head)
		case CT_DIRECTED_REVERSED:
			return fmt.Sprintf("%v<-%v", conn.Tail, conn.Head)
		case CT_NONE:
			return fmt.Sprintf("%v><%v", conn.Tail, conn.Head)
	}
	return fmt.Sprintf("%v!!%v", conn.Tail, conn.Head)
}

func styleMapToString(style map[string]string) string {
	chunks := make([]string, len(style))
	i := 0
	for k, v := range style {
		chunks[i] = fmt.Sprintf("%v=\"%v\"", k, v)
		i++
	}
	return "[" + strings.Join(chunks, ",") + "]"
}

type DotNodeStyleFunc func(node VertexId) map[string]string
type DotConnectionStyleFunc func(conn TypedConnection) map[string]string

// Basic dot style function for vertexes.
//
// Generates only one style property: vertex label, which match vertex id.
func SimpleNodeStyle(node VertexId) map[string]string {
	style := make(map[string]string)
	style["label"] = node.String()
	return style
}

// Basic dot style function for connections.
//
// Empty style, except undirected connection. For undirected connection 
// "dir=both" style is set.
func SimpleConnectionStyle(conn TypedConnection) map[string]string {
	res := make(map[string]string)
	if conn.Type==CT_UNDIRECTED {
		// for undirected connection by default set "both" direction
		res["dir"] = "both"
	}
	return res
}

// Plot graph vertexes to dot format.
//
// nodesIter -- iterable object over graph vertexes
// wr -- writer interface
// styleFunc -- style function for vertexes. For example, see SimpleNodeStyle()
func PlotVertexesToDot(nodesIter VertexesIterable, wr io.Writer, styleFunc DotNodeStyleFunc) {
	if styleFunc==nil {
		styleFunc = SimpleNodeStyle
	}
	for node := range nodesIter.VertexesIter() {
		wr.Write([]byte("n" + node.String() + styleMapToString(styleFunc(node)) + ";\n"))
	}
}

// Plot graph connections to dot format.
//
// connIter -- iterable object over graph connections
// separator -- string betweet two vertexes. "->" for digraph and "--" for graph in graphviz file format
// wr -- writer interface
// styleFunc -- style function for typed connections. For example, see SimpleConnectionStyle()
func PlotConnectionsToDot(connIter TypedConnectionsIterable, separator string, wr io.Writer, styleFunc DotConnectionStyleFunc) {
	if styleFunc==nil {
		styleFunc = SimpleConnectionStyle
	}
	for conn := range connIter.TypedConnectionsIter() {
		wr.Write([]byte(fmt.Sprintf("n%v" + separator + "n%v%v;\n", 
			conn.Tail.String(),
			conn.Head.String(),
			styleMapToString(styleFunc(conn)))))
	}
}

// Plot directed graph to dot format.
//
// gr -- directed graph interface
// wr -- writer interface
// nodeStyleFunc -- style function for vertexes. For example, see SimpleNodeStyle()
// connStyleFunc -- style function for typed connections. For example, see SimpleConnectionStyle()
func PlotDgraphToDot(gr DirectedGraphReader, wr io.Writer, nodeStyleFunc DotNodeStyleFunc, connStyleFunc DotConnectionStyleFunc) {
	wr.Write([]byte("digraph messages {\n"))
	PlotVertexesToDot(gr, wr, nodeStyleFunc)
	PlotConnectionsToDot(ArcsToTypedConnIterable(gr), "->", wr, connStyleFunc)
	wr.Write([]byte("}\n"))
}

// Plot mixed graph to dot format.
//
// gr -- directed graph interface
// wr -- writer interface
// nodeStyleFunc -- style function for vertexes. For example, see SimpleNodeStyle()
// connStyleFunc -- style function for typed connections. For example, see SimpleConnectionStyle()
func PlotMgraphToDot(gr MixedGraphReader, wr io.Writer, nodeStyleFunc DotNodeStyleFunc, connStyleFunc DotConnectionStyleFunc) {
	wr.Write([]byte("digraph messages {\n"))
	PlotVertexesToDot(gr, wr, nodeStyleFunc)
	PlotConnectionsToDot(gr, "->", wr, connStyleFunc)
	wr.Write([]byte("}\n"))
}

// Plot undirected graph to dot format.
//
// gr -- directed graph interface
// wr -- writer interface
// nodeStyleFunc -- style function for vertexes. For example, see SimpleNodeStyle()
// connStyleFunc -- style function for typed connections. For example, see SimpleConnectionStyle()
func PlotUgraphToDot(gr UndirectedGraphReader, wr io.Writer, nodeStyleFunc DotNodeStyleFunc, connStyleFunc DotConnectionStyleFunc) {
	wr.Write([]byte("graph messages {\n"))
	PlotVertexesToDot(gr, wr, nodeStyleFunc)
	PlotConnectionsToDot(EdgesToTypedConnIterable(gr), "--", wr, connStyleFunc)
	wr.Write([]byte("}\n"))
}
