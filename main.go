package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cook-a-doodle-do/trace-comparator/graph"
)

type hoge struct {
	name string
}

func (h *hoge) Name() string {
	return h.name
}

func (h *hoge) Dot() string {
	return fmt.Sprintf("\"%s\";\n", h.name)
}

type edge struct {
	hidden bool
	start  string
	end    string
	name   string
}

func (e *edge) Start() string {
	return e.start
}

func (e *edge) End() string {
	return e.end
}

func (e *edge) Name() string {
	return e.name
}

func (e *edge) Dot() string {
	if e.hidden {
		return fmt.Sprintf("\"%s\"->\"%s\"[label=\"%s\", style=\"dashed\"];\n", e.start, e.end, e.name)
	}
	return fmt.Sprintf("\"%s\"->\"%s\"[label=\"%s\"];\n", e.start, e.end, e.name)
}

type pair struct {
	specification  string
	implementation string
}

func main() {

	//specification ================================================================
	spc := graph.NewGraph()
	spc.AddNode(&hoge{name: "0"})
	spc.AddNode(&hoge{name: "1"})
	spc.AddNode(&hoge{name: "2"})
	spc.AddNode(&hoge{name: "3"})
	spc.AddEdge(&edge{start: "0", end: "1", name: "go1"})
	spc.AddEdge(&edge{start: "1", end: "1", name: "go1", hidden: true})
	spc.AddEdge(&edge{start: "1", end: "2", name: "go2"})
	spc.AddEdge(&edge{start: "2", end: "2", name: "go2"})
	spc.AddEdge(&edge{start: "2", end: "3", name: "go3", hidden: true})
	spc.AddEdge(&edge{start: "3", end: "3", name: "go3"})
	spc.AddEdge(&edge{start: "3", end: "0", name: "go0"})
	str, err := spc.ExportDot()
	if err != nil {
		log.Fatal(err)
	}
	f0, err := os.Create("spc.dot")
	if err != nil {
		log.Fatal(err)
	}
	defer f0.Close()
	f0.WriteString(str)

	//implementation ===============================================================
	imp := graph.NewGraph() //spc.Clone()
	imp.AddNode(&hoge{name: "2"})
	imp.AddNode(&hoge{name: "3"})
	spc.AddEdge(&edge{start: "2", end: "3", name: "go2"})
	str, err = imp.ExportDot()
	if err != nil {
		log.Fatal(err)
	}
	f1, err := os.Create("imp.dot")
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()
	f1.WriteString(str)

	//result
	result := graph.NewGraph()
	f2, err := os.Create("result.dot")
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	//make queue & check list
	var field []*pair
	checked := make(map[pair]bool)
	field = append(field, &pair{specification: "0", implementation: "0"})

	for len(field) > 0 {
		// pull one value from queue
		cur := field[0]
		field = field[1:]
		// check cur pair
		checked[*cur] = true
		result.AddNode(&hoge{name: fmt.Sprintf("{%s,%s}", cur.implementation, cur.specification)})

		fmt.Println(cur)
		for _, impEdges := range imp.Edges(cur.implementation) {
			impEdge, _ := impEdges.(*edge)
			impEvent := impEdge.Name()
			impNextN := impEdge.End()

			for _, spcEdges := range spc.Edges(cur.specification) {
				spcEdge, _ := spcEdges.(*edge)
				spcEvent := spcEdge.Name()
				spcNextN := spcEdge.End()

				if spcEvent == impEvent {
					p := &pair{
						implementation: impNextN,
						specification:  spcNextN,
					}
					result.AddEdge(&edge{
						start: fmt.Sprintf("{%s,%s}", cur.implementation, cur.specification),
						end:   fmt.Sprintf("{%s,%s}", p.implementation, p.specification),
						name:  impEvent,
					})
					if _, ok := checked[*p]; !ok {
						field = append(field, p)
					}
				}
				fmt.Println(spcEvent, spcNextN)
			}
		}
	}
	//export result
	str, err = result.ExportDot()
	if err != nil {
		log.Fatal(err)
	}
	f2.WriteString(str)
}
