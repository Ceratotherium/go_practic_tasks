//go:build task_template

package main

type Drink struct {
	Name string
	Ice  bool
}

func MakeDrink() {
	var wg sync.WaitGroup
	wg.Add(2)
	go putIce(Drink{Name: "Coke"}, wg)
	go putIce(Drink{Name: "Fanta"}, wg)
	wg.Wait()
}

func putIce(d Drink, wg sync.WaitGroup) bool {
	d.Ice = true
	fmt.Println(d.Name)
	wg.Done()
	return true
}
