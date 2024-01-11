package smmoth_weight_rb

import (
	"fmt"
	"sync"
)

type WeightRoundRobinWithRing2 struct {
	Nodes       []*Node
	mu          sync.RWMutex //when add node or del node ,lock
	totalWeight int
	x           []*Node
	i           int
}

func NewWeightRoundRobinWithRing2() *WeightRoundRobinWithRing2 {
	return &WeightRoundRobinWithRing2{
		Nodes: []*Node{},
	}

}
func (r *WeightRoundRobinWithRing2) AddNode(node *Node) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	//check
	for _, n := range r.Nodes {
		if n.Host == node.Host {
			return DuplicateNodeError
		}
	}
	r.Nodes = append(r.Nodes, node)
	r.totalWeight += node.Weight
	r.x = make([]*Node, r.totalWeight)

	for i := 0; i < r.totalWeight; i++ {
		n := r.next()
		r.x[i] = n
	}
	return nil
}

func (r *WeightRoundRobinWithRing2) DelNode(node *Node) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, n := range r.Nodes {
		if n == node {
			r.Nodes = append(r.Nodes[:i], r.Nodes[i+1:]...)
			r.totalWeight -= node.Weight
			break
		}
	}
	r.x = make([]*Node, r.totalWeight)
	for i := 0; i < r.totalWeight; i++ {
		n := r.next()
		r.x[i] = n
	}
}

func (r *WeightRoundRobinWithRing2) Print() {
	for _, node := range r.Nodes {
		fmt.Print(node.CurrentWeight, " ")
	}
}

func (r *WeightRoundRobinWithRing2) next() *Node {

	n := len(r.Nodes)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return r.Nodes[0]
	}
	flag := 0
	m := 0
	for i := 0; i < n; i++ {
		r.Nodes[i].CurrentWeight += r.Nodes[i].Weight
		if r.Nodes[i].CurrentWeight > m {
			m = r.Nodes[i].CurrentWeight
			flag = i
		}
	}
	r.Nodes[flag].CurrentWeight -= r.totalWeight
	return r.Nodes[flag]
}

func (r *WeightRoundRobinWithRing2) Next() *Node {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.i == r.totalWeight {
		r.i = 0
	}
	n := r.x[r.i]
	r.i++
	return n
}
