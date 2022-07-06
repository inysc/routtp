package routtp

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
	i := longestPrefix(n.Path, path)

	if i < len(n.Path) { // 大地的裂变
		oldNode := NewNode(n.Path[i:], n.Handlers)
		oldNode.Children = n.Children

		n.Path = n.Path[:i]
		n.Handlers = []HandlerFunc{}
		n.Children = []*Node{oldNode}
		if i < len(path) {
			newNode := NewNode(path[i:], handlers)
			n.Children = append(n.Children, newNode)
		}
		return
	}
	n.InsChild(path[i:], "", handlers) // i == len(n.Path)
}

func (n *Node) InsChild(path string, fullPath string, handlers HandlersChain) {
	for _, v := range n.Children {
		if v.Path[0] == path[0] {
			v.AddRoute(path, handlers...)
			return
		}
	}
	n.Children = append(n.Children, NewNode(path, handlers))
}
