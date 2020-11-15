package main

import "github.com/sidav/golibrl/geometry"

func (p *pawn) ai_timeoutState() {
	if p.ai.currentStateTimeoutTurn < CURRENT_TURN {
		if p.ai.currentState == AI_SEARCHING {
			// reset to one if calm states
			if p.ai.route != nil {
				p.ai.currentState = AI_PATROLLING
			} else {
				p.ai.currentState = AI_ROAM
			}
		} else if p.ai.currentState == AI_ALERTED {
			p.ai.currentState = AI_SEARCHING
		}
		p.ai.currentStateTimeoutTurn = CURRENT_TURN+25*10
	}
}

//func (p *pawn) ai_raiseAlarmLevel() {
//	newState := p.ai.currentState
//	switch p.ai.currentState {
//	case AI_ROAM, AI_PATROLLING:
//		newState = AI_SEARCHING
//	case AI_SEARCHING:
//		newState = AI_ALERTED
//	}
//	p.ai.currentState = newState
//}

func (p *pawn) ai_isCalm() bool {
	return p.ai.currentState == AI_PATROLLING || p.ai.currentState == AI_ROAM
}

func (p *pawn) ai_canSeePlayer() bool {
	x, y := p.getCoords()
	px, py := CURRENT_MAP.player.getCoords()
	if CURRENT_MAP.currentPlayerVisibilityMap[x][y] {
		usedSightRange := p.getStaticData().sightRangeCalm
		if p.ai_isCalm() {
			return CURRENT_MAP.player.isInLight() && geometry.AreCoordsInRange(px, py, x, y, usedSightRange)
		} else {
			if CURRENT_MAP.player.isInLight() {
				usedSightRange = p.getStaticData().sightRangeAlerted
			} else {
				usedSightRange = p.getStaticData().sightRangeAlertedDark
			}
		}
		return geometry.AreCoordsInRange(px, py, x, y, usedSightRange)
	}
	return false
}

// returns true if action is done
func (p *pawn) ai_TryMoveOrOpenDoorOrAlert(dirx, diry int) bool {
	ai := p.ai
	newx, newy := p.x+dirx, p.y+diry
	if CURRENT_MAP.isTilePassable(newx, newy) || CURRENT_MAP.isTileADoor(newx, newy){
		pawnAt := CURRENT_MAP.getPawnAt(newx, newy)
		if pawnAt == CURRENT_MAP.player {
			ai.targetPawn = pawnAt
			ai.currentState = AI_ALERTED
		}
		if pawnAt == nil {
			// close the door behind if needed
			if CURRENT_MAP.isTileADoor(p.x, p.y) && CURRENT_MAP.tiles[p.x][p.y].isOpened {
				CURRENT_MAP.tiles[p.x][p.y].isOpened = false
			}
			CURRENT_MAP.movePawnOrOpenDoorByVector(p, true, dirx, diry)
		}
		return true
	}
	return false
}
