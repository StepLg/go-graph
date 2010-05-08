package graph

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
