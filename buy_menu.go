package main

import cw "github.com/sidav/golibrl/console"

type buyMenu struct {
	currentGold    int
	cursorPosition int
	itemsNames     []string
	itemsCosts     []int
	itemsBought    []int
}

func initBuyMenu(i *inventory) *buyMenu {
	b := buyMenu{}
	b.currentGold = i.gold
	b.itemsNames = make([]string, 4)
	b.itemsCosts = make([]int, 4)
	b.itemsBought = make([]int, 4)
	b.itemsNames[0] = "Water arrow"
	b.itemsCosts[0] = 120
	b.itemsNames[1] = "Noise arrow"
	b.itemsCosts[1] = 90
	b.itemsNames[2] = "Gas arrow"
	b.itemsCosts[2] = 250
	b.itemsNames[3] = "Explosive arrow"
	b.itemsCosts[3] = 800
	return &b
}

func (b *buyMenu) accessBuyMenu(pInv *inventory) {
	menuActive := true
	for menuActive {
		renderer.renderBuyMenu(b)
		key := cw.ReadKeyAsync()
		switch key {
		case "ENTER", "ESCAPE": menuActive = false
		case "DOWN": b.cursorPosition++
		case "UP": b.cursorPosition--
		case "RIGHT":
			if b.currentGold > b.itemsCosts[b.cursorPosition] {
				b.itemsBought[b.cursorPosition] ++
				b.currentGold -= b.itemsCosts[b.cursorPosition]
			}
		case "LEFT":
			if b.itemsBought[b.cursorPosition] > 0 {
				b.itemsBought[b.cursorPosition]--
				b.currentGold += b.itemsCosts[b.cursorPosition]
			}
		}
		if b.cursorPosition < 0 {
			b.cursorPosition = len(b.itemsNames)-1
		}
		if b.cursorPosition == len(b.itemsNames) {
			b.cursorPosition = 0
		}
	}
	pInv.gold = 0
	boughtInv := inventory{
		gold:   b.currentGold,
	}
	for i := range b.itemsBought {
		boughtInv.arrows = append(boughtInv.arrows, arrow{
			name:   b.itemsNames[i],
			amount: b.itemsBought[i],
		})
	}
	pInv.grabEverythingFromInventory(&boughtInv)
}
