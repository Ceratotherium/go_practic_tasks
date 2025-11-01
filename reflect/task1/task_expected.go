package main

import (
	"fmt"
	"reflect"
)

func CallMethod(obj interface{}, methodName string, args ...interface{}) ([]interface{}, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		// Разыменовываем указатель, если он указывает на структуру
		val = val.Elem()
	}

	// Проверяем, что переданный объект - структура
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct or pointer to struct, got %T", obj)
	}

	// Получаем метод по имени
	method := val.MethodByName(methodName)
	if !method.IsValid() {
		// Если метод не найден у значения, проверяем у указателя
		ptrVal := reflect.ValueOf(obj)
		if ptrVal.Kind() == reflect.Ptr {
			method = ptrVal.MethodByName(methodName)
		}

		if !method.IsValid() {
			return nil, fmt.Errorf("method '%s' not found", methodName)
		}
	}

	// Проверяем количество аргументов
	methodType := method.Type()
	if methodType.NumIn() != len(args) {
		return nil, fmt.Errorf("method expects %d arguments, got %d",
			methodType.NumIn(), len(args))
	}

	// Подготавливаем аргументы
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		argType := methodType.In(i)
		argValue := reflect.ValueOf(arg)

		// Проверяем тип аргумента
		if !argValue.Type().ConvertibleTo(argType) {
			return nil, fmt.Errorf("argument %d has type %T, but method expects %v",
				i, arg, argType)
		}

		// Конвертируем тип, если необходимо
		in[i] = argValue.Convert(argType)
	}

	// Вызываем метод
	results := method.Call(in)

	// Конвертируем результаты в interface{}
	out := make([]interface{}, len(results))
	for i, result := range results {
		out[i] = result.Interface()
	}

	return out, nil
}
