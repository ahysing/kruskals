package main

import (
	"container/heap"
	"container/list"
	"log"
	"sort"
)

// Edge is a connection bitween two vertecies. Every edge has a weight between the vertecies.
type Edge struct {
	source string
	sink   string
	weight float32
}

type edgeSlice []Edge

func (e edgeSlice) Len() int {
	return len(e)
}

func (e edgeSlice) Less(i, j int) bool {
	return e[i].weight < e[j].weight
}

func (e edgeSlice) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func sortEdges(edges []Edge) {
	var readyToSort edgeSlice = edges
	sort.Sort(readyToSort)
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
func (g *Graph) AddEdge(source string, sink string, weight float32) {
	edge := Edge{source, sink, weight}
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

func graphContainsNoCycles(g traverseGraph, from, to string) bool {
	return depthFirstSearch(g, from, to) == false
}

func depthFirstSearch(g traverseGraph, from, to string) bool {
	verteciesNext := make(vertexQueue, 0)
	heap.Init(&verteciesNext)
	heap.Push(&verteciesNext, from)

	for verteciesNext.Len() > 0 {
		vertex := heap.Pop(&verteciesNext).(string)
		edges, hasVertex := g.vertecies[vertex]
		if hasVertex {
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

func addUndirectedEdge(tg *traverseGraph, edge Edge) {
	addEdge(tg, edge)
	reverseEdge := Edge{edge.sink, edge.source, edge.weight}
	addEdge(tg, reverseEdge)
}
func addEdge(tg *traverseGraph, x Edge) {
	edgesForVertex, hasVertex := tg.vertecies[x.source]
	if !hasVertex {
		edgesForVertex = list.New()
		tg.vertecies[x.source] = edgesForVertex
	}

	edgesForVertex.PushBack(x)
}

func getUniqueEdges(g *Graph, vertecies []string) []Edge {
	edges := make([]Edge, 0)
	for _, vertex := range vertecies {
		edgeInVertex, hasVertex := g.getEdges(vertex)
		if hasVertex {
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
	}

	return edges
}

// Kruskals performs kruskal's algorithm
func Kruskals(g *Graph) []Edge {
	vertecies := g.getVertecies()
	numVertecies := len(vertecies)

	edges := getUniqueEdges(g, vertecies)
	sortEdges(edges)

	var gCopy = createEmptyGraphWithVertecies(*g)

	traversedVertecies := make(map[string]bool)
	a := make([]Edge, numVertecies-1) // resulting set
	it := 0
	for _, currentEdge := range edges {
		_, hasSource := traversedVertecies[currentEdge.source]
		_, hasSink := traversedVertecies[currentEdge.sink]
		if hasSource && hasSink {
			if graphContainsNoCycles(gCopy, currentEdge.source, currentEdge.sink) {
				a[it] = currentEdge
				traversedVertecies[currentEdge.source] = true
				traversedVertecies[currentEdge.sink] = true
				addUndirectedEdge(&gCopy, currentEdge)
				it++
			} else {
				log.Printf("Dropped %v -> %v", currentEdge.source, currentEdge.sink)
			}
		} else {
			a[it] = currentEdge
			traversedVertecies[currentEdge.source] = true
			traversedVertecies[currentEdge.sink] = true
			addUndirectedEdge(&gCopy, currentEdge)
			it++
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

func createEmptyGraphWithVertecies(g Graph) traverseGraph {
	vertecies := make(map[string]*list.List)
	var gCopy = traverseGraph{vertecies}
	for k := range g.vertecies {
		traverseEdges := list.New()
		gCopy.vertecies[k] = traverseEdges
	}

	return gCopy
}

func buildExampleGraph() Graph {
	var g = New()
	var vertecies = []string{"a", "b", "c", "d"}
	for _, vertex := range vertecies {
		g.AddVertex(vertex)
	}

	g.AddEdge("a", "b", 1)
	g.AddEdge("b", "c", 1)
	g.AddEdge("c", "a", 1)
	g.AddEdge("c", "d", 1)
	g.AddEdge("d", "a", 1)

	return g
}

func main() {
	g := buildExampleGraph()
	edges := Kruskals(&g)
	for _, edge := range edges {
		log.Printf("Edge from %s to %s with cost %v", edge.source, edge.sink, edge.weight)
	}
}
