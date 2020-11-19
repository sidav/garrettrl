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
		{name: "Water arrow", amount: 1},
		{name: "Noise arrow", amount: 1},
		{name: "Gas arrow", amount: 0},
		{name: "Explosive arrow", amount: 0},
	}
}

func (i *inventory) grabEverythingFromInventory(i2 *inventory) {
	i.gold += i2.gold
	for t := range i2.arrows {
		i.arrows[t].amount += i2.arrows[t].amount
	}
}
