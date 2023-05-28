package models

import "sync"

type ValueHeap struct {
	Elements []*NodeWithValue
	sync.Mutex
}

func (h *ValueHeap) Size() int {
	h.Lock()
	defer h.Unlock()
	return len(h.Elements)
}

func (h *ValueHeap) Push(Element *NodeWithValue) {
	h.Lock()
	defer h.Unlock()
	h.Elements = append(h.Elements, Element)
	i := len(h.Elements) - 1
	for ; h.Elements[i].Value < h.Elements[parent(i)].Value; i = parent(i) {
		h.swap(i, parent(i))
	}
}

func (h *ValueHeap) Pop() *NodeWithValue {
	h.Lock()
	defer h.Unlock()
	mn := h.Elements[0]
	h.Elements[0] = h.Elements[len(h.Elements)-1]
	h.Elements = h.Elements[:len(h.Elements)-1]
	h.rearrange(0)
	return mn
}

func (h *ValueHeap) rearrange(i int) {
	smallest := i
	left, right, size := leftChild(i), rightChild(i), len(h.Elements)
	if left < size && h.Elements[left].Value < h.Elements[smallest].Value {
		smallest = left
	}
	if right < size && h.Elements[right].Value < h.Elements[smallest].Value {
		smallest = right
	}
	if smallest != i {
		h.swap(i, smallest)
		h.rearrange(smallest)
	}
}

func (h *ValueHeap) swap(i, j int) {
	h.Elements[i], h.Elements[j] = h.Elements[j], h.Elements[i]
}

func parent(i int) int {
	return (i - 1) / 2
}

func leftChild(i int) int {
	return 2*i + 1
}

func rightChild(i int) int {
	return 2*i + 2
}
