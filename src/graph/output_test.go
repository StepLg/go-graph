package graph

/**
import (
	"testing"
	"fmt"

	"github.com/orfjackal/gospec/src/gospec"
	// . "github.com/orfjackal/gospec/src/gospec"
)

func PlotDirectedGraphToDotSpec(c gospec.Context) {
	gr := NewDirectedMap()
	gr.AddArc(1, 2)
	gr.AddArc(2, 3)
	gr.AddArc(3, 4)
	gr.AddArc(2, 4)
	gr.AddArc(4, 5)
	gr.AddArc(1, 6)
	gr.AddArc(2, 6)
	
	writer := &StringWriter{}
	PlotDirectedGraphToDot(gr, writer, SimpleNodeStyle, SimpleArcStyle)
	fmt.Println(writer.Str)
}

func TestOutput(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(PlotDirectedGraphToDotSpec)
	gospec.MainGoTest(r, t)
}
*/
