package main

import "fmt"

type (
	Vertix struct {
		from, to int
		data     interface{}
	}
	Graph struct {
		// this is the number of nodes that the graph can contain
		// the number of vertices can be larger
		size int

		// this contains the matrix, all in a single slice
		// the size of it is size*size
		// the values in the matrix are offsets into the vertices slice with information about this relationship
		// direction and whether a vertix exists can be known by just checking for non-zero values in this slice
		matrix []int

		// keeps the nodes
		nodes []interface{}

		vertices []Vertix
	}

	intPath []int
	Path    []Vertix
)

func NewGraph(size int) *Graph {
	return &Graph{
		size:   size,
		matrix: make([]int, size*size),
		nodes:  make([]interface{}, size),
	}
}

func (g *Graph) AddNode(i int, node interface{}) error {
	if g.size < i {
		return fmt.Errorf("node too large")
	}
	if g.nodes[i] != nil {
		return fmt.Errorf("node already taken")
	}
	g.nodes[i] = node
	return nil
}

func (g *Graph) matrixOffset(from, to int) int {
	return from*g.size + to
}

func (g *Graph) GetVertix(from, to int) interface{} {
	vIdx := g.matrix[from*g.size+to]
	if vIdx == 0 {
		return nil
	}
	return g.vertices[vIdx-1]
}

func (g *Graph) AddVertix(from, to int, data interface{}) error {
	if g.size <= from || g.size <= to {
		return fmt.Errorf("node too large")
	}

	idx := g.matrixOffset(from, to)
	if g.matrix[idx] != 0 {
		return nil
	}
	g.vertices = append(g.vertices, Vertix{
		from: from,
		to:   to,
		data: data,
	})
	g.matrix[idx] = len(g.vertices)
	return nil
}

func (g *Graph) NeighboursFrom(from int) (result []int) {
	for idx := from * g.size; idx < (from+1)*g.size; idx++ {
		if g.matrix[idx] != 0 {
			result = append(result, idx-from*g.size)
		}
	}
	return
}

func (g *Graph) createVertixPath(ints intPath) (result []Vertix) {
	last := -1
	for _, node := range ints {
		if last != -1 {
			vIdx := g.matrix[g.matrixOffset(last, node)]
			result = append(result, g.vertices[vIdx-1])
		}
		last = node
	}
	return
}

func (i intPath) contains(node int) bool {
	for _, x := range i {
		if node == x {
			return true
		}
	}
	return false
}

func (g *Graph) PathsBetween(from, to []int) (result []Path) {
	// TODO: lots of optimization possibilities here
	var todo []intPath
	for _, node := range from {
		todo = append(todo, intPath{node})
	}

	for len(todo) > 0 {
		current := todo[0]
		todo = todo[1:]
		lastNode := current[len(current)-1]
	neighbours:
		for _, otherNode := range g.NeighboursFrom(lastNode) {
			if current.contains(otherNode) {
				// we don't want loops
				continue
			}
			thisPath := append(current, otherNode)
			for _, dst := range to {
				if otherNode == dst {
					// we found a path!
					result = append(result, g.createVertixPath(thisPath))
					continue neighbours
				}
			}
			todo = append(todo, thisPath)
		}
	}
	return
}

func (g *Graph) GetNode(node int) interface{} {
	return g.nodes[node]
}
