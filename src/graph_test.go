package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {

	g := NewGraph(3)
	require.NoError(t,
		g.AddNode(0, "node 0"))
	require.NoError(t,
		g.AddNode(1, "node 1"))
	require.NoError(t,
		g.AddNode(2, "node 2"))

	require.NoError(t,
		g.AddVertix(0, 1, "0->1"))
	require.NoError(t,
		g.AddVertix(1, 2, "1->2"))

	apa := g.PathsBetween([]int{0}, []int{2})
	fmt.Printf("%v\n", apa)
}

func TestBigGraph(t *testing.T) {
	g := NewGraph(9)
	addNode(t, g, 0, "node 0")
	addNode(t, g, 1, "node 1")
	addNode(t, g, 2, "node 2")
	addNode(t, g, 3, "node 3")
	addNode(t, g, 4, "node 4")
	addNode(t, g, 5, "node 5")
	addNode(t, g, 6, "node 6")
	addNode(t, g, 7, "node 7")
	addNode(t, g, 8, "node 8")

	addV(t, g, 0, 3)
	addV(t, g, 0, 4)
	addV(t, g, 1, 4)
	addV(t, g, 4, 1)
	addV(t, g, 2, 5)
	addV(t, g, 4, 7)
	addV(t, g, 4, 8)

	paths := g.PathsBetween([]int{0, 1, 2}, []int{6, 7, 8})

	fmt.Printf("%v\n", paths)
}

func addNode(t *testing.T, g *Graph, i int, node string) {
	require.NoError(t,
		g.AddNode(i, node))
}
func addV(t *testing.T, g *Graph, from, to int) {
	require.NoError(t,
		g.AddVertix(from, to, ""))
}

func TestNeigbour(t *testing.T) {
	g := NewGraph(3)
	require.NoError(t,
		g.AddNode(0, "node 0"))
	require.NoError(t,
		g.AddNode(1, "node 1"))
	require.NoError(t,
		g.AddNode(2, "node 2"))

	require.NoError(t,
		g.AddVertix(0, 1, "0->1"))
	require.NoError(t,
		g.AddVertix(0, 2, "0->2"))

	assert.Equal(t, []int{1, 2}, g.NeighboursFrom(0))
	assert.Empty(t, g.NeighboursFrom(1))
}
