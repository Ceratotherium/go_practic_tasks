//go:build task_template

package main

type Set struct{}

func NewSet(length int) *Set {
	// Создаем новый set с заданной капасити
}

func (s *Set) Add(value ...int) {
	// Добавляем новые значения
}

func (s *Set) Remove(value int) {
	// Удаляем значение
}

func (s *Set) Contains(value int) bool {
	// Проверяем наличие элемента
}

func (s *Set) Len() int {
	// Количество элементов
}

func (s *Set) ToSlice() []int {
	// Получить список элементов (порядок не важен)
}
