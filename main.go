package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cook-a-doodle-do/trace-comparator/graph"
)

type hoge struct {
	name  string
	color string
}

func (h *hoge) Name() string {
	return h.name
}

func (h *hoge) Dot() string {
	str := fmt.Sprintf("\"%s\"[shape=\"ellipse\"", h.name)
	if h.color != "" {
		str = fmt.Sprintf("%s, color=\"%s\"", str, h.color)
	}
	return str + "];\n"
}

type edge struct {
	hidden bool
	start  string
	end    string
	name   string
	color  string
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
	str := fmt.Sprintf("\"%s\"->\"%s\"[label=\"%s\"", e.start, e.end, e.name)
	if e.hidden {
		str = fmt.Sprintf("%s, style=\"dashed\"", str)
	}
	if e.color != "" {
		str = fmt.Sprintf("%s, color=\"%s\"", str, e.color)
	}
	return fmt.Sprintf("%s];\n", str)
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
	spc.AddEdge(&edge{start: "0", end: "1", name: "a"})
	//	spc.AddNode(&hoge{name: "0"})
	//	spc.AddNode(&hoge{name: "1"})
	//	spc.AddNode(&hoge{name: "2"})
	//	spc.AddEdge(&edge{start: "0", end: "1", name: "a"})
	//	spc.AddEdge(&edge{start: "2", end: "1", name: "a"})
	//	spc.AddEdge(&edge{start: "1", end: "1", name: "a"})
	//	spc.AddEdge(&edge{start: "1", end: "2", name: "b"})
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
	imp.AddNode(&hoge{name: "0"})
	imp.AddNode(&hoge{name: "3"})
	imp.AddEdge(&edge{start: "0", end: "3", name: "a", hidden: false})
	//	imp.AddNode(&hoge{name: "0"})
	//	imp.AddNode(&hoge{name: "1"})
	//	imp.AddNode(&hoge{name: "2"})
	//	imp.AddEdge(&edge{start: "0", end: "1", name: "a", hidden: false})
	//	imp.AddEdge(&edge{start: "1", end: "2", name: "a", hidden: false})
	//	imp.AddEdge(&edge{start: "2", end: "0", name: "b", hidden: false})
	//	imp.AddEdge(&edge{start: "2", end: "2", name: "a", hidden: false})
	//	imp.AddEdge(&edge{start: "1", end: "0", name: "b", hidden: false})
	//imp.AddEdge(&edge{start: "0", end: "2", name: "a", hidden: false})

	//	imp.AddEdge(&edge{start: "1", end: "1", name: "b"})

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
	root := &pair{specification: "0", implementation: "0"}
	field = append(field, root)
	checked[*root] = true
	violated := false

	for len(field) > 0 {
		// pull one value from queue
		cur := field[0]
		field = field[1:]
		fmt.Println(cur)
		//add new node to result
		result.AddNode(&hoge{name: "{imp, spc}"})
		result.AddNode(&hoge{name: fmt.Sprintf("{%s, %s}", cur.implementation, cur.specification)})

		for _, impEdges := range imp.Edges(cur.implementation) {
			impEdge, _ := impEdges.(*edge)
			impEvent := impEdge.Name()
			impNextN := impEdge.End()

			//get pairs
			pairs := []*pair{}
			if impEdge.hidden {
				pairs = append(pairs, &pair{implementation: impNextN, specification: cur.specification})
			} else {
				for _, spcEdges := range spc.Edges(cur.specification) {
					spcEdge, _ := spcEdges.(*edge)
					spcEvent := spcEdge.Name()
					spcNextN := spcEdge.End()
					if spcEvent != impEvent {
						continue
					}
					pairs = append(pairs, &pair{implementation: impNextN, specification: spcNextN})
				}
			}

			for _, p := range pairs {
				result.AddEdge(&edge{
					start: fmt.Sprintf("{%s, %s}", cur.implementation, cur.specification),
					end:   fmt.Sprintf("{%s, %s}", p.implementation, p.specification),
					name:  impEvent,
				})
				if _, ok := checked[*p]; !ok {
					checked[*p] = true
					field = append(field, p)
				}
			}

			if len(pairs) > 0 {
				continue
			}
			violated = true
			fmt.Printf("Out of Specification imp:{cur:%s -> next:%s}\n", cur.implementation, impNextN)
			result.AddEdge(
				&edge{
					start: fmt.Sprintf("{%s, %s}", cur.implementation, cur.specification),
					end:   "Out of Specification",
					name:  fmt.Sprintf("{%s: %s}", impEvent, impNextN),
					color: "red",
				})
		}
	}

	if violated {
		result.AddNode(&hoge{name: fmt.Sprintf("Out of Specification"), color: "red"})
	}

	//export result
	str, err = result.ExportDot()
	if err != nil {
		log.Fatal(err)
	}
	f2.WriteString(str)
}
