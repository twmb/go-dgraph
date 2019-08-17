// Package `dgraph` provides a very simple directed graph with the sole purpose
// of returning nodes and cycles in dependency order with the strongly
// connected components (Tarjan's) algorithm.
package dgraph

// Graph is a directed graph.
type Graph struct {
	out [][]int
	in  [][]int
}

// New returns a new graph with no nodes.
func New() *Graph {
	return NewSize(0)
}

// NewSize returns a new graph with a set amount of starting nodes.
func NewSize(nodes int) *Graph {
	return &Graph{
		out: make([][]int, nodes),
		in:  make([][]int, nodes),
	}
}

// Add adds a node to the graph if it does not exist, extending the graph as
// much as necessary to fit the node. Thus, this will add nodes as necessary
// if the graph is too small.
func (g *Graph) Add(node int) {
	need := node + 1
	if len(g.out) < need {
		new := make([][]int, need)
		copy(new, g.out)
		g.out = new
		new = make([][]int, need)
		copy(new, g.in)
		g.in = new
	}
}

// Link adds an edge from src to dst, creating the nodes if they do not exist.
func (g *Graph) Link(src, dst int) {
	if src > dst {
		g.Add(src)
	} else {
		g.Add(dst)
	}

	g.out[src] = append(g.out[src], dst)
	g.in[dst] = append(g.in[dst], src)
}

// StrongComponents returns all strong components of the graph in dependency
// order. If a graph has no cycles, this will return each node individually.
//
// Note that this is a topological sort: if a graph has no cycles, the order of
// the returned single element components is the dependency order of the graph.
func (g *Graph) StrongComponents() [][]int {
	flip := make([]bool, len(g.out))
	saw := func(node int) bool {
		r := flip[node]
		flip[node] = true
		return r
	}

	order := make([]int, 0, len(g.out))
	for node := range g.out {
		if !saw(node) {
			dfs(g.out, node, saw, func(ord int) { order = append(order, ord) })
		}
	}

	saw = func(node int) bool {
		r := flip[node]
		flip[node] = false
		return !r
	}

	components := make([][]int, 0, len(order))
	nodes := make([]int, 0, len(order))
	for i := len(order) - 1; i >= 0; i-- {
		if node := order[i]; !saw(node) {
			dfs(g.in, node, saw, func(ord int) { nodes = append(nodes, ord) })
			used := len(nodes)
			component := nodes[:used:used]
			components = append(components, component)
			nodes = nodes[used:]
		}
	}

	return components
}

func dfs(
	graph [][]int,
	node int,
	saw func(int) bool,
	fn func(int),
) {
	for _, neighbor := range graph[node] {
		if !saw(neighbor) {
			dfs(graph, neighbor, saw, fn)
		}
	}
	fn(node)
}
