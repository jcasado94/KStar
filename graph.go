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

	// FValue returns the heuristic cost from node n to T().
	FValue(n int) float64
}

// Edge represents an Edge defined in Graph.Connections(), specifically the ith from u to v.
type Edge struct {
	U, V, I int
}

// RemoveLoopPaths removes the paths which contain loops.
func RemoveLoopPaths(paths [][]*Edge) (modifiedPaths [][]*Edge) {
	for _, path := range paths {
		admittable := true
		beenTo := map[int]bool{path[0].U: true}
		for _, edge := range path {
			if beenTo[edge.V] {
				admittable = false
				break
			}
			beenTo[edge.V] = true
		}
		if admittable {
			modifiedPaths = append(modifiedPaths, path)
		}
	}
	return modifiedPaths
}
