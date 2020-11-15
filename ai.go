package main

import "parcelcreationtool/parcel"

type aiState uint8

const (
	AI_ROAM aiState = iota
	AI_PATROLLING
	AI_ALERTED
	AI_ENGAGING
)

type aiData struct {
	currentState aiState
	targetPawn *pawn
	dirx, diry int // for roaming

	route                *parcel.Route // for patrol
	currentWaypointIndex int
}

func (p *pawn) ai_checkSituation() {
	switch p.ai.currentState {
	case AI_ROAM, AI_PATROLLING:
		p.ai_checkRoam()
	case AI_ALERTED:
		p.ai_checkAlerted()
	default:
		log.AppendMessage("No CHECK func for some ai state!")
	}
}

func (p *pawn) ai_act() {
	switch p.ai.currentState {
	case AI_ROAM:
		p.ai_actRoam()
	case AI_PATROLLING:
		p.ai_actPatrolling()
	case AI_ALERTED:
		p.ai_actAlerted()
	default:
		log.AppendMessage("No ACT func for some ai state!")
	}
}

func (p *pawn) ai_resetStateToCalm() {
	if p.ai.route != nil {
		p.ai.currentState = AI_PATROLLING
	} else {
		p.ai.currentState = AI_ROAM
	}
}

func (p *pawn) ai_checkRoam() {
	x, y := p.getCoords()
	px, py := CURRENT_MAP.player.getCoords()
	if CURRENT_MAP.currentPlayerVisibilityMap[x][y] {
		if CURRENT_MAP.tiles[px][py].lightLevel > 0 {
			p.ai.targetPawn = CURRENT_MAP.player
			p.ai.currentState = AI_ALERTED
			return
		}
	}
	p.ai_resetStateToCalm()
}

func (p *pawn) ai_actRoam() {
	ai := p.ai
	for ai.dirx == 0 && ai.diry == 0 {
		ai.dirx, ai.diry = rnd.RandomUnitVectorInt()
	}
	if !p.ai_TryMoveOrOpenDoorOrAlert(ai.dirx, ai.diry) {
		ai.dirx = 0
		ai.diry = 0
	}
}

func (p *pawn) ai_actPatrolling() {
	ai := p.ai
	currWaypoint := ai.route.Waypoints[ai.currentWaypointIndex]
	px, py := p.getCoords()
	if px == currWaypoint.X && py == currWaypoint.Y {
		ai.currentWaypointIndex++
	}
	if ai.currentWaypointIndex >= len(ai.route.Waypoints) {
		ai.currentWaypointIndex = 0
	}
	currWaypoint = ai.route.Waypoints[ai.currentWaypointIndex]
	path := CURRENT_MAP.getPathFromTo(p.x, p.y, currWaypoint.X, currWaypoint.Y, false)
	dirx, diry := path.GetNextStepVector()
	p.ai_TryMoveOrOpenDoorOrAlert(dirx, diry)
}

func (p *pawn) ai_checkAlerted() {
	x, y := p.getCoords()
	px, py := CURRENT_MAP.player.getCoords()
	if CURRENT_MAP.currentPlayerVisibilityMap[x][y] {
		if CURRENT_MAP.tiles[px][py].lightLevel > 0 {
			p.ai.targetPawn = CURRENT_MAP.player
			p.ai.currentState = AI_ALERTED
			return
		}
	}
	p.ai_resetStateToCalm()
}

func (p *pawn) ai_actAlerted() {
	ai := p.ai
	for ai.dirx == 0 && ai.diry == 0 {
		ai.dirx, ai.diry = rnd.RandomUnitVectorInt()
	}
	newx, newy := p.x + ai.dirx, p.y + ai.diry
	if CURRENT_MAP.isTilePassable(newx, newy) {
		pawnAt := CURRENT_MAP.getPawnAt(newx, newy)
		if pawnAt == CURRENT_MAP.player {
			ai.targetPawn = pawnAt
			ai.currentState = AI_ALERTED
		}
		if pawnAt == nil {
			CURRENT_MAP.movePawnOrOpenDoorByVector(p, true, ai.dirx, ai.diry)
		}
	} else {
		ai.dirx = 0
		ai.diry = 0
	}
}

// returns true if action is done
func (p *pawn) ai_TryMoveOrOpenDoorOrAlert(dirx, diry int) bool {
	ai := p.ai
	newx, newy := p.x + dirx, p.y + diry
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
