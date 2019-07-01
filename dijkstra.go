package kstar

import (
	"container/heap"
)

type dijkstra struct {
	pq        []*dijkstraNode
	formerTop *dijkstraNode
}

func newDijkstra(rn *rNode) (d *dijkstra) {
	r := newDijkstraNode(rn, 0, []*dijkstraNode{}, false, true)
	d = &dijkstra{
		pq:        []*dijkstraNode{r},
		formerTop: nil,
	}
	heap.Init(d)
	return
}

type dijkstraNode struct {
	n       pathGraphNode
	cost    float64
	path    []*dijkstraNode
	isCross bool
	isR     bool
}

func newDijkstraNode(n pathGraphNode, cost float64, path []*dijkstraNode, isCross, isR bool) *dijkstraNode {
	dn := dijkstraNode{
		n:       n,
		cost:    cost,
		path:    make([]*dijkstraNode, len(path), cap(path)),
		isCross: isCross,
		isR:     isR,
	}
	for i, elem := range path {
		dn.path[i] = elem
	}
	dn.path = append(dn.path, &dn)
	return &dn
}

func (d *dijkstra) step() (path []*dijkstraNode, empty bool) {

	path = make([]*dijkstraNode, 0)
	current := heap.Pop(d).(*dijkstraNode)

	d.pushChildren(current)

	if d.Empty() {
		d.formerTop = current
		return current.path, true
	}

	return current.path, false

}

func (d *dijkstra) pushChildren(current *dijkstraNode) (hasChildren bool) {
	c := current.n.CrossEdgeChild()
	if c != nil {
		heap.Push(d, newDijkstraNode(c, current.cost+c.D(), current.path, true, false))
		hasChildren = true
	}

	for _, c := range current.n.HeapEdgeChildren() {
		heap.Push(d, newDijkstraNode(c, current.cost+c.D()-current.n.D(), current.path, false, false))
		hasChildren = true
	}

	return hasChildren
}

func (d *dijkstra) resume() (end bool) {
	return !d.pushChildren(d.formerTop)
}

func (d dijkstra) Len() int { return len(d.pq) }

func (d dijkstra) Empty() bool { return len(d.pq) == 0 }

func (d dijkstra) Less(i, j int) bool {
	return d.pq[i].cost < d.pq[j].cost
}

func (d dijkstra) Swap(i, j int) {
	d.pq[i], d.pq[j] = d.pq[j], d.pq[i]
}

func (d *dijkstra) Push(x interface{}) {
	d.pq = append(d.pq, x.(*dijkstraNode))
}

func (d *dijkstra) Pop() interface{} {
	old := d.pq
	l := len(old)
	n := old[l-1]
	d.pq = old[0 : l-1]
	return n
}

func (d dijkstra) Top() interface{} {
	return d.pq[0]
}
