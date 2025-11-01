package main

var privateTestCases = []TestCase{
	{
		name:        "MethodNotFound",
		object:      &Calculator{},
		method:      "NonExistentMethod",
		expectedErr: true,
	},
	{
		name:        "WrongNumberOfArguments",
		object:      &Calculator{},
		method:      "Add",
		input:       []interface{}{1},
		expectedErr: true,
	},
	{
		name:        "WrongArgumentType",
		object:      &Calculator{},
		method:      "Add",
		input:       []interface{}{"string", 2},
		expectedErr: true,
	},
}
