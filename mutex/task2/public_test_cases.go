package main

import (
	"sync"
	"time"
)

type TestCase struct {
	name    string
	prepare func() (Cache, func())
	check   func(Cache) bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Добавление и получение элемента",
		prepare: func() (Cache, func()) {
			c, stop := New(0)
			c.Set("key1", "value1", time.Minute)
			return c, stop
		},
		check: func(c Cache) bool {
			val, ok := c.Get("key1")
			return ok && val == "value1"
		},
	},
	{
		name: "Получение несуществующего элемента",
		prepare: func() (Cache, func()) {
			return New(0)
		},
		check: func(c Cache) bool {
			val, ok := c.Get("nonexistent")
			return !ok && val == nil
		},
	},
	{
		name: "Элемент с истекшим TTL",
		prepare: func() (Cache, func()) {
			c, stop := New(0)
			c.Set("expired", "value", time.Millisecond)
			time.Sleep(2 * time.Millisecond)
			return c, stop
		},
		check: func(c Cache) bool {
			val, ok := c.Get("expired")
			return !ok && val == nil
		},
	},
	// Тесткейсы в помощь
	{
		name: "Автоматическая очистка просроченных элементов",
		prepare: func() (Cache, func()) {
			c, stop := New(time.Millisecond)
			c.Set("expired1", "value1", time.Millisecond)
			c.Set("valid1", "value2", time.Minute)
			time.Sleep(2 * time.Millisecond)
			return c, stop
		},
		check: func(c Cache) bool {
			_, okExpired := c.Get("expired1")
			val, okValid := c.Get("valid1")
			return !okExpired && okValid && val == "value2"
		},
	},
	{
		name: "Параллельные операции записи",
		prepare: func() (Cache, func()) {
			c, stop := New(0)
			var wg sync.WaitGroup

			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					c.Set(string(rune('a'+idx)), idx, time.Minute)
				}(i)
			}

			wg.Wait()
			return c, stop
		},
		check: func(c Cache) bool {
			for i := 0; i < 100; i++ {
				val, ok := c.Get(string(rune('a' + i)))
				if !ok || val != i {
					return false
				}
			}
			return true
		},
	},
	{
		name: "Параллельные операции чтения и записи",
		prepare: func() (Cache, func()) {
			c, stop := New(0)
			c.Set("shared", 0, time.Minute)

			var wg sync.WaitGroup
			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					val, _ := c.Get("shared")
					c.Set("shared", val.(int)+1, time.Minute)
				}()
			}

			wg.Wait()
			return c, stop
		},
		check: func(c Cache) bool {
			val, ok := c.Get("shared")
			typedValue := val.(int)
			return ok && typedValue > 0 && typedValue <= 100
		},
	},
}
