package main

import "github.com/sidav/golibrl/console"

func main() {
	console.Init_console("SQU@D", console.TCellRenderer)
	defer console.Close_console()
	game := game{}
	game.runGame()
}
