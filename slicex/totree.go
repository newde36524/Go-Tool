package slicex

type Node[T any] struct {
	Value    T
	Parent   map[string]*Node[T] `json:"-"`
	Children map[string]*Node[T]
}

func (n Node[T]) Val() T {
	return n.Value
}

// ToTree 数组转树返回根节点和map
func ToTree[S ~[]T, T any](s S, getIdAndParentId func(in T) (id, parentId string)) (root []*Node[T], query Query[T]) {
	var (
		cache S
		flag  bool
		fn    func(s S)
	)
	query = make(map[string]*Node[T], 0)
	fn = func(s S) {
		cache = make(S, 0)
		for _, v := range s {
			id, parantId := getIdAndParentId(v)
			node := &Node[T]{
				Value:    v,
				Parent:   make(map[string]*Node[T]),
				Children: make(map[string]*Node[T]),
			}
			query[id] = node
			if parentNode, ok := query[parantId]; ok {
				node.Parent[parantId] = parentNode
				parentNode.Children[id] = node
			} else {
				if flag {
					continue
				}
				if len(parantId) > 0 {
					cache = append(cache, v)
				} else {
					root = append(root, node)
				}
			}
		}
		if len(cache) == 0 {
			return
		}
		flag = true
		fn(cache)
	}
	fn(s)
	return
}

type Query[T any] map[string]*Node[T]

func (q Query[T]) GetNode(id string) *Node[T] {
	return q[id]
}

func (q Query[T]) GetParentFirst(id string) *Node[T] {
	for _, v := range q[id].Parent {
		return v
	}
	return nil
}

func (q Query[T]) GetParentList(id string) Query[T] {
	return q[id].Parent
}

func (q Query[T]) GetChildrenList(id string) Query[T] {
	return q[id].Children
}
