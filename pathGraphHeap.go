package kstar

import "container/heap"

const undefinedPos = -1

type pathGraphHeap struct {
	pq    []pathGraphNode
	nodes set
}

func newPathGraphHeap() *pathGraphHeap {
	pgh := &pathGraphHeap{
		pq:    make([]pathGraphNode, 0),
		nodes: set{},
	}
	heap.Init(pgh)
	return pgh
}

func copyHt(src *pathGraphHeap, dst *pathGraphHeap) {
	pq := make([]pathGraphNode, 0)
	for _, htn := range src.pq {
		pq = append(pq, htNode{
			hinNode: htn.(htNode).hinNode,
			ht:      dst,
		})
	}
	nodes := src.nodes.copy()
	dst.pq, dst.nodes = pq, nodes
}

func (h pathGraphHeap) Len() int { return len(h.pq) }

func (h pathGraphHeap) Empty() bool { return h.Len() == 0 }

func (h pathGraphHeap) Less(i, j int) bool {
	return h.pq[i].D() < h.pq[j].D()
}

func (h pathGraphHeap) Swap(i, j int) {
	h.pq[i], h.pq[j] = h.pq[j], h.pq[i]

	u1, v1, i1 := h.pq[i].EdgeKeys()
	u2, v2, i2 := h.pq[j].EdgeKeys()
	h.nodes.swap(u1, v1, i1, u2, v2, i2)
}

func (h *pathGraphHeap) Push(x interface{}) {
	n := x.(pathGraphNode)
	pos := h.Len()
	h.pq = append(h.pq, n)

	u, v, i := n.EdgeKeys()
	h.nodes.put(pos, u, v, i)
}

func (h *pathGraphHeap) Pop() interface{} {
	old := h.pq
	n := len(old)
	rel := old[n-1]
	h.pq = old[0 : n-1]

	h.nodes.remove(rel.EdgeKeys())

	return rel
}

func (h pathGraphHeap) Top() interface{} {
	if len(h.pq) == 0 {
		return nil
	}
	return h.pq[0]
}

func (h pathGraphHeap) exists(u, v, i int) (exists bool, pos int) {
	if h.nodes.exists(u, v, i) {
		return true, h.nodes[u][v][i]
	}
	return false, -1
}

func (h *pathGraphHeap) replace(oldNode, newNode pathGraphNode) {
	uOld, vOld, iOld := oldNode.EdgeKeys()
	uNew, vNew, iNew := newNode.EdgeKeys()
	pos := h.nodes[uOld][vOld][iOld]
	h.nodes.put(pos, uNew, vNew, iNew)
	h.pq[pos] = newNode
}

// TODO: change name
type set map[int]map[int]map[int]int

func (s *set) copy() (dst set) {
	dst = make(map[int]map[int]map[int]int)
	for k1, v1 := range *s {
		dst[k1] = make(map[int]map[int]int)
		for k2, v2 := range v1 {
			dst[k1][k2] = make(map[int]int)
			for k3, v3 := range v2 {
				dst[k1][k2][k3] = v3
			}
		}
	}
	return dst
}

func (s set) put(pos, u, v, i int) {
	if _, ok := s[u]; !ok {
		s[u] = make(map[int]map[int]int)
	}
	if _, ok := s[u][v]; !ok {
		s[u][v] = make(map[int]int)
	}
	s[u][v][i] = pos
}

func (s set) remove(u, v, i int) {
	s[u][v][i] = undefinedPos
}

func (s set) exists(u, v, i int) bool {
	if _, ok := s[u]; ok {
		if _, ok = s[u][v]; ok {
			pos := s[u][v][i]
			return pos != undefinedPos
		}
	}
	return false
}

func (s set) swap(u1, v1, i1, u2, v2, i2 int) {
	s[u1][v1][i1], s[u2][v2][i2] = s[u2][v2][i2], s[u1][v1][i1]
}
