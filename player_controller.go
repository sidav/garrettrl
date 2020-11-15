package main

import (
	"github.com/sidav/golibrl/console"
)

type playerController struct {
}

func (p *playerController) playerControl(d *gameMap) {
	// p := d.player
	valid_key_pressed := false
	movex := 0
	movey := 0
	for !valid_key_pressed {
		key_pressed := console.ReadKey()
		valid_key_pressed = true
		movex, movey = p.keyToDirection(key_pressed)
		if movex == 0 && movey == 0 {
			switch key_pressed {
			case "5":
				CURRENT_MAP.player.spendTurnsForAction(10)
			case "r":
				CURRENT_MAP.player.isRunning = !CURRENT_MAP.player.isRunning
			case "ESCAPE":
				GAME_IS_RUNNING = false
			case "INSERT":
				RENDER_DISABLE_LOS = !RENDER_DISABLE_LOS
			default:
				valid_key_pressed = false
				log.AppendMessagef("Unknown key %s (Wrong keyboard layout?)", key_pressed)
				renderLevel(d, true)
			}
		}
	}
	// move player's pawn here and something
	if movex != 0 || movey != 0 {
		CURRENT_MAP.movePawnOrOpenDoorByVector(CURRENT_MAP.player, true, movex, movey)
	}
}

func (p *playerController) keyToDirection(keyPressed string) (int, int) {
	switch keyPressed {
	case "s", "2":
		return 0, 1
	case "w", "8":
		return 0, -1
	case "a", "4":
		return -1, 0
	case "d", "6":
		return 1, 0
	case "7":
		return -1, -1
	case "9":
		return 1, -1
	case "1":
		return -1, 1
	case "3":
		return 1, 1
	default:
		return 0, 0
	}
}
