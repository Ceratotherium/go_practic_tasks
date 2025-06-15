package main

import (
	"fmt"
	"sync"
	"time"
)

var privateTestCases = []TestCase{
	{
		name: "Остановка фоновой очистки",
		prepare: func() (Cache, func()) {
			c, stop := New(time.Millisecond)
			c.Set("test", "value", time.Millisecond)

			stop() // останавливаем очистку

			time.Sleep(2 * time.Millisecond)
			return c, func() {}
		},
		check: func(c Cache) bool {
			_, ok := c.Get("test")
			return !ok // Очистки не должно было быть, но элемент все равно недоступен
		},
	},
	{
		name: "Независимость двух кешей",
		prepare: func() (Cache, func()) {
			c1, stop1 := New(time.Minute)
			c1.Set("test", "value1", time.Minute)

			c2, stop2 := New(time.Minute)
			defer stop2()
			c2.Set("test", "value2", time.Minute)

			return c1, stop1
		},
		check: func(c Cache) bool {
			value, ok := c.Get("test")
			return ok && value.(string) == "value1"
		},
	},
	{
		name: "Массовое добавление элементов (10 000)",
		prepare: func() (Cache, func()) {
			c, stop := New(0)
			for i := 0; i < 10000; i++ {
				key := fmt.Sprintf("key%d", i)
				c.Set(key, i, time.Minute)
			}
			return c, stop
		},
		check: func(c Cache) bool {
			for i := 0; i < 10000; i++ {
				key := fmt.Sprintf("key%d", i)
				val, ok := c.Get(key)
				if !ok || val != i {
					return false
				}
			}
			return true
		},
	},
	{
		name: "Параллельное массовое добавление (1 000 элементов)",
		prepare: func() (Cache, func()) {
			c, stop := New(0)
			var wg sync.WaitGroup

			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					key := fmt.Sprintf("key%d", idx)
					c.Set(key, idx, time.Minute)
				}(i)
			}

			wg.Wait()
			return c, stop
		},
		check: func(c Cache) bool {
			for i := 0; i < 1000; i++ {
				key := fmt.Sprintf("key%d", i)
				val, ok := c.Get(key)
				if !ok || val != i {
					return false
				}
			}
			return true
		},
	},
	{
		name: "Смешанная нагрузка (чтение/запись 5 000 элементов)",
		prepare: func() (Cache, func()) {
			c, stop := New(time.Minute)
			var wg sync.WaitGroup

			// Заполняем начальные данные
			for i := 0; i < 5000; i++ {
				key := fmt.Sprintf("key%d", i)
				c.Set(key, i, time.Minute)
			}

			// Параллельные чтения и записи
			for i := 0; i < 5000; i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					key := fmt.Sprintf("key%d", idx)

					// Читаем
					val, ok := c.Get(key)
					if !ok || val != idx {
						return
					}

					// Обновляем
					c.Set(key, idx+1000, time.Minute)
				}(i)
			}

			wg.Wait()
			return c, stop
		},
		check: func(c Cache) bool {
			for i := 0; i < 5000; i++ {
				key := fmt.Sprintf("key%d", i)
				val, ok := c.Get(key)
				if !ok || val != i+1000 {
					return false
				}
			}
			return true
		},
	},
	{
		name: "Долгоживущие и краткоживущие элементы (10 000)",
		prepare: func() (Cache, func()) {
			c, stop := New(time.Second)

			for i := 0; i < 10000; i++ {
				key := fmt.Sprintf("key%d", i)
				if i%2 == 0 {
					c.Set(key, i, time.Minute) // Долгоживущие
				} else {
					c.Set(key, i, time.Millisecond) // Краткоживущие
				}
			}

			time.Sleep(2 * time.Millisecond) // Ждем истечения краткоживущих
			return c, stop
		},
		check: func(c Cache) bool {
			count := 0
			for i := 0; i < 10000; i++ {
				key := fmt.Sprintf("key%d", i)
				val, ok := c.Get(key)

				if i%2 == 0 {
					// Должны существовать долгоживущие
					if !ok || val != i {
						return false
					}
					count++
				} else {
					// Краткоживущие должны исчезнуть
					if ok {
						return false
					}
				}
			}
			return count == 5000
		},
	},
	{
		name: "Тест на утечку памяти (100 000 элементов)",
		prepare: func() (Cache, func()) {
			c, stop := New(time.Millisecond)

			for i := 0; i < 100000; i++ {
				key := fmt.Sprintf("key%d", i)
				c.Set(key, i, time.Microsecond) // Очень короткое время жизни
			}

			time.Sleep(10 * time.Millisecond) // Ждем очистки
			return c, stop
		},
		check: func(c Cache) bool {
			count := 0
			for i := 0; i < 100000; i++ {
				key := fmt.Sprintf("key%d", i)
				if _, ok := c.Get(key); ok {
					count++
				}
			}
			return count == 0 // Все элементы должны быть удалены
		},
	},
}
