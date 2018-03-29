package main

import (
	"container/heap"
	"container/list"
	"log"
	"sort"
)

// Edge is a connection bitween two points. Every edge has two terminal vertecies and a weight between the vertecies.
type Edge struct {
	source   string
	sink     string
	capacity float32
}

type edgeSlice []Edge

func (e edgeSlice) Len() int {
	return len(e)
}

func (e edgeSlice) Less(i, j int) bool {
	return e[i].capacity < e[j].capacity
}

func (e edgeSlice) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// Graph is a complete graph with vertecies and edges between them.
type Graph struct {
	vertecies map[string][]Edge
}

// AddVertex adds a vertex to the graph
func (g *Graph) AddVertex(vertex string) {
	g.vertecies[vertex] = make([]Edge, 0)
}

// AddEdge adds an edge to the graph
func (g *Graph) AddEdge(source string, sink string, capacity float32) {
	edge := Edge{source, sink, capacity}
	g.vertecies[source] = append(g.vertecies[source], edge)
}

func (g *Graph) getVertecies() []string {
	var edges = make([]string, len(g.vertecies))
	i := 0
	for k := range g.vertecies {
		edges[i] = k
		i++
	}

	return edges
}

func (g Graph) getEdges(vertex string) ([]Edge, bool) {
	edges, ok := g.vertecies[vertex]
	return edges, ok
}

type vertexQueue []string

func (vq vertexQueue) Len() int            { return len(vq) }
func (vq vertexQueue) Less(i, j int) bool  { return vq[i] < vq[j] }
func (vq vertexQueue) Swap(i, j int)       { vq[i], vq[j] = vq[j], vq[i] }
func (vq *vertexQueue) Push(x interface{}) { *vq = append(*vq, x.(string)) }
func (vq *vertexQueue) Pop() interface{} {
	items := *vq
	n := len(*vq)
	x := items[n-1]
	*vq = items[0 : n-1]
	return x
}

func depthFirstSearch(g traverseGraph, from, to string) bool {
	verteciesNext := make(vertexQueue, 1)
	heap.Init(&verteciesNext)
	heap.Push(&verteciesNext, from)

	for verteciesNext.Len() > 0 {
		vertex := heap.Pop(&verteciesNext).(string)
		edges, hasVertex := g.vertecies[vertex]
		for hasVertex {
			for e := edges.Front(); e != nil; e = e.Next() {
				edge := e.Value.(Edge)
				if edge.sink == to {
					return true
				}

				heap.Push(&verteciesNext, edge.sink)
			}
		}
	}

	return false
}

// Kruskals performs kruskal's algorithm
func Kruskals(g *Graph) []Edge {
	vertecies := g.getVertecies()
	numVertecies := len(vertecies)

	edges := make([]Edge, 0)

	for _, vertex := range vertecies {
		edgeInVertex, _ := g.getEdges(vertex)
		for _, edge := range edgeInVertex {
			has := false
			for _, e := range edges {
				has = has || e == edge
			}

			if has == false {
				edges = append(edges, edge)
			}
		}
	}

	var readyToSort edgeSlice = edges
	sort.Sort(readyToSort)

	var gCopy = copy(*g)

	pastVertecies := make(map[string]bool)
	a := make([]Edge, numVertecies-1) // resulting set
	it := 0
	for i := 0; i < len(readyToSort); i++ {
		x := readyToSort[i]
		_, hasSource := pastVertecies[x.source]
		_, hasSink := pastVertecies[x.sink]
		if !hasSource || !hasSink {
			a[it] = x

			pastVertecies[x.source] = true
			pastVertecies[x.sink] = true

			it++

			edgesForVertex, hasVertex := gCopy.vertecies[x.sink]
			if hasVertex {
				var e *list.Element
				for e = edgesForVertex.Front(); e != nil; e = e.Next() {
					value := e.Value
					edge := value.(Edge) // TODO: runtime exception
					if edge.source == x.source && edge.sink == x.sink {
						break
					}
				}

				edgesForVertex.Remove(e)
			}
		} else {
			if depthFirstSearch(gCopy, x.source, x.sink) {
				a[it] = x
				pastVertecies[x.source] = true
				pastVertecies[x.sink] = true
				it++

				if it == numVertecies-1 {
					break
				}
			} else {
				log.Print("Dropped")
				log.Print(x.source)
				log.Println(x.sink)
			}
		}

		if it == numVertecies-1 {
			break
		}
	}

	return a
}

// New creates a new graph and assigns its values
func New() Graph {
	vertecies := make(map[string][]Edge)
	g := Graph{vertecies}
	return g
}

type traverseGraph struct {
	vertecies map[string]*list.List
}

// Copy copies all content from g to a new graph
func copy(g Graph) traverseGraph {
	vertecies := make(map[string]*list.List)
	var gCopy = traverseGraph{vertecies}
	for k, edges := range g.vertecies {
		traverseEdges := list.New()
		for edge := range edges {
			traverseEdges.PushBack(edge)
		}

		gCopy.vertecies[k] = traverseEdges
	}

	return gCopy
}

func buildExampleGraph() Graph {
	var g = New()
	var vertecies = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	for _, vertex := range vertecies {
		g.AddVertex(vertex)
	}

	g.AddEdge("a", "b", 4)
	g.AddEdge("a", "h", 8)
	g.AddEdge("b", "c", 8)
	g.AddEdge("b", "h", 11)
	g.AddEdge("c", "d", 7)
	g.AddEdge("c", "i", 2)
	g.AddEdge("c", "f", 4)
	g.AddEdge("d", "e", 9)
	g.AddEdge("d", "f", 14)
	g.AddEdge("e", "f", 10)
	g.AddEdge("f", "g", 2)
	g.AddEdge("g", "h", 1)
	g.AddEdge("g", "i", 6)
	g.AddEdge("h", "i", 7)
	return g
}

func main() {
	g := buildExampleGraph()
	edges := Kruskals(&g)
	for _, edge := range edges {
		log.Printf("Edge from %s to %s with cost %v", edge.source, edge.sink, edge.capacity)
	}

	// Edges in MST
	// (d, e)
	// (a, h)
	// (c, d)
	// (c, f)
	// (a, b)
	// (f, g)
	// (c, i)
	// (g, h)
}
