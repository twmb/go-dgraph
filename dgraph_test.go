package dgraph

import (
	"reflect"
	"sort"
	"testing"
)

// CLRS page 616 ed. 3
// This linkage is taken directly from my old algoimpl repo.
func testgraph() *Graph {
	g := New()

	g.Link(0, 1)
	g.Link(0, 4)
	g.Link(1, 2)
	g.Link(1, 3)
	g.Link(2, 1)
	g.Link(3, 3)
	g.Link(4, 3)
	g.Link(4, 0)
	g.Link(5, 6)
	g.Link(5, 0)
	g.Link(5, 2)
	g.Link(6, 2)
	g.Link(6, 7)
	g.Link(7, 5)

	return g
}

func TestSCC(t *testing.T) {
	g := testgraph()
	got := g.StrongComponents()

	exp := [][]int{
		{5, 6, 7},
		{0, 4},
		{1, 2},
		{3},
	}
	if len(got) != len(exp) {
		t.Errorf("got %d comps != exp %d", len(got), len(exp))
	}

	for i, gotComp := range got {
		expComp := exp[i]
		sort.Ints(gotComp)
		if !reflect.DeepEqual(gotComp, expComp) {
			t.Errorf("comp %d: got %v != exp %v", i, gotComp, expComp)
		}
	}
}

func BenchmarkSCC(b *testing.B) {
	g := testgraph()
	for i := 0; i < b.N; i++ {
		g.StrongComponents()
	}
}
