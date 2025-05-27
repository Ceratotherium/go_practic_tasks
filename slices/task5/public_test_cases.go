package main

import "errors"

type TestCase struct {
	name    string
	prepare func() *TreeNode
	check   func(root *TreeNode) bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Дерево из одного узла",
		prepare: func() *TreeNode {
			return NewNode("root")
		},
		check: func(root *TreeNode) bool {
			return root.Value() == "root" && len(root.Children()) == 0
		},
	},
	{
		name: "Дерево с несколькими дочерними узлами",
		prepare: func() *TreeNode {
			root := NewNode("root")
			child1 := NewNode("child1")
			child2 := NewNode("child2")
			root.AddChild(child1)
			root.AddChild(child2)
			return root
		},
		check: func(root *TreeNode) bool {
			if root.Value() != "root" || len(root.Children()) != 2 {
				return false
			}
			return root.Children()[0].Value() == "child1" &&
				root.Children()[1].Value() == "child2"
		},
	},
	{
		name: "Удаление узла",
		prepare: func() *TreeNode {
			root := NewNode("root")
			child1 := NewNode("child1")
			child2 := NewNode("child2")
			root.AddChild(child1)
			root.AddChild(child2)
			root.DetachChild(child1)
			return root
		},
		check: func(root *TreeNode) bool {
			return len(root.Children()) == 1 &&
				root.Children()[0].Value() == "child2"
		},
	},
	{
		name: "Обход всего дерева",
		prepare: func() *TreeNode {
			root := NewNode("root")
			child1 := NewNode("child1")
			child2 := NewNode("child2")
			grandchild := NewNode("grandchild")

			child1.AddChild(grandchild)
			root.AddChild(child1)
			root.AddChild(child2)
			return root
		},
		check: func(root *TreeNode) bool {
			var nodes []string
			err := Enumerate(root, func(n *TreeNode) error {
				nodes = append(nodes, n.Value())
				return nil
			})

			if err != nil {
				return false
			}

			expected := []string{"root", "child1", "grandchild", "child2"}
			if len(nodes) != len(expected) {
				return false
			}

			for i, val := range expected {
				if nodes[i] != val {
					return false
				}
			}
			return true
		},
	},
	// Тесткейсы в помощь
	{
		name: "Дерево с тремя уровнями",
		prepare: func() *TreeNode {
			root := NewNode("root")
			child1 := NewNode("child1")
			child2 := NewNode("child2")
			grandchild := NewNode("grandchild")

			child1.AddChild(grandchild)
			root.AddChild(child1)
			root.AddChild(child2)
			return root
		},
		check: func(root *TreeNode) bool {
			if len(root.Children()) != 2 {
				return false
			}
			child1 := root.Children()[0]
			if len(child1.Children()) != 1 {
				return false
			}
			return child1.Children()[0].Value() == "grandchild"
		},
	},
	{
		name: "Удаление среднего узла",
		prepare: func() *TreeNode {
			root := NewNode("root")
			child1 := NewNode("child1")
			child2 := NewNode("child2")
			child3 := NewNode("child3")
			root.AddChild(child1)
			root.AddChild(child2)
			root.AddChild(child3)
			root.DetachChild(child2)
			return root
		},
		check: func(root *TreeNode) bool {
			return len(root.Children()) == 2 &&
				root.Children()[0].Value() == "child1" &&
				root.Children()[1].Value() == "child3"
		},
	},
	{
		name: "Удаление правого узла",
		prepare: func() *TreeNode {
			root := NewNode("root")
			child1 := NewNode("child1")
			child2 := NewNode("child2")
			child3 := NewNode("child3")
			root.AddChild(child1)
			root.AddChild(child2)
			root.AddChild(child3)
			root.DetachChild(child3)
			return root
		},
		check: func(root *TreeNode) bool {
			return len(root.Children()) == 2 &&
				root.Children()[0].Value() == "child1" &&
				root.Children()[1].Value() == "child2"
		},
	},
	{
		name: "Прерывание обхода",
		prepare: func() *TreeNode {
			root := NewNode("root")
			child1 := NewNode("child1")
			child2 := NewNode("child2")
			root.AddChild(child1)
			root.AddChild(child2)
			return root
		},
		check: func(root *TreeNode) bool {
			testErr := errors.New("test error")
			err := Enumerate(root, func(n *TreeNode) error {
				if n.Value() == "child1" {
					return testErr
				}
				return nil
			})
			return errors.Is(err, testErr)
		},
	},
}
