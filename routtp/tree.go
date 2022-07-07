package routtp

import (
	"fmt"
	"strings"
)

type nodeType uint8

const (
	// 普通的中间节点
	NodeNormal nodeType = 1 << iota
	// 叶子节点
	// 这里的叶子不是将它是树的叶子
	// 而是说，它是某个路由的结尾
	NodeLeaf
	// 包含通配符节点
	NodeWild
	// 包含全匹配节点
	NodeAllin
)

type Node struct {
	Path     string
	Type     nodeType
	Children []*Node
	Handlers HandlersChain `json:"-"`
	// FullPath string
}

func NewNode(path string, handlers HandlersChain) *Node {
	var typ nodeType
	for i := 0; i < len(path); i++ {
		switch path[i] {
		case ':':
			typ = NodeWild
		case '*':
			typ = NodeAllin
		}
	}
	return &Node{
		Path:     path,
		Type:     typ,
		Children: []*Node{},
		Handlers: handlers,
	}
}

func (n *Node) AddRoute(path string, handlers ...HandlerFunc) {
	if n.Path == "" {
		n.Path = path
		n.Handlers = handlers
		return
	}
	idx := longestPrefix(n.Path, path)

	if idx < len(n.Path) { // 大地的裂变
		for i := idx - 1; i >= 0 && n.Path[i] != '/'; i-- {
			assert(n.Path[i] == '*', "") // all in 节点后面不可以有子节点
			if n.Path[i] == ':' {
				assert(n.Path[idx] != '/', fmt.Sprintf("n.Path<%s> path<%s> idx<%d> i<%d>", n.Path, path, idx, i)) // 同一个路由节只能有一个通配节点
				if idx < len(path) {
					assert(path[idx] != '/', "") // 同一个路由节只能有一个通配节点
				}
			}
		}

		oldNode := NewNode(n.Path[idx:], n.Handlers)
		oldNode.Children = n.Children

		n.Path = n.Path[:idx]
		n.Handlers = []HandlerFunc{}
		n.Children = []*Node{oldNode}
		if idx < len(path) {
			newNode := NewNode(path[idx:], handlers)
			n.Children = append(n.Children, newNode)
		}
		return
	}

	n.insChild(path[idx:], "", handlers) // i == len(n.Path)
}

func (n *Node) insChild(path string, fullPath string, handlers HandlersChain) {
	for _, v := range n.Children {
		if v.Path[0] == path[0] {
			v.AddRoute(path, handlers...)
			return
		}
	}
	n.appendChild(NewNode(path, handlers))
}

func (n *Node) appendChild(newNode *Node) {
	l := len(n.Children)
	switch l {
	case 0:
		n.Children = append(n.Children, newNode)
	default:
		firstChar := n.Children[len(n.Children)-1].Path[0]
		if firstChar == ':' || firstChar == '*' { // 如果最后一个节点是通配节点
			assert(newNode.Path[0] == ':', "can not")
			assert(newNode.Path[0] == '*', "can not")
			n.Children = append(n.Children[:l-1], newNode, n.Children[l-1])
		} else {
			n.Children = append(n.Children, newNode)
		}
	}
}

func (n *Node) Get(ctx *Context) {
	uri := ctx.Request.RequestURI
	if !strings.HasPrefix(uri, n.Path) {
		panic("404 Not Found")
	}
	ctx.Fns = append(ctx.Fns, n.Handlers...)
	for _, v := range n.Children {
		if v.Path[0] == uri[0] {
			v.Get(ctx)
			goto then
		}
	}
	panic("404 Not Found")
	// 后续处理
then:
}
