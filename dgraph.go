// Package `dgraph` provides a very simple directed graph with the sole purpose
// of returning nodes and cycles in dependency order with the strongly
// connected components (Tarjan's) algorithm.
package dgraph

// Graph is a directed graph.
type Graph struct {
	out map[int][]int
	in  map[int][]int
}

// New returns a new Graph.
func New() *Graph {
	return &Graph{
		out: make(map[int][]int),
		in:  make(map[int][]int),
	}
}

// Add adds a node to the graph if it does not exist.
func (g *Graph) Add(node int) {
	if _, exists := g.out[node]; !exists {
		g.out[node] = []int{}
		g.in[node] = []int{}
	}
}

// Remove removes a node from the graph, unlinking everything it was linked to.
func (g *Graph) Remove(node int) {
	for _, in := range g.in[node] {
		out := g.out[in]
		g.out[in] = rmIdx(out, intIdx(out, node))
	}
	for _, out := range g.out[node] {
		in := g.in[out]
		g.in[out] = rmIdx(in, intIdx(in, node))
	}
	delete(g.out, node)
	delete(g.in, node)
}

// Link adds an edge from src to dst, creating the nodes if they do not exist.
func (g *Graph) Link(src, dst int) {
	out, exists := g.out[src]
	g.out[src] = append(out, dst)
	if !exists {
		g.in[src] = []int{}
	}

	in, exists := g.in[dst]
	g.in[dst] = append(in, src)
	if !exists {
		g.out[dst] = []int{}
	}
}

func rmIdx(in []int, idx int) []int {
	in[idx] = in[len(in)-1]
	in = in[:len(in)-1]
	return in
}

func intIdx(in []int, needle int) int {
	r := -1
	for i, v := range in {
		if v == needle {
			r = i
			break
		}
	}
	return r
}

// Unlink removes an edge from src to dst if the edge exists.
func (g *Graph) Unlink(src, dst int) {
	out := g.out[src]
	if idx := intIdx(out, dst); idx != -1 {
		g.out[src] = rmIdx(out, idx)
		ins := g.in[dst]
		g.in[dst] = rmIdx(ins, intIdx(ins, src))
	}
}

// StrongComponents returns all strong components of the graph in dependency
// order. If a graph has no cycles, this will return each node individually.
//
// Note that this is a topological sort: if a graph has no cycles, the order of
// the returned single element components is the dependency order of the graph.
func (g *Graph) StrongComponents() [][]int {
	seen := make(map[int]struct{}, len(g.out))
	saw := func(node int) bool {
		if _, saw := seen[node]; !saw {
			seen[node] = struct{}{}
			return false
		}
		return true
	}

	order := make([]int, 0, len(g.out))
	for node := range g.out {
		if !saw(node) {
			dfs(g.out, node, saw, func(ord int) { order = append(order, ord) })
		}
	}

	saw = func(node int) bool {
		if _, remains := seen[node]; remains {
			delete(seen, node)
			return false
		}
		return true
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
	graph map[int][]int,
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
