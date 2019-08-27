package kstar

const rEdgeKey = -1

type pathGraphNode interface {
	CrossEdgeChild() pathGraphNode
	HeapEdgeChildren() []pathGraphNode
	D() float64
	EdgeKeys() (int, int, int)
}

func getHeapLeftChild(h pathGraphHeap, i, shifted int) pathGraphNode {

	iChild := 2*i + 1 - shifted
	if iChild >= h.Len() {
		return nil
	}

	return h.pq[iChild].(pathGraphNode)
}

func getHeapRightChild(h pathGraphHeap, i, shifted int) pathGraphNode {

	iChild := 2*i + 2 - shifted
	if iChild >= h.Len() {
		return nil
	}

	return h.pq[iChild].(pathGraphNode)
}

type hinNode struct {
	u, v, i int
	d       float64
	vHin    *pathGraphHeap
	hts     *map[int]*pathGraphHeap
}

func (n hinNode) EdgeKeys() (u, v, i int) {
	return n.u, n.v, n.i
}

func (n hinNode) D() float64 {
	return n.d
}

func (n hinNode) CrossEdgeChild() pathGraphNode {
	ht := (*n.hts)[n.u]
	if ht == nil {
		return nil
	}
	htRoot := ht.Top()
	if htRoot == nil {
		return nil
	}
	return htRoot.(pathGraphNode)
}

func (n hinNode) HeapEdgeChildren() []pathGraphNode {

	leftChild := n.getLeftChild()
	rightChild := n.getRightChild()

	children := []pathGraphNode{}

	if leftChild != nil {
		children = append(children, leftChild)
	}

	if rightChild != nil {
		children = append(children, rightChild)
	}

	return children
}

func (n hinNode) getLeftChild() pathGraphNode {

	pos := n.vHin.nodes[n.u][n.v][n.i]
	if pos == 0 {
		if n.vHin.Len() > 1 {
			return (*n.vHin).pq[1].(pathGraphNode)
		}
		return nil
	}

	return getHeapLeftChild(*n.vHin, pos, 1)
}

func (n hinNode) getRightChild() pathGraphNode {

	pos := n.vHin.nodes[n.u][n.v][n.i]
	if pos == 0 {
		return nil
	}

	return getHeapRightChild(*n.vHin, pos, 1)
}

func (n hinNode) equals(n2 hinNode) bool {
	return n.u == n2.u && n.v == n2.v && n.i == n2.i && n.d == n2.d
}

type htNode struct {
	hinNode *hinNode
	ht      *pathGraphHeap
}

func (n htNode) EdgeKeys() (u, v, i int) {
	return n.hinNode.EdgeKeys()
}

func (n htNode) D() float64 {
	return n.hinNode.D()
}

func (n htNode) CrossEdgeChild() pathGraphNode {
	return n.hinNode.CrossEdgeChild()
}

func (n htNode) HeapEdgeChildren() []pathGraphNode {

	children := n.hinNode.HeapEdgeChildren()
	leftChild := n.getLeftChild()
	rightChild := n.getRightChild()

	if leftChild != nil {
		children = append(children, leftChild)
	}

	if rightChild != nil {
		children = append(children, rightChild)
	}

	return children
}

func (n htNode) getLeftChild() pathGraphNode {
	pos := n.ht.nodes[n.hinNode.u][n.hinNode.v][n.hinNode.i]
	return getHeapLeftChild(*n.ht, pos, 0)
}

func (n htNode) getRightChild() pathGraphNode {
	pos := n.ht.nodes[n.hinNode.u][n.hinNode.v][n.hinNode.i]
	return getHeapRightChild(*n.ht, pos, 0)
}

type rNode struct {
	tHt *pathGraphHeap
}

func (n rNode) CrossEdgeChild() pathGraphNode {
	top := n.tHt.Top()
	if top == nil {
		return nil
	}
	return n.tHt.Top().(pathGraphNode)
}

func (n rNode) HeapEdgeChildren() []pathGraphNode {
	return []pathGraphNode{}
}

func (n rNode) D() float64 {
	return 0
}

// Dummy functions. Either this or separate EdgeKeys into another interface.
func (n rNode) EdgeKeys() (int, int, int) {
	return rEdgeKey, rEdgeKey, rEdgeKey
}

func (n rNode) SetIndex(i int) {
	// empty
}
