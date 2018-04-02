package main

import (
	"container/list"
	"testing"
)

func Test_depthFirstSearch_inputLacksTerminalVertex(t *testing.T) {
	a := list.New()
	a.PushBack(Edge{"a", "b", 1})

	b := list.New()
	b.PushBack(Edge{"b", "c", 1})

	c := list.New()
	vertecies := make(map[string]*list.List)
	vertecies["a"] = a
	vertecies["b"] = b
	vertecies["c"] = c

	var g = traverseGraph{vertecies}

	result := depthFirstSearch(g, "a", "d")
	if result {
		t.Error("Expected to find no path from a->d, but there is one")
	}
}
func Test_depthFirstSearch(t *testing.T) {
	a := list.New()
	a.PushBack(Edge{"a", "b", 1})

	b := list.New()
	b.PushBack(Edge{"b", "c", 1})

	c := list.New()
	vertecies := make(map[string]*list.List)
	vertecies["a"] = a
	vertecies["b"] = b
	vertecies["c"] = c

	var g = traverseGraph{vertecies}

	result := depthFirstSearch(g, "a", "c")
	if !result {
		t.Error("Expected to find a path from a->b->c")
	}
}

func Test_createEmptyGraphWithVertecies(t *testing.T) {
	a := make([]Edge, 1)
	a[0] = Edge{"a", "b", 1}

	b := make([]Edge, 1)
	b[0] = Edge{"b", "c", 1}

	var c []Edge
	vertecies := make(map[string][]Edge)
	vertecies["a"] = a
	vertecies["b"] = b
	vertecies["c"] = c
	g := Graph{vertecies}
	result := createEmptyGraphWithVertecies(g)

	if result.vertecies == nil {
		t.Error("Expected type traverseGraph.vertecies but got nil")
	}

	_, hasA := result.vertecies["a"]
	_, hasB := result.vertecies["b"]
	_, hasC := result.vertecies["c"]
	if !hasA {
		t.Error("result is lacking vertex a")
	}

	if !hasB {
		t.Error("result is lacking vertex b")
	}

	if !hasC {
		t.Error("result is lacking vertex c")
	}
}

func Test_Kruskals(t *testing.T) {
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

	result := Kruskals(&g)

	if len(result) != 3 {
		t.Errorf("Expected graph size to be number of vertecies N - 1: 3. Actually got %v", len(result))
	}
}

func Test_Kruskals_InputHasCycle(t *testing.T) {
	var g = New()
	var vertecies = []string{"a", "b", "c"}
	for _, vertex := range vertecies {
		g.AddVertex(vertex)
	}

	g.AddEdge("a", "b", 1)
	g.AddEdge("b", "c", 1)
	g.AddEdge("c", "b", 1)

	result := Kruskals(&g)

	if len(result) != 2 {
		t.Errorf("Expected graph size to be number of vertecies N - 1: 2. Actually got %v", len(result))
	}
}

func Test_Kruskals_InputHasCycle2(t *testing.T) {
	var g = New()
	var vertecies = []string{"a", "b", "c", "d"}
	for _, vertex := range vertecies {
		g.AddVertex(vertex)
	}

	g.AddEdge("a", "b", 1)
	g.AddEdge("c", "d", 1)
	g.AddEdge("b", "c", 1)
	g.AddEdge("c", "b", 1)

	result := Kruskals(&g)

	if len(result) != 3 {
		t.Errorf("Expected graph size to be number of vertecies N - 1: 3. Actually got %v", len(result))
	}
}
