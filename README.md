go-dgraph
=========

Package `dgraph` provides a very simple directed graph with the sole purpose of
returning nodes and cycles in dependency order with the strongly connected
components (Tarjan's) algorithm.

The StrongComponents returns graph components in dependency order. If the graph
has no cycles, each component will have a single element. Otherwise, all nodes
in a cycle are grouped into one "strong" component.

Documentation
-------------

[![GoDoc](https://godoc.org/github.com/twmb/go-dgraph?status.svg)](https://godoc.org/github.com/twmb/go-dgraph)
