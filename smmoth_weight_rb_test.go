package smmoth_weight_rb

import (
	"testing"
)

func TestAddNode(t *testing.T) {
	//s := []string{"A", "B", "A", "C", "A", "B", "A", "A", "B", "A", "A", "B", "A"}
	r := NewWeightRoundRobinWithRing2()
	A := NewNode("A", 5, nil)
	B := NewNode("B", 2, nil)
	C := NewNode("C", 1, nil)
	r.AddNode(A)
	r.AddNode(B)
	r.AddNode(C)

	for i := 0; i < 10; i++ {
		n := r.Next()
		t.Logf("[%d]n:%s", i, n.Host)
		//if n.Host != s[i] {
		//	t.Errorf("[%d]:expected %s, got %s", i, s[i], n.Host)
		//}
	}
}
func BenchmarkWeightRoundRobinWithRing_Next(b *testing.B) {
	r := NewWeightRoundRobinWithRing2()
	A := NewNode("A", 4, nil)
	B := NewNode("B", 2, nil)
	C := NewNode("C", 3, nil)
	D := NewNode("D", 2, nil)
	E := NewNode("E", 1, nil)
	r.AddNode(A)
	r.AddNode(B)
	r.AddNode(C)
	r.AddNode(D)
	r.AddNode(E)
	for i := 0; i < b.N; i++ {
		r.Next()
	}
}

func BenchmarkWeightRoundRobinWithRing2_Next(b *testing.B) {
	r := NewWeightRoundRobinWithRing()
	A := NewNode("A", 4, nil)
	B := NewNode("B", 2, nil)
	C := NewNode("C", 3, nil)
	D := NewNode("D", 2, nil)
	E := NewNode("E", 1, nil)
	r.AddNode(A)
	r.AddNode(B)
	r.AddNode(C)
	r.AddNode(D)
	r.AddNode(E)
	for i := 0; i < b.N; i++ {
		r.Next()
	}
}

func BenchmarkWeightRoundRobin_Next2(b *testing.B) {
	r := NewWeightRoundRobin()
	A := NewNode("A", 4, nil)
	B := NewNode("B", 2, nil)
	C := NewNode("C", 3, nil)
	D := NewNode("D", 2, nil)
	E := NewNode("E", 1, nil)
	r.AddNode(A)
	r.AddNode(B)
	r.AddNode(C)
	r.AddNode(D)
	r.AddNode(E)
	for i := 0; i < b.N; i++ {
		r.Next()
	}
}
