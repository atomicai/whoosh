package models

import "sync"

type HevHeap struct {
	Elements []*NodeWithValue
	sync.Mutex
}

func (h *HevHeap) Size() int {
	h.Lock()
	defer h.Unlock()
	return len(h.Elements)
}

func (h *HevHeap) Push(Element *NodeWithValue) {
	h.Lock()
	defer h.Unlock()
	h.Elements = append(h.Elements, Element)
	i := len(h.Elements) - 1
	for ; h.Elements[i].HValue < h.Elements[parent(i)].HValue; i = parent(i) {
		h.swap(i, parent(i))
	}
}

func (h *HevHeap) Pop() *NodeWithValue {
	h.Lock()
	defer h.Unlock()
	mn := h.Elements[0]
	h.Elements[0] = h.Elements[len(h.Elements)-1]
	h.Elements = h.Elements[:len(h.Elements)-1]
	h.rearrange(0)
	return mn
}

func (h *HevHeap) rearrange(i int) {
	smallest := i
	left, right, size := leftChild(i), rightChild(i), len(h.Elements)
	if left < size && h.Elements[left].HValue < h.Elements[smallest].HValue {
		smallest = left
	}
	if right < size && h.Elements[right].HValue < h.Elements[smallest].HValue {
		smallest = right
	}
	if smallest != i {
		h.swap(i, smallest)
		h.rearrange(smallest)
	}
}

func (h *HevHeap) swap(i, j int) {
	h.Elements[i], h.Elements[j] = h.Elements[j], h.Elements[i]
}
