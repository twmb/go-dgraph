// Package dgraph provides a very simple directed graph with the sole purpose
// of returning nodes and cycles in dependency order with the strongly
// connected components (Tarjan's) algorithm.
//
// The StrongComponents returns graph components in dependency order. If the
// graph has no cycles, each component will have a single element. Otherwise,
// all nodes in a cycle are grouped into one "strong" component.
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
// much as necessary to fit the node.
func (g *Graph) Add(node int) {
	need := node + 1
	if len(g.out) < need {
		g.out = append(g.out, make([][]int, need-len(g.out))...)
		g.in = append(g.in, make([][]int, need-len(g.in))...)
	}
}

// Link adds an edge from src to dst, creating the nodes if they do not exist
// and extending the graph as much as necessary to fit those nodes.
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
	d := dfser{
		graph: g.out,
		order: make([]int, 0, len(g.out)),
		flip:  make([]bool, len(g.out)),
	}

	for node := range g.out {
		if !d.sawTrue(node) {
			d.dfs1(node)
		}
	}

	order := d.order

	d.order = make([]int, 0, len(order))
	d.graph = g.in

	components := make([][]int, 0, len(order))
	for i := len(order) - 1; i >= 0; i-- {
		if node := order[i]; !d.sawFalse(node) {
			d.dfs2(node)
			used := len(d.order)
			component := d.order[:used:used]
			components = append(components, component)
			d.order = d.order[used:]
		}
	}

	return components
}

// dfser performs depth first search on the graph, appending bottom nodes to
// order.
//
// Since SCC requires two dfs runs, we keep track of nodes seen with flip: on
// the first run, we swap flip nodes to true, on the second back to false.
//
// The original impl just passed graph and two closures to one dfs function,
// but switching to methods on a struct proved a big (>30%) performance
// increase.
type dfser struct {
	graph [][]int
	order []int
	flip  []bool
}

func (d *dfser) sawTrue(node int) bool {
	r := d.flip[node]
	d.flip[node] = true
	return r
}
func (d *dfser) sawFalse(node int) bool {
	r := d.flip[node]
	d.flip[node] = false
	return !r
}

func (d *dfser) dfs1(node int) {
	for _, neighbor := range d.graph[node] {
		if !d.sawTrue(neighbor) {
			d.dfs1(neighbor)
		}
	}
	d.order = append(d.order, node)
}

func (d *dfser) dfs2(node int) {
	for _, neighbor := range d.graph[node] {
		if !d.sawFalse(neighbor) {
			d.dfs2(neighbor)
		}
	}
	d.order = append(d.order, node)
}
