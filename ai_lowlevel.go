package main

import "github.com/sidav/golibrl/geometry"

func (p *pawn) ai_resetStateToCalm() {
	if p.ai.route != nil {
		p.ai.currentState = AI_PATROLLING
	} else {
		p.ai.currentState = AI_ROAM
	}
}

func (p *pawn) ai_isCalm() bool {
	return p.ai.currentState == AI_PATROLLING || p.ai.currentState == AI_ROAM
}

func (p *pawn) ai_canSeePlayer() bool {
	x, y := p.getCoords()
	px, py := CURRENT_MAP.player.getCoords()
	if CURRENT_MAP.currentPlayerVisibilityMap[x][y] {
		if p.ai_isCalm() {
			if CURRENT_MAP.tiles[px][py].lightLevel > 0 && geometry.AreCoordsInRange(px, py, x, y, p.getStaticData().sightRangeCalm){
				return true
			}
		} else {
			if CURRENT_MAP.tiles[px][py].lightLevel > 0 &&
				geometry.AreCoordsInRange(px, py, x, y, p.getStaticData().sightRangeAlerted) ||
				geometry.AreCoordsInRange(px, py, x, y, p.getStaticData().sightRangeCalm){
				return true
			}
		}
	}
	return false
}

// returns true if action is done
func (p *pawn) ai_TryMoveOrOpenDoorOrAlert(dirx, diry int) bool {
	ai := p.ai
	newx, newy := p.x+dirx, p.y+diry
	if CURRENT_MAP.isTilePassable(newx, newy) {
		pawnAt := CURRENT_MAP.getPawnAt(newx, newy)
		if pawnAt == CURRENT_MAP.player {
			ai.targetPawn = pawnAt
			ai.currentState = AI_ALERTED
		}
		if pawnAt == nil {
			CURRENT_MAP.movePawnOrOpenDoorByVector(p, true, dirx, diry)
		}
		return true
	}
	return false
}
