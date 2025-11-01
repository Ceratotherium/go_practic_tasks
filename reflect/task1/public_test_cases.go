package main

import (
	"fmt"
)

type TestCase struct {
	name        string
	object      interface{}
	method      string
	input       []interface{}
	expected    []interface{}
	expectedErr bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:     "MethodWithPointerReceiver",
		object:   &Calculator{offset: 10},
		method:   "Add",
		input:    []interface{}{5, 3},
		expected: []interface{}{18},
	},
	{
		name:     "MethodWithStringArgsAndReturn",
		object:   &Calculator{offset: 42},
		method:   "Greet",
		input:    []interface{}{"Alice"},
		expected: []interface{}{"Hello, Alice! Offset is 42"},
	},
	{
		name:     "MethodWithValueReceiver",
		object:   Calculator{offset: 5},
		method:   "Multiply",
		input:    []interface{}{4, 3},
		expected: []interface{}{12},
	},
	// Тесткейсы в помощь
	{
		name:     "MethodWithValueReceiverAndPointerToObject",
		object:   &Calculator{offset: 5},
		method:   "Multiply",
		input:    []interface{}{4, 3},
		expected: []interface{}{12},
	},
	{
		name:     "MethodWithNoArgs",
		object:   &Calculator{},
		method:   "NoArgs",
		expected: []interface{}{"No arguments method"},
	},
	{
		name:        "WithNilObject",
		object:      (*Calculator)(nil),
		method:      "Add",
		input:       []interface{}{1, 2},
		expectedErr: true,
	},
}

type Calculator struct {
	offset int
}

func (c *Calculator) Add(a, b int) int {
	return a + b + c.offset
}

func (c Calculator) Multiply(a, b int) int {
	return a * b
}

func (c *Calculator) Greet(name string) string {
	return fmt.Sprintf("Hello, %s! Offset is %d", name, c.offset)
}

func (c Calculator) NoArgs() string {
	return "No arguments method"
}
