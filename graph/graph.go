package graph

import (
	"bytes"
	"fmt"
)

type Node interface {
	Name() string
	Dot() string
}

type Edge interface {
	Start() string
	End() string
	//	Name() string
	Dot() string
}

type Graph struct {
	nodes map[string]Node
	edges map[string][]Edge
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]Node),
		edges: make(map[string][]Edge),
	}
}

func (g *Graph) AddEdge(e Edge) {
	g.edges[e.Start()] = append(
		g.edges[e.Start()],
		e,
	)
}

func (g *Graph) Clone() *Graph {
	ng := &Graph{
		nodes: make(map[string]Node),
		edges: make(map[string][]Edge),
	}
	for key, val := range g.nodes {
		ng.nodes[key] = val
	}
	for key, val := range g.edges {
		nval := make([]Edge, len(val))
		copy(nval, val)
		ng.edges[key] = nval
	}
	return ng
}

func (g *Graph) AddNode(n Node) {
	g.nodes[n.Name()] = n
}

func (g *Graph) Node(name string) Node {
	return g.nodes[name]
}

func (g *Graph) Edges(node string) []Edge {
	return g.edges[node]
}

func (g *Graph) ExportDot() (string, error) {
	var buf bytes.Buffer
	head := "digraph \"test\" {\n\tgraph[\n\t\tcharset=\"UTF-8\";\n\t];\n\n"
	tail := "}\n"
	buf.WriteString(head)

	for _, node := range g.nodes {
		_, err := buf.WriteString(fmt.Sprintf("\t%s", node.Dot()))
		if err != nil {
			return "", err
		}
	}
	for _, edges := range g.edges {
		for _, edge := range edges {
			_, err := buf.WriteString(fmt.Sprintf("\t%s", edge.Dot()))
			if err != nil {
				return "", err
			}
		}
	}
	buf.WriteString(tail)
	return string(buf.Bytes()), nil
}
