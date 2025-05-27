//go:build task_template

package main

func SafeExecute(fn func() error) (err error) {
	// Ваш код здесь
	// Используйте defer с recover для обработки паники
}
