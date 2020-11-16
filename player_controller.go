package main

import (
	"github.com/sidav/golibrl/console"
)

type playerController struct {
}

func (p *playerController) playerControl(d *gameMap) {
	// p := d.player
	if p.checkGameState() {
		return
	}

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
			case "n":
				newNoise := CURRENT_MAP.player.createMovementNoise()
				newNoise.intensity = 15
				newNoise.suspicious = true
				CURRENT_MAP.createNoise(newNoise)
				log.AppendMessage("*Whistle*")
				CURRENT_MAP.player.spendTurnsForAction(10)
			case "c":
				pc.doCloseDoor()
			case "ESCAPE":
				GAME_IS_RUNNING = false
			case "INSERT":
				RENDER_DISABLE_LOS = !RENDER_DISABLE_LOS
			case "HOME":
				CURRENT_MAP.player.inv.gold += 111
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

func (pc *playerController) doCloseDoor() {
	px, py := CURRENT_MAP.player.getCoords()
	doorsAround := CURRENT_MAP.getNumberOfTilesOfTypeAround(TILE_DOOR, px, py)
	if doorsAround == 1 {
		for x := px-1; x <= px+1; x++ {
			for y := py-1; y <= py+1; y++{
				if CURRENT_MAP.isTileADoor(x, y) && CURRENT_MAP.tiles[x][y].isOpened {
					CURRENT_MAP.tiles[x][y].isOpened = false
					CURRENT_MAP.player.spendTurnsForAction(10)
				}
			}
		}
	} else if doorsAround > 1 {
		log.AppendMessage("Which direction?")
		renderLog(true)
		dirx, diry := pc.keyToDirection(console.ReadKey())
		if dirx != 0 || diry != 0 {
			if CURRENT_MAP.isTileADoor(px+dirx, py+diry) {
				CURRENT_MAP.tiles[px+dirx][py+diry].isOpened = false
				CURRENT_MAP.player.spendTurnsForAction(10)
			}
		}
	}
}

func (p *playerController) checkGameState() bool {
	plr := CURRENT_MAP.player
	if plr.hp <= 0 {
		GAME_IS_RUNNING = false
		gameover()
		return true
	}
	w, h := CURRENT_MAP.getSize()
	if plr.x == 0 || plr.x == w-1 || plr.y == 0 || plr.y == h-1 {
		if plr.inv.gold >= 1000 {
			GAME_IS_RUNNING = false
			gamewon()
			return true
		} else {
			log.AppendMessage("You need to collect at least 1000 gold before exfiltration!")
		}
	}
	return false
}
