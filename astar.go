package kstar

import (
	"container/heap"
)

// astar keeps the needed structure for the K* astar algorithm. The algorithm does not assume a monotonic heuristic function is provided in g.
type astar struct {
	g Graph

	pq                 []int
	open               map[int]int // pq position, -1 if closed
	gScore             map[int]float64
	searchTreeParents  map[int]*Edge
	searchTreeChildren map[int]map[int]interface{}

	c expansionConditionChecker
}

// newAstar generates a new Astar instance given a Graph implementation
func newAstar(g Graph) *astar {

	var as astar

	as.g = g
	as.open = make(map[int]int, 0)
	as.gScore = make(map[int]float64, 0)
	as.searchTreeParents = make(map[int]*Edge, 0)
	as.searchTreeChildren = make(map[int]map[int]interface{}, 0)
	arrivingEdges := make(map[int]int, 0)

	initNode(g.S(), &as, arrivingEdges)

	heap.Init(&as)
	heap.Push(&as, g.S())

	as.c = expansionConditionChecker{start: true, arrivingEdges: arrivingEdges}

	return &as
}

func initNode(n int, as *astar, arrivingEdges map[int]int) {
	as.open[n] = -1
	as.gScore[n] = 0
	as.searchTreeParents[n] = nil
	as.searchTreeChildren[n] = make(map[int]interface{}, 0)
	arrivingEdges[n] = 0
}

func (as *astar) run() (newEdges []*Edge, empty bool) {

	newEdges = make([]*Edge, 0)

	for !as.Empty() {

		if as.c.shouldStop(as.Top().(int), as.g.T()) {
			return newEdges, false
		}

		current := heap.Pop(as).(int)
		reopening := as.c.expand(current)

		for neighbor, edges := range as.g.Connections(current) {

			if _, ok := as.open[neighbor]; !ok {
				initNode(neighbor, as, as.c.arrivingEdges)
			}

			if len(edges) == 0 {
				continue
			}

			as.c.hit(neighbor, len(edges))
			minEdge, minCost := as.processEdges(current, neighbor, edges, &newEdges, reopening)

			tentativeScore := as.gScore[current] + minCost
			isOpen := as.open[neighbor] != -1
			hasParent := as.searchTreeParents[neighbor] != nil

			e := Edge{current, neighbor, minEdge}
			if hasParent {
				if tentativeScore >= as.gScore[neighbor] {
					newEdges = appendIf(newEdges, &e, !reopening)
					continue
				}
				newEdges = appendIf(newEdges, as.searchTreeParents[neighbor], !reopening)
			}

			if neighbor == as.g.S() {
				newEdges = appendIf(newEdges, &e, !reopening)
			} else {
				if hasParent {
					oldParent := as.searchTreeParents[neighbor].U
					delete(as.searchTreeChildren[oldParent], neighbor)
				}
				as.searchTreeParents[neighbor] = &e
				as.searchTreeChildren[current][neighbor] = true
			}

			as.gScore[neighbor] = tentativeScore
			if isOpen {
				heap.Fix(as, as.open[neighbor])
			} else {
				heap.Push(as, neighbor)
			}

		}

		as.open[current] = -1
		as.c.close(current)

	}

	return newEdges, true

}

func appendIf(newEdges []*Edge, e *Edge, should bool) []*Edge {
	if should {
		return append(newEdges, e)
	}
	return newEdges
}

func (as astar) processEdges(current, neighbor int, edges []float64, newEdges *[]*Edge, reopening bool) (minEdge int, minCost float64) {

	minCost, minEdge = edges[0], 0
	e := 1
	for _, cost := range edges[1:] {
		if cost < minCost {
			*newEdges = appendIf(*newEdges, &Edge{current, neighbor, minEdge}, !reopening)
			minEdge, minCost = e, cost
		} else {
			*newEdges = appendIf(*newEdges, &Edge{current, neighbor, e}, !reopening)
		}
		e++
	}

	return
}

func (as astar) fScore(n int) float64 {
	return as.gScore[n] + as.g.FValue(n)
}

func (as astar) minPathCost() (cost float64) {

	node := as.g.T()
	for node != as.g.S() {
		e := as.searchTreeParents[node]
		if e == nil {
			break
		}
		cost += as.g.Connections(e.U)[node][e.I]
		node = e.U
	}

	if node != as.g.S() {
		return -1
	}

	return cost

}

func (as astar) dValue(e *Edge) float64 {
	cost := as.g.Connections(e.U)[e.V][e.I]
	return as.gScore[e.U] + cost - as.gScore[e.V]
}

func (as astar) Len() int { return len(as.pq) }

func (as astar) Empty() bool { return len(as.pq) == 0 }

func (as astar) Less(i, j int) bool {
	return as.fScore(as.pq[i]) < as.fScore(as.pq[j])
}

func (as astar) Swap(i, j int) {
	as.pq[i], as.pq[j] = as.pq[j], as.pq[i]
	as.open[as.pq[i]], as.open[as.pq[j]] = i, j
}

func (as *astar) Push(x interface{}) {
	as.open[x.(int)] = len(as.pq)
	as.pq = append(as.pq, x.(int))
}

func (as *astar) Pop() interface{} {
	old := as.pq
	l := len(old)
	n := old[l-1]
	as.pq = old[0 : l-1]
	return n
}

func (as astar) Top() interface{} {
	return as.pq[0]
}

type expansionConditionChecker struct {
	arrivingEdges                   map[int]int // nr of arriving edges per open node, -1 if node is closed
	innerEdges, expandedNodes       int
	oldInnerEdges, oldExpandedNodes int
	start                           bool
}

func (c *expansionConditionChecker) shouldStop(top, t int) bool {

	stop := false

	if c.start {
		if top == t {
			c.start = false
			stop = true
		}
	} else {
		innerEdgesDoubled := c.oldInnerEdges == 0 && float32(c.innerEdges+2)/float32(c.oldInnerEdges+1) >= 2 || float32(c.innerEdges)/float32(c.oldInnerEdges) >= 2
		expandedNodesDoubled := c.oldExpandedNodes == 0 && float32(c.expandedNodes+2)/float32(c.oldExpandedNodes+1) >= 2 || float32(c.expandedNodes)/float32(c.oldExpandedNodes) >= 2
		stop = innerEdgesDoubled && expandedNodesDoubled
	}

	if stop {
		c.oldInnerEdges = c.innerEdges
		c.oldExpandedNodes = c.expandedNodes
	}

	return stop

}

func (c *expansionConditionChecker) hit(n, hits int) {
	// mind if heuristic not monotonic and we are opening a closed node, we might be counting edges twice.
	if c.arrivingEdges[n] > -1 {
		c.arrivingEdges[n] += hits
	} else {
		c.innerEdges += hits
	}
}

func (c *expansionConditionChecker) close(n int) {
	c.innerEdges += c.arrivingEdges[n]
	c.arrivingEdges[n] = -1
}

func (c *expansionConditionChecker) expand(n int) bool {
	c.expandedNodes++
	if c.arrivingEdges[n] == -1 {
		c.arrivingEdges[n] = 0
		return true
	}
	return false
}

func (c *expansionConditionChecker) opened(n int) bool {
	return c.arrivingEdges[n] != -1
}
