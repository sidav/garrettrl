package main

import (
	"github.com/sidav/golibrl/console"
	"time"
)

type playerController struct {
	currentSelectedArrowIndex int
	previousHp                int
}

func (pc *playerController) playerControl(d *gameMap) {
	p := d.player
	if pc.previousHp > p.hp {
		for i := pc.previousHp; i >= p.hp; i-- {
			renderer.renderDamageFlash()
			time.Sleep(25 * time.Millisecond)
		}
	}
	pc.previousHp = p.hp
	renderer.renderGameScreen(true)
	if pc.checkGameStop() {
		return
	}

	valid_key_pressed := false
	movex := 0
	movey := 0
	for !valid_key_pressed {
		key_pressed := console.ReadKeyAsync()
		valid_key_pressed = true
		movex, movey = pc.keyToDirection(key_pressed)
		if movex == 0 && movey == 0 {
			switch key_pressed {
			case "5":
				p.spendTurnsForAction(10)
			case "r":
				p.isRunning = !p.isRunning
			case "n":
				newNoise := p.createMovementNoise()
				newNoise.intensity = 15
				newNoise.suspicious = true
				CURRENT_MAP.createNoise(newNoise)
				log.AppendMessage("*Whistle*")
				p.spendTurnsForAction(10)
			case "a": // select next arrow 
				initialArrowIndex := pc.currentSelectedArrowIndex
				pc.currentSelectedArrowIndex++
				for {
					if pc.currentSelectedArrowIndex >= len(p.inv.arrows) {
						pc.currentSelectedArrowIndex = 0
					}
					if pc.currentSelectedArrowIndex == initialArrowIndex {
						break
					}
					if p.inv.arrows[pc.currentSelectedArrowIndex].amount == 0 {
						pc.currentSelectedArrowIndex++
					} else {
						break
					}
				}
				break
			case "c":
				pc.doCloseDoor()
			case "f": // fire arrow
				log.AppendMessage("Select a target.")
				if p.inv.arrows[pc.currentSelectedArrowIndex].amount > 0 {
					sx, sy := pc.selectCoords(true)
					if sx == CURRENT_MAP.player.x && sy == CURRENT_MAP.player.y {
						log.AppendMessage("Trying to suicide with " + p.inv.arrows[pc.currentSelectedArrowIndex].name + ", huh?")
					} else if sx != -1 && sy != -1 {
						applyArrowEffect(p.inv.arrows[pc.currentSelectedArrowIndex].name, sx, sy)
						p.inv.arrows[pc.currentSelectedArrowIndex].amount--
						p.spendTurnsForAction(10)
					}
				} else {
					log.AppendMessage("No such arrow in the quiver.")
				}
			case "ESCAPE":
				GAME_IS_RUNNING = false
			case "INSERT":
				renderer.RENDER_DISABLE_LOS = !renderer.RENDER_DISABLE_LOS
			case "HOME":
				p.inv.gold += rnd.Rand(44)
				for a := range p.inv.arrows {
					p.inv.arrows[a].amount++
				}
			case "DELETE":
				p.hp += rnd.Rand(5)
			case "NOTHING", "NON-KEY":
				time.Sleep(25 * time.Millisecond)
				valid_key_pressed = false
			default:
				valid_key_pressed = false
				log.AppendMessagef("Unknown key %s (Wrong keyboard layout?)", key_pressed)
				renderer.renderGameScreen(true)
			}
		}
	}
	// move player's pawn here and something
	if movex != 0 || movey != 0 {
		CURRENT_MAP.defaultMovementActionByVector(p, true, movex, movey)
	}
}

func (p *playerController) keyToDirection(keyPressed string) (int, int) {
	switch keyPressed {
	case "2":
		return 0, 1
	case "8":
		return 0, -1
	case "4":
		return -1, 0
	case "6":
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

func (pc *playerController) selectCoords(forceVisible bool) (int, int) {
	sx, sy := CURRENT_MAP.player.getCoords()
	for {
		renderer.renderGameScreen(false)
		renderer.renderCursor(sx, sy, true)
		key := console.ReadKeyAsync()
		if key == "ENTER" || key == "f" {
			if !forceVisible || CURRENT_MAP.currentPlayerVisibilityMap[sx][sy] {
				break
			} else {
				log.AppendMessage("Select visible tile!")
			}
		}
		if key == "ESCAPE" {
			log.AppendMessage("Fine, then.")
			return -1, -1
		}
		dx, dy := pc.keyToDirection(key)
		sx += dx
		sy += dy
	}
	return sx, sy
}

func (pc *playerController) doCloseDoor() {
	px, py := CURRENT_MAP.player.getCoords()
	doorsAround := CURRENT_MAP.getNumberOfTilesOfTypeAround(TILE_DOOR, px, py)
	if doorsAround == 1 {
		for x := px - 1; x <= px+1; x++ {
			for y := py - 1; y <= py+1; y++ {
				if CURRENT_MAP.isTileADoor(x, y) && CURRENT_MAP.tiles[x][y].isOpened {
					CURRENT_MAP.tiles[x][y].isOpened = false
					CURRENT_MAP.player.spendTurnsForAction(10)
				}
			}
		}
	} else if doorsAround > 1 {
		log.AppendMessage("Which direction?")
		renderer.renderLog(true)
		dirx, diry := pc.keyToDirection(console.ReadKeyAsync())
		if dirx != 0 || diry != 0 {
			if CURRENT_MAP.isTileADoor(px+dirx, py+diry) {
				CURRENT_MAP.tiles[px+dirx][py+diry].isOpened = false
				CURRENT_MAP.player.spendTurnsForAction(10)
			}
		}
	}
}

func (p *playerController) checkGameStop() bool {
	plr := CURRENT_MAP.player
	if plr.hp <= 0 {
		GAME_IS_RUNNING = false
		gameover()
		return true
	}
	if p.checkWinning() {
		GAME_IS_RUNNING = false
		gamewon()
		return true
	}
	return false
}

func (pc *playerController) checkWinning() bool {
	won := false
	plr := CURRENT_MAP.player
	w, h := CURRENT_MAP.getSize()
	isPlayerOnEdge := plr.x == 0 || plr.x == w-1 || plr.y == 0 || plr.y == h-1
	switch currMission.MissionType {
	case MISSION_STEAL_MINIMUM_LOOT:
		won = isPlayerOnEdge && plr.inv.gold >= currMission.TargetNumber[currDifficultyNumber]
	case MISSION_STEAL_TARGET_ITEMS:
		won = isPlayerOnEdge && len(plr.inv.targetItems) >= currMission.TargetNumber[currDifficultyNumber]
	}
	return won
}
