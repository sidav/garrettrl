package main

import (
	cw "github.com/sidav/golibrl/console"
	"github.com/sidav/golibrl/random/additive_random"
	log2 "gorltemplate/game_log"
)

var (
	levelsizex, levelsizey int // TODO: remove as redundant, use dung.getSize() instead
)

var (
	GAME_IS_RUNNING bool
	log             log2.GameLog
	rnd 			additive_random.FibRandom
	pc 				playerController
	CURRENT_TURN    int
	CURRENT_MAP 	gameMap
)

type game struct {
}

func areCoordinatesValid(x, y int) bool {
	return x >= 0 && y >= 0 && x < levelsizex && y < levelsizey
}

func areCoordinatesInRangeFrom(fx, fy, tx, ty, srange int) bool {
	return (tx-fx)*(tx-fx) + (ty-fy)*(ty-fy) < srange * srange 
}

func (g *game) runGame() {
	log = log2.GameLog{}
	log.Init(5)
	rnd = additive_random.FibRandom{}
	rnd.InitDefault()
	updateBoundsIfNeccessary(true)

	GAME_IS_RUNNING = true
	CURRENT_MAP = gameMap{}
	CURRENT_MAP.generateAndInitMap() // applyRuneMap(&testMap)

	for GAME_IS_RUNNING {
		g.mainLoop()
	}
}

func (g *game) mainLoop() {
	CURRENT_MAP.recalculateLights()
	CURRENT_MAP.currentPlayerVisibilityMap = *CURRENT_MAP.getFieldOfVisionFor(CURRENT_MAP.player)

	renderLevel(&CURRENT_MAP, true)

	if CURRENT_MAP.player.isTimeToAct() {
		pc.playerControl(&CURRENT_MAP)
	}

	// check if pawns should be removed
	for i := 0; i < len(CURRENT_MAP.pawns); i++ {
		if CURRENT_MAP.pawns[i].isTimeToAct() {
			// ai_act for pawns here
			if CURRENT_MAP.pawns[i].ai != nil {
				CURRENT_MAP.pawns[i].ai_checkSituation()
				CURRENT_MAP.pawns[i].ai_act()
			}
		}
	}
	CURRENT_MAP.cleanupNoises()
	CURRENT_TURN++
}

func gameover() {
	cw.Clear_console()
	cw.PutString("You are dead! Press ENTER to exit.", 0, 0)
	cw.Flush_console()
	for cw.ReadKey() != "ENTER" {

	}
}

func gamewon() {
	cw.Clear_console()
	cw.PutString("You won! Press ENTER to exit.", 0, 0)
	cw.Flush_console()
	for cw.ReadKey() != "ENTER" {

	}
}
