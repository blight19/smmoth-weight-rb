package smmoth_weight_rb

import (
	"container/ring"
	"errors"
	"fmt"
	"sync"
)

var DuplicateNodeError = errors.New("duplicate node")

type Node struct {
	Host          string
	Weight        int
	CurrentWeight int
	Meta          map[string]string //other info that you want to store
}

func NewNode(name string, weight int, meta map[string]string) *Node {
	return &Node{
		Host:   name,
		Weight: weight,
		Meta:   meta,
	}
}

type WeightRoundRobinWithRing struct {
	Nodes       []*Node
	mu          sync.RWMutex //when add node or del node ,lock
	totalWeight int
	rr          *ring.Ring //ring
}

func NewWeightRoundRobinWithRing() *WeightRoundRobinWithRing {
	return &WeightRoundRobinWithRing{
		Nodes: []*Node{},
	}

}
func (r *WeightRoundRobinWithRing) AddNode(node *Node) error {
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
	r.rr = ring.New(r.totalWeight)

	for i := 0; i < r.totalWeight; i++ {
		n := r.next()
		r.rr.Value = n
		r.rr = r.rr.Next()
	}
	return nil
}

func (r *WeightRoundRobinWithRing) DelNodeByName(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, n := range r.Nodes {
		if n.Host == name {
			r.Nodes = append(r.Nodes[:i], r.Nodes[i+1:]...)
			r.totalWeight -= n.Weight
			break
		}
	}
	r.rr = ring.New(r.totalWeight)
	for i := 0; i < r.totalWeight; i++ {
		n := r.next()
		r.rr.Value = n
		r.rr = r.rr.Next()
	}
}

func (r *WeightRoundRobinWithRing) DelNode(node *Node) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, n := range r.Nodes {
		if n == node {
			r.Nodes = append(r.Nodes[:i], r.Nodes[i+1:]...)
			r.totalWeight -= node.Weight
			break
		}
	}
	r.rr = ring.New(r.totalWeight)
	for i := 0; i < r.totalWeight; i++ {
		n := r.next()
		r.rr.Value = n
		r.rr = r.rr.Next()
	}
}

func (r *WeightRoundRobinWithRing) Print() {
	for _, node := range r.Nodes {
		fmt.Print(node.CurrentWeight, " ")
	}
}

func (r *WeightRoundRobinWithRing) next() *Node {

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

func (r *WeightRoundRobinWithRing) Next() *Node {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n := r.rr.Value.(*Node)
	r.rr = r.rr.Next()
	return n
}

type WeightRoundRobin struct {
	Nodes       []*Node
	mu          sync.Mutex
	totalWeight int
}

func NewWeightRoundRobin() *WeightRoundRobin {
	return &WeightRoundRobin{
		Nodes: []*Node{},
	}
}
func (r *WeightRoundRobin) AddNode(node *Node) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Nodes = append(r.Nodes, node)
	r.totalWeight += node.Weight
}
func (r *WeightRoundRobin) Print() {
	for _, node := range r.Nodes {
		fmt.Print(node.CurrentWeight, " ")
	}
}

func (r *WeightRoundRobin) Next() *Node {
	r.mu.Lock()
	defer r.mu.Unlock()
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
