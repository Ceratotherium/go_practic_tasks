#!/bin/bash

read -p "Package name: " new_package
read -p "Task name: " task_name
read -p "Testing function: " test_func
read -p "Testing function sign: " test_func_sign

mkdir "./${new_package}";
mkdir "./${new_package}/${task_name}" && \
cp ./testing/* ./${new_package}/${task_name}/ && \
cp ./run_tools/* ./${new_package}/${task_name}/ && \
cp ./templates/* ./${new_package}/${task_name}/ && \
chmod +x ./${new_package}/${task_name}/compile.sh && \
chmod +x ./${new_package}/${task_name}/compile.sh

echo "package main

 func main() {
 	tests := append(testCases, privateTestCases...)

 	for _, tt := range tests {
 		AssertEqual(tt.name, tt.expected, ${test_func}, tt.input)
 	}
 }
" > "./${new_package}/${task_name}/main.go"


echo "//go:build task_template

 package main

 func ${test_func}${test_func_sign} {
 	return
 }
" > "./${new_package}/${task_name}/task.go"


echo "package main

func ${test_func}${test_func_sign} {
return
}
" > "./${new_package}/${task_name}/task_expected.go"