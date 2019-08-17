package dgraph_test

import (
	"fmt"
	"sort"

	"github.com/twmb/go-dgraph"
)

func ExampleGraph() {
	nodeVals := map[int]string{
		0: "foo",
		1: "bar",
		2: "baz",
		3: "these",
		4: "are",
		5: "my",
		6: "nodes",
	}
	g := dgraph.NewSize(len(nodeVals))
	g.Link(0, 2)
	g.Link(2, 1)
	g.Link(1, 3)
	g.Link(3, 4)
	g.Link(4, 5)
	g.Link(5, 6)
	g.Link(6, 4)
	// creates graph 0->2->1->3->[4 5 6]

	order := g.StrongComponents()

	// we can see 5 components: 4 non cycles, 1 cycle
	fmt.Printf("graph has %d components\n", len(order))

	for i := 0; i < len(order); i++ {
		component := order[i]
		if len(component) == 1 {
			node := component[0]
			fmt.Printf("[%d] => %s\n", node, nodeVals[node])
			continue
		}

		// a cycle can have many orderings; we sort for output ordering
		sort.Ints(component)
		var vals []string
		for _, node := range component {
			vals = append(vals, nodeVals[node])
		}
		fmt.Printf("cycle %v => %v\n", component, vals)
	}

	// Output:
	// graph has 5 components
	// [0] => foo
	// [2] => baz
	// [1] => bar
	// [3] => these
	// cycle [4 5 6] => [are my nodes]
}
