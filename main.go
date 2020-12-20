package main

import "github.com/sidav/golibrl/console"

func main() {
	console.Init_console("TaffeRL", console.SDLRenderer)
	USE_ALT_RUNES = true
	defer console.Close_console()
	game := game{}
	game.runGame()
}
