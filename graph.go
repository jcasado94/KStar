package kstar

// Graph defines the graph interface used by K*.
type Graph interface {

	// Connections is the actual graph representation.
	// It returns a map matrix with the positive costs of the edges between any i and j (0 <= i, j!), or empty slice if there's no such connection.
	// Edge costs must be strictly positive. Loops allowed. Keep complexity on O(1).
	Connections() map[int]map[int][]float64

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
