package graph

import (
	"fmt"
	"os"
	"strings"
	"strconv"
)

func (node NodeId) String() string {
	return strconv.Itoa(int(uint(node)))
}

func (conn Connection) String() string {
	return fmt.Sprintf("%v->%v", conn.Tail, conn.Head)
}

type IWriter interface {
	Write(s string)
}

type StringWriter struct {
	Str string
}

func (wr *StringWriter) Write(s string) {
	wr.Str += s
}

type IOsStringWriter interface {
	WriteString(s string) (ret int, err os.Error)
}

type IoWriterAdapter struct {
	writer IOsStringWriter
}

func NewIoWriter(writer IOsStringWriter) *IoWriterAdapter {
	return &IoWriterAdapter{writer}
}

func (wr *IoWriterAdapter) Write(s string) {
	slen := len(s)
	wrote := 0
	for wrote<slen {
		curWrote, err := wr.writer.WriteString(s[wrote:])
		if err!=nil {
			panic(err)
		}
		wrote += curWrote
	}
}

func styleMapToString(style map[string]string) string {
	chunks := make([]string, len(style))
	i := 0
	for k, v := range style {
		chunks[i] = fmt.Sprintf("%v=%v", k, v)
		i++
	}
	return "[" + strings.Join(chunks, ",") + "]"
}

type DotNodeStyleFunc func(node NodeId) map[string]string
type DotConnectionStyleFunc func(conn Connection) map[string]string

func SimpleNodeStyle(node NodeId) map[string]string {
	style := make(map[string]string)
	style["label"] = node.String()
	return style
}

func SimpleArcStyle(conn Connection) map[string]string {
	return make(map[string]string)
}

func PlotNodesToDot(nodesIter NodesIterable, wr IWriter, styleFunc DotNodeStyleFunc) {
	for node := range nodesIter.NodesIter() {
		wr.Write("n" + node.String() + styleMapToString(styleFunc(node)) + ";\n")
	}
}

func PlotArcsToDot(connIter ConnectionsIterable, wr IWriter, styleFunc DotConnectionStyleFunc) {
	for conn := range connIter.ConnectionsIter() {
		wr.Write(fmt.Sprintf("n%v->n%v%v;\n", 
			conn.Tail.String(),
			conn.Head.String(),
			styleMapToString(styleFunc(conn))))
	}
}

func PlotDirectedGraphToDot(gr DirectedGraphReader, wr IWriter, nodeStyleFunc DotNodeStyleFunc, arcStyleFunc DotConnectionStyleFunc) {
	wr.Write("digraph messages {\n")
	PlotNodesToDot(gr, wr, nodeStyleFunc)
	PlotArcsToDot(gr, wr, arcStyleFunc)
	wr.Write("}\n")
}

func PlotMixedGraphToDot(gr MixedGraph, wr IWriter, nodeStyleFunc DotNodeStyleFunc, connStyleFunc DotConnectionStyleFunc) {
	wr.Write("digraph messages {\n")
	PlotNodesToDot(gr, wr, nodeStyleFunc)
	PlotArcsToDot(gr, wr, connStyleFunc)
	wr.Write("}\n")
}
