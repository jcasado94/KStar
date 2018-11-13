package kstar

import (
	"encoding/json"
	"fmt"
	"testing"

	testutils "github.com/jcasado94/kstar/testUtils"
)

const (
	datasetPath   = "datasets/graph/"
	astarTestPath = "test/astar/"
)

// TestOutputAstar represents an astar test output file structure
type TestOutputAstar struct {
	TestName string
	InEdges  []*Edge
	MinPath  float64
}

func (to TestOutputAstar) Name() string {
	return to.TestName
}

func (to TestOutputAstar) Marshall() ([]byte, error) {
	return json.MarshalIndent(to, "", "    ")
}

func (to *TestOutputAstar) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, to)
}

func (to *TestOutputAstar) New(args ...interface{}) {
	to.TestName = args[0].(string)
	to.InEdges = args[1].([]*Edge)
	to.MinPath = args[2].(float64)
}

func (to *TestOutputAstar) TestFolderPath() string {
	return astarTestPath
}

func TestAllAstarInstances(t *testing.T) {
	tgs := testutils.GenerateTests(datasetPath)
	for _, tg := range tgs {
		as := newAstar(tg)
		inEdges, _ := as.run()
		minPathCost := as.minPathCost()
		to := new(TestOutputAstar)
		found := testutils.ReadTestOutput(to, tg.TestName, tg.TestName, inEdges, minPathCost)
		if found {
			// test
			expectedMinPath, expectedInEdges := to.MinPath, to.InEdges
			if as.minPathCost() > expectedMinPath {
				t.Errorf("Test %s failed! Min path cost higher. Expected %f, but found %f.", tg.TestName, expectedMinPath, minPathCost)
			} else if len(inEdges) != len(expectedInEdges) {
				t.Errorf("Test %s failed! Lengths of incoming edges differ. Expected %d, but found %d.", tg.TestName, len(expectedInEdges), len(inEdges))
			}
			m := mapEdges(inEdges)
			for _, edge := range expectedInEdges {
				if _, ok := m[edge.U][edge.V][edge.I]; !ok {
					t.Errorf("Test %s failed! Different incoming edges.", tg.TestName)
					break
				}
			}
		}
	}
}

func mapEdges(edges []*Edge) (m map[int]map[int]map[int]bool) {
	m = make(map[int]map[int]map[int]bool)
	for _, edge := range edges {
		u, v, i := edge.U, edge.V, edge.I
		if _, present := m[u]; !present {
			m[u] = make(map[int]map[int]bool)
		}
		if _, present := m[u][v]; !present {
			m[u][v] = make(map[int]bool)
		}
		m[u][v][i] = true
	}
	return
}

type mockGraph struct {
	connections map[int]map[int][]float64
	s, t        int
	fValues     map[int]float64
}

func (g mockGraph) Nodes() []int {
	nodes := make([]int, len(g.fValues))
	i := 0
	for k := range g.fValues {
		nodes[i] = k
		i++
	}
	return nodes
}

func (g mockGraph) Connections() map[int]map[int][]float64 {
	return g.connections
}

func (g mockGraph) S() int {
	return g.s
}

func (g mockGraph) T() int {
	return g.t
}

func (g mockGraph) FValues() map[int]float64 {
	return g.fValues
}

func newMockGraph(s, t int) mockGraph {
	return mockGraph{
		connections: make(map[int]map[int][]float64, 0),
		s:           s,
		t:           t,
		fValues:     make(map[int]float64, 0),
	}
}

func TestNewAstar(t *testing.T) {
	g := newMockGraph(0, 1)
	g.connections[0] = make(map[int][]float64, 0)
	g.connections[0][1] = []float64{1}

	as := newAstar(g)
	for node := range g.connections {
		if node != g.S() && as.open[node] != -1 {
			t.Errorf("%d is open after initialization.", node)
		} else if as.gScore[node] != 0 {
			t.Errorf("%d has non-zero gScore after initialization.", node)
		} else if as.searchTreeParents[node] != nil {
			t.Errorf("%d has a parent after initialization.", node)
		} else if as.c.arrivingEdges[node] != 0 {
			t.Errorf("%d has arriving edges marked after initialization.", node)
		}
	}

	if as.Len() != 1 || as.Top() != g.s {
		t.Error("Astar heap not well initialized.")
	} else if !as.c.start {
		t.Error("Astar expansionChecker set as not start after initialization.")
	}
}

func TestProcessEdges(t *testing.T) {
	current, neighbor := 0, 1
	edges := []float64{0.2, 0.1, 0.3}
	newEdges := make([]*Edge, 0)
	g := newMockGraph(0, 1)
	as := newAstar(g)

	minEdge, minCost := as.processEdges(current, neighbor, edges, &newEdges, false)
	expectedNewEdges := []Edge{Edge{0, 1, 0}, Edge{0, 1, 2}}
	if minEdge != 1 || minCost != 0.1 {
		t.Errorf("minimum edge not correctly processed. Expected %d (cost %f), but was %d (cost %f)", 1, 0.1, minEdge, minCost)
	} else {
		msg := fmt.Sprintf("newEdges not correctly processed. Should be %v, but was %v.", expectedNewEdges, newEdges)
		if len(newEdges) != len(expectedNewEdges) {
			t.Errorf(msg)
		} else {
			for i, edge := range expectedNewEdges {
				if newEdges[i].I != edge.I || newEdges[i].U != edge.U || newEdges[i].V != edge.V {
					t.Errorf(msg)
				}
			}
		}
	}
}

func TestExpansionConditionChecker(t *testing.T) {
	c := expansionConditionChecker{
		start: true,
	}
	if !c.shouldStop(1, 1) {
		t.Error("Should stop when finding t.")
	}

	c.start = false
	c.innerEdges, c.oldInnerEdges = 12, 4
	c.expandedNodes, c.oldExpandedNodes = 8, 3
	if !c.shouldStop(1, 1) {
		t.Error("Should stop if expanded nodes and inner edges are doubled.")
	}

	c.innerEdges = 7
	if c.shouldStop(1, 1) {
		t.Error("Should not stop if inner edges is not doubled.")
	}

	c.innerEdges = 12
	c.expandedNodes = 5
	if c.shouldStop(1, 1) {
		t.Error("Should not stop if expanded nodes is not doubles.")
	}

}
