package main

var privateTestCases = []TestCase{
	{
		name: "Удаление несуществующего узла",
		prepare: func() *TreeNode {
			root := NewNode("root")
			child1 := NewNode("child1")
			child2 := NewNode("child2")
			root.AddChild(child1)
			root.AddChild(child2)
			return root
		},
		check: func(root *TreeNode) bool {
			allNodes := CollectNodes(root)

			absentChild := NewNode("absent")
			root.DetachChild(absentChild)

			currentNodes := CollectNodes(root)

			return compareSliceValues(allNodes, currentNodes)
		},
	},
	{
		name: "Добавление nil узла",
		prepare: func() *TreeNode {
			root := NewNode("root")
			return root
		},
		check: func(root *TreeNode) bool {
			return !AssertPanic(func() {
				root.AddChild(nil)
			})
		},
	},
	{
		name: "Удаление nil узла",
		prepare: func() *TreeNode {
			root := NewNode("root")
			return root
		},
		check: func(root *TreeNode) bool {
			return !AssertPanic(func() {
				root.DetachChild(nil)
			})
		},
	},
	{
		name: "Обход nil дерева",
		prepare: func() *TreeNode {
			return nil
		},
		check: func(root *TreeNode) bool {
			called := false
			err := Enumerate(root, func(n *TreeNode) error {
				called = true
				return nil
			})
			return err == nil && !called
		},
	},
	{
		name: "Добавление поддерева",
		prepare: func() *TreeNode {
			root := NewNode("root")
			return root
		},
		check: func(root *TreeNode) bool {
			child1 := NewNode("child1")
			grandChild1 := NewNode("grandChild1")
			grandChild2 := NewNode("grandChild2")
			child1.AddChild(grandChild1)
			child1.AddChild(grandChild2)
			root.AddChild(child1)

			return compareSliceValues(
				[]string{"root", "child1", "grandChild1", "grandChild2"},
				CollectNodes(root),
			)
		},
	},
	{
		name: "Удаление поддерева",
		prepare: func() *TreeNode {
			root := NewNode("root")
			child1 := NewNode("child1")
			grandChild1 := NewNode("grandChild1")
			grandChild2 := NewNode("grandChild2")
			child1.AddChild(grandChild1)
			child1.AddChild(grandChild2)
			root.AddChild(child1)

			child2 := NewNode("child2")
			grandChild3 := NewNode("grandChild3")
			grandChild4 := NewNode("grandChild4")
			child2.AddChild(grandChild3)
			child2.AddChild(grandChild4)
			root.AddChild(child2)
			return root
		},
		check: func(root *TreeNode) bool {
			root.DetachChild(root.Children()[0])

			return compareSliceValues(
				[]string{"root", "child2", "grandChild3", "grandChild4"},
				CollectNodes(root),
			)
		},
	},
}

func CollectNodes(root *TreeNode) []string {
	result := make([]string, 0)

	_ = Enumerate(root, func(n *TreeNode) error {
		result = append(result, n.Value())
		return nil
	})

	return result
}
