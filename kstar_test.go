package kstar

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	testutils "github.com/jcasado94/kstar/testUtils"
)

const (
	kstarTestPath = "test/kstar/"
	k             = "k"
	inExtension   = ".in"
)

// TestOutputKstar represents a kstar test output file structure
type TestOutputKstar struct {
	TestName string
	Paths    []TestPath
}

func (to TestOutputKstar) Name() string {
	return to.TestName
}

func (to TestOutputKstar) Marshall() ([]byte, error) {
	return json.MarshalIndent(to, "", "    ")
}

func (to *TestOutputKstar) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, to)
}

func (to *TestOutputKstar) New(args ...interface{}) {
	to.TestName = args[0].(string)
	to.Paths = args[1].([]TestPath)
}

func (to *TestOutputKstar) TestFolderPath() string {
	return kstarTestPath
}

type TestPath struct {
	Edges []*Edge
	Cost  float64
}

func getPathCost(path []*Edge, connections map[int]map[int][]float64) float64 {
	cost := 0.0
	for _, edge := range path {
		cost += connections[edge.U][edge.V][edge.I]
	}
	return cost
}

type kstarTest struct {
	tg testutils.TestGraph
	k  int
}

func TestAllKstarInstances(t *testing.T) {
	tgs := generateTests()
	for _, tg := range tgs {
		paths := Run(tg.tg, tg.k)
		tPaths := make([]TestPath, 0)
		for _, path := range paths {
			tPaths = append(tPaths, TestPath{
				Edges: path,
				Cost:  getPathCost(path, tg.tg.Connections()),
			})
		}
		to := new(TestOutputKstar)
		found := testutils.ReadTestOutput(to, tg.tg.TestName, tg.tg.TestName, tPaths)
		if found {
			// test
			expectedPaths := to.Paths
			for i, expectedPath := range expectedPaths {
				path := paths[i]
				for j, expectedEdge := range expectedPath.Edges {
					edge := path[j]
					if !expectedEdge.equals(*edge) {
						t.Fatalf("Test %s failed! Path %d was\n%s\n, but expected\n%s", tg.tg.TestName, i, printPath(path), printPath(expectedPath.Edges))
					}
				}
			}
		}
	}
}

func generateTests() (kstgs []kstarTest) {
	tgs := testutils.GenerateTests(datasetPath)
	kstgs = make([]kstarTest, 0)
	for _, tg := range tgs {
		file, err := os.Open(kstarTestPath + tg.TestName + inExtension)
		if err != nil {
			continue
		}
		defer file.Close()

		baseName := tg.TestName
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			vals := strings.Fields(line)
			if vals[0] == k {
				k, err := strconv.Atoi(vals[1])
				if err != nil {
					log.Fatal(err)
				}
				tg.TestName = fmt.Sprintf("%s.%s", baseName, vals[1])
				kstgs = append(kstgs, kstarTest{
					k:  k,
					tg: tg,
				})
			}
		}
	}
	return kstgs
}

func printPath(path []*Edge) (p string) {
	for _, edge := range path {
		p += fmt.Sprintf(" %v", *edge)
	}
	return p
}
