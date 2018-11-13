package kstar

import "container/heap"

type pathGraph struct {
	hin map[int]*pathGraphHeap
	ht  map[int]*pathGraphHeap
	r   rNode
}

func newPathGraph() *pathGraph {
	var pg pathGraph

	pg.hin = make(map[int]*pathGraphHeap)
	pg.ht = make(map[int]*pathGraphHeap)
	pg.r = rNode{}

	return &pg
}

func (pg *pathGraph) updateHinNodes(edges []*Edge, as *astar) {
	for _, e := range edges {
		if _, ok := pg.hin[e.V]; !ok {
			pg.hin[e.V] = newPathGraphHeap()
		}
		hin := pg.hin[e.V]
		n := hinNode{
			u:    e.U,
			v:    e.V,
			i:    e.I,
			d:    as.dValue(e),
			vHin: hin,
			hts:  &pg.ht,
		}

		currentTop := hin.Top()
		exists, pos := hin.exists(n.EdgeKeys())
		if exists {
			hin.replace(hin.pq[pos], n)
			heap.Fix(hin, pos)
		} else {
			heap.Push(hin, n)
		}
		top := hin.Top().(hinNode)

		if currentTop != nil && !currentTop.(hinNode).equals(top) {
			currentTopHin := currentTop.(hinNode)
			pg.propagateHinTopChange(&currentTopHin, &top, as)
		}

	}
}

func (pg *pathGraph) propagateHinTopChange(oldNode, newNode *hinNode, as *astar) {
	current := as.g.T()
	currentHt := pg.ht[current]
	for ok, pos := currentHt.exists(oldNode.u, oldNode.v, oldNode.i); ok; {
		currentHt.replace(oldNode, newNode)
		heap.Fix(currentHt, pos)
		current = as.searchTreeParents[current].U
		currentHt = pg.ht[current]
	}
}

func (pg *pathGraph) generateHts(as *astar) {
	pg.generateHt(as.g.S(), as.g.S(), as)
}

func (pg *pathGraph) generateHt(n, s int, as *astar) {

	if n == s {
		pg.ht[n] = newPathGraphHeap()
	} else {
		parent := as.searchTreeParents[n]
		htParent := pg.ht[parent.U]
		pg.ht[n] = copyPathGraphHeap(htParent)
	}

	ht := pg.ht[n]
	hin := pg.hin[n]
	if hin != nil {
		hinRoot := hin.Top().(hinNode)
		heap.Push(ht, htNode{
			hinNode: &hinRoot,
			ht:      ht,
		})
	}

	for _, child := range as.searchTreeChildren[n] {
		pg.generateHt(child, s, as)
	}
}
