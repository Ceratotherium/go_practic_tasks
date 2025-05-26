//go:build task_template

package main

type Stack struct{}

func (s Stack) Push(v int) {
	// Добавление нового элемента
}

func (s Stack) Pop() int {
	// Удаление элемента из стека и возвращение из функции
}

func (s Stack) Len() int {
	// Количество элементов в стеке
}

func (s Stack) Values() []int {
	// Возвращение всех значений стека
}

func (s Stack) Max() int {
	// Возвращает наибольшее хранимое значение
}
