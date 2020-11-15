package main

type inventory struct {
	gold   int
	arrows []arrow
}

type arrow struct {
	name   string
	amount int
}

func (i *inventory) init() {
	i.gold = 0
	i.arrows = []arrow{
		{name: "Water arrow", amount: 3},
		{name: "Gas arrow", amount: 2},
		{name: "Explosive arrow", amount: 1},
	}
}
