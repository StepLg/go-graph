package graph

import (
	"fmt"
	"os"
	"testing"
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
)

func DirectedGraphEquals(actual interface{}, expected interface{}) (match bool, pos os.Error, neg os.Error, err os.Error) {
	match = false
	neg = Errorf("Didn't expect that directed graphs are equal.")
	if aGr, ok := actual.(DirectedGraph); ok {
		if eGr, ok1 := expected.(DirectedGraph); ok1 {
			match = true
			missed := ""
			for arrow := range eGr.ArrowsIter() {
				isExist, err := aGr.CheckArrow(arrow.From, arrow.To)
				if err!=nil || !isExist {
					match = false
					if missed != "" {
						missed += ", "
					}
					missed += string(arrow.From) + "->" + string(arrow.To) 
				}
			}
			
			phantom := ""
			for arrow := range aGr.ArrowsIter() {
				isExist, err := eGr.CheckArrow(arrow.From, arrow.To)
				if err!=nil || !isExist {
					match = false
					if phantom!="" {
						phantom += ", "
					}
					phantom += fmt.Sprintf("%v->%v", arrow.From, arrow.To) 
				}
			}
			
			errorText := "Actual graph"
			if missed!="" {
				errorText += " miss " + missed + " arrows"
			}
			if missed!="" && phantom!="" {
				errorText += " and"
			}
			if phantom!="" {
				errorText += " contain " + phantom + " phantom arrows"
			}
			errorText += "."
			pos = Errorf(errorText)
		} else {
			err = Errorf("Expected DirectedGraph in actual, but was '%v' of type '%T'", expected, expected) 
		}
	} else {
		err = Errorf("Expected DirectedGraph in actual, but was '%v' of type '%T'", actual, actual) 
	}
	
	return
}

func ArrowsIteratorSpec(c gospec.Context) {
	gr := NewDirectedMap()
	
	c.Specify("Copy empty graph", func() {
		gr1 := NewDirectedMap()
		c.Expect(CopyDirectedGraph(gr, gr1), IsNil)
		c.Expect(gr1, DirectedGraphEquals, gr)
	})
	
	c.Specify("Copy simple directed graph", func() {
		gr1 := NewDirectedMap()
		gr.AddArrow(1, 2)
		gr.AddArrow(2, 3)
		gr.AddArrow(1, 4)
		gr.AddArrow(5, 1)

		c.Expect(CopyDirectedGraph(gr, gr1), IsNil)
		c.Expect(gr1, DirectedGraphEquals, gr)
	})
}

func TestArrowsIteratorSpec(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec("ArrowsIteratorSpec", ArrowsIteratorSpec)
	gospec.MainGoTest(r, t)
}
