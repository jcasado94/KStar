package kstar

// Graph defines the graph interface used by K*.
type Graph interface {

	// Connections is the implicit representation of our graph.
	// Given a graph node represented by non-negative integer n, it returns costs of the edges from n to any other node.
	// Edge costs must be strictly positive. Loops allowed. Keep complexity on O(1).
	Connections(n int) map[int][]float64

	// S returns the departure node.
	S() int

	//T returns the arrival node.
	T() int

	// FValues returns the heuristic cost from node i to T(), for any i.
	FValues() map[int]float64
}

// Edge represents an Edge defined in Graph.Connections(), specifically the ith from u to v.
type Edge struct {
	U, V, I int
}

func (e1 Edge) equals(e2 Edge) bool {
	return e1.U == e2.U && e1.V == e2.V && e1.I == e2.I
}
