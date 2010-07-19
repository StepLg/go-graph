package graph

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	
	"github.com/StepLg/go-erx/src/erx"
)

type graphWriterGeneric interface {
	AddNode(vertex VertexId)
	AddConnection(tail, head VertexId)
}

type graphWriterGeneric_ugraph struct {
	gr UndirectedGraphWriter
}

func (writer *graphWriterGeneric_ugraph) AddNode(vertex VertexId) {
	writer.gr.AddNode(vertex)
}

func (writer *graphWriterGeneric_ugraph) AddConnection(tail, head VertexId) {
	writer.gr.AddEdge(tail, head)
}

type graphWriterGeneric_dgraph struct {
	gr DirectedGraphWriter
}

func (writer *graphWriterGeneric_dgraph) AddNode(vertex VertexId) {
	writer.gr.AddNode(vertex)
}

func (writer *graphWriterGeneric_dgraph) AddConnection(tail, head VertexId) {
	writer.gr.AddArc(tail, head)
}

func readGraphLine(gr graphWriterGeneric, line string, connectionDelimiter string) {
	line = strings.Trim(line, " \t\n")
	if commentPos := strings.Index(line, "#"); commentPos!=-1 {
		// truncate comments
		line = strings.Trim(line[0:commentPos], " \t\n")
	}
	
	if line=="" {
		// skip empty lines
		return
	}
	
	if isMatch, err := regexp.MatchString("^[0-9]+$", line); isMatch && err==nil {
		// only one number - it's vertex id
		vertexId, err := strconv.Atoi(line)
		if err!=nil {
			panic(err)
		}
		
		gr.AddNode(VertexId(vertexId))
	} else if err!=nil {
		panic(err)
	}

	{
		// removing spaces between delimiter and vertexes
		reg1 := regexp.MustCompile("[ \t]*" + connectionDelimiter + "[ \t]*")
		line = reg1.ReplaceAllString(line, connectionDelimiter)
		
		// replace spaces to delimiter between vertexes
		reg2 := regexp.MustCompile("[ \t]+")
		line = reg2.ReplaceAllString(line, " ")
		line = strings.Replace(line, " ", connectionDelimiter, -1)
	}

	var prevVertexId VertexId
	hasPrev := false
	for _, nodeAsStr := range strings.Split(line, connectionDelimiter, -1) {
		nodeAsStr = strings.Trim(nodeAsStr, " \t\n")
		nodeAsInt, err := strconv.Atoi(nodeAsStr)
		if err!=nil {
			errErx := erx.NewSequent("Can't parse node id.", err)
			errErx.AddV("chunk", nodeAsStr)
			panic(errErx)
		}
		
		VertexId := VertexId(nodeAsInt)
		if hasPrev {
			gr.AddConnection(prevVertexId, VertexId)
		} else {
			hasPrev = true
		}
		prevVertexId = VertexId
	}
	return
}

func ReadUgraphLine(gr UndirectedGraphWriter, line string) {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Parsing graph edges from line.", e)
			err.AddV("line", line)
			panic(err)
		}
	}()
	
	readGraphLine(&graphWriterGeneric_ugraph{gr:gr}, line, "-")
}

func ReadDgraphLine(gr DirectedGraphWriter, line string) {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Parsing graph arcs from line.", e)
			err.AddV("line", line)
			panic(err)
		}
	}()
	
	readGraphLine(&graphWriterGeneric_dgraph{gr:gr}, line, ">")
}

func ReadMgraphLine(gr MixedGraphWriter, line string) {
	defer func() {
		if e:=recover(); e!=nil {
			err := erx.NewSequent("Parsing graph arcs and edges from line.", e)
			err.AddV("line", line)
			panic(err)
		}
	}()
	
	line = strings.Trim(line, " \t\n")
	if commentPos := strings.Index(line, "#"); commentPos!=-1 {
		// truncate comments
		line = strings.Trim(line[0:commentPos], " \t\n")
	}
	
	if line=="" {
		// skip empty lines
		return
	}
	
	if isMatch, err := regexp.MatchString("^[0-9]+$", line); isMatch && err==nil {
		// only one number - it's vertex id
		vertexId, err := strconv.Atoi(line)
		if err!=nil {
			panic(err)
		}
		
		gr.AddNode(VertexId(vertexId))
	} else if err!=nil {
		panic(err)
	}
	
	var prevVertexId VertexId
	hasPrev := false
	for _, nodeAsStr := range strings.Split(line, "-", -1) {
		nodeAsStr = strings.Trim(nodeAsStr, " \t\n")
		
		if strings.Index(nodeAsStr, ">")!=-1 {
			for index, nodeAsStr1 := range strings.Split(nodeAsStr, ">", -1) {
				nodeAsStr1 = strings.Trim(nodeAsStr1, " \t\n")
				nodeAsInt, err := strconv.Atoi(nodeAsStr1)
				if err!=nil {
					errErx := erx.NewSequent("Can't parse node id.", err)
					errErx.AddV("chunk", nodeAsStr1)
					panic(errErx)
				}
				
				VertexId := VertexId(nodeAsInt)
				if hasPrev {
					if index!=0 {
						gr.AddArc(prevVertexId, VertexId)
					} else {
						gr.AddEdge(prevVertexId, VertexId)
					}
				} else {
					hasPrev = true
				}
				prevVertexId = VertexId				
			}
		} else {
			nodeAsInt, err := strconv.Atoi(nodeAsStr)
			if err!=nil {
				errErx := erx.NewSequent("Can't parse node id.", err)
				errErx.AddV("chunk", nodeAsStr)
				panic(errErx)
			}
			
			VertexId := VertexId(nodeAsInt)
			if hasPrev {
				gr.AddEdge(prevVertexId, VertexId)
			} else {
				hasPrev = true
			}
			prevVertexId = VertexId
		}
	}
	
}

func readGraphFile(f io.Reader, lineParser func(string)) {
	reader := bufio.NewReader(f)
	var err os.Error
	var line string
	line, err = reader.ReadString('\n');
	for err==nil || err==os.EOF {
		lineParser(line)
		if err==os.EOF {
			break
		}
		line, err = reader.ReadString('\n');
	}
	if err!=nil && err!=os.EOF {
		erxErr := erx.NewSequent("Error while reading file.", err)
		panic(erxErr)
	}
}

func ReadUgraphFile(f io.Reader, gr UndirectedGraphWriter) {
	readGraphFile(f, func(line string) {
		ReadUgraphLine(gr, line)
	})
}

func ReadDgraphFile(f io.Reader, gr DirectedGraphWriter) {
	readGraphFile(f, func(line string) {
		ReadDgraphLine(gr, line)
	})
}

func ReadMgraphFile(f io.Reader, gr MixedGraphWriter) {
	readGraphFile(f, func(line string) {
		ReadMgraphLine(gr, line)
	})	
}
