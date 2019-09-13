package kstar

type kstar struct {
	pg          *pathGraph
	as          *astar
	d           *dijkstra
	paths       [][]Edge
	asExhausted bool
}

func newKstar(g *Graph, pg *pathGraph) kstar {
	return kstar{
		pg:    pg,
		as:    newAstar(*g),
		d:     newDijkstra(&pg.r),
		paths: make([][]Edge, 0),
	}
}

// Run returns the k shortest paths given a Graph implementation and k.
func Run(g Graph, k int) (paths [][]Edge) {
	pg := newPathGraph()
	ks := newKstar(&g, pg)
	paths = make([][]Edge, 0)
	tReached := ks.startAstar()
	if tReached {
		for len(paths) < k {
			sigmaPath, empty := ks.d.step()
			edgeSeq := buildSeq(sigmaPath)
			path := buildPath(edgeSeq, ks.as.searchTreeParents, g.S(), g.T())
			paths = append(paths, path)
			if empty {
				if ks.asExhausted {
					break
				}
				ks.resumeAstar()
				end := ks.d.resume()
				if end {
					break
				}
			}
		}
	}

	return paths
}

func (ks *kstar) startAstar() (tReached bool) {
	newEdges, end := ks.as.run()
	if !end {
		ks.pg.updateHinNodes(newEdges, ks.as)
		ks.pg.generateHts(ks.as)
		tHt := ks.pg.ht[ks.as.g.T()]
		ks.pg.r.tHt = tHt
	}
	return !end
}

func (ks *kstar) resumeAstar() {
	newEdges, end := ks.as.run()
	ks.asExhausted = end
	ks.pg.updateHinNodes(newEdges, ks.as)
}

// Transforms a dijkstra path into a sequence of sidetrack edges
func buildSeq(pgPath []*dijkstraNode) (seq []Edge) {
	seq = make([]Edge, 0)
	if len(pgPath) < 2 {
		// length 1 is just R
		return seq
	}

	vdn := pgPath[len(pgPath)-1]
	u, v, i := (*vdn).n.EdgeKeys()
	seq = append(seq, Edge{U: u, V: v, I: i})

	for j := len(pgPath) - 1; j > 0; j-- {
		udn := pgPath[j-1]
		if vdn.isCross && !udn.isR {
			u, v, i = (*udn).n.EdgeKeys()
			seq = append(seq, Edge{U: u, V: v, I: i})
		}
		vdn = udn
	}

	return seq

}

// Adds tree nodes to the sidetrack edges to complete the path
func buildPath(seq []Edge, spTree map[int]Edge, s, t int) (path []Edge) {
	path = make([]Edge, 0)
	current := t
	for current != s || len(seq) != 0 {
		if len(seq) > 0 && seq[len(seq)-1].V == current {
			e := seq[len(seq)-1]
			path = append(path, e)
			seq = seq[:len(seq)-1]
			current = e.U
		} else {
			path = append(path, spTree[current])
			current = spTree[current].U
		}
	}
	return path
}
