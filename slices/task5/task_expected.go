package main

type TreeNode struct {
	value    string
	children []*TreeNode
}

func (n *TreeNode) Value() string {
	return n.value
}

func (n *TreeNode) AddChild(child *TreeNode) {
	n.children = append(n.children, child)
}

func (n *TreeNode) DetachChild(child *TreeNode) {
	for i, c := range n.children {
		if c == child {
			n.children = append(n.children[:i], n.children[i+1:]...)
		}
	}
}

func (n *TreeNode) Children() []*TreeNode {
	return n.children
}

func NewNode(value string) *TreeNode {
	return &TreeNode{
		children: make([]*TreeNode, 0),
		value:    value,
	}
}

func Enumerate(root *TreeNode, cb func(n *TreeNode) error) error {
	if root == nil {
		return nil
	}

	if err := cb(root); err != nil {
		return err
	}

	for _, child := range root.Children() {
		if err := Enumerate(child, cb); err != nil {
			return err
		}
	}

	return nil
}
