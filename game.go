package main

import (
	"github.com/sidav/golibrl/console"
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

	GAME_IS_RUNNING = true
	CURRENT_MAP = gameMap{}
	CURRENT_MAP.initialize_level()
	CURRENT_MAP.MakeMapFromGenerated(&testMap)


	CURRENT_MAP.player = &pawn{
		ccell:         &consoleCell{
			appearance: '@',
			color:      console.WHITE,
			inverse:    false,
		},
		hp:            0,
		maxhp:         0,
		x:             1,
		y:             1,
		nextTurnToAct: 0,
		sightRange:    10,
		name:          "",
	}

	for GAME_IS_RUNNING {
		CURRENT_MAP.recalculateLights()
		CURRENT_MAP.currentPlayerVisibilityMap = *CURRENT_MAP.getFieldOfVisionFor(CURRENT_MAP.player)
		renderLevel(&CURRENT_MAP, true)
		pc.playerControl(&CURRENT_MAP)

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
		CURRENT_TURN++
	}
}
