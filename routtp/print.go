package routtp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func genName(name string) string {
	return fmt.Sprintf("%s<%d>", name, rand.Int())
}

func genNode(graph *cgraph.Graph, v *Node) *cgraph.Node {
	childNode, err := graph.CreateNode(genName(v.Path))
	if err != nil {
		panic(err)
	}
	childNode.SetLabel(v.Path)
	return childNode
}

func PrintNode(n *Node) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			panic(err)
		}
		g.Close()
	}()

	pNode, err := graph.CreateNode(genName(n.Path))
	if err != nil {
		panic(err)
	}
	pNode.SetLabel(n.Path)

	var printNode func(root *cgraph.Node, node *Node)

	printNode = func(root *cgraph.Node, node *Node) {
		for _, v := range node.Children {
			childNode := genNode(graph, v)
			childNode.SetTooltip(node.Path + v.Path)
			graph.CreateEdge("", root, childNode)
			printNode(childNode, v)
		}
	}

	printNode(pNode, n)

	var buf bytes.Buffer
	err = g.Render(graph, graphviz.XDOT, &buf)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("d.dot", buf.Bytes(), 0644)
}
