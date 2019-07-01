package testutils

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	commentMark   = "c"
	problemMark   = "p"
	heuristicMark = "h"
	edgeMark      = "e"
)

const graphExtension = ".graph"

// TestGraph represents a test graph read from a file. Implements kstar.Graph.
type TestGraph struct {
	TestName string
	s, t     int
	fValues  map[int]float64
	graph    map[int]map[int][]float64
}

// Nodes returns the set of nodes of the graph.
func (tg TestGraph) Nodes() []int {
	nodes := make([]int, len(tg.fValues))
	i := 0
	for k := range tg.fValues {
		nodes[i] = k
		i++
	}
	return nodes
}

// Connections returns the problem connections
func (tg TestGraph) Connections(n int) map[int][]float64 {
	return tg.graph[n]
}

// S returns the problem s node
func (tg TestGraph) S() int {
	return tg.s
}

// T returns the problem t node
func (tg TestGraph) T() int {
	return tg.t
}

// FValue returns the heuristic cost from node n to t.
func (tg TestGraph) FValue(n int) float64 {
	return tg.fValues[n]
}

// GenerateTests generates all the TestGraph instances from all the files at the specified path.
func GenerateTests(path string) []TestGraph {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var tgs []TestGraph
	for _, file := range files {
		fLen := len(file.Name())
		if file.Name()[fLen-len(graphExtension):fLen] == graphExtension {
			tgs = append(tgs, generateTest(path+file.Name(), file.Name()[0:fLen-len(graphExtension)]))
		}
	}

	return tgs

}

func generateTest(filePath, testName string) (tg TestGraph) {

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	tg.TestName = testName
	tg.fValues = make(map[int]float64, 0)
	tg.graph = make(map[int]map[int][]float64, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		vals := strings.Fields(line)
		if vals[0] == problemMark {
			tg.s, err = strconv.Atoi(vals[3])
			tg.t, err = strconv.Atoi(vals[4])
			if err != nil {
				log.Fatal(err)
			}
			tg.graph = make(map[int]map[int][]float64, 0)
		} else if vals[0] == heuristicMark {
			nodeInt, err := strconv.Atoi(vals[1])
			fValue, err := strconv.ParseFloat(vals[2], 64)
			if err != nil {
				log.Fatal(err)
			}
			tg.fValues[nodeInt] = fValue
		} else if vals[0] == edgeMark {
			u, err := strconv.Atoi(vals[1])
			v, err := strconv.Atoi(vals[2])
			c, err := strconv.ParseFloat(vals[3], 64)
			if err != nil {
				log.Fatal(err)
			}

			if _, ok := tg.graph[u]; !ok {
				tg.graph[u] = make(map[int][]float64, 0)
			}
			if _, ok := tg.graph[u][v]; !ok {
				tg.graph[u][v] = make([]float64, 0)
			}

			tg.graph[u][v] = append(tg.graph[u][v], c)

		}

	}

	return

}

func fillFValues(tg *TestGraph) {
	for node := range tg.graph {
		if _, ok := tg.fValues[node]; !ok {
			tg.fValues[node] = 0
		}
	}
}
