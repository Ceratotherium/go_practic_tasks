//go:build task_template

package main

type TreeNode struct{}

func (n TreeNode) Value() string {
	// Возвращает значение узла
}

func (n TreeNode) AddChild(child TreeNode) {
	// Добавляет потомка
}

func (n TreeNode) DetachChild(child TreeNode) {
	// Удаляет потомка
}

func (n TreeNode) Children() []TreeNode {
	// Возвращает дочерние узлы
}

func Enumerate(root TreeNode, cb func(n TreeNode) error) error {
	return nil // Если cb возвращает ошибку - прерываем обход
}
